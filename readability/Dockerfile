FROM oven/bun:latest

WORKDIR /app

RUN <<EOF
apt-get update
apt-get -y install curl
EOF

COPY package.json ./
COPY bun.lockb ./
COPY src ./

RUN bun install

CMD ["bun", "run", "main.ts"]
