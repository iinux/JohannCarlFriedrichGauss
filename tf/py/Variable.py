#!/usr/bin/python

import tensorflow as tf
initial = tf.truncated_normal(shape=[10,10],mean=0,stddev=1)
W=tf.Variable(initial)
list = [[1.,1.],[2.,2.]]
X = tf.Variable(list,dtype=tf.float32)
init_op = tf.global_variables_initializer()
with tf.Session() as sess:
    sess.run(init_op)
    print ("##################(1)################")
    print (sess.run(W))
    print ("##################(2)################")
    print (sess.run(W[:2,:2]))
    op = W[:2,:2].assign(22.*tf.ones((2,2)))
    print ("###################(3)###############")
    print (sess.run(op))
    print ("###################(4)###############")
    print (W.eval(sess)) #computes and returns the value of this variable
    print ("####################(5)##############")
    print (W.eval())  #Usage with the default session
    print ("#####################(6)#############")
    print (W.dtype)
    print (sess.run(W.initial_value))
    print (sess.run(W.op))
    print (W.shape)
    print ("###################(7)###############")
    print (sess.run(X))
