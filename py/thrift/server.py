import thriftpy2
pingpong_thrift = thriftpy2.load("pingpong.thrift", module_name="pingpong_thrift")

from thriftpy2.rpc import make_server

class Dispatcher(object):
    def ping(self):
        return "pong"

server = make_server(pingpong_thrift.PingPong, Dispatcher(), '127.0.0.1', 6000)
server.serve()
