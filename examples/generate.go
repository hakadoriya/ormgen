package example

//go:generate sh -cx "cd .. && pwd && go run -trimpath ./cmd/ormgen generate ./examples/model --debug --dialect postgres  --go-orm-output-path ./examples/generated/postgres"
//go:generate sh -cx "cd .. && pwd && go run -trimpath ./cmd/ormgen generate ./examples/model --debug --dialect cockroach --go-orm-output-path ./examples/generated/cockroach"
//go:generate sh -cx "cd .. && pwd && go run -trimpath ./cmd/ormgen generate ./examples/model --debug --dialect mysql     --go-orm-output-path ./examples/generated/mysql"
//go:generate sh -cx "cd .. && pwd && go run -trimpath ./cmd/ormgen generate ./examples/model --debug --dialect sqlite3   --go-orm-output-path ./examples/generated/sqlite3"
//go:generate sh -cx "cd .. && pwd && go run -trimpath ./cmd/ormgen generate ./examples/model --debug --dialect spanner   --go-orm-output-path ./examples/generated/spanner"
