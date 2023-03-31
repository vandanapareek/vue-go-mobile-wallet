## Introduction
The mobile wallet with login and funds transfer features.

## Walkthrough
This is a 3-tier application consisting of:

- A VueJs/Vuetify frontend. The implementation is under the [web] folder.
- A Golang API. The implementation is under the [api] folder.
- A Postgres database. The relevant folder is [db] folder.

### Frontend
Once a user navigates to the frontend (by default found at http://localhost:8080), if she is not logged in, she will be redirected to a login page where she will have to provide her username and password. There are a few test users registered already:

| **Username** | **Password** |
|--------------|--------------|
| Alice        | password123  |
| Bob          | password123  |
| Charlie      | password123  |
| David        | password123  |

Once logged in, she will be redirected to the page that needs to be implemented.

### API
The API exposes a single endpoint **/login**, which is expecting a POST request with a JSON payload of the form: 

    {"username": "Alice", "password": "password123"}
    
If the credentials are valid, it responds with 200 OK, returning a JWT token in the response body which is valid for 60 minutes. The JWT token returned by the API is stored in the browser's local storage

If the credentials are invalid, it responds with 401 Unauthorized.

If any errors happen during the processing of the request, a 500 Internal Server Error response is returned.

### Steps to run the app

    docker-compose up
Spins up the application

