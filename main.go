package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Book  []Book
}

type Book struct {
	gorm.Model

	Title      string
	Arthor     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var (
	person = &Person{Name: "Abas Boazar", Email: "abasboazar@email.com"}
	books  = []Book{
		{Title: "Math", Arthor: "Farabi", CallNumber: 1000, PersonID: 1},
		{Title: "physic", Arthor: "Asemi", CallNumber: 1001, PersonID: 1},
	}
)

var db *gorm.DB
var err error

func main() {
	//loading enviroment variables

	dialet := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	//database connection string
	dbURL := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	//openinng connection to database
	db, err = gorm.Open(dialet, dbURL)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("succsesfuly connected to database")
	}

	// Closing connection to database when the main function is finshes
	defer db.Close()

	// Make migrations to the database if they have not already been created
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// API router
	router := mux.NewRouter()

	router.HandleFunc("/poeple", getPeople).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")

	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")

	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/create/book", createBook).Methods("POST")

	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	router.HandleFunc("/delete/book/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// API contorollers

// Books controllers

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.Find(&book, params["id"])

	json.NewEncoder(w).Encode(book)

}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}

// people controllers

func getPeople(w http.ResponseWriter, r *http.Request) {
	var poeple []Person

	db.Find(&poeple)

	json.NewEncoder(w).Encode(&poeple)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Model(&person).Related(&books)

	person.Book = books

	json.NewEncoder(w).Encode(person)

}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}
