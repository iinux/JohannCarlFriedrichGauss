function myfunction ()
print(debug.traceback("Stack trace"))
print(debug.getinfo(1))
print("Stack trace end")
    return 10
end
myfunction ()
print(debug.getinfo(1))


function newCounter ()
    local n = 0
    local k = 0
    return function ()
        k = n
        n = n + 1
        return n
        end
end

counter = newCounter ()
print(counter())
print(counter())


local i = 1

repeat
    name, val = debug.getupvalue(counter, i)
    if name then
        print ("index", i, name, "=", val)
        if(name == "n") then
            debug.setupvalue (counter,2,10)
        end
        i = i + 1
    end -- if
until not name

while( true )
do
    print(i)
    goto world
    print("hello")
    ::world::
    print("world")
    i = i + 1
    if i > 10 then
        break
    end
end

print(counter())

-- type

print(type("hello"))            --> string
print(type(10.4*3))             --> number
print(type(print))              --> function
print(type(type))               --> function
print(type(true))               --> boolean
print(type(nil))                --> nil
print(type(type(X)))            --> string
print(type(a))                  --> nil
print(type({}))

-- table
-- 不同于其他语言的数组把 0 作为数组的初始索引，在 Lua 里表的默认初始索引一般以 1 开始。
tab1 = { key1 = "val1", key2 = "val2", "val3" }
for k, v in pairs(tab1) do
    print(k .. " - " .. v)
end

tab1.key1 = nil
for k, v in pairs(tab1) do
    print(k .. " - " .. v)
end

print(tab1['ok'])
tab1.ok = 1
print(tab1['ok'])
print(tab1.ok)

-- true / false

if false or nil then
    print("至少有一个是 true")
else
    print("false 和 nil 都为 false")
end

if 0 then
    print("数字 0 是 true")
else
    print("数字 0 为 false")
end

-- assign

