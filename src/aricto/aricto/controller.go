package aricto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"strconv"

	sitDatatype "aricto/datatypes"
	DB "aricto/dbconnection"

	"github.com/gorilla/context"
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
    fmt.Println( context.Get(r, "user_id"))
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
    				"message" : "All issue list successfully retrieved",
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

/**
* Get Issue Info by issue id
*/
var GetIssueInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
	issue_id := r.FormValue("issue_id")

	db,_ = DB.DbConnection(w)
    defer db.Close()
    // query
    res1 := &sitDatatype.IssuesTable{}

	err := db.QueryRow("SELECT * FROM issues WHERE id=?",issue_id).Scan(&res1.Id,&res1.Title,&res1.Description,&res1.AssignedTo,&res1.CreatedBy,&res1.Status)
	if err != nil {
        DB.CheckError(errors.New("Wrong Issue ID: The given issue id is not a valid one."), w)
    } else {
    	result, err := json.Marshal(map[string]interface{}{
			"forcePageRefresh" : 0,
			"message" : "Issue information successfully retrieved",
	    	"status" : true,
	    	"data" : res1,
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
})

/**
* Create a issue
*/
var CreateIssue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	issue_info := sitDatatype.CreateIssueType{}
	_ = json.NewDecoder(r.Body).Decode(&issue_info)
	// create issue	
	created_by := context.Get(r, "user_id")
	db,_ = DB.DbConnection(w)
    defer db.Close()

	_, err := db.Exec("INSERT INTO issues(title, description, assigned_to, created_by, status) VALUES(?, ?, ?, ?, ?)", issue_info.Title, issue_info.Description, issue_info.AssignedTo, created_by, issue_info.Status)
    if err != nil {
    	fmt.Println(err)
    	DB.CheckError(errors.New("Server error, unable to create your issue."), w)
    } else {
    	result, err := json.Marshal(map[string]interface{}{
			"forcePageRefresh" : 0,
			"message" : "New Issue successfully created",
	    	"status" : true,
	    	"data" : "New Issue successfully created",
	    	"error" : &sitDatatype.ErrorType{
			    		Exists : false,
			    		Errors : "",
				    },
			"statusCode" : 201,
			"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
    	})

	    if err != nil {
	    	DB.CheckError(err, w)
	    }
	    w.Header().Set("Content-Type", "application/json")
	    w.Write([]byte(result))
    } 
})

/**
* Update a issue
*/
var UpdateIssue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	issue_info := sitDatatype.CreateIssueType{}
	_ = json.NewDecoder(r.Body).Decode(&issue_info)
	// create issue	
	created_by := context.Get(r, "user_id")
	db,_ = DB.DbConnection(w)
    defer db.Close()

    // get issue info by id 
    var issue_created_by_id int
    var issue_id_response int
	err := db.QueryRow("SELECT id,created_by FROM issues WHERE id=? AND created_by=?",issue_info.Id, created_by).Scan(&issue_id_response, &issue_created_by_id)

	if issue_id_response == issue_info.Id {
		_, err = db.Exec("UPDATE issues SET title=?, description=?, assigned_to=?, status=? WHERE id=? AND created_by=?", issue_info.Title, issue_info.Description, issue_info.AssignedTo, issue_info.Status, issue_id_response, created_by)
	    
	    if err != nil {
	    	fmt.Println(err)
	    	DB.CheckError(errors.New("You don't have access to update/edit this issue information or the issue id is wrong."), w)
	    } else {
	    	result, err := json.Marshal(map[string]interface{}{
				"forcePageRefresh" : 0,
				"message" : "Issue information successfully updated",
		    	"status" : true,
		    	"data" : "Issue information successfully updated",
		    	"error" : &sitDatatype.ErrorType{
				    		Exists : false,
				    		Errors : "",
					    },
				"statusCode" : 201,
				"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
	    	})

		    if err != nil {
		    	DB.CheckError(err, w)
		    }
		    w.Header().Set("Content-Type", "application/json")
		    w.Write([]byte(result))
	    } 
	} else {
    	DB.CheckError(errors.New("You don't have access to update/edit this issue information or the issue id is wrong."), w)
	}
})

/**
* Delete issue
*/
var DeleteIssue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	issue_id, _ := strconv.Atoi(r.FormValue("issue_id"))
	created_by := context.Get(r, "user_id")

	db,_ = DB.DbConnection(w)
    defer db.Close()

    // get issue info by id 
    var issue_id_response int
	err := db.QueryRow("SELECT id FROM issues WHERE id=? AND created_by=?", issue_id, created_by).Scan(&issue_id_response)

	if issue_id_response == issue_id {
		_, err = db.Exec("DELETE FROM issues WHERE id=? AND created_by=?", issue_id, created_by)
	    
	    if err != nil {
	    	fmt.Println(err)
	    	DB.CheckError(errors.New("You don't have access to delete this issue or the issue id is wrong."), w)
	    } else {
	    	result, err := json.Marshal(map[string]interface{}{
				"forcePageRefresh" : 0,
				"message" : "Issue successfully removed",
		    	"status" : true,
		    	"data" : "Issue successfully removed",
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
	} else {
    	DB.CheckError(errors.New("You don't have access to delete this issue or the issue id is wrong."), w)
	}
})

/**
* get all the issues created by me
*/
var GetAllIssuesByMe = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	created_by := context.Get(r, "user_id")

	db,_ = DB.DbConnection(w)
    defer db.Close()
    // query
	rows, err := db.Query("SELECT * FROM issues WHERE created_by=?",created_by)
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
		"message" : "All issue list created by you successfully retrieved",
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

/**
* get all the list of issues assigned to me
*/
var GetAllIssuesAssignedToMe = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		created_by := context.Get(r, "user_id")

	db,_ = DB.DbConnection(w)
    defer db.Close()
    // query
	rows, err := db.Query("SELECT * FROM issues WHERE assigned_to=?",created_by)
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
		"message" : "All issue assigned to you successfully retrieved",
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
