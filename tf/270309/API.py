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
#    subprocess.call(c, shell=True)

mnist = input_data.read_data_sets('./', one_hot=False)

train_x = mnist.train.images.astype('float32')
train_y = mnist.train.labels.astype('int32')
test_x = mnist.test.images.astype('float32')
test_y = mnist.test.labels.astype('int32')

print('训练数据量 %s, %s' %(train_x.shape,train_y.shape))
print('测试数据量 %s, %s' %(test_x.shape,test_y.shape))

import tensorflow.contrib.eager as tfe

# 使用 keras 定义模型（与之前的代码等价）
model = tf.keras.Sequential([
          tf.keras.layers.Dense(512, activation=tf.nn.relu, input_shape=(784,)),
          tf.keras.layers.Dense(256, activation=tf.nn.relu),
          tf.keras.layers.Dense(10)
        ])

from sklearn.metrics import accuracy_score
# 误差函数
def loss(logits, label):
    loss = tf.losses.sparse_softmax_cross_entropy(labels=label, logits=logits)
    return loss
# 日志消息按问题严重性升序排列分别是 DEBUG，INFO，WARN，ERROR，FATAL 。
tf.logging.set_verbosity(tf.logging.ERROR)
# 更新方式
def train(model, x, y, learning_rate, batch_size, epoch):
    # 更新次数
    for e in range(epoch):
        # 创建用于记录训练集误差和准确率的方法
        epoch_loss_avg = tfe.metrics.Mean()
        epoch_accuracy = tfe.metrics.Accuracy()
        # 批量更新
        for b in range(0,len(x),batch_size):
            # 计算梯度
            with tf.GradientTape() as tape:
                loss_value = loss(model(np.array(x[b:b+batch_size])), np.array(y[b:b+batch_size]))
                grads = tape.gradient(loss_value, model.variables)
            # 训练更新
            optimizer = tf.train.GradientDescentOptimizer(learning_rate=learning_rate)
            optimizer.apply_gradients(zip(grads, model.variables),
                            global_step=tf.train.get_or_create_global_step())
            # 因为记忆量限制，记录训练集每个 batch 的误差和准确率
            epoch_loss_avg(loss_value)
            epoch_accuracy(tf.argmax(model(np.array(x[b:b+batch_size])), axis=1, output_type=tf.int32), np.array(y[b:b+batch_size]))
        # 显示
        test_p = model(test_x)
        print("Epoch: %03d | train loss: %.3f | train acc: %.3f | test loss: %.3f | test acc: %.3f" 
              %(e, epoch_loss_avg.result(), epoch_accuracy.result(),
                   loss(test_p, test_y), accuracy_score(tf.argmax(test_p,1), test_y)))

# 训练
r = np.random.permutation(len(train_y))
# 送入打乱顺序的 train_x 和 train_y
train(model, [train_x[i] for i in r], [train_y[i] for i in r], learning_rate = 0.001, batch_size = 128, epoch = 10)

# 评估
test_p = model(test_x)
print("Final Test Loss: %s" %accuracy_score(tf.argmax(test_p,1), test_y))
