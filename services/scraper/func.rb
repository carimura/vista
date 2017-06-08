require_relative 'bundle/bundler/setup'

require 'flickraw'
require 'json'
require 'rest-client'

FlickRaw.api_key = "7658b46c6b2c677c40a359bea13b8ec9"
FlickRaw.shared_secret = "5368630b54244b6f"

search_text = "baby smile"

photos = flickr.photos.search(
	:text => search_text, 
	:per_page => 5, 
	:extras => 'original_format', 
	:safe_search => 1,
	:content_type => 1
)

image_urls = Array.new

photos.each do |photo|
	image_url = FlickRaw.url_c(photo)
	image_urls.push(image_url)

  payload = {:id => photo.id, :image_url => image_url}

  #post_url = "https://requestb.in/yk4wumyk"
  post_url = "http://129.146.10.253/r/myapp/detect"
  RestClient.post(post_url, payload, headers={})
end



