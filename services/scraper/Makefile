deps:
	docker run --rm -v "$(PWD)":/func -w /func funcy/ruby:dev bundle install --standalone --clean

test: 
	docker run --rm -it -v "$(PWD)":/func -w /func funcy/ruby ruby scraper.rb


