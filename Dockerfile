FROM golang:alpine AS builder
WORKDIR /bread
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOEXPERIMENT==fieldtrack,boringcrypto,greenteagc,randomizedheapbase64,simd CGO_ENABLED=0 go build -trimpath -o /noai .

FROM scratch
COPY --from=builder /noai /noai
EXPOSE 8080
ENTRYPOINT ["/noai"]