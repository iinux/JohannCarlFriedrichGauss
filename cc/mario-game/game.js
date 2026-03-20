// 超级玛丽游戏 - 主游戏引擎

const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');

// 游戏常量
const GRAVITY = 0.5;
const FRICTION = 0.8;
const TILE_SIZE = 40;

// 游戏状态
let gameState = {
    score: 0,
    lives: 3,
    level: 1,
    gameOver: false,
    victory: false
};

// 按键状态
const keys = {
    left: false,
    right: false,
    up: false
};

// 玩家对象
let player = {
    x: 100,
    y: 100,
    width: 30,
    height: 40,
    velX: 0,
    velY: 0,
    speed: 5,
    jumpStrength: 12,
    grounded: false,
    color: '#ff0000',
    facingRight: true,
    invincible: false,
    invincibleTimer: 0
};

// 游戏对象数组
let platforms = [];
let enemies = [];
let coins = [];
let particles = [];
let goal = null;

// 关卡设计
const levels = [
    // 第 1 关
    {
        platforms: [
            // 地面
            { x: 0, y: 560, width: 800, height: 40, type: 'ground' },
            // 平台
            { x: 200, y: 450, width: 120, height: 20, type: 'brick' },
            { x: 400, y: 380, width: 120, height: 20, type: 'brick' },
            { x: 600, y: 300, width: 120, height: 20, type: 'brick' },
            { x: 50, y: 350, width: 100, height: 20, type: 'brick' },
            // 台阶
            { x: 720, y: 480, width: 40, height: 80, type: 'block' },
            { x: 760, y: 520, width: 40, height: 40, type: 'block' }
        ],
        enemies: [
            { x: 300, y: 520, width: 30, height: 30, speed: 2, range: 100 },
            { x: 500, y: 520, width: 30, height: 30, speed: 2.5, range: 150 }
        ],
        coins: [
            { x: 240, y: 410 },
            { x: 280, y: 410 },
            { x: 440, y: 340 },
            { x: 480, y: 340 },
            { x: 640, y: 260 },
            { x: 680, y: 260 },
            { x: 90, y: 310 },
            { x: 130, y: 310 }
        ],
        goal: { x: 750, y: 480, width: 40, height: 80 },
        playerStart: { x: 50, y: 500 }
    },
    // 第 2 关
    {
        platforms: [
            { x: 0, y: 560, width: 300, height: 40, type: 'ground' },
            { x: 350, y: 560, width: 200, height: 40, type: 'ground' },
            { x: 600, y: 560, width: 200, height: 40, type: 'ground' },
            { x: 150, y: 450, width: 100, height: 20, type: 'brick' },
            { x: 300, y: 380, width: 100, height: 20, type: 'brick' },
            { x: 450, y: 300, width: 100, height: 20, type: 'brick' },
            { x: 600, y: 400, width: 100, height: 20, type: 'brick' },
            { x: 50, y: 280, width: 150, height: 20, type: 'brick' }
        ],
        enemies: [
            { x: 400, y: 520, width: 30, height: 30, speed: 3, range: 120 },
            { x: 650, y: 520, width: 30, height: 30, speed: 2, range: 100 },
            { x: 200, y: 410, width: 30, height: 30, speed: 2.5, range: 80 }
        ],
        coins: [
            { x: 180, y: 410 },
            { x: 220, y: 410 },
            { x: 330, y: 340 },
            { x: 370, y: 340 },
            { x: 480, y: 260 },
            { x: 520, y: 260 },
            { x: 630, y: 360 },
            { x: 670, y: 360 },
            { x: 100, y: 240 },
            { x: 140, y: 240 }
        ],
        goal: { x: 750, y: 480, width: 40, height: 80 },
        playerStart: { x: 50, y: 500 }
    },
    // 第 3 关
    {
        platforms: [
            { x: 0, y: 560, width: 200, height: 40, type: 'ground' },
            { x: 250, y: 500, width: 80, height: 20, type: 'brick' },
            { x: 380, y: 450, width: 80, height: 20, type: 'brick' },
            { x: 510, y: 400, width: 80, height: 20, type: 'brick' },
            { x: 640, y: 350, width: 80, height: 20, type: 'brick' },
            { x: 720, y: 560, width: 80, height: 40, type: 'ground' },
            { x: 100, y: 350, width: 100, height: 20, type: 'brick' },
            { x: 300, y: 280, width: 100, height: 20, type: 'brick' },
            { x: 500, y: 220, width: 150, height: 20, type: 'brick' }
        ],
        enemies: [
            { x: 270, y: 460, width: 30, height: 30, speed: 2, range: 60 },
            { x: 400, y: 410, width: 30, height: 30, speed: 2.5, range: 60 },
            { x: 530, y: 360, width: 30, height: 30, speed: 3, range: 60 },
            { x: 660, y: 310, width: 30, height: 30, speed: 3.5, range: 60 }
        ],
        coins: [
            { x: 280, y: 460 },
            { x: 410, y: 410 },
            { x: 540, y: 360 },
            { x: 670, y: 310 },
            { x: 130, y: 310 },
            { x: 330, y: 240 },
            { x: 530, y: 180 },
            { x: 570, y: 180 },
            { x: 610, y: 180 }
        ],
        goal: { x: 750, y: 480, width: 40, height: 80 },
        playerStart: { x: 50, y: 500 }
    }
];

