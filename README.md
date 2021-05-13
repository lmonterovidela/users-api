# users-api

Users-api allows the creation, modification, deletion and search by user name, it also allows to obtain more details about the address of a user.

Run api:
- Download the code
- Run in a terminal: 'docker-compose -f docker-compose.db.yml up' (that up a mysqldatabase)
- Run make up or go build cmd/*.go


Run whole api on docker:
- Download the code
- Run in a terminal: docker-compose up

Run linter:
- golangci-lint run ./...

Run test:
- docker-compose -f docker-compose.test.yml up
- @go test ./...
- docker-compose -f docker-compose.test.yml down

----
**Endpoints:**
----

**Show User**
  Returns json data about a single user.

* **URL**

  :version/users/:id

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `id=[integer]`

**List User**
  Returns json data about a users.

* **URL**

  :version/users?:name&:size&:offset

* **Method:**

  `GET`
  
*  **Query Params**

   **Required:**
 
   `name=[string]`
   
   **Option:**
 
   `size=[integer]`
   `offset=[integer]`
   
 **Create User**
  Create new user.

* **URL**

  :version/users

* **Method:**

  `POST`
  
* **Data Params**
 
   `name=[string]`
   `address=[string]`
   `dob=[datetime]`
   
   
 **Update User**
  Update an existent user.

* **URL**

  :version/users/:id

* **Method:**

  `PUT`
  
*  **URL Params**

 **Required:**

 `id=[integer]`

* **Data Params**
 
   `name=[string]`
   `address=[string]`
   `dob=[datetime]`


**Delete User**
  Delete an existent user.

* **URL**

  :version/users/:id

* **Method:**

  `DELETE`
  
*  **URL Params**

   **Required:**
 
   `id=[integer]`


**Show Location information**
  Returns json data about a location of user.

* **URL**

  :version/users/:id/locations

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `id=[integer]`
