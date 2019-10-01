package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	DB_DRIVER = "postgres"

	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgreadmin"
	PASSWORD = "postgres"
	DBNAME   = "auth"
)

type GrantObject struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	NameUkr        string `json:"nameUkr"`
	StandardObject int    `json:"standardObject"`
	IsAdmin        int    `json:"isAdmin"`
}

type GrantOperation struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	StandardOperation int    `json:"standardOperation"`
}

type GrantMatrix struct {
	Id             int
	GrantObject    GrantObject
	GrantOperation GrantOperation
}

type Role struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Standard_role   int       `json:"standardRole"`
	Create_time     time.Time `json:"create"`
	Expiration_date time.Time `json:"expirationDate"`
	Status          int       `json:"update"`
	IsAdmin         int       `json:"isAdmin"`
	UpdateTime      time.Time `json:"whoChanged"`
	WhoChanged      int       `json:"id"`
}

func dbInit() (db *sql.DB, error error) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open(DB_DRIVER, dbinfo)
	checkErr(err)

	fmt.Println("# Успішно підключилися до бази!")
	return db, err
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func GetAllGrantObject(w http.ResponseWriter, r *http.Request) {
	db, err := dbInit()

	grantObjects := []GrantObject{}

	rows, err := db.Query(`SELECT id, name, standard_object as standardObject, name_ukr as nameUkr, is_admin as isAdmin FROM grant_object order by id`)

	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var id int
			var name string
			var standardObject int
			var nameUkr string
			var isAdmin int

			err = rows.Scan(&id, &name, &standardObject, &nameUkr, &isAdmin)
			if err == nil {
				currentGrantObject := GrantObject{Id: id, Name: name, StandardObject: standardObject, NameUkr: nameUkr, IsAdmin: isAdmin}
				checkErr(err)
				grantObjects = append(grantObjects, currentGrantObject)
			} else {
				panic(err.Error())
			}
		}
		response, err := json.Marshal(grantObjects)
		checkErr(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func GetAllGrantOperation(w http.ResponseWriter, r *http.Request) {
	db, err := dbInit()

	grantOperation := []GrantOperation{}

	rows, err := db.Query(`SELECT id, name, standard_operation as standardOperation FROM grant_operation order by id`)

	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var id int
			var name string
			var standardOperation int

			err = rows.Scan(&id, &name, &standardOperation)
			if err == nil {
				currentGrantObject := GrantOperation{Id: id, Name: name, StandardOperation: standardOperation}
				checkErr(err)
				grantOperation = append(grantOperation, currentGrantObject)
			} else {
				panic(err.Error())
			}
		}
		response, err := json.Marshal(grantOperation)
		checkErr(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func GetAllGrantMatrix(w http.ResponseWriter, r *http.Request) {
	db, err := dbInit()

	grantMatrix := []GrantMatrix{}

	rows, err := db.Query(`SELECT id,  grant_operation_id as operation, grant_object_id as object   FROM grant_matrix order by id`)

	defer rows.Close()
	if err == nil {
		for rows.Next() {

			var id int
			var operation int
			var object int
			Object := GrantObject{}
			Operation := GrantOperation{}

			err = rows.Scan(&id, &operation, &object)
			if err == nil {

				err := db.QueryRow(`SELECT id, name, standard_object as standardObject, name_ukr as nameUkr, is_admin as isAdmin FROM grant_object where id = $1 limit 1`, object).Scan(&Object.Id, &Object.Name, &Object.StandardObject, &Object.NameUkr, &Object.IsAdmin)

				if err != nil {
					panic(err.Error())
				}

				err = db.QueryRow(`SELECT id, name, standard_operation as standardOperation FROM grant_operation where id = $1 limit 1`, operation).Scan(&Operation.Id, &Operation.Name, &Operation.StandardOperation)

				if err != nil {
					panic(err.Error())
				}
				currentGrantMatrix := GrantMatrix{Id: id, GrantObject: Object, GrantOperation: Operation}
				checkErr(err)
				grantMatrix = append(grantMatrix, currentGrantMatrix)
			} else {
				panic(err.Error())
			}
		}

		response, err := json.Marshal(grantMatrix)
		checkErr(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func GetAllGrantMatrixJoin(w http.ResponseWriter, r *http.Request) {
	db, err := dbInit()

	grantMatrix := []GrantMatrix{}

	rows, err := db.Query(`SELECT m.id  , ob.id, ob.name, ob.name_ukr,ob.standard_object, ob.is_admin, op.id, op.name, op.standard_operation   FROM grant_matrix as m
join grant_object ob on m.grant_object_id = ob.id
join grant_operation op on m.grant_operation_id = op.id
order by  m.id `)

	defer rows.Close()
	if err == nil {
		for rows.Next() {
			Matrix := GrantMatrix{}
			//GrantMatrix{Id: id, GrantObject: Object, GrantOperation: Operation}
			Object := GrantObject{}
			Operation := GrantOperation{}

			err = rows.Scan(&Matrix.Id, &Object.Id, &Object.Name, &Object.NameUkr, &Object.StandardObject, &Object.IsAdmin, &Operation.Id, &Operation.Name, &Operation.StandardOperation)
			if err == nil {
				Matrix.GrantObject = Object
				Matrix.GrantOperation = Operation
				checkErr(err)
				grantMatrix = append(grantMatrix, Matrix)
			} else {
				panic(err.Error())
			}
		}

		response, err := json.Marshal(grantMatrix)
		checkErr(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func GetAllGrantMatrixJoin2(w http.ResponseWriter, r *http.Request) {
	db, err := dbInit()

	grantMatrix := []GrantMatrix{}

	rows, err := db.Query(`SELECT m.id  , ob.id, ob.name, ob.name_ukr,ob.standard_object, ob.is_admin, op.id, op.name, op.standard_operation   FROM grant_matrix as m
join grant_object ob on m.grant_object_id = ob.id
join grant_operation op on m.grant_operation_id = op.id
order by  m.id `)

	defer rows.Close()
	if err == nil {
		for rows.Next() {
			Matrix := GrantMatrix{Id: 0, GrantObject: GrantObject{}, GrantOperation: GrantOperation{}}

			err = rows.Scan(&Matrix.Id, &Matrix.GrantObject.Id, &Matrix.GrantObject.Name, &Matrix.GrantObject.NameUkr, &Matrix.GrantObject.StandardObject, &Matrix.GrantObject.IsAdmin, &Matrix.GrantOperation.Id, &Matrix.GrantOperation.Name, &Matrix.GrantOperation.StandardOperation)
			if err == nil {
				checkErr(err)
				grantMatrix = append(grantMatrix, Matrix)
			} else {
				panic(err.Error())
			}
		}

		response, err := json.Marshal(grantMatrix)
		checkErr(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	} else {
		panic(err.Error())
	}

	panic(err.Error())
}

func main() {

	log.Println("Server started on: http://localhost:8081")
	http.HandleFunc("/api/grant/object/", GetAllGrantObject)
	http.HandleFunc("/api/grant/operation/", GetAllGrantOperation)
	http.HandleFunc("/api/grant/matrix/", GetAllGrantMatrix)
	http.HandleFunc("/api/grant/matrix/join/", GetAllGrantMatrixJoin)
	http.HandleFunc("/api/grant/matrix/join/2/", GetAllGrantMatrixJoin2)
	http.ListenAndServe(":8081", nil)
}