html = [[
<html>
<head></head>
<body>
        <a href="http://www.runoob.com/">菜鸟教程</a>
</body>
</html>
]]
print(string.upper(html))
print(#html)
print(string.find(html, "html"))
--print("error" + 1)

a, b, c = 0, 1
print(a,b,c)             --> 0   1   nil

a, b = a+1, b+1, b+2     -- value of b+2 is ignored
print(a,b)               --> 1   2

function add(...)
    local s = 0
    for i, v in ipairs{...} do   --> {...} 表示一个由所有变长参数构成的数组
        s = s + v
    end
    return s
end
print(add(3,4,5,6,7))  --->25

function average(...)
    result = 0
    local arg={...}    --> arg 为一个表，局部变量
    for i,v in ipairs(arg) do
        result = result + v
    end
    print("总共传入 " .. #arg .. " 个数")
    return result/#arg
end

print("平均值为",average(10,5,3,4,5,6))

function average(...)
    print(select(2,...))
    result = 0
    local arg={...}
    for i,v in ipairs(arg) do
       result = result + v
    end
    print("总共传入 " .. select("#",...) .. " 个数")
    return result/select("#",...)
end

print("平均值为",average(10,5,3,4,5,6))

function fwrite(fmt, ...)  ---> 固定的参数fmt
    return io.write(string.format(fmt, ...))
end

fwrite("run\n")                --->fmt = "runoob", 没有变长参数
fwrite("%d %d\n", 1, 1^2^10)   --->fmt = "%d%d", 变长参数
print(string.gsub("aaaa","a","z",3))
print(string.reverse("Lua"))
print(string.char(97,98,99,100))
print(string.byte("ABCD",4))
print(string.byte("ABCD"))
print(string.byte("Lua",-1))
print(string.len("abc"))
print(string.rep("abcd",2))
print(string.format("%d, %q", string.match("I have 2 questions for you.", "(%d+) (%a+)")))
s = "Deadline is 30/05/1999%, firm"
date = "%d%d/%d%d/%d%d%d%d%%"
-- Lua 中的匹配模式直接用常规的字符串来描述。
-- 它用于模式匹配函数 string.find, string.gmatch, string.gsub, string.match。
--[[
    单个字符(除 ^$()%.[]*+-? 外): 与该字符自身配对

.(点): 与任何字符配对
%a: 与任何字母配对
%c: 与任何控制符配对(例如\n)
%d: 与任何数字配对
%l: 与任何小写字母配对
%p: 与任何标点(punctuation)配对
%s: 与空白字符配对
%u: 与任何大写字母配对
%w: 与任何字母/数字配对
%x: 与任何十六进制数配对
%z: 与任何代表0的字符配对
%x(此处x是非字母非数字字符): 与字符x配对. 主要用来处理表达式中有功能的字符(^$()%.[]*+-?)的配对问题, 例如%%与%配对
[数个字符类]: 与任何[]中包含的字符类配对. 例如[%w_]与任何字母/数字, 或下划线符号(_)配对
[^数个字符类]: 与任何不包含在[]中的字符类配对. 例如[^%s]与任何非空白字符配对
当上述的字符类用大写书写时, 表示与非此字符类的任何字符配对. 例如, %S表示与任何非空白字符配对.例如，'%A'非字母的字符:
'%' 用作特殊字符的转义字符，因此 '%.' 匹配点；'%%' 匹配字符 '%'。转义字符 '%'不仅可以用来转义特殊字符
]]
print(string.sub(s, string.find(s, date)))    --> 30/05/1999

for word in string.gmatch("Hello Lua user", "%a+")
do
    print(word)
end

if 1~=2 then
    print("1~=2")
end

-- 字符串
local sourcestr = "prefix--runoobgoogletaobao--suffix"
print("原始字符串", string.format("%q", sourcestr))

-- 截取部分，第4个到第15个
local first_sub = string.sub(sourcestr, 4, 15)
print("第一次截取", string.format("%q", first_sub))

-- 截取最后10个
local third_sub = string.sub(sourcestr, -10)
print("第二次截取", string.format("%q", third_sub))

-- 索引越界，输出原始字符串
local fourth_sub = string.sub(sourcestr, -100)
print("第三次截取", string.format("%q", fourth_sub))

-- 初始化数组
array = {}
maxRows = 3
maxColumns = 3
for row=1,maxRows do
    for col=1,maxColumns do
        array[row*maxColumns +col] = row*col
    end
end

-- 访问数组
for row=1,maxRows do
    for col=1,maxColumns do
        print(array[row*maxColumns +col])
    end
end

function square(iteratorMaxCount,currentNumber)
    if currentNumber<iteratorMaxCount
    then
        currentNumber = currentNumber+1
        return currentNumber, currentNumber*currentNumber
    end
 end

 for i,n in square,3,0
 do
    print(i,n)
 end

fruits = {"banana","orange","apple"}
-- 返回 table 连接后的字符串
print("连接后的字符串 ",table.concat(fruits))

-- 指定连接字符
print("连接后的字符串 ",table.concat(fruits,", "))

-- 指定索引来连接 table
print("连接后的字符串 ",table.concat(fruits,", ", 2,3))

fruits = {"banana","orange","apple"}

-- 在末尾插入
table.insert(fruits,"mango")
print("索引为 4 的元素为 ",fruits[4])

-- 在索引为 2 的键处插入
table.insert(fruits,2,"grapes")
print("索引为 2 的元素为 ",fruits[2])

print("最后一个元素为 ",fruits[5])
table.remove(fruits)
print("移除后最后一个元素为 ",fruits[5])

fruits = {"banana","orange","apple","grapes"}
print("排序前")
for k,v in ipairs(fruits) do
        print(k,v)
end

table.sort(fruits)
print("排序后")
for k,v in ipairs(fruits) do
        print(k,v)
end

function table_maxn(t)
    local mn=nil;
    for k, v in pairs(t) do
        if(mn==nil) then
            mn=v
        end
        if mn < v then
            mn = v
        end
    end
    return mn
end
tbl = {[1] = 2, [2] = 6, [3] = 34, [26] =5}
print("tbl 最大值：", table_maxn(tbl))
print("tbl 长度 ", #tbl)
function table_leng(t)
    local leng=0
    for k, v in pairs(t) do
        leng=leng+1
    end
    return leng;
end
print(table_leng(tbl))
