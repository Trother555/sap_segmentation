.PHONY: local
local:
	docker-compose up -d


.PHONY: run-local
run:
	DB_USER=admin DB_PASSWORD=admin DB_NAME=postgres DB_HOST=localhost DB_PORT=5432 \
	LOG_PATH=log.txt CONN_URI=http://localhost:3001/ords/bsm/segmentation/get_segmentation \
	IMPORT_BATCH_SIZE=5 \
	go run ./cmd/sap_segmentation


.PHONY: test
test: 
	go test ./...

.PHONY: test-all
test-all: 
	go test --tags=postgres ./...


