package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(res http.ResponseWriter, req *http.Request) { // <-- handler to handel response
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partial/navbar.tmpl.html",
	}
	ts,err := template.ParseFiles(files...)
	if err != nil{
		log.Println(err.Error())
		http.Error(res,"internal server error",http.StatusInternalServerError)
		return
	}
	error := ts.ExecuteTemplate(res,"base",nil)
	if error !=nil{
		log.Println(error.Error())
		http.Error(res,"internal server error",http.StatusInternalServerError)
	}
}

func snippetView(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.Header().Set("Allow", http.MethodGet)
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed) // <-- this works as both .write and .writeHeader
		return
	}
	qry, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || qry < 0 {
		http.NotFound(res,req)
		return
	}
	res.Header().Set("content-type", "application/json")         // <-- this does case sanitation
	res.Header()["X-XSS-PROTECTION"] = []string{"1; mode=block"} // <- if you dont want case sanitation
	fmt.Fprintf(res, "The number you query is %d...", qry)
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
