BookMate application is a prototype of an online bookstore

The application consists of two databases:

Postgres database is used for storage
    command to build docker image Postgres:
        docker run --rm -d --net=book_m --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres
        -e POSTGRES_DB=BookMarket postgres
For Redis cache
     command to build docker image Redis:
        docker run --rm -d --name redis redis 

The application also consists of two microservices:

    Microservice Api-GW
Represents the entrance to the application.  A Dockerfile is described inside the service
There in the description there are instructions for assembling the image.
The microservice acts as a gateway and receives commands from the client via the http protocol using the rest api. 
Also in the api-gw microservice, authorization by token is implemented.

    MicroserviceThe Book  
Is a storage of data about books. It receives data from the api-gw service via the GRPS protocol and works directly with postgres and redis to perform CRUD operations on the bookstore.
The Dockerfile and description of how to create the image will also be inside the Book project

The pens are attached in the file: BooksMarket.postman_collection, which can be exported to Postman.

    CRUD methods
"/books"      http.MethodPost   --send information about the book
"/books"      http.MethodGet    --get information about a book or author
"/books/{id}" http.MethodPatch  --edit information about a book or author
"/books/{id}" http.MethodDelete --delete information about a book or author

    registration methods
"/register"   http.MethodPost   --user registration
"/login"      http.MethodPost   --getting a token after authentication
    
"/test"       http.MethodGet    --method for checking application operation

There is a docker-compose file in the deploy folder
Locally deployed application running on port 80