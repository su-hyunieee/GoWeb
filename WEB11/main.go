package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Email string
	Age int
}

func (u User) IsOld() bool {
	return u.Age > 30
}

func main() {
	user := User{Name:"sh", Email:"sh@naver.com", Age:23}
	user2 := User{Name:"aaa", Email:"sss@ga.mcom", Age:40}
	users := []User{user, user2}

	tmpl, err := template.New("Tmpl1").ParseFiles("templates/tmpl1.tmpl", "templates/tmpl2.tmpl")
	if err !=nil {
		panic(err)
	}

	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user)
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl1.tmpl", user2)
}