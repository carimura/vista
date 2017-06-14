require 'flickraw'
require 'json'
require 'rest-client'

payload_in = JSON.parse(STDIN.read)

FlickRaw.api_key = payload_in["flickr_api_key"]
FlickRaw.shared_secret = payload_in["flickr_api_secret"]

search_text = payload_in["query"] || "baby smile"
num_results = payload_in["num"] || 5

puts "Querying Flickr for \"#{search_text}\" limiting results to #{num_results}"
photos = flickr.photos.search(
	:text => search_text,
	:per_page => num_results,
	:extras => 'original_format',
	:safe_search => 1,
	:content_type => 1
)

puts "Found #{photos.size} images, posting to #{payload_in["func_server_url"]}/detect"
photos.each do |photo|
	image_url = FlickRaw.url_c(photo)

  payload = {:id => photo.id, 
             :image_url => image_url, 
             :func_server_url => payload_in["func_server_url"],
             :aws_bucket => payload_in["aws_bucket"],
             :aws_access => payload_in["aws_access"],
             :aws_secret => payload_in["aws_secret"]
  }

  RestClient.post(payload_in["func_server_url"] + "/detect", payload.to_json, headers={content_type: :json, accept: :json})
end

puts "done"


