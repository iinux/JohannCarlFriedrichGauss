import dbus

# 创建一个DBus会话连接
bus = dbus.SessionBus()

# 获取要通信的目标应用程序的代理对象
target_app = bus.get_object('org.example.targetapp', '/org/example/targetapp')
# 'org.example.targetapp'是目标应用程序的DBus名称，'/org/example/targetapp'是目标对象的路径

# 调用目标应用程序的方法
response = target_app.HelloWorld('Hello from DBus!')
# 'HelloWorld'是目标应用程序提供的方法名称，'Hello from DBus!'是要传递的参数

# 打印返回的响应
print(response)
