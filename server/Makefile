SERVER_TGT = dist/server
.PHONY = init build serve clean

init:
	@echo "done"

build:
	go build -o $(SERVER_TGT) .

run:
	go run main.go --dir=pb_data_test

serve:
	dist/server --dir pb_data_test/ serve

clean:
	rm $(SERVER_TGT)

