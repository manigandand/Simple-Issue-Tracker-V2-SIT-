# Steps to run this repo
1. DB and Tables<br />
	a. Have not created any migrations to create table and database, did manually in my local.<br />
	b. Database Name: `aricto`<br />
 	c. Tables: `user`, `issues`<br /><br />

2. DB Connection<br />
	a. Give the proper db connection address in dbconnection.go file.<br />
	b. `sql.Open("mysql", "homestead:secret@tcp([192.168.11.11]:3306)/aricto")`<br />
	c. Replace your db connection address here.<br /><br />

3. Run the code.<br />
	a. Nacicate to `mani@Mani:~/Simple-Issue-Tracker-V2-SIT-/src$` project directory.<br />
	b. Run `go run aricto/main.go`<br />
	b. Now the server started to serve in `http://localhost:3011/api/`<br /><br />

4. Available API to test.<br />
	```shell
	# a. The following API's no need of any authendication/access token<br />
		1. http://localhost:3011/api/<br />
		2. http://localhost:3011/api/user/all-user-list<br />
		3. http://localhost:3011/api/login<br />
		<br />
	# b. After login access token will be provided for the user, need to send access token along with the request<br />
		1. http://localhost:3011/api/issues/all-issues-list<br />
		2. http://localhost:3011/api/issues/issue-info?issue_id=1<br />
		3. http://localhost:3011/api/issues/create-issue<br />
		4. http://localhost:3011/api/issues/update-issue<br />
		5. http://localhost:3011/api/issues/delete-issue?issue_id=6<br />
		6. http://localhost:3011/api/issues/issues-by-me<br />
		7. http://localhost:3011/api/issues/issues-for-me
	```	
		<br /><br />

5. Need more info on the API, please refer the following Postman Document.<br />
	Note: Limited lifetime to this document.<br />
	1. https://documenter.getpostman.com/view/1310922/aircto-api-test/6n5yt3s<br />
	Postman Collection:<br />
	2. https://www.getpostman.com/collections/7c8f1844ca96f5e1b859<br /><br />

Please reach me out for more clarification @ manigandan.jeff@gmail.com, 9578628779, Skype: manigandan.dharmalingam 






# Simple-Issue-Tracker-V2-SIT-

Design:
System will have two models called User and Issue. With following information

#1. User <br />

a. Email
b. Username
c. FirstName
d. LastName
e. Password
f. AccessToken

#2. Issue
1. Title
2. Description
3. AssignedTo (User relation)
4. Createdby (User relation)
5. Status (Open, Closed)

Problem Statement:
Expose a RESTful API to make CRUD operation of Issue resource.
1. Every endpoint need user authentication
2. Authentication should be stateless (access_token)
3. User who created the issue only should be able to edit or delete that issue

Note:
1. Whenever an Issue is created or assigned to different user(in case of update), an email
should be triggered exactly after 12 mins to the particular user saying issue has been
assigned to him/her.
2. Every 24 hours an email should be triggered to every users with details of all the issues
assigned to him/her. Here 24 hours should be configurable.(for e.g we may ask you to
send emails for every 10 hours or even every 10 secs)
