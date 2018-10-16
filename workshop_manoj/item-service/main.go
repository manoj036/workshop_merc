package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
)

var (
	mysqlUser = os.Getenv("MYSQL_USER")
	mysqlPass = os.Getenv("MYSQL_PASS")
	mysqlDB   = os.Getenv("MYSQL_DB")
	mysqlHost = os.Getenv("MYSQL_HOST")
	mysqlPort = os.Getenv("MYSQL_PORT")
)


type item struct {
	Id				int		`json:"id"`
	Name			string	`json:"name"`
	Price			int		`json:"price"`
	Description 	string	`json:"description"`
}

func LoadDatabase()(*sql.DB,error)  {
	// Set defaults
	if mysqlUser == "" {
		mysqlUser = "root"
	}
	if mysqlDB == "" {
		mysqlDB = "workshopdb"
	}
	if mysqlHost == "" {
		mysqlHost = "127.0.0.1"
	}
	if mysqlPass == "" {
		mysqlPass = "User123$"
	}
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDB)
	fmt.Printf("Connecting mysql at: %s", dsn)

	db, err := sql.Open("mysql", dsn)
	return db,err
}

func ErrorCheck(err error)  {
	if err != nil {
		panic(err.Error())
	}
}

func GetItems(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db, err := LoadDatabase()
	ErrorCheck(err)
	rows, err := db.Query("SELECT * FROM items")
	ErrorCheck(err)
	defer rows.Close()

	items:=make([]item,0)
	for rows.Next(){
		tmpItem:=new(item)
		err:= rows.Scan(&tmpItem.Id,&tmpItem.Name,&tmpItem.Description,&tmpItem.Price)
		ErrorCheck(err)

		items=append(items,*tmpItem)
	}
	json.NewEncoder(w).Encode(items)

	defer db.Close()
}


func GetItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	idGiven := params.ByName("id")
	fmt.Fprintf(w, "GET %s\n", idGiven)
	db, err := LoadDatabase()
	ErrorCheck(err)
	rows, err := db.Prepare("SELECT * FROM items where id=?")
	ErrorCheck(err)

	defer rows.Close()

	tmpItem:=new(item)
	err=rows.QueryRow(idGiven).Scan(&tmpItem.Id,&tmpItem.Name,&tmpItem.Description,&tmpItem.Price)
	ErrorCheck(err)

	json.NewEncoder(w).Encode(tmpItem)

	defer db.Close()
}

func PostItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	tmpItem:=new(item)
	tmpItem.Name= r.FormValue("name")
	tmpItem.Description=r.FormValue("description")
	tmpItem.Price,_=strconv.Atoi(r.FormValue("price"))
	db, err := LoadDatabase()
	ErrorCheck(err)
	defer db.Close()
	rows, err := db.Prepare("INSERT INTO items(name,price,description) VALUES(?,?,?) ")
	ErrorCheck(err)
	defer rows.Close()

	priceTemp,_:=strconv.Atoi(string(tmpItem.Price))
	rows.Exec(tmpItem.Name, priceTemp,tmpItem.Description)
	fmt.Fprintf(w, "POST %s\n", tmpItem.Name)
}

func UpdateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params){
	idGiven := params.ByName("id")
	fmt.Fprintf(w, "GET %s\n", idGiven)

	tmpItem:=new(item)
	tmpItem.Name= r.FormValue("name")
	tmpItem.Description=r.FormValue("description")
	tmpItem.Price,_=strconv.Atoi(r.FormValue("price"))

	db, err := LoadDatabase()
	ErrorCheck(err)
	defer db.Close()

	name, err := db.Prepare("UPDATE items SET name=? where id=?")
	ErrorCheck(err)
	defer name.Close()

	price, err := db.Prepare("UPDATE items SET price=? where id=?")
	ErrorCheck(err)
	defer price.Close()

	description, err := db.Prepare("UPDATE items SET description=? where id=?")
	ErrorCheck(err)
	defer description.Close()

	if tmpItem.Name!=""{
		name.Exec(tmpItem.Name,idGiven)
	}

	if tmpItem.Description!=""{
		description.Exec(tmpItem.Description,idGiven)
	}

	if r.FormValue("price")!=""{
		price.Exec(r.FormValue("price"),idGiven)
	}
}

func main(){
	router := httprouter.New()

	router.GET("/items", GetItems)
	router.GET("/items/:id", GetItem)
	router.POST("/items", PostItem)
	router.PUT("/items/:id", UpdateItem)
	http.ListenAndServe(":8080", router)
}