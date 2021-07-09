# address_book
Simple server API application for an address book in golang The entries are stored in system memory so as soon as the application closes everything is wiped. 

Unfortunately the testing is non functional because the built in go get client will not allow connections to localhost and colons are not allowed in the path

To build and run the server, simply download the server.go file, navigate to its location in the command line and type "go run server.go"

Every call to the server needs to be made to "localhost:8080/endpoint"

Endpoints:
GET /entriesfn 
This endpoint will return a json object containing every entry in the address book ordered alphabetically by first name

GET /entriessn 
This endpoint will return a json object containing every entry in the address book ordered alphabetically by surname

POST /entry
This endpoint will add a new entry to the address book
Params ; "{"forename": "FirstName", "surname": "LastName"}"
optionally a phone number can be added like this;
"{"forename": "FirstName", "surname": "LastName", "phonenumber": 07928401283}"

DELETE /entry
This endpoint will delete an entry in the address book matching any parameters sent in the request It must contain a first name and last name
Param: "{"forename": "Alex", "surname": "Dingleby"}"

PUT /search
This endpoint will search the address books first name and last name entries for any entries that have matching 
letters to the parameter sent via the endpoint
Param: "{"searchterm": "Al"}"


