package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
)

type UserPSQL struct {
	db *sql.DB
}

func NewUserPSQL(db *sql.DB) *UserPSQL {
	return &UserPSQL{
		db: db,
	}
}

func (r *UserPSQL) Create(u *entity.User) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO users (id, username, email, pswd, userrole) 
		values($1, $2, $3, sha256($4), $5)`)
	if err != nil {
		return u.ID, err
	}
	_, err = stmt.Exec(
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.UserRole,
	)
	if err != nil {
		return u.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return u.ID, err
	}
	return u.ID, nil
}

func (r *UserPSQL) GetbyUsername(username string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, pswd, userrole from users where username = $1`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(username)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.UserRole)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPSQL) GetbyID(ID string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, userrole from users where ID = $1`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(ID)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.UserRole)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPSQL) Update(u *entity.User) error {
	_, err := r.db.Exec("UPDATE users SET pswd = sha256($1),  username = $2, email = $3, userrole = $4 where id = $5",
		u.Password, u.Username, u.Email, u.UserRole, u.Username)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPSQL) Search(query string) ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, userrole FROM users WHERE username like $1`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.UserRole)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserPSQL) List() ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT ID, username, email, userrole FROM users`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID,
			&u.Username, &u.Email, &u.UserRole)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserPSQL) Delete(username string) error {
	_, err := r.db.Exec("DELETE FROM users where username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserPSQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
