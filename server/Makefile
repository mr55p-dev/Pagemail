SERVER_TGT = dist/server
.PHONY = init build serve clean

init:
	@echo "done"

build-local:
	go build -o $(SERVER_TGT) .

build:
	GOOS=linux GOARCH=amd64 go build -o $(SERVER_TGT) .

deploy:
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo service pagemail stop; sudo rm -f /home/ec2-user/server"
	scp -i ../pagemail-prod-rsa.pem -r dist/server ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com:/home/ec2-user/server
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo service pagemail start; sudo rm -f /home/ec2-user/server"

run:
	go run main.go --debug --dir=pb_data_test serve

serve:
	dist/server --dir pb_data_test/ serve

clean:
	rm $(SERVER_TGT)

