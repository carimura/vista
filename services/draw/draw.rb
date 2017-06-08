require_relative 'bundle/bundler/setup'

require 'rubygems'
require 'open-uri'
require 'aws'
require 'mini_magick'
require 'pubnub'

def download_image
  payload = ""

  filename = 'input.jpg'
  unless payload['disable_network']
    filepath = filename
    File.open(filepath, 'wb') do |fout|
      open(payload['image_url']) do |fin|
        IO.copy_stream(fin, fout)
      end
    end
  end
  filename
end

def upload_file(filename)
  payload = ""

  link = nil
  unless payload['disable_network']
    filepath = filename
    puts "\nUploading the file to s3..."
    s3 = Aws::S3Interface.new(payload['aws_access'], payload['aws_secret'])
    #s3.create_bucket(payload['aws_s3_bucket_name'])
    response = s3.put(payload['aws_s3_bucket_name'], filename, File.open(filepath))
    if response == true
      puts "Uploading succesful."
      link = s3.get_link(payload['aws_s3_bucket_name'], filename)
      puts "\nYou can view the file here on s3:", link
    else
      puts "Error placing the file in s3."
    end
    puts "-"*60
  end
  link
end

payload = ""
puts "payload: #{payload}"
puts "Downloading image from " + payload['image_url']

filename = download_image()

img = MiniMagick::Image.new("input.jpg")

payload["rectangles"].each do |coords|
  img.combine_options do |c|                                                                                  
    draw_string = "rectangle #{coords["startx"]}, #{coords["starty"]}, #{coords["endx"]}, #{coords["endy"]}"
    puts "draw string: " + draw_string
    c.fill('none')

    puts "is_nude: " + payload["is_nude"]
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

filename = "#{payload["id"]}.jpg"

img.write(filename)

link = upload_file(filename)

puts "link: #{link}"
