.PHONY: setup test clean

setup:
	go mod tidy

test:
	go test ./...

clean:
	rm -f kortex_flight_recorder.jsonl
