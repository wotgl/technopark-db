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
	"os"
	"reflect"
	"strconv"
	"strings"
	rs "technopark-db/response"

	mysql "github.com/go-sql-driver/mysql"
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
		log.Panic(err)
	}

	var parsed map[string]interface{}
	json.Unmarshal([]byte(body), &parsed)
	ir.json = parsed

	// GET Query
	ir.query = r.URL.Query()
}

func createResponse(code int, response map[string]interface{}) (string, error) {
	content := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(content)
	if err != nil {
		log.Println("Error encoding JSON")
		return "", err
	}

	return string(str), nil
}

func _createResponse(code int, response rs.RespStruct) (string, error) {
	content := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(content)
	if err != nil {
		log.Println("Error encoding JSON")
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
		log.Println("Error encoding JSON")
		return "", err
	}

	return string(str), nil
}

func _createResponseFromArray(code int, response []interface{}) (string, error) {
	cacheContent := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(cacheContent)
	if err != nil {
		log.Println("Error encoding JSON")
		return "", err
	}

	return string(str), nil
}

func errorExecParse(err error) (int, map[string]interface{}) {
	log.Println("Error:\t", err)
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
			// fmt.Println("errorExecParse() default")
			panic(err.Error())
			responseCode = 4
			errorMessage = map[string]interface{}{"msg": "Unknown Error"}
		}

		return responseCode, errorMessage
	}
	panic(err.Error()) // proper error handling instead of panic in your app
}

func createNotFoundForArray() string {
	var resp string
	responseCode := 1
	errorMessage := map[string]interface{}{
		"msg": "Not found",
	}

	var responseMsg []map[string]interface{}

	responseMsg = append(responseMsg, errorMessage)

	resp, _ = createResponseFromArray(responseCode, responseMsg)
	return resp
}

func createInvalidQuery() string {
	responseCode := 3
	errorMessage := map[string]interface{}{"msg": "Invalid query"}

	resp, err := createResponse(responseCode, errorMessage)
	if err != nil {
		log.Panic(err)
	}

	return resp
}

func createInvalidJsonResponse(json *map[string]interface{}) string {
	responseCode := 3
	errorMessage := map[string]interface{}{"msg": "Invalid json"}

	resp, err := createResponse(responseCode, errorMessage)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Invalid JSON\t", json)

	return resp
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
		log.Panic(err)
	}

	return resp
}

