from dataclasses import dataclass
from pathlib import Path

from app.schemas import Bounds
from app.services.storage import resolve_input_path


@dataclass
class DetectResult:
    mask_path: str
    bounds: Bounds


def _read_pnm(path: str) -> tuple[str, int, int, int, list[int]]:
    tokens = resolve_input_path(path).read_text().split()
    magic = tokens[0]
    width = int(tokens[1])
    height = int(tokens[2])
    max_value = int(tokens[3])
    values = [int(token) for token in tokens[4:]]
    return magic, width, height, max_value, values


def _write_pgm(path: Path, width: int, height: int, values: list[int]) -> None:
    lines = ["P2", f"{width} {height}", "255"]
    for row in range(height):
        start = row * width
        lines.append(" ".join(str(value) for value in values[start:start + width]))
    path.write_text("\n".join(lines) + "\n")


def detect_watermark(image_path: str, output_dir: Path) -> DetectResult:
    magic, width, height, _max_value, values = _read_pnm(image_path)
    if magic != "P3":
        raise ValueError("Only ASCII PPM images are supported in the test scaffold")

    mask: list[int] = []
    bright_pixels: list[tuple[int, int]] = []
    for index in range(0, len(values), 3):
        x = (index // 3) % width
        y = (index // 3) // width
        pixel = values[index:index + 3]
        is_bright = all(channel >= 235 for channel in pixel)
        mask_value = 255 if is_bright else 0
        mask.append(mask_value)
        if is_bright:
            bright_pixels.append((x, y))

    if bright_pixels:
        min_x = min(point[0] for point in bright_pixels)
        min_y = min(point[1] for point in bright_pixels)
        max_x = max(point[0] for point in bright_pixels)
        max_y = max(point[1] for point in bright_pixels)
        bounds = Bounds(x=min_x, y=min_y, width=max_x - min_x + 1, height=max_y - min_y + 1)
    else:
        bounds = Bounds(x=0, y=0, width=0, height=0)

    output_dir.mkdir(parents=True, exist_ok=True)
    mask_path = output_dir / "detected-mask.pgm"
    _write_pgm(mask_path, width, height, mask)
    return DetectResult(mask_path=str(mask_path), bounds=bounds)
