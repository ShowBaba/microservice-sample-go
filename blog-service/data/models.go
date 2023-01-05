package data

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/showbaba/microservice-sample-go/shared"
)

type Post struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	UserID    int64       `json:"userId"`
	User      shared.User `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type User shared.User

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

func (u *Post) Insert(post *Post) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var id int
	query := `INSERT INTO Posts (title, body, userId, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := db.QueryRowContext(ctx, query,
		&post.Title, &post.Body,
		&post.UserID,
		time.Now(), time.Now()).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (u *Post) Update() (*Post, error) {
	return nil, nil
}

func (u *Post) Delete() error {
	return nil
}

func (u *Post) DeleteByID(id int) error {
	return nil
}

func (u *Post) GetAll() ([]*Post, error) {
	return nil, nil
}

func (u *Post) GetByID(id int) (*Post, error) {
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
		tableExist, err := shared.CheckTableExist(ctx, db, "posts")
		if err != nil {
			errorCh <- err
		}
		if !tableExist {
			query := `CREATE TABLE POSTS(
					ID SERIAL PRIMARY KEY, TITLE VARCHAR(500) NOT NULL,
					BODY TEXT NOT NULL, USERID INT NOT NULL, CREATED_AT DATE, UPDATED_AT DATE
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
