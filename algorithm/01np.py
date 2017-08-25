import random
import sys

number = random.randint(1, 10)
print ('number')
print (number)
p_arr = [0]
v_arr = [0]
total_v = 0
total_p = 0

for t in range(number):
    p = random.randint(1, 10)
    p_arr.append(p)
    v = random.randint(1, 10)
    v_arr.append(v)
    total_p += p
    total_v += v

# f_arr = [([0] * (total_v + 1)) for i in range(number + 1)]
f_arr = [0] * (total_v + 1)

print ('total volume')
print (total_v)
print ('total price')
print (total_p)
print ('price')
print (p_arr)
print ('volume')
print (v_arr)
print ('')

for i in range(1, number + 1):
    for j in range(total_v, v_arr[i] - 1, -1):
        f_arr[j] = max(f_arr[j], f_arr[j - v_arr[i]] + p_arr[i])

print (f_arr)
