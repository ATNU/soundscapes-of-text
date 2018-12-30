FROM node as builder

WORKDIR /app

COPY ./client/package*.json /app/

RUN npm install

COPY ./client /app/

ARG configuration=production

RUN npm run build -- --output-path=./dist/out --configuration $configuration


FROM nginx:latest

COPY --from=builder /app/dist/out /usr/share/nginx/html

COPY ./client/nginx-custom.conf /etc/nginx/conf.d/default.conf