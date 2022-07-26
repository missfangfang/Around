/* Unlike Spring, Go cannot support annotation-based HTTP routing, so we need to define our own method to map request URL to HTTP handler */

package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// uploadHandler will trigger only when POST request is sent to "/upload" url
	router.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST")

	return router
}
