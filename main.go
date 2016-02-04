package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

func initLog() {
	fmt.Println("Hey, initLog() here")
	f, err := os.OpenFile("main2.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	for {
	}
}

// =================
// Main here
// =================

func makeHandler(db *sql.DB, fn func(http.ResponseWriter, *http.Request, *InputRequest, *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inputRequest := new(InputRequest)
		inputRequest.parse(r)

		fn(w, r, inputRequest, db)
	}
}

func main() {
	db, err := sql.Open("mysql", "sasha1003:10031995@/mydb")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("db ok")
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// args here
	argsWithProg := os.Args[1:]
	MAX_DB_CONNECTIONS := int(stringToInt64(argsWithProg[1]))
	db.SetMaxOpenConns(MAX_DB_CONNECTIONS)
	PORT := ":" + argsWithProg[0]

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/db/api/user/", makeHandler(db, userHandler))
	http.HandleFunc("/db/api/forum/", makeHandler(db, forumHandler))
	http.HandleFunc("/db/api/thread/", makeHandler(db, threadHandler))
	http.HandleFunc("/db/api/post/", makeHandler(db, postHandler))
	http.HandleFunc("/db/api/status/", makeHandler(db, statusHandler))
	http.HandleFunc("/db/api/clear/", makeHandler(db, clearHandler))

	http.ListenAndServe(PORT, nil)
}
