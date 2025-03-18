All:

UpTest:
	docker-compose -f docker-compose.test.yml --env-file config.env up --abort-on-container-exit --exit-code-from test

DownTest:
	docker-compose -f docker-compose.test.yml --env-file config.env down

Up:
	docker-compose --env-file config.env up

Down:
	docker-compose --env-file config.env down