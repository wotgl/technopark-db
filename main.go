package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	mysql "github.com/go-sql-driver/mysql"
)

type InputRequest struct {
	method string
	url    string
	path   string
	json   map[string]interface{}
	query  map[string][]string
}

func (ir *InputRequest) parse(r *http.Request) {
	ir.method = r.Method
	ir.url = fmt.Sprintf("%v", r.URL)
	ir.path = r.URL.EscapedPath()

	// POST JSON
	body, err := ioutil.ReadAll(r.Body) // ReadAll reads from r until an error or EOF and returns the data it read
	if err != nil {
		panic(err)
	}

	var parsed map[string]interface{}
	json.Unmarshal([]byte(body), &parsed)
	ir.json = parsed

	// GET Query
	ir.query = r.URL.Query()
}

func createResponse(code int, response map[string]interface{}) (string, error) {
	cacheContent := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(cacheContent)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return "", err
	}

	return string(str), nil
}

func errorExecParse(err error) (int, map[string]interface{}) {
	if driverErr, ok := err.(*mysql.MySQLError); ok { // Now the error number is accessible directly
		var responseCode int
		var errorMessage map[string]interface{}

		switch driverErr.Number {
		case 1062:
			responseCode = 5
			errorMessage = map[string]interface{}{"msg": "Exist"}

		// Error 1452: Cannot add or update a child row: a foreign key constraint fails
		case 1452:
			responseCode = 1
			errorMessage = map[string]interface{}{"msg": "Not found"}

		default:
			fmt.Println("errorExecParse() default")
			panic(err.Error())
			responseCode = 4
			errorMessage = map[string]interface{}{"msg": "Unknown Error"}
		}

		return responseCode, errorMessage
	}
	panic(err.Error()) // proper error handling instead of panic in your app
}

// ======================
// Database queries here
// ======================

func createArgs(stringArgs []string) *[]interface{} {
	args := make([]interface{}, len(stringArgs))
	for i, s := range stringArgs {
		args[i] = s
	}

	return &args
}

func createArgsFromJson(stringArgs map[string]interface{}) *[]interface{} {
	// args := make([]interface{}, len(stringArgs))
	var args []interface{}
	for _, s := range stringArgs {
		args = append(args, s)
	}

	return &args
}

type ExecResponse struct {
	lastId   int64
	rowCount int64
}

