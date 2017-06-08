FROM treeder/ruby-imagemagick

WORKDIR /func
ADD . /func

# DEV
#CMD PAYLOAD_FILE=/worker/payload.json ruby draw.rb

# PROD
ENTRYPOINT ["ruby", "func.rb"]
