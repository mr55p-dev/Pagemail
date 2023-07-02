.PHONY = init run

init:
	npm install

run:
	VITE_PAGEMAIL_API_HOST=http://localhost:8090 npm run dev

build:
	VITE_PAGEMAIL_API_HOST=https://pagemail.io npm run build

deploy:
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo rm -rf /var/www/pagemail/* /home/ec2-user/dist"
	scp -i ../pagemail-prod-rsa.pem -r dist/* ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com:/home/ec2-user/dist
	ssh -i ../pagemail-prod-rsa.pem ec2-user@ec2-18-169-123-238.eu-west-2.compute.amazonaws.com "sudo cp -r ./dist/* /var/www/pagemail/"
