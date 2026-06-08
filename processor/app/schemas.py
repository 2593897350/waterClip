from pydantic import BaseModel


class Bounds(BaseModel):
    x: int
    y: int
    width: int
    height: int


class DetectRequest(BaseModel):
    image_path: str


class DetectResponse(BaseModel):
    mask_path: str
    bounds: Bounds


class InpaintRequest(BaseModel):
    image_path: str
    mask_path: str
    mode: str


class InpaintResponse(BaseModel):
    output_path: str
