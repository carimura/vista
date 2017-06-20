require 'json'
require 'rubygems'
require 'open-uri'
require 'aws-sdk'
require 'mini_magick'

def download_image(payload_in)
  payload = payload_in

  temp_image_name = "temp_image_#{payload["id"]}.jpg"

  unless payload['disable_network']
    File.open(temp_image_name, "wb") do |fout|
      open(payload["image_url"]) do |fin|
        IO.copy_stream(fin, fout)
      end
    end
  end
  
  temp_image_name
end

def upload_file(image_name, payload_in)
  payload = payload_in
 
  s3 = Aws::S3::Resource.new(region: "us-east-1", 
                             credentials: Aws::Credentials.new(payload["access"], 
                                                               payload["secret"]))
  link = nil
  puts "\nUploading the file to s3..."

	name = File.basename(image_name)
  obj = s3.bucket(payload['bucket_name']).object(name)
	obj.upload_file(image_name)

	link = obj.public_url()

	link
end

std_in = STDIN.read
STDERR.puts "std_in --------> " + std_in
payload = JSON.parse(std_in)
puts "payload: " + payload.inspect
puts "Downloading image from " + payload['image_url']

temp_image_name = download_image(payload)

img = MiniMagick::Image.new(temp_image_name)

payload["rectangles"].each do |coords|
  img.combine_options do |c|
    draw_string = "rectangle #{coords["startx"]}, #{coords["starty"]}, #{coords["endx"]}, #{coords["endy"]}"
    puts "draw string: " + draw_string
    c.fill('none')

    is_nude = payload["is_nude"] || "false"
    puts "is_nude: " + is_nude
    if payload["is_nude"] == "true"
      puts "NUDE!!"
      c.stroke('red')
    else
      puts "not nude"
      c.stroke('yellow')
    end

    c.strokewidth(8)
    c.draw draw_string
  end 
end

image_name = "image_#{payload["id"]}.jpg"
img.write(image_name)

link = upload_file(image_name, payload)

puts "link: #{link}"
