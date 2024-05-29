from abc import ABCMeta, abstractmethod
from collections import abc, MutableSequence


class Payment(metaclass=ABCMeta):
    @abstractmethod  # 调用@abstractmethod规定子类必须有pay方法
    def pay(self, money):
        pass


class Wechatpay(Payment):
    def pay(self, money):
        print('微信支付了%s元' % money)


class Struggle:
    def __len__(self):
        return 23


class Foo(MutableSequence):
    pass


if __name__ == '__main__':
    obj = Wechatpay()
    obj.pay(1)
    # True
    print(isinstance(Struggle(), abc.Sized))
    a = Foo()
