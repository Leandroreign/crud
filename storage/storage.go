package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Leandroreign/crud/pkg/product"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// Driver of storage
type Driver string

// Drivers
const (
	MySQL    Driver = "MySQL"
	Postgres Driver = "Postgres"
)

// New create the connection with database
func New(d Driver) {
	switch d {
	case MySQL:
		newMySQLDB()
	case Postgres:
		newPostgresDB()
	}
}

func newPostgresDB() {
	once.Do(
		func() {
			var err error
			db, err = sql.Open("postgres", "postgres://openpg:openpgpwd@localhost:5432/godb?sslmode=disable")
			if err != nil {
				log.Fatalf("can't open database: %v", err)
			}
			// no vamos a cerrar la conexion porque es unica
			// defer db.Close()

			if err = db.Ping(); err != nil {
				log.Fatalf("can't ping database: %v", err)
			}
			fmt.Println("database postgresql connected")
		},
	)
}

func newMySQLDB() {
	once.Do(
		func() {
			var err error
			db, err = sql.Open("mysql", "root:leandro@tcp(localhost:3306)/godb")
			if err != nil {
				log.Fatalf("can't open database: %v", err)
			}
			// no vamos a cerrar la conexion porque es unica
			// defer db.Close()

			if err = db.Ping(); err != nil {
				log.Fatalf("can't ping database: %v", err)
			}
			fmt.Println("database mysql connected")
		},
	)
}

// Pool return a unique instance od db
func Pool() *sql.DB {
	return db
}

func stringToNull(s string) sql.NullString {
	null := sql.NullString{String: s}
	if null.String != "" {
		null.Valid = true
	}
	return null
}

func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}

// DAOProduct factory of product.storage
func DAOProduct(driver Driver) (product.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlProduct(db), nil
	case MySQL:
		return newMysqlProduct(db), nil
	default:
		return nil, fmt.Errorf("Driver no implementado")
	}

}
