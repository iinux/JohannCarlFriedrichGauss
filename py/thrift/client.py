import thriftpy2
pingpong_thrift = thriftpy2.load("pingpong.thrift", module_name="pingpong_thrift")

from thriftpy2.rpc import make_client

client = make_client(pingpong_thrift.PingPong, '127.0.0.1', 6000)
print(client.ping())
