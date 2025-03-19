# UserService

## How to Run

Run <code>go run server/main.go</code> and call grpc endpoints

## Endpoints

### <em>CreateUser</em>
Creates a new User with the given fields <br>

Example Request:

    {
        "FirstName" : "Critiano",
        "LastName"  : "Ronaldo",
        "Nickname"  : "BuenosDiasMatosinhos",
        "Password"  : <password>,
        "Email"     : "ronaldo@cr7.com"
        "Country"   : "PT", 
    }

### <em>GetUser</em>
Fetches an existing Users <br>

Example Request:

    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }

### <em>UpdateUser</em>
Updates an existing User with the given fields, omitted fields stay the same <br>

Example Request:
    
    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
        "FirstName" : "Critiano",                         //optional
        "LastName"  : "Ronaldo",                          //optional
        "Nickname"  : "BuenosDiasMatosinhos",             //optional
        "Password"  : <password>,                       //optional
        "Email"     : "ronaldo@cr7.com"                 //optional
        "Country"   : "PT",                             //optional
    }

### <em>DeleteUser</em>
Deletes an existing User <br>

Example Request:
    
    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }
    
### <em>ListUsers</em>
List all Users that match the given filters, if no filters are present returns all Users <br>

Example Request:

    {
        "Page"       1,
        "PageSize"  : 10,
        "Country"   : "PT",         //optional
        "LastName"  : "Ronaldo",    //optional
        "FromDate"  : "2025-02-01", //optional
        "ToDate"    : "2025-03-01", //optional
    }
    
#### Known Issues
* <em>ListUsers</em> doesn't return IDs
* <em>UpdateUser</em> doesn't update Fields


#### To-Do
* fix known issues
* add filters to <em>ListUsers</em>
* remove password from <em>GetUser</em> and <em>ListUsers</em> responses
* finish unit tests
* add healthcheck
* add API Gateway + Docker
* prep Kafka messages
* add logs
* add comments where needed