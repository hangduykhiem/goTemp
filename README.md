# goTemp

goTemp is my attempt to create a RESTful service to query
the current status of the temperature sensors built with 
NodeMCU

### How to run on your local machine

- Clone the application
- Create a Postgres user name `test` with password `test`
- Run the application

### Current feature
- Query all values at `localhost:3333/value`
- Add new temperature value via POST request to 
`localhost:3333/value` with model

```
{
  "index" : 1, // sensor index
  "temp" : "23", // current temperature
  "hum" : "43" // current humidity
}
```

### Future feature
- Implement paging for all value

### Technology use
- PostgreSQL driver
- Chi
