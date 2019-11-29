package nozzle

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// var router = mux.NewRouter()
var routes = []Route{}

// Handler is a delegate to concrete handler
type Handler func(w http.ResponseWriter, r *http.Request)

// Route to be matched for the url
type Route struct {
	path    string
	verb    string
	handler Handler
}

// LogData for logging
func (route *Route) LogData() map[string]string {

	return map[string]string{
		"path": route.path,
		"verb": route.verb,
	}
}

// AddRoute adds a specific route to the router
func AddRoute(verb string, path string, handler Handler) {

	route := Route{
		path:    path,
		verb:    verb,
		handler: handler,
	}

	routes = append(routes, route)
	logger.Info(LogEntry{Action: "Added Route", Map: route.LogData()})
}

// Serve handle a request and using the router redirects the traffic
func Serve(w http.ResponseWriter, r *http.Request) {

	router := mux.NewRouter()

	for _, route := range routes {
		router.HandleFunc(route.path, route.handler).Methods(route.verb)
		logger.Info(LogEntry{Action: "Initialized Route", Map: route.LogData()})
	}

	// Always handle OPTIONS, on any path
	// This will need improving if we want to handle differently per end-point
	router.Methods(http.MethodOptions).HandlerFunc(OptionsHandler)

	router.ServeHTTP(w, r)
}

// GetRequest from the HTTP request
func GetRequest(r *http.Request) (Request, error) {

	body, err := readBytes(r.Body)
	if err != nil {
		logger.Error("Reading Request", err)
		return Request{}, err
	}

	return Request{
		Params: mux.Vars(r),
		Body:   body,
	}, nil
}

func readBytes(reader io.Reader) ([]byte, error) {

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(reader)
	if err != nil {
		logger.Error("Reading Request", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}