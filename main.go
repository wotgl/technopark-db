package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type InputRequest struct {
	method string
	url    string
	path   string
	json   string
	query  map[string][]string
}

func (ir *InputRequest) parse(r *http.Request) {
	ir.method = r.Method
	ir.url = fmt.Sprintf("%v", r.URL)
	ir.path = r.URL.EscapedPath()

	// ReadAll reads from r until an error or EOF and returns the data it read
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	ir.json = string(body)

	ir.query = r.URL.Query()
}

func api(w http.ResponseWriter, r *http.Request) {
	inputRequest := new(InputRequest)
	inputRequest.parse(r)

	fmt.Println(inputRequest)

	io.WriteString(w, "Hello world!")
	fmt.Println("")
}

func main() {
	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/db/api/", api)

	http.ListenAndServe(PORT, nil)

}
