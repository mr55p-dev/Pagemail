url := https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio
init:
	npm install
	python3 -m venv venv
	venv/bin/python -m pip install -r requirements.txt

build:
	npx tsc 

run:
	@npx tsc 
	curl -s $(url) | node dist/main.js --url $(url) | /Users/ellis/Git/pagemail/readability/venv/bin/python test.py
