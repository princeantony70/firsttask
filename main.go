package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	bcrypt "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error
var tpl *template.Template

type user struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
	Password  []byte
}

func init() {
	db, err = sql.Open("mysql", "root:nfn@/nfn")
	checkErr(err)
	err = db.Ping()
	checkErr(err)
//	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	defer db.Close()
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", index)
	http.HandleFunc("/userForm", userForm)
	http.HandleFunc("/createUsers", createUsers)
	//http.HandleFunc("/editUsers", editUsers)
	//http.HandleFunc("/deleteUsers", deleteUsers)
	//http.HandleFunc("/updateUsers", updateUsers)
	log.Println("Server is up on  port")
	log.Fatalln(http.ListenAndServe(":5053", nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	rows, e := db.Query(
		`SELECT id,
			username,
			first_name,
			last_name,
			password
		FROM users;
		`)
	checkErr(e)

	users := make([]user, 0)
	for rows.Next() {
		usr := user{}
		rows.Scan(&usr.ID, &usr.Username, &usr.FirstName, &usr.LastName, &usr.Password)
		users = append(users, usr)
	}
	log.Println(users)
//	tpl.ExecuteTemplate(w, "index.gohtml", users)
tpl,err :=template.ParseFiles("index.gohtml")
if err != nil{
	log.Fatalln("error parsing template",err)
}

err =tpl.ExecuteTemplate(w,"index.gohtml",users)
if err !=nil{
	log.Fatalln("error executing template",err)
}
}

func userForm(w http.ResponseWriter, req *http.Request) {
	//err = tpl.ExecuteTemplate(w, "userForm.gohtml", nil)
	tpl,err :=template.ParseFiles("userForm.gohtml")
	if err != nil{
		log.Fatalln("error parsing template",err)
	}

	err =tpl.ExecuteTemplate(w,"userForm.gohtml",nil)
	if err !=nil{
		log.Fatalln("error executing template",err)
	}
	checkErr(err)
}

func createUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		usr := user{}
		usr.Username = req.FormValue("username")
		usr.FirstName = req.FormValue("firstName")
		usr.LastName = req.FormValue("lastName")
		bPass, e := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.MinCost)
		checkErr(e)
		usr.Password = bPass
		_, e = db.Exec(
			"INSERT INTO users (username, first_name, last_name, password) VALUES (?, ?, ?, ?)",
			usr.Username,
			usr.FirstName,
			usr.LastName,
			usr.Password,
		)
		checkErr(e)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
}
