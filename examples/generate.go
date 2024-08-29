package example

//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./examples/model --dialect postgres  --go-orm-output-path ./examples/generated/postgres"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./examples/model --dialect cockroach --go-orm-output-path ./examples/generated/cockroach"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./examples/model --dialect mysql     --go-orm-output-path ./examples/generated/mysql"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./examples/model --dialect sqlite3   --go-orm-output-path ./examples/generated/sqlite3"
//go:generate sh -cx "cd .. && pwd && go run ./cmd/ormgen generate ./examples/model --dialect spanner   --go-orm-output-path ./examples/generated/spanner"
