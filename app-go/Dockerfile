FROM golang:1.23.3-alpine AS build-stage

# copy project
WORKDIR /home/app/app-go
COPY ./ ./

# build app-go's binary 
WORKDIR /home/app/app-go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app-go

# rebuild image from scratch 
FROM scratch AS build-release-stage

WORKDIR /

COPY --from=build-stage /home/app/app-go/bin/app-go /app-go

EXPOSE 80
