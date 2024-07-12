package storage

import (
	"database/sql"

	"github.com/Leandroreign/crud/pkg/invoice"
	"github.com/Leandroreign/crud/pkg/invoiceLines"
	"github.com/Leandroreign/crud/pkg/invoiceRelLines"
)

type PsqlInvoiceRelLines struct {
	db            *sql.DB
	storageHeader invoice.Storage
	storageLines  invoiceLines.Storage
}

func NewPsqlInvoiceRelLines(db *sql.DB, h invoice.Storage, l invoiceLines.Storage) *PsqlInvoiceRelLines {
	return &PsqlInvoiceRelLines{
		db:            db,
		storageHeader: h,
		storageLines:  l,
	}
}

func (p *PsqlInvoiceRelLines) Create(m *invoiceRelLines.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	// if err := p.storageHeader.CreateTx(tx, m.Invoice); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// if err := p.storageLines.CreateTx(tx, m.Invoice.Id, m.InvoiceLines); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	return tx.Commit()

}
