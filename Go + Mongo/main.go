package main

import (
	"encoding/json" //allow us to write our struct as a json
	"fmt"           //allow us to print
	"html/template" //allows us to output html file
	"log"           //allows us to have a logger (using fatal in this case)
	"net/http"      // allows us to create a wejsonFileerver
	"os"            //operating system. Used to create and read files.

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2"
)

type Person struct {
	UserName  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Age       string
}

func main() {
	//Creating our server on port 51553.
	http.HandleFunc("/", form)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

//Function to take user input, and write to a JSON file we create.
func form(w http.ResponseWriter, r *http.Request) {

	//Setting up MongoDB connection
	session, err := mgo.Dial("mongodb://testUser:password@ds051553.mlab.com:51553/test_v1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if r.Method == "GET" {
		template, _ := template.ParseFiles("index.html")
		template.Execute(w, nil)
	} else {

		r.ParseForm()
		//Storing the user's info
		var userName = r.FormValue("userName")
		var email = r.FormValue("email")
		var password = r.FormValue("password")
		var firstName = r.FormValue("fName")
		var lastName = r.FormValue("lName")
		var userAge = r.FormValue("age")

		//hashing password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		//storing user data into struct
		p1 := Person{UserName: userName, Email: email, Password: string(hashedPassword), FirstName: firstName, LastName: lastName, Age: userAge}

		//Create test.json file.
		dataFile, err := os.Create("test.json")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		//Convert struct to json.
		jsonFile, err := json.MarshalIndent(p1, "", "    ")
		if err != nil {
			fmt.Println(err)
		}
		//Write to json file.
		fmt.Fprint(dataFile, string(jsonFile))

		//Inserting user into MongoDB
		c := session.DB("test_v1").C("Col1")
		err = c.Insert(p1)
		if err != nil {
			log.Fatal(err)
		}

	}
}
