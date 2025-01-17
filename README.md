# todo-goland-api
Todo app built along while learning with Mastering Go With GoLand course

# How to run 
1. Create a database in postgres  
2. Run [todo-schema.sql](scripts/todo-schema.sql) to create table
3. Create .env file and configure DB_URL 
4. Run go api server
```go
go run main.go
```
app will be running at `http://localhost:8080`

5. API documentation check [todo-api-spec.http](todo-api-spec.http)


# Todo/Pending 
- [x] .env configuration via godotenv  
- [ ] Update todo status
- [ ] Delete todo item
- [ ] Frontend with htmx 
- [ ] Docker setup