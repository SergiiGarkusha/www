package main

// Для подлючения драйвера mySql (с версии 1.6.0) нужно прописать в командной строке:
// 1. go mod init main
// 2. go get github.com/go-sql-driver/mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

func createDB(dbName string) {

	// Подключение к базе данных mysql
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создание базы данных с заданным именем
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		panic(err)
	}

	// Использование базы данных с заданным именем
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		panic(err)
	}

	// Создание таблицы example с полями id и data в выбранной базе данных
	_, err = db.Exec("CREATE TABLE example (id  integer, data varchar(32))")
	if err != nil {
		panic(err)
	}
}

func dropDB(dbName string) {
	// Подключение к базе данных mysql
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Использование базы данных с заданным именем
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		panic(err)
	}

	// Использование базы данных с заданным именем
	_, err = db.Exec("DROP DATABASE " + dbName)
	if err != nil {
		panic(err)
	}
}

func main() {

	// Создание новой базы данных
	createDB("test")
	// Удаление созданной базы данных
	dropDB("test")
	// Подключение к базе данных mysql
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}

	fmt.Println("Подключено к MySQL")

	defer db.Close()

	// Ввод данных в базу данных mysql
	// insert, err := db.Query("INSERT INTO `users` (`name`,`age`) VALUES('Anna', 37)")
	// if err != nil {
	// 	panic(err)
	// }
	// defer insert.Close()

	res, err := db.Query("SELECT `name`,`age` FROM `users`")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var user User
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		//fmt.Println(fmt.Sprintf("User: %s with age %d", user.Name, user.Age))
		fmt.Printf("User: %s with age %d \n", user.Name, user.Age)
	}

	fmt.Println("Отключено от MySQL")
}
