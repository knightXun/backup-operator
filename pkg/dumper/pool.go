package dumper

import (
	"k8s.io/klog"
	"sync"

	"github.com/xelabs/go-mysqlstack/driver"

	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
)

// Pool tuple.
type Pool struct {
	mu    sync.RWMutex
	conns chan *Connection
	address string
	user string
	password string
	vars string
}

// Connection tuple.
type Connection struct {
	ID     int
	client driver.Conn
}

// Execute used to executes the query.
func (conn *Connection) Execute(query string) error {
	return conn.client.Exec(query)
}

// Fetch used to fetch the results.
func (conn *Connection) Fetch(query string) (*sqltypes.Result, error) {
	return conn.client.FetchAll(query, -1)
}

// StreamFetch used to the results with streaming.
func (conn *Connection) StreamFetch(query string) (driver.Rows, error) {
	return conn.client.Query(query)
}

// NewPool creates the new pool.
func NewPool(cap int, address string, user string, password string, vars string) (*Pool, error) {
	conns := make(chan *Connection, cap)
	for i := 0; i < cap; i++ {
		client, err := driver.NewConn(user, password, address, "", "utf8")
		if err != nil {
			return nil, err
		}
		conn := &Connection{ID: i, client: client}
		if vars != "" {
			klog.Info("session vars: ", vars)
			err = conn.Execute(vars)
			if err != nil {
				klog.Errorf("Set Session Failed: ", err)
			}
		}
		conns <- conn
	}

	return &Pool{
		conns: conns,
		address: address,
		user: user,
		password: password,
		vars: vars,
	}, nil
}

// Get used to get one connection from the pool.
func (p *Pool) Get() *Connection {
	conns := p.getConns()
	if conns == nil {
		return nil
	}

	conn := <-conns

	if conn.client.Closed() {
		client, err := driver.NewConn(p.user, p.password, p.address, "", "utf8")
		if err != nil {

		}
		conn = &Connection{ID: conn.ID, client: client}
		if p.vars != "" {
			err = conn.Execute(p.vars)
			if err != nil {
				klog.Errorf("Set Session Failed:  ", err)
			}
		}
	}

	return conn
}

// Put used to put one connection to the pool.
func (p *Pool) Put(conn *Connection) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.conns == nil {
		return
	}
	p.conns <- conn
}

// Close used to close the pool and the connections.
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)
	for conn := range p.conns {
		klog.Info("Closing Connection")
		conn.client.Close()
	}
	p.conns = nil
}

func (p *Pool) getConns() chan *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}
