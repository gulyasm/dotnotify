package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sq "github.com/lann/squirrel"
)

type File struct {
	Name, User            string
	Hash                  int64
	CreatedAt, ModifiedAt time.Time
}

type RealDb struct {
	Database *sql.DB
}

func (db *RealDb) Init(params map[string]string) error {
	user := params["user"]
	password := params["password"]
	host := params["host"]
	dbname := params["database"]
	cs := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, host, dbname)
	log.Println(cs)
	var err error
	db.Database, err = sql.Open("mysql", cs)
	return err
}

func (db *RealDb) InsertFile(f File) error {
	qb := sq.Insert("Files").Columns("User", "Name", "Hash", "ModifiedAt", "CreatedAt")
	qb = qb.Values(f.User, f.Name, f.Hash, time.Now(), time.Now())
	_, err := qb.RunWith(db.Database).Exec()
	if err != nil {
		log.Println(err)
		err = errors.New(fmt.Sprintf("Failed to insert file %s", f))
		return err
	}
	return nil
}

func (db *RealDb) UpdateFile(f File) error {
	qb := sq.Update("Files")
	qb = qb.Set("Hash", f.Hash)
	qb = qb.Set("ModifiedAt", f.ModifiedAt)
	qb = qb.Where(sq.Eq{"User": f.User, "Name": f.Name})
	_, err := qb.RunWith(db.Database).Exec()
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to update file %s", f))
		log.Println(err)
		return err
	}
	return nil
}

func (db *RealDb) GetFile(user, fn string) ([]File, error) {
	qb := sq.Select("*").From("Files")
	qb = qb.Where(sq.Eq{"User": user})
	qb = qb.Where(sq.Eq{"Name": fn})
	rows, err := qb.RunWith(db.Database).Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	files := []File{}
	for rows.Next() {
		f := File{}
		err := rows.Scan(
			&f.Name,
			&f.User,
			&f.Hash,
			&f.CreatedAt,
			&f.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
