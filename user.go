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
// User handler here
// =================

type User struct {
	inputRequest *InputRequest
	db           *sql.DB
}

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

	resp = createResponse(responseCode, responseMsg)

	log.Printf("User '%s' created", responseMsg.Email)

	return resp
}

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

	return createResponse(responseCode, responseMsg)
}

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

		return createResponse(responseCode, errorMessage)
	}

	u.inputRequest.query["user"] = append(u.inputRequest.query["user"], u.inputRequest.json["follower"].(string))
	return u.getDetails()
}

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

	return createResponseFromArray(responseCode, responseInterface)
}

func (u *User) listFollowers() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.follower WHERE followee = ?"

	return u.listBasic(query)
}

func (u *User) listFollowing() string {
	query := "SELECT u.* FROM user u JOIN follow f ON u.email = f.followee WHERE follower = ?"

	return u.listBasic(query)
}

func (u *User) listPosts() string {
	delete(u.inputRequest.query, "forum")

	p := Post{inputRequest: u.inputRequest, db: u.db}
	return p.list()
}

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
