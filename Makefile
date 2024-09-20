.PHONY = build clean run

server := pagemail
ifndef PM_ENV
	PM_ENV := dev
endif

$(server):
	templ generate
	go build ./cmd/pagemail/

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build: $(server)
clean: 
	rm -f $(server)

run: $(server)
	./pagemail
