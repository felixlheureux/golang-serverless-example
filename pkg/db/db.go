package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"net"
	"net/url"
	"reflect"
	"time"
)

const (
	DefaultLimit = 50
	DefaultPage  = 1
)

func WithTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()

	if err != nil {
		return err
	}

	err = fn(tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// JoinSuffix returns the `ON` part of a `JOIN` statement for two table that shares the same key name
//
// example:
// JoinSuffix("companies", "projects", "company_id") => "companies ON companies.company_id = projects.company_id"
func JoinSuffix(table1, table2, column string) string {
	return table1 + " ON " + table1 + "." + column + " = " + table2 + "." + column
}

func QueryExecuteError(err error, query string, args []interface{}) error {
	if len(query) > 200 {
		query = query[:200] + "..."
	}

	return fmt.Errorf("failed to execute query: %w\n%s\n%s", err, query, args)
}

func UnmarshalError(err error, query string, args []interface{}) error {
	return fmt.Errorf("failed to unmarshal struct: %w\n%s\n%s", err, query, args)
}

// PostgresURL returns a connection URL for postgres
func PostgresURL(host, port, user, pass, database string) string {
	// postgres wants connection url to be percentage quoted
	// https://www.postgresql.org/docs/11/libpq-connect.html#id-1.7.3.8.3.6
	pass = url.QueryEscape(pass)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, database)
}

// Postgres creates a database connection with pgx and sqlx
func Postgres(host, port, user, pass, database string, clients ...*ssh.Client) (*sqlx.DB, error) {
	u := PostgresURL(host, port, user, pass, database)

	config, err := pgx.ParseConfig(u)
	if err != nil {
		return nil, fmt.Errorf("invalid url %s: %w", u, err)
	}

	// assumes the global logger is already configured correctly
	config.Logger = zapadapter.NewLogger(zap.L())

	// dialer is used to transmit data via SSH
	if len(clients) > 0 {
		client := clients[0]

		dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return client.Dial(network, addr)
		}

		config.DialFunc = dialer
	}

	// open connection
	db := stdlib.OpenDB(*config)
	db.SetMaxIdleConns(0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// make sure connection is successful
	err = db.PingContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("unable to ping db %s: %w", u, err)
	}

	return sqlx.NewDb(db, "pgx"), nil
}

func WithRetry(period time.Duration, limit int, fn func() (*sqlx.DB, error)) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < limit; i++ {
		if db, err = fn(); err == nil {
			return db, nil
		}
		<-time.After(period)
	}

	return nil, fmt.Errorf("retry failed: %w", err)
}

func SSH(sshUser, sshHost, sshPort string, sshKey []byte) (*ssh.Client, error) {
	key, err := ssh.ParsePrivateKey(sshKey)

	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	var auth ssh.AuthMethod
	auth = ssh.PublicKeys(key)

	cfg := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), cfg)

	if err != nil {
		return nil, fmt.Errorf("unable to open ssh connection: %w", err)
	}

	return client, nil
}

// GetDBColumns returns all the field defined by `db` tags of a struct
func GetDBColumns(v interface{}) []string {
	t := reflect.TypeOf(v)
	var columns []string

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if column, ok := f.Tag.Lookup("db"); ok {
			columns = append(columns, column)
		}
	}

	return columns
}

// GetOffsetAndLimit returns offset and limit for a query, if params are invalid return default values
func GetOffsetAndLimit(page int, limit int) (int, int) {
	if limit < 1 || limit > DefaultLimit {
		limit = DefaultLimit
	}

	if page < 1 {
		page = DefaultPage
	}

	// calculate offset
	offset := page*limit - limit

	return offset, limit
}

func PrefixColumns(table string, columns []string) []string {
	var prefixedColumns []string

	for _, c := range columns {
		prefixedColumns = append(prefixedColumns, table+"."+c+" AS "+c)
	}

	return prefixedColumns
}
