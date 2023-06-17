import torch
print(torch.cuda.is_available())
x = torch.rand(5, 3)
print(x)
x = torch.empty(5, 3)
print(x)
zero_x = torch.zeros(5, 3, dtype=torch.long)
print(zero_x)
one_x = torch.ones(5, 3, dtype=torch.long)
print(one_x)

tensor1 = torch.tensor([5.5, 3])
print(tensor1)
tensor2 = tensor1.new_ones(5, 3, dtype=torch.double)  # new_* 方法需要输入 tensor 大小
print(tensor2)

tensor3 = torch.randn_like(tensor2, dtype=torch.float)
print('tensor3: ', tensor3)
