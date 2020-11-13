package dumper

import (
	"fmt"
	"github.com/backup-operator/pkg/storage"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"
	"k8s.io/klog"
)

func writeMetaData(args *Args, writer storage.StorageReadWriter) error {
	klog.Info("Write Meta Data")
	file := fmt.Sprintf("%s/metadata", args.Outdir)
	err := writer.WriteFile(file, "")

	return err
}

func dumpDatabaseSchema(conn *Connection, args *Args, database string, writer storage.StorageReadWriter) error {
	err := conn.Execute(fmt.Sprintf("USE `%s`", database))
	if err != nil {
		klog.Errorf("Dump %s Tables Schema Failed: %v", database, err)
	}

	schema := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", database)
	file := fmt.Sprintf("%s/%s-schema-create.sql", args.Outdir, database)
	err = writer.WriteFile(file, schema)

	return err
}

func dumpTableSchema(conn *Connection, args *Args, database string, table string, writer storage.StorageReadWriter) error {
	qr, err := conn.Fetch(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", database, table))
	if err != nil {
		klog.Errorf("Dump %s Tables %s Schema Failed: %v", database, table, err)
		return err
	}

	schema := qr.Rows[0][1].String() + ";\n"

	file := fmt.Sprintf("%s/%s.%s-schema.sql", args.Outdir, database, table)
	err = writer.WriteFile(file, schema)

	return err
}

func dumpTable(conn *Connection, args *Args, database string, table string, writer storage.StorageReadWriter) error {
	var allBytes uint64
	var allRows uint64
	var where string
	var selfields []string

	fields := make([]string, 0, 16)
	{
		cursor, err := conn.StreamFetch(fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT 1", database, table))
		if err != nil {
			klog.Errorf("DumpTable Failed: %v", err)
			return err
		}

		flds := cursor.Fields()
		for _, fld := range flds {
			klog.Info("dump -- %#v, %s, %s", args.Filters, table, fld.Name)
			if _, ok := args.Filters[table][fld.Name]; ok {
				continue
			}

			fields = append(fields, fmt.Sprintf("`%s`", fld.Name))
			replacement, ok := args.Selects[table][fld.Name]
			if ok {
				selfields = append(selfields, fmt.Sprintf("%s AS `%s`", replacement, fld.Name))
			} else {
				selfields = append(selfields, fmt.Sprintf("`%s`", fld.Name))
			}
		}
		err = cursor.Close()
		if err != nil {
			klog.Errorf("DumpTable: Close Cursor Failed: %v", err)
			return err
		}
	}

	if v, ok := args.Wheres[table]; ok {
		where = fmt.Sprintf(" WHERE %v", v)
	}

	cursor, err := conn.StreamFetch(fmt.Sprintf("SELECT %s FROM `%s`.`%s` %s", strings.Join(selfields, ", "), database, table, where))
	if err != nil {
		return err
	}

	fileNo := 1
	stmtsize := 0
	chunkbytes := 0
	rows := make([]string, 0, 256)
	inserts := make([]string, 0, 256)
	for cursor.Next() {
		row, err := cursor.RowValues()
		if err != nil {
			klog.Errorf("DumpTable: Fetching Data Failed: %v", err)
			return err
		}

		values := make([]string, 0, 16)
		for _, v := range row {
			if v.Raw() == nil {
				values = append(values, "NULL")
			} else {
				str := v.String()
				switch {
				case v.IsSigned(), v.IsUnsigned(), v.IsFloat(), v.IsIntegral(), v.Type() == querypb.Type_DECIMAL:
					values = append(values, str)
				default:
					values = append(values, fmt.Sprintf("\"%s\"", EscapeBytes(v.Raw())))
				}
			}
		}
		r := "(" + strings.Join(values, ",") + ")"
		rows = append(rows, r)

		allRows++
		stmtsize += len(r)
		chunkbytes += len(r)
		allBytes += uint64(len(r))
		atomic.AddUint64(&args.Allbytes, uint64(len(r)))
		atomic.AddUint64(&args.Allrows, 1)

		if stmtsize >= args.StmtSize {
			insertone := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES\n%s", table, strings.Join(fields, ","), strings.Join(rows, ",\n"))
			inserts = append(inserts, insertone)
			rows = rows[:0]
			stmtsize = 0
		}

		if (chunkbytes / 1024 / 1024) >= args.ChunksizeInMB {
			query := strings.Join(inserts, ";\n") + ";\n"
			file := fmt.Sprintf("%s/%s.%s.%05d.sql", args.Outdir, database, table, fileNo)
			err := writer.WriteFile(file, query)
			if err != nil {
				klog.Errorf("DumpTable: Writing Data Failed: %v", err)
				return err
			}

			klog.Info("dumping.table[%s.%s].rows[%v].bytes[%vMB].part[%v].thread[%d]", database, table, allRows, (allBytes / 1024 / 1024), fileNo, conn.ID)
			inserts = inserts[:0]
			chunkbytes = 0
			fileNo++
		}
	}
	if chunkbytes > 0 {
		if len(rows) > 0 {
			insertone := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES\n%s", table, strings.Join(fields, ","), strings.Join(rows, ",\n"))
			inserts = append(inserts, insertone)
		}

		query := strings.Join(inserts, ";\n") + ";\n"
		file := fmt.Sprintf("%s/%s.%s.%05d.sql", args.Outdir, database, table, fileNo)
		err := writer.WriteFile(file, query)
		if err != nil {
			klog.Errorf("DumpTable: Writing Data Failed: %v", err)
			return err
		}
	}
	err = cursor.Close()
	if err != nil {
		klog.Errorf("DumpTable Failed: %v", err)
		return err
	}

	klog.Info("dumping.table[%s.%s].done.allrows[%v].allbytes[%vMB].thread[%d]...", database, table, allRows, (allBytes / 1024 / 1024), conn.ID)
	return nil
}

func allTables(conn *Connection, database string) []string {
	qr, err := conn.Fetch(fmt.Sprintf("SHOW TABLES FROM `%s`", database))
	if err != nil {
		klog.Fatalf("Fetch Tables Failed: %v", err)
	}

	tables := make([]string, 0, 128)
	for _, t := range qr.Rows {
		tables = append(tables, t[0].String())
	}
	return tables
}

func allDatabases(conn *Connection) []string {
	qr, err := conn.Fetch("SHOW DATABASES")
	if err != nil {
		klog.Fatalf("Fetch Databases Failed: %v", err)
	}

	databases := make([]string, 0, 128)
	for _, t := range qr.Rows {
		databases = append(databases, t[0].String())
	}
	return databases
}

func filterDatabases(conn *Connection, filter *regexp.Regexp, invert bool) []string {
	qr, err := conn.Fetch("SHOW DATABASES")
	if err != nil {
		klog.Fatalf("Filter Databases Failed: %v", err)
	}

	databases := make([]string, 0, 128)
	for _, t := range qr.Rows {
		if (!invert && filter.MatchString(t[0].String())) || (invert && !filter.MatchString(t[0].String())) {
			databases = append(databases, t[0].String())
		}
	}
	return databases
}

// Dumper used to start the dumper worker.
func Dumper(args *Args) {
	pool, err := NewPool(args.Threads, args.Address, args.User, args.Password, args.SessionVars)
	if err != nil {
		klog.Fatalf("Make Mysql Pool Failed: %v", err)
	}
	defer pool.Close()

	var writer storage.StorageReadWriter
	if args.Local != nil {
		writer , err = storage.NewLocalReadWriter(args.Local.Outdir)
		if err != nil {
			klog.Fatalf("Create LocalReadWriter Failed: %v", err)
		}
	} else if args.S3 != nil {
		writer, err = storage.NewS3ReadWriter(args.S3.S3Endpoint,
			args.S3.S3Region,
			args.S3.S3Bucket,
			args.S3.S3AccessKey,
			args.S3.S3SecretAccessKey,
			args.S3.BackupDir)

		if err != nil {
			klog.Fatalf("Create S3ReadWriter Failed: %v", err)
		}
	} else {
		klog.Fatalf("Invalid Writer Config")
	}

	// Meta data.
	err = writeMetaData(args, writer)

	if err != nil {
		klog.Fatalf("Write Meta Data Failed: %v", err)
	}
	// database.
	var wg sync.WaitGroup
	conn := pool.Get()
	var databases []string
	t := time.Now()
	if args.DatabaseRegexp != "" {
		r := regexp.MustCompile(args.DatabaseRegexp)
		databases = filterDatabases(conn, r, args.DatabaseInvertRegexp)
	} else {
		if args.Database != "" {
			databases = strings.Split(args.Database, ",")
		} else {
			databases = allDatabases(conn)
		}
	}
	for _, database := range databases {
		err = dumpDatabaseSchema(conn, args, database, writer)
		if err != nil {
			klog.Fatalf("Dump Database Schema Failed: %v", err)
		}
	}

	// tables.
	tables := make([][]string, len(databases))
	for i, database := range databases {
		if args.Table != "" {
			tables[i] = strings.Split(args.Table, ",")
		} else {
			tables[i] = allTables(conn, database)
		}
	}
	pool.Put(conn)

	for i, database := range databases {
		for _, table := range tables[i] {
			conn := pool.Get()
			err = dumpTableSchema(conn, args, database, table, writer)
			if err != nil {
				klog.Fatalf("Dump Table Schema Failed: %v", err)
			}

			wg.Add(1)
			go func(conn *Connection, database string, table string) {
				defer func() {
					wg.Done()
					pool.Put(conn)
				}()
				klog.Info("dumping.table[%s.%s].datas.thread[%d]...", database, table, conn.ID)
				err = dumpTable(conn, args, database, table, writer)
				if err != nil {
					klog.Fatalf("Dump Database Schema Failed: %v", err)
				}
				klog.Info("dumping.table[%s.%s].datas.thread[%d].done...", database, table, conn.ID)
			}(conn, database, table)
		}
	}

	tick := time.NewTicker(time.Millisecond * time.Duration(args.IntervalMs))
	defer tick.Stop()
	go func() {
		for range tick.C {
			diff := time.Since(t).Seconds()
			allbytesMB := float64(atomic.LoadUint64(&args.Allbytes) / 1024 / 1024)
			allrows := atomic.LoadUint64(&args.Allrows)
			rates := allbytesMB / diff
			klog.Info("dumping.allbytes[%vMB].allrows[%v].time[%.2fsec].rates[%.2fMB/sec]...", allbytesMB, allrows, diff, rates)
		}
	}()

	wg.Wait()
	elapsed := time.Since(t).Seconds()
	klog.Info("dumping.all.done.cost[%.2fsec].allrows[%v].allbytes[%v].rate[%.2fMB/s]", elapsed, args.Allrows, args.Allbytes, (float64(args.Allbytes/1024/1024) / elapsed))
}