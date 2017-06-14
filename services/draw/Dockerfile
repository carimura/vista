FROM ruby:2.4.0-alpine

RUN apk add --update alpine-sdk ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --update imagemagick-dev imagemagick && rm -rf /var/cache/apk/*
WORKDIR /func
ADD . /func/
RUN bundle install
ENTRYPOINT ["ruby", "func.rb"]

