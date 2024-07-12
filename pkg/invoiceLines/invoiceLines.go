package invoiceLines

import (
	"database/sql"
	"time"
)

type Model struct {
	Id         uint
	InvoiceId  uint
	ProductId  uint
	CreateDate time.Time
	UpdateDate time.Time
}

func NewModel(productId uint) *Model {
	return &Model{
		ProductId: productId,
	}
}

type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, uint, Models) error
}

type Models []*Model

// Service of invoiceLines
type Service struct {
	storage Storage
}

// NewService return a pointer of service
func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Migrate() error {
	return s.storage.Migrate()
}

func (s *Service) CreateTx(tx *sql.Tx, invoice_id uint, ms Models) error {
	return s.storage.CreateTx(tx, invoice_id, ms)
}
