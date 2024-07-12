package invoice

import (
	"database/sql"
	"time"
)

// Model of invoice
type Model struct {
	Id         uint
	Client     string
	CreateDate time.Time
	UpdateDate time.Time
}

func NewModel(client string) *Model {
	return &Model{
		Client: client,
	}
}

type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, *Model) error
}

// Service of invoice
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

func (s *Service) CreateTx(tx *sql.Tx, m *Model) error {
	return s.storage.CreateTx(tx, m)
}
