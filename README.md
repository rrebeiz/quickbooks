# Quickbooks
## Welcome to Quickbook

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
`/v1/users/logout` `logs out a user, by deleting token from DB (Required authentication) <br>

## POST
`/v1/users/login` logs in a user <br>
`/v1/users` Creates a user <br>
`/v1/books` Creates a book (Requires authentication) <br>

## PATCH
`/v1/users/:id` Updates a user (Requires authentication) <br>
`/v1/books/:id` Updates a book (Requires authentication) <br>

## DELETE
`/v1/users/:id` Deletes a user (Requires Admin privileges) <br>
`/v1/users/logout/:id`Force logout a user by destroying their token (Requires Admin privileges) <br>

## Endpoint examples
* Coming soon