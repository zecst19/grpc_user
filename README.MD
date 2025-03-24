# UserService

## How to Run

With Kafka server running in port <code>:9092</code>, run <code>go run server/main.go</code> and call grpc endpoints in port <code>:50051</code>

## Endpoints

### <em>CreateUser</em>
Creates a new User with the given fields <br>

<b>Example Request:</b>

    {
        "FirstName" : "Critiano",
        "LastName"  : "Ronaldo",
        "Nickname"  : "BuenosDiasMatosinhos",
        "Password"  : <password>,
        "Email"     : "ronaldo@cr7.com"
        "Country"   : "PT", 
    }

### <em>GetUser</em>
Fetches an existing User <br>

<b>Example Request:</b>

    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }

### <em>UpdateUser</em>
Updates an existing User with the given fields, omitted fields stay the same <br>

<b>Example Request:</b>
    
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

<b>Example Request:</b>
    
    {
        "Id" : "26ef0140-c436-4838-a271-32652c72f6f2",
    }
    
### <em>ListUsers</em>
List all Users that match the given filters, if no filters are present returns all Users <br>

<b>Example Request:</b>

    {
        "Page"       1,
        "PageSize"  : 10,
        "Country"   : "PT",         //optional
        "LastName"  : "Ronaldo",    //optional
    }
    
#### To-Do
* remove password from <em>GetUser</em> and <em>ListUsers</em> responses
* add creation date filters in <em>ListUsers</em>
* add checks in <em>CreateUser</em> and <em>UpdateUser</em>
* add API Gateway + Containerization

#### Notes:
* Used MongoDB because it's the DB EFG uses plus it was good to use it again after a while to refresh my memory 
* Used sarama for it's ease of use and because it's written in Go (plus the support it has compared to other kafka go libraries)
* The next steps would be to containerize with Docker to have cleaner deployment process, with a mongo service and a kafka service running in parallel
* After that an HTTP client was on the TO-DO list
* After all this I'd focus on the details like strengthening the API with checks on Create and Update, maybe an endpoint to update the password (not including it in the Update was on purpose), further test cases, etc