func createNotExistResponse() string {
	responseCode := 1
	errorMessage := &rs.ErrorMsg{
		Msg: "Not exist",
	}

	resp, err := _createResponse(responseCode, errorMessage)
	if err != nil {
		log.Panic(err)
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

func clearQuery(query *map[string][]string) {
	for k := range *query {
		delete(*query, k)
	}
}

type Args struct {
	data []interface{}
}

func (args *Args) generateFromJson(json *map[string]interface{}, params ...string) {
	for _, value := range params {
		args.data = append(args.data, (*json)[value])
	}
}

func (args *Args) addFromJson(json *map[string]interface{}, params ...string) {
	for _, value := range params {
		args.data = append(args.data, (*json)[value])
	}
}

func (args *Args) append(newData ...interface{}) {
	for data := range newData {
		args.data = append(args.data, newData[data])
	}
}

func (args *Args) clear() {
	args.data = args.data[0:0]
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

// +
func (u *User) create() string {
	var resp string
	args := Args{}

	query := "INSERT INTO user (username, about, name, email, isAnonymous) VALUES(?, ?, ?, ?, ?)"

	resp, err := validateJson(u.inputRequest, "email")
	if err != nil {
		return createInvalidJsonResponse(&u.inputRequest.json)
	}

	args.generateFromJson(&u.inputRequest.json, "username", "about", "name", "email")

	// Validate isAnonymous param
	isAnonymous := u.inputRequest.json["isAnonymous"]
	if isAnonymous == nil {
		args.append(false)
	} else {
		if isAnonymous != false && isAnonymous != true {
			return createInvalidResponse()
		}
		args.addFromJson(&u.inputRequest.json, "isAnonymous")
	}

	dbResp, err := execQuery(query, &args.data, u.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	query = "SELECT * FROM user WHERE id = ?"
	args.clear()
	args.append(dbResp.lastId)
	newUser, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	responseCode := 0
	responseMsg := &rs.UserCreate{
		About:       newUser.values[0]["about"],
		Email:       newUser.values[0]["email"],
		Id:          dbResp.lastId,
		IsAnonymous: stringToBool(newUser.values[0]["isAnonymous"]),
		Name:        newUser.values[0]["name"],
		Username:    newUser.values[0]["username"],
	}

	resp, err = _createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("User '%s' created", responseMsg.Email)

	return resp
}

// +
func (u *User) _getUserDetails(args Args) (int, *rs.UserDetails) {
	query := "SELECT * FROM user WHERE email = ?"

	getUser, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	if getUser.rows == 0 {
		responseCode := 1
		errorMessage := &rs.UserDetails{}

		return responseCode, errorMessage
	}

	// followers here
	listFollowers := u.getUserFollowers(args.data[0].(string))

	// following here
	listFollowing := u.getUserFollowing(args.data[0].(string))

	// subscriptions here
	listSubscriptions := u.getUserSubscriptions(args.data[0].(string))

	respAbout := getUser.values[0]["about"]
	respName := getUser.values[0]["name"]
	respUsername := getUser.values[0]["username"]

	responseCode := 0
	responseMsg := &rs.UserDetails{
		About:         &respAbout,
		Email:         getUser.values[0]["email"],
		Followers:     listFollowers,
		Following:     listFollowing,
		Id:            stringToInt64(getUser.values[0]["id"]),
		IsAnonymous:   stringToBool(getUser.values[0]["isAnonymous"]),
		Name:          &respName,
		Subscriptions: listSubscriptions,
		Username:      &respUsername,
	}

	if respAbout == "NULL" {
		responseMsg.About = nil
	}
	if respName == "NULL" {
		responseMsg.Name = nil
	}
	if respUsername == "NULL" {
		responseMsg.Username = nil
	}

	return responseCode, responseMsg
}

// -
func (u *User) getUserDetails() (int, map[string]interface{}) {
	query := "SELECT * FROM user WHERE email = ?"
	var args []interface{}
	args = append(args, u.inputRequest.query["user"][0])

	getUser, err := selectQuery(query, &args, u.db)
	if err != nil {
		log.Panic(err)
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
	respAbout := getUser.values[0]["about"]
	respName := getUser.values[0]["name"]
	respUsername := getUser.values[0]["username"]

	responseCode := 0
	responseMsg := map[string]interface{}{
		"about":         respAbout,
		"email":         getUser.values[0]["email"],
		"followers":     listFollowers,
		"following":     listFollowing,
		"id":            respId,
		"isAnonymous":   respIsAnonymous,
		"name":          respName,
		"subscriptions": listSubscriptions,
		"username":      respUsername,
	}

	if respAbout == "NULL" {
		responseMsg["about"] = nil
	}
	if respName == "NULL" {
		responseMsg["name"] = nil
	}
	if respUsername == "NULL" {
		responseMsg["username"] = nil
	}

	return responseCode, responseMsg
}

// +
func (u *User) getUserFollowers(followee string) []string {
	args := Args{}
	args.append(followee)

	// followers here
	query := "SELECT follower FROM follow WHERE followee = ?"
	getUserFollowers, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	listFollowers := make([]string, 0)
	for _, value := range getUserFollowers.values {
		listFollowers = append(listFollowers, value["follower"])
	}

	return listFollowers
}

// +
func (u *User) getUserFollowing(follower string) []string {
	args := Args{}
	args.append(follower)

	// following here
	query := "SELECT followee FROM follow WHERE follower = ?"
	getUserFollowing, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	listFollowing := make([]string, 0)
	for _, value := range getUserFollowing.values {
		listFollowing = append(listFollowing, value["followee"])
	}

	return listFollowing
}

// +
func (u *User) getUserSubscriptions(user string) []int {
	args := Args{}
	args.append(user)

	// subscriptions here
	query := "SELECT thread FROM subscribe WHERE user = ? ORDER BY thread asc"
	getUserSubscriptions, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	listSubscriptions := make([]int, 0)
	for _, value := range getUserSubscriptions.values {
		listSubscriptions = append(listSubscriptions, stringToInt(value["thread"]))
	}

	return listSubscriptions
}

// +
func (u *User) getDetails() string {
	var resp string
	if len(u.inputRequest.query["user"]) != 1 {
		return createInvalidResponse()
	}

	args := Args{}
	args.append(u.inputRequest.query["user"][0])

	responseCode, responseMsg := u._getUserDetails(args)
	if responseCode == 1 {
		return createNotExistResponse()
	}

	resp, err := _createResponse(responseCode, responseMsg)
	if err != nil {
		log.Panic(err)
	}

	// log.Printf("User '%s' get details", responseMsg)

	return resp
}

// +
func (u *User) follow() string {
	var resp string
	query := "INSERT INTO follow (follower, followee) VALUES(?, ?)"

	args := Args{}

	resp, err := validateJson(u.inputRequest, "follower", "followee")
	if err != nil {
		return createInvalidJsonResponse(&u.inputRequest.json)
	}

	args.generateFromJson(&u.inputRequest.json, "follower", "followee")

	_, err = execQuery(query, &args.data, u.db)
	if err != nil {
		// return exist
		if checkError1062(err) == true {
			clearQuery(&u.inputRequest.query)
			u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))

			return u.getDetails()
		}
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	return u.getDetails()
}

// +
func (u *User) listBasic(query string) string {
	var resp string
	args := Args{}

	// Validate query values
	if len(u.inputRequest.query["user"]) == 1 {
		args.append(u.inputRequest.query["user"][0])
	} else {
		return createInvalidResponse()
	}

	// Check and validate optional params
	if len(u.inputRequest.query["since_id"]) >= 1 {
		query += " AND id >= ?"
		args.append(u.inputRequest.query["since_id"][0])
	}
	if len(u.inputRequest.query["order"]) >= 1 {
		orderType := u.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			return createInvalidResponse()
		}

		query += fmt.Sprintf(" ORDER BY date %s", orderType)
	}
	if len(u.inputRequest.query["limit"]) >= 1 {
		limitValue := u.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			return createInvalidResponse()
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Prepare users
	getUserFollowers, err := selectQuery(query, &args.data, u.db)
	if err != nil {
		log.Panic(err)
	}

	responseCode := 0
	responseArray := make([]rs.UserDetails, 0)
	responseMsg := &rs.UserListBasic{Users: responseArray}

	for _, value := range getUserFollowers.values {
		// followers here
		listFollowers := u.getUserFollowers(value["email"])

		// following here
		listFollowing := u.getUserFollowing(value["email"])

		// subscriptions here
		listSubscriptions := u.getUserSubscriptions(value["email"])

		respAbout := value["about"]
		respName := value["name"]
		respUsername := value["username"]

		tempUser := &rs.UserDetails{
			About:         &respAbout,
			Email:         value["email"],
			Followers:     listFollowers,
			Following:     listFollowing,
			Id:            stringToInt64(value["id"]),
			IsAnonymous:   stringToBool(value["isAnonymous"]),
			Name:          &respName,
			Subscriptions: listSubscriptions,
			Username:      &respUsername,
		}

		if respAbout == "NULL" {
			tempUser.About = nil
		}
		if respName == "NULL" {
			tempUser.Name = nil
		}
		if respUsername == "NULL" {
			tempUser.Username = nil
		}

		responseMsg.Users = append(responseMsg.Users, *tempUser)
	}

	responseInterface := make([]interface{}, len(responseMsg.Users))
	for i, v := range responseMsg.Users {
		responseInterface[i] = v
	}

	resp, err = _createResponseFromArray(responseCode, responseInterface)
	if err != nil {
		log.Panic(err)
	}

	return resp
}

// +
func (u *User) listFollowers() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.follower WHERE followee = ?"

	resp := u.listBasic(query)
	return resp
}

// +
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

// +
func (u *User) unfollow() string {
	var resp string
	query := "DELETE FROM follow WHERE follower = ? AND followee = ?"

	args := Args{}

	resp, err := validateJson(u.inputRequest, "follower", "followee")
	if err != nil {
		return createInvalidJsonResponse(&u.inputRequest.json)
	}

	args.generateFromJson(&u.inputRequest.json, "follower", "followee")

	_, err = execQuery(query, &args.data, u.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	clearQuery(&u.inputRequest.query)
	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	return u.getDetails()
}

// +
func (u *User) updateProfile() string {
	var resp string
	query := "UPDATE user SET about = ?, name = ? WHERE email =  ?"

	args := Args{}

	resp, err := validateJson(u.inputRequest, "about", "name", "user")
	if err != nil {
		return createInvalidJsonResponse(&u.inputRequest.json)
	}

	args.generateFromJson(&u.inputRequest.json, "about", "name", "user")

	_, err = execQuery(query, &args.data, u.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	clearQuery(&u.inputRequest.query)
	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["user"].(string))
	return u.getDetails()
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
	args := Args{}

	query := "INSERT INTO forum (name, short_name, user) VALUES(?, ?, ?)"

	resp, err := validateJson(f.inputRequest, "name", "short_name", "user")
	if err != nil {
		return createInvalidJsonResponse(&f.inputRequest.json)
	}

	args.generateFromJson(&f.inputRequest.json, "name", "short_name", "user")

	dbResp, err := execQuery(query, &args.data, f.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	query = "SELECT * FROM forum WHERE id = ?"
	args.clear()
	args.append(dbResp.lastId)
	newForum, err := selectQuery(query, &args.data, f.db)
	if err != nil {
		log.Panic(err)
	}

	responseCode := 0
	responseMsg := &rs.ForumCreate{
		Name:       newForum.values[0]["name"],
		Short_Name: newForum.values[0]["short_name"],
		Id:         dbResp.lastId,
		User:       newForum.values[0]["user"],
	}

	resp, err = _createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Forum '%s' created", responseMsg.Short_Name)

	return resp
}

func (f *Forum) getForumDetails() (int, map[string]interface{}) {
	query := "SELECT * FROM forum WHERE short_name = ?"
	var args []interface{}
	args = append(args, f.inputRequest.query["forum"][0])

	getForum, err := selectQuery(query, &args, f.db)
	if err != nil {
		log.Panic(err)
	}

	if getForum.rows == 0 {
		responseCode := 1
		errorMessage := map[string]interface{}{"msg": "Not found"}

		return responseCode, errorMessage
	}

	respId, _ := strconv.ParseInt(getForum.values[0]["id"], 10, 64)

	responseCode := 0
	responseMsg := map[string]interface{}{
		"id":         respId,
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

func (f *Forum) listThreads() string {
	relatedUser := false
	relatedForum := false

	t := Thread{inputRequest: f.inputRequest, db: f.db}
	var resp, query string
	var args []interface{}

	// Validate query values
	if len(t.inputRequest.query["forum"]) == 1 {
		query = "SELECT t.*, (SELECT COUNT(*) FROM post p WHERE p.thread = t.id AND p.isDeleted = false) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.forum = ?"
		args = append(args, t.inputRequest.query["forum"][0])
	} else {
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
	} else {
		query += " ORDER BY t.date desc"
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

	// related params
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("user", f.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("forum", f.inputRequest.query["related"]) {
		relatedForum = true
	}
	// Response here
	responseCode, responseMsg := t.getArrayThreadsDetails(query, args)

	if responseCode != 0 {
		// E6ANYI KOSTYL`
		test := make(map[string]interface{})
		resp, _ = createResponse(0, test)

		return resp
		resp, _ = createResponse(responseCode, responseMsg[0])
		return resp
	}
	for _, value := range responseMsg {
		if relatedUser {
			u := User{inputRequest: f.inputRequest, db: f.db}
			u.inputRequest.query["user"] = u.inputRequest.query["user"][0:0]
			u.inputRequest.query["user"] = append(u.inputRequest.query["user"], value["user"].(string))

			_, responseUser := u.getUserDetails()
			value["user"] = responseUser
		}

		if relatedForum {

			f := Forum{inputRequest: f.inputRequest, db: f.db}
			f.inputRequest.query["forum"] = append(f.inputRequest.query["forum"], value["forum"].(string))

			_, responseForum := f.getForumDetails()
			value["forum"] = responseForum
		}
	}
	resp, err := createResponseFromArray(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}
	return resp
}

func (f *Forum) listPosts() string {
	var query, order, resp string
	relatedUser := false
	relatedThread := false
	relatedForum := false
	var args []interface{}

	p := Post{inputRequest: f.inputRequest, db: f.db}

	// Validate query values
	if len(p.inputRequest.query["forum"]) == 1 {
		query = "SELECT * FROM post p WHERE p.forum = ?"
		args = append(args, p.inputRequest.query["forum"][0])

	} else {
		resp = createInvalidResponse()
		return resp
	}

	// related params
	// var join string
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("user", f.inputRequest.query["related"]) {
		relatedUser = true
		// query += " "
		// join += " LEFT JOIN user u ON p.user=u.email"
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("thread", f.inputRequest.query["related"]) {
		relatedThread = true
		// query += " "
		// join += " LEFT JOIN thread t ON p.thread=t.id"
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("forum", f.inputRequest.query["related"]) {
		relatedForum = true
		// query += " "
		// join += " LEFT JOIN forum f ON p.forum=f.short_name"
	}

	// order by here
	if len(p.inputRequest.query["order"]) >= 1 {
		orderType := p.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			return createInvalidResponse()
		}

		order = fmt.Sprintf(" ORDER BY date %s", orderType)
	} else {
		order = " ORDER BY date DESC"
	}

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
			return createInvalidResponse()
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	getPost, err := selectQuery(query, &args, p.db)
	if err != nil {
		log.Panic(err)
	}

	if getPost.rows == 0 {
		// E6ANYI KOSTYL`
		test := make(map[string]interface{})
		resp, _ = createResponse(0, test)

		return resp
		return createNotFoundForArray()
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

		if relatedUser {
			u := User{inputRequest: f.inputRequest, db: f.db}
			u.inputRequest.query["user"] = u.inputRequest.query["user"][0:0]
			u.inputRequest.query["user"] = append(u.inputRequest.query["user"], value["user"])

			_, responseUser := u.getUserDetails()
			tempMsg["user"] = responseUser
		}

		if relatedThread {
			t := Thread{inputRequest: f.inputRequest, db: f.db}
			t.inputRequest.query["thread"] = t.inputRequest.query["thread"][0:0]
			t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], value["thread"])

			_, responseThread := t.getThreadDetails()
			tempMsg["thread"] = responseThread
		}

		if relatedForum {

			f := Forum{inputRequest: f.inputRequest, db: f.db}
			f.inputRequest.query["forum"] = f.inputRequest.query["forum"][0:0]
			f.inputRequest.query["forum"] = append(f.inputRequest.query["forum"], value["forum"])

			_, responseForum := f.getForumDetails()
			tempMsg["forum"] = responseForum
		}

		responseMsg = append(responseMsg, tempMsg)
	}

	resp, _ = createResponseFromArray(responseCode, responseMsg)
	return resp
}

func (f *Forum) listUsers() string {
	var resp, query string
	var args []interface{}

	// Validate query values
	if len(f.inputRequest.query["forum"]) == 1 {
		query = "SELECT DISTINCT p.user FROM post p JOIN user u ON p.user=u.email WHERE p.forum = ?"
		args = append(args, f.inputRequest.query["forum"][0])
	} else {
		resp = createInvalidResponse()
		return resp
	}

	// Check and validate optional params
	if len(f.inputRequest.query["since_id"]) >= 1 {
		query += " AND u.id >= ?"
		args = append(args, f.inputRequest.query["since_id"][0])
	}

	if len(f.inputRequest.query["order"]) >= 1 {
		orderType := f.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			resp = createInvalidResponse()
			return resp
		}

		// query += " ORDER BY u.name desc"
		query += fmt.Sprintf(" ORDER BY u.name %s", orderType)
	} else {
		query += " ORDER BY u.name desc"
	}
	if len(f.inputRequest.query["limit"]) >= 1 {
		limitValue := f.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			resp = createInvalidResponse()
			return resp
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	fmt.Println(query)

	getUser, err := selectQuery(query, &args, f.db)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(getUser)

	if getUser.rows == 0 {
		// E6ANYI KOSTYL`
		test := make(map[string]interface{})
		resp, _ = createResponse(0, test)

		return resp
		return createNotFoundForArray()
	}

	responseCode := 0
	var responseMsg []map[string]interface{}

	for _, value := range getUser.values {
		fmt.Println(value)
		u := User{inputRequest: f.inputRequest, db: f.db}
		u.inputRequest.query["user"] = u.inputRequest.query["user"][0:0]
		u.inputRequest.query["user"] = append(u.inputRequest.query["user"], value["user"])

		_, responseUser := u.getUserDetails()

		fmt.Println(responseUser)

		responseMsg = append(responseMsg, responseUser)
	}

	resp, err = createResponseFromArray(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func forumHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	forum := Forum{inputRequest: inputRequest, db: db}

	var result string

	if inputRequest.method == "GET" {

		if inputRequest.path == "/db/api/forum/details/" {
			result = forum.details()
		} else if inputRequest.path == "/db/api/forum/listPosts/" {
			result = forum.listPosts()
		} else if inputRequest.path == "/db/api/forum/listThreads/" {
			result = forum.listThreads()
		} else if inputRequest.path == "/db/api/forum/listUsers/" {
			result = forum.listUsers()
		}
	} else if inputRequest.method == "POST" {

		// Like Router
		if inputRequest.path == "/db/api/forum/create/" {
			result = forum.create()
		}
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
	args := Args{}

	resp, err := validateJson(t.inputRequest, "thread")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		return createInvalidResponse()
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args.append(value, threadId)

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
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
	responseMsg := &rs.ThreadBoolBasic{
		Thread: threadId,
	}

	resp, err = _createResponse(responseCode, responseMsg)
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
	args := Args{}

	query := "INSERT INTO thread (forum, title, isClosed, user, date, message, slug, isDeleted) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	resp, err := validateJson(t.inputRequest, "forum", "title", "isClosed", "user", "date", "message", "slug")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	args.generateFromJson(&t.inputRequest.json, "forum", "title", "isClosed", "user", "date", "message", "slug")

	// Validate isDeleted param
	isDeleted := t.inputRequest.json["isDeleted"]
	if isDeleted == nil {
		args.append(false)
	} else {
		if isDeleted != false && isDeleted != true {
			return createInvalidResponse()
		}
		args.append(isDeleted)
	}

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	query = "SELECT * FROM thread WHERE id = ?"
	args.clear()
	args.append(dbResp.lastId)
	newThread, err := selectQuery(query, &args.data, t.db)
	if err != nil {
		log.Panic(err)
	}

	responseCode := 0
	responseMsg := &rs.ThreadCreate{
		Forum:     newThread.values[0]["forum"],
		Title:     newThread.values[0]["title"],
		Id:        dbResp.lastId,
		User:      newThread.values[0]["user"],
		Date:      newThread.values[0]["date"],
		Message:   newThread.values[0]["message"],
		Slug:      newThread.values[0]["slug"],
		IsClosed:  stringToBool(newThread.values[0]["isClosed"]),
		IsDeleted: stringToBool(newThread.values[0]["isDeleted"]),
	}

	resp, err = _createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Thread '#%d' created", responseMsg.Id)

	return resp
}

func (t *Thread) getThreadDetails() (int, map[string]interface{}) {
	// query := "SELECT t.*, COUNT(*) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.id = ?"		// FIX
	query := "SELECT t.*, (SELECT COUNT(*) FROM post p WHERE p.thread = ? AND p.isDeleted = false) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.id = ?"
	var args []interface{}
	args = append(args, t.inputRequest.query["thread"][0])
	args = append(args, t.inputRequest.query["thread"][0])

	getThread, err := selectQuery(query, &args, t.db)
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
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

	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("thread", t.inputRequest.query["related"]) {
		return createInvalidQuery()
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
		query = "SELECT t.*, (SELECT COUNT(*) FROM post p WHERE p.thread = t.id AND p.isDeleted = false) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.user = ?"
		args = append(args, t.inputRequest.query["user"][0])

		f = true
	}
	if len(t.inputRequest.query["forum"]) == 1 && f == false {
		query = "SELECT t.*, (SELECT COUNT(*) FROM post p WHERE p.thread = t.id AND p.isDeleted = false) posts FROM thread t LEFT JOIN post p ON t.id=p.thread WHERE t.forum = ?"
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
		if responseCode == 1 {
			// E6ANYI KOSTYL`
			test := make(map[string]interface{})
			resp, _ = createResponse(0, test)
			return resp
		}
		resp, _ = createResponse(responseCode, responseMsg[0])
		return resp
	}
	resp, err := createResponseFromArray(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	return resp
}

func (t *Thread) parentTree(order string) (int, []map[string]interface{}) {
	var args []interface{}
	query := "SELECT parent FROM post WHERE thread = ?"
	order = " ORDER BY parent " + order

	args = append(args, t.inputRequest.query["thread"][0])

	// Check and validate optional params
	if len(t.inputRequest.query["since"]) >= 1 {
		query += " AND date > ?"
		args = append(args, t.inputRequest.query["since"][0])
	}

	query += " GROUP BY id HAVING LENGTH(parent) = 5"
	query += order

	if len(t.inputRequest.query["limit"]) >= 1 {
		limitValue := t.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			return 100500, nil
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	//
	// query here
	//

	getPost, err := selectQuery(query, &args, t.db)
	if err != nil {
		log.Panic(err)
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
		subQuery := "SELECT * FROM post WHERE thread = ? AND parent LIKE ? ORDER BY parent"
		var subArgs []interface{}
		subArgs = append(subArgs, t.inputRequest.query["thread"][0])
		subArgs = append(subArgs, value["parent"]+"%")

		getSubPost, err := selectQuery(subQuery, &subArgs, t.db)
		if err != nil {
			log.Panic(err)
		}

		for _, subValue := range getSubPost.values {
			respId, _ := strconv.ParseInt(subValue["id"], 10, 64)
			respLikes, _ := strconv.ParseInt(subValue["likes"], 10, 64)
			respDislikes, _ := strconv.ParseInt(subValue["dislikes"], 10, 64)
			respPoints, _ := strconv.ParseInt(subValue["points"], 10, 64)
			respThread, _ := strconv.ParseInt(subValue["thread"], 10, 64)
			respIsApproved, _ := strconv.ParseBool(subValue["isApproved"])
			respIsDeleted, _ := strconv.ParseBool(subValue["isDeleted"])
			respIsEdited, _ := strconv.ParseBool(subValue["isEdited"])
			respIsHighlighted, _ := strconv.ParseBool(subValue["isHighlighted"])
			respIsSpam, _ := strconv.ParseBool(subValue["isSpam"])

			tempMsg := map[string]interface{}{
				"date":          subValue["date"],
				"dislikes":      respDislikes,
				"forum":         subValue["forum"],
				"id":            respId,
				"isApproved":    respIsApproved,
				"isDeleted":     respIsDeleted,
				"isEdited":      respIsEdited,
				"isHighlighted": respIsHighlighted,
				"isSpam":        respIsSpam,
				"likes":         respLikes,
				"message":       subValue["message"],
				"parent":        nil,
				"points":        respPoints,
				"thread":        respThread,
				"user":          subValue["user"],
			}

			p := Post{inputRequest: t.inputRequest, db: t.db}
			parent := p.getParentId(respId, subValue["parent"])
			if parent == int(respId) {
				tempMsg["parent"] = nil
			} else {
				tempMsg["parent"] = parent
			}
			responseMsg = append(responseMsg, tempMsg)
		}
	}

	return responseCode, responseMsg
}

func (t *Thread) listPosts() string {
	var query, order, sort, resp string

	parentTree := false
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
		order = orderType
	} else {
		order = "DESC"
	}

	// sort here
	if len(t.inputRequest.query["sort"]) >= 1 {
		sortType := t.inputRequest.query["sort"][0]

		switch sortType {
		case "flat":
			sort = " ORDER BY date " + order
		case "tree":
			sort = " ORDER BY SUBSTRING(parent, 1, 5) " + order + ", parent asc "
		case "parent_tree":
			parentTree = true
		default:
			resp = createInvalidResponse()
			return resp
		}
	} else {
		sort = " ORDER BY date " + order
	}

	var responseCode int
	var responseMsg []map[string]interface{}
	// simple sort
	if parentTree == false {
		p := Post{inputRequest: t.inputRequest, db: t.db}
		responseCode, responseMsg = p.getList(query, sort, args)
	} else {
		// parent_tree sort
		responseCode, responseMsg = t.parentTree(order)
	}

	// check responseCode
	if responseCode == 100500 {
		resp = createInvalidResponse()
		return resp
	} else if responseCode == 1 {
		// E6ANYI KOSTYL`
		test := make(map[string]interface{})
		resp, _ = createResponse(0, test)
		return resp

		// resp, _ = createResponse(responseCode, responseMsg[0])
		// return resp
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

	return resp
}

func (t *Thread) open() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, false)

	return resp
}

func (t *Thread) remove() string {
	query := "UPDATE thread SET isDeleted = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, true)

	query = "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, true)

	return resp
}

func (t *Thread) restore() string {
	query := "UPDATE thread SET isDeleted = ? WHERE id = ?"

	resp := t.updateBoolBasic(query, false)

	query = "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, false)

	return resp
}

func (t *Thread) subscribe() string {
	var resp string
	args := Args{}

	query := "INSERT INTO subscribe (thread, user) VALUES(?, ?)"

	resp, err := validateJson(t.inputRequest, "thread", "user")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	args.generateFromJson(&t.inputRequest.json, "thread", "user")

	_, err = execQuery(query, &args.data, t.db)
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
			log.Panic(err)
		}

		return resp
	}

	// else return info

	responseCode := 0
	responseMsg := &rs.ThreadSubscribe{
		Thread: int64(t.inputRequest.json["thread"].(float64)),
		User:   t.inputRequest.json["user"].(string),
	}

	resp, err = _createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("User '%s' subscribe to thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

func (t *Thread) unsubscribe() string {
	var resp string
	args := Args{}

	query := "DELETE FROM subscribe WHERE thread = ? AND user = ?"

	resp, err := validateJson(t.inputRequest, "thread", "user")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	args.generateFromJson(&t.inputRequest.json, "thread", "user")

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
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
	responseMsg := &rs.ThreadSubscribe{
		Thread: int64(t.inputRequest.json["thread"].(float64)),
		User:   t.inputRequest.json["user"].(string),
	}

	resp, err = _createResponse(responseCode, responseMsg)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("User '%s' unsubscribe from thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

func (t *Thread) update() string {
	var resp string
	args := Args{}

	query := "UPDATE thread SET message = ?, slug = ? WHERE id = ?"

	resp, err := validateJson(t.inputRequest, "thread", "message", "slug")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args.generateFromJson(&t.inputRequest.json, "message", "slug", "thread")

	_, err = execQuery(query, &args.data, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
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
	args := Args{}

	resp, err := validateJson(t.inputRequest, "thread", "vote")
	if err != nil {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	if checkFloat64Type(t.inputRequest.json["thread"]) == false || checkFloat64Type(t.inputRequest.json["vote"]) == false {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	threadId := t.inputRequest.json["thread"].(float64)
	vote := t.inputRequest.json["vote"].(float64)

	if vote == 1 {
		query = "UPDATE thread SET likes = likes + 1, points = points + 1 WHERE id = ?"
	} else if vote == -1 {
		query = "UPDATE thread SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?"
	} else {
		return createInvalidJsonResponse(&t.inputRequest.json)
	}

	args.append(threadId)

	_, err = execQuery(query, &args.data, t.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
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
		} else if inputRequest.path == "/db/api/thread/listPosts/" {
			result = thread.listPosts()
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
	fmt.Println("post.create()\t", p.inputRequest)
	var resp string
	query := "INSERT INTO post (thread, message, user, forum, date, isApproved, isHighlighted, isEdited, isSpam, isDeleted, parent) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	// 	ir.url = fmt.Sprintf("%v", r.URL)

	var args []interface{}

	resp, err := validateJson(p.inputRequest, "thread", "message", "user", "forum", "date")
	if err != nil {
		return resp
	}

	args = append(args, p.inputRequest.json["thread"])
	args = append(args, p.inputRequest.json["message"])
	args = append(args, p.inputRequest.json["user"])
	args = append(args, p.inputRequest.json["forum"])
	args = append(args, p.inputRequest.json["date"])

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
			log.Panic(err)
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
			log.Panic(err)
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
			log.Panic(err)
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
		args := Args{}
		args.append(parentId)

		getParent, _ := selectQuery(query, &args.data, p.db)

		respId, _ := strconv.ParseInt(getParent.values[0]["id"], 10, 64)

		return int(respId)
	}
}

func (p *Post) updateBoolBasic(query string, value bool) string {
	var resp string

	args := Args{}

	resp, err := validateJson(p.inputRequest, "post")
	if err != nil {
		return createInvalidJsonResponse(&p.inputRequest.json)
	}

	if checkFloat64Type(p.inputRequest.json["post"]) == false {
		return createInvalidJsonResponse(&p.inputRequest.json)
	}

	postId := p.inputRequest.json["post"].(float64)

	args.append(value, postId)

	dbResp, err := execQuery(query, &args.data, p.db)
	if err != nil {
		log.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	if dbResp.rowCount == 0 {
		p.inputRequest.query["post"] = append(p.inputRequest.query["post"], intToString(int(postId)))
		responseCode, responseMsg := p.getPostDetails()

		if responseCode != 0 {
			return createNotExistResponse()
		}

		resp, err = createResponse(responseCode, responseMsg)
		if err != nil {
			panic(err.Error())
		}
		return resp
	}

	responseCode := 0
	responseMsg := &rs.PostBoolBasic{
		Post: postId,
	}

	resp, err = _createResponse(responseCode, responseMsg)
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
		log.Panic(err)
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
	fmt.Println("post.details()\t", p.inputRequest)
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
		t.inputRequest.query["thread"] = t.inputRequest.query["thread"][0:0]
		t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], int64ToString(responseMsg["thread"].(int64)))
		_, threadDetails := t.getThreadDetails()

		// threadDetails["posts"] = 1

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
	fmt.Println(resp)

	return resp
}

func (p *Post) getArrayPostDetails(query string, args []interface{}) (int, []map[string]interface{}) {
	getPost, err := selectQuery(query, &args, p.db)
	if err != nil {
		log.Panic(err)
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
	fmt.Println("post.list()\t", p.inputRequest)
	var query, order, resp string
	f := false
	var args []interface{}

	// Validate query values
	if len(p.inputRequest.query["thread"]) == 1 {
		query = "SELECT * FROM post WHERE thread = ?"
		args = append(args, p.inputRequest.query["thread"][0])

		f = true
	}
	if len(p.inputRequest.query["user"]) == 1 && f == false {
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
	} else if responseCode == 1 {
		// E6ANYI KOSTYL`
		test := make(map[string]interface{})
		resp, _ = createResponse(0, test)
		// resp, _ = createResponse(responseCode, responseMsg[0])

	} else if responseCode == 0 {
		resp, _ = createResponseFromArray(responseCode, responseMsg)
	} else {
		resp = createInvalidResponse()
	}

	fmt.Println(resp)
	return resp
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
	query := "UPDATE post SET message = ? WHERE id = ?"

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

	dbResp, err := execQuery(query, &args, p.db)
	if err != nil {
		fmt.Println(err)
		responseCode, errorMessage := errorExecParse(err)

		resp, err = createResponse(responseCode, errorMessage)
		if err != nil {
			log.Panic(err)
		}

		return resp
	}

	if dbResp.rowCount == 0 {
		query := "UPDATE post SET isEdited = false WHERE id = ?"
		_, _ = execQuery(query, &args, p.db)
	} else {
		query := "UPDATE post SET isEdited = true WHERE id = ?"
		_, _ = execQuery(query, &args, p.db)
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
			log.Panic(err)
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

// ==========================
// Information methods here
// ==========================
func statusHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	if inputRequest.method == "GET" {
		var args []interface{}

		query := "SELECT COUNT(*) count FROM user"
		dbResp, _ := selectQuery(query, &args, db)
		respUsers, _ := strconv.ParseInt(dbResp.values[0]["count"], 10, 64)

		query = "SELECT COUNT(*) count FROM thread"
		dbResp, _ = selectQuery(query, &args, db)
		respThreads, _ := strconv.ParseInt(dbResp.values[0]["count"], 10, 64)

		query = "SELECT COUNT(*) count FROM forum"
		dbResp, _ = selectQuery(query, &args, db)
		respForums, _ := strconv.ParseInt(dbResp.values[0]["count"], 10, 64)

		query = "SELECT COUNT(*) count FROM post"
		dbResp, _ = selectQuery(query, &args, db)
		respPosts, _ := strconv.ParseInt(dbResp.values[0]["count"], 10, 64)

		responseCode := 0
		responseMsg := map[string]interface{}{
			"user":   respUsers,
			"thread": respThreads,
			"forum":  respForums,
			"post":   respPosts,
		}
		resp, _ := createResponse(responseCode, responseMsg)

		io.WriteString(w, resp)
	}
}

func clearHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	if inputRequest.method == "POST" {
		var args []interface{}

		query := "DELETE FROM follow"
		_, _ = execQuery(query, &args, db)
		query = "DELETE FROM subscribe"
		_, _ = execQuery(query, &args, db)

		query = "DELETE FROM post WHERE id > 0"
		_, _ = execQuery(query, &args, db)
		query = "DELETE FROM thread WHERE id > 0"
		_, _ = execQuery(query, &args, db)
		query = "DELETE FROM forum WHERE id > 0"
		_, _ = execQuery(query, &args, db)
		query = "DELETE FROM user WHERE id > 0"
		_, _ = execQuery(query, &args, db)

		responseCode := 0

		cacheContent := map[string]interface{}{
			"code":     responseCode,
			"response": "OK",
		}

		str, _ := json.Marshal(cacheContent)

		io.WriteString(w, string(str))
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
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	PORT := ":8000"

	fmt.Printf("The server is running on http://localhost%s\n", PORT)

	// go initLog()

	http.HandleFunc("/db/api/user/", makeHandler(db, userHandler))
	http.HandleFunc("/db/api/forum/", makeHandler(db, forumHandler))
	http.HandleFunc("/db/api/thread/", makeHandler(db, threadHandler))
	http.HandleFunc("/db/api/post/", makeHandler(db, postHandler))
	http.HandleFunc("/db/api/status/", makeHandler(db, statusHandler))
	http.HandleFunc("/db/api/clear/", makeHandler(db, clearHandler))

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
		log.Panic(err)
	}
	return result
}

func stringToInt64(inputStr string) int64 {
	return int64(stringToInt(inputStr))
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

func stringToBool(inputString string) (result bool) {
	result, err := strconv.ParseBool(inputString)

	if err != nil {
		log.Panic(err)
	}

	return
}

// =================
// Future here
// =================

/*
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
*/
