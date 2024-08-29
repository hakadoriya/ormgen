package example

//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --dialect postgres  --go-orm-output-path ./example/generated/postgres"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --dialect cockroach --go-orm-output-path ./example/generated/cockroach"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --dialect mysql     --go-orm-output-path ./example/generated/mysql"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --dialect sqlite3   --go-orm-output-path ./example/generated/sqlite3"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./example/model --dialect spanner   --go-orm-output-path ./example/generated/spanner"
