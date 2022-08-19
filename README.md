# Quickbooks
## Welcome to Quickbooks

### About Quickbooks
Quickbooks is the backend which allows you to create and manage your own library of books & users.

### Setup
Please keep in mind that it's still work in progress. 

Database: you can either install Postgresql locally or use the docker version. A docker compose file is available. <br>
For now it only sets up the database but later on it will contain the backend & frontend as well. 

### Docker Version
* Please create the following folders in the root level of the project by running `mkdir -p db-data/postgres` <br>
* Change the default if needed.
* Then simply run `docker-compose up -d` to start the Postgresql container.

### Database Dump
After creating the database either manually or via docker, you can either restore an empty database (schemas only) <br> 
or a database with some data, it includes an admin user, a couple of normal users as well as a few books to play around with in the API. <br>
The DB dumps can be found in the database folder

* run `psql go_books < go_books_db_dump.sql` to create the tables with some sample data.
* run `psql go_books < go_books_schema_only.sql` to just create the tables.

### Starting the server
There are several flags that can be passed to change things like the default port, environment, database connection info ect.<br>
It is best to configure these directly in the provided makefile, which currently uses the defaults.

* `make start` will start the server.
* `make restart` will restart the server.
* `make stop` will stop the server. 

Once the server is up you can use Postman, or curl to send requests. A frontend written in Vue is also being worked on & will also be committed soon. 

## Available endpoints (WIP, more endpoints will be added and or endpoints changed.)

## GET
`/healthcheck` returns status info <br>
`/v1/users` returns all registered users. (Requires admin privileges) <br>
`/v1/users/authenticated` returns all currently logged-in users (Requires admin privileges) <br>
`/v1/users/:id` returns a single user. (Requires authentication) <br>
`/v1/users/auth` authenticates a user, by checking their token (Requires authentication) <br>
`/v1/users/logout` logs out a user, by deleting token from DB (Required authentication) <br>

`/v1/books/` returns all books <br>
`/v1/books/:id` returns a book by ID <br>
`/v1/books/slug` returns a book by slug <br>
`/v1/books/authors` returns all authors <br>
`/v1/books/reviews` returns all reviews <br>
`/v1/books/reviews/:id` returns a review by ID <br>


## POST
`/v1/users/login` logs in a user <br>
`/v1/users` Creates a user <br>
`/v1/books` Creates a book (Requires authentication) <br>
`/v1/books/reviews` Creates a review (Requires authentication) <br>

## PATCH
`/v1/users/:id` Updates a user (Requires authentication) <br>
`/v1/books/:id` Updates a book (Requires authentication) <br>
`/v1/books/reviews/:id` Updates a review (Requires authentication) <br>

## DELETE
`/v1/users/:id` Deletes a user (Requires Admin privileges) <br>
`/v1/users/logout/:id`Force logout a user by destroying their token (Requires Admin privileges) <br>
`/v1/books/:id` Deletes a book (Requires authentication) <br>
`/v1/books/reviews/:id` Deletes a review (Requires authentication) <br>

## Endpoints WIP
### Show Users
Returns json data about all users, requires admin privileges
* URL: `/v1/users`
* Method: GET
* URL Params: None
* Body Params: None
* Headers: Bearer $token
* Success Response: 
  * Code: 200
  * Content: {"users":[{"id": 1, "name":"test", "email":"test@email.com"...}]}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Show Authenticated Users
