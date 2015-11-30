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
	"strings"
	"technopark-db/response"

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

func createResponseFromArray(code int, response []map[string]interface{}) (string, error) {
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

func checkError1062(err error) bool {
	if driverErr, ok := err.(*mysql.MySQLError); ok { // Now the error number is accessible directly
		if driverErr.Number == 1062 {
			return true
		}
	}
	return false
}

func createInvalidResponse() string {
	responseCode := 2
	errorMessage := map[string]interface{}{"msg": "Invalid"}

	resp, err := createResponse(responseCode, errorMessage)
	if err != nil {
		panic(err)
	}

	return resp
}

func createNotExistResponse() string {
	responseCode := 1
	errorMessage := map[string]interface{}{"msg": "Not exist"}

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
			resp = createInvalidResponse()
			return resp, errors.New("Invalid json")
		}
	}

	return "", nil
}

func validateBoolParams(json map[string]interface{}, args *[]interface{}, params ...string) (string, error) {
	var resp string

	for _, value := range params {

		if json[value] == nil {
			*args = append(*args, false)
		} else {
			if reflect.TypeOf(json[value]).Kind() != reflect.Bool {
				resp := createInvalidResponse()
				return resp, errors.New("Invalid json")
			}
			*args = append(*args, json[value])
		}
	}

	return resp, nil
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
			resp = createInvalidResponse()
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
	listFollowers := u.getUserFollowers(args[0].(string))

	// following here
	listFollowing := u.getUserFollowing(args[0].(string))

	// subscriptions here
	listSubscriptions := u.getUserSubscriptions(args[0].(string))

	respIsAnonymous, _ := strconv.ParseBool(getUser.values[0]["isAnonymous"])
	respId, _ := strconv.ParseInt(getUser.values[0]["id"], 10, 64)

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":         getUser.values[0]["about"],
		"email":         getUser.values[0]["email"],
		"followers":     listFollowers,
		"following":     listFollowing,
		"id":            respId,
		"isAnonymous":   respIsAnonymous,
		"name":          getUser.values[0]["name"],
		"subscriptions": listSubscriptions,
		"username":      getUser.values[0]["username"],
	}

	return responseCode, responseMsg
}

func (u *User) getUserFollowers(followee string) []string {
	var args []interface{}
	args = append(args, followee)

	// followers here
	query := "SELECT follower FROM follow WHERE followee = ?"
	getUserFollowers, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	listFollowers := make([]string, 0)
	for _, value := range getUserFollowers.values {
		listFollowers = append(listFollowers, value["follower"])
	}

	return listFollowers
}

func (u *User) getUserFollowing(follower string) []string {
	var args []interface{}
	args = append(args, follower)

	// following here
	query := "SELECT followee FROM follow WHERE follower = ?"
	getUserFollowing, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	listFollowing := make([]string, 0)
	for _, value := range getUserFollowing.values {
		listFollowing = append(listFollowing, value["followee"])
	}

	return listFollowing
}

func (u *User) getUserSubscriptions(user string) []int {
	var args []interface{}
	args = append(args, user)

	// subscriptions here
	query := "SELECT thread FROM subscribe WHERE user = ? ORDER BY thread asc"
	getUserSubscriptions, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	listSubscriptions := make([]int, 0)
	for _, value := range getUserSubscriptions.values {
		listSubscriptions = append(listSubscriptions, stringToInt(value["thread"]))
	}

	return listSubscriptions
}

func (u *User) getDetails() string {
	var resp string
	if len(u.inputRequest.query["user"]) != 1 {
		resp = createInvalidResponse()
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
		// return exist
		if checkError1062(err) == true {
			for k := range u.inputRequest.query {
				delete(u.inputRequest.query, k)
			}
			u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))

			resp = u.getDetails()
			return resp
		}
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

func (u *User) listBasic(query string) string {
	var resp string
	var args []interface{}

	// Validate query values
	if len(u.inputRequest.query["user"]) == 1 {
		args = append(args, u.inputRequest.query["user"][0])
	} else {
		resp = createInvalidResponse()
		return resp
	}

	// Check and validate optional params
	if len(u.inputRequest.query["since_id"]) >= 1 {
		query += " AND id >= ?"
		args = append(args, u.inputRequest.query["since_id"][0])
	}
	if len(u.inputRequest.query["order"]) >= 1 {
		orderType := u.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			resp = createInvalidResponse()
			return resp
		}

		query += fmt.Sprintf(" ORDER BY date %s", orderType)
	}
	if len(u.inputRequest.query["limit"]) >= 1 {
		limitValue := u.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			resp = createInvalidResponse()
			return resp
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Prepare users
	getUserFollowers, err := selectQuery(query, &args, u.db)
	if err != nil {
		panic(err)
	}

	responseCode := 0
	var responseMsg []map[string]interface{}
	for _, value := range getUserFollowers.values {
		respIsAnonymous, _ := strconv.ParseBool(value["isAnonymous"])
		respId, _ := strconv.ParseInt(value["id"], 10, 64)

		// followers here
		listFollowers := u.getUserFollowers(value["email"])

		// following here
		listFollowing := u.getUserFollowing(value["email"])

		// subscriptions here
		listSubscriptions := u.getUserSubscriptions(value["email"])

		tempMsg := map[string]interface{}{
			"about":         value["about"],
			"email":         value["email"],
			"followers":     listFollowers,
			"following":     listFollowing,
			"id":            respId,
			"isAnonymous":   respIsAnonymous,
			"name":          value["name"],
			"subscriptions": listSubscriptions,
			"username":      value["username"],
		}

		responseMsg = append(responseMsg, tempMsg)
	}

	resp, err = createResponseFromArray(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (u *User) listFollowers() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.follower WHERE followee = ?"

	resp := u.listBasic(query)
	return resp
}

func (u *User) listFollowing() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.followee WHERE follower = ?"

	resp := u.listBasic(query)
	return resp
}

