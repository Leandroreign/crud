package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Leandroreign/crud/pkg/product"
)

const (
	mysqlMigrateProduct = `create table if not exists products(
		id int auto_increment not null primary key,
		name varchar(25) not null,
		observations varchar(100),
		price int not null,
		createDate timestamp not null default now(),
		updateDate timestamp
	)`
)

// PsqlProduct used for work with postgres - product
type MysqlProduct struct {
	db *sql.DB
}

func newMysqlProduct(db *sql.DB) *MysqlProduct {
	return &MysqlProduct{db: db}
}

// Migrate implements the interface product storage
func (p *MysqlProduct) Migrate() error {
	stmt, err := p.db.Prepare(mysqlMigrateProduct)
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

// Delete implement the interface product storafe
func (p *MysqlProduct) Delete(id uint) error {
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
func (p *MysqlProduct) Update(m *product.Model) error {
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
func (p *MysqlProduct) GetById(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(psqlGetProductById)
	if err != nil {
		return &product.Model{}, err
	}
	defer stmt.Close()

	return scanRowProduct(stmt.QueryRow(id))
}

// GetAll implements the interface product storage
func (p *MysqlProduct) GetAll() (product.Models, error) {
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

// Create implements the interface product storage
func (p *MysqlProduct) Create(m *product.Model) error {
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

// func scanRowProduct(s scanner) (*product.Model, error) {
// 	m := &product.Model{}
// 	observationNull := sql.NullString{}
// 	updateDateNull := sql.NullTime{}

// 	err := s.Scan(
// 		&m.Id, &m.Name, &observationNull, &m.Price, &m.CreateDate, &updateDateNull,
// 	)
// 	if err != nil {
// 		return &product.Model{}, err
// 	}
// 	m.Observations = observationNull.String
// 	m.UpdateDate = updateDateNull.Time

// 	return m, nil
// }
