package repo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/myntdeveloper/agromir_boiler_server/internal/config"
)

func NewDB() *sql.DB {
	conf := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database unavailable:", err)
	}
	log.Println("Connected to database!")
	return db
}
