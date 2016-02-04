package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"

	rs "technopark-db/response"
)

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

	tempCounter := 6
	responseCode := 0
	responseMsg := &rs.PostCreate{
		Date:          p.inputRequest.json["date"].(string),
		Forum:         p.inputRequest.json["forum"].(string),
		Id:            dbResp.lastId,
		IsApproved:    args.data[tempCounter].(bool),
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

	return createResponse(responseCode, responseMsg)
}

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

		return false, createResponse(responseCode, responseMsg)
	}

	responseCode := 0
	responseMsg := &rs.PostBoolBasic{
		Post: postId,
	}

	return true, createResponse(responseCode, responseMsg)
}

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
		return createResponse(responseCode, responseMsg)
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

	return createResponse(responseCode, responseMsg)
}

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
		resp = createResponseFromArray(responseCode, responseInterface)
	} else if responseCode == 1 {
		resp = becauseAPI()
	} else {
		resp = createInvalidResponse()
	}

	return resp
}

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

	return createResponse(responseCode, responseMsg)
}

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

	return createResponse(responseCode, responseMsg)
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
