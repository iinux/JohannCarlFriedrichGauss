#!/usr/bin/env python3
import random

def generate_math_problems(num_problems=30, max_num=20):
    """生成指定数量的加减法题目

    Args:
        num_problems: 题目数量，默认30道
        max_num: 最大数字，默认20以内
    """
    problems = []

    for i in range(num_problems):
        # 随机选择运算符：+ 或 -
        operator = random.choice(['+', '-'])

        # 随机生成两个数字
        num1 = random.randint(0, max_num)
        num2 = random.randint(0, max_num)

        # 确保减法结果不为负数
        if operator == '-' and num1 < num2:
            num1, num2 = num2, num1

        # 生成题目
        problem = f"{i+1}. {num1} {operator} {num2} ="
        problems.append(problem)

    return problems

def main():
    print("=" * 50)
    print("20以内加减法练习题")
    print("=" * 50)

    # 生成30道题目
    problems = generate_math_problems(30, 20)

    # 打印题目
    for problem in problems:
        print(problem)

    print("\n" + "=" * 50)
    print("祝你好运！")
    print("=" * 50)

if __name__ == "__main__":
    main()