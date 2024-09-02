package example

//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --go-orm-output-path ./example/generated/postgres --dialect postgres"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --go-orm-output-path ./example/generated/cockroach --dialect cockroach"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --go-orm-output-path ./example/generated/mysql --dialect mysql"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --go-orm-output-path ./example/generated/sqlite3 --dialect sqlite3"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --go-orm-output-path ./example/generated/spanner --dialect spanner"
