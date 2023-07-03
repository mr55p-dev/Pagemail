.PHONY = halt clean install-nginx install-stage start 

BASE_DIR := /home/ec2-user/pagemail
PROD_DIR := /home/ec2-user/prod
STAGE_DIR := /home/ec2-user/stage
SERVICES_TARGET := /etc/systemd
PROD_WEB_TARGET := /var/www/pagemail
STAGE_WEB_TARGET := /var/www/pagemail-staging

install-nginx:
	sudo rm /etc/nginx/conf.d/*
	sudo cp $(BASE_DIR)/nginx/* /etc/nginx/conf.d/
	sudo certbot certonly --nginx -d "v2.pagemail.io,www.pagemail.io,pagemail.io,staging.pagemail.io" --expand --non-interactive
	sudo systemctl restart nginx 

install-stage-frontend:
	rm -rf $(STAGE_WEB_TARGET)/*
	cp -r $(BASE_DIR)/client/dist/* $(STAGE_WEB_TARGET)/

install-stage-backend:
	sudo systemctl stop pagemail.staging
	sudo cp $(BASE_DIR)/services/pagemail.staging.service $(SERVICES_TARGET)/pagemail.staging.service
	cp $(BASE_DIR)/server/dist/server $(STAGE_DIR)/server
	sudo chmod a+x $(STAGE_DIR)/server
	systemctl --user daemon-reload
	sudo systemctl start pagemail.staging

install-prod-frontend:
	sudo systemctl stop pagemail
	sudo cp $(BASE_DIR)/services/pagemail.service $(SERVICES_TARGET)/pagemail.service
	rm -rf $(PROD_WEB_TARGET)/*
	cp -r $(BASE_DIR)/client/dist/* $(PROD_WEB_TARGET)/
	systemctl --user daemon-reload
	sudo systemctl start pagemail
	
install-prod-backend:
	cp $(BASE_DIR)/server/dist/server $(PROD_DIR)/server
	sudo chmod a+x $(PROD_DIR)/server

install-stage: install-stage-frontend install-stage-backend

install-prod: install-prod-frontend install-prod-backend

pre-install:
	if [ -d $(BASE_DIR) ]; then rm -rf $(BASE_DIR)/*; fi


post-install:
	rm -rf $(BASE_DIR)

clean:
	sudo rm -f /home/ec2-user/server
	sudo rm -rf /var/www/pagemail/* /home/ec2-user/dist

install:
	sudo cp -r /home/ec2-user/dist/* /var/www/pagemail/
	sudo chmod a+x /home/ec2-user/server

start:
	sudo service pagemail start
	sudo service nginx restart
