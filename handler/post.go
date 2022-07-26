package handler // Handler matches controller in OnlineOrder

import (
	// Go packages
	"encoding/json" // for the encoding and decoding of JSON，equivalent to Jackson in OnlineOrder
	"fmt"
	"net/http"

	// our own package
	"around/model"
	"around/service"
)

/* process user uploads
 * if a user sends a HTTP request with a body as:
 * {
 *     "user": "tiffany",
 *     "message": "this is a post from tiffany",
 * }
 * it will automatically construct a Post object p
 * and update its values to be p.User = "tiffany" and p.Message = "this is a post from tiffany"
 */

// Objective: parse from body of request to get a JSON object
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	/* r *http.Request is a reference, it takes up less space than http.Request (type Request is a struct). I.e., needs 1kB space to copy a request vs. a few bytes of reference address) -> program runs faster by not making a copy of the Request */

	/* ResponseWriter is an Go interface, you can declare but not implement, you can't "new" a ResponseWriter object because it's an interface. Since you cannot implement it, it does not support pointers. */

	fmt.Println("Received one upload request") // helps with debugging

	// read the Request Body to the "decoder" object
	decoder := json.NewDecoder(r.Body)

	// (equivalent to Post p = new Post() in Java) declare a Post object p and send it to the Decode method to make data (json format) into a Go object
	var p model.Post

	// use the json Decode method to convert the Request Body to a Post object
	// &p because we need to make changes to the p object
	// if decode fails and error is not nil -> panic (throw runtime exception) -> program crashes and restarts
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	// Fprintf = don't print to console, print content to ResponseWriter w (response buffer that will stream the message data to the response body -> mechanism to catch any exceptions, such as message is to long resulting in the program to crash)
	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts) // Convert Go data into JSON
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}
