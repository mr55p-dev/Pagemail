FROM node:18 AS init
RUN mkdir /app
WORKDIR /app
COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json
COPY Makefile /app/Makefile
RUN make init

COPY . /app
RUN make build

FROM node:18
COPY --from=init /app/dist /dist

