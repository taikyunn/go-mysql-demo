# 開発環境
FROM golang:1.15.2-alpine as dev

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 3000
# ホットリロード機能を入れるためにgo get(こいつはmain.goでimportできないので注意)
RUN go get -u github.com/cosmtrek/air
# airでgo runさせるためにCMDを下記に変更させる。
CMD ["air", "-c", ".air.toml"]

# 本番環境
FROM golang:1.15.7-alpine as builder

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download

COPY . ${ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary


FROM scratch as prod

ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder ${ROOT}/binary ${ROOT}

EXPOSE 3000
CMD ["/go/src/app/binary"]
