# Start from golang base image
FROM golang:1.22.6-alpine3.20 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

COPY . .

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.10

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

#Command to run the executable
CMD swag init -g cmd/main.go\
  && chmod -R 777 /app/docs \
  && /wait \
  && CompileDaemon --build="go build cmd/main.go"  --command="./main" --color
