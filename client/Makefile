ifeq ($(ENV), "production")
	NODE_ENV := production
	VITE_PAGEMAIL_API_HOST := "https://pagemail.io"
else ifeq ($(ENV), "staging")
	NODE_ENV := production
	VITE_PAGEMAIL_API_HOST := "https://staging.pagemail.io"
else
	NODE_ENV := development
	VITE_PAGEMAIL_API_HOST := "http://localhost:8090"
endif

.PHONY = init run test

init:
	npm install

run:
	npm run dev

build:
	npm run build

deploy:
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo rm -rf /var/www/pagemail/* /home/ec2-user/dist"
	scp -i ../pagemail-prod-rsa.pem -r dist/* ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com:/home/ec2-user/dist
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo cp -r ./dist/* /var/www/pagemail/"
