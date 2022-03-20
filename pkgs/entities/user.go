package entities

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId   int    `json:"user_id"`  // user_id
	Username string `json:"username"` // username
	// password
}

// find an existing user
// the password will be checked
// returns pgx.ErrNoRows if the user isn't found
func FindUser(
	ctx context.Context,
	dbConn *pgx.Conn,
	username string,
	password string,
) (*User, error) {
	// get sql string
	sql, args, err := dialect.
		From("users").
		Select("user_id", "username", "password").
		Where(goqu.Ex{"username": username}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	// query db
	user := User{}
	passwordHash := ""
	err = dbConn.
		QueryRow(ctx, sql, args).
		Scan(&user.UserId, &user.Username, &passwordHash)
	if err != nil {
		return nil, err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// creates a new user
// the password will be hashed
func CreateUser(
	ctx context.Context,
	dbConn *pgx.Conn,
	username string,
	password string,
) (*User, error) {
	// I don't really care how complex the password is right now
	// hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	sql, _, err := dialect.
		Insert("users").
		Rows(goqu.Record{"username": username, "password": passwordHash}).
		Returning("user_id", "username").
		ToSQL()
	if err != nil {
		return nil, err
	}

	user := User{}
	err = dbConn.QueryRow(ctx, sql).Scan(&user.UserId, &user.Username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
