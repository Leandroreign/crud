package product

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrIdNotSended = errors.New("el producto no contiene un id")
)

// Model of products
type Model struct {
	Id           uint
	Name         string
	Observations string
	Price        float64
	CreateDate   time.Time
	UpdateDate   time.Time
}

func (m *Model) String() string {
	return fmt.Sprintf("%02d | %20s | %20s | %5f | %10s | %10s",
		m.Id, m.Name, m.Observations, m.Price, m.CreateDate.Format("2006-01-02"),
		m.UpdateDate.Format("2006-01-02"))
}

func NewProduct(name, observations string, price float64) *Model {
	return &Model{
		Name:         name,
		Observations: observations,
		Price:        price,
	}
}

// Models slice de Model de products
type Models []*Model

type Storage interface {
	Migrate() error
	Create(*Model) error
	Update(*Model) error
	GetAll() (Models, error)
	GetById(uint) (*Model, error)
	Delete(uint) error
}

// Service of product
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

func (s *Service) Delete(id uint) error {
	if id == 0 {
		return ErrIdNotSended
	}
	return s.storage.Delete(id)
}

func (s *Service) Update(m *Model) error {
	if m.Id == 0 {
		return ErrIdNotSended
	}
	m.UpdateDate = time.Now()
	return s.storage.Update(m)
}

func (s *Service) GetById(id uint) (*Model, error) {
	return s.storage.GetById(id)
}

func (s *Service) GetAll() (Models, error) {
	return s.storage.GetAll()
}

func (s *Service) Create(p *Model) error {
	p.CreateDate = time.Now()
	return s.storage.Create(p)
}