// 初始化关卡
function initLevel(levelNum) {
    const levelIndex = (levelNum - 1) % levels.length;
    const level = levels[levelIndex];
    
    platforms = level.platforms.map(p => ({ ...p }));
    enemies = level.enemies.map(e => ({ 
        ...e, 
        startX: e.x, 
        direction: 1,
        color: '#8b4513'
    }));
    coins = level.coins.map(c => ({ ...c, width: 20, height: 20, collected: false }));
    goal = { ...level.goal };
    
    player.x = level.playerStart.x;
    player.y = level.playerStart.y;
    player.velX = 0;
    player.velY = 0;
}

// 粒子效果
function createParticles(x, y, color, count = 10) {
    for (let i = 0; i < count; i++) {
        particles.push({
            x: x,
            y: y,
            velX: (Math.random() - 0.5) * 10,
            velY: (Math.random() - 0.5) * 10,
            life: 1,
            color: color,
            size: Math.random() * 5 + 2
        });
    }
}

// 更新粒子
function updateParticles() {
    for (let i = particles.length - 1; i >= 0; i--) {
        const p = particles[i];
        p.x += p.velX;
        p.y += p.velY;
        p.velY += GRAVITY;
        p.life -= 0.02;
        
        if (p.life <= 0) {
            particles.splice(i, 1);
        }
    }
}

// 绘制粒子
function drawParticles() {
    particles.forEach(p => {
        ctx.globalAlpha = p.life;
        ctx.fillStyle = p.color;
        ctx.beginPath();
        ctx.arc(p.x, p.y, p.size, 0, Math.PI * 2);
        ctx.fill();
    });
    ctx.globalAlpha = 1;
}

// 更新玩家
function updatePlayer() {
    // 水平移动
    if (keys.left) {
        player.velX = -player.speed;
        player.facingRight = false;
    }
    if (keys.right) {
        player.velX = player.speed;
        player.facingRight = true;
    }
    
    // 跳跃
    if (keys.up && player.grounded) {
        player.velY = -player.jumpStrength;
        player.grounded = false;
    }
    
    // 应用物理
    player.velX *= FRICTION;
    player.velY += GRAVITY;
    
    // 限制速度
    player.velX = Math.max(-player.speed, Math.min(player.speed, player.velX));
    
    // 更新位置
    player.x += player.velX;
    player.y += player.velY;
    
    // 边界检查
    if (player.x < 0) player.x = 0;
    if (player.x + player.width > canvas.width) player.x = canvas.width - player.width;
    
    // 掉落死亡
    if (player.y > canvas.height) {
        playerDeath();
    }
    
    // 无敌时间
    if (player.invincible) {
        player.invincibleTimer--;
        if (player.invincibleTimer <= 0) {
            player.invincible = false;
        }
    }
}

// 碰撞检测
function checkCollision(rect1, rect2) {
    return rect1.x < rect2.x + rect2.width &&
           rect1.x + rect1.width > rect2.x &&
           rect1.y < rect2.y + rect2.height &&
           rect1.y + rect1.height > rect2.y;
}

// 处理平台碰撞
function handlePlatformCollisions() {
    player.grounded = false;
    
    platforms.forEach(platform => {
        // 简单的 AABB 碰撞检测
        const collision = checkCollision(player, platform);
        
        if (collision) {
            // 计算重叠量
            const overlapX = (player.width + platform.width) / 2 - Math.abs((player.x + player.width / 2) - (platform.x + platform.width / 2));
            const overlapY = (player.height + platform.height) / 2 - Math.abs((player.y + player.height / 2) - (platform.y + platform.height / 2));
            
            if (overlapX < overlapY) {
                // 水平碰撞
                if (player.x < platform.x) {
                    player.x = platform.x - player.width;
                } else {
                    player.x = platform.x + platform.width;
                }
                player.velX = 0;
            } else {
                // 垂直碰撞
                if (player.y < platform.y) {
                    player.y = platform.y - player.height;
                    player.velY = 0;
                    player.grounded = true;
                } else {
                    player.y = platform.y + platform.height;
                    player.velY = 0;
                }
            }
        }
    });
}

// 更新敌人
function updateEnemies() {
    enemies.forEach(enemy => {
        // 来回移动
        enemy.x += enemy.speed * enemy.direction;
        
        // 检查是否到达范围端点
        if (enemy.x > enemy.startX + enemy.range || enemy.x < enemy.startX - enemy.range) {
            enemy.direction *= -1;
        }
        
        // 与玩家碰撞
        if (checkCollision(player, enemy)) {
            // 如果玩家从上方落下，消灭敌人
            if (player.velY > 0 && player.y + player.height - player.velY < enemy.y + enemy.height / 2) {
                enemy.x = -1000; // 移除敌人
                player.velY = -player.jumpStrength / 2; // 弹起
                gameState.score += 100;
                createParticles(enemy.x + enemy.width / 2, enemy.y + enemy.height / 2, '#8b4513', 15);
                updateUI();
            } else if (!player.invincible) {
                // 玩家受伤
                playerDeath();
            }
        }
    });
}

// 玩家死亡
function playerDeath() {
    gameState.lives--;
    updateUI();
    
    if (gameState.lives <= 0) {
        gameState.gameOver = true;
    } else {
        // 重生
        player.invincible = true;
        player.invincibleTimer = 120; // 2 秒无敌时间
        const level = levels[(gameState.level - 1) % levels.length];
        player.x = level.playerStart.x;
        player.y = level.playerStart.y;
        player.velX = 0;
        player.velY = 0;
    }
}

// 更新金币
function updateCoins() {
    coins.forEach(coin => {
        if (!coin.collected && checkCollision(player, coin)) {
            coin.collected = true;
            gameState.score += 50;
            createParticles(coin.x + coin.width / 2, coin.y + coin.height / 2, '#ffd700', 8);
            updateUI();
        }
    });
}

// 检查胜利条件
function checkVictory() {
    if (goal && checkCollision(player, goal)) {
        gameState.level++;
        if (gameState.level > levels.length) {
            gameState.victory = true;
        } else {
            initLevel(gameState.level);
        }
        updateUI();
    }
}

// 更新 UI
function updateUI() {
    document.getElementById('score').textContent = `得分：${gameState.score}`;
    document.getElementById('lives').textContent = `生命：${gameState.lives}`;
    document.getElementById('level').textContent = `关卡：${gameState.level}`;
}

// 绘制平台
function drawPlatforms() {
    platforms.forEach(platform => {
        switch(platform.type) {
            case 'ground':
                ctx.fillStyle = '#8b4513';
                break;
            case 'brick':
                ctx.fillStyle = '#cd853f';
                break;
            case 'block':
                ctx.fillStyle = '#daa520';
                break;
            default:
                ctx.fillStyle = '#8b4513';
        }
        
        ctx.fillRect(platform.x, platform.y, platform.width, platform.height);
        
        // 添加边框
        ctx.strokeStyle = '#000';
        ctx.lineWidth = 2;
        ctx.strokeRect(platform.x, platform.y, platform.width, platform.height);
    });
}

// 绘制玩家
function drawPlayer() {
    if (player.invincible && Math.floor(Date.now() / 100) % 2 === 0) {
        return; // 闪烁效果
    }
    
    // 身体
    ctx.fillStyle = player.color;
    ctx.fillRect(player.x, player.y, player.width, player.height);
    
    // 帽子
    ctx.fillStyle = '#cc0000';
    ctx.fillRect(player.x - 5, player.y, player.width + 10, 10);
    
    // 脸
    ctx.fillStyle = '#ffcc99';
    const faceX = player.facingRight ? player.x + 15 : player.x + 5;
    ctx.fillRect(faceX, player.y + 10, 15, 15);
    
    // 眼睛
    ctx.fillStyle = '#000';
    const eyeX = player.facingRight ? faceX + 8 : faceX + 2;
    ctx.fillRect(eyeX, player.y + 14, 4, 4);
    
    // 背带裤
    ctx.fillStyle = '#0000cc';
    ctx.fillRect(player.x + 5, player.y + 25, player.width - 10, 15);
}

// 绘制敌人
function drawEnemies() {
    enemies.forEach(enemy => {
        if (enemy.x < -100) return; // 已死亡的敌人不绘制
        
        // 身体
        ctx.fillStyle = enemy.color;
        ctx.fillRect(enemy.x, enemy.y, enemy.width, enemy.height);
        
        // 眼睛
        ctx.fillStyle = '#fff';
        ctx.fillRect(enemy.x + 5, enemy.y + 8, 8, 8);
        ctx.fillRect(enemy.x + 17, enemy.y + 8, 8, 8);
        
        ctx.fillStyle = '#000';
        ctx.fillRect(enemy.x + 7, enemy.y + 10, 4, 4);
        ctx.fillRect(enemy.x + 19, enemy.y + 10, 4, 4);
        
        // 脚
        ctx.fillStyle = '#000';
        const footOffset = Math.sin(Date.now() / 100) * 3;
        ctx.fillRect(enemy.x + 3 + footOffset, enemy.y + enemy.height - 5, 10, 5);
        ctx.fillRect(enemy.x + 17 - footOffset, enemy.y + enemy.height - 5, 10, 5);
    });
}

// 绘制金币
function drawCoins() {
    coins.forEach(coin => {
        if (coin.collected) return;
        
        // 金币主体
        ctx.fillStyle = '#ffd700';
        ctx.beginPath();
        ctx.arc(coin.x + coin.width / 2, coin.y + coin.height / 2, 10, 0, Math.PI * 2);
        ctx.fill();
        
        // 高光
        ctx.fillStyle = '#ffff00';
        ctx.beginPath();
        ctx.arc(coin.x + coin.width / 2 - 3, coin.y + coin.height / 2 - 3, 4, 0, Math.PI * 2);
        ctx.fill();
        
        // 边框
        ctx.strokeStyle = '#daa520';
        ctx.lineWidth = 2;
        ctx.stroke();
    });
}

// 绘制终点
function drawGoal() {
    if (!goal) return;
    
    // 传送门/旗杆
    ctx.fillStyle = '#00ff00';
    ctx.globalAlpha = 0.6 + Math.sin(Date.now() / 200) * 0.2;
    ctx.fillRect(goal.x, goal.y, goal.width, goal.height);
    
    // 闪光效果
    ctx.fillStyle = '#ffffff';
    ctx.globalAlpha = 0.3;
    ctx.fillRect(goal.x + 10, goal.y + 10, goal.width - 20, goal.height - 20);
    
    ctx.globalAlpha = 1;
}

