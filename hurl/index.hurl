GET http://localhost:8080/
HTTP 200
[Asserts]
xpath "string(/html/head/title)" == "Pagemail"
xpath "string(/html/body/main/h1)" == "Hello, pagemail!"
