# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import book_store_pb2 as book__store__pb2


class BookStoreStub(object):
  """(Method definitions not shown)
  """

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.GetBook = channel.unary_unary(
        '/BookStore/GetBook',
        request_serializer=book__store__pb2.BookReference.SerializeToString,
        response_deserializer=book__store__pb2.Book.FromString,
        )


class BookStoreServicer(object):
  """(Method definitions not shown)
  """

  def GetBook(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_BookStoreServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'GetBook': grpc.unary_unary_rpc_method_handler(
          servicer.GetBook,
          request_deserializer=book__store__pb2.BookReference.FromString,
          response_serializer=book__store__pb2.Book.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'BookStore', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
