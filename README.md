# interview

B2B team Software Engineer Interview Challenge

## Introduction
In a parallel universe early nineteen nineties, software and the internet are booming. Many
programming languages and technologies that we know of in our world today exist and are
production ready. You are a Software Engineer at MobileWallet2023 - a mobile wallet company
based in Singapore - that provides a mobile wallet as a service.

The mobile wallet is in its early stages and so far only the login feature has been implemented.

You have been tasked with implementing a new feature: Funds Transfers!

## Existing Solution Walkthrough
The existing solution is a 3-tier application consisting of:

- A VueJs/Vuetify frontend. The implementation is under the [web](https://github.com/dd-cs/interview/tree/main/web) folder.
- A Golang API. The implementation is under the [api](https://github.com/dd-cs/interview/tree/main/api) folder.
- A Postgres database. The relevant folder is [db](https://github.com/dd-cs/interview/tree/main/db)

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

## Requirements
- Provide a new page on the frontend, where a logged in user can select a beneficiary to send funds to.
- Provide a new endpoint in the API that allows for funds to be transferred from one user to another.
- Make any necessary changes to the DB schema in order to support the requirements above.
- Give some thoughts as to how the application would be monitored, instrumentation etc,
considering that it would be expected to provide logging downstream to an external
platform such as Prometheus (To be discussed during the interview)
- Please state any assumptions that you make.

## Deliverables
- Fully functioning code
- Unit tests (partial coverage is perfectly fine, no need to unit test everything)
- You may raise a Merge Request with your changes so that a code review session can ensue with the rest of the team
