package data

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/showbaba/microservice-sample-go/shared"
)

type User shared.User

func (u *User) Insert(user *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var id int
	query := `INSERT INTO Users (email, firstname, lastname, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	if err := db.QueryRowContext(ctx, query,
		&user.Email, &user.Firstname,
		&user.Lastname, &user.Password,
		time.Now(), time.Now()).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) Update() (*User, error) {
	return nil, nil
}

func (u *User) Delete() error {
	return nil
}

func (u *User) DeleteByID(id int) error {
	return nil
}

func (u *User) GetAll() ([]*User, error) {
	return nil, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT id, email, firstname, lastname, password, created_at, updated_at FROM Users WHERE email = $1`
	var user User
	row := db.QueryRowContext(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func (u *User) GetByID(id int) (*User, error) {
	return nil, nil
}

func Migrate() {
	const TOTAL_WORKERS = 1
	var (
		wg      sync.WaitGroup
		errorCh = make(chan error, TOTAL_WORKERS)
	)
	wg.Add(TOTAL_WORKERS)
	log.Println("running db migration")

	go func() {
		defer wg.Done()
		// Users
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		tableExist, err := shared.CheckTableExist(ctx, db, "users")
		if err != nil {
			errorCh <- err
		}
		if !tableExist {
			query := `CREATE TABLE USERS(
					ID SERIAL PRIMARY KEY, EMAIL CHAR(500) UNIQUE NOT NULL,
					FIRSTNAME CHAR(500) NOT NULL, LASTNAME CHAR(500) NOT NULL,
					PASSWORD CHAR(500) NOT NULL, CREATED_AT DATE, UPDATED_AT DATE
				)`
			_, err := db.ExecContext(ctx, query)
			if err != nil {
				errorCh <- err
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	for err := range errorCh {
		if err != nil {
			panic(err)
		}
	}

	log.Println("complete db migration")
}
