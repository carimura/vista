FROM funcy/ruby:dev as build-stage
WORKDIR /function
ADD Gemfile* /function/
RUN bundle install
FROM funcy/ruby
WORKDIR /function
COPY --from=build-stage /usr/lib/ruby/gems/ /usr/lib/ruby/gems/
ADD . /function/
ENTRYPOINT ["ruby", "func.rb"]
