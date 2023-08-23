init:
	mkdir dist
	npm install
	python3 -m venv dist/venv
	dist/venv/bin/python -m pip install -r requirements.txt

build:
	npx tsc 
	cp \
		main.py \
		requirements.txt \
		package.json \
		package-lock.json \
		dist/

clean:
	rm -rf dist/

run:
	@npx tsc 
	curl -s $(url) | node dist/main.js --url $(url) --check