Returns json data currently authenticated users, requires admin privileges.
* URL: `/v1/users/authenticated`
* Method: GET
* URL Params: None
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"users":[{"id":1. "name":"test", "email":"test@email.com"...}]}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Show User
Returns json data about a single user, requires authentication
* URL: `/v1/users/:id`
* Method: GET
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"user":{"id":1, "name":"test", "email":"test@email.com"...}}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Authenticate User
Returns json data about a single user, requires authentication. (for frontends to authenticate)
* URL: `/v1/users/auth`
* Method: GET
* URL Params: None
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"user":{"id":1, "name":"test", "email":"test@email.com"...}}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Logout User
Logs out a user and destroys their token in the DB.
* URL: `/v1/users/logout`
* Method: GET
* URL Params: None
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"message":"token destroyed"}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Show all books
Returns json data all books
* URL: `/v1/books`
* Method: GET
* URL Params:
  * Optional: 
    * title=[string] filter by title default ""
    * sort=[string] sort by (id, title, publication_year, -id, -title, -publication_year) default id
    * page=[int] limit default 1
    * page_size[int] offset default 20
* Body Params: None
* Success Response:
  * Code: 200
  * Content: {"books":[{"id":1, "title":"book", "author_id":1...}], "metadata": {"current_page":1, "page_size":20, "first_page": 1, "last_page":1, "total_records":1}}
* Error Response:
  * Code: 500
  * Content: {"error": "internal server error"}

### Show Book
Returns json data about a single book
* URL: `/v1/books/:id`
* Method: GET
* URL Params:
  * Required: id=[int]
