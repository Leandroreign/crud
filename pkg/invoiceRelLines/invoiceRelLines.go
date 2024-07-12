package invoiceRelLines

import (
	"github.com/Leandroreign/crud/pkg/invoice"
	"github.com/Leandroreign/crud/pkg/invoiceLines"
)

type Model struct {
	Invoice      *invoice.Model
	InvoiceLines invoiceLines.Models
}

func NewModel(invoice *invoice.Model, invoiceLines invoiceLines.Models) *Model {
	return &Model{
		Invoice:      invoice,
		InvoiceLines: invoiceLines,
	}
}

// Storage interface that must implement a db storage
type Storage interface {
	Create(*Model) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
