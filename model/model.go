package model

// Postman test: Around Backend Upload with JSON
type Post struct {
	// capital 'Id' = public
	// lowercase 'id' = private
	// `` = raw string
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}
