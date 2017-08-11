package main

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"github.com/gorilla/mux"
)

/**
 * [Goriila mux router]
 * @return {[type]} [request response route handler]
 */
func main() {
	rtr := mux.NewRouter()
	str := rtr.PathPrefix("/api").Subrouter()
	str.HandleFunc("/", HandleIndex)
	http.Handle("/", rtr)

	log.Println(http.ListenAndServe(":3011", nil))
}

/**
 * [HandleIndex description]
 * @param {[type]} w http.ResponseWriter [description]
 * @param {[type]} r *http.Request       [description]
 */
func HandleIndex(w http.ResponseWriter, r *http.Request) {	
	fmt.Println("Hi Aircto")
	io.WriteString(w, "Hi Aircto!")
}