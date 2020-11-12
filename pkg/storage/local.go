package storage

import (
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
	"io"
	"io/ioutil"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"strings"
)

type LocalReadWriter struct {
	OutDir   string
}

func NewLocalReadWriter(outdir string) (*LocalReadWriter, error) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		x := os.MkdirAll(outdir, 0777)
		if x != nil {
			return nil, x
		}
	} else {
		return nil, err
	}

	return &LocalReadWriter{
		OutDir: outdir,
	}, nil
}

func (local *LocalReadWriter) WriteFile(name string, data string)  error {
	flag := os.O_RDWR | os.O_TRUNC
	if _, err := os.Stat(name); os.IsNotExist(err) {
		flag |= os.O_CREATE
	}
	f, err := os.OpenFile(name, flag, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(common.StringToBytes(data))
	if err != nil {
		return err
	}
	if n != len(data) {
		return io.ErrShortWrite
	}
	return nil
}

func (local *LocalReadWriter) ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func (local *LocalReadWriter) LoadFiles(dir string) *Files {
	files := &Files{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			klog.Fatalf("loader.file.walk.error:%+v", err)
		}

		if !info.IsDir() {
			switch {
			case strings.HasSuffix(path, dbSuffix):
				files.Databases = append(files.Databases, path)
			case strings.HasSuffix(path, schemaSuffix):
				files.Schemas = append(files.Schemas, path)
			default:
				if strings.HasSuffix(path, tableSuffix) {
					files.Tables = append(files.Tables, path)
				}
			}
		}
		return nil
	}); err != nil {
		klog.Fatalf("loader.file.walk.error:%+v", err)
	}
	return files
}