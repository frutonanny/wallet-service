FROM golang:1.18.7 as build

ENV BIN_FILE /opt/wallet-service/wallet-app
ENV CODE_DIR /go/src/github.com/frutonanny/wallet-service

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/service/*

LABEL SERVICE="wallet-service"

ENV CONFIG_FILE /etc/wallet-service/config.json
COPY /config/config.dev.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
