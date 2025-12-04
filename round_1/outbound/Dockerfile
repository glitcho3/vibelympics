FROM cgr.dev/chainguard/python:latest-dev as builder

ENV LANG=C.UTF-8
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/movieguess/venv/bin:$PATH"

WORKDIR /movieguess

RUN python -m venv /movieguess/venv
# RUN pip install --no-cache-dir -r requirements.txt

FROM cgr.dev/chainguard/python:latest

ENV TZ="Europa/Madrid"

WORKDIR /movieguess

ENV PYTHONUNBUFFERED=1
ENV PATH="/venv/bin:$PATH"


COPY main.py ./
COPY --from=builder /movieguess/venv /venv

USER 65532


ENTRYPOINT [ "python", "/movieguess/main.py" ]

