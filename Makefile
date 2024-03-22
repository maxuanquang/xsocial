new_migration:
	migrate create -ext sql -dir scripts/migration/ -seq $(MESSAGE_NAME)
up_migration:
	migrate -path scripts/migration/ -database "mysql://root:123456@tcp(0.0.0.0:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose up
down_migration:
	migrate -path scripts/migration/ -database "mysql://root:123456@tcp(0.0.0.0:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose down
proto_aap:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/authen_and_post --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/authen_and_post --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	authen_and_post.proto
proto_newsfeed:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/newsfeed --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/newsfeed --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	newsfeed.proto
proto_newsfeed_publishing:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/newsfeed_publishing --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/newsfeed_publishing --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional \
	newsfeed_publishing.proto

.PHONY: vendor
vendor:
	go mod vendor -v
docker_clear:
	docker volume rm $(docker volume ls -qf dangling=true) & docker rmi $(docker images -qf "dangling=true")
compose_up_rebuild:
	docker compose up --build --force-recreate
compose_up:
	docker compose up
caddy:
	cp ./Caddyfile /etc/caddy/Caddyfile
	systemctl restart caddy
gen_swagger:
	swag init -g cmd/web_app/main.go
build_web:
	cp .env web/
	cd web && npm install && npm run build && cd ..