func execQuery(query string, args *[]interface{}, db *sql.DB) (*ExecResponse, error) {
	resp := new(ExecResponse)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(*args...)
	if err != nil {
		return nil, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	resp.lastId = lastId
	resp.rowCount = rowCnt

	return resp, nil
}

type SelectResponse struct {
	rows    int
	columns []string
	values  []map[string]string
}

func selectQuery(query string, args *[]interface{}, db *sql.DB) (*SelectResponse, error) {
	resp := new(SelectResponse)

	rows, err := db.Query(query, *args...)
	if err != nil {
		panic(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
		return nil, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var respValues []map[string]string

	// Fetch rows
	for rows.Next() {
		x := make(map[string]string)

		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
			return nil, err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			x[columns[i]] = value
			// x = append(x, value)
		}

		respValues = append(respValues, x)
	}

	resp.columns = columns
	resp.values = respValues
	resp.rows = len(resp.values)

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}

	return resp, nil
}

// =================
// User handler here
// =================

type User struct {
	inputRequest *InputRequest
	db           *sql.DB
	userResponse struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		About       string `json:"about"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		IsAnonymous bool   `json:"isAnonymous"`
		Date        string `json:"date"`
	}
}

func (u *User) checkUser(email string) error {
	stmtOut, err := u.db.Prepare("SELECT * FROM user WHERE email = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// user := map[string]interface{}{
	// 	"Id":          0,
	// 	"Username":    "",
	// 	"About":       "",
	// 	"Name":        "",
	// 	"Email":       "",
	// 	"IsAnonymous": false,
	// 	"Date":        "",
	// }

	// user := map[string]interface{}{
	// 	"Id":          int,
	// 	"Username":    string,
	// 	"About":       "",
	// 	"Name":        "",
	// 	"Email":       "",
	// 	"IsAnonymous": "",
	// 	"Date":        "",
	// }

	err = stmtOut.QueryRow(email).Scan(&u.userResponse.Id, &u.userResponse.Username, &u.userResponse.About, &u.userResponse.Name, &u.userResponse.Email, &u.userResponse.IsAnonymous, &u.userResponse.Date) // WHERE number = 13
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return err
}

func (u *User) create() string {
	var resp string
	query := "INSERT INTO user (username, about, name, email, isAnonymous) VALUES(?, ?, ?, ?, ?)"

	var args []interface{}

	args = append(args, u.inputRequest.json["username"])
	args = append(args, u.inputRequest.json["about"])
	args = append(args, u.inputRequest.json["name"])
	args = append(args, u.inputRequest.json["email"])
	if u.inputRequest.json["isAnonymous"] == nil {
		args = append(args, false)
	} else {
		args = append(args, u.inputRequest.json["isAnonymous"])
	}

	dbResp, err := execQuery(query, &args, u.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	query = "SELECT * FROM user WHERE id = ?"
	args = args[0:0]
	args = append(args, dbResp.lastId)
	newUser, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":       newUser.values[0]["about"],
		"email":       newUser.values[0]["email"],
		"id":          dbResp.lastId,
		"isAnonymous": newUser.values[0]["isAnonymous"],
		"name":        newUser.values[0]["name"],
		"username":    newUser.values[0]["username"],
	}

	response, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.create()")

	return response

	// args := createArgs()
	/*
		// Prepare statement for inserting data
		stmtIns, err := u.db.Prepare("INSERT INTO user (username, about, name, email, isAnonymous) VALUES(?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates



		res, err := stmtIns.Exec(username, about, name, email, isAnonymous)
		if err != nil {
			responseCode, errorMessage := errorExecParse(err)

			response, err := createResponse(responseCode, errorMessage)
			if err != nil {
				panic(err.Error())
			}

			return response

			// panic(err.Error()) // proper error handling instead of panic in your app
		}

		insertId, err := res.LastInsertId()

		responseCode := 0
		responseMsg := map[string]interface{}{
			"about":       about,
			"email":       email,
			"id":          insertId,
			"isAnonymous": isAnonymous,
			"name":        name,
			"username":    username,
		}

		response, err := createResponse(responseCode, responseMsg)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("user.create()")

		return response
	*/
}

func (u *User) getDetails() string {
	// Prepare statement for reading data
	stmtOut, err := u.db.Prepare("SELECT * FROM user WHERE email = ?;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	type User struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		About       string `json:"about"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		IsAnonymous bool   `json:"isAnonymous"`
		Date        string `json:"date"`
	}

	user := new(User)

	// Query the square-number of 13
	err = stmtOut.QueryRow(u.inputRequest.query["user"][0]).Scan(&user.Id, &user.Username, &user.About, &user.Name, &user.Email, &user.IsAnonymous, &user.Date) // WHERE number = 13
	if err != nil {
		if err == sql.ErrNoRows {
			var responseCode int
			var errorMessage map[string]interface{}

			responseCode = 1
			errorMessage = map[string]interface{}{"msg": "Doesn`t exist"}

			response, err := createResponse(responseCode, errorMessage)
			if err != nil {
				panic(err.Error())
			}

			fmt.Println("user.getDetails()")
			return response
		}
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	type Response struct {
		Code     int   `json:"code"`
		Response *User `json:"response"`
	}

	tempResponse := Response{0, user}

	response, err := json.Marshal(tempResponse)
	if err != nil {
		panic(err)
	}

	fmt.Println("user.getDetails()")

	return string(response)
}

func (u *User) follow() string {
	stmtIns, err := u.db.Prepare("INSERT INTO follow (follower, followee) VALUES(?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	follower := u.inputRequest.json["follower"]
	followee := u.inputRequest.json["followee"]

	_, err = stmtIns.Exec(follower, followee)

	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		response, err := createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err.Error())
		}

		return response

		panic(err.Error()) // proper error handling instead of panic in your app
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about": "lol",
	}

	response, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.follow()")

	return response
}

