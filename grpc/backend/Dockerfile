FROM alpine:3.8

# Install python, pip and virtualenv
RUN apk add --update \
    python \
    python-dev \
    py-pip \
    build-base \
    && pip install -U pip \
    && pip install virtualenv \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY . /app
RUN virtualenv /env && /env/bin/pip install -r /app/requirements.txt

EXPOSE 50051
CMD ["/env/bin/python", "backend.py"]