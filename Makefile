SHELL := /bin/bash

.PHONY: up down backend frontend test

up:
	docker compose up --build -d

down:
	docker compose down

backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...