func (u *User) listPosts() string {
	delete(u.inputRequest.query, "forum")

	p := Post{inputRequest: u.inputRequest, db: u.db}
	return p.list()
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
		if inputRequest.path == "/db/api/user/details/" {
			result = user.getDetails()
		} else if inputRequest.path == "/db/api/user/listFollowers/" {
			result = user.listFollowers()
		} else if inputRequest.path == "/db/api/user/listFollowing/" {
			result = user.listFollowing()
		} else if inputRequest.path == "/db/api/user/listPosts/" {
			result = user.listPosts()
		}
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
		fmt.Println(err)
		return resp
	}

	args = append(args, f.inputRequest.json["name"])
	args = append(args, f.inputRequest.json["short_name"])
	args = append(args, f.inputRequest.json["user"])

	dbResp, err := execQuery(query, &args, f.db)
	if err != nil {
		fmt.Println(err)
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
	var relatedUser bool
	if len(f.inputRequest.query["forum"]) != 1 {
		resp = createInvalidResponse()
		return resp
	}

	if len(f.inputRequest.query["related"]) == 1 && f.inputRequest.query["related"][0] == "user" {
		relatedUser = true
	}

	responseCode, responseMsg := f.getForumDetails()

	if relatedUser {
		u := User{inputRequest: f.inputRequest, db: f.db}
		u.inputRequest.query["user"] = append(u.inputRequest.query["user"], responseMsg["user"].(string))
		_, userDetails := u.getUserDetails()

		responseMsg["user"] = userDetails
	}

	resp, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

// DO IT
func (f *Forum) listThreads() string {

	return "DO IT"
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
// Thread handler here
// =================

type Thread struct {
	inputRequest *InputRequest
	db           *sql.DB
}

func (t *Thread) updateBoolBasic(query string, value bool) string {
	var resp string

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "thread")
	if err != nil {
		return resp
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args = append(args, value)
	args = append(args, t.inputRequest.json["thread"])

	dbResp, err := execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	if dbResp.rowCount == 0 {
		t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], intToString(int(threadId)))
		responseCode, responseMsg := t.getThreadDetails()

		if responseCode != 0 {
			resp = createNotExistResponse()
			return resp
		}

		resp, err = createResponse(responseCode, responseMsg)
		if err != nil {
			panic(err.Error())
		}
		return resp
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"thread": threadId,
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (t *Thread) close() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, true)

	return resp
}

