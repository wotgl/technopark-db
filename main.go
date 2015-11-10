package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"

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
			responseCode = 5
			errorMessage = map[string]interface{}{"msg": "Exist [Error 1452]"}

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

func createIvalidResponse() string {
	responseCode := 2
	errorMessage := map[string]interface{}{"msg": "Invalid"}

	resp, err := createResponse(responseCode, errorMessage)
	if err != nil {
		panic(err)
	}

	return resp
}

func validateJson(ir *InputRequest, args ...string) (string, error) {
	var resp string

	for _, value := range args {
		if reflect.TypeOf(ir.json[value]) == nil {
			resp = createIvalidResponse()
			return resp, errors.New("Invalid json")
		}
	}

	return "", nil
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
	defer stmt.Close()

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
}

func (u *User) create() string {
	var resp string
	query := "INSERT INTO user (username, about, name, email, isAnonymous) VALUES(?, ?, ?, ?, ?)"

	var args []interface{}

	resp, err := validateJson(u.inputRequest, "username", "about", "name", "email")
	if err != nil {
		return resp
	}

	args = append(args, u.inputRequest.json["username"])
	args = append(args, u.inputRequest.json["about"])
	args = append(args, u.inputRequest.json["name"])
	args = append(args, u.inputRequest.json["email"])

	// Validate isAnonymous param
	isAnonymous := u.inputRequest.json["isAnonymous"]
	if isAnonymous == nil {
		args = append(args, false)
	} else {
		if isAnonymous != false && isAnonymous != true {
			resp = createIvalidResponse()
			return resp
		}
		args = append(args, isAnonymous)
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

	respIsAnonymous, _ := strconv.ParseBool(newUser.values[0]["isAnonymous"])

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":       newUser.values[0]["about"],
		"email":       newUser.values[0]["email"],
		"id":          dbResp.lastId,
		"isAnonymous": respIsAnonymous,
		"name":        newUser.values[0]["name"],
		"username":    newUser.values[0]["username"],
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.create()")

	return resp
}

func (u *User) getUserDetails() (int, map[string]interface{}) {
	query := "SELECT * FROM user WHERE email = ?"
	var args []interface{}
	args = append(args, u.inputRequest.query["user"][0])

	getUser, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	if getUser.rows == 0 {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		return responseCode, errorMessage
	}

	// followers here
	query = "SELECT follower FROM follow WHERE followee = ?"
	getUserFollowers, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	listFollowers := make([]string, 0)
	for _, value := range getUserFollowers.values {
		listFollowers = append(listFollowers, value["follower"])
	}

	// following here
	query = "SELECT followee FROM follow WHERE follower = ?"
	getUserFollowing, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	listFollowing := make([]string, 0)
	for _, value := range getUserFollowing.values {
		listFollowing = append(listFollowing, value["followee"])
	}

	respIsAnonymous, _ := strconv.ParseBool(getUser.values[0]["isAnonymous"])

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":       getUser.values[0]["about"],
		"email":       getUser.values[0]["email"],
		"followers":   listFollowers,
		"following":   listFollowing,
		"id":          getUser.values[0]["id"],
		"isAnonymous": respIsAnonymous,
		"name":        getUser.values[0]["name"],
		"username":    getUser.values[0]["username"],
	}

	return responseCode, responseMsg
}

func (u *User) getDetails() string {
	var resp string
	if len(u.inputRequest.query["user"]) != 1 {
		resp = createIvalidResponse()
		return resp
	}

	responseCode, responseMsg := u.getUserDetails()

	resp, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("user.getDetails()")

	return resp
}

