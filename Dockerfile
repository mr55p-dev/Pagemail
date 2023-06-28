FROM golang:latest AS BACKEND_BUILD
COPY ./server /app
WORKDIR /app
RUN make build

FROM node:latest AS FRONTEND_BUILD
ARG PAGEMAIL_API_HOST=http://localhost
COPY client /app
WORKDIR /app
RUN make init
RUN PAGEMAIL_API_HOST=$PAGEMAIL_API_HOST make build

FROM nginx:stable-alpine3.17-slim
COPY --from=FRONTEND_BUILD /app/dist /app/pagemail/site
COPY --from=BACKEND_BUILD /app/dist/server /app/pagemail/server 
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY docker_entry.sh /start.sh

ENTRYPOINT "/start.sh"

# Environment="PAGEMAIL_TOKEN_SIGNING_PUBLIC_KEY=/home/ec2-user/pagemail-jwt.key.pub"
# Environment="PAGEMAIL_TOKEN_SIGNING_PRIVATE_KEY=/home/ec2-user/pagemail-jwt.key"

# [Service]
# ExecStart=/home/ec2-user/server serve --dir=/home/ec2-user/pb_data --http="0.0.0.0:4000"
# Restart=always
# 
# [Install]
# WantedBy=multi-user.target
