# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /server

COPY ./.bin/app .
COPY ./config/ ./config/