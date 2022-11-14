import tensorflow as tf
hello = tf.constant('hello')
sess = tf.Session()
print(sess.run(hello))
