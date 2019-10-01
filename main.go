package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Employee struct {
	Id   int
	Name string
	City string
}

type GrantObject struct {
	Id      int
	Name    string
	IsAdmin int
	NameUkr string
}

type GrantOperation struct {
	Id                int
	Name              string
	StandardOperation int
}

type grantMatrix struct {
	Id             int
	GrantObject    GrantObject
	GrantOperation GrantOperation
}

const (
	DB_USER     = "postgreadmin"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

func dbConn() (db *sql.DB) {
	dbDriver := "postgres"

	host := "localhost"
	port := 5432
	user := "postgreadmin"
	password := "postgres"
	dbName := "auth"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	db, err := sql.Open(dbDriver, psqlInfo)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Успішно підключилися до бази даних")
	}
	return db
}

func getAllGranObject(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	grantObjects := []GrantObject{}

	rows, err := db.Query(`SELECT id, name, is_admin as isAdmin, name_ukr FROM grant_object order by id`)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var id int
			var name string
			var isAdmin int
			var nameUkr string

			err = rows.Scan(&id, &name, &isAdmin, &nameUkr)
			if err == nil {
				currentGrantObject := GrantObject{Id: id, Name: name, IsAdmin: isAdmin, NameUkr: nameUkr}

				grantObjects = append(grantObjects, currentGrantObject)
			} else {
				panic(err.Error())
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grantObjects)
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func dbConnMySQL() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "test"
	dbPass := "Fkg7h4f3$"
	dbName := "go"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(api-o.aptekar.ua:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8081")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/api/grant/object/", getAllGranObject)
	http.ListenAndServe(":8081", nil)
}
