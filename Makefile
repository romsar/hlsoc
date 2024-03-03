.PHONY: proto
proto:
	protoc -I grpc/proto grpc/proto/*.proto --go_out=./grpc/gen --go_opt=paths=source_relative --go-grpc_out=./grpc/gen --go-grpc_opt=paths=source_relative

.PHONY: serve
serve:
	go run ./cmd/serve

.PHONY: docker-run
docker-run:
	docker-compose up -d --build