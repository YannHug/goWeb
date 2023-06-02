package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var schema = `
CREATE TABLE IF NOT EXISTS book
(
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 title TEXT,
 author TEXT,
 page_count INTEGER
)
`

type Book struct {
	Id        int
	Title     string
	Author    string
	PageCount int
}

func main() {
	db, err := sql.Open("sqlite3", "goweb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(schema)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	// map to values
	var id int
	var title string
	var author string
	var pageCount int

	rows, _ := db.Query("SELECT * FROM book")

	for rows.Next() {
		rows.Scan(&id, &title, &author, &pageCount)
		fmt.Printf("id=%v, title=%v, author=%v, pagecount=%v \n",
			id, title, author, pageCount)
	}

	// map to struct
	rows, _ = db.Query("SELECT * FROM book")
	b := Book{}
	for rows.Next() {
		rows.Scan(&b.Id, &b.Title, &b.Author, &b.PageCount)
		fmt.Printf("book=%v\n", b)
	}

	fmt.Println("Ping to database successful")

	// Insert into (prepared statements

	b = Book{
		Title:     "Harry Potter",
		Author:    "JK ROWLING",
		PageCount: 308,
	}

	stmt, _ = db.Prepare("INSERT INTO book (title, author, page_count) VALUES (?,?,?)")
	_, err = stmt.Exec(b.Title, b.Author, b.PageCount)
	if err != nil {
		log.Fatal(err)
	}
}
