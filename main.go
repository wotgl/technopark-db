package main

import (
	"encoding/json"
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

func user(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest) {
	if inputRequest.method == "GET" {
		fmt.Println(inputRequest.query)
	} else if inputRequest.method == "POST" {
		type InputJson struct {
			Key1 string
			Key2 string
			Key3 string
		}

		jsonByteArray := make([]byte, len(inputRequest.json))
		copy(jsonByteArray[:], inputRequest.json)

		var inputJson InputJson

		err := json.Unmarshal(jsonByteArray, &inputJson)
		if err != nil {
			fmt.Println("error:\t", err)
		}

		fmt.Println(inputJson)
	}

	io.WriteString(w, "Hello world!")
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *InputRequest)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inputRequest := new(InputRequest)
		inputRequest.parse(r)

		fn(w, r, inputRequest)
	}
}

func main() {
	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/db/api/user/", makeHandler(user))

	http.ListenAndServe(PORT, nil)
}
