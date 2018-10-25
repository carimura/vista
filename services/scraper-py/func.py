import fdk
import ujson
import os
import random
import flickrapi
import ssl

from fdk import fixtures

ssl._create_default_https_context = ssl._create_unverified_context


flickr = flickrapi.FlickrAPI(
    os.environ.get("FLICKR_API_KEY"),
    os.environ.get("FLICKR_API_SECRET"),
    token_cache_location='/tmp',
    format='parsed-json'
)

PHOTO_SOURCE_URL = 'https://farm{0}.staticflickr.com/{1}/{2}_{3}{4}.{5}'


async def test_override_content_type(aiohttp_client):
    with open("payload.json", "r") as payload_file:
        call = await fixtures.setup_fn_call(
            aiohttp_client, handler, json=ujson.load(payload_file))
        content, status, headers = await call
        data = ujson.loads(content)
        assert 200 == status
        assert "result" in data
        assert len(data.get("result")) > 0


def get_image_url(photo_dict):
    return PHOTO_SOURCE_URL.format(
        photo_dict['farm'], photo_dict['server'],
        photo_dict['id'], photo_dict['secret'],
        '_c', 'jpg'
    )


def photo_to_payload(body, photo_dict):
    return {
        "id": photo_dict.get('id'),
        "image_url": get_image_url(photo_dict),
        "countrycode": body.get("countrycode"),
        "bucket": body.get("bucket", "")
    }


def handler(ctx, data=None, loop=None):
    payloads = []
    if data and len(data) > 0:
        body = ujson.loads(data)
        photos = flickr.photos.search(
            text=body.get("query", "baby smile"),
            per_page=int(body.get("num", "5")),
            page=int(body.get("page", int(random.uniform(1, 50)))),
            extras="original_format",
            safe_search="1",
            content_type="1",
        )

        for p in photos.get('photos', {'photo': {}}).get('photo', []):
            payloads.append(photo_to_payload(body, p))

    return {"result": payloads}


if __name__ == "__main__":
    fdk.handle(handler)
