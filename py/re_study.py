import re

s = "a asda '123 345' sadd"
result = re.split("\s+(?=(?:[^']*'[^']*')*[^']*$)", s)
print(result)
