.PHONY = clean-all build-all clean-css clean-templates clean-sql clean test

# get the correct os and arch for installing tailwind
ifeq ($(shell uname), Darwin)
	OS := macos
else
	OS := linux
endif
ifeq ($(shell arch), arm64)
	ARCH := arm64
else
	ARCH := x64
endif

server := pagemail
all: clean-all $(server) 

# Install
tailwindcss := ./bin/tailwindcss
templ := ${GOBIN}/templ
sqlc := ${GOBIN}/sqlc
air := ${GOBIN}/air
dbmate := ./bin/dbmate
tools := $(tailwindcss) $(templ) $(sqlc) $(air) $(dbmate)
$(tailwindcss): 
	mkdir -p ./bin
	curl -fsSL \
		-o ./bin/tailwindcss \
		https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$(OS)-$(ARCH)
	chmod +x ./bin/tailwindcss
$(dbmate):
	mkdir -p ./bin
	curl -fsSL -o ./bin/dbmate \
		https://github.com/amacneil/dbmate/releases/latest/download/dbmate-$(OS)-$(ARCH)
	chmod +x ./bin/dbmate
$(templ):
	go install github.com/a-h/templ/cmd/templ@latest
$(sqlc):
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
$(air):
	go install github.com/air-verse/air@latest
install: $(tools)
uninstall: 
	rm -f $(tools)

# Templates 
templates := $(shell ls render/**/*.templ | sed "s/\.templ/_templ.go/")
$(templates): $(templ)
	$(templ) generate

templates: $(templates)
watch-templates: clean-templates
	$(templ) generate -watch
clean-templates:
	rm -f \
		./render/**/*_templ.go \
		./render/**/*_templ.txt

# Sql
sql := db/queries/db.go db/queries/models.go $(shell ls db/query.*.sql | sed "s/query.*\.sql/queries\/&.go/")
$(sql): $(sqlc)
	$(sqlc) generate

sql: clean-sql $(sql)
clean-sql:
	rm -f ./db/queries/*.go

css := assets/css/main.css
css-input := tailwind.base.css
$(css): $(tailwindcss) $(css-input)
	$(tailwindcss) --input $(css-input) --output $(css) --minify

css: clean-css $(css)
watch-css: clean-css
	$(tailwindcss) --input $(css-input) --output $(css) --watch
clean-css:
	rm -f $(css)

# Server
$(server): $(templates) $(sql) $(css)
	go build -o $(server) ./cmd/pagemail 

server: clean $(server)
watch-server: clean $(air)
	air
clean: 
	rm -f $(server)

# Shared
prerequisites: $(templates) $(sql) $(css)
clean-all: clean clean-css clean-sql clean-templates
test: $(server)
	go test ./...
