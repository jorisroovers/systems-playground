from __future__ import print_function

from concurrent import futures
from datetime import datetime
import time

import grpc
import os

import book_store_pb2_grpc
import book_store_pb2


_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class BookStore(book_store_pb2_grpc.BookStoreServicer):
    """ Class to actually handle Bookstore requests"""

    def __init__(self, backend_name, hostname):
        self.backend_name = backend_name
        self.hostname = hostname

    def GetBook(self, request, context):
        book_title = "My Book - Served from backend '{}' (hostname: {}) at {:%Y-%m-%d %H:%M:%S}".format(
            self.backend_name, self.hostname, datetime.now())
        return book_store_pb2.Book(title=book_title, author="My Author", year=2018)


def serve():
    print("Initializing...\n")
    backend_name = os.environ.get("BACKEND_NAME", "Default backend")
    print("Backend Name ($BACKEND_NAME):", backend_name)
    hostname = os.environ.get("HOSTNAME", "Default backend")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    book_store_pb2_grpc.add_BookStoreServicer_to_server(BookStore(backend_name, hostname), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Starting server (on port 50051)...")
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
