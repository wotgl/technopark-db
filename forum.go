package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	rs "technopark-db/response"
)

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

	resp = createResponse(responseCode, responseMsg)

	log.Printf("Forum '%s' created", responseMsg.Short_Name)

	return resp
}

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

	return createResponse(responseCode, responseMsg)
}

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

	return createResponseFromArray(responseCode, responseInterface)
}

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

	return createResponseFromArray(responseCode, responseInterface)
}

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

	return createResponseFromArray(responseCode, responseInterface)
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
