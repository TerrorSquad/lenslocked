FROM golang
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server
#TODO: Use env variable for port
EXPOSE 3000
CMD ./server
