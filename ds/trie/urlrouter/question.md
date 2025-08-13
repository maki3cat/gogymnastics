// url path
- /users/{:userID:}/location, GET/DELETE
- /users, GET
- /users/{:userID:}/account, GET/PATCH


// api
```
func register(path string, method string, handler func())
func route(path string, method string, handler func())
```
