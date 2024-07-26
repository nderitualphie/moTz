package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

type payload struct {
	Msisdn           string `json:"msisdn"`
	ShortCode        string `json:"shortcode"`
	Message          string `json:"message"`
	MessageTimestamp string `json:"messageTimestamp"`
}

func Process(c echo.Context, db *sql.DB) error {
	var p payload
	if err := c.Bind(&p); err != nil {
		return err
	}
	tableName := os.Getenv("TABLE_NAME")
	msisdn, _ := strconv.ParseInt(p.Msisdn, 10, 64)
	shortcode, _ := strconv.ParseInt(p.ShortCode, 10, 64)
	message := p.Message

	// Use placeholders in the SQL query to avoid SQL injection
	query := fmt.Sprintf("INSERT INTO %v (message, sender_address, dest_address,date_time, created_on) VALUES (?, ?, ?,NOW(),NOW())", tableName)

	_, err := db.Exec(query, message, msisdn, shortcode)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return c.String(http.StatusOK, "Payload processed successfully")
}

func main() {
	e := echo.New()

	// Set up your MySQL database connection here
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", username, password, host, port, dbName)
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	e.POST("/moTz", func(c echo.Context) error {
		return Process(c, db)
	})
	h := "0.0.0.0"
	p := "1235"

	hst := fmt.Sprintf("%s:%s", h, p)
	e.Logger.Fatal(e.Start(hst))
}
