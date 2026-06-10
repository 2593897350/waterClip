FROM python:3.12-slim

WORKDIR /app

ARG PIP_INDEX_URL=https://pypi.org/simple
ARG PIP_TRUSTED_HOSTS=

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY pyproject.toml ./
RUN set -eux; \
    set -- pip install --index-url "$PIP_INDEX_URL"; \
    for host in $PIP_TRUSTED_HOSTS; do \
      set -- "$@" --trusted-host "$host"; \
    done; \
    "$@" fastapi==0.136.3 uvicorn==0.38.0 pydantic==2.13.4

COPY app ./app

EXPOSE 8000

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
