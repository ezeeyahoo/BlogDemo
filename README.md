## BlogDemo
A grpc implementation of blog CRUD APIs

## Steps
## Start the grpc server
### Run from project root folder
> go mod tidy
> go run server/main.go


### Run client tests
#### Run from project root folder
> go test -v client/main_test.go

#### Server Logs
```
2024/04/13 01:23:58 starting the server at: [::]:7878
2024/04/13 01:24:06 recvd create request for blog - Title:"blog1"  Content:"this is blog1"  Author:"alice"  PubDate:"01/01/2009"  Tags:"blog"  Tags:"grpc"
2024/04/13 01:24:06 new blog generated with ID: 69208532-0784-4476-b60a-3f25865fe80d
2024/04/13 01:24:06 recvd create request for blog - Title:"blog1"  Content:"this is blog2"  Author:"joe"  PubDate:"01/12/2009"  Tags:"blog"  Tags:"grpc"
2024/04/13 01:24:06 new blog generated with ID: 2571a3c4-a5d3-472c-9ded-ee62e17db1ea
2024/04/13 01:24:06 recvd GET request for blog ID: 2571a3c4-a5d3-472c-9ded-ee62e17db1ea
2024/04/13 01:24:06 Recvd update for ID: 2571a3c4-a5d3-472c-9ded-ee62e17db1ea
2024/04/13 01:24:06 deleted 2571a3c4-a5d3-472c-9ded-ee62e17db1ea
2024/04/13 01:24:06 update blog with new ID: d7f3e132-07c7-40a4-904d-44480f30835a
2024/04/13 01:24:06 update completed, new ID is: d7f3e132-07c7-40a4-904d-44480f30835a
2024/04/13 01:24:06 Recvd delete for blog ID - d7f3e132-07c7-40a4-904d-44480f30835a
```

#### Client Logs
```
eshant@Eshants-MacBook-Pro demoBlogServiceInGrpc % go test -v client/main_test.go
=== RUN   Test_main
    main_test.go:77: sending request for new blog- {Title:blog1 Content:this is blog1 Author:alice PubDate:01/01/2009 Tags:[blog grpc]}
    main_test.go:77: sending request for new blog- {Title:blog1 Content:this is blog2 Author:joe PubDate:01/12/2009 Tags:[blog grpc]}
    main_test.go:104: read blog output of ID 42b2d0bf-5a86-4f50-925a-698ad5670357 : Title:"blog1" Content:"this is blog2" Author:"joe" PubDate:"01/12/2009" Tags:"blog" Tags:"grpc"
    main_test.go:106: Title:"blog1" Content:"this is blog2" Author:"joe" PubDate:"01/12/2009" Tags:"blog" Tags:"grpc" 
         Title:"blog1" Content:"this is blog2" Author:"joe" PubDate:"01/12/2009" Tags:"blog" Tags:"grpc"
--- PASS: Test_main (0.01s)
PASS
ok      command-line-arguments  0.154s
```
