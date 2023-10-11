package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const limitPerTransaction = 10000

func main() {
	var users_number int
	var messages_number int
	flag.IntVar(&users_number, "users", 1000, "number of users to add")
	flag.IntVar(&messages_number, "messages", 100000, "number of messages to add")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	connURL := os.Getenv("CONNURL")
	if connURL == "" {
		log.Printf("env CONNURL is not set")
		return
	}

	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	//var greeting string
	//err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	//	os.Exit(1)
	//}
	_, err = conn.Exec(context.Background(), `TRUNCATE TABLE messages CASCADE`)
	if err != nil {
		log.Printf("got error truncation messages table", err)
		return
	}
	log.Println("Successfully truncated messages table")

	_, err = conn.Exec(context.Background(), `TRUNCATE TABLE users CASCADE`)
	if err != nil {
		log.Printf("got error truncation users table", err)
		return
	}
	log.Println("Successfully truncated messages table")

	if err := addUsers(conn, users_number); err != nil {
		log.Printf("got error adding users:%s", err)
		return
	}
	log.Printf("Added %d users", users_number)

	if err := addMessages(conn, messages_number); err != nil {
		log.Printf("got error adding messages:%s", err)
		return
	}
	log.Printf("Added %d messages", messages_number)

}

func addUsers(conn *pgx.Conn, numberusers int) (err error) {
	tx, err := conn.Begin(context.Background())
	defer func() {
		e := tx.Rollback(context.Background())
		if err != nil {
			err = errors.Join(e, err)
		}
	}()
	for i := 0; i < numberusers; i++ {
		_, err := conn.Exec(context.Background(), `INSERT INTO users (username) VALUES ($1)`, faker.FirstName())
		if err != nil {
			log.Printf("got error addin user", err)
			return fmt.Errorf("got error addin user: %w", err)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func addMessages(conn *pgx.Conn, numberMessages int) (err error) {
	var numberUsers int
	err = conn.QueryRow(context.Background(), `select count(user_id) from users`).Scan(&numberUsers)
	if err != nil {
		return err
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	users := make([]int, 0, limitPerTransaction)
	friends := make([]int, 0, limitPerTransaction)
	for i := 0; i < numberMessages; i++ {
		user := r1.Intn(numberUsers)
		friend := r1.Intn(numberUsers)
		var user_id int
		var friend_id int

		query := `select user_id from users OFFSET $1 limit 1`
		err = conn.QueryRow(context.Background(), query, user).Scan(&user_id)
		if err != nil {
			return err
		}
		err = conn.QueryRow(context.Background(), query, friend).Scan(&friend_id)
		if err != nil {
			return err
		}
		users = append(users, user_id)
		friends = append(friends, friend_id)
		if (i+1)%limitPerTransaction == 0 {
			err = addMessagesTX(conn, users, friends)
			log.Printf("Added %d messages", i+1)

			if err != nil {
				return err
			}
			users = make([]int, 0, limitPerTransaction)
			friends = make([]int, 0, limitPerTransaction)
		}
	}

	if len(users) != 0 {
		err = addMessagesTX(conn, users, friends)
		if err != nil {
			return err
		}
	}
	return nil
}

func addMessagesTX(conn *pgx.Conn, users []int, friends []int) error {
	tx, err := conn.Begin(context.Background())
	defer func() {
		e := tx.Rollback(context.Background())
		if err != nil {
			err = errors.Join(e, err)
		}
	}()
	for i := 0; i < len(users); i++ {
		if users[i] == friends[i] {
			continue
		}
		_, err := conn.Exec(context.Background(), `INSERT INTO messages (user_id, friend_id, message) VALUES ($1,$2,$3)`,
			users[i], friends[i],
			fmt.Sprintf("%s %s %s %s %s", faker.Word(), faker.Word(), faker.Word(), faker.Word(), faker.Word()))
		if err != nil {
			log.Printf("user_id=%d, friend_id=%d", users[i], friends[i])
			log.Printf("got error adding message", err)
			return fmt.Errorf("got error adding message: %w", err)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
