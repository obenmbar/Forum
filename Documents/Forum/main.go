package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	forimino "forumino/database"
	forumino "forumino/handlers"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatal("USAGE: go run main.go")
	}
	db := forimino.InitDB()
	defer db.Close()
	http.HandleFunc("/", Home_handler)
	// http.HandleFunc("/Registre/", Registre_handler)
	// http.HandleFunc("/login", Login_andlder)
	// http.HandleFunc("/comment", HandleCommente)
	http.HandleFunc("/CreatePost", func(w http.ResponseWriter, r *http.Request){
		forumino.Creat_PostGlobale(w,r,db)
	} )
	 http.HandleFunc("/Reaction",  func(w http.ResponseWriter, r *http.Request){
		forumino.ReactHandler(w,r,db)
	 })
	// http.HandleFunc("/Logout", Logout)
	fmt.Println("server started on: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("failed to start server: ", err)
		return
	}
}

func Home_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"template/index.html")
}

