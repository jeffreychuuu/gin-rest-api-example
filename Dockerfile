#FROM golang:alpine
#ENV POSTGRES_HOST=host.docker.internal
#WORKDIR /app
#ADD . /app
#RUN cd /app && go build -o app
#EXPOSE 8080
#ENTRYPOINT ./app

# build stage
FROM golang:alpine AS build-env
ADD . /app
RUN cd /app && go build -o app

# final stage
FROM alpine
ENV POSTGRES_HOST=host.docker.internal
WORKDIR /app
COPY --from=build-env /app /app
EXPOSE 8080
ENTRYPOINT ./app