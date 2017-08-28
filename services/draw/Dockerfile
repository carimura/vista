FROM ruby:2.4.0-alpine

RUN apk add --update alpine-sdk ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --update imagemagick-dev imagemagick && rm -rf /var/cache/apk/*

WORKDIR /func
ADD Gemfile* /func/
RUN bundle install

ADD . /func/
ENTRYPOINT ["ruby", "func.rb"]
