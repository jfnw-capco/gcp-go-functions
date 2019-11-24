package nozzle

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

// Request is an abstract of the key request info
type Request struct {
	Params map[string]string
	Body   []byte
}

// Response is an abstract of the key response info
type Response struct {
	Code    int
	Headers map[string]string
	Data    interface{}
}

// Error object used to pass back msg
type Error struct {
	Message string `json:"id"`
}

// ErrorBody represents a structured error
type ErrorBody struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// Ping is provided by default to test the service is up
// and connectivity in place
func Ping(w http.ResponseWriter, r *http.Request) {

	logger.Debug(LogEntry{Action: "Ping Started"})

	_, err := pingDB()
	if err != nil {
		logger.Error("Service Ping Failed", err)
		handleErr(w, "Service Unavailable", http.StatusBadRequest, err)
	} else {
		logger.Debug(LogEntry{Action: "Ping Successful"})
		WriteResponse(w, NewResponse(http.StatusOK, nil))
	}
}

// Handle processes a HTTP request
func Handle(w http.ResponseWriter, object Entity, sql string, params ...interface{}) {

	response := Response{}

	err := ReturnOne(object, sql, params...)

	if err != nil {
		status, message := lookupError(err)
		response = NewResponse(status, Error{Message: message})

	} else {
		response = NewResponse(http.StatusOK, object)
	}

	WriteResponse(w, response)
}

// HandleBadRequestErr handles a bad request from the client
func HandleBadRequestErr(w http.ResponseWriter, err error) {

	handleErr(w, "Bad Request", http.StatusBadRequest, err)
}

func handleErr(w http.ResponseWriter, message string, code int, err error) {

	logger.Error(message, err)

	body := ErrorBody{
		Code:        code,
		Description: message,
	}

	WriteResponse(w, NewResponse(code, body))
}

// NewResponse creates an initialized Response
func NewResponse(code int, data interface{}) Response {

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, PATCH",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Max-Age":       "3600",
		"Content-Type":                 "application/json"}

	return Response{
		Headers: headers,
		Code:    code,
		Data:    data,
	}
}

// Serialize object to byte array
func Serialize(data interface{}) ([]byte, error) {

	json, err := json.Marshal(data)
	if err != nil {
		logger.Error("Serializing Body", err)
	}

	return json, err
}

// Deserialize JSON byte array to an object
func Deserialize(request Request, object interface{}) error {

	err := json.Unmarshal(request.Body, &object)
	if err != nil {
		logger.Error("Deserializing Body", err)
	}

	logger.Info(LogEntry{Action: "Deserialized Body", Message: string(request.Body)})
	return err
}

// WriteResponse writes a response in a consistent way
func WriteResponse(w http.ResponseWriter, response Response) {

	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(response.Code)

	logger.Debug(LogEntry{Action: "HTTP Headers", Map: response.Headers})
	logger.Debug(LogEntry{Action: "HTTP Code", Message: strconv.Itoa(response.Code)})

	if response.Data != nil {

		json, err := Serialize(response.Data)
		if err != nil {
			logger.Error("Serializing Body", err)
			WriteErrorToResponse(w, http.StatusInternalServerError)
		}

		logger.Debug(LogEntry{Action: "HTTP Body", Message: string(json)})

		_, err = w.Write(json)
		if err != nil {
			logger.Error("Writing Body", err)
			WriteErrorToResponse(w, http.StatusInternalServerError)
		}

		logger.Debug(LogEntry{Action: "HTTP Body", Message: string(json)})
	}
}

// WriteErrorToResponse writes in a consistent way
func WriteErrorToResponse(w http.ResponseWriter, code int) {

	errorResponse := NewResponse(code, nil)
	WriteResponse(w, errorResponse)
}

func lookupError(err error) (int, string) {
	if err != nil {
		switch err.(type) {
		case *pq.Error:
			logger.Error("Mapping Error To Db Error", err)
			return lookupDBError(err)
		default:
			logger.Error("Mapping Error To Internal Error", err)
			return InternalError, "Internal Error"
		}
	}
	logger.Error("Failed To Map Error", err)
	return InternalError, "Internal Error"
}
