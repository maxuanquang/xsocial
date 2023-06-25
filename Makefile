# make new_migration MESSAGE_NAME="message name"
new_migration:
	migrate create -ext sql -dir scripts/migration/ -seq $(MESSAGE_NAME)
up_migration:
	migrate -path scripts/migration/ -database "mysql://quangmx:2511@tcp(192.168.0.103:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose up
down_migration:
	migrate -path scripts/migration/ -database "mysql://quangmx:2511@tcp(192.168.0.103:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose down
proto_aap:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/authen_and_post --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/authen_and_post --go-grpc_opt=paths=source_relative \
	authen_and_post.proto
proto_newsfeed:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/newsfeed --go_opt=paths=source_relative \
	--go-grpc_out=pkg/types/proto/pb/newsfeed --go-grpc_opt=paths=source_relative \
	newsfeed.proto