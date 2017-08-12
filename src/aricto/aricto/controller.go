package aricto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	sitDatatype "aricto/datatypes"
	DB "aricto/dbconnection"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/dgrijalva/jwt-go"
	
)

/**
* Global Variables
*/
var db *sql.DB 
var err error
var mySigningSecretKey = []byte("qwerty123")

/**
 * [HandleIndex description]
 * @param {[type]} w http.ResponseWriter [description]
 * @param {[type]} r *http.Request       [description]
 */
func HandleIndex(w http.ResponseWriter, r *http.Request) {	
	result, _ := json.Marshal(map[string]interface{}{
	    				"forcePageRefresh" : 0,
	    				"message" : "Welcome to AirCTO!",
				    	"status" : true,
				    	"data" : "Welcome to AirCTO Test - Simple Issue Tracker - SIT V.0",
				    	"error" : &sitDatatype.ErrorType{
						    		Exists : false,
						    		Errors : "",
							    },
						"statusCode" : 200,
						"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
			    	})
	w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result))
}

/**
* Login handler
*/
func PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	db,_ = DB.DbConnection(w)
    defer db.Close()

	// Grab from the user info 
    res := sitDatatype.UserTable{}
    err = db.QueryRow("SELECT * FROM user WHERE email=? AND password=?", email,password).Scan(&res.Id,&res.Email,&res.UserName,&res.FirstName,&res.LastName,&res.Password,&res.AccessToken)
    // If not then redirect to the login page
    if err != nil {
        DB.CheckError(errors.New("Unauthorized: Wrong Credentials. Unfortunately, your login credentials do not yet have access to the app."), w)
    } else {
    	if res.Id > 0 {
    		// valid user, generate access token 
    		acstkn := getTokenHandler(res)
    		data := struct{AccessToken string `json:"access_token"`}{acstkn}
    		result, _ := json.Marshal(map[string]interface{}{
				"forcePageRefresh" : 0,
				"message" : "You have successfully logged in.",
		    	"status" : true,
		    	"data" : data,
		    	"error" : &sitDatatype.ErrorType{
				    		Exists : false,
				    		Errors : "",
					    },
				"statusCode" : 200,
				"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
	    	})
		   
			w.Header().Set("Content-Type", "application/json")
		    w.Write([]byte(result))
		}
    }    	
}

/**
* generate JWT Access Token
*/
func getTokenHandler(res sitDatatype.UserTable) string {
	/* Create the token */
    token := jwt.New(jwt.SigningMethodHS256)
    // Create a map to store our claims
    claims := token.Claims.(jwt.MapClaims)
    // Set token claims 
    claims["id"] 		 = res.Id
    claims["email"] 	 = res.Email
    claims["user_name"]  = res.UserName
    claims["first_name"] = res.FirstName
    claims["last_name"]  = res.LastName
    claims["exp"]   	 = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    tokenString, _ := token.SignedString(mySigningSecretKey)

    /* Finally, write the token to the browser window */
    return tokenString
}

/**
 * [getAllUserList description]
 * @param {[type]} w http.ResponseWriter [description]
 * @param {[type]} r *http.Request       [description]
 */
func GetAllUserList(w http.ResponseWriter, r *http.Request) {	
	fmt.Println("Hi Aircto")
	db,_ = DB.DbConnection(w)   
    defer db.Close()

    // query
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
        DB.CheckError(err, w)
    }
    defer rows.Close()

    var userRes []*sitDatatype.UserTable
    
    for rows.Next() {
        res1 := &sitDatatype.UserTable{}
        err = rows.Scan(&res1.Id,&res1.Email,&res1.UserName,&res1.FirstName,&res1.LastName,&res1.Password,&res1.AccessToken)
        if err != nil {
	        DB.CheckError(err, w)
	    } 

	    userRes = append(userRes, res1) 
    }

    result, err := json.Marshal(map[string]interface{}{
    				"forcePageRefresh" : 0,
    				"message" : "All user list successfully retrieved",
			    	"status" : true,
			    	"data" : userRes,
			    	"error" : &sitDatatype.ErrorType{
					    		Exists : false,
					    		Errors : "",
						    },
					"statusCode" : 200,
					"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
			    	})

    if err != nil {
    	DB.CheckError(err, w)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result))
}

/**
 * [getAllIssuesList description]
 * @param {[type]} w http.ResponseWriter [description]
 * @param {[type]} r *http.Request       [description]
 */
var GetAllIssuesList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
	db,_ = DB.DbConnection(w)
    defer db.Close()

    // query
	rows, err := db.Query("SELECT * FROM issues")
	if err != nil {
        DB.CheckError(err, w)
    }
    defer rows.Close()

    var issuesRes []*sitDatatype.IssuesTable
    
    for rows.Next() {
        res1 := &sitDatatype.IssuesTable{}
        err = rows.Scan(&res1.Id,&res1.Title,&res1.Description,&res1.AssignedTo,&res1.CreatedBy,&res1.Status)
        if err != nil {
	        DB.CheckError(err, w)
	    }

	    issuesRes = append(issuesRes, res1) 
    }

    result, err := json.Marshal(map[string]interface{}{
    				"forcePageRefresh" : 0,
    				"message" : "All user list successfully retrieved",
			    	"status" : true,
			    	"data" : issuesRes,
			    	"error" : &sitDatatype.ErrorType{
					    		Exists : false,
					    		Errors : "",
						    },
					"statusCode" : 200,
					"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
			    	})

    if err != nil {
    	DB.CheckError(err, w)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result))
})