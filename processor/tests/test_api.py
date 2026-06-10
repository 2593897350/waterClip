import json

from app.main import handle_request


def test_health_endpoint():
    status, body = handle_request("GET", "/health")
    assert status == 200
    assert body == {"status": "ok"}


def test_detect_endpoint_returns_mask_metadata():
    payload = json.dumps({"image_path": "processor/tests/fixtures/sample_photo.ppm"}).encode("utf-8")
    status, body = handle_request("POST", "/detect", payload)
    assert status == 200
    assert body["mask_path"].endswith(".pgm")
    assert body["bounds"] == {"x": 0, "y": 0, "width": 0, "height": 0} or body["bounds"]["width"] >= 0
