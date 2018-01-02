package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        string `bson:"_id,omitempty"`
	UserName  string
	Email     string
	FirstName string
	LastName  string
	Age       string
}

// Load the login.html template.
var tmpl = template.Must(template.New("tmpl").ParseFiles("login.html"))

func main() {	
	// Serve / with the login.html file.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	session, err := mgo.Dial("mongodb://127.0.0.1/GoTesting");
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//variable of type User (struct)
	var results User

	c := session.DB("test_v1").C("Col1")
	//Mongo query to select the user's account info matching our struct (everything but the password)
	err = c.Find(bson.M{"username": "Test User"}).Select(results).One(&results)
	// , "password": hashedPassword}
	if err != nil {
		log.Fatal(err)
	}

	//Store field values as string for output

	accountID := (bson.ObjectId(results.ID).Hex())
	accountName := string(results.UserName)
	accountFirst := string(results.FirstName)
	accountLast := string(results.LastName)
	accountEmail := string(results.Email)
	accountAge := string(results.Age)

	//Write back to the login.html file with our user's info
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<b>Account #:</b> "+accountID+"<br>")
		fmt.Fprintln(w, "<b>Acc Name:</b> "+accountName+"<br>")
		fmt.Fprintln(w, "<b>Email:</b> "+accountEmail+"<br>")
		fmt.Fprintln(w, "<b>First:</b> "+accountFirst+"<br>")
		fmt.Fprintln(w, "<b>Last:</b> "+accountLast+"<br>")
		fmt.Fprintln(w, "<b>Age:</b> "+accountAge+"<br>")
	})

	log.Fatal(http.ListenAndServe(":5000", nil))

}
