# UserService

## How to Run

Run <code>go run server/main.go</code> and call grpc endpoints

## Endpoints

### CreateUser
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

### GetUser
Fetches an existing Users <br>

Example Request:

    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }

### UpdateUser
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

### DeleteUser
Deletes an existing User <br>

Example Request:
    
    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }
    
### ListUsers
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
    