* Body Params: None
* Success Response:
  * Code: 200
  * Content: {"book":{"id":1, "title":"book", "author_id":"1...}}
* Error Response:
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Show Book
Returns json data about a single book by slug
* URL: `/v1/books/:slug`
* Method: GET
* URL Params:
  * Required: slug=[string]
* Body Params: None
* Success Response:
  * Code: 200
  * Content: {"book":{"id":1, "title":"book", "author_id":"1...}}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
* Error Response:
  * Code: 500
  * Content: {"error": "internal server error"}

### Show all Authors
Returns json data about all authors
* URL: `/v1/books/authors`
* Method: GET
* URL Params:
  * Optional:
    * author=[string] filter by author default ""
    * sort=[string] sort by (id, author_name, publication_year, -id, -author_name, -publication_year) default id
    * page=[int] limit default 1
    * page_size[int] offset default 20
* Body Params: None
* Success Response:
  * Code: 200
  * Content: {"books":[{"id":1, "title":"book", "author_id":1...}], "metadata": {"current_page":1, "page_size":20, "first_page": 1, "last_page":1, "total_records":1}}
* Error Response:
  * Code: 500
  * Content: {"error": "internal server error"}

### Show Reviews
Returns json data about all reviews, requires authentication
* URL: `/v1/books/reviews`
* Method: GET
* URL Params:
  * Optional:
    * user=[string] filter reviews by user default ""
    * sort=[string] sort by (id, user, -id, -user) default id
    * page=[int] limit default 1
    * page_size=[int] offset default 20
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"reviews":[{"id": 1, "rating":5, "review":"test review"}], "metadata":{"current_page":1, "page_size":20, "first_page":1, "last_page":1, "total_records":4}}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Show Review
Returns json data about a review by ID
* URL: `/v1/books/reviews/:id`
* Method: GET
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: None
* Success Response:
  * Code: 200
  * Content: {"reviews":[{"id": 1, "rating":5, "review":"test review"}]}}
* Error Response:
  * Code: 400
  * Content: {"error": "no authorization header received"}
  * Code: 404
  * Content: {"error": "the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Login User
Logs in a new using their credentials.
* URL: `/v1/users`
* Method: POST
* URL Params: None
* Body Params:
  * Required:
    * `{"email":"test@test.com", "password":"password"}`
* Success Response:
  * Code: 200
  * Content: {"user":{"id":1, "name":"test", "email":"test@email.com"...}}
* Error Response:
  * Code: 422
  * Content: {"error": {"email":"should not be empty", "password":"should not be empty"}}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Create User
Creates a new user.
* URL: `/v1/users`
* Method: POST
* URL Params: None
* Body Params:
  * Required:
    * `{"name": "test", "email":"test@test.com", "password":"password"}`
* Success Response:
  * Code: 200
  * Content: {"user":{"id":1, "name":"test", "email":"test@email.com"...}}
* Error Response:
  * Code: 400
  * Content: {"error": "email address already taken"}
  * Code: 422
  * Content: {"error": {"name":"should not be empty","email":"should not be empty", "password":"should not be empty"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Create Book
Creates a new book, requires authentication.
* URL: `/v1/books`
* Method: POST
* URL Params: None
* Body Params:
  * Required:
    * `{"title": "book", "author_id":1, "publication_year":2015, "description":"Some book", "genres":["Science Fiction","Fantasy"]}`
* Success Response:
  * Code: 200
  * Content: {"book":{"id":1, "title":"book", "author_id":1...}}
* Error Response:
  * Code: 422
  * Content: {"error": {"title":"should not be empty","author_id":"should not be empty", "publication_year":"should not be empty", "description":"should not be empty", "genres":"should not be empty"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Update User
Updates a user.
* URL: `/v1/users/:id`
* Method: PATCH
* URL Params:
  * Required: id=[int]
* Body Params:
  * Optional:
    * `{"name": "test", "email":"test@test.com", "password":"password"}`
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"user":{"id":1, "name":"test", "email":"test@email.com"...}}
* Error Response:
  * Code: 400
  * Content: {"error": "email address already taken"}
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 422
  * Content: {"error": {"name":"should not be empty","email":"should not be empty", "password":"should not be empty"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Update Book
Updates a book.
* URL: `/v1/books/:id`
* Method: PATCH
* URL Params:
  * Required: id=[int]
* Body Params:
  * Optional:
    * `{"title": "test", "author_id":1, "publication_year":2015, "description":"Some book", "genres":["Science Fiction","Fantasy"]}`
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"book":{"id":1, "title":"test", "publication_year":2015...}}
* Error Response:
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 422
  * Content: {"error": {"title":"should not be empty","author_id":"should not be empty", "publication_year":"should not be empty", "description":"should not be empty", "genres":"should not be empty"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Create Review
Creates a new review, requires authentication.
* URL: `/v1/books/reviews`
* Method: POST
* URL Params: None
* Body Params:
  * Required:
    * `{"rating": 1, "review":"test review", "book_id":1, "user_id":1}`
* Success Response:
  * Code: 200
  * Content: {"review":{"id":1, "rating":1", "review":"test review"}}
* Error Response:
  * Code: 422
  * Content: {"error": {"rating":"should not be empty","review":"should not be empty", "rating":"should not be more than 5"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Update Review
Updates a review.
* URL: `/v1/books/reviews/:id`
* Method: PATCH
* URL Params:
  * Required: id=[int]
* Body Params:
  * Optional:
    * `{"rating":3, "review":updated review}`
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"review":{"id":1, "rating":3, "review":updated review}}
* Error Response:
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 422
  * Content: {"error": {"rating":"should not be empty","review":"should not be empty", "rating":"should not be more than 5"}}
  * Code: 500
  * Content: {"error": "internal server error"}

### Delete User
Deletes a user, requires admin privileges.
* URL: `/v1/users/:id`
* Method: DELETE
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"message":"user with id: 3 has been deleted"}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Delete User Token
Force logout a user by destroying their token, requires admin privileges
* URL: `/v1/users/logout:id`
* Method: DELETE
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"message":"token with id: 50 destroyed"}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Delete Book
Deletes a book, requires authentication.
* URL: `/v1/books/:id`
* Method: DELETE
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"message":"book with id: 3 has been deleted"}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}

### Delete Review
Deletes a review, requires authentication.
* URL: `/v1/books/reviews/:id`
* Method: DELETE
* URL Params:
  * Required: id=[int]
* Body Params: None
* Headers: Bearer $token
* Success Response:
  * Code: 200
  * Content: {"message":"review with id: 3 has been deleted"}
* Error Response:
  * Code: 401
  * Content: {"error": "you are not authorized to view this content"}
  * Code: 404
  * Content: {"error":"the requested resource could not be found"}
  * Code: 500
  * Content: {"error": "internal server error"}