// 绘制游戏结束画面
function drawGameOver() {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    ctx.fillStyle = '#ff0000';
    ctx.font = 'bold 48px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('游戏结束', canvas.width / 2, canvas.height / 2 - 30);
    
    ctx.fillStyle = '#ffffff';
    ctx.font = '24px Arial';
    ctx.fillText(`最终得分：${gameState.score}`, canvas.width / 2, canvas.height / 2 + 20);
    ctx.fillText('按 R 重新开始', canvas.width / 2, canvas.height / 2 + 60);
}

// 绘制胜利画面
function drawVictory() {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    ctx.fillStyle = '#ffd700';
    ctx.font = 'bold 48px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('恭喜通关！', canvas.width / 2, canvas.height / 2 - 30);
    
    ctx.fillStyle = '#ffffff';
    ctx.font = '24px Arial';
    ctx.fillText(`最终得分：${gameState.score}`, canvas.width / 2, canvas.height / 2 + 20);
    ctx.fillText('🎉 你完成了所有关卡！ 🎉', canvas.width / 2, canvas.height / 2 + 60);
    ctx.fillText('按 R 重新开始', canvas.width / 2, canvas.height / 2 + 100);
}

// 绘制背景
function drawBackground() {
    // 天空渐变
    const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
    gradient.addColorStop(0, '#5c94fc');
    gradient.addColorStop(1, '#87ceeb');
    ctx.fillStyle = gradient;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    // 云朵
    ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
    drawCloud(100, 80, 60);
    drawCloud(300, 120, 80);
    drawCloud(550, 60, 70);
    drawCloud(700, 100, 50);
    
    // 山丘
    ctx.fillStyle = '#228b22';
    ctx.beginPath();
    ctx.moveTo(0, 560);
    ctx.quadraticCurveTo(200, 450, 400, 560);
    ctx.fill();
    
    ctx.beginPath();
    ctx.moveTo(300, 560);
    ctx.quadraticCurveTo(500, 480, 700, 560);
    ctx.fill();
}

// 绘制云朵
function drawCloud(x, y, size) {
    ctx.beginPath();
    ctx.arc(x, y, size / 3, 0, Math.PI * 2);
    ctx.arc(x + size / 3, y - size / 6, size / 4, 0, Math.PI * 2);
    ctx.arc(x + size / 2, y, size / 3, 0, Math.PI * 2);
    ctx.arc(x + size / 4, y + size / 6, size / 5, 0, Math.PI * 2);
    ctx.fill();
}

// 游戏主循环
function gameLoop() {
    // 清空画布
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    if (!gameState.gameOver && !gameState.victory) {
        // 更新游戏逻辑
        updatePlayer();
        handlePlatformCollisions();
        updateEnemies();
        updateCoins();
        updateParticles();
        checkVictory();
    }
    
    // 绘制
    drawBackground();
    drawPlatforms();
    drawCoins();
    drawGoal();
    drawEnemies();
    drawPlayer();
    drawParticles();
    
    if (gameState.gameOver) {
        drawGameOver();
    }
    
    if (gameState.victory) {
        drawVictory();
    }
    
    requestAnimationFrame(gameLoop);
}

// 键盘事件监听
document.addEventListener('keydown', (e) => {
    switch(e.key) {
        case 'ArrowLeft':
            keys.left = true;
            break;
        case 'ArrowRight':
            keys.right = true;
            break;
        case 'ArrowUp':
            keys.up = true;
            e.preventDefault();
            break;
        case 'r':
        case 'R':
            resetGame();
            break;
    }
});

document.addEventListener('keyup', (e) => {
    switch(e.key) {
        case 'ArrowLeft':
            keys.left = false;
            break;
        case 'ArrowRight':
            keys.right = false;
            break;
        case 'ArrowUp':
            keys.up = false;
            break;
    }
});

// 重置游戏
function resetGame() {
    gameState = {
        score: 0,
        lives: 3,
        level: 1,
        gameOver: false,
        victory: false
    };
    initLevel(1);
    updateUI();
}

// 初始化游戏
initLevel(1);
updateUI();
gameLoop();

console.log('🍄 超级玛丽游戏已启动！使用方向键控制，R 键重新开始。');
