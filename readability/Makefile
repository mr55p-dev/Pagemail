init:
	mkdir dist
	npm install
	python3 -m venv venv
	venv/bin/python -m pip install -r requirements.txt

build:
	npx tsc 
	cp \
		main.py \
		requirements.txt \
		package.json \
		package-lock.json \
		dist/

clean-venv:
	rm -rf venv/

clean:
	rm -rf dist/

run:
	@npx tsc 
	curl -s $(url) | node dist/main.js --url $(url) --check
