package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Database struct {
	db           *pgxpool.Pool
	connString   string
	errCount     []string
	lockErrors   sync.Mutex
	errorsPingDB *int32
}

func New(connString string) (*Database, error) {
	db, err := newPool(connString)
	if err != nil {
		log.Printf("error creating pgxpool %s", err)
		return nil, err
	}
	return &Database{db: db,
		connString:   connString,
		errorsPingDB: new(int32)}, nil
}

func newPool(connString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}
	poolConfig.MaxConns = 50
	poolConfig.MinConns = 10
	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (d *Database) addErrorCount(err error) {
	d.lockErrors.Lock()
	d.lockErrors.Unlock()
	d.errCount = append(d.errCount, err.Error())
}

func (d *Database) checkConnect() error {
	if err := d.db.Ping(context.Background()); err != nil {
		atomic.AddInt32(d.errorsPingDB, 1)
		return err
	}
	return nil
}

func (d *Database) Close() error {
	d.db.Close()
	return nil
}

func (d *Database) GetErrors() []string {
	return d.errCount
}

func (d *Database) GetErrorsCount() int {
	return len(d.errCount)
}

func (d *Database) GetErrorsPingDB() int {
	return int(atomic.LoadInt32(d.errorsPingDB))
}

func (d *Database) SimpleQueryReturnRandomUserName() (username string, err error) {
	defer func() {
		if err != nil {
			d.addErrorCount(err)
		}
	}()
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("recoverd from panic:%v", e)
		}
	}()

	if err = d.checkConnect(); err != nil {
		return "", err
	}

	var numberUsers int
	err = d.db.QueryRow(context.Background(), `select count(username) from users`).Scan(&numberUsers)
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
	err = d.db.QueryRow(context.Background(), query, userOffset).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
