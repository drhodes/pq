package pq

import (
	"database/sql"
	"log"
	"math/rand"
	"testing"
	"runtime"
)

const createperson = `
create table person (
        did SERIAL,
        unique(did),

        name varchar(255),
        email varchar(255),
        unique(email)

        );
`

func TestMaxProcs(t *testing.T) {
	runtime.GOMAXPROCS(4)

	db := openTestConn(t)
	defer db.Close()

	// Drop the table if it exists
	_, err := db.Exec(`drop table if exists person cascade;`)
	if err != nil {
		log.Fatal(err, "\ncan't drop table person")
	}

	// Create the table
	_, err = db.Exec(createperson)
	if err != nil {
		log.Fatal(err)
	}

	// Insert a bunch of records. 
	for i := 0; i < 10000; i++ {
		log.Print(i, " ")
		// This is where it will fail, if it does.
		err := maxProcsInsert(db, randomName(), randomEmail())
		if err != nil {
			log.Fatal(err, "\nProblem inserting person")
		}
	}
}

func randomName() string {
	alphas := `abcdefghijklmnopqrstuvwxyz`
	name := ``
	for i := 0; i < 10; i++ {
		name += string(alphas[rand.Int()%26])
	}
	return name
}

func randomEmail() string {
	return randomName() + "@" + randomName() + ".com"
}

func maxProcsInsert(db *sql.DB, name, email string) error {
	log.Println("Inserting: ", name, email)
	cols := `(name, email) values ($1,$2) `
	_, err := db.Exec(`INSERT INTO person `+cols, name, email)
	if err != nil {
		return err
	}
	return nil
}
