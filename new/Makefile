server := ./tmp/main

$(server):
	GOARCH=amd64 GOOS=linux go build -o ./tmp/main .

.PHONY = build clean run

build: $(server)
clean: $(server)
	rm -rf ./tmp/

run: $(server)
	./tmp/main

