from dowhen import do
from dowhen import bp
from dowhen import goto
from dowhen import when

# https://github.com/gaogaotiantian/dowhen

def f(x):
    x += 100
    # Let's change the value of x before return
    return x

# do("x = 1") is the callback
# when(f, "return x") is the trigger
# This is equivalent to:
# handler = when(f, "return x").do("x = 1")
handler = do("x = 1").when(f, "return x")
# x = 1 is executed before "return x"
print(f(0))

# You can remove the handler
handler.remove()
print(f(0))

# This will skip the line of `x += 100`
# The handler will be removed after the with context
with goto("return x").when(f, "x += 100"):
    print(f(0))

# You can chain callbacks and they'll run in order at the trigger
# You don't need to store the handler if you don't use it
when(f, "x += 100").goto("return x").do("x = 42")
print(f(0))

# bp() is another callback that brings up pdb
# handler = bp().when(f, "return x")
# This will enter pdb
# f(0)
# You can temporarily disable the handler
# handler.enable() will enable it again
# handler.disable()



if __name__ == '__main__':
    pass
