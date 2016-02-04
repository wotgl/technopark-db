package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	rs "technopark-db/response"

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
		log.Panic(err)
	}

	var parsed map[string]interface{}
	json.Unmarshal([]byte(body), &parsed)
	ir.json = parsed

	// GET Query
	ir.query = r.URL.Query()
}

func createResponse(code int, response rs.RespStruct) string {
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

func createResponseFromArray(code int, response []interface{}) string {
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

	return createResponse(responseCode, errorMessage)
}

func createInvalidQuery() string {
	responseCode := 3
	errorMessage := &rs.ErrorMsg{
		Msg: "Invalid query",
	}

	return createResponse(responseCode, errorMessage)
}

func createInvalidJsonResponse(inputRequest *InputRequest) string {
	responseCode := 3
	errorMessage := &rs.ErrorMsg{
		Msg: "Invalid json",
	}

	log.Println("Invalid JSON:\turl=\tjson=", inputRequest.url, inputRequest.json)

	return createResponse(responseCode, errorMessage)
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

	return createResponse(responseCode, errorMessage)
}

func createNotExistResponse() string {
	responseCode := 1
	errorMessage := &rs.ErrorMsg{
		Msg: "Not exist",
	}

	return createResponse(responseCode, errorMessage)
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
// Utils here
// =================

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
