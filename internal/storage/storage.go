package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const TimeReconnect = time.Second * 2

type Database struct {
	db         *sql.DB
	connString string
	errCount   []string
	lock       sync.Mutex
}

func New(conn string) (*Database, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		e := db.Close()
		return nil, errors.Join(err, e)
	}
	return &Database{db: db, connString: conn}, nil
}

func (d *Database) addErrorCount(err error) {
	d.lock.Lock()
	d.lock.Unlock()
	d.errCount = append(d.errCount, err.Error())
}

func (d *Database) CheckConnect() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (d *Database) ReConnect() error {
	_ = d.db.Close()
	time.Sleep(TimeReconnect)
	db, err := sql.Open("postgres", d.connString)
	if err != nil {
		d.addErrorCount(err)
		return err
	}

	err = d.CheckConnect()
	if err != nil {
		e := db.Close()
		d.addErrorCount(err)
		return errors.Join(err, e)
	}
	d.db = db
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) SimpleQuery() (username string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("recoverd from panic:%v", e)
		}
	}()

	if err = d.CheckConnect(); err != nil {
		return "", err
	}

	var numberUsers int
	err = d.db.QueryRow(`select count(username) from users`).Scan(&numberUsers)
	if err != nil {
		return "", err
	}
	if numberUsers == 0 {
		return "", errors.New("0 users")
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	userOffset := r1.Intn(numberUsers)

	query := `select username from users OFFSET $1 limit 1`
	err = d.db.QueryRow(query, userOffset).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
