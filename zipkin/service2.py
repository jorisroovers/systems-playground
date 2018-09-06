from flask import Flask, request
from py_zipkin.zipkin import zipkin_span, ZipkinAttrs
import requests
import time


app = Flask(__name__)


def http_transport(encoded_span):
    # The collector expects a thrift-encoded list of spans.
    requests.post(
        'http://localhost:9411/api/v1/spans',
        data=encoded_span,
        headers={'Content-Type': 'application/x-thrift'},
    )


def do_stuff(zipkin_attrs):
    with zipkin_span(
        service_name='hello_world_service',
        span_name='doing_stuff',
        zipkin_attrs=zipkin_attrs,
        transport_handler=http_transport,
        port=6001,
        sample_rate=100,  # Value between 0.0 and 100.0
    ):
        print "DOING STUFF"


@app.route("/")
def world():
    zipkin_attrs = ZipkinAttrs(
        trace_id=request.headers['X-B3-TraceID'],
        span_id=request.headers['X-B3-SpanID'],
        parent_span_id=request.headers['X-B3-ParentSpanID'],
        flags=request.headers['X-B3-Flags'],
        is_sampled=request.headers['X-B3-Sampled'],
    )
    with zipkin_span(
        service_name='hello_world_service',
        span_name='world',
        zipkin_attrs=zipkin_attrs,
        transport_handler=http_transport,
        port=6001,
        sample_rate=100,  # Value between 0.0 and 100.0
    ):
        time.sleep(0.4)
        # do_stuff(zipkin_attrs)
        print "WORLD!"

        return "World!\n", 200


app.run(port=6001)
