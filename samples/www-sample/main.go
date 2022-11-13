package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name  string
	Age   uint16
	Money int16

	Avg_grades float64
	Happiness  float64
	Hobbies    []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("User naime is: %s, Не is: %d and hi has maney equal %d", u.Name, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Football", "Skate", "Dance"}}
	// bob.setNewName("Alex")
	// fmt.Fprintf(w, bob.getAllInfo())
	// fmt.Fprintf(w, `<b>Main text</b>
	// 				<h1>MAIN TEXT</h1>`)

	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(w, bob)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ours contacts")
}

func handlRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts", contacts_page)
	http.ListenAndServe(":8080", nil)
}

func main() {
	//bob := User{name: "Bob", age: 25, money: -50, avg_grades: 4.2, happiness: 0.8}
	handlRequest()
}
