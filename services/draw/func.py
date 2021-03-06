import os
import ujson
import boto3
import fdk

from PIL import Image
from PIL import ImageDraw

from fdk import fixtures
from urllib import request

import ssl

ssl._create_default_https_context = ssl._create_unverified_context

s3 = boto3.resource(
    's3',
    aws_access_key_id=os.environ.get("STORAGE_ACCESS_KEY"),
    aws_secret_access_key=os.environ.get("STORAGE_SECRET_KEY"),
    region_name=os.environ.get("S3_REGION", "us-east-1"),
    endpoint_url=os.environ.get("MINIO_SERVER_URL")
)

service_s3 = boto3.client('s3',
  aws_access_key_id=os.environ.get("STORAGE_ACCESS_KEY"),
  aws_secret_access_key=os.environ.get("STORAGE_SECRET_KEY"),
  region_name=os.environ.get("S3_REGION", "us-east-1"),
  endpoint_url=os.environ.get("MINIO_SERVER_URL")
)


async def test_override_content_type(aiohttp_client):
    with open("payload.json", "r") as payload_file:
        call = await fixtures.setup_fn_call(
            aiohttp_client, handle, json=ujson.load(payload_file))
        content, status, headers = await call
        assert 200 == status
        assert "image_url" in ujson.loads(content)


def upload_file(file_name):
    """Upload a file to s3 compatible storage."""
    bucket = os.environ.get("STORAGE_BUCKET", "oracle-vista-out")
    key = os.path.basename(file_name)
    with open(file_name, 'rb') as f:
        s3.Bucket(bucket).put_object(Key=key, Body=f)

    return service_s3.generate_presigned_url(
        ClientMethod='get_object',
        Params={
            'Bucket': bucket,
            'Key': key
        }
    )

def download_image(image_url, id):
    """Download an image from an http url."""
    file_name = "/tmp/temp_image"+id+".jpg"
    request.urlretrieve(image_url, file_name)
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


def handle(ctx, data=None, **kwargs):
    if data and len(data) > 0:
        payload = ujson.loads(data)
        file_name = download_image(payload.get("image_url"), payload.get("id"))
        draw_rects(payload["rectangles"], file_name)
        file_url = upload_file(file_name)
        return {
            "id": payload.get("id"),
            "image_url": file_url,
            "plate": payload.get("plate")
        }


if __name__ == "__main__":
    fdk.handle(handle)
