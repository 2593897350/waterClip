FROM python:3.12-slim

WORKDIR /app

ARG PIP_INDEX_URL=https://pypi.org/simple
ARG PIP_TRUSTED_HOST=

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY pyproject.toml ./
RUN if [ -n "$PIP_TRUSTED_HOST" ]; then \
      pip install --index-url "$PIP_INDEX_URL" --trusted-host "$PIP_TRUSTED_HOST" fastapi==0.136.3 uvicorn==0.38.0 pydantic==2.13.4; \
    else \
      pip install --index-url "$PIP_INDEX_URL" fastapi==0.136.3 uvicorn==0.38.0 pydantic==2.13.4; \
    fi

COPY app ./app

EXPOSE 8000

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
