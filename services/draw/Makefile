deps:
	docker run --rm -v "$(PWD)":/func -w /func funcy/ruby:dev bundle install --standalone --clean

test:
	docker run --rm -e "PAYLOAD_FILE=payload.json" -v "$(PWD)":/func -w /func treeder/ruby-imagemagick ruby /func/func.rb
