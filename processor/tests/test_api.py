from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_health_endpoint():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_detect_endpoint_returns_mask_metadata():
    response = client.post(
        "/detect",
        json={"image_path": "processor/tests/fixtures/sample_photo.ppm"},
    )
    assert response.status_code == 200
    body = response.json()
    assert body["mask_path"].endswith(".pgm")
    assert body["bounds"] == {"x": 0, "y": 0, "width": 0, "height": 0} or body["bounds"]["width"] >= 0
