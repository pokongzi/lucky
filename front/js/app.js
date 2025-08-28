// 应用状态管理
class LotteryApp {
    constructor() {
        this.currentGame = 'ssq'; // 默认双色球
        this.currentPage = 'home';
        this.gameConfig = {
            ssq: {
                name: '双色球',
                redBalls: 33,
                blueBalls: 16,
                redSelect: 6,
                blueSelect: 1
            },
            dlt: {
                name: '大乐透',
                redBalls: 35,
                blueBalls: 12,
                redSelect: 5,
                blueSelect: 2
            }
        };
        this.components = {};
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadComponents();
        this.updateGameUI();
    }

    // 绑定事件
    bindEvents() {
        // 游戏切换
        document.querySelectorAll('.game-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const game = e.target.dataset.game;
                this.switchGame(game);
            });
        });

        // 页面导航
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', (e) => {
                const page = e.target.closest('.nav-item').dataset.page;
                this.navigateTo(page);
            });
        });

        // 键盘导航
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey || e.metaKey) {
                switch(e.key) {
                    case '1':
                        e.preventDefault();
                        this.navigateTo('home');
                        break;
                    case '2':
                        e.preventDefault();
                        this.navigateTo('history');
                        break;
                    case '3':
                        e.preventDefault();
                        this.navigateTo('mynumbers');
                        break;
                }
            }
        });
    }

    // 切换游戏类型
    switchGame(game) {
        if (this.currentGame === game) return;
        
        this.currentGame = game;
        this.updateGameUI();
        
        // 通知所有组件游戏类型已更改
        Object.values(this.components).forEach(component => {
            if (component.onGameChange) {
                component.onGameChange(game, this.gameConfig[game]);
            }
        });
    }

    // 更新游戏UI
    updateGameUI() {
        // 更新游戏切换按钮状态
        document.querySelectorAll('.game-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.game === this.currentGame);
        });
    }

    // 页面导航
    navigateTo(page) {
        if (this.currentPage === page) return;

        // 隐藏当前页面
        const currentPageEl = document.getElementById(`${this.currentPage}-page`);
        if (currentPageEl) {
            currentPageEl.classList.remove('active');
        }

        // 更新导航状态
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.toggle('active', item.dataset.page === page);
        });

        // 显示新页面
        const newPageEl = document.getElementById(`${page}-page`);
        if (newPageEl) {
            newPageEl.classList.add('active');
        }

        this.currentPage = page;

        // 通知组件页面已切换
        if (this.components[page]) {
            this.components[page].onPageShow?.();
        }
    }

    // 加载组件
    loadComponents() {
        // 组件将在各自的文件中初始化
        // 这里只是预留接口
    }

    // 注册组件
    registerComponent(name, component) {
        this.components[name] = component;
        component.app = this;
    }

    // 获取当前游戏配置
    getCurrentGameConfig() {
        return this.gameConfig[this.currentGame];
    }

    // 工具方法：生成随机号码
    generateRandomNumbers(max, count, exclude = []) {
        const numbers = [];
        const available = [];
        
        for (let i = 1; i <= max; i++) {
            if (!exclude.includes(i)) {
                available.push(i);
            }
        }
        
        if (available.length < count) {
            throw new Error('可用号码不足');
        }
        
        while (numbers.length < count) {
            const randomIndex = Math.floor(Math.random() * available.length);
            const number = available.splice(randomIndex, 1)[0];
            numbers.push(number);
        }
        
        return numbers.sort((a, b) => a - b);
    }

    // 工具方法：格式化号码显示
    formatNumbers(redNumbers, blueNumbers) {
        const redBalls = redNumbers.map(num => 
            `<span class="number-ball red">${num.toString().padStart(2, '0')}</span>`
        ).join('');
        
        const blueBalls = blueNumbers.map(num => 
            `<span class="number-ball blue">${num.toString().padStart(2, '0')}</span>`
        ).join('');
        
        return redBalls + blueBalls;
    }

    // 工具方法：保存到本地存储
    saveToStorage(key, data) {
        try {
            localStorage.setItem(`lottery_${key}`, JSON.stringify(data));
            return true;
        } catch (e) {
            console.error('保存失败:', e);
            return false;
        }
    }

    // 工具方法：从本地存储读取
    loadFromStorage(key, defaultValue = null) {
        try {
            const data = localStorage.getItem(`lottery_${key}`);
            return data ? JSON.parse(data) : defaultValue;
        } catch (e) {
            console.error('读取失败:', e);
            return defaultValue;
        }
    }

    // 工具方法：显示通知
    showNotification(message, type = 'info', duration = 3000) {
        // 创建通知元素
        const notification = document.createElement('div');
        notification.className = `notification notification-${type}`;
        notification.innerHTML = `
            <div class="notification-content">
                <span class="notification-message">${message}</span>
                <button class="notification-close" onclick="this.parentElement.parentElement.remove()">×</button>
            </div>
        `;
        
        // 添加样式（如果还没有）
        if (!document.querySelector('#notification-styles')) {
            const style = document.createElement('style');
            style.id = 'notification-styles';
            style.textContent = `
                .notification {
                    position: fixed;
                    top: 20px;
                    right: 20px;
                    z-index: 1000;
                    min-width: 300px;
                    padding: 16px;
                    border-radius: 8px;
                    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
                    animation: slideInRight 0.3s ease;
                }
                .notification-info { background: #3182ce; color: white; }
                .notification-success { background: #38a169; color: white; }
                .notification-warning { background: #ed8936; color: white; }
                .notification-error { background: #e53e3e; color: white; }
                .notification-content {
                    display: flex;
                    align-items: center;
                    justify-content: space-between;
                    gap: 12px;
                }
                .notification-close {
                    background: rgba(255,255,255,0.2);
                    border: none;
                    color: white;
                    width: 24px;
                    height: 24px;
                    border-radius: 50%;
                    cursor: pointer;
                    font-size: 16px;
                    line-height: 1;
                }
                @keyframes slideInRight {
                    from { transform: translateX(100%); opacity: 0; }
                    to { transform: translateX(0); opacity: 1; }
                }
            `;
            document.head.appendChild(style);
        }
        
        document.body.appendChild(notification);
        
        // 自动移除
        setTimeout(() => {
            if (notification.parentElement) {
                notification.style.animation = 'slideInRight 0.3s ease reverse';
                setTimeout(() => notification.remove(), 300);
            }
        }, duration);
    }

    // 工具方法：确认对话框
    confirm(message, callback) {
        const overlay = document.createElement('div');
        overlay.className = 'confirm-overlay';
        overlay.innerHTML = `
            <div class="confirm-dialog">
                <div class="confirm-message">${message}</div>
                <div class="confirm-actions">
                    <button class="btn btn-secondary confirm-cancel">取消</button>
                    <button class="btn btn-primary confirm-ok">确定</button>
                </div>
            </div>
        `;
        
        // 添加样式
        if (!document.querySelector('#confirm-styles')) {
            const style = document.createElement('style');
            style.id = 'confirm-styles';
            style.textContent = `
                .confirm-overlay {
                    position: fixed;
                    top: 0;
                    left: 0;
                    right: 0;
                    bottom: 0;
                    background: rgba(0,0,0,0.5);
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    z-index: 1000;
                    animation: fadeIn 0.3s ease;
                }
                .confirm-dialog {
                    background: white;
                    border-radius: 12px;
                    padding: 24px;
                    min-width: 300px;
                    max-width: 90vw;
                    animation: scaleIn 0.3s ease;
                }
                .confirm-message {
                    margin-bottom: 20px;
                    font-size: 16px;
                    color: var(--text-primary);
                    text-align: center;
                }
                .confirm-actions {
                    display: flex;
                    gap: 12px;
                    justify-content: center;
                }
                @keyframes scaleIn {
                    from { transform: scale(0.8); opacity: 0; }
                    to { transform: scale(1); opacity: 1; }
                }
            `;
            document.head.appendChild(style);
        }
        
        document.body.appendChild(overlay);
        
        // 绑定事件
        overlay.querySelector('.confirm-cancel').onclick = () => {
            overlay.remove();
            callback(false);
        };
        
        overlay.querySelector('.confirm-ok').onclick = () => {
            overlay.remove();
            callback(true);
        };
        
        overlay.onclick = (e) => {
            if (e.target === overlay) {
                overlay.remove();
                callback(false);
            }
        };
    }
}

// 全局应用实例
window.lotteryApp = new LotteryApp();

// DOM 加载完成后的初始化
document.addEventListener('DOMContentLoaded', () => {
    console.log('彩票号码生成器已加载');
}); 