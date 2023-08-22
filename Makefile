.PHONY = halt clean install-nginx install-stage start 

BASE_DIR := /home/ec2-user/pagemail
PROD_DIR := /home/ec2-user/prod
STAGE_DIR := /home/ec2-user/stage
SERVICES_TARGET := /etc/systemd/system
PROD_WEB_TARGET := /var/www/pagemail
STAGE_WEB_TARGET := /var/www/pagemail.staging
TEST_WEB_TARGET := /var/www/testsites

# Static files
install-nginx:
	sudo rm /etc/nginx/conf.d/*
	sudo cp $(BASE_DIR)/nginx/* /etc/nginx/conf.d/
	sudo certbot certonly --nginx -d "v2.pagemail.io,www.pagemail.io,pagemail.io,staging.pagemail.io" --expand --non-interactive
	sudo systemctl restart nginx 

install-test-sites:
	cp $(BASE_DIR)/test_pages/* $(TEST_WEB_TARGET)/

install-stage-templates:
	rm -rf $(STAGE_DIR)/templates/*
	cp $(BASE_DIR)/templates/* $(STAGE_DIR)/templates/

install-prod-templates:
	rm -rf $(PROD_DIR)/templates/*
	cp $(BASE_DIR)/templates/* $(PROD_DIR)/templates/

# Frontend installations
install-stage-frontend:
	rm -rf $(STAGE_WEB_TARGET)/*
	cp -r $(BASE_DIR)/client/dist/* $(STAGE_WEB_TARGET)/

install-prod-frontend:
	rm -rf $(PROD_WEB_TARGET)/*
	cp -r $(BASE_DIR)/client/dist/* $(PROD_WEB_TARGET)/
	
# Backend installations
install-stage-backend:
	sudo cp $(BASE_DIR)/services/pagemail.staging.service $(SERVICES_TARGET)/pagemail.staging.service
	sudo chmod a+x $(BASE_DIR)/server/dist/server
	cp $(BASE_DIR)/server/dist/server $(STAGE_DIR)/server
	sudo systemctl daemon-reload

install-prod-backend:
	sudo cp $(BASE_DIR)/services/pagemail.service $(SERVICES_TARGET)/pagemail.service
	cp $(BASE_DIR)/server/dist/server $(PROD_DIR)/server
	sudo chmod a+x $(PROD_DIR)/server
	sudo systemctl daemon-reload

# Readability installations
install-stage-readability:
	sudo cp $(BASE_DIR)/readability/dist/* $(STAGE_DIR)/readability/
	npm --prefix $(STAGE_DIR)/readability/ ci
	python3 -m venv $(STAGE_DIR)/readability/venv
	$(STAGE_DIR)/readability/venv/bin/pip install -r $(STAGE_DIR)/readability/requirements.txt

install-prod-readability:
	sudo cp \
		$(BASE_DIR)/readability/dist/main.py \
		$(BASE_DIR)/readability/dist/main.js \
		$(PROD_DIR)/readability/

# Full installations
install-stage: install-stage-frontend install-stage-backend install-stage-readability

install-prod: install-prod-frontend install-prod-backend install-prod-readability

# Pre install scripts
pre-install-stage:
	if [ -d $(BASE_DIR) ]; then rm -rf $(BASE_DIR)/*; fi
	sudo systemctl stop pagemail.staging

# Post install script
pre-install-prod:
	if [ -d $(BASE_DIR) ]; then rm -rf $(BASE_DIR)/*; fi
	sudo systemctl stop pagemail

# Post install scripts
post-install-stage:
	# rm -rf $(BASE_DIR)
	sudo systemctl start pagemail.staging

post-install-prod:
	# rm -rf $(BASE_DIR)
	sudo systemctl start pagemail

# Legacy
clean:
	sudo rm -f /home/ec2-user/server
	sudo rm -rf /var/www/pagemail/* /home/ec2-user/dist

install:
	sudo cp -r /home/ec2-user/dist/* /var/www/pagemail/
	sudo chmod a+x /home/ec2-user/server

start:
	sudo service pagemail start
	sudo service nginx restart
