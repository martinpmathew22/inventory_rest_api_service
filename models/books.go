package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Writer string  `json:"writer"`
	Price  float32 `json:"price"`
}

// BooksByWriter queries for Books that have the specified Writer name.
func BooksByWriter(db *sql.DB, name string) ([]Book, error) {
	// An Books slice to hold data from returned rows.
	var Books []Book

	rows, err := db.Query("SELECT * FROM book WHERE writer = ?", name)
	if err != nil {
		return nil, fmt.Errorf("BooksByWriter %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Book
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Writer, &alb.Price); err != nil {
			return nil, fmt.Errorf("BooksByWriter %q: %v", name, err)
		}
		Books = append(Books, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("BooksByWriter %q: %v", name, err)
	}
	return Books, nil
}

// BookByID queries for the Book with the specified ID.
func GetBooks(db *sql.DB) ([]Book, error) {
	// An Book to hold data from the returned row.
	var albs []Book

	results, err := db.Query("SELECT * FROM book")
	fmt.Printf("results %v", results)
	if err != nil {
		return nil, fmt.Errorf("Err", err.Error())
	}

	for results.Next() {
		var alb Book
		// for each row, scan into the Product struct
		err = results.Scan(&alb.ID, &alb.Title, &alb.Writer, &alb.Price)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// append the product into products array
		albs = append(albs, alb)
	}

	return albs, nil
}

// BookByID queries for the Book with the specified ID.
func BookByID(db *sql.DB, id int64) (Book, error) {
	// An Book to hold data from the returned row.
	var alb Book

	row := db.QueryRow("SELECT * FROM book WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Writer, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("BooksById %d: no such Book", id)
		}
		return alb, fmt.Errorf("BooksById %d: %v", id, err)
	}
	return alb, nil
}

// addBook adds the specified Book to the database,
// returning the Book ID of the new entry
func AddBook(db *sql.DB, alb Book) (int64, error) {
	result, err := db.Exec("INSERT INTO book (title, writer, price) VALUES (?, ?, ?)", alb.Title, alb.Writer, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addBook: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addBook: %v", err)
	}
	return id, nil
}
