# REST-API in Go

<h2>Configuration</h2>

```
sudo docker run --name=postgres -e POSTGRES_PASSWORD=12345 -p 5430:5432 -d postgres:16.2
```

<h3>Export env variables</h3>

```
export DB="postgres:12345@localhost:5430/postgres?sslmode=disable" 
export KEY="imagine your own secret key"
```

<h2>How to run</h2>

``` make run ``` - in default way

``` make up ``` - in docker-compose way

``` make migrate ``` - make the migrations for database

<h3>Users/Auth EndPoints</h3>

```
http://127.0.0.1:8000/swagger/doc.json - swagger endpoint

/sign-up - auth handler to pass registration 
--> in request params: email="...." username="..." pass="..."

/sign-in - auth handler to get in service 
--> in request params: email="...." pass="..."

/users/{id:[0-9]+} - Get the account. GET - method

/users/{id:[0-9]+} - Update the account. PUT, POST - methods

/users/name/{id:[0-9]+} - Part update account name. POST, PATCH - methods
	
/users/email/{id:[0-9]+} - Part update account email. POST, PATCH - methods
	
/users/pass/{id:[0-9]+} - Part update account password. POST, PATCH - methods
	
/users/{id:[0-9]+} - Delete the account. DELETE - method
	
/users/logout/{id:[0-9]+} - Log out from account. GET-method

```
