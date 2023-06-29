FROM golang:1.20.4 AS BACKEND_BUILD
RUN mkdir /app
WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ /app
RUN --mount=type=cache,target=/root/.cache/go-build make build


FROM node:18 AS FRONTEND_BUILD
RUN mkdir /app
WORKDIR /app
COPY client/package.json /app/package.json
COPY client/package-lock.json /app/package-lock.json
COPY client/Makefile /app/Makefile
RUN make init
COPY client/ /app
ENV PAGEMAIL_API_HOST=http://localhost:5001
RUN npm run build


FROM nginx:stable-alpine3.17-slim
COPY --from=FRONTEND_BUILD /app/dist /app/pagemail/site
COPY --from=BACKEND_BUILD /app/dist/server /app/pagemail/server 
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY docker_entry.sh /start.sh


EXPOSE 80
EXPOSE 443
EXPOSE 8080
EXPOSE 4000


ENTRYPOINT ["/start.sh"]
