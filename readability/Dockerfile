FROM python:3.8-slim

RUN <<EOF
apt-get update
apt-get -y install curl
EOF

WORKDIR /app
COPY requirements.txt .

RUN pip install -r requirements.txt

COPY . .

CMD ["flask", "run", "--host=0.0.0.0"]
