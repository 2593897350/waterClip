from pathlib import Path

from app.services.detect import detect_watermark


def test_detect_watermark_returns_mask_file(tmp_path: Path):
    result = detect_watermark("processor/tests/fixtures/sample_photo.ppm", tmp_path)
    assert result.mask_path.endswith(".pgm")
    assert result.bounds.width >= 0
