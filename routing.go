package gcpfunctions

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
	logger.InfoData("Added Route", route.LogData())
}

// Serve handle a request and using the router redirects the traffic
func Serve(w http.ResponseWriter, r *http.Request) {

	http.DefaultServeMux = new(http.ServeMux)
	router := mux.NewRouter()

	for _, route := range routes {
		router.HandleFunc(route.path, route.handler).Methods(route.verb)
		logger.InfoData("Initialized Route", route.LogData())
	}
	http.Handle("/", router)

	router.ServeHTTP(w, r)
}

// GetRequest from the HTTP request
func GetRequest(r *http.Request) (Request, error) {

	body, err := readBytes(r.Body)
	if err != nil {
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
		return nil, err
	}

	return buffer.Bytes(), nil
}
