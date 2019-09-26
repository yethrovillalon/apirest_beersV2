# Clone
FROM alpine/git as clone
WORKDIR /api
RUN git clone https://github.com/yethrovillalon/apirest_beers

# Build
FROM golang:1.13.1-alpine as build
WORKDIR /api
COPY --from=clone /api/apirest_beers /api
RUN apk add git
# RUN go get -v -u github.com/gorilla/mux
# RUN go get -v -u gopkg.in/mgo.v2

RUN go get -d -v ./...
RUN go build

# Run
FROM golang:1.13.1-alpine
WORKDIR /api
EXPOSE 8081
COPY --from=build /api/api /api/apirest_beers
CMD ["/api/apirest_beers"]
