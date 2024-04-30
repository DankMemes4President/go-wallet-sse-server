FROM golang as base

###
# Dev image
###
FROM base as dev
RUN go install github.com/cosmtrek/air@latest
WORKDIR /opt/app/api
CMD ["air"]