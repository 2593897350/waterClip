from dataclasses import dataclass
from pathlib import Path

from app.services.storage import resolve_input_path


@dataclass
class InpaintResult:
    output_path: str


def _read_pnm(path: str) -> tuple[str, int, int, int, list[int]]:
    tokens = resolve_input_path(path).read_text().split()
    magic = tokens[0]
    width = int(tokens[1])
    height = int(tokens[2])
    max_value = int(tokens[3])
    values = [int(token) for token in tokens[4:]]
    return magic, width, height, max_value, values


def _write_ppm(path: Path, width: int, height: int, values: list[int]) -> None:
    lines = ["P3", f"{width} {height}", "255"]
    for row in range(height):
        start = row * width * 3
        lines.append(" ".join(str(value) for value in values[start:start + width * 3]))
    path.write_text("\n".join(lines) + "\n")


def _neighbor_values(
    image_values: list[int],
    mask_values: list[int],
    width: int,
    height: int,
    x: int,
    y: int,
    radius: int,
) -> list[list[int]]:
    colors: list[list[int]] = []
    for delta_y in range(-radius, radius + 1):
        for delta_x in range(-radius, radius + 1):
            nx = x + delta_x
            ny = y + delta_y
            if not (0 <= nx < width and 0 <= ny < height):
                continue
            if mask_values[ny * width + nx] != 0:
                continue
            start = (ny * width + nx) * 3
            colors.append(image_values[start:start + 3])
    return colors


def inpaint_image(image_path: str, mask_path: str, mode: str, output_dir: Path) -> InpaintResult:
    image_magic, width, height, _max_value, image_values = _read_pnm(image_path)
    mask_magic, mask_width, mask_height, _mask_max_value, mask_values = _read_pnm(mask_path)
    if image_magic != "P3" or mask_magic != "P2":
        raise ValueError("Only ASCII PPM/PGM fixtures are supported in the test scaffold")
    if width != mask_width or height != mask_height:
        raise ValueError("Image and mask dimensions must match")

    radius = 1 if mode == "fast" else 2
    output_values = image_values[:]
    for y in range(height):
        for x in range(width):
            if mask_values[y * width + x] == 0:
                continue
            neighbors = _neighbor_values(image_values, mask_values, width, height, x, y, radius)
            fill = [255, 255, 255]
            if neighbors:
                channels = list(zip(*neighbors))
                fill = [sum(channel_values) // len(channel_values) for channel_values in channels]
            start = (y * width + x) * 3
            output_values[start:start + 3] = fill

    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / f"output-{mode}.ppm"
    _write_ppm(output_path, width, height, output_values)
    return InpaintResult(output_path=str(output_path))
