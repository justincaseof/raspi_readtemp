package database

import (
	 "database/sql"
	 "fmt"

	 _ "github.com/lib/pq"
)

const (
	host     = "192.168.171.34"
	port     = 5432
	user     = "pitemp"
	password = "pitemp"
	dbname   = "pitemp"
)

func InitDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	ensureTableExists("foobar123")
}

/**
 tableIdentifier should be the raspi's mac address
 */
func ensureTableExists(tableIdentifier string) {

}