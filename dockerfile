FROM golang:1.22-alpine3.20 AS builder

WORKDIR /consultant-microservices

COPY go.mod go.sum ./

RUN go mod download

COPY ./gen ./gen

COPY ./libs ./libs

#chat

FROM builder AS builder_chat

WORKDIR /consultant-microservices

COPY ./services/chat ./services/chat

RUN go build -o app ./services/chat/cmd/app/

FROM builder_chat AS chat

COPY --from=builder_chat . .

CMD [ "./app" ]

#gateway

FROM builder AS builder_gateway

WORKDIR /consultant-microservices

COPY ./services/gateway ./services/gateway

RUN go build -o app ./services/gateway/cmd/app/

FROM builder_gateway AS gateway

COPY --from=builder_gateway . .

CMD [ "./app" ]

#matchmaking

FROM builder AS builder_matchmaking

WORKDIR /consultant-microservices

COPY ./services/matchmaking ./services/matchmaking

RUN go build -o app ./services/matchmaking/cmd/app/

FROM builder_gateway AS matchmaking

COPY --from=builder_matchmaking . .

CMD [ "./app" ]

#auth

FROM builder AS builder_auth

WORKDIR /consultant-microservices

COPY ./libs/publicauth ./libs/publicauth

COPY ./services/auth ./services/auth

RUN go build -o app ./services/auth/cmd/app/

FROM builder_auth AS auth

COPY --from=builder_auth . .

CMD [ "./app" ]

#test

FROM builder AS test

WORKDIR /consultant-microservices

COPY . .

CMD [ "go","test","-v","./..." ]