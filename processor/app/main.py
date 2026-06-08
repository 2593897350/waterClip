from fastapi import FastAPI

from app.schemas import DetectRequest, DetectResponse, InpaintRequest, InpaintResponse
from app.services.detect import detect_watermark
from app.services.inpaint import inpaint_image
from app.services.storage import processor_output_dir

app = FastAPI()


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/detect", response_model=DetectResponse)
def detect(payload: DetectRequest) -> DetectResponse:
    result = detect_watermark(payload.image_path, processor_output_dir())
    return DetectResponse(mask_path=result.mask_path, bounds=result.bounds)


@app.post("/inpaint", response_model=InpaintResponse)
def inpaint(payload: InpaintRequest) -> InpaintResponse:
    result = inpaint_image(
        image_path=payload.image_path,
        mask_path=payload.mask_path,
        mode=payload.mode,
        output_dir=processor_output_dir(),
    )
    return InpaintResponse(output_path=result.output_path)
