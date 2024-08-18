from py2neo import Graph, Node, Relationship
import neo4j_config

# 连接到Neo4j数据库
graph = Graph(neo4j_config.address, auth=(neo4j_config.user, neo4j_config.password))

# 创建用户节点
user = Node("User", name="Alice")

# 创建交易节点
transaction = Node("Transaction", id="tx123", amount=100, currency="USD")

# 创建设备节点
device = Node("Device", id="dev001", type="mobile")

user2 = Node("User", name="user2")

# 创建关系
user_transaction = Relationship(user, "MADE", transaction)
transaction_device = Relationship(transaction, "USING", device)
kill_event = Relationship(user, "KILL", user2)

# 将节点和关系添加到图数据库
#graph.create(user_transaction)
#graph.create(transaction_device)
#graph.create(kill_event)

# 查询图数据库
result = graph.run("MATCH (u:User)-[:MADE]->(t:Transaction) RETURN u, t")
for record in result:
    print(record)

result = graph.run("MATCH (start:User {name: 'Alice'}), (end:Device {id: 'dev001'}) \
MATCH path = shortestPath((start)-[*]-(end)) \
RETURN path")
print(result)

result = graph.run("MATCH (start:User {name: 'Alice'}), (end:User {name: 'user2'}) \
MATCH path = shortestPath((start)-[*]-(end)) \
RETURN path")
print(result)

if __name__ == '__main__':
    pass
