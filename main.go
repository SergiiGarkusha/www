package main

/*
	Нужно сменить каталог при помощи команды (при необходимости):
		cd www
	Установить go в текущий каталог (при необходимости):
		go install
	Запустить на исполнение программу:
		go run main.go
	Установим дополнительную библиотеку для динамического отслеживания url gorilla/mux

		go get github.com/gorilla/mux
		go mod init
		go mod tidy

	Включение отключение брандмауэра
		netsh advfirewall set allprofiles state off
		netsh advfirewall set allprofiles state on

		git init
		git add .
		git commit -m "to heroku app"
		git push

		git push heroku master

*/

import (
	"database/sql"
	"fmt"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Arctical struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Arctical{}
var showPost = Arctical{}

func index(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter - запись данных на страницу,
	// r *http.Request - получение данных от страницы
	// Загружаем шаблоны с необходимой разметкой html при создании глваной страницы / (для ввода данных)
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	// Подключение к базе данных mysql
	// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	db, err := sql.Open("mysql", "ezyro_32990016:chvz7tmk@tcp(sql302.ezyro.com)/ezyro_32990016_golang")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articales`")
	if err != nil {
		panic(err)
	}

	posts = []Arctical{}
	for res.Next() {
		var post Arctical
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	// "index" - наименование подключаемого блока
	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter - запись данных на страницу,
	// r *http.Request - получение данных от страницы
	// Загружаем шаблоны с необходимой разметкой html при создании страницы по url адрессу /create (для ввода данных)
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	// "create" - наименование подключаемого блока
	t.ExecuteTemplate(w, "create", nil)
}

func save_articale(w http.ResponseWriter, r *http.Request) {

	// Считываем значения из ячеек формы ввода в переменные
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "не все данные заполнены")
	} else {

		// Подключение к базе данных mysql
		// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
		db, err := sql.Open("mysql", "ezyro_32990016:chvz7tmk@tcp(sql302.ezyro.com)/ezyro_32990016_golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		// Ввод данных в базу данных mysql
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articales` (`title`,`anons`,`full_text`) VALUES('%s', '%s','%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		// переопределяем url адресс с текуйщей страницы на главную, после нажатия кнопки
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	vars := mux.Vars(r)
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "ID: %v\n", vars["id"])

	// Подключение к базе данных mysql
	// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	db, err := sql.Open("mysql", "ezyro_32990016:chvz7tmk@tcp(sql302.ezyro.com)/ezyro_32990016_golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articales` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Arctical{}
	for res.Next() {
		var post Arctical
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	// "show" - наименование подключаемого блока
	t.ExecuteTemplate(w, "show", showPost)

}

func HandleFunc() {
	// Подключение gorilla/mux
	// mux.NewRouter() - создание объекта роутер и инициализация переменной rtr
	rtr := mux.NewRouter()

	// обработчик url адресса главнойй страницы
	rtr.HandleFunc("/", index).Methods("GET")

	// обработчик url адресса для перехода на страницу ввода данных
	rtr.HandleFunc("/create", create).Methods("GET")

	// обработчик url адресса при обработке ввода данных в базу данных
	rtr.HandleFunc("/save_articale", save_articale).Methods("POST")

	// обработка динамического url
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)

	// Подключение статических файлов
	// <body class="d-flex h-100 text-center text-bg-dark"> в header.html нужно убрать класс text-bg-dark тогда фон будет переопределен
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// // Запускаем локальный сервер
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("web server started")
	HandleFunc()
}
