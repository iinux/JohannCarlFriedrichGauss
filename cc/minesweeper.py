#!/usr/bin/env python3
import tkinter as tk
from tkinter import messagebox
import random

class Minesweeper:
    def __init__(self, root, rows=9, cols=9, mines=10):
        self.root = root
        self.root.title("扫雷游戏")
        self.root.resizable(False, False)

        # 游戏参数
        self.rows = rows
        self.cols = cols
        self.mines = mines
        self.cell_size = 30

        # 游戏状态
        self.game_over = False
        self.first_click = True
        self.mines_left = mines
        self.timer_running = False
        self.time_elapsed = 0

        # 创建游戏界面
        self.setup_ui()
        self.new_game()

    def setup_ui(self):
        """设置用户界面"""
        # 顶部信息栏
        info_frame = tk.Frame(self.root, bg='#c0c0c0')
        info_frame.pack(fill=tk.X, padx=5, pady=5)

        # 剩余地雷数显示
        self.mines_label = tk.Label(
            info_frame,
            text=f"💣 {self.mines_left:03d}",
            font=('Arial', 12, 'bold'),
            bg='#c0c0c0',
            width=10
        )
        self.mines_label.pack(side=tk.LEFT, padx=10)

        # 重新开始按钮
        self.restart_button = tk.Button(
            info_frame,
            text="😊",
            font=('Arial', 16),
            command=self.new_game,
            width=3,
            height=1
        )
        self.restart_button.pack(side=tk.LEFT, padx=20)

        # 计时器显示
        self.timer_label = tk.Label(
            info_frame,
            text="⏱ 000",
            font=('Arial', 12, 'bold'),
            bg='#c0c0c0',
            width=10
        )
        self.timer_label.pack(side=tk.LEFT, padx=10)

        # 游戏画布
        canvas_frame = tk.Frame(self.root, bg='#c0c0c0')
        canvas_frame.pack(padx=5, pady=5)

        self.canvas = tk.Canvas(
            canvas_frame,
            width=self.cols * self.cell_size,
            height=self.rows * self.cell_size,
            bg='#c0c0c0',
            highlightthickness=2,
            highlightbackground='#808080'
        )
        self.canvas.pack()

        # 绑定鼠标事件
        self.canvas.bind("<Button-1>", self.left_click)
        self.canvas.bind("<Button-3>", self.right_click)

    def new_game(self):
        """开始新游戏"""
        self.game_over = False
        self.first_click = True
        self.mines_left = self.mines
        self.time_elapsed = 0
        self.timer_running = False

        # 更新显示
        self.mines_label.config(text=f"💣 {self.mines_left:03d}")
        self.timer_label.config(text="⏱ 000")
        self.restart_button.config(text="😊")

        # 初始化游戏板
        self.board = [[0 for _ in range(self.cols)] for _ in range(self.rows)]
        self.revealed = [[False for _ in range(self.cols)] for _ in range(self.rows)]
        self.flagged = [[False for _ in range(self.cols)] for _ in range(self.rows)]
        self.buttons = [[None for _ in range(self.cols)] for _ in range(self.rows)]

        # 重新绘制游戏板
        self.draw_board()

    def place_mines(self, avoid_row, avoid_col):
        """放置地雷，避开第一次点击的位置"""
        mines_placed = 0
        while mines_placed < self.mines:
            row = random.randint(0, self.rows - 1)
            col = random.randint(0, self.cols - 1)

            # 避开第一次点击的位置及其周围
            if abs(row - avoid_row) <= 1 and abs(col - avoid_col) <= 1:
                continue

            if self.board[row][col] != -1:  # -1 表示地雷
                self.board[row][col] = -1
                mines_placed += 1

        # 计算每个格子周围的地雷数
        for row in range(self.rows):
            for col in range(self.cols):
                if self.board[row][col] != -1:
                    count = self.count_adjacent_mines(row, col)
                    self.board[row][col] = count

    def count_adjacent_mines(self, row, col):
        """计算指定格子周围的地雷数量"""
        count = 0
        for dr in [-1, 0, 1]:
            for dc in [-1, 0, 1]:
                if dr == 0 and dc == 0:
                    continue
                new_row, new_col = row + dr, col + dc
                if (0 <= new_row < self.rows and
                    0 <= new_col < self.cols and
                    self.board[new_row][new_col] == -1):
                    count += 1
        return count

    def draw_board(self):
        """绘制游戏板"""
        self.canvas.delete("all")
        for row in range(self.rows):
            for col in range(self.cols):
                x1 = col * self.cell_size
                y1 = row * self.cell_size
                x2 = x1 + self.cell_size
                y2 = y1 + self.cell_size

                # 创建按钮效果
                self.buttons[row][col] = self.canvas.create_rectangle(
                    x1+2, y1+2, x2-2, y2-2,
                    fill='#bdbdbd',
                    outline='white',
                    width=2
                )

                # 如果已标记，显示旗帜
                if self.flagged[row][col]:
                    self.canvas.create_text(
                        x1 + self.cell_size//2,
                        y1 + self.cell_size//2,
                        text='🚩',
                        font=('Arial', 14)
                    )
                # 如果已揭开且不是地雷，显示数字
                elif self.revealed[row][col] and self.board[row][col] >= 0:
                    self.canvas.create_rectangle(
                        x1, y1, x2, y2,
                        fill='#e0e0e0',
                        outline='#808080'
                    )
                    if self.board[row][col] > 0:
                        colors = ['', 'blue', 'green', 'red', 'purple',
                                'maroon', 'turquoise', 'black', 'gray']
                        color = colors[min(self.board[row][col], 8)]
                        self.canvas.create_text(
                            x1 + self.cell_size//2,
                            y1 + self.cell_size//2,
                            text=str(self.board[row][col]),
                            font=('Arial', 12, 'bold'),
                            fill=color
                        )
                # 如果是地雷且游戏结束
                elif self.game_over and self.board[row][col] == -1:
                    self.canvas.create_rectangle(
                        x1, y1, x2, y2,
                        fill='#ff0000',
                        outline='black'
                    )
                    self.canvas.create_text(
                        x1 + self.cell_size//2,
                        y1 + self.cell_size//2,
                        text='💣',
                        font=('Arial', 14)
                    )

    def left_click(self, event):
        """处理左键点击"""
        if self.game_over:
            return

        col = event.x // self.cell_size
        row = event.y // self.cell_size

        if 0 <= row < self.rows and 0 <= col < self.cols:
            if not self.flagged[row][col] and not self.revealed[row][col]:
                # 第一次点击时放置地雷
                if self.first_click:
                    self.place_mines(row, col)
                    self.first_click = False
                    self.start_timer()

                self.reveal_cell(row, col)
                self.draw_board()
                self.check_win()

    def right_click(self, event):
        """处理右键点击"""
        if self.game_over:
            return

        col = event.x // self.cell_size
        row = event.y // self.cell_size

        if 0 <= row < self.rows and 0 <= col < self.cols:
            if not self.revealed[row][col]:
                # 切换旗帜标记
                if self.flagged[row][col]:
                    self.flagged[row][col] = False
                    self.mines_left += 1
                else:
                    self.flagged[row][col] = True
                    self.mines_left -= 1

                self.mines_label.config(text=f"💣 {self.mines_left:03d}")
                self.draw_board()

    def reveal_cell(self, row, col):
        """揭开指定格子"""
        if not (0 <= row < self.rows and 0 <= col < self.cols):
            return
        if self.revealed[row][col] or self.flagged[row][col]:
            return

        self.revealed[row][col] = True

        # 如果是地雷，游戏结束
        if self.board[row][col] == -1:
            self.game_over = True
            self.timer_running = False
            self.restart_button.config(text="😵")
            messagebox.showinfo("游戏结束", "你踩到地雷了！游戏结束！")
            return

        # 如果是空格（周围没有地雷），自动揭开周围的格子
        if self.board[row][col] == 0:
            for dr in [-1, 0, 1]:
                for dc in [-1, 0, 1]:
                    if dr == 0 and dc == 0:
                        continue
                    self.reveal_cell(row + dr, col + dc)

    def check_win(self):
        """检查是否获胜"""
        cells_to_reveal = 0
        for row in range(self.rows):
            for col in range(self.cols):
                if self.board[row][col] != -1 and not self.revealed[row][col]:
                    cells_to_reveal += 1

        if cells_to_reveal == 0:
            self.game_over = True
            self.timer_running = False
            self.restart_button.config(text="😎")

            # 自动标记所有地雷
            for row in range(self.rows):
                for col in range(self.cols):
                    if self.board[row][col] == -1:
                        self.flagged[row][col] = True

            self.mines_left = 0
            self.mines_label.config(text=f"💣 {self.mines_left:03d}")
            self.draw_board()
            messagebox.showinfo("恭喜", f"恭喜你赢了！用时{self.time_elapsed}秒！")

    def start_timer(self):
        """开始计时"""
        self.timer_running = True
        self.update_timer()

    def update_timer(self):
        """更新计时器"""
        if self.timer_running and not self.game_over:
            self.time_elapsed += 1
            self.timer_label.config(text=f"⏱ {self.time_elapsed:03d}")
            self.root.after(1000, self.update_timer)

def main():
    root = tk.Tk()
    game = Minesweeper(root, rows=9, cols=9, mines=10)
    root.mainloop()

if __name__ == "__main__":
    main()