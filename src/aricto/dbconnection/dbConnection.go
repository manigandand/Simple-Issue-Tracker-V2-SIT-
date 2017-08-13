package dbconnection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sitDatatype "aricto/datatypes"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/**
* Global Variables
*/
var db *sql.DB 
var err error

/**
* db connection
*/
func DbConnection(w http.ResponseWriter) (*sql.DB, error){
	// Create an sql.DB and check for errors
    db, err = sql.Open("mysql", "homestead:secret@tcp([192.168.11.11]:3306)/aricto")
    if err != nil {
    	CheckError(err, w)
    }

    // Test the connection to the database
    err = db.Ping()
    if err != nil {
        CheckError(err, w)
    }

    return db,nil
}

/**
* error handler
*/
func CheckError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	result, _ := json.Marshal(map[string]interface{}{
		"forcePageRefresh" : 0,
		"message" : "Something went wrong.",
		"status" : false,
		"data" : "",
		"error" : &sitDatatype.ErrorType{
		    		Exists : true,
		    		Errors : err.Error(),
			    },
		"statusCode" : 500,
		"time" : struct{UnixTime int32}{int32(time.Now().Unix())},
	})
	
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result))
}