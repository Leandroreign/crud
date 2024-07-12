package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Leandroreign/crud/pkg/product"
)

type scanner interface {
	Scan(dest ...any) error
}

const (
	psqlMigrateProduct = `create table if not exists products(
		id serial not null,
		name varchar(25) not null,
		observations varchar(100),
		price int not null,
		createDate timestamp not null default now(),
		updateDate timestamp,
		constraint products_id_pk primary key (id)
	)`

	psqlInsertProduct = `
		insert into products(
			name, observations, price, createDate
		)
			values(
				$1, $2, $3, $4
			) returning id
	`

	psqlGetProducts = `
		select * from products
	`

	psqlGetProductById = psqlGetProducts + " where id = $1"

	psqlUpdateProduct = `
		update products set 
			name = $1,
			observations = $2,
			price = $3,
			updateDate = $4
		where id = $5  
	`

	psqlDeleteProduct = `
		delete from products where id = $1
	`
)

// PsqlProduct used for work with postgres - product
type PsqlProduct struct {
	db *sql.DB
}

func newPsqlProduct(db *sql.DB) *PsqlProduct {
	return &PsqlProduct{db: db}
}

// Delete implement the interface product storafe
func (p *PsqlProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(psqlDeleteProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("no existe el id: %d", id)
	}
	fmt.Printf("El producto con id %v fue eliminado correctamente", id)
	return nil
}

// Update implements the interface product storage
func (p *PsqlProduct) Update(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlUpdateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(m.Name, stringToNull(m.Observations), m.Price, timeToNull(time.Now()), m.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("no existe el id: %d", m.Id)
	}

	fmt.Printf("El producto %v fue actualizado correctamente", m.Name)
	return nil

}

// GerById implements the interface product storafe
func (p *PsqlProduct) GetById(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(psqlGetProductById)
	if err != nil {
		return &product.Model{}, err
	}
	defer stmt.Close()

	return scanRowProduct(stmt.QueryRow(id))
}

// GetAll implements the interface product storage
func (p *PsqlProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(psqlGetProducts)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ms := make(product.Models, 0)
	for rows.Next() {
		m, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil

}

// Migrate implements the interface product storage
func (p *PsqlProduct) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migracion de la tabla productos ejecutados correctamente")
	return nil
}

// Create implements the interface product storage
func (p *PsqlProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlInsertProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(m.Name, stringToNull(m.Observations), m.Price, m.CreateDate).Scan(&m.Id)
	if err != nil {
		return err
	}
	fmt.Printf("Producto %v creado correctamente\n", m.Name)
	return nil
}

func scanRowProduct(s scanner) (*product.Model, error) {
	m := &product.Model{}
	observationNull := sql.NullString{}
	updateDateNull := sql.NullTime{}

	err := s.Scan(
		&m.Id, &m.Name, &observationNull, &m.Price, &m.CreateDate, &updateDateNull,
	)
	if err != nil {
		return &product.Model{}, err
	}
	m.Observations = observationNull.String
	m.UpdateDate = updateDateNull.Time

	return m, nil
}
