package main

import (
	"fmt"
	"net/http"
	"os"
	"aricto/aricto"
	"aricto/middleware"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

/**
 * [Goriila mux router]
 * @return {[type]} [request response route handler]
 */
func main() {
	rtr := mux.NewRouter()

	str := rtr.PathPrefix("/api").Subrouter()
	str.HandleFunc("/", aricto.HandleIndex)
	str.HandleFunc("/login", aricto.PostLogin).Methods("POST")

	//user subroute api
	usr_str := str.PathPrefix("/user").Subrouter()
	usr_str.HandleFunc("/all-user-list", aricto.GetAllUserList).Methods("GET")

	//issue subroute api
	isu_str := str.PathPrefix("/issues").Subrouter()
	isu_str.Handle("/all-issues-list", middleware.JwtMiddleware(aricto.GetAllIssuesList)).Methods("GET")
	isu_str.Handle("/issue-info", middleware.JwtMiddleware(aricto.GetIssueInfo)).Methods("GET")
	isu_str.Handle("/create-issue", middleware.JwtMiddleware(aricto.CreateIssue)).Methods("POST")
	isu_str.Handle("/update-issue", middleware.JwtMiddleware(aricto.UpdateIssue)).Methods("PUT")
	isu_str.Handle("/delete-issue", middleware.JwtMiddleware(aricto.DeleteIssue)).Methods("DELETE")
	isu_str.Handle("/issues-by-me", middleware.JwtMiddleware(aricto.GetAllIssuesByMe)).Methods("GET")
	isu_str.Handle("/issues-for-me", middleware.JwtMiddleware(aricto.GetAllIssuesAssignedToMe)).Methods("GET")

	http.Handle("/", rtr)

	fmt.Println("********************************************************************")
	fmt.Println("*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*")
	fmt.Println("#                                                                  #")
	fmt.Println("#                                                                  #")
	fmt.Println("#                                                                  #")	
	fmt.Println("********** AirCTO SIT API - Hit http://localhost:3011/api **********")
	fmt.Println("#                                                                  #")
	fmt.Println("#                                                                  #")
	fmt.Println("#                                                                  #")	
	fmt.Println("*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*")
	fmt.Println("********************************************************************")

	http.ListenAndServe(":3011", handlers.LoggingHandler(os.Stdout, rtr))
}
