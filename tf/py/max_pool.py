import tensorflow as tf

a = tf.constant([1,3,2,1,2,9,1,1,1,3,2,3,5,6,1,2],dtype=tf.float32,shape=[1,4,4,1])
b = tf.nn.max_pool(a,ksize=[1, 2, 2, 1],strides=[1, 2, 2, 1],padding='VALID')
c = tf.nn.max_pool(a,ksize=[1, 2, 2, 1],strides=[1, 2, 2, 1],padding='SAME')
with tf.Session() as sess:
    print ("b shape:")
    print (b.shape)
    print ("b value:")
    print (sess.run(b))
    print ("c shape:")
    print (c.shape)
    print ("c value:")
    print (sess.run(c))
