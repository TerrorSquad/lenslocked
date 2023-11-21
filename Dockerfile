FROM node:alpine AS tailwind_builder
WORKDIR /app
# Install dependencies
RUN npm init -y && \
       npm install tailwindcss postcss-cli autoprefixer && \
         npx tailwindcss init;
COPY ./templates /templates
COPY ./tailwind/tailwind.config.js ./tailwind.config.js
COPY ./tailwind/styles.css ./styles.css

RUN npx tailwindcss -c /app/tailwind.config.js -i /app/styles.css -o /styles.css --minify

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server

FROM alpine AS app
WORKDIR /app
COPY .env .env
COPY --from=builder /app/server ./server
COPY --from=tailwind_builder /styles.css /app/assets/styles.css
COPY --from=builder /app/.fly/scripts /app/.fly/scripts
COPY --from=builder /app/.fly/script_runner.sh /app/.fly/script_runner.sh
RUN /app/.fly/script_runner.sh
CMD ./server
