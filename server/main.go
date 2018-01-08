package main

import (
	"net/http"
)


func PageHandler(w http.ResponseWriter,r *http.Request)  {
	w.Write([]byte("hammer time"))
}


func main() {
	http.HandleFunc("/",PageHandler)
	http.ListenAndServe(":6061", nil)
}