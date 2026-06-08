from pathlib import Path

from app.services.inpaint import inpaint_image


def test_inpaint_fast_mode_creates_output(tmp_path: Path):
    output = inpaint_image(
        image_path="processor/tests/fixtures/sample_photo.ppm",
        mask_path="processor/tests/fixtures/sample_mask.pgm",
        mode="fast",
        output_dir=tmp_path,
    )
    assert output.output_path.endswith(".ppm")
