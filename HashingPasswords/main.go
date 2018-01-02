package main

import (
	"encoding/json" //allow us to write our struct as a json
	"fmt"           //allow us to print
	"html/template" //allows us to output html file
	"log"           //allows us to have a logger (using fatal in this case)
	"net/http"      // allows us to create a webserver
	"os"            //operating system. Used to create and read files.

	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Name           string
	Pass           string
	HashedPassword string
}

//Creating our server on port 9000.
func main() {
	http.HandleFunc("/", login)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Function to take user input, and write to a JSON file we create.
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		template, _ := template.ParseFiles("login.html")
		template.Execute(w, nil)
	} else {
		r.ParseForm()
		//Storing the user's name.
		var userName = r.FormValue("name")

		//Storing the user's password
		var userPass = r.FormValue("pass")

		//Encrypting the user password so it can be stored safely.
		bsPass, _ := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.MinCost)
		stringPass := string(bsPass)
		fmt.Println(stringPass)

		//storing user data
		p1 := Person{Name: userName, Pass: userPass, HashedPassword: stringPass}

		//Output user's data that has been placed in our Person Struct to the console.
		fmt.Printf("%+v", p1)

		//Create test.json file.
		dataFile, err := os.Create("userInfo.json")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}

		//Convert struct to json and format to look pretty.
		b, err := json.MarshalIndent(p1, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}

		//Write to json file.
		fmt.Fprint(dataFile, string(b))
	}

}
