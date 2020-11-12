package dumper

import (
	"fmt"
	"strings"

	ini "github.com/dlintw/goconf"
)

func ParseDumperConfig(file string) (*Args, error) {
	args := &Args{
		Wheres: make(map[string]string, 0),
	}

	cfg, err := ini.ReadConfigFile(file)
	if err != nil {
		return nil, err
	}

	host, err := cfg.GetString("mysql", "host")
	if err != nil {
		return nil, err
	}
	port, err := cfg.GetInt("mysql", "port")
	if err != nil {
		return nil, err
	}
	user, err := cfg.GetString("mysql", "user")
	if err != nil {
		return nil, err
	}
	password, err := cfg.GetString("mysql", "password")
	if err != nil {
		return nil, err
	}
	database, _ := cfg.GetString("mysql", "database")
	sessionVars, err := cfg.GetString("mysql", "vars")
	if err != nil {
		return nil, err
	}
	chunksizemb, err := cfg.GetInt("mysql", "chunksize")
	if err != nil {
		return nil, err
	}
	table, _ := cfg.GetString("mysql", "table")

	// Options
	if err := loadOptions(cfg, "where", args.Wheres); err != nil {
		return nil, err
	}

	var selects []string
	if selects, err = cfg.GetOptions("select"); err != nil {
		return nil, err
	}
	for _, tblcol := range selects {
		var table, column string
		split := strings.Split(tblcol, ".")
		table = split[0]
		column = split[1]

		if args.Selects == nil {
			args.Selects = make(map[string]map[string]string)
		}
		if args.Selects[table] == nil {
			args.Selects[table] = make(map[string]string, 0)
		}
		if args.Selects[table][column], err = cfg.GetString("select", tblcol); err != nil {
			return nil, err
		}
	}

	database_regexp, _ := cfg.GetString("database", "regexp")
	database_invert_regexp, err := cfg.GetBool("database", "invert_regexp")

	if cfg.HasSection("s3") {
		args.S3 = &S3Backend{}
		endpoint, err := cfg.GetString("s3", "endpoint")
		if err != nil {
			return nil, err
		}
		args.S3.S3Endpoint = endpoint

		region, err := cfg.GetString("s3", "region")
		if err != nil {
			return nil, err
		}
		args.S3.S3Region = region

		accessKey, err := cfg.GetString("s3", "accessKey")
		if err != nil {
			return nil, err
		}
		args.S3.S3AccessKey = accessKey

		secretAccessKey, err := cfg.GetString("s3", "secretAccessKey")
		if err != nil {
			return nil, err
		}
		args.S3.S3SecretAccessKey = secretAccessKey

		bucket, err := cfg.GetString("s3", "bucket")
		if err != nil {
			return nil, err
		}
		args.S3.S3Bucket = bucket

		backupdir, err := cfg.GetString("s3", "backupdir")
		if err != nil {
			return nil, err
		}
		args.S3.BackupDir = backupdir

		args.Outdir = backupdir
	} else if cfg.HasSection("local") {
		args.Local = &LocalBackend{}
		outdir, err := cfg.GetString("local", "outdir")
		if err != nil {
			return nil, err
		}
		args.Local.Outdir = outdir

		args.Outdir = outdir
	} else {
		return nil, fmt.Errorf("Unknown Storage Backend")
	}

	if err != nil {
		database_invert_regexp = false
	}

	var filters []string
	if filters, err = cfg.GetOptions("filter"); err != nil {
		return nil, err
	}
	for _, tblcol := range filters {
		var table, column string
		split := strings.Split(tblcol, ".")
		table = split[0]
		column = split[1]

		if args.Filters == nil {
			args.Filters = make(map[string]map[string]string)
		}
		if args.Filters[table] == nil {
			args.Filters[table] = make(map[string]string, 0)
		}
		if args.Filters[table][column], err = cfg.GetString("filter", tblcol); err != nil {
			return nil, err
		}
	}

	args.Address = fmt.Sprintf("%s:%d", host, port)
	args.User = user
	args.Password = password
	args.Database = database
	args.DatabaseRegexp = database_regexp
	args.DatabaseInvertRegexp = database_invert_regexp
	args.Table = table
	args.ChunksizeInMB = chunksizemb
	args.SessionVars = sessionVars
	args.Threads = 16
	args.StmtSize = 1000000
	args.IntervalMs = 10 * 1000
	return args, nil
}

func loadOptions(cfg *ini.ConfigFile, section string, optMap map[string]string) error {
	var err error
	var opts []string

	if opts, err = cfg.GetOptions(section); err != nil {
		return err
	}

	for _, key := range opts {
		if optMap[key], err = cfg.GetString(section, key); err != nil {
			return err
		}
	}
	return nil
}