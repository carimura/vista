require 'json'
require 'rubygems'
require 'open-uri'
require 'aws-sdk'
require 'mini_magick'

def download_image(payload_in)
  payload = payload_in

  temp_image_name = "temp_image_#{payload["id"]}.jpg"

  File.open(temp_image_name, "wb") do |fout|
    open(payload["image_url"]) do |fin|
      IO.copy_stream(fin, fout)
    end
  end

  temp_image_name
end

def upload_file(image_name, payload_in)
  payload = payload_in

  Aws.config.update({
    endpoint: ENV["MINIO_SERVER_URL"],
    credentials: Aws::Credentials.new(ENV["STORAGE_ACCESS_KEY"], ENV["STORAGE_SECRET_KEY"]),
    force_path_style: true,
    region: 'us-east-1'
  })

  s3 = Aws::S3::Resource.new

  link = nil

	name = File.basename(image_name)
  obj = s3.bucket(ENV["STORAGE_BUCKET"]).object(name)
	obj.upload_file(image_name)

	link = obj.public_url()

	link
end

std_in = STDIN.read
payload = JSON.parse(std_in)

temp_image_name = download_image(payload)

img = MiniMagick::Image.new(temp_image_name)

payload["rectangles"].each do |coords|
  img.combine_options do |c|
    draw_string = "rectangle #{coords["startx"]}, #{coords["starty"]}, #{coords["endx"]}, #{coords["endy"]}"
    c.fill('none')
    is_nude = payload["is_nude"] || "false"
    c.stroke('yellow')
    c.strokewidth(10)
    c.draw draw_string
  end
end

image_name = "image_#{payload["id"]}.jpg"
if payload["resize"]
   img.resize payload["resize"]
end

img.write(image_name)

link = upload_file(image_name, payload)

STDERR.puts "Image link: #{link}"

if ENV["NO_CHAIN"]
    result = { :id => payload["id"], :image_url => link }
    puts result.to_json
end
