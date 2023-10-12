build:
	docker-compose -f deployment/docker-compose.yml build

run: build
	docker-compose -f deployment/docker-compose.yml up -d

down:
	docker-compose -f deployment/docker-compose.yml down \
        --rmi local \
        --volumes \
        --remove-orphans \
        --timeout 60; \


