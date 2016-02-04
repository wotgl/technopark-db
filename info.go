package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	rs "technopark-db/response"
)

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

		io.WriteString(w, createResponse(responseCode, responseMsg))
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
