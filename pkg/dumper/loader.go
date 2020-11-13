package dumper

import (
	"fmt"
	"github.com/backup-operator/pkg/storage"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"k8s.io/klog"

	"github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
)

// Files tuple.
type Files struct {
	databases []string
	schemas   []string
	tables    []string
}

var (
	dbSuffix     = "-schema-create.sql"
	schemaSuffix = "-schema.sql"
	tableSuffix  = ".sql"
)


func restoreDatabaseSchema(dbs []string, conn *Connection, reader storage.StorageReadWriter) error {
	for _, db := range dbs {
		base := filepath.Base(db)
		name := strings.TrimSuffix(base, dbSuffix)

		data, err := reader.ReadFile(db)
		if err != nil {
			klog.Errorf("ReadFile Failed: %v", err)
			return err
		}
		sql := common.BytesToString(data)

		err = conn.Execute(sql)
		if err != nil {
			klog.Errorf("Execute SQL Failed: %v", err)
			return err
		}

		klog.Info("restoring.database[%s]", name)
	}

	return nil
}

func restoreTableSchema(overwrite bool, tables []string, conn *Connection,reader storage.StorageReadWriter) error {
	for _, table := range tables {
		// use
		base := filepath.Base(table)
		name := strings.TrimSuffix(base, schemaSuffix)
		db := strings.Split(name, ".")[0]
		tbl := strings.Split(name, ".")[1]
		name = fmt.Sprintf("`%v`.`%v`", db, tbl)

		klog.Info("working.table[%s.%s]", db, tbl)

		err := conn.Execute(fmt.Sprintf("USE `%s`", db))
		if err != nil {
			klog.Errorf("Execute SQL Failed: %v", err)
			return err
		}

		err = conn.Execute("SET FOREIGN_KEY_CHECKS=0")
		if err != nil {
			klog.Errorf("Execute SQL Failed: %v", err)
			return err
		}

		data, err := reader.ReadFile(table)
		if err != nil {
			klog.Errorf("Read SQL File Failed: %v", err)
			return err
		}

		query1 := common.BytesToString(data)
		querys := strings.Split(query1, ";\n")
		for _, query := range querys {
			if !strings.HasPrefix(query, "/*") && query != "" {
				if overwrite {
					klog.Info("drop(overwrite.is.true).table[%s.%s]", db, tbl)
					dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
					err = conn.Execute(dropQuery)
					if err != nil {
						klog.Errorf("Execute SQL Failed: %v", err)
						return err
					}
				}
				err = conn.Execute(query)
				if err != nil {
					klog.Errorf("Execute SQL Failed: %v", err)
					return err
				}
			}
		}
		klog.Info("restoring.schema[%s.%s]", db, tbl)
	}

	return nil
}

func restoreTable(table string, conn *Connection, reader storage.StorageReadWriter) (int, error) {
	bytes := 0
	part := "0"
	base := filepath.Base(table)
	name := strings.TrimSuffix(base, tableSuffix)
	splits := strings.Split(name, ".")
	db := splits[0]
	tbl := splits[1]
	if len(splits) > 2 {
		part = splits[2]
	}

	klog.Info("restoring.tables[%s.%s].parts[%s].thread[%d]", db, tbl, part, conn.ID)
	err := conn.Execute(fmt.Sprintf("USE `%s`", db))
	if err != nil {
		klog.Errorf("Execute SQL Failed: %v", err)
		return 0, err
	}

	err = conn.Execute("SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		klog.Errorf("Execute SQL Failed: %v", err)
		return 0, err
	}


	data, err := reader.ReadFile(table)
	if err != nil {
		klog.Errorf("Read File Failed: %v", err)
		return 0, err
	}

	query1 := common.BytesToString(data)
	querys := strings.Split(query1, ";\n")
	bytes = len(query1)
	for _, query := range querys {
		if !strings.HasPrefix(query, "/*") && query != "" {
			err = conn.Execute(query)
			if err != nil {
				klog.Errorf("Execute Query Failed: %v", err)
				return 0, err
			}
		}
	}
	klog.Info("restoring.tables[%s.%s].parts[%s].thread[%d].done...", db, tbl, part, conn.ID)
	return bytes, nil
}

// Loader used to start the loader worker.
func Loader(args *Args) {
	pool, err := NewPool(args.Threads, args.Address, args.User, args.Password, args.SessionVars)
	if err != nil {
		klog.Fatalf("Make Mysql Connection Pool Failed: %v", err)
	}

	defer pool.Close()

	var reader storage.StorageReadWriter
	if args.Local != nil {
		reader, err  = storage.NewLocalReadWriter(args.Local.Outdir)
		if err != nil {
			klog.Fatalf("Create LocalReadWriter Failed: %v", err)
		}
	} else if args.S3 != nil {
		reader, err = storage.NewS3ReadWriter(args.S3.S3Endpoint,
			args.S3.S3Region,
			args.S3.S3Bucket,
			args.S3.S3AccessKey,
			args.S3.S3SecretAccessKey,
			args.S3.BackupDir)

		if err != nil {
			klog.Fatalf("Create S3ReadWriter Failed: %v", err)
		}
	} else {
		log.Fatalf("Invalid Reader")
	}

	files := reader.LoadFiles(args.Outdir)

	if files == nil {
		log.Fatalf("Cannot fetch any files")
	}
	// database.
	conn := pool.Get()
	err = restoreDatabaseSchema(files.Databases, conn, reader)
	if err != nil {
		klog.Fatalf("Restore Database Schema Failed: %v", err)
	}
	pool.Put(conn)

	// tables.
	conn = pool.Get()
	err = restoreTableSchema(args.OverwriteTables, files.Schemas, conn, reader)
	if err != nil {
		klog.Fatalf("Restore Table Schema Failed: %v", err)
	}
	pool.Put(conn)

	// Shuffle the tables
	for i := range files.Tables {
		j := rand.Intn(i + 1)
		files.Tables[i], files.Tables[j] = files.Tables[j], files.Tables[i]
	}

	var wg sync.WaitGroup
	var bytes uint64
	t := time.Now()
	for _, table := range files.Tables {
		conn := pool.Get()
		wg.Add(1)
		go func(conn *Connection, table string) {
			defer func() {
				wg.Done()
				pool.Put(conn)
			}()
			r , err := restoreTable(table, conn, reader)
			if err != nil {
				klog.Fatalf("Restore Table Failed: %v", err)
			}
			atomic.AddUint64(&bytes, uint64(r))
		}(conn, table)
	}

	tick := time.NewTicker(time.Millisecond * time.Duration(args.IntervalMs))
	defer tick.Stop()
	go func() {
		for range tick.C {
			diff := time.Since(t).Seconds()
			bytes := float64(atomic.LoadUint64(&bytes) / 1024 / 1024)
			rates := bytes / diff
			klog.Info("restoring.allbytes[%vMB].time[%.2fsec].rates[%.2fMB/sec]...", bytes, diff, rates)
		}
	}()

	wg.Wait()
	elapsed := time.Since(t).Seconds()
	klog.Info("restoring.all.done.cost[%.2fsec].allbytes[%.2fMB].rate[%.2fMB/s]", elapsed, float64(bytes/1024/1024), (float64(bytes/1024/1024) / elapsed))
}
