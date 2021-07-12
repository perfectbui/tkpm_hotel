FROM golang:1.16 as build-stage 

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 && go build -o app-exe ./cmd/...

FROM golang:1.16 as production-stage

COPY --from=build-stage /app/app-exe /bin
