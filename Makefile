run:
	docker run -d --name forum -p4000:4000 forum && echo "server started at http://localhost:4000/" --rm
build:
	docker build -t forum .