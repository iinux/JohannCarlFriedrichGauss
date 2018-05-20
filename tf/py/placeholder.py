#!/usr/bin/python

import tensorflow as tf
import numpy as np

x = tf.placeholder(tf.float32,[None,3])
y = tf.matmul(x,x)
with tf.Session() as sess:
    rand_array = np.random.rand(3,3)
    print(sess.run(y,feed_dict={x:rand_array}))
