package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

/* steps:
create the routes through the http package
create the function for the routes
in those functions get the values from the databases
*/
type User struct {
	UserId    int
	FirstName string
	LastName  string
	Email     string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":8000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	psqlConnection := fmt.Sprintf("host=localhost port=5432 user=vivek password=%v dbname=vivek sslmode=disable", os.Getenv("PG_PASS"))
	db, err := sql.Open("postgres", psqlConnection)
	checkError(err)
	defer db.Close()
	rows, err := db.Query("select * from users")
	var (
		id         int
		first_name string
		last_name  string
		email      string
		user       User
		users      []User
	)
	for rows.Next() {
		err := rows.Scan(&id, &first_name, &last_name, &email)
		checkError(err)
		user = User{id, first_name, last_name, email}
		users = append(users, user)
	}
	tmpl, err := template.ParseFiles("templates/home.gohtml")
	checkError(err)
	tmpl.Execute(w, users)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
