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

// Handle processes a HTTP request
func Handle(w http.ResponseWriter, object Entity, sql string, params ...interface{}) {

	response := Response{}

	err := ReturnOne(object, sql, params...)

	if err != nil {
		status, message := lookupError(err)
		response = newResponse(status, Error{Message: message})

	} else {
		response = newResponse(http.StatusOK, object)
	}

	writeResponse(w, response)
}

// HandleBadRequestErr handles a bad request from the client
func HandleBadRequestErr(w http.ResponseWriter, err error) {

	logger.Error("Bad Request", err)

	body := ErrorBody{
		Code:        http.StatusBadRequest,
		Description: "Bad Request",
	}

	writeResponse(w, newResponse(http.StatusBadRequest, body))
}

// NewResponse creates an initialized Response
func newResponse(code int, data interface{}) Response {

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET",
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

func writeResponse(w http.ResponseWriter, response Response) {

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

	errorResponse := newResponse(code, nil)
	writeResponse(w, errorResponse)
}

func lookupError(err error) (int, string) {

	if err != nil {
		switch err.(type) {
		case *pq.Error:
			return lookupDBError(err)
		default:
			return InternalError, "Internal Error"
		}
	}

	return InternalError, "Internal Error"
}
