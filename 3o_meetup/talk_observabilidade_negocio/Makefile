.PHONY: all build up down logs clean

COMPOSE=docker compose

all: build up

build:
	$(COMPOSE) build

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

clean:
	$(COMPOSE) down -v --remove-orphans
	rm -rf bad_metrics/cmd/app good_metrics/cmd/app