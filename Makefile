build:
	@go build -o bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 5000:5000 api	