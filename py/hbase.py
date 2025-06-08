import happybase

connection = happybase.Connection(host='10.190.14.16', port=2181)
table = connection.table('nn:t2')

table.put('row_key', {'f1:name': 'value'})

row = table.row('row_key')
print(row)

scanner = table.scan()
for key, data in scanner:
    print(key, data)
connection.close()
