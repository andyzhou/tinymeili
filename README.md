# tiny meiLi
search client api for meiLi search, second dev code.

# base rule
- Use low level access the api and call the web api access meiLi search service internal.

# example
pls see sub dir of `example`

#3rd depend
- meilisearch service v1.9.1
- meilisearch go client v0.27.2

#testing
go test -v -run="AddDoc"
go test -bench="AddDoc"
go test -bench="AddDoc" -benchmem -benchtime=10s