func (u *User) unfollow() string {
	stmtIns, err := u.db.Prepare("DELETE FROM follow WHERE follower = ? AND followee = ?") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	follower := u.inputRequest.json["follower"]
	followee := u.inputRequest.json["followee"]

	_, err = stmtIns.Exec(follower, followee)

	if err != nil {
		if err == sql.ErrNoRows {
			var responseCode int
			var errorMessage map[string]interface{}

			responseCode = 1
			errorMessage = map[string]interface{}{"msg": "Doesn`t exist"}

			response, err := createResponse(responseCode, errorMessage)
			if err != nil {
				panic(err.Error())
			}

			fmt.Println("user.unfollow()")
			return response
		}
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about": "lol",
	}

	response, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.unfollow()")

	return response
}

func (u *User) updateProfile() string {
	email, ok := u.inputRequest.json["email"].(string)
	if !ok {
		panic(ok)
	}
	about, ok := u.inputRequest.json["about"].(string)
	if !ok {
		panic(ok)
	}
	name, ok := u.inputRequest.json["name"].(string)
	if !ok {
		panic(ok)
	}

	// Check user exist
	err := u.checkUser(email)
	if err != nil {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		response, err := createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err.Error())
		}

		return response
	}

	stmtIns, err := u.db.Prepare("UPDATE user SET about = ?, name = ? WHERE email =  ?") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// about := u.inputRequest.json["about"]
	// name := u.inputRequest.json["name"]
	// email := u.inputRequest.json["email"]

	fmt.Println("LOLCA1")

	_, err = stmtIns.Exec(about, name, email)
	fmt.Println("LOLCA2")

	// insertId, err := res.LastInsertId()

	// responseMsg := map[string]interface{}{
	// 	"about":       about,
	// 	"email":       email,
	// 	"id":          insertId,
	// 	"isAnonymous": isAnonymous,
	// 	"name":        name,
	// 	"username":    username,
	// }

	// u.userResponse.About = about
	// u.userResponse.Name = name

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":       u.userResponse.About,
		"email":       u.userResponse.Email,
		"id":          u.userResponse.Id,
		"isAnonymous": u.userResponse.IsAnonymous,
		"name":        u.userResponse.Name,
		"username":    u.userResponse.Username,
	}

	response, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.updateProfile()")

	return response
}

func userHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	user := User{inputRequest: inputRequest, db: db}

	var result string

	if inputRequest.method == "GET" {
		result = user.getDetails()
	} else if inputRequest.method == "POST" {

		// Like Router
		if inputRequest.path == "/db/api/user/create/" {
			result = user.create()
		} else if inputRequest.path == "/db/api/user/follow/" {
			result = user.follow()
		} else if inputRequest.path == "/db/api/user/unfollow/" {
			result = user.unfollow()
		} else if inputRequest.path == "/db/api/user/updateProfile/" {
			result = user.updateProfile()
		}
	}

	io.WriteString(w, result)
}

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
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	/*
		stmtOut, err := db.Prepare("SELECT name FROM user WHERE id > ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()

		var squareNum string

		for i := 0; i < 1; i++ {
			fmt.Println("1")
			// Query the square-number of 13
			err = stmtOut.QueryRow("'SELECT").Scan(&squareNum) // WHERE number = 13
			fmt.Println("2")
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			fmt.Println("3")
			fmt.Printf("The square number of 13 is: %v\n", squareNum)
		}
	*/

	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/db/api/user/", makeHandler(db, userHandler))

	http.ListenAndServe(PORT, nil)

	fmt.Println(reflect.TypeOf("str"))
}
