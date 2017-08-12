package main

import (
	"fmt"
	"net/http"
	"os"
	"aricto/aricto"
	"aricto/middleware"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)




/**
* Global Variables
*/
var db *sql.DB 
var err error
var mySigningSecretKey = []byte("qwerty123")


/**
 * [Goriila mux router]
 * @return {[type]} [request response route handler]
 */
func main() {
	rtr := mux.NewRouter()

	str := rtr.PathPrefix("/api").Subrouter()
	str.HandleFunc("/", aricto.HandleIndex)
	str.HandleFunc("/login", aricto.PostLogin).Methods("GET")

	//user subroute api
	usr_str := str.PathPrefix("/user").Subrouter()
	usr_str.HandleFunc("/all-user-list", aricto.GetAllUserList).Methods("GET")

	//issue subroute api
	isu_str := str.PathPrefix("/issues").Subrouter()
	isu_str.Handle("/all-issues-list", middleware.JwtMiddleware(aricto.GetAllIssuesList)).Methods("GET")

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
