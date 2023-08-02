# url := https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio
url := https://tracking.tldrnewsletter.com/CL0/https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html?utm_source=tldrnewsletter/1/010001899c017aa2-fc5a345e-100d-4d13-850c-2ab5deaadc35-000000/8UEzeRvMLI-LqFs5MOwoODhx8d5d6pWNPNSn50bzzV8=311
init:
	npm install
	python3 -m venv venv
	venv/bin/python -m pip install -r requirements.txt

build:
	npx tsc 

run:
	@npx tsc 
	curl -s $(url) | node dist/main.js --url $(url) --check
