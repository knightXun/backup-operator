package storage

var (
	dbSuffix     = "-schema-create.sql"
	schemaSuffix = "-schema.sql"
	tableSuffix  = ".sql"
)
// Files tuple.
type Files struct {
	Databases []string
	Schemas   []string
	Tables    []string
}

type StorageReadWriter interface {
	WriteFile(name string, data string) error
	ReadFile(name string) ([]byte, error)
	LoadFiles(dir string) *Files
}