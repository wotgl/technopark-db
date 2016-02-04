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
// Thread handler here
// =================

type Thread struct {
	inputRequest *InputRequest
	db           *sql.DB
}

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

		return createResponse(responseCode, responseMsg)
	}

	responseCode := 0
	responseMsg := &rs.ThreadBoolBasic{
		Thread: threadId,
	}

	return createResponse(responseCode, responseMsg)
}

func (t *Thread) close() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	return t.updateBoolBasic(query, true)
}

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

	resp = createResponse(responseCode, responseMsg)

	log.Printf("Thread '#%d' created", responseMsg.Id)

	return resp
}

// Rewrite subquery
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

	return createResponse(responseCode, responseMsg)
}

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
func (t *Thread) list() string {
	responseCode, responseMsg := t.listBasic()

	if responseCode != 0 {
		return becauseAPI()
	}

	responseInterface := make([]interface{}, len(responseMsg.Threads))
	for i, v := range responseMsg.Threads {
		responseInterface[i] = v
	}

	return createResponseFromArray(responseCode, responseInterface)
}

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
		resp = createResponseFromArray(responseCode, responseInterface)
	} else if responseCode == 1 {
		resp = becauseAPI()
	} else {
		resp = createInvalidResponse()
	}

	return resp
}

func (t *Thread) open() string {
	query := "UPDATE thread SET isClosed = ? WHERE id = ?"

	return t.updateBoolBasic(query, false)
}

func (t *Thread) remove() string {
	query := "UPDATE thread SET isDeleted = ?, posts = 0 WHERE id = ?"

	resp := t.updateBoolBasic(query, true)

	query = "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, true)

	return resp
}

func (t *Thread) restore() string {
	query := "UPDATE post SET isDeleted = ? WHERE thread = ?"

	_ = t.updateBoolBasic(query, false)

	query = "UPDATE thread t SET t.isDeleted = ?, t.posts = (SELECT COUNT(*) FROM post p WHERE p.thread = t.id AND p.isDeleted = false) WHERE t.id = ?"

	return t.updateBoolBasic(query, false)
}

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

	resp = createResponse(responseCode, responseMsg)

	log.Printf("User '%s' subscribe to thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

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

	resp = createResponse(responseCode, responseMsg)

	log.Printf("User '%s' unsubscribe from thread '#%d'", responseMsg.User, responseMsg.Thread)

	return resp
}

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

	return createResponse(responseCode, responseMsg)
}

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

	return createResponse(responseCode, responseMsg)
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
