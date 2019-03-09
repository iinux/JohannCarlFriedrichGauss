#!/usr/bin/python
#coding:utf-8
# 国内用户
import subprocess
import tensorflow as tf
import numpy as np
tf.enable_eager_execution()
from tensorflow.examples.tutorials.mnist import input_data

# 下载 mnist 数据集
commands = ['wget https://devlab-1251520893.cos.ap-guangzhou.myqcloud.com/t10k-images-idx3-ubyte.gz',
            'wget https://devlab-1251520893.cos.ap-guangzhou.myqcloud.com/t10k-labels-idx1-ubyte.gz',
            'wget https://devlab-1251520893.cos.ap-guangzhou.myqcloud.com/train-images-idx3-ubyte.gz',
            'wget https://devlab-1251520893.cos.ap-guangzhou.myqcloud.com/train-labels-idx1-ubyte.gz']
#for c in commands:
    #subprocess.call(c, shell=True)

mnist = input_data.read_data_sets('./', one_hot=False)

train_x = mnist.train.images.astype('float32')
train_y = mnist.train.labels.astype('int32')
test_x = mnist.test.images.astype('float32')
test_y = mnist.test.labels.astype('int32')

print('训练数据量 %s, %s' %(train_x.shape,train_y.shape))
print('测试数据量 %s, %s' %(test_x.shape,test_y.shape))

import tensorflow.contrib.eager as tfe
# 定义模型类
class Model(object):
    # 初始化方法
    def inits(self, shape):
        return tf.random_uniform(shape,
            minval=-np.sqrt(5) * np.sqrt(1.0 / shape[0]),
            maxval=np.sqrt(5) * np.sqrt(1.0 / shape[0]))
    def __init__(self):
        # 参数初始化
        self.W1 = tfe.Variable(self.inits([784,256]))
        self.b1 = tfe.Variable(self.inits([256]))
        self.W2 = tfe.Variable(self.inits([256,128]))
        self.b2 = tfe.Variable(self.inits([128]))
        self.W = tfe.Variable(self.inits([128,10]))
        self.b = tfe.Variable(self.inits([10]))
    def __call__(self, x):
        # 正向传递
        # tf.nn.relu 是非线性函数
        # tf.matmul 是矩阵乘法，输入放在前，参数放在后
        h1 = tf.nn.relu(tf.matmul(x, self.W1) + self.b1)
        h2 = tf.nn.relu(tf.matmul(h1, self.W2) + self.b2)
        y = tf.matmul(h2, self.W) + self.b
        return y
# 实例模型
model = Model()

from sklearn.metrics import accuracy_score
# 误差函数
def loss(logits, label):
    loss = tf.losses.sparse_softmax_cross_entropy(labels=label, logits=logits)
    return loss
# 日志消息按问题严重性升序排列分别是 DEBUG，INFO，WARN，ERROR，FATAL 。
tf.logging.set_verbosity(tf.logging.ERROR)
from sklearn.metrics import accuracy_score

# 日志消息按问题严重性升序排列分别是 DEBUG，INFO，WARN，ERROR，FATAL 。
tf.logging.set_verbosity(tf.logging.ERROR)
# 更新方式
def train(model, x, y, learning_rate, batch_size, epoch):
    # 更新次数
    for e in range(epoch):
        # 批量更新
        for b in range(0,len(x),batch_size):
            # 计算梯度
            with tf.GradientTape() as tape:
                loss_value = loss(model(np.array(x[b:b+batch_size])), np.array(y[b:b+batch_size]))
                dW1, db1, dW2, db2, dW, db = tape.gradient(loss_value, 
                                   [model.W1, model.b1, model.W2, model.b2, model.W, model.b])
            # 训练更新
            model.W1.assign_sub(dW1 * learning_rate)
            model.b1.assign_sub(db1 * learning_rate)
            model.W2.assign_sub(dW2 * learning_rate)
            model.b2.assign_sub(db2 * learning_rate)
            model.W.assign_sub(dW * learning_rate)
            model.b.assign_sub(db * learning_rate)

        # 显示
        # 训练集太大，内存不够，所以每次以 batch size 个来的计算预测值
        train_p = tf.concat([model(train_x[b:b+batch_size]) for b in range(0,len(train_x),batch_size)],axis=0)
        test_p = model(test_x)
        print("Epoch: %03d | train loss: %.3f | train acc: %.3f | test loss: %.3f | test acc: %.3f" 
              %(e, loss(train_p, train_y), accuracy_score(tf.argmax(train_p,1), train_y),
                   loss(test_p, test_y), accuracy_score(tf.argmax(test_p,1), test_y)))

# 训练
r = np.random.permutation(len(train_y))
# 送入打乱顺序的 train_x 和 train_y
train(model, [train_x[i] for i in r], [train_y[i] for i in r], learning_rate = 0.001, batch_size = 128, epoch = 10)

# 评估
test_p = model(test_x)
print("Final Test Loss: %s" %accuracy_score(tf.argmax(test_p,1), test_y))
