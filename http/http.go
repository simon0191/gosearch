package http

import (
	"net/http"
	"strconv"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gosearch/gosearch/service"
	"io"
	"io/ioutil"
	"log"
)

const indexPath = "/index"

// Server holds the configuration for the HTTP server.
type Server struct {
	Index service.IndexService
}

// Listen starts the http server on the given port.
func (server *Server) Listen(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/{index}/{id}", createIndex(server.Index)).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func createIndex(indexService service.IndexService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		index := vars["index"]
		id := vars["id"]
		data, err := bodyToJSON(r)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		indexService.Create(index, id, data)
		w.WriteHeader(http.StatusCreated)
	}
}

func bodyToJSON(r *http.Request) (interface{}, error) {
	var data interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}
