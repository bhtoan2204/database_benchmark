generate_data:
	k6 run k6/fake_data.js
.PHONY: generate_data

run:
	go run cmd/main/main.go
.PHONY: run
