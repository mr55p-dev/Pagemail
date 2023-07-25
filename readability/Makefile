init:
	npm install
	python3 -m venv venv
	# venv/bin/python -m pip install -r requirements.txt

build:
	npx tsc 

run:
	@npx tsc 
	@curl http://testsite.pagemail.io/long_everything.html | node dist/main.js http://testsite.pagemail.io/long_everything.html
