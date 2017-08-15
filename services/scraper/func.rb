require 'flickraw'
require 'json'
require 'rest-client'

payload_in = JSON.parse(STDIN.read)

FlickRaw.api_key = ENV["FLICKR_API_KEY"]
FlickRaw.shared_secret = ENV["FLICKR_API_SECRET"]

search_text = payload_in["query"] || "baby smile"
num_results = payload_in["num"] || 5
service_to_call = payload_in["service_to_call"] || "detect-faces"
page = payload_in["page"] || rand(50)

puts "Querying Flickr for \"#{search_text}\" grabbing from page #{page} limiting results to #{num_results}"

photos = flickr.photos.search(
	:text => search_text,
	:per_page => num_results,
  :page => page,
	:extras => 'original_format',
	:safe_search => 1,
	:content_type => 1
)

puts "Found #{photos.size} images, posting to #{ENV["FUNC_SERVER_URL"]}/#{service_to_call}"
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
             :countrycode => payload_in["countrycode"],
             :bucket => payload_in["bucket"]
  }

  RestClient.post(ENV["FUNC_SERVER_URL"] + "/" + service_to_call, payload.to_json, headers={content_type: :json, accept: :json})
end

puts "done"


