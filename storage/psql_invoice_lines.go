package storage

import (
	"database/sql"
	"fmt"

	"github.com/Leandroreign/crud/pkg/invoiceLines"
)

const (
	psqlMigrateInvoiceLines = `
		create table if not exists invoiceLines (
			id serial not null,
			invoice_id int not null,
			product_id int not null,
			createDate timestamp not null default now(),
			updateDate timestamp,
			constraint invoiceLines_id_pk primary key (id),
			constraint invoiceLines_invoice_id_fk foreign key 
			(invoice_id) references invoices (id) on update restrict
			on delete restrict,
			constraint invoiceLines_products_id_fk foreign key
			(product_id) references products (id) on update restrict
			on delete restrict
		)
	`

	psqlCreateInvoiceLines = `
		insert into invoicelines (invoice_id, product_id) values ($1, $2) returning id, createDate
	`
)

// PsqlProduct used for work with postgres - invoice
type PsqlInvoicesLines struct {
	db *sql.DB
}

func NewPsqlInvoicesLines(db *sql.DB) *PsqlInvoicesLines {
	return &PsqlInvoicesLines{db: db}
}

// Migrate implements the interface invoice storage
func (p *PsqlInvoicesLines) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoiceLines)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migracion de la tabla invoicesLines ejecutados correctamente")
	return nil
}

func (p *PsqlInvoicesLines) CreateTx(tx *sql.Tx, invoice_id uint, ms invoiceLines.Models) error {
	stmt, err := tx.Prepare(psqlCreateInvoiceLines)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, line := range ms {
		err := stmt.QueryRow(invoice_id, line.ProductId).Scan(&line.Id, &line.CreateDate)
		if err != nil {
			return err
		}
	}

	return nil
}
