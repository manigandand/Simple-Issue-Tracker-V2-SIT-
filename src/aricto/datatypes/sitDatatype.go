package datatypes

/**
* Datatypes
*/
type UserTable struct {
	Id  		int    `json:"id"`
	Email 		string `json:"email"`	
	UserName 	string `json:"user_name"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Password 	string `json:"password"`
	AccessToken string `json:"access_token"`
}

type IssuesTable struct {
	Id  		int    `json:"id"`
	Title 		string `json:"title"`
	Description string `json:"description"`
	AssignedTo 	string `json:"assigned_to"`
	CreatedBy 	string `json:"created_by"`
	Status 		string `json:"status"`

}

type ErrorType struct {
	Exists 		bool  	 `json:"exists"`
	Errors      string   `json:"errors"`
}