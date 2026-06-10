from dataclasses import asdict, dataclass
from typing import Any


@dataclass
class Bounds:
    x: int
    y: int
    width: int
    height: int

    def to_dict(self) -> dict[str, int]:
        return asdict(self)


@dataclass
class DetectRequest:
    image_path: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]) -> "DetectRequest":
        return cls(image_path=str(payload["image_path"]))


@dataclass
class DetectResponse:
    mask_path: str
    bounds: Bounds

    def to_dict(self) -> dict[str, Any]:
        return {"mask_path": self.mask_path, "bounds": self.bounds.to_dict()}


@dataclass
class InpaintRequest:
    image_path: str
    mask_path: str
    mode: str

    @classmethod
    def from_dict(cls, payload: dict[str, Any]) -> "InpaintRequest":
        return cls(
            image_path=str(payload["image_path"]),
            mask_path=str(payload["mask_path"]),
            mode=str(payload["mode"]),
        )


@dataclass
class InpaintResponse:
    output_path: str

    def to_dict(self) -> dict[str, str]:
        return {"output_path": self.output_path}
