package system

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	proxy "github.com/shogo82148/go-sql-proxy"
	"golang.org/x/crypto/ssh"
)

type DatabaseAccessPoint struct {
	Username string
	Password string
	Host     string
	Port     int
	Schema   string
}

func (p *DatabaseAccessPoint) ToDescriptiveString() string {
	return fmt.Sprintf("Host: %s\nPort: %d\nDatabase: %s\nUsername: %s\nPassword: %s", p.Host, p.Port, p.Schema, p.Username, p.Password)
}

func (p *DatabaseAccessPoint) toString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", p.Username, p.Password, p.Host, p.Port, p.Schema)
}

func (p *DatabaseAccessPoint) connect() (*sql.DB, error) {
	sql.Register("mysql-proxy", proxy.NewProxyContext(&mysql.MySQLDriver{}, &proxy.HooksContext{
		Open: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			log.Println("Open")
			return nil
		},
		Exec: func(_ context.Context, _ interface{}, stmt *proxy.Stmt, args []driver.NamedValue, result driver.Result) error {
			log.Printf("Exec: %s; args = %v\n", stmt.QueryString, args)
			return nil
		},
		Query: func(_ context.Context, _ interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows) error {
			log.Printf("Query: %s; args = %v\n", stmt.QueryString, args)
			return nil
		},
		Begin: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			log.Println("Begin")
			return nil
		},
		Commit: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			log.Println("Commit")
			return nil
		},
		Rollback: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			log.Println("Rollback")
			return nil
		},
	}))

	db, err := sql.Open("mysql", p.toString())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

type DatabaseConnection struct {
	connection    *sql.DB
	sshConnection *ssh.Client
}

func (dap *DatabaseAccessPoint) CreateDatabaseConnection() (*DatabaseConnection, error) {
	db, err := dap.connect()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &DatabaseConnection{connection: db}, nil
}

func (c *DatabaseConnection) Close() {
	c.connection.Close()
	if c.sshConnection != nil {
		c.sshConnection.Close()
	}
}

type DatabaseAccessPointOnSSH struct {
	DB  *DB
	SSH *SSH
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type SSH struct {
	Key        string
	Host       string
	Port       string
	User       string
	Passphrase string
}

func createSSHConnection(conf *SSH) (*ssh.Client, error) {
	sshKey, err := ioutil.ReadFile(conf.Key)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer
	if conf.Passphrase == "" {
		signer, err = ssh.ParsePrivateKey(sshKey)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(conf.Passphrase))
	}
	if err != nil {
		return nil, err
	}
	hostKeyCallbackFunc := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	sshConf := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallbackFunc,
	}
	return ssh.Dial("tcp", conf.Host+":"+conf.Port, sshConf)
}

func createDBConnection(conf *DB, sshc *ssh.Client) (*sql.DB, error) {
	mysqlNet := "tcp"
	if sshc != nil {
		mysqlNet = "mysql+tcp"
		dialFunc := func(addr string) (net.Conn, error) {
			return sshc.Dial("tcp", addr)
		}
		mysql.RegisterDial(mysqlNet, dialFunc)
	}
	dbConf := &mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Password,
		Addr:                 conf.Host + ":" + conf.Port,
		Net:                  mysqlNet,
		DBName:               conf.DBName,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", dbConf.FormatDSN())
}

func (dapOnSsh *DatabaseAccessPointOnSSH) CreateDatabaseConnectionOnSSH() (*DatabaseConnection, error) {
	sshConnection, err := createSSHConnection(dapOnSsh.SSH)
	if err != nil {
		return nil, err
	}

	dbConnection, err := createDBConnection(dapOnSsh.DB, sshConnection)
	if err != nil {
		return nil, err
	}
	return &DatabaseConnection{connection: dbConnection, sshConnection: sshConnection}, nil
}
