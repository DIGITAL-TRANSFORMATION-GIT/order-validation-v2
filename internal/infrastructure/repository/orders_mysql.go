package repository

import (
	"database/sql"

	"clean/internal/entity"
)

type OrdersMySQL struct {
	db *sql.DB
}

//NewBookMySQL create new repository
func NewOrdersMySQL(db *sql.DB) *OrdersMySQL {
	return &OrdersMySQL{
		db: db,
	}
}

//Create a book
func (r *OrdersMySQL) Create(e *entity.Orders) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO orders (id, title, description, deadline) 
		values(?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Title,
		e.Description,
		e.Deadline,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

//Get a book
func (r *OrdersMySQL) Get(id entity.ID) (*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders where id = ?`)
	if err != nil {
		return nil, err
	}
	var b entity.Orders
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&b.ID, &b.Title, &b.Description, &b.Deadline)
		if err != nil {
			return nil, err
		}
	}
	return &b, nil
}

//Update a book
func (r *OrdersMySQL) Update(e *entity.Orders) error {
	_, err := r.db.Exec("UPDATE orders SET title = ?, description = ?, deadline = ? where id = ?",
		e.Title, e.Description, e.Deadline, e.ID)
	if err != nil {
		return err
	}
	return nil
}

//Search books
func (r *OrdersMySQL) Search(query string) ([]*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders WHERE title like ?`)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Orders
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o entity.Orders
		err = rows.Scan(&o.ID, &o.Title, &o.Description, &o.Deadline)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

//List books
func (r *OrdersMySQL) List() ([]*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Orders
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o entity.Orders
		err = rows.Scan(&o.ID, &o.Title, &o.Description, &o.Deadline)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

//Delete a book
func (r *OrdersMySQL) Delete(id entity.ID) error {
	_, err := r.db.Exec("DELETE FROM orders where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersMySQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
