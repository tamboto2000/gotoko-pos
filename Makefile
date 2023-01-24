build-docker:
	docker build -t "gotoko-pos" .

run-docker:	
	docker stop gotoko-pos
	docker rm gotoko-pos
	docker run --name=gotoko-pos -e MYSQL_HOST=host.docker.internal -e MYSQL_PORT=3306 -e MYSQL_USER=root -e MYSQL_PASSWORD=kepler22b -e MYSQL_DBNAME=gotoko-pos -p 3030:3030 gotoko-pos
