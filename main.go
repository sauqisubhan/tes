package main

import (
	"html/template"
	"net/http"
)

type User struct {
	ID    int
	Name  string
	Email string
	Nomor string
}

var tmpl = template.Must(template.ParseGlob("templates/*"))

func main() {
	InitDB()

	http.HandleFunc("/", Index)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, name, email, nomor FROM users ORDER BY id ASC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Nomor)
		users = append(users, u)
	}

	tmpl.ExecuteTemplate(w, "index.html", users)
}

func Add(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "form.html", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		nomor := r.FormValue("nomor")
		_, err := DB.Exec("INSERT INTO users(name, email, nomor) VALUES($1, $2, $3)", name, email, nomor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := DB.QueryRow("SELECT id, name, email, nomor FROM users WHERE id=$1", id)
	var u User
	row.Scan(&u.ID, &u.Name, &u.Email, &u.Nomor)
	tmpl.ExecuteTemplate(w, "form.html", u)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		nomor := r.FormValue("nomor")
		_, err := DB.Exec("UPDATE users SET name=$1, email=$2, nomor=$3 WHERE id=$4", name, email, nomor, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
