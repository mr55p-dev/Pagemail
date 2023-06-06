SERVER_TGT = dist/server
.PHONY = init build serve clean

init:
	@echo "done"

build-local:
	go build -o $(SERVER_TGT) .

build:
	GOOS=linux GOARCH=amd64 go build -o $(SERVER_TGT) .

deploy:
	aws 

run:
	go run main.go --debug --dir=pb_data_test serve

serve:
	dist/server --dir pb_data_test/ serve

clean:
	rm $(SERVER_TGT)

