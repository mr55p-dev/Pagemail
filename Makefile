server := pagemail
ifndef PM_ENV
	PM_ENV := dev
endif

css:
	npx tailwindcss -i ./input.css -o ./internal/assets/public/css/main.css

sqlc:
	sqlc generate

templ:
	templ generate

$(server): templ css
	go build -v ./cmd/$(server)

clean:
	rm -f $(server)

.PHONY = build clean run templ css

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build: $(server)

run: paegemail
	./$(server)
