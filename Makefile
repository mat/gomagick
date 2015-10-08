
run_server:
	PORT=8080 go run magickserver/magickserver.go

docker_build:
	docker build -t magickserver .

docker_run:
	docker run -it -p 8080:8080 --name magickserver magickserver
