server := pagemail
ifndef PM_ENV
	PM_ENV := dev
endif

$(server):
	templ generate
	GOARCH=arm64 GOOS=linux go build .

.PHONY = build clean run

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build: $(server)
clean: $(server)
	rm -f $(server)

run: $(server)
	./tmp/main
