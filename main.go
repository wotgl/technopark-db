package main

import (
	"database/sql"
	"encoding/json"
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

func _createResponse(code int, response rs.RespStruct) string {
	content := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(content)
	if err != nil {
		log.Println("Error encoding JSON")
		log.Panic(err)
	}

	return string(str)
}

func _createResponseFromArray(code int, response []interface{}) string {
	cacheContent := map[string]interface{}{
		"code":     code,
		"response": response,
	}

	str, err := json.Marshal(cacheContent)
	if err != nil {
		log.Println("Error encoding JSON")
		log.Panic(err)
	}

	return string(str)
}

func errorExecParse(err error) (int, string) {
	log.Println("Error:\t", err)
	if driverErr, ok := err.(*mysql.MySQLError); ok { // Now the error number is accessible directly
		var responseCode int
		var errorMessage string

		switch driverErr.Number {
		case 1062:
			responseCode = 5
			errorMessage = "Exist"

		// Error 1452: Cannot add or update a child row: a foreign key constraint fails
		case 1452:
			responseCode = 5
			errorMessage = "Exist [Error 1452]"

		default:
			// fmt.Println("errorExecParse() default")
			panic(err.Error())
			responseCode = 4
			errorMessage = "Unknown Error"
		}

		return responseCode, errorMessage
	}
	panic(err.Error()) // proper error handling instead of panic in your app
}

func createErrorResponse(err error) string {
	responseCode, msg := errorExecParse(err)
	errorMessage := &rs.ErrorMsg{
		Msg: msg,
	}

	return _createResponse(responseCode, errorMessage)
}

func createInvalidQuery() string {
	responseCode := 3
	errorMessage := &rs.ErrorMsg{
		Msg: "Invalid query",
	}

	return _createResponse(responseCode, errorMessage)
}

func createInvalidJsonResponse(inputRequest *InputRequest) string {
	responseCode := 3
	errorMessage := &rs.ErrorMsg{
		Msg: "Invalid json",
	}

	log.Println("Invalid JSON:\turl=\tjson=", inputRequest.url, inputRequest.json)

	return _createResponse(responseCode, errorMessage)
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
	errorMessage := &rs.ErrorMsg{
		Msg: "Invalid",
	}

	return _createResponse(responseCode, errorMessage)
}

func createNotExistResponse() string {
	responseCode := 1
	errorMessage := &rs.ErrorMsg{
		Msg: "Not exist",
	}

	return _createResponse(responseCode, errorMessage)
}

// KOSTYL` API
func becauseAPI() string {
	kostyl := make(map[string]interface{})

	content := map[string]interface{}{
		"code":     0,
		"response": kostyl,
	}

	str, err := json.Marshal(content)
	if err != nil {
		log.Println("Error encoding JSON")
		log.Panic(err)
	}

	return string(str)
}

func validateJson(ir *InputRequest, args ...string) bool {
	for _, value := range args {
		if reflect.TypeOf(ir.json[value]) == nil {
			return false
		}
	}

	return true
}

func validateBoolParams(json map[string]interface{}, args *Args, params ...string) bool {
	for _, value := range params {

		if json[value] == nil {
			args.append(false)
		} else {
			if reflect.TypeOf(json[value]).Kind() != reflect.Bool {
				return false
			}
			args.append(json[value])
		}
	}

	return true
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

func (args *Args) clear() { args.data = args.data[0:0] }

// ======================
// Database queries here
// ======================

type ExecResponse struct {
	lastId   int64
	rowCount int64
}

func execQuery(query string, args *[]interface{}, db *sql.DB) (*ExecResponse, error) {
	resp := new(ExecResponse)

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(*args...)
	if err != nil {
		return nil, err
	}

	// err = tx.Commit()

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	// log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	resp.lastId = lastId
	resp.rowCount = rowCnt

	return resp, nil
}

type SelectResponse struct {
	rows    int
	columns []string
	values  []map[string]string
}

func selectQuery(query string, args *[]interface{}, db *sql.DB) *SelectResponse {
	resp := new(SelectResponse)

	/*
		f := false
		_query := query
		for _, value := range *args {
			var tmp string
			switch reflect.TypeOf(value).Kind() {
			case reflect.Float64:
				tmp = floatToString(value.(float64))
			case reflect.String:
				tmp = value.(string)
			case reflect.Bool:
				tmp = strconv.FormatBool(value.(bool))
			case reflect.Int64:
				tmp = int64ToString(value.(int64))
			}

			if strings.ContainsAny(tmp, ` !"#$&'()*+,-.\/:;<=>?@[\]^`+"`"+`{|}~`) {
				f = true
				break
			} else {
				if reflect.TypeOf(value).Kind() == reflect.String {
					tmp = "\"" + value.(string) + "\""
				}

				fmt.Println(_query)
				_query = strings.Replace(_query, "?", tmp, 1)
			}
		}

		var rows *sql.Rows
		var err error

		if f {
			rows, err = db.Query(query, *args...)
		} else {
			rows, err = db.Query(_query)
		}
	*/
	rows, err := db.Query(query, *args...)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Panic(err)
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
			log.Panic(err)
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
		log.Panic(err)
	}

	return resp
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

	if !validateJson(u.inputRequest, "email") {
		return createInvalidJsonResponse(u.inputRequest)
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
		return createErrorResponse(err)
	}

	query = "SELECT * FROM user WHERE id = ?"
	args.clear()
	args.append(dbResp.lastId)
	newUser := selectQuery(query, &args.data, u.db)

	responseCode := 0
	responseMsg := &rs.UserCreate{
		About:       newUser.values[0]["about"],
		Email:       newUser.values[0]["email"],
		Id:          dbResp.lastId,
		IsAnonymous: stringToBool(newUser.values[0]["isAnonymous"]),
		Name:        newUser.values[0]["name"],
		Username:    newUser.values[0]["username"],
	}

	resp = _createResponse(responseCode, responseMsg)

	log.Printf("User '%s' created", responseMsg.Email)

	return resp
}

// +
func (u *User) _getUserDetails(args Args) (int, *rs.UserDetails) {
	query := "SELECT * FROM user WHERE email = ?"

	getUser := selectQuery(query, &args.data, u.db)

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

// +
func (u *User) getUserFollowers(followee string) []string {
	args := Args{}
	args.append(followee)

	// followers here
	query := "SELECT follower FROM follow WHERE followee = ?"
	getUserFollowers := selectQuery(query, &args.data, u.db)

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
	getUserFollowing := selectQuery(query, &args.data, u.db)

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
	getUserSubscriptions := selectQuery(query, &args.data, u.db)

	listSubscriptions := make([]int, 0)
	for _, value := range getUserSubscriptions.values {
		listSubscriptions = append(listSubscriptions, stringToInt(value["thread"]))
	}

	return listSubscriptions
}

// +
func (u *User) getDetails() string {
	if len(u.inputRequest.query["user"]) != 1 {
		return createInvalidResponse()
	}

	args := Args{}
	args.append(u.inputRequest.query["user"][0])

	responseCode, responseMsg := u._getUserDetails(args)
	if responseCode == 1 {
		return createNotExistResponse()
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (u *User) follow() string {
	query := "INSERT INTO follow (follower, followee) VALUES(?, ?)"
	args := Args{}

	if !validateJson(u.inputRequest, "follower", "followee") {
		return createInvalidJsonResponse(u.inputRequest)
	}

	args.generateFromJson(&u.inputRequest.json, "follower", "followee")

	_, err := execQuery(query, &args.data, u.db)
	if err != nil {
		// return exist
		if checkError1062(err) == true {
			clearQuery(&u.inputRequest.query)
			u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))

			return u.getDetails()
		}

		responseCode, msg := errorExecParse(err)
		errorMessage := &rs.ErrorMsg{
			Msg: msg,
		}

		return _createResponse(responseCode, errorMessage)
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	return u.getDetails()
}

// +
func (u *User) listBasic(query string) string {
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
	getUserFollowers := selectQuery(query, &args.data, u.db)

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

	return _createResponseFromArray(responseCode, responseInterface)
}

// +
func (u *User) listFollowers() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.follower WHERE followee = ?"

	return u.listBasic(query)
}

// +
func (u *User) listFollowing() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.followee WHERE follower = ?"

	return u.listBasic(query)
}

// +
func (u *User) listPosts() string {
	delete(u.inputRequest.query, "forum")

	p := Post{inputRequest: u.inputRequest, db: u.db}
	return p.list()
}

// +
func (u *User) unfollow() string {
	query := "DELETE FROM follow WHERE follower = ? AND followee = ?"
	args := Args{}

	if !validateJson(u.inputRequest, "follower", "followee") {
		return createInvalidJsonResponse(u.inputRequest)
	}

	args.generateFromJson(&u.inputRequest.json, "follower", "followee")

	_, err := execQuery(query, &args.data, u.db)
	if err != nil {
		return createErrorResponse(err)
	}

	clearQuery(&u.inputRequest.query)
	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	return u.getDetails()
}

// +
func (u *User) updateProfile() string {
	query := "UPDATE user SET about = ?, name = ? WHERE email =  ?"
	args := Args{}

	if !validateJson(u.inputRequest, "about", "name", "user") {
		return createInvalidJsonResponse(u.inputRequest)
	}

	args.generateFromJson(&u.inputRequest.json, "about", "name", "user")

	_, err := execQuery(query, &args.data, u.db)
	if err != nil {
		return createErrorResponse(err)
	}

	clearQuery(&u.inputRequest.query)
	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["user"].(string))
	return u.getDetails()
}

func userHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	//t0 := time.Now()
	user := User{inputRequest: inputRequest, db: db}
	var result string

	if inputRequest.method == "GET" {
		switch inputRequest.path {
		case "/db/api/user/details/":
			result = user.getDetails()
		case "/db/api/user/listFollowers/":
			result = user.listFollowers()
		case "/db/api/user/listFollowing/":
			result = user.listFollowing()
		case "/db/api/user/listPosts/":
			result = user.listPosts()
		}
	} else if inputRequest.method == "POST" {
		switch inputRequest.path {
		case "/db/api/user/create/":
			result = user.create()
		case "/db/api/user/follow/":
			result = user.follow()
		case "/db/api/user/unfollow/":
			result = user.unfollow()
		case "/db/api/user/updateProfile/":
			result = user.updateProfile()
		}
	}

	// t1 := time.Now()
	// fmt.Printf("User handler: %v\n", t1.Sub(t0))
	io.WriteString(w, result)
}

// =================
// Forum handler here
// =================

type Forum struct {
	inputRequest *InputRequest
	db           *sql.DB
}

// +
func (f *Forum) create() string {
	var resp string
	args := Args{}

	query := "INSERT INTO forum (name, short_name, user) VALUES(?, ?, ?)"

	if !validateJson(f.inputRequest, "name", "short_name", "user") {
		return createInvalidJsonResponse(f.inputRequest)
	}

	args.generateFromJson(&f.inputRequest.json, "name", "short_name", "user")

	dbResp, err := execQuery(query, &args.data, f.db)
	if err != nil {
		return createErrorResponse(err)
	}

	responseCode := 0
	responseMsg := &rs.ForumCreate{
		Name:       f.inputRequest.json["name"].(string),
		Short_Name: f.inputRequest.json["short_name"].(string),
		Id:         dbResp.lastId,
		User:       f.inputRequest.json["user"].(string),
	}

	resp = _createResponse(responseCode, responseMsg)

	log.Printf("Forum '%s' created", responseMsg.Short_Name)

	return resp
}

