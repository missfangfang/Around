/* Unlike Spring, Go cannot support annotation-based HTTP routing, so we need to define our own method to map request URL to HTTP handler */

package handler

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware" // Parse token
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()

	// uploadHandler will trigger only when POST request is sent to "/upload" url
	// Goes through jwtMiddleware first before going through handler
	router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
	router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	return router
}