func (t *Thread) create() string {
	var resp string
	query := "INSERT INTO thread (forum, title, isClosed, user, date, message, slug, isDeleted) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "forum", "title", "isClosed", "user", "date", "message", "slug")
	if err != nil {
		return resp
	}

	args = append(args, t.inputRequest.json["forum"])
	args = append(args, t.inputRequest.json["title"])
	args = append(args, t.inputRequest.json["isClosed"])
	args = append(args, t.inputRequest.json["user"])
	args = append(args, t.inputRequest.json["date"])
	args = append(args, t.inputRequest.json["message"])
	args = append(args, t.inputRequest.json["slug"])

	// Validate isDeleted param
	isDeleted := t.inputRequest.json["isDeleted"]
	if isDeleted == nil {
		args = append(args, false)
	} else {
		if isDeleted != false && isDeleted != true {
			resp = createInvalidResponse()
			return resp
		}
		args = append(args, isDeleted)
	}

	dbResp, err := execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	query = "SELECT * FROM thread WHERE id = ?"
	args = args[0:0]
	args = append(args, dbResp.lastId)
	newThread, err := selectQuery(query, &args, t.db)
	if err != nil {
		panic(err)
	}

	respIsClosed, _ := strconv.ParseBool(newThread.values[0]["isClosed"])
	respIsDeleted, _ := strconv.ParseBool(newThread.values[0]["isDeleted"])

	responseCode := 0
	responseMsg := map[string]interface{}{
		"forum":     newThread.values[0]["forum"],
		"title":     newThread.values[0]["title"],
		"id":        dbResp.lastId,
		"user":      newThread.values[0]["user"],
		"date":      newThread.values[0]["date"],
		"message":   newThread.values[0]["message"],
		"slug":      newThread.values[0]["slug"],
		"isClosed":  respIsClosed,
		"isDeleted": respIsDeleted,
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("thread.create()")

	return resp
}

func (t *Thread) getThreadDetails() (int, map[string]interface{}) {
	query := "SELECT t.*, COUNT(*) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.id = ?"
	var args []interface{}
	args = append(args, t.inputRequest.query["thread"][0])

	getThread, err := selectQuery(query, &args, t.db)
	if err != nil {
		panic(err)
	}

	if getThread.rows == 0 {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		return responseCode, errorMessage
	}

	respId, _ := strconv.ParseInt(getThread.values[0]["id"], 10, 64)
	respLikes, _ := strconv.ParseInt(getThread.values[0]["likes"], 10, 64)
	respDislikes, _ := strconv.ParseInt(getThread.values[0]["dislikes"], 10, 64)
	respPoints, _ := strconv.ParseInt(getThread.values[0]["points"], 10, 64)
	respPosts, _ := strconv.ParseInt(getThread.values[0]["posts"], 10, 64)
	respIsClosed, _ := strconv.ParseBool(getThread.values[0]["isClosed"])
	respIsDeleted, _ := strconv.ParseBool(getThread.values[0]["isDeleted"])

	responseCode := 0
	responseMsg := map[string]interface{}{
		"date":      getThread.values[0]["date"],
		"dislikes":  respDislikes,
		"forum":     getThread.values[0]["forum"],
		"id":        respId,
		"isClosed":  respIsClosed,
		"isDeleted": respIsDeleted,
		"likes":     respLikes,
		"message":   getThread.values[0]["message"],
		"points":    respPoints,
		"posts":     respPosts,
		"slug":      getThread.values[0]["slug"],
		"title":     getThread.values[0]["title"],
		"user":      getThread.values[0]["user"],
	}

	return responseCode, responseMsg
}

func (t *Thread) getArrayThreadsDetails(query string, args []interface{}) (int, []map[string]interface{}) {
	getThread, err := selectQuery(query, &args, t.db)
	if err != nil {
		panic(err)
	}

	if getThread.rows == 0 {
		var responseMsg []map[string]interface{}
		responseCode := 1
		errorMessage := map[string]interface{}{
			"msg": "Not found",
		}
		responseMsg = append(responseMsg, errorMessage)

		// fmt.Println(responseCode, errorMessage)
		return responseCode, responseMsg
	}

	responseCode := 0
	var responseMsg []map[string]interface{}

	for _, value := range getThread.values {
		respId, _ := strconv.ParseInt(value["id"], 10, 64)
		respLikes, _ := strconv.ParseInt(value["likes"], 10, 64)
		respDislikes, _ := strconv.ParseInt(value["dislikes"], 10, 64)
		respPoints, _ := strconv.ParseInt(value["points"], 10, 64)
		respPosts, _ := strconv.ParseInt(value["posts"], 10, 64)
		respIsClosed, _ := strconv.ParseBool(value["isClosed"])
		respIsDeleted, _ := strconv.ParseBool(value["isDeleted"])

		tempMsg := map[string]interface{}{
			"date":      value["date"],
			"dislikes":  respDislikes,
			"forum":     value["forum"],
			"id":        respId,
			"isClosed":  respIsClosed,
			"isDeleted": respIsDeleted,
			"likes":     respLikes,
			"message":   value["message"],
			"points":    respPoints,
			"posts":     respPosts,
			"slug":      value["slug"],
			"title":     value["title"],
			"user":      value["user"],
		}

		responseMsg = append(responseMsg, tempMsg)
	}

	return responseCode, responseMsg
}

func (t *Thread) details() string {
	var resp string
	var relatedUser, relatedForum bool
	if len(t.inputRequest.query["thread"]) != 1 {
		resp = createInvalidResponse()
		return resp
	}

	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("user", t.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("forum", t.inputRequest.query["related"]) {
		relatedForum = true
	}

	responseCode, responseMsg := t.getThreadDetails()

	if relatedUser {
		u := User{inputRequest: t.inputRequest, db: t.db}
		u.inputRequest.query["user"] = append(u.inputRequest.query["user"], responseMsg["user"].(string))
		_, userDetails := u.getUserDetails()

		responseMsg["user"] = userDetails
	}

	if relatedForum {
		f := Forum{inputRequest: t.inputRequest, db: t.db}
		f.inputRequest.query["user"] = f.inputRequest.query["user"][:0]
		f.inputRequest.query["forum"] = append(f.inputRequest.query["forum"], responseMsg["forum"].(string))
		_, forumDetails := f.getForumDetails()

		responseMsg["forum"] = forumDetails
	}

	resp, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (t *Thread) list() string {
	var resp, query string
	f := false
	var args []interface{}

	// Validate query values
	if len(t.inputRequest.query["user"]) == 1 {
		query = "SELECT t.*, COUNT(*) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.user = ?"
		args = append(args, t.inputRequest.query["user"][0])

		f = true
	}
	if len(t.inputRequest.query["forum"]) == 1 && f == false {
		query = "SELECT t.*, COUNT(*) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.forum = ?"
		args = append(args, t.inputRequest.query["forum"][0])

		f = true
	}
	if f == false {
		resp = createInvalidResponse()
		return resp
	}

	// Check and validate optional params
	if len(t.inputRequest.query["since"]) >= 1 {
		query += " AND t.date > ?"
		args = append(args, t.inputRequest.query["since"][0])
	}

	query = query + " GROUP BY t.id"
	if len(t.inputRequest.query["order"]) >= 1 {
		orderType := t.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			resp = createInvalidResponse()
			return resp
		}

		query += fmt.Sprintf(" ORDER BY t.date %s", orderType)
	}
	if len(t.inputRequest.query["limit"]) >= 1 {
		limitValue := t.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			resp = createInvalidResponse()
			return resp
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Response here
	responseCode, responseMsg := t.getArrayThreadsDetails(query, args)
	if responseCode != 0 {
		resp, _ = createResponse(responseCode, responseMsg[0])
		return resp
	}
	resp, err := createResponseFromArray(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

// DO IT
func (t *Thread) listPosts() string {
	/*
		var query, order, resp string
		f := false
		var args []interface{}

		// Validate query values
		if len(t.inputRequest.query["thread"]) == 1 {
			query = "SELECT * FROM post WHERE thread = ?"
			args = append(args, t.inputRequest.query["thread"][0])

			f = true
		}
		if f == false {
			resp = createInvalidResponse()
			return resp
		}

		// order by here
		if len(t.inputRequest.query["order"]) >= 1 {
			orderType := t.inputRequest.query["order"][0]
			if orderType != "desc" && orderType != "asc" {
				resp = createInvalidResponse()
				return resp
			}

			order = fmt.Sprintf(" ORDER BY date %s", orderType)
		} else {
			order = " ORDER BY date DESC"
		}

		responseCode, responseMsg := p.getList(query, order, args)

		// check responseCode
		if responseCode == 100500 {
			resp = createInvalidResponse()
			return resp
		} else if responseCode == 1 {
			resp, _ = createResponse(responseCode, responseMsg[0])
			return resp
		} else if responseCode == 0 {
			resp, err := createResponseFromArray(responseCode, responseMsg)
			if err != nil {
				panic(err.Error())
			}
			return resp
		} else {
			resp = createInvalidResponse()
			return resp
		}
	*/

	return "DO IT"
}

func (t *Thread) open() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, false)

	return resp
}

func (t *Thread) remove() string {
	query := "UPDATE thread SET isDeleted = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, true)

	return resp
}

func (t *Thread) restore() string {
	query := "UPDATE thread SET isDeleted = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, false)

	return resp
}

func (t *Thread) subscribe() string {
	var resp string
	query := "INSERT INTO subscribe (thread, user) VALUES(?, ?)"

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "thread", "user")
	if err != nil {
		return resp
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	args = append(args, t.inputRequest.json["thread"])
	args = append(args, t.inputRequest.json["user"])

	_, err = execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)

		// return exist
		if checkError1062(err) == true {
			for k := range t.inputRequest.query {
				delete(t.inputRequest.query, k)
			}
			t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], floatToString(t.inputRequest.json["thread"].(float64)))

			resp = t.details()
			return resp
		}

		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	// else return info

	responseCode := 0
	responseMsg := map[string]interface{}{
		"thread": int(t.inputRequest.json["thread"].(float64)),
		"user":   t.inputRequest.json["user"],
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("thread.subscribe()")

	return resp
}

func (t *Thread) unsubscribe() string {
	var resp string
	query := "DELETE FROM subscribe WHERE thread = ? AND user = ?"

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "thread", "user")
	if err != nil {
		return resp
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	args = append(args, t.inputRequest.json["thread"])
	args = append(args, t.inputRequest.json["user"])

	dbResp, err := execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	if dbResp.rowCount == 0 {
		for k := range t.inputRequest.query {
			delete(t.inputRequest.query, k)
		}
		t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], floatToString(t.inputRequest.json["thread"].(float64)))

		resp = t.details()
		return resp
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"thread": int(t.inputRequest.json["thread"].(float64)),
		"user":   t.inputRequest.json["user"],
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (t *Thread) update() string {
	var resp string
	query := "UPDATE thread SET message = ?, slug = ? WHERE id = ?"

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "thread", "message", "slug")
	if err != nil {
		return resp
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args = append(args, t.inputRequest.json["message"])
	args = append(args, t.inputRequest.json["slug"])
	args = append(args, t.inputRequest.json["thread"])

	_, err = execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], intToString(int(threadId)))
	responseCode, responseMsg := t.getThreadDetails()

	if responseCode != 0 {
		resp = createNotExistResponse()
		return resp
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func (t *Thread) vote() string {
	var resp string
	var query string

	var args []interface{}

	resp, err := validateJson(t.inputRequest, "thread", "vote")
	if err != nil {
		return resp
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false || checkFloat64Type(t.inputRequest.json["vote"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	threadId := t.inputRequest.json["thread"].(float64)
	vote := t.inputRequest.json["vote"].(float64)

	if vote == 1 {
		query = "UPDATE thread SET likes = likes + 1, points = points + 1 WHERE id = ?"
	} else if vote == -1 {
		query = "UPDATE thread SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?"
	} else {
		resp = createInvalidResponse()
		return resp
	}

	args = append(args, t.inputRequest.json["thread"])

	_, err = execQuery(query, &args, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], intToString(int(threadId)))
	responseCode, responseMsg := t.getThreadDetails()

	if responseCode != 0 {
		resp = createNotExistResponse()
		return resp
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func threadHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	thread := Thread{inputRequest: inputRequest, db: db}

	var result string

	if inputRequest.method == "GET" {

		// Like Router
		if inputRequest.path == "/db/api/thread/details/" {
			result = thread.details()
		} else if inputRequest.path == "/db/api/thread/list/" {
			result = thread.list()
		}
	} else if inputRequest.method == "POST" {

		// Like Router
		if inputRequest.path == "/db/api/thread/create/" {
			result = thread.create()
		} else if inputRequest.path == "/db/api/thread/close/" {
			result = thread.close()
		} else if inputRequest.path == "/db/api/thread/restore/" {
			result = thread.restore()
		} else if inputRequest.path == "/db/api/thread/vote/" {
			result = thread.vote()
		} else if inputRequest.path == "/db/api/thread/remove/" {
			result = thread.remove()
		} else if inputRequest.path == "/db/api/thread/open/" {
			result = thread.open()
		} else if inputRequest.path == "/db/api/thread/update/" {
			result = thread.update()
		} else if inputRequest.path == "/db/api/thread/subscribe/" {
			result = thread.subscribe()
		} else if inputRequest.path == "/db/api/thread/unsubscribe/" {
			result = thread.unsubscribe()
		}
		// else if inputRequest.path == "/db/api/user/unfollow/" {
		// 	result = user.unfollow()
		// } else if inputRequest.path == "/db/api/user/updateProfile/" {
		// 	result = user.updateProfile()
		// }
	}

	io.WriteString(w, result)
}

// =================
// Post handler here
// =================
type Post struct {
	inputRequest *InputRequest
	db           *sql.DB
}

func (p *Post) create() string {
	var resp string
	query := "INSERT INTO post (date, thread, message, user, forum, isApproved, isHighlighted, isEdited, isSpam, isDeleted, parent) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// 	ir.url = fmt.Sprintf("%v", r.URL)

	var args []interface{}

	resp, err := validateJson(p.inputRequest, "date", "thread", "message", "user", "forum")
	if err != nil {
		return resp
	}

	args = append(args, p.inputRequest.json["date"])
	args = append(args, p.inputRequest.json["thread"])
	args = append(args, p.inputRequest.json["message"])
	args = append(args, p.inputRequest.json["user"])
	args = append(args, p.inputRequest.json["forum"])

	result, err := validateBoolParams(p.inputRequest.json, &args, "isApproved", "isHighlighted", "isEdited", "isSpam", "isDeleted")
	if err != nil {
		return result
	}

	boolParent := false

	// parent here
	if p.inputRequest.json["parent"] != nil {
		if checkFloat64Type(p.inputRequest.json["parent"]) == false {
			return createInvalidResponse()
		}
		//
		// search parent here
		//

		parent := p.inputRequest.json["parent"].(float64)

		// find parent and last chils
		parentQuery := "SELECT id, parent FROM post WHERE id = ? && isDeleted = false"
		var parentArgs []interface{}

		parentArgs = append(parentArgs, parent)

		getThread, err := selectQuery(parentQuery, &parentArgs, p.db)
		if err != nil {
			panic(err)
		}

		// check query
		if getThread.rows == 0 {
			return createNotExistResponse()
		}

		parentArgs = parentArgs[0:0] // clear args

		//
		// search place for child
		//

		var child string
		getParent := getThread.values[0]["parent"]

		parentQuery = "SELECT parent FROM post WHERE parent LIKE ?"
		parentArgs = append(parentArgs, getParent+"%")

		getThread, err = selectQuery(parentQuery, &parentArgs, p.db)
		if err != nil {
			panic(err)
		}

		// because the first element is parent
		if getThread.rows == 1 {
			newParent := getParent
			newChild := toBase95(1)

			child = newParent + newChild
		} else {
			lastChild := getThread.values[getThread.rows-1]["parent"]
			newParent := getParent
			oldChild := fromBase95(lastChild[len(lastChild)-5:])

			oldChild++

			newChild := toBase95(oldChild)

			child = newParent + newChild
		}

		args = append(args, child)
	} else {
		args = append(args, nil)
		boolParent = true
	}

	dbResp, err := execQuery(query, &args, p.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	if boolParent {
		boolParentQuery := "UPDATE post SET parent = ? WHERE id = ?"
		var boolParentArgs []interface{}

		parent := toBase95(int(dbResp.lastId))
		boolParentArgs = append(boolParentArgs, parent)
		boolParentArgs = append(boolParentArgs, dbResp.lastId)

		_, _ = execQuery(boolParentQuery, &boolParentArgs, p.db)
	}

	tempCounter := 5
	responseCode := 0
	responseMsg := map[string]interface{}{
		"date":          p.inputRequest.json["date"],
		"forum":         p.inputRequest.json["forum"],
		"id":            dbResp.lastId,
		"isApproved":    args[incInt(&tempCounter)],
		"isHighlighted": args[tempCounter+0],
		"isEdited":      args[tempCounter+1],
		"isSpam":        args[tempCounter+2],
		"isDeleted":     args[tempCounter+3],
		"message":       p.inputRequest.json["message"],
		"parent":        args[tempCounter+4],
		"thread":        p.inputRequest.json["thread"].(float64),
		"user":          p.inputRequest.json["user"],
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("post.create()")

	return resp
}

func (p *Post) getParentId(id int64, path string) int {
	if len(path) == 5 {
		return int(id)
	} else if len(path) == 10 {
		return fromBase95(path[len(path)-10 : len(path)-5])
	} else {
		parentId := path[:len(path)-5]

		query := "SELECT id FROM post WHERE parent = ?"
		var args []interface{}
		args = append(args, parentId)

		getParent, _ := selectQuery(query, &args, p.db)

		respId, _ := strconv.ParseInt(getParent.values[0]["id"], 10, 64)

		return int(respId)
	}
}

func (p *Post) updateBoolBasic(query string, value bool) string {
	var resp string

	var args []interface{}

	resp, err := validateJson(p.inputRequest, "post")
	if err != nil {
		return resp
	}

	if checkFloat64Type(p.inputRequest.json["post"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	postId := p.inputRequest.json["post"].(float64)

	args = append(args, value)
	args = append(args, p.inputRequest.json["post"])

	dbResp, err := execQuery(query, &args, p.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	if dbResp.rowCount == 0 {
		p.inputRequest.query["post"] = append(p.inputRequest.query["post"], intToString(int(postId)))
		responseCode, responseMsg := p.getPostDetails()

		if responseCode != 0 {
			resp = createNotExistResponse()
			return resp
		}

		resp, err = createResponse(responseCode, responseMsg)
		if err != nil {
			panic(err.Error())
		}
		return resp
	}

	responseCode := 0
	responseMsg := map[string]interface{}{
		"post": postId,
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (p *Post) getPostDetails() (int, map[string]interface{}) {
	query := "SELECT * FROM post WHERE id = ?"
	var args []interface{}
	args = append(args, p.inputRequest.query["post"][0])

	getPost, err := selectQuery(query, &args, p.db)
	if err != nil {
		panic(err)
	}

	if getPost.rows == 0 {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		return responseCode, errorMessage
	}

	respId, _ := strconv.ParseInt(getPost.values[0]["id"], 10, 64)
	respLikes, _ := strconv.ParseInt(getPost.values[0]["likes"], 10, 64)
	respDislikes, _ := strconv.ParseInt(getPost.values[0]["dislikes"], 10, 64)
	respPoints, _ := strconv.ParseInt(getPost.values[0]["points"], 10, 64)
	respThread, _ := strconv.ParseInt(getPost.values[0]["thread"], 10, 64)
	respIsApproved, _ := strconv.ParseBool(getPost.values[0]["isApproved"])
	respIsDeleted, _ := strconv.ParseBool(getPost.values[0]["isDeleted"])
	respIsEdited, _ := strconv.ParseBool(getPost.values[0]["isEdited"])
	respIsHighlighted, _ := strconv.ParseBool(getPost.values[0]["isHighlighted"])
	respIsSpam, _ := strconv.ParseBool(getPost.values[0]["isSpam"])

	responseCode := 0
	responseMsg := map[string]interface{}{
		"date":          getPost.values[0]["date"],
		"dislikes":      respDislikes,
		"forum":         getPost.values[0]["forum"],
		"id":            respId,
		"isApproved":    respIsApproved,
		"isDeleted":     respIsDeleted,
		"isEdited":      respIsEdited,
		"isHighlighted": respIsHighlighted,
		"isSpam":        respIsSpam,
		"likes":         respLikes,
		"message":       getPost.values[0]["message"],
		"parent":        nil,
		"points":        respPoints,
		"thread":        respThread,
		"user":          getPost.values[0]["user"],
	}

	parent := p.getParentId(respId, getPost.values[0]["parent"])
	if parent == int(respId) {
		responseMsg["parent"] = nil
	} else {
		responseMsg["parent"] = parent
	}

	return responseCode, responseMsg
}

func (p *Post) details() string {
	var resp string
	var relatedUser, relatedThread, relatedForum bool
	if len(p.inputRequest.query["post"]) != 1 {
		resp = createInvalidResponse()
		return resp
	}

	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("user", p.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("thread", p.inputRequest.query["related"]) {
		relatedThread = true
	}
	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("forum", p.inputRequest.query["related"]) {
		relatedForum = true
	}

	responseCode, responseMsg := p.getPostDetails()

	if relatedUser {
		u := User{inputRequest: p.inputRequest, db: p.db}
		u.inputRequest.query["user"] = append(u.inputRequest.query["user"], responseMsg["user"].(string))
		_, userDetails := u.getUserDetails()

		responseMsg["user"] = userDetails
	}

	if relatedThread {
		t := Thread{inputRequest: p.inputRequest, db: p.db}
		t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], int64ToString(responseMsg["thread"].(int64)))
		_, threadDetails := t.getThreadDetails()

		responseMsg["thread"] = threadDetails
	}

	if relatedForum {
		f := Forum{inputRequest: p.inputRequest, db: p.db}
		f.inputRequest.query["user"] = f.inputRequest.query["user"][:0]
		f.inputRequest.query["forum"] = append(f.inputRequest.query["forum"], responseMsg["forum"].(string))
		_, forumDetails := f.getForumDetails()

		responseMsg["forum"] = forumDetails
	}

	resp, err := createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (p *Post) getArrayPostDetails(query string, args []interface{}) (int, []map[string]interface{}) {
	getPost, err := selectQuery(query, &args, p.db)
	if err != nil {
		panic(err)
	}

	if getPost.rows == 0 {
		var responseMsg []map[string]interface{}
		responseCode := 1
		errorMessage := map[string]interface{}{
			"msg": "Not found",
		}
		responseMsg = append(responseMsg, errorMessage)

		// fmt.Println(responseCode, errorMessage)
		return responseCode, responseMsg
	}

	responseCode := 0
	var responseMsg []map[string]interface{}

	for _, value := range getPost.values {
		respId, _ := strconv.ParseInt(value["id"], 10, 64)
		respLikes, _ := strconv.ParseInt(value["likes"], 10, 64)
		respDislikes, _ := strconv.ParseInt(value["dislikes"], 10, 64)
		respPoints, _ := strconv.ParseInt(value["points"], 10, 64)
		respThread, _ := strconv.ParseInt(value["thread"], 10, 64)
		respIsApproved, _ := strconv.ParseBool(value["isApproved"])
		respIsDeleted, _ := strconv.ParseBool(value["isDeleted"])
		respIsEdited, _ := strconv.ParseBool(value["isEdited"])
		respIsHighlighted, _ := strconv.ParseBool(value["isHighlighted"])
		respIsSpam, _ := strconv.ParseBool(value["isSpam"])

		tempMsg := map[string]interface{}{
			"date":          value["date"],
			"dislikes":      respDislikes,
			"forum":         value["forum"],
			"id":            respId,
			"isApproved":    respIsApproved,
			"isDeleted":     respIsDeleted,
			"isEdited":      respIsEdited,
			"isHighlighted": respIsHighlighted,
			"isSpam":        respIsSpam,
			"likes":         respLikes,
			"message":       value["message"],
			"parent":        nil,
			"points":        respPoints,
			"thread":        respThread,
			"user":          value["user"],
		}

		parent := p.getParentId(respId, value["parent"])
		if parent == int(respId) {
			tempMsg["parent"] = nil
		} else {
			tempMsg["parent"] = parent
		}

		responseMsg = append(responseMsg, tempMsg)
	}

	return responseCode, responseMsg
}

func (p *Post) getList(query string, order string, args []interface{}) (int, []map[string]interface{}) {
	// Check and validate optional params
	if len(p.inputRequest.query["since"]) >= 1 {
		query += " AND date > ?"
		args = append(args, p.inputRequest.query["since"][0])
	}

	query += order

	if len(p.inputRequest.query["limit"]) >= 1 {
		limitValue := p.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			return 100500, nil
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Response here
	responseCode, responseMsg := p.getArrayPostDetails(query, args)
	if responseCode != 0 {
		return responseCode, responseMsg
	}
	return responseCode, responseMsg
}

func (p *Post) list() string {
	var query, order, resp string
	f := false
	var args []interface{}

	// Validate query values
	if len(p.inputRequest.query["user"]) == 1 {
		query = "SELECT * FROM post WHERE user = ?"
		args = append(args, p.inputRequest.query["user"][0])

		f = true
	}
	if len(p.inputRequest.query["forum"]) == 1 && f == false {
		query = "SELECT * FROM post WHERE forum = ?"
		args = append(args, p.inputRequest.query["forum"][0])

		f = true
	}
	if f == false {
		resp = createInvalidResponse()
		return resp
	}

	// order by here
	if len(p.inputRequest.query["order"]) >= 1 {
		orderType := p.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			resp = createInvalidResponse()
			return resp
		}

		order = fmt.Sprintf(" ORDER BY date %s", orderType)
	} else {
		order = " ORDER BY date DESC"
	}

	responseCode, responseMsg := p.getList(query, order, args)

	// check responseCode
	if responseCode == 100500 {
		resp = createInvalidResponse()
		return resp
	} else if responseCode == 1 {
		resp, _ = createResponse(responseCode, responseMsg[0])
		return resp
	} else if responseCode == 0 {
		resp, err := createResponseFromArray(responseCode, responseMsg)
		if err != nil {
			panic(err.Error())
		}
		return resp
	} else {
		resp = createInvalidResponse()
		return resp
	}
}

func (p *Post) remove() string {
	query := "UPDATE post SET isDeleted = ? WHERE id = ?"

	resp := p.updateBoolBasic(query, true)

	return resp
}

func (p *Post) restore() string {
	query := "UPDATE post SET isDeleted = ? WHERE id = ?"

	resp := p.updateBoolBasic(query, false)

	return resp
}

func (p *Post) update() string {
	var resp string
	query := "UPDATE post SET message = ?, isEdited = true WHERE id = ?"

	var args []interface{}

	resp, err := validateJson(p.inputRequest, "post", "message")
	if err != nil {
		return resp
	}

	if checkFloat64Type(p.inputRequest.json["post"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	threadId := p.inputRequest.json["post"].(float64)

	args = append(args, p.inputRequest.json["message"])
	args = append(args, p.inputRequest.json["post"])

	_, err = execQuery(query, &args, p.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	p.inputRequest.query["post"] = append(p.inputRequest.query["post"], intToString(int(threadId)))
	responseCode, responseMsg := p.getPostDetails()

	if responseCode != 0 {
		resp = createNotExistResponse()
		return resp
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func (p *Post) vote() string {
	var resp string
	var query string

	var args []interface{}

	resp, err := validateJson(p.inputRequest, "post", "vote")
	if err != nil {
		return resp
	}

	if checkFloat64Type(p.inputRequest.json["post"]) == false || checkFloat64Type(p.inputRequest.json["vote"]) == false {
		resp = createInvalidResponse()
		return resp
	}

	postId := p.inputRequest.json["post"].(float64)
	vote := p.inputRequest.json["vote"].(float64)

	if vote == 1 {
		query = "UPDATE post SET likes = likes + 1, points = points + 1 WHERE id = ?"
	} else if vote == -1 {
		query = "UPDATE post SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?"
	} else {
		resp = createInvalidResponse()
		return resp
	}

	args = append(args, p.inputRequest.json["post"])

	_, err = execQuery(query, &args, p.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			panic(err)
		}

		return resp
	}

	p.inputRequest.query["post"] = append(p.inputRequest.query["post"], intToString(int(postId)))
	responseCode, responseMsg := p.getPostDetails()

	if responseCode != 0 {
		resp = createNotExistResponse()
		return resp
	}

	resp, err = createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func postHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	post := Post{inputRequest: inputRequest, db: db}

	var result string

	if inputRequest.method == "GET" {
		// Like Router
		if inputRequest.path == "/db/api/post/details/" {
			result = post.details()
		} else if inputRequest.path == "/db/api/post/list/" {
			result = post.list()
		}
	} else if inputRequest.method == "POST" {

		// Like Router
		if inputRequest.path == "/db/api/post/create/" {
			result = post.create()
		} else if inputRequest.path == "/db/api/post/restore/" {
			result = post.restore()
		} else if inputRequest.path == "/db/api/post/vote/" {
			result = post.vote()
		} else if inputRequest.path == "/db/api/post/remove/" {
			result = post.remove()
		} else if inputRequest.path == "/db/api/post/update/" {
			result = post.update()
		}
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
	http.HandleFunc("/db/api/thread/", makeHandler(db, threadHandler))
	http.HandleFunc("/db/api/post/", makeHandler(db, postHandler))

	http.ListenAndServe(PORT, nil)
}

// =================
// Utils here
// =================

func floatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func intToString(inputNum int) string {
	// to convert a float number to a string
	return strconv.Itoa(inputNum)
}

func int64ToString(inputNum int64) string {
	// to convert a float number to a string
	return strconv.FormatInt(inputNum, 10)
}

func checkFloat64Type(inputNum interface{}) bool {
	if reflect.TypeOf(inputNum).Kind() != reflect.Float64 {
		return false
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringToInt(inputStr string) int {
	result, err := strconv.Atoi(inputStr)
	if err != nil {
		panic(err)
	}
	return result
}

func incInt(value *int) int {
	*value += 1
	return *value
}

func toBase95(value int) string {
	BASE95 := ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`
	length := 5
	base := len(BASE95)
	result := make([]byte, 5)

	for i := 0; i < length; i++ {
		result[i] = BASE95[0]
	}

	if value == 0 {
		return string(result)
	}

	counter := 0
	for value != 0 {
		mod := value % base

		result[length-counter-1] = BASE95[mod]
		counter++

		value = value / base
	}

	return string(result)
}

func fromBase95(value string) int {
	BASE95 := ` !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`
	length := 5
	base := len(BASE95)
	var result int

	counter := 0
	step := 1
	for i := 0; i < length; i++ {
		index := strings.Index(BASE95, string(value[length-counter-1]))
		counter++

		result = result + index*step

		step = step * base
	}

	return result
}

// =================
// Future here
// =================

func testArrayJson(code int, response []response.Foo) string {
	cacheContent := map[string]interface{}{
		"code":     code,
		"response": response,
	}
	result, _ := json.Marshal(cacheContent)

	return string(result)
}

func testJson(code int, response response.Foo) string {
	cacheContent := map[string]interface{}{
		"code":     code,
		"response": response,
	}
	result, _ := json.Marshal(cacheContent)

	return string(result)
}