// +
func (f *Forum) _getForumDetails(args Args) (int, *rs.ForumDetails) {
	query := "SELECT * FROM forum WHERE short_name = ?"

	getForum := selectQuery(query, &args.data, f.db)

	if getForum.rows == 0 {
		responseCode := 1
		errorMessage := &rs.ForumDetails{}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseMsg := &rs.ForumDetails{
		Id:         stringToInt64(getForum.values[0]["id"]),
		Short_Name: getForum.values[0]["short_name"],
		Name:       getForum.values[0]["name"],
		User:       getForum.values[0]["user"],
	}

	return responseCode, responseMsg
}

// +
func (f *Forum) details() string {
	var relatedUser bool
	args := Args{}

	if len(f.inputRequest.query["forum"]) != 1 {
		return createInvalidResponse()
	}
	if len(f.inputRequest.query["related"]) == 1 && f.inputRequest.query["related"][0] == "user" {
		relatedUser = true
	}

	args.append(f.inputRequest.query["forum"][0])

	responseCode, responseMsg := f._getForumDetails(args)

	if relatedUser {
		u := User{inputRequest: f.inputRequest, db: f.db}
		clearQuery(&u.inputRequest.query)
		args := Args{}
		args.append(responseMsg.User.(string))

		_, userDetails := u._getUserDetails(args)

		responseMsg.User = userDetails
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (f *Forum) listThreads() string {
	relatedUser := false
	relatedForum := false

	t := Thread{inputRequest: f.inputRequest, db: f.db}

	responseCode, responseMsg := t.listBasic()

	if responseCode == 1 {
		return becauseAPI()
	} else if responseCode != 0 {
		return createInvalidResponse()
	}

	// related params
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("user", f.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("forum", f.inputRequest.query["related"]) {
		relatedForum = true
	}

	// Response here
	for key, _ := range responseMsg.Threads {
		if relatedUser {
			u := User{inputRequest: f.inputRequest, db: f.db}
			userArgs := Args{}
			userArgs.append(responseMsg.Threads[key].User)

			_, responseUser := u._getUserDetails(userArgs)
			responseMsg.Threads[key].User = responseUser
		}

		if relatedForum {
			f := Forum{inputRequest: f.inputRequest, db: f.db}
			forumArgs := Args{}
			forumArgs.append(responseMsg.Threads[key].Forum)

			_, responseForum := f._getForumDetails(forumArgs)
			responseMsg.Threads[key].Forum = responseForum
		}
	}

	responseInterface := make([]interface{}, len(responseMsg.Threads))
	for i, v := range responseMsg.Threads {
		responseInterface[i] = v
	}

	return _createResponseFromArray(responseCode, responseInterface)
}

// +
func (f *Forum) listPosts() string {
	var query, order string
	relatedUser := false
	relatedThread := false
	relatedForum := false
	args := Args{}

	p := Post{inputRequest: f.inputRequest, db: f.db}

	// Validate query values
	if len(p.inputRequest.query["forum"]) == 1 {
		query = "SELECT * FROM post p WHERE p.forum = ?"
		args.append(p.inputRequest.query["forum"][0])
	} else {
		return createInvalidResponse()
	}

	// related params
	// var join string
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("user", f.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("thread", f.inputRequest.query["related"]) {
		relatedThread = true
	}
	if len(f.inputRequest.query["related"]) >= 1 && stringInSlice("forum", f.inputRequest.query["related"]) {
		relatedForum = true
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

	responseCode, responseMsg := p._getList(query, order, args)

	if responseCode == 1 {
		return becauseAPI()
	} else if responseCode != 0 {
		return createInvalidResponse()
	}

	for key, _ := range responseMsg.Posts {
		if relatedUser {
			u := User{inputRequest: f.inputRequest, db: f.db}
			userArgs := Args{}
			userArgs.append(responseMsg.Posts[key].User)

			_, responseUser := u._getUserDetails(userArgs)
			responseMsg.Posts[key].User = responseUser
		}

		if relatedThread {
			t := Thread{inputRequest: f.inputRequest, db: f.db}
			threadArgs := Args{}
			threadArgs.append(responseMsg.Posts[key].Thread)

			_, responseThread := t._getThreadDetails(threadArgs)
			responseMsg.Posts[key].Thread = responseThread
		}

		if relatedForum {
			f := Forum{inputRequest: f.inputRequest, db: f.db}
			forumArgs := Args{}
			forumArgs.append(responseMsg.Posts[key].Forum)

			_, responseForum := f._getForumDetails(forumArgs)
			responseMsg.Posts[key].Forum = responseForum
		}
	}

	responseInterface := make([]interface{}, len(responseMsg.Posts))
	for i, v := range responseMsg.Posts {
		responseInterface[i] = v
	}

	return _createResponseFromArray(responseCode, responseInterface)
}

// +
func (f *Forum) listUsers() string {
	var query string
	args := Args{}

	query = "SELECT u.email FROM user u WHERE email IN (SELECT DISTINCT p.user FROM post p WHERE p.forum = ?)"

	// Validate query values
	if len(f.inputRequest.query["forum"]) != 1 {
		return createInvalidResponse()
	}

	args.append(f.inputRequest.query["forum"][0])

	// Optional params
	if len(f.inputRequest.query["since_id"]) >= 1 {
		query += " AND u.id >= ?"
		args.append(f.inputRequest.query["since_id"][0])
	}

	if len(f.inputRequest.query["order"]) >= 1 {
		orderType := f.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			return createInvalidResponse()
		}

		query += fmt.Sprintf(" ORDER BY u.name %s", orderType)
	} else {
		query += " ORDER BY u.name desc"
	}

	if len(f.inputRequest.query["limit"]) >= 1 {
		i, err := strconv.Atoi(f.inputRequest.query["limit"][0])
		if err != nil || i < 0 {
			return createInvalidResponse()
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Query
	users := selectQuery(query, &args.data, f.db)

	if users.rows == 0 {
		return becauseAPI()
	}

	// Generate response
	responseCode := 0
	responseArray := make([]rs.UserDetails, 0)
	responseMsg := &rs.UserListBasic{Users: responseArray}

	for _, value := range users.values {
		u := User{inputRequest: f.inputRequest, db: f.db}
		userArgs := Args{}
		userArgs.append(value["email"])

		_, responseUser := u._getUserDetails(userArgs)

		responseMsg.Users = append(responseMsg.Users, *responseUser)
	}

	// Convert to interface
	responseInterface := make([]interface{}, len(responseMsg.Users))
	for i, v := range responseMsg.Users {
		responseInterface[i] = v
	}

	return _createResponseFromArray(responseCode, responseInterface)
}

func forumHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	//t0 := time.Now()
	forum := Forum{inputRequest: inputRequest, db: db}
	var result string

	if inputRequest.method == "GET" {
		switch inputRequest.path {
		case "/db/api/forum/details/":
			result = forum.details()
		case "/db/api/forum/listPosts/":
			result = forum.listPosts()
		case "/db/api/forum/listThreads/":
			result = forum.listThreads()
		case "/db/api/forum/listUsers/":
			result = forum.listUsers()
		}
	} else if inputRequest.method == "POST" {
		switch inputRequest.path {
		case "/db/api/forum/create/":
			result = forum.create()
		}
	}

	// t1 := time.Now()
	// fmt.Printf("Forum handler: %v\n", t1.Sub(t0))
	io.WriteString(w, result)
}

// =================
// Thread handler here
// =================

type Thread struct {
	inputRequest *InputRequest
	db           *sql.DB
}

// +
func (t *Thread) updateBoolBasic(query string, value bool) string {
	args := Args{}

	if !validateJson(t.inputRequest, "thread") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	if !checkFloat64Type(t.inputRequest.json["thread"]) {
		return createInvalidJsonResponse(t.inputRequest)
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args.append(value, threadId)

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		return createErrorResponse(err)
	}

	if dbResp.rowCount == 0 {
		threadArgs := Args{}
		threadArgs.append(threadId)
		responseCode, responseMsg := t._getThreadDetails(threadArgs)

		if responseCode != 0 {
			return createNotExistResponse()
		}

		return _createResponse(responseCode, responseMsg)
	}

	responseCode := 0
	responseMsg := &rs.ThreadBoolBasic{
		Thread: threadId,
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (t *Thread) close() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	return t.updateBoolBasic(query, true)
}

// +
func (t *Thread) create() string {
	var resp string
	args := Args{}

	query := "INSERT INTO thread (forum, title, isClosed, user, date, message, slug, isDeleted) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	if !validateJson(t.inputRequest, "forum", "title", "isClosed", "user", "date", "message", "slug") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	args.generateFromJson(&t.inputRequest.json, "forum", "title", "isClosed", "user", "date", "message", "slug")

	// Validate isDeleted param
	if !validateBoolParams(t.inputRequest.json, &args, "isDeleted") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		return createErrorResponse(err)
	}

	responseCode := 0
	responseMsg := &rs.ThreadCreate{
		Forum:     t.inputRequest.json["forum"].(string),
		Title:     t.inputRequest.json["title"].(string),
		Id:        dbResp.lastId,
		User:      t.inputRequest.json["user"].(string),
		Date:      t.inputRequest.json["date"].(string),
		Message:   t.inputRequest.json["message"].(string),
		Slug:      t.inputRequest.json["slug"].(string),
		IsClosed:  t.inputRequest.json["isClosed"].(bool),
		IsDeleted: t.inputRequest.json["isDeleted"].(bool),
	}

	resp = _createResponse(responseCode, responseMsg)

	log.Printf("Thread '#%d' created", responseMsg.Id)

	return resp
}

// Rewrite subquery
// +
func (t *Thread) _getThreadDetails(args Args) (int, *rs.ThreadDetails) {
	query := "SELECT t.* FROM thread t WHERE t.id = ?"

	getThread := selectQuery(query, &args.data, t.db)

	if getThread.rows == 0 {
		responseCode := 1
		errorMessage := &rs.ThreadDetails{}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseMsg := &rs.ThreadDetails{
		Date:      getThread.values[0]["date"],
		Dislikes:  stringToInt64(getThread.values[0]["dislikes"]),
		Forum:     getThread.values[0]["forum"],
		Id:        stringToInt64(getThread.values[0]["id"]),
		IsClosed:  stringToBool(getThread.values[0]["isClosed"]),
		IsDeleted: stringToBool(getThread.values[0]["isDeleted"]),
		Likes:     stringToInt64(getThread.values[0]["likes"]),
		Message:   getThread.values[0]["message"],
		Points:    stringToInt64(getThread.values[0]["points"]),
		Posts:     stringToInt64(getThread.values[0]["posts"]),
		Slug:      getThread.values[0]["slug"],
		Title:     getThread.values[0]["title"],
		User:      getThread.values[0]["user"],
	}

	return responseCode, responseMsg
}

// +
func (t *Thread) _getArrayThreadsDetails(query string, args Args) (int, *rs.ThreadList) {
	getThread := selectQuery(query, &args.data, t.db)

	if getThread.rows == 0 {
		responseCode := 1
		errorMessage := &rs.ThreadList{}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseArray := make([]rs.ThreadDetails, 0)
	responseMsg := &rs.ThreadList{Threads: responseArray}

	for _, value := range getThread.values {
		tempMsg := &rs.ThreadDetails{
			Date:      value["date"],
			Dislikes:  stringToInt64(value["dislikes"]),
			Forum:     value["forum"],
			Id:        stringToInt64(value["id"]),
			IsClosed:  stringToBool(value["isClosed"]),
			IsDeleted: stringToBool(value["isDeleted"]),
			Likes:     stringToInt64(value["likes"]),
			Message:   value["message"],
			Points:    stringToInt64(value["points"]),
			Posts:     stringToInt64(value["posts"]),
			Slug:      value["slug"],
			Title:     value["title"],
			User:      value["user"],
		}

		responseMsg.Threads = append(responseMsg.Threads, *tempMsg)
	}

	return responseCode, responseMsg
}

// +
// Проверка на пустой ответ с кодом 1 после _getThreadDetails()
func (t *Thread) details() string {
	var relatedUser, relatedForum bool

	if len(t.inputRequest.query["thread"]) != 1 {
		return createInvalidResponse()
	}
	args := Args{}
	args.append(t.inputRequest.query["thread"][0])

	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("user", t.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("forum", t.inputRequest.query["related"]) {
		relatedForum = true
	}
	if len(t.inputRequest.query["related"]) >= 1 && stringInSlice("thread", t.inputRequest.query["related"]) {
		return createInvalidQuery()
	}

	responseCode, responseMsg := t._getThreadDetails(args)

	if relatedUser {
		u := User{inputRequest: t.inputRequest, db: t.db}
		clearQuery(&u.inputRequest.query)
		args := Args{}
		args.append(responseMsg.User)

		_, userDetails := u._getUserDetails(args)

		responseMsg.User = userDetails
	}

	if relatedForum {
		f := Forum{inputRequest: t.inputRequest, db: t.db}
		clearQuery(&f.inputRequest.query)
		args := Args{}
		args.append(responseMsg.Forum)

		_, forumDetails := f._getForumDetails(args)
		responseMsg.Forum = forumDetails
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (t *Thread) listBasic() (int, *rs.ThreadList) {
	var query string
	args := Args{}

	// Validate query values
	if len(t.inputRequest.query["user"]) == 1 {
		query = "SELECT t.* FROM thread t WHERE t.user = ?"
		args.append(t.inputRequest.query["user"][0])
	} else if len(t.inputRequest.query["forum"]) == 1 {
		query = "SELECT t.* FROM thread t WHERE t.forum = ?"
		args.append(t.inputRequest.query["forum"][0])
	} else {
		return 100500, nil
	}

	// Check and validate optional params
	if len(t.inputRequest.query["since"]) >= 1 {
		query += " AND t.date > ?"
		args.append(t.inputRequest.query["since"][0])
	}

	if len(t.inputRequest.query["order"]) >= 1 {
		orderType := t.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			return 100500, nil
		}

		query += fmt.Sprintf(" ORDER BY t.date %s", orderType)
	}

	if len(t.inputRequest.query["limit"]) >= 1 {
		limitValue := t.inputRequest.query["limit"][0]
		i, err := strconv.Atoi(limitValue)
		if err != nil || i < 0 {
			return 100500, nil
		}

		query += fmt.Sprintf(" LIMIT %d", i)
	}

	// Response here
	responseCode, responseMsg := t._getArrayThreadsDetails(query, args)

	return responseCode, responseMsg
}

// Rewrite subquery
// +
func (t *Thread) list() string {
	responseCode, responseMsg := t.listBasic()

	if responseCode != 0 {
		return becauseAPI()
	}

	responseInterface := make([]interface{}, len(responseMsg.Threads))
	for i, v := range responseMsg.Threads {
		responseInterface[i] = v
	}

	return _createResponseFromArray(responseCode, responseInterface)
}

// +
func (t *Thread) parentTree(order string) (int, *rs.PostList) {
	query := "SELECT parent FROM post WHERE thread = ?"
	order = " ORDER BY parent " + order
	args := Args{}

	args.append(t.inputRequest.query["thread"][0])

	// Check and validate optional params
	if len(t.inputRequest.query["since"]) >= 1 {
		query += " AND date > ?"
		args.append(t.inputRequest.query["since"][0])
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

	getPost := selectQuery(query, &args.data, t.db)

	if getPost.rows == 0 {
		responseCode := 1
		errorMessage := &rs.PostList{}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseArray := make([]rs.PostDetails, 0)
	responseMsg := &rs.PostList{Posts: responseArray}

	for _, value := range getPost.values {
		subQuery := "SELECT * FROM post WHERE thread = ? AND parent LIKE ? ORDER BY parent"
		subArgs := Args{}
		subArgs.append(t.inputRequest.query["thread"][0])
		subArgs.append(value["parent"] + "%")

		getSubPost := selectQuery(subQuery, &subArgs.data, t.db)

		for _, subValue := range getSubPost.values {
			respId := stringToInt64(subValue["id"])

			tempMsg := &rs.PostDetails{
				Date:          subValue["date"],
				Dislikes:      stringToInt64(subValue["dislikes"]),
				Forum:         subValue["forum"],
				Id:            respId,
				IsApproved:    stringToBool(subValue["isApproved"]),
				IsHighlighted: stringToBool(subValue["isHighlighted"]),
				IsEdited:      stringToBool(subValue["isEdited"]),
				IsSpam:        stringToBool(subValue["isSpam"]),
				IsDeleted:     stringToBool(subValue["isDeleted"]),
				Likes:         stringToInt64(subValue["likes"]),
				Message:       subValue["message"],
				Parent:        nil,
				Points:        stringToInt64(subValue["points"]),
				Thread:        stringToInt64(subValue["thread"]),
				User:          subValue["user"],
			}

			p := Post{inputRequest: t.inputRequest, db: t.db}
			parent := p.getParentId(respId, subValue["parent"])
			if parent == int(respId) {
				tempMsg.Parent = nil
			} else {
				tempParent := int64(parent)
				tempMsg.Parent = &tempParent
			}

			responseMsg.Posts = append(responseMsg.Posts, *tempMsg)
		}
	}

	fmt.Println("Thread.parentTree()")

	return responseCode, responseMsg
}

// +
func (t *Thread) listPosts() string {
	var query, order, sort, resp string
	parentTree := false
	args := Args{}

	// Validate query values
	if len(t.inputRequest.query["thread"]) == 1 {
		query = "SELECT * FROM post WHERE thread = ?"
		args.append(t.inputRequest.query["thread"][0])
	} else {
		return createInvalidResponse()
	}

	// order by here
	if len(t.inputRequest.query["order"]) >= 1 {
		orderType := t.inputRequest.query["order"][0]
		if orderType != "desc" && orderType != "asc" {
			return createInvalidResponse()
		}
		order = orderType
	} else {
		order = "DESC"
	}

	// sort here
	if len(t.inputRequest.query["sort"]) >= 1 {
		switch sortType := t.inputRequest.query["sort"][0]; sortType {
		case "flat":
			sort = " ORDER BY date " + order
		case "tree":
			sort = " ORDER BY SUBSTRING(parent, 1, 5) " + order + ", parent asc "
		case "parent_tree":
			parentTree = true
		default:
			return createInvalidResponse()
		}
	} else {
		sort = " ORDER BY date " + order
	}

	var responseCode int
	responseMsg := &rs.PostList{}

	// simple sort
	if parentTree == false {
		p := Post{inputRequest: t.inputRequest, db: t.db}
		responseCode, responseMsg = p._getList(query, sort, args)
	} else {
		// parent_tree sort
		responseCode, responseMsg = t.parentTree(order)
	}

	// check responseCode
	if responseCode == 0 {
		responseInterface := make([]interface{}, len(responseMsg.Posts))
		for i, v := range responseMsg.Posts {
			responseInterface[i] = v
		}
		resp = _createResponseFromArray(responseCode, responseInterface)
	} else if responseCode == 1 {
		resp = becauseAPI()
	} else {
		resp = createInvalidResponse()
	}

	return resp
}

// +
func (t *Thread) open() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	return t.updateBoolBasic(query, false)
}

// +
func (t *Thread) remove() string {
	query := "UPDATE thread SET isDeleted = ?, posts = 0 WHERE id = ?"

	resp := t.updateBoolBasic(query, true)

	query = "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, true)

	return resp
}

// +
func (t *Thread) restore() string {
	query := "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, false)

	query = "UPDATE thread t SET t.isDeleted = ?, t.posts = (SELECT COUNT(*) FROM post p WHERE p.thread = t.id AND p.isDeleted = false) WHERE t.id = ?"

	return t.updateBoolBasic(query, false)
}

// +
func (t *Thread) subscribe() string {
	var resp string
	args := Args{}

	query := "INSERT INTO subscribe (thread, user) VALUES(?, ?)"

	if !validateJson(t.inputRequest, "thread", "user") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	if !checkFloat64Type(t.inputRequest.json["thread"]) {
		return createInvalidJsonResponse(t.inputRequest)
	}

	args.generateFromJson(&t.inputRequest.json, "thread", "user")

	_, err := execQuery(query, &args.data, t.db)
	if err != nil {
		fmt.Println(err)

		// return exist
		if checkError1062(err) == true {
			clearQuery(&t.inputRequest.query)
			t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], floatToString(t.inputRequest.json["thread"].(float64)))

			return t.details()
		}

		return createErrorResponse(err)
	}

	// else return info

	responseCode := 0
	responseMsg := &rs.ThreadSubscribe{
		Thread: int64(t.inputRequest.json["thread"].(float64)),
		User:   t.inputRequest.json["user"].(string),
	}

	resp = _createResponse(responseCode, responseMsg)

	log.Printf("User '%s' subscribe to thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

// +
func (t *Thread) unsubscribe() string {
	var resp string
	args := Args{}

	query := "DELETE FROM subscribe WHERE thread = ? AND user = ?"

	if !validateJson(t.inputRequest, "thread", "user") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	if !checkFloat64Type(t.inputRequest.json["thread"]) {
		return createInvalidJsonResponse(t.inputRequest)
	}

	args.generateFromJson(&t.inputRequest.json, "thread", "user")

	dbResp, err := execQuery(query, &args.data, t.db)
	if err != nil {
		return createErrorResponse(err)
	}

	if dbResp.rowCount == 0 {
		clearQuery(&t.inputRequest.query)
		t.inputRequest.query["thread"] = append(t.inputRequest.query["thread"], floatToString(t.inputRequest.json["thread"].(float64)))

		return t.details()
	}

	responseCode := 0
	responseMsg := &rs.ThreadSubscribe{
		Thread: int64(t.inputRequest.json["thread"].(float64)),
		User:   t.inputRequest.json["user"].(string),
	}

	resp = _createResponse(responseCode, responseMsg)

	log.Printf("User '%s' unsubscribe from thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

// +
func (t *Thread) update() string {
	args := Args{}

	query := "UPDATE thread SET message = ?, slug = ? WHERE id = ?"

	if !validateJson(t.inputRequest, "thread", "message", "slug") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	if !checkFloat64Type(t.inputRequest.json["thread"]) {
		return createInvalidJsonResponse(t.inputRequest)
	}

	threadId := t.inputRequest.json["thread"].(float64)

	args.generateFromJson(&t.inputRequest.json, "message", "slug", "thread")

	_, err := execQuery(query, &args.data, t.db)
	if err != nil {
		return createErrorResponse(err)
	}

	threadArgs := Args{}
	threadArgs.append(threadId)
	responseCode, responseMsg := t._getThreadDetails(threadArgs)

	if responseCode != 0 {
		return createNotExistResponse()
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (t *Thread) vote() string {
	var query string
	args := Args{}

	if !validateJson(t.inputRequest, "thread", "vote") {
		return createInvalidJsonResponse(t.inputRequest)
	}

	if !checkFloat64Type(t.inputRequest.json["thread"]) || !checkFloat64Type(t.inputRequest.json["vote"]) {
		return createInvalidJsonResponse(t.inputRequest)
	}

	threadId := t.inputRequest.json["thread"].(float64)
	vote := t.inputRequest.json["vote"].(float64)

	if vote == 1 {
		query = "UPDATE thread SET likes = likes + 1, points = points + 1 WHERE id = ?"
	} else if vote == -1 {
		query = "UPDATE thread SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?"
	} else {
		return createInvalidJsonResponse(t.inputRequest)
	}

	args.append(threadId)

	_, err := execQuery(query, &args.data, t.db)
	if err != nil {
		return createErrorResponse(err)
	}

	threadArgs := Args{}
	threadArgs.append(threadId)
	responseCode, responseMsg := t._getThreadDetails(threadArgs)

	if responseCode != 0 {
		return createNotExistResponse()
	}

	return _createResponse(responseCode, responseMsg)
}

func threadHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	//t0 := time.Now()
	thread := Thread{inputRequest: inputRequest, db: db}
	var result string

	if inputRequest.method == "GET" {
		switch inputRequest.path {
		case "/db/api/thread/details/":
			result = thread.details()
		case "/db/api/thread/list/":
			result = thread.list()
		case "/db/api/thread/listPosts/":
			result = thread.listPosts()
		}
	} else if inputRequest.method == "POST" {
		switch inputRequest.path {
		case "/db/api/thread/create/":
			result = thread.create()
		case "/db/api/thread/close/":
			result = thread.close()
		case "/db/api/thread/restore/":
			result = thread.restore()
		case "/db/api/thread/vote/":
			result = thread.vote()
		case "/db/api/thread/remove/":
			result = thread.remove()
		case "/db/api/thread/open/":
			result = thread.open()
		case "/db/api/thread/update/":
			result = thread.update()
		case "/db/api/thread/subscribe/":
			result = thread.subscribe()
		case "/db/api/thread/unsubscribe/":
			result = thread.unsubscribe()
		}
	}

	// t1 := time.Now()
	// fmt.Printf("Thread handler: %v\n", t1.Sub(t0))
	io.WriteString(w, result)
}

// =================
// Post handler here
// =================
type Post struct {
	inputRequest *InputRequest
	db           *sql.DB
}

func (p *Post) threadCounter(operation string, thread string, isDeleted bool) {
	var query string
	args := Args{}
	args.append(thread)

	switch operation {
	case "create":
		if isDeleted == false {
			query = "UPDATE thread SET posts = posts + 1 WHERE id = ?"
		}
	case "remove":
		query = "UPDATE thread SET posts = posts - 1 WHERE id = ?"
	case "restore":
		query = "UPDATE thread SET posts = posts + 1 WHERE id = ?"

	}
	_, _ = execQuery(query, &args.data, p.db)
}

// +
func (p *Post) create() string {
	args := Args{}

	query := "INSERT INTO post (thread, message, user, forum, date, isApproved, isHighlighted, isEdited, isSpam, isDeleted, parent) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	if !validateJson(p.inputRequest, "thread", "message", "user", "forum", "date") {
		return createInvalidJsonResponse(p.inputRequest)
	}

	args.generateFromJson(&p.inputRequest.json, "thread", "message", "user", "forum", "date")

	if !validateBoolParams(p.inputRequest.json, &args, "isApproved", "isHighlighted", "isEdited", "isSpam", "isDeleted") {
		return createInvalidJsonResponse(p.inputRequest)
	}
	var boolParent bool

	// parent here
	if p.inputRequest.json["parent"] != nil {
		if checkFloat64Type(p.inputRequest.json["parent"]) == false {
			return createInvalidResponse()
		}

		// find parent and last child
		parent := p.inputRequest.json["parent"].(float64)

		fmt.Println("Need parent:\t", parent)

		parentQuery := "SELECT id, parent FROM post WHERE id = ?"

		parentArgs := Args{}
		parentArgs.append(parent)

		getThread := selectQuery(parentQuery, &parentArgs.data, p.db)

		// check query
		if getThread.rows == 0 {
			return createNotExistResponse()
		}

		parentArgs.clear()

		// search place for child
		var child string
		getParent := getThread.values[0]["parent"]

		parentQuery = "SELECT parent FROM post WHERE parent LIKE ? ORDER BY parent desc LIMIT 1"
		parentArgs.append(getParent + "%")

		getThread = selectQuery(parentQuery, &parentArgs.data, p.db)

		if getThread.values[0]["parent"] == getParent {
			newParent := getParent
			newChild := toBase92(1)

			child = newParent + newChild
		} else {
			lastChild := getThread.values[0]["parent"]
			newParent := getParent
			oldChild := fromBase92(lastChild[len(lastChild)-5:])

			oldChild++

			newChild := toBase92(oldChild)

			child = newParent + newChild
		}

		args.append(child)
		boolParent = false
	} else {
		args.append(nil)
		boolParent = true
	}

	dbResp, err := execQuery(query, &args.data, p.db)
	if err != nil {
		return createErrorResponse(err)
	}

	fmt.Println("Create post with id:\t", dbResp.lastId)

	// maybe goroutine?
	if boolParent {
		boolParentQuery := "UPDATE post SET parent = ? WHERE id = ?"
		boolParentArgs := Args{}

		parent := toBase92(int(dbResp.lastId))
		boolParentArgs.append(parent)
		boolParentArgs.append(dbResp.lastId)

		_, _ = execQuery(boolParentQuery, &boolParentArgs.data, p.db)
	}

	// thread + isDeleted
	isDeleted := p.inputRequest.json["isDeleted"].(bool)
	thread := floatToString(p.inputRequest.json["thread"].(float64))
	p.threadCounter("create", thread, isDeleted)

	tempCounter := 5
	responseCode := 0
	responseMsg := &rs.PostCreate{
		Date:          p.inputRequest.json["date"].(string),
		Forum:         p.inputRequest.json["forum"].(string),
		Id:            dbResp.lastId,
		IsApproved:    args.data[incInt(&tempCounter)].(bool),
		IsHighlighted: args.data[tempCounter+0].(bool),
		IsEdited:      args.data[tempCounter+1].(bool),
		IsSpam:        args.data[tempCounter+2].(bool),
		IsDeleted:     args.data[tempCounter+3].(bool),
		Message:       p.inputRequest.json["message"].(string),
		Parent:        nil,
		Thread:        p.inputRequest.json["thread"].(float64),
		User:          p.inputRequest.json["user"].(string),
	}

	if !boolParent {
		tempParent := floatToString(p.inputRequest.json["parent"].(float64))
		responseMsg.Parent = &tempParent
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (p *Post) getParentId(id int64, path string) int {
	if len(path) == 5 {
		return int(id)
	} else if len(path) == 10 {
		return fromBase92(path[len(path)-10 : len(path)-5])
	} else {
		parentId := path[:len(path)-5]

		query := "SELECT id FROM post WHERE parent = ?"
		args := Args{}
		args.append(parentId)

		getParent := selectQuery(query, &args.data, p.db)

		respId := stringToInt64(getParent.values[0]["id"])

		return int(respId)
	}
}

// +
// Return true if need threadCounter()
func (p *Post) updateBoolBasic(query string, value bool) (bool, string) {
	args := Args{}

	if !validateJson(p.inputRequest, "post") {
		return false, createInvalidJsonResponse(p.inputRequest)
	}

	if checkFloat64Type(p.inputRequest.json["post"]) == false {
		return false, createInvalidJsonResponse(p.inputRequest)
	}

	postId := p.inputRequest.json["post"].(float64)

	args.append(value, postId)

	dbResp, err := execQuery(query, &args.data, p.db)
	if err != nil {
		return false, createErrorResponse(err)
	}

	if dbResp.rowCount == 0 {
		postArgs := Args{}
		postArgs.append(postId)

		responseCode, responseMsg := p._getPostDetails(postArgs)

		if responseCode != 0 {
			return false, createNotExistResponse()
		}

		return false, _createResponse(responseCode, responseMsg)
	}

	responseCode := 0
	responseMsg := &rs.PostBoolBasic{
		Post: postId,
	}

	return true, _createResponse(responseCode, responseMsg)
}

// +
func (p *Post) _getPostDetails(args Args) (int, *rs.PostDetails) {
	query := "SELECT * FROM post WHERE id = ?"

	getPost := selectQuery(query, &args.data, p.db)

	if getPost.rows == 0 {
		responseCode := 1
		errorMessage := &rs.PostDetails{}

		return responseCode, errorMessage
	}

	respId := stringToInt64(getPost.values[0]["id"])

	responseCode := 0
	responseMsg := &rs.PostDetails{
		Date:          getPost.values[0]["date"],
		Dislikes:      stringToInt64(getPost.values[0]["dislikes"]),
		Forum:         getPost.values[0]["forum"],
		Id:            respId,
		IsApproved:    stringToBool(getPost.values[0]["isApproved"]),
		IsHighlighted: stringToBool(getPost.values[0]["isHighlighted"]),
		IsEdited:      stringToBool(getPost.values[0]["isEdited"]),
		IsSpam:        stringToBool(getPost.values[0]["isSpam"]),
		IsDeleted:     stringToBool(getPost.values[0]["isDeleted"]),
		Likes:         stringToInt64(getPost.values[0]["likes"]),
		Message:       getPost.values[0]["message"],
		Parent:        nil,
		Points:        stringToInt64(getPost.values[0]["points"]),
		Thread:        stringToInt64(getPost.values[0]["thread"]),
		User:          getPost.values[0]["user"],
	}

	if parent := p.getParentId(respId, getPost.values[0]["parent"]); parent == int(respId) {
		responseMsg.Parent = nil
	} else {
		tempParent := int64(parent)
		responseMsg.Parent = &tempParent
	}

	return responseCode, responseMsg
}

// +
func (p *Post) details() string {
	var relatedUser, relatedThread, relatedForum bool
	args := Args{}

	if len(p.inputRequest.query["post"]) != 1 {
		return createInvalidResponse()
	}
	args.append(p.inputRequest.query["post"][0])

	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("user", p.inputRequest.query["related"]) {
		relatedUser = true
	}
	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("thread", p.inputRequest.query["related"]) {
		relatedThread = true
	}
	if len(p.inputRequest.query["related"]) >= 1 && stringInSlice("forum", p.inputRequest.query["related"]) {
		relatedForum = true
	}

	responseCode, responseMsg := p._getPostDetails(args)
	if responseCode != 0 {
		return _createResponse(responseCode, responseMsg)
	}

	if relatedUser {
		u := User{inputRequest: p.inputRequest, db: p.db}

		userArgs := Args{}
		userArgs.append(responseMsg.User)
		_, userDetails := u._getUserDetails(userArgs)

		responseMsg.User = userDetails
	}

	if relatedThread {
		t := Thread{inputRequest: p.inputRequest, db: p.db}

		threadArgs := Args{}
		threadArgs.append(responseMsg.Thread)
		_, threadDetails := t._getThreadDetails(threadArgs)

		responseMsg.Thread = threadDetails
	}

	if relatedForum {
		f := Forum{inputRequest: p.inputRequest, db: p.db}

		forumArgs := Args{}
		forumArgs.append(responseMsg.Forum)
		_, forumDetails := f._getForumDetails(forumArgs)

		responseMsg.Forum = forumDetails
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (p *Post) _getArrayPostDetails(query string, args Args) (int, *rs.PostList) {
	getPost := selectQuery(query, &args.data, p.db)

	if getPost.rows == 0 {
		responseCode := 1
		errorMessage := &rs.PostList{}

		return responseCode, errorMessage
	}

	responseCode := 0
	responseArray := make([]rs.PostDetails, 0)
	responseMsg := &rs.PostList{Posts: responseArray}

	for _, value := range getPost.values {
		respId := stringToInt64(value["id"])

		tempMsg := &rs.PostDetails{
			Date:          value["date"],
			Dislikes:      stringToInt64(value["dislikes"]),
			Forum:         value["forum"],
			Id:            respId,
			IsApproved:    stringToBool(value["isApproved"]),
			IsHighlighted: stringToBool(value["isHighlighted"]),
			IsEdited:      stringToBool(value["isEdited"]),
			IsSpam:        stringToBool(value["isSpam"]),
			IsDeleted:     stringToBool(value["isDeleted"]),
			Likes:         stringToInt64(value["likes"]),
			Message:       value["message"],
			Parent:        nil,
			Points:        stringToInt64(value["points"]),
			Thread:        stringToInt64(value["thread"]),
			User:          value["user"],
		}

		parent := p.getParentId(respId, value["parent"])
		if parent == int(respId) {
			tempMsg.Parent = nil
		} else {
			tempParent := int64(parent)
			tempMsg.Parent = &tempParent
		}

		responseMsg.Posts = append(responseMsg.Posts, *tempMsg)
	}

	return responseCode, responseMsg
}

// +
func (p *Post) _getList(query string, order string, args Args) (int, *rs.PostList) {
	// Check and validate optional params
	if len(p.inputRequest.query["since"]) >= 1 {
		query += " AND date > ?"
		args.append(p.inputRequest.query["since"][0])
	}

	query += order

	if len(p.inputRequest.query["limit"]) >= 1 {
		i, err := strconv.Atoi(p.inputRequest.query["limit"][0])
		if err != nil || i < 0 {
			return 100500, nil
		}
		query += fmt.Sprintf(" LIMIT %d", i)
	}

	fmt.Println("_getList: ", query, args.data)

	responseCode, responseMsg := p._getArrayPostDetails(query, args)

	return responseCode, responseMsg
}

// +
func (p *Post) list() string {
	var query, order, resp string
	args := Args{}

	// Validate query values
	if len(p.inputRequest.query["thread"]) == 1 {
		query = "SELECT * FROM post WHERE thread = ?"
		args.append(p.inputRequest.query["thread"][0])
	} else if len(p.inputRequest.query["user"]) == 1 {
		query = "SELECT * FROM post WHERE user = ?"
		args.append(p.inputRequest.query["user"][0])
	} else if len(p.inputRequest.query["forum"]) == 1 {
		query = "SELECT * FROM post WHERE forum = ?"
		args.append(p.inputRequest.query["forum"][0])
	} else {
		return createInvalidResponse()
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

	responseCode, responseMsg := p._getList(query, order, args)

	// check responseCode
	if responseCode == 0 {
		responseInterface := make([]interface{}, len(responseMsg.Posts))
		for i, v := range responseMsg.Posts {
			responseInterface[i] = v
		}
		resp = _createResponseFromArray(responseCode, responseInterface)
	} else if responseCode == 1 {
		resp = becauseAPI()
	} else {
		resp = createInvalidResponse()
	}

	return resp
}

// +
func (p *Post) remove() string {
	query := "UPDATE post SET isDeleted = ? WHERE id = ?"

	check, resp := p.updateBoolBasic(query, true)
	if check {
		args := Args{}
		args.append(p.inputRequest.json["post"])
		_, responseMsg := p._getPostDetails(args)
		thread := int64ToString(responseMsg.Thread.(int64))
		p.threadCounter("remove", thread, true)
	}

	return resp
}

// +
func (p *Post) restore() string {
	query := "UPDATE post SET isDeleted = ? WHERE id = ?"

	check, resp := p.updateBoolBasic(query, false)
	if check {
		args := Args{}
		args.append(p.inputRequest.json["post"])
		_, responseMsg := p._getPostDetails(args)
		thread := int64ToString(responseMsg.Thread.(int64))
		p.threadCounter("restore", thread, false)
	}

	return resp
}

// +
func (p *Post) update() string {
	query := "UPDATE post SET message = ? WHERE id = ?"

	args := Args{}

	if !validateJson(p.inputRequest, "message", "post") {
		return createInvalidJsonResponse(p.inputRequest)
	}

	if !checkFloat64Type(p.inputRequest.json["post"]) {
		return createInvalidJsonResponse(p.inputRequest)
	}

	threadId := p.inputRequest.json["post"].(float64)
	args.generateFromJson(&p.inputRequest.json, "message", "post")

	dbResp, err := execQuery(query, &args.data, p.db)
	if err != nil {
		return createErrorResponse(err)
	}

	if dbResp.rowCount == 0 {
		query := "UPDATE post SET isEdited = false WHERE id = ?"
		_, _ = execQuery(query, &args.data, p.db)
	} else {
		query := "UPDATE post SET isEdited = true WHERE id = ?"
		_, _ = execQuery(query, &args.data, p.db)
	}

	postArgs := Args{}
	postArgs.append(threadId)
	responseCode, responseMsg := p._getPostDetails(postArgs)

	if responseCode != 0 {
		return createNotExistResponse()
	}

	return _createResponse(responseCode, responseMsg)
}

// +
func (p *Post) vote() string {
	var query string
	args := Args{}

	if !validateJson(p.inputRequest, "post", "vote") {
		return createInvalidJsonResponse(p.inputRequest)
	}

	if !checkFloat64Type(p.inputRequest.json["post"]) || !checkFloat64Type(p.inputRequest.json["vote"]) {
		return createInvalidJsonResponse(p.inputRequest)
	}

	postId := p.inputRequest.json["post"].(float64)
	vote := p.inputRequest.json["vote"].(float64)

	if vote == 1 {
		query = "UPDATE post SET likes = likes + 1, points = points + 1 WHERE id = ?"
	} else if vote == -1 {
		query = "UPDATE post SET dislikes = dislikes + 1, points = points - 1 WHERE id = ?"
	} else {
		return createInvalidJsonResponse(p.inputRequest)
	}

	args.append(p.inputRequest.json["post"])

	_, err := execQuery(query, &args.data, p.db)
	if err != nil {
		return createErrorResponse(err)
	}

	postArgs := Args{}
	postArgs.append(postId)
	responseCode, responseMsg := p._getPostDetails(postArgs)

	if responseCode != 0 {
		return createNotExistResponse()
	}

	return _createResponse(responseCode, responseMsg)
}

func postHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	//t0 := time.Now()
	post := Post{inputRequest: inputRequest, db: db}
	var result string

	if inputRequest.method == "GET" {
		switch inputRequest.path {
		case "/db/api/post/details/":
			result = post.details()
		case "/db/api/post/list/":
			result = post.list()
		}
	} else if inputRequest.method == "POST" {
		switch inputRequest.path {
		case "/db/api/post/create/":
			result = post.create()
		case "/db/api/post/restore/":
			result = post.restore()
		case "/db/api/post/vote/":
			result = post.vote()
		case "/db/api/post/remove/":
			result = post.remove()
		case "/db/api/post/update/":
			result = post.update()
		}
	}

	// t1 := time.Now()
	// fmt.Printf("Post handler: %v\n", t1.Sub(t0))
	io.WriteString(w, result)
}

// ==========================
// Information methods here
// ==========================
func statusHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	if inputRequest.method == "GET" {
		args := Args{}

		query := "SELECT COUNT(*) count FROM user"
		dbResp := selectQuery(query, &args.data, db)
		respUsers := stringToInt64(dbResp.values[0]["count"])

		query = "SELECT COUNT(*) count FROM thread"
		dbResp = selectQuery(query, &args.data, db)
		respThreads := stringToInt64(dbResp.values[0]["count"])

		query = "SELECT COUNT(*) count FROM forum"
		dbResp = selectQuery(query, &args.data, db)
		respForums := stringToInt64(dbResp.values[0]["count"])

		query = "SELECT COUNT(*) count FROM post"
		dbResp = selectQuery(query, &args.data, db)
		respPosts := stringToInt64(dbResp.values[0]["count"])

		responseCode := 0
		responseMsg := &rs.StatusHandler{
			User:   respUsers,
			Thread: respThreads,
			Forum:  respForums,
			Post:   respPosts,
		}

		io.WriteString(w, _createResponse(responseCode, responseMsg))
	}
}

func clearHandler(w http.ResponseWriter, r *http.Request, inputRequest *InputRequest, db *sql.DB) {
	if inputRequest.method == "POST" {
		args := Args{}

		query := "DELETE FROM follow"
		_, _ = execQuery(query, &args.data, db)
		query = "DELETE FROM subscribe"
		_, _ = execQuery(query, &args.data, db)
		query = "DELETE FROM post WHERE id > 0"
		_, _ = execQuery(query, &args.data, db)
		query = "DELETE FROM thread WHERE id > 0"
		_, _ = execQuery(query, &args.data, db)
		query = "DELETE FROM forum WHERE id > 0"
		_, _ = execQuery(query, &args.data, db)
		query = "DELETE FROM user WHERE id > 0"
		_, _ = execQuery(query, &args.data, db)

		responseCode := 0

		cacheContent := &rs.ClearHandler{
			Code:     int64(responseCode),
			Response: "OK",
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
		//t0 := time.Now()
		inputRequest := new(InputRequest)
		inputRequest.parse(r)
		// t1 := time.Now()
		// fmt.Printf("Parse time: %v\n", t1.Sub(t0))

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

	// config here
	f, err := os.Open("app.conf")
	check(err)
	conf := make([]byte, 10)
	_, err = f.Read(conf)
	check(err)

	if strings.Split(string(conf), "=")[1][:1] == "1" {
		http.HandleFunc("/db/api/clear/", makeHandler(db, clearHandler))
	}

	http.HandleFunc("/db/api/user/", makeHandler(db, userHandler))
	http.HandleFunc("/db/api/forum/", makeHandler(db, forumHandler))
	http.HandleFunc("/db/api/thread/", makeHandler(db, threadHandler))
	http.HandleFunc("/db/api/post/", makeHandler(db, postHandler))
	http.HandleFunc("/db/api/status/", makeHandler(db, statusHandler))

	http.ListenAndServe(PORT, nil)
}

// =================
// Utils here
// =================

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func floatToString(inputNum float64) string { return strconv.FormatFloat(inputNum, 'f', 6, 64) }

func int64ToString(inputNum int64) string { return strconv.FormatInt(inputNum, 10) }

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

func stringToInt64(inputStr string) int64 { return int64(stringToInt(inputStr)) }

func incInt(value *int) int { return *value + 1 }

func toBase92(value int) string {
	BASE92 := ` !"#$&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`
	length := 5
	base := len(BASE92)
	result := make([]byte, 5)

	for i := 0; i < length; i++ {
		result[i] = BASE92[0]
	}

	if value == 0 {
		return string(result)
	}

	counter := 0
	for value != 0 {
		mod := value % base

		result[length-counter-1] = BASE92[mod]
		counter++

		value = value / base
	}

	return string(result)
}

func fromBase92(value string) int {
	BASE92 := ` !"#$&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`
	length := 5
	base := len(BASE92)
	var result int

	counter := 0
	step := 1
	for i := 0; i < length; i++ {
		index := strings.Index(BASE92, string(value[length-counter-1]))
		counter++

		result = result + index*step

		step = step * base
	}

	return result
}

func stringToBool(inputString string) bool {
	result, err := strconv.ParseBool(inputString)

	if err != nil {
		log.Panic(err)
	}

	return result
}
