# @zipkin_span(service_name='webapp', span_name='do_stuff')
# def do_stuff():
#     time.sleep(5)
#     headers = create_http_headers_for_new_span()
#     requests.get('http://localhost:6000/service1/', headers=headers)
#     return 'OK'

# @app.route('/')
# def index():
#     with zipkin_span(
#         service_name='webapp',
#         span_name='index',
#         transport_handler=http_transport,
#         port=5000,
#         sample_rate=100, #0.05, # Value between 0.0 and 100.0
#     ):
#         do_stuff()
#         time.sleep(10)
#     return 'OK', 200


from flask import Flask
from py_zipkin.zipkin import zipkin_span, create_http_headers_for_new_span
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


@app.route("/")
def hello():
    with zipkin_span(
        service_name='hello_world_service',
        span_name='hello',
        transport_handler=http_transport,
        port=6000,
        sample_rate=100,  # Value between 0.0 and 100.0
    ):
        zipkin_headers = create_http_headers_for_new_span()
        print zipkin_headers
        time.sleep(0.2)
        requests.get('http://localhost:6001', headers=zipkin_headers)
        time.sleep(0.3)
        return "Hello World!\n", 200


app.run(port=6000)
