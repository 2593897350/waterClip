import json
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from app.schemas import DetectRequest, DetectResponse, InpaintRequest, InpaintResponse
from app.services.detect import detect_watermark
from app.services.inpaint import inpaint_image
from app.services.storage import processor_output_dir


def handle_request(method: str, path: str, body: bytes | None = None) -> tuple[int, dict[str, object]]:
    if method == "GET" and path == "/health":
        return HTTPStatus.OK, {"status": "ok"}

    if method == "POST" and path == "/detect":
        payload = DetectRequest.from_dict(json.loads((body or b"{}").decode("utf-8")))
        result = detect_watermark(payload.image_path, processor_output_dir())
        response = DetectResponse(mask_path=result.mask_path, bounds=result.bounds)
        return HTTPStatus.OK, response.to_dict()

    if method == "POST" and path == "/inpaint":
        payload = InpaintRequest.from_dict(json.loads((body or b"{}").decode("utf-8")))
        result = inpaint_image(
            image_path=payload.image_path,
            mask_path=payload.mask_path,
            mode=payload.mode,
            output_dir=processor_output_dir(),
        )
        response = InpaintResponse(output_path=result.output_path)
        return HTTPStatus.OK, response.to_dict()

    return HTTPStatus.NOT_FOUND, {"error": "not found"}


class ProcessorHandler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        self._handle()

    def do_POST(self) -> None:
        self._handle()

    def _handle(self) -> None:
        content_length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(content_length) if content_length else b""

        try:
            status, payload = handle_request(self.command, self.path, body)
        except (KeyError, ValueError, json.JSONDecodeError) as exc:
            status, payload = HTTPStatus.BAD_REQUEST, {"error": str(exc)}
        except Exception as exc:  # pragma: no cover
            status, payload = HTTPStatus.INTERNAL_SERVER_ERROR, {"error": str(exc)}

        encoded = json.dumps(payload).encode("utf-8")
        self.send_response(int(status))
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)

    def log_message(self, format: str, *args: object) -> None:
        return


def main() -> None:
    server = ThreadingHTTPServer(("0.0.0.0", 8000), ProcessorHandler)
    server.serve_forever()


if __name__ == "__main__":
    main()
