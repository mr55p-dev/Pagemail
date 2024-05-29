FROM public.ecr.aws/docker/library/node:lts-slim

RUN <<EOF
apt-get update
apt-get -y install curl
EOF

WORKDIR /app
COPY package.json package-lock.json tsconfig.json .

RUN npm install

COPY ./src ./src
RUN npx tsc

CMD ["node", "dist/main.js"]
