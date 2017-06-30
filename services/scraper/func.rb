require 'flickraw'
require 'json'
require 'rest-client'

payload_in = JSON.parse(STDIN.read)

FlickRaw.api_key = payload_in["flickr_api_key"]
FlickRaw.shared_secret = payload_in["flickr_api_secret"]

search_text = payload_in["query"] || "baby smile"
num_results = payload_in["num"] || 5
service_to_call = payload_in["service_to_call"] || "detect-faces"

puts "Querying Flickr for \"#{search_text}\" limiting results to #{num_results}"
photos = flickr.photos.search(
	:text => search_text,
	:per_page => num_results,
	:extras => 'original_format',
	:safe_search => 1,
	:content_type => 1
)

puts "Found #{photos.size} images, posting to #{payload_in["func_server_url"]}/#{service_to_call}"
threads = []

blacklist_photos = ['35331390846']

photos.each do |photo|
  if blacklist_photos.include?(photo.id)
    image_url = "https://farm3.staticflickr.com/2175/5714544755_e5dc8e6ede_b.jpg"
  else
    image_url = FlickRaw.url_c(photo)
  end
  payload = {:id => photo.id, 
             :image_url => image_url, 
             :func_server_url => payload_in["func_server_url"],
             :bucket => payload_in["bucket"],
             :access => payload_in["access"],
             :secret => payload_in["secret"],
             :flickr_api_key => payload_in["flickr_api_key"],
             :flickr_api_secret => payload_in["flickr_api_secret"],
             :pubnub_subscribe_key => payload_in["pubnub_subscribe_key"],
             :pubnub_publish_key => payload_in["pubnub_publish_key"],
             :algorithmia_key => payload_in["algorithmia_key"]
  }

  threads <<  Thread.new(payload, payload_in) { |payload, payload_in| 
    RestClient.post(payload_in["func_server_url"] + "/" + service_to_call, payload.to_json, headers={content_type: :json, accept: :json})
  }
end

threads.each do |t|
    t.join
end

puts "done"


