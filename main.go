package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"example.com/inventory_rest_api_service/models"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	{ID: "4", Title: "Martin", Artist: "Martin P Mathew", Price: 99.99},
}

func main() {
	log.SetOutput(os.Stdout)

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Addr:   os.Getenv("DB_HOST"),
		Net:    "tcp",
		DBName: os.Getenv("DB_NAME"),
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	_, err = db.Exec("CREATE TABLE book (id         INT AUTO_INCREMENT NOT NULL, title      VARCHAR(128) NOT NULL, writer     VARCHAR(255) NOT NULL, price      DECIMAL(5,2) NOT NULL, PRIMARY KEY (`id`));")
	if err != nil {
		panic(err)
	}

	// // Create
	// res, err = db.Exec("INSERT INTO mytable (some_text) VALUES (?)", "hello world")
	// if err != nil {
	//     panic(err)
	// }
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/books", getBooks)
	router.GET("/book/:code", getBook)
	router.POST("/book", addBook)

	router.Run("0.0.0.0:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getBooks(c *gin.Context) {
	books, err := models.GetBooks(db)

	fmt.Printf("err %v", err)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, books)
	}
}

func getBook(c *gin.Context) {
	code := c.Param("code")

	id, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	fmt.Printf("Hello, %v with type %s!\n", id, reflect.TypeOf(id))

	book, err := models.BookByID(db, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, book)
	}
}

func addBook(c *gin.Context) {
	var prod models.Book

	if err := c.BindJSON(&prod); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		id, err := models.AddBook(db, prod)
		prod.ID = id
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.IndentedJSON(http.StatusCreated, prod)
	}
}
