FROM alpine:3.7

RUN apk --no-cache add curl

ADD sensormockery ./

ARG CONTAINER_PORT
ENV PORT $CONTAINER_PORT
EXPOSE $CONTAINER_PORT

ENTRYPOINT [ "/sensormockery" ]