FROM node:22-alpine AS node_base
RUN npm install -g pnpm

WORKDIR /app/frontend
COPY package.json pnpm-lock.yaml .
RUN pnpm i
COPY . .
RUN pnpm build

FROM golang:1.22-alpine AS build_base

RUN apk add git
RUN apk add build-base

WORKDIR /app/backend
COPY server/go.mod server/go.sum .
RUN go mod download

COPY server/* .
WORKDIR /app/backend
RUN CGO_ENABLED=1 GOOS=linux go build -o amvillage .

FROM alpine:3.20

WORKDIR /app
COPY --from=build_base /app/backend/amvillage /app/amvillage
COPY --from=node_base /app/frontend/build/ /app/ui/
COPY server/config.json .

EXPOSE 8080

CMD ["/app/amvillage", "-addr", ":8080", "-serve", "/app/ui"]
