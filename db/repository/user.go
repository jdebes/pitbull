package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"firstname"`
	Surname   string `db:"surname"`
	UserName  string `db:"username"`
	Password  []byte `db:"password"`
}

func InsertUser(sqlDB *sqlx.DB, user *User) (sql.Result, error) {
	return sq.Insert("user").
		Columns("firstname", "surname", "username", "password").
		Values(user.FirstName, user.Surname, user.UserName, user.Password).
		RunWith(sqlDB).
		Exec()
}

func UserByUserName(sqlDB *sqlx.DB, userName string, user *User) error {
	selectSq, params, err := sq.
		Select("*").
		From("user").
		Where(sq.Eq{"username": userName}).
		ToSql()
	if err != nil {
		return err
	}

	return sqlx.Get(sqlDB, user, selectSq, params...)
}