func (u *User) follow() string {
	var resp string
	query := "INSERT INTO follow (follower, followee) VALUES(?, ?)"

	var args []interface{}

	resp, err := validateJson(u.inputRequest, "follower", "followee")
	if err != nil {
		return resp
	}

	args = append(args, u.inputRequest.json["follower"])
	args = append(args, u.inputRequest.json["followee"])

	_, err = execQuery(query, &args, u.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	resp = u.getDetails()

	return resp
}

func (u *User) unfollow() string {
	var resp string
	query := "DELETE FROM follow WHERE follower = ? AND followee = ?"

	var args []interface{}

	resp, err := validateJson(u.inputRequest, "follower", "followee")
	if err != nil {
		return resp
	}

	args = append(args, u.inputRequest.json["follower"])
	args = append(args, u.inputRequest.json["followee"])

	_, err = execQuery(query, &args, u.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	resp = u.getDetails()

	return resp
}

func (u *User) updateProfile() string {
	var resp string
	query := "UPDATE user SET about = ?, name = ? WHERE email =  ?"

	var args []interface{}

	resp, err := validateJson(u.inputRequest, "about", "user", "name")
	if err != nil {
		return resp
	}

	args = append(args, u.inputRequest.json["about"])
	args = append(args, u.inputRequest.json["name"])
	args = append(args, u.inputRequest.json["user"])

	_, err = execQuery(query, &args, u.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["user"].(string))
	resp = u.getDetails()

	return resp
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

// =================
// Forum handler here
// =================

type Forum struct {
	inputRequest *InputRequest
	db           *sql.DB
}

func (f *Forum) create() string {
	var resp string
	query := "INSERT INTO forum (name, short_name, user) VALUES(?, ?, ?)"

	var args []interface{}

	resp, err := validateJson(f.inputRequest, "name", "short_name", "user")
	if err != nil {
		return resp
	}

	args = append(args, f.inputRequest.json["name"])
	args = append(args, f.inputRequest.json["short_name"])
	args = append(args, f.inputRequest.json["user"])

	dbResp, err := execQuery(query, &args, f.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	query = "SELECT * FROM forum WHERE id = ?"
	args = args[0:0]
	args = append(args, dbResp.lastId)
	newForum, err := selectQuery(query, &args, f.db)
	if err != nil {
		panic(err)
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"name":       newForum.values[0]["name"],
		"short_name": newForum.values[0]["short_name"],
		"id":         dbResp.lastId,
		"user":       newForum.values[0]["user"],
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("forum.create()")

	return resp
}

func (f *Forum) getForumDetails() (int, map[string]interface{}) {
	query := "SELECT * FROM forum WHERE short_name = ?"
	var args []interface{}
	args = append(args, f.inputRequest.query["forum"][0])

	getForum, err := selectQuery(query, &args, f.db)
	if err != nil {
		panic(err)
	}

	if getForum.rows == 0 {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"id":         getForum.values[0]["id"],
		"short_name": getForum.values[0]["short_name"],
		"name":       getForum.values[0]["name"],
		"user":       getForum.values[0]["user"],
	}

	return responseCode, responseMsg
}

func (f *Forum) details() string {
	var resp string
	var related bool
	if len(f.inputRequest.query["forum"]) != 1 {
		resp = createIvalidResponse()
		return resp
	}

	if len(f.inputRequest.query["related"]) == 1 && f.inputRequest.query["related"][0] == "user" {
		related = true
	}

	responseCode, responseMsg := f.getForumDetails()

	if related {
		u := User{inputRequest: f.inputRequest, db: f.db}
		u.inputRequest.query["user"] = append(u.inputRequest.query["user"], "sasha@localhost")
		_, userDetails := u.getUserDetails()

		responseMsg["user"] = userDetails
	}

	resp, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func forumHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	forum := Forum{inputRequest: inputRequest, db: db}

	var result string

	if inputRequest.method == "GET" {
		result = forum.details()
	} else if inputRequest.method == "POST" {

		// Like Router
		if inputRequest.path == "/db/api/forum/create/" {
			result = forum.create()
		}
		// else if inputRequest.path == "/db/api/user/follow/" {
		// 	result = user.follow()
		// } else if inputRequest.path == "/db/api/user/unfollow/" {
		// 	result = user.unfollow()
		// } else if inputRequest.path == "/db/api/user/updateProfile/" {
		// 	result = user.updateProfile()
		// }
	}

	io.WriteString(w, result)
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
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	http.HandleFunc("/db/api/user/", makeHandler(db, userHandler))
	http.HandleFunc("/db/api/forum/", makeHandler(db, forumHandler))

	http.ListenAndServe(PORT, nil)

	fmt.Println(reflect.TypeOf("str"))
}
