import os
import json
import sys
import boto3
import urllib
from PIL import Image, ImageDraw


def upload_file(file_name):
    """Upload a file to s3 comptiable storage."""
    s3 = boto3.resource(
        's3',
        aws_access_key_id=os.environ.get("STORAGE_ACCESS_KEY"),
        aws_secret_access_key=os.environ.get("STORAGE_SECRET_KEY"),
        region_name="us-phoenix-1",
        endpoint_url=os.environ.get("MINIO_SERVER_URL"),)

    s3.meta.client.upload_file(
        file_name, os.environ.get("STORAGE_BUCKET"),
        os.path.basename(file_name))


def download_image(image_url, id):
    """Download an image from an http url."""
    file_name = "temp_image"+id+".jpg"
    urllib.urlretrieve(image_url, file_name)
    return file_name


# Draw Rects
def draw_rect(draw, rect, fill=None, width=None):
    """Take draw object and x,y x,y to draw a rect of specific width."""
    cor = (int(rect["startx"]), int(rect["starty"]),
           int(rect["endx"]), int(rect["endy"]))  # (x1,y1, x2,y2)
    line = (cor[0], cor[1], cor[0], cor[3])
    draw.line(line, fill=fill, width=width)
    line = (cor[0], cor[1], cor[2], cor[1])
    draw.line(line, fill=fill, width=width)
    line = (cor[0], cor[3], cor[2], cor[3])
    draw.line(line, fill=fill, width=width)
    line = (cor[2], cor[1], cor[2], cor[3])
    draw.line(line, fill=fill, width=width)


def draw_rects(rects, image_file):
    """Draw the provided rectangles on the image."""
    source_img = Image.open(image_file).convert("RGBA")
    draw = ImageDraw.Draw(source_img)
    for rect in rects:
        draw_rect(draw, rect, "yellow", 5)
    source_img.save(image_file, "PNG")


payload = json.load(sys.stdin)

file_name = download_image(payload["image_url"], payload["id"])

draw_rects(payload["rectangles"], file_name)

upload_file(file_name)
