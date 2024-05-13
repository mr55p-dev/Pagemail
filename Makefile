server := ./tmp/main
ifndef PM_ENV
	PM_ENV := dev
endif

$(server):
	mkdir -p ./tmp
	templ generate
	GOARCH=amd64 GOOS=linux go build -o ./tmp/main .

.PHONY = build clean run

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build: $(server)
clean: $(server)
	rm -rf ./tmp/

run: $(server)
	./tmp/main
