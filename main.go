package main

import (
	"fmt"
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/", hello)
	http.ListenAndServe(PORT, nil)

}
