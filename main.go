package main

import (
	"log"
	"net/http"
)

// routes handeler
func home(res http.ResponseWriter, req *http.Request) { // <-- handler to handel response
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	res.Write([]byte("hello world"))
}
func snippetView(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.Header().Set("ALLOW", http.MethodGet)
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed) // <-- this works as both .write and .writeHeader
		return
	}
	res.Write([]byte("Snippets view..."))
}
func snippetCreate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("Allow", http.MethodPost) // <-- must be called before .Write for custom status code
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte("Method not allowed"))
		return
	}
	res.Write([]byte("Snippet created...")) // <- only this will send 200 status code
}
func main() {
	mux := http.NewServeMux() //<-- multiplexr to handel the respose
	mux.HandleFunc("/", home) // <-- router
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Println("Starting server at 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
