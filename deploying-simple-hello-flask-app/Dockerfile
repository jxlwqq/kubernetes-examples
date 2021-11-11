# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM python:alpine

WORKDIR /
ADD . /
RUN pip install -r requirements.txt

EXPOSE 5000
CMD ["python", "app.py"]