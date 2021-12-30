build:
	@go build -o bin/main main.go
run:
	@go run main.go
docker:
	@docker image build -f Dockerfile -t goo .
	@docker container run -p 8080:8080 -d --name forum goo

prune:
	@docker system prune
