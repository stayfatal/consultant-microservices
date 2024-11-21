all: run

runTest:
	go test -v ./...

stop:
	docker-compose down

delChat: stop
	docker rmi consultant-microservices-chat

delTest: stop
	docker rmi consultant-microservices-test

delSso: stop
	docker rmi consultant-microservices-sso

delGatewayHttp: stop
	docker rmi consultant-microservices-gatewayhttp

delGatewayWebsocket: stop
	docker rmi consultant-microservices-gatewaywebsocket

chat: delChat
	docker-compose up

test: delTest
	docker-compose up
	
sso: delSso
	docker-compose up

gateh: delGatewayHttp
	docker-compose up

gatew: delGatewayWebsocket
	docker-compose up

run: delChat delSso delGatewayHttp delGatewayWebsocket
	docker-compose up
