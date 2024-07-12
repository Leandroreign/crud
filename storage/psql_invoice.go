package storage

import (
	"database/sql"
	"fmt"

	"github.com/Leandroreign/crud/pkg/invoice"
)

const (
	psqlMigrateInvoice = `
		create table if not exists Invoices (
			id serial not null,
			client varchar(100) not null,
			createDate timestamp not null default now(),
			updateDate timestamp,
			constraint invoice_id_pk primary key (id)
		)
	`

	psqlCreateInvoice = `
		insert into invoices (client) values ($1) returning id, createDate
	`
)

// PsqlProduct used for work with postgres - invoice
type PsqlInvoices struct {
	db *sql.DB
}

func NewPsqlInvoice(db *sql.DB) *PsqlInvoices {
	return &PsqlInvoices{db: db}
}

// Migrate implements the interface invoice storage
func (p *PsqlInvoices) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoice)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migracion de la tabla invoices ejecutados correctamente")
	return nil
}

func (p *PsqlInvoices) CreateTx(tx *sql.Tx, m *invoice.Model) error {
	stmt, err := tx.Prepare(psqlCreateInvoice)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(m.Client).Scan(&m.Id, &m.CreateDate)
}
