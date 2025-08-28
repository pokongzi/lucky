// 首页组件
class HomeComponent {
    constructor() {
        this.container = document.getElementById('home-page');
        this.quantity = 5;
        this.maxQuantity = 10;
        this.minQuantity = 1;
        this.advancedVisible = false;
        this.settings = {
            redLocked: [],
            redExcluded: [],
            blueLocked: [],
            blueExcluded: []
        };
        this.results = [];
        this.pendingList = [];
        this.init();
    }

    init() {
        this.render();
        this.bindEvents();
        this.loadPendingList();
        
        // 注册到应用
        if (window.lotteryApp) {
            window.lotteryApp.registerComponent('home', this);
        }
    }

    render() {
        const gameConfig = window.lotteryApp.getCurrentGameConfig();
        
        this.container.innerHTML = `
            <div class="home-container">
                <!-- 生成控制区 -->
                <div class="card generation-control">
                    <div class="card-body">
                        <div class="control-header">
                            <h2 class="control-title">机选注数</h2>
                            <div class="main-action-row">
                                <button class="btn btn-primary btn-lg generate-btn">
                                    🎲 一键生成幸运号码
                                </button>
                                <div class="quantity-control">
                                    <button class="quantity-btn decrease-btn">-</button>
                                    <span class="quantity-display">${this.quantity}</span>
                                    <button class="quantity-btn increase-btn">+</button>
                                </div>
                            </div>
                            <button class="advanced-toggle">
                                胆码/排除等高级设置
                                <span class="arrow">▼</span>
                            </button>
                        </div>

                        <!-- 高级设置区 -->
                        <div class="advanced-panel">
                            <div class="advanced-content">
                                <div class="section-title">
                                    🔴 红球选择 (点击号码选择状态)
                                </div>
                                <div class="number-grid red-numbers" id="red-numbers">
                                    ${this.generateNumberButtons('red', gameConfig.redBalls)}
                                </div>

                                <div class="section-title">
                                    🔵 蓝球选择 (点击号码选择状态)
                                </div>
                                <div class="number-grid blue-numbers" id="blue-numbers">
                                    ${this.generateNumberButtons('blue', gameConfig.blueBalls)}
                                </div>

                                <div class="clear-settings">
                                    <button class="btn btn-secondary btn-sm clear-btn">清空所有设置</button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 机选结果池 -->
                <div class="card results-pool">
                    <div class="card-header">
                        <h3 class="card-title">机选结果 (点击 + 选用号码)</h3>
                    </div>
                    <div class="card-body">
                        <div id="results-container">
                            <div class="empty-state">
                                <div class="icon">🎯</div>
                                <div class="message">点击上方按钮生成号码</div>
                                <div class="hint">从结果中挑选您心仪的组合</div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 分隔线 -->
                <hr style="border: none; height: 1px; background: var(--border-color); margin: 20px 0;">

                <!-- 我的待选清单 -->
                <div class="card pending-list">
                    <div class="card-header">
                        <h3 class="card-title">我的待选清单</h3>
                    </div>
                    <div class="card-body">
                        <div id="pending-container">
                            ${this.renderPendingList()}
                        </div>
                    </div>
                </div>

                <!-- 底部操作栏 -->
                <div class="final-actions ${this.pendingList.length === 0 ? 'disabled' : ''}">
                    <div class="action-buttons">
                        <button class="btn btn-secondary copy-btn">📋 一键复制清单</button>
                        <button class="btn btn-primary save-btn">⭐ 一键收藏清单</button>
                        <button class="btn btn-outline clear-pending-btn">🗑️ 清空清单</button>
                    </div>
                </div>
            </div>
        `;
    }

    generateNumberButtons(type, max) {
        let html = '';
        for (let i = 1; i <= max; i++) {
            html += `<button class="number-btn ${type}" data-number="${i}">${i.toString().padStart(2, '0')}</button>`;
        }
        return html;
    }

    bindEvents() {
        // 生成按钮
        this.container.querySelector('.generate-btn').addEventListener('click', () => {
            this.generateNumbers();
        });

        // 数量控制
        this.container.querySelector('.decrease-btn').addEventListener('click', () => {
            this.changeQuantity(-1);
        });

        this.container.querySelector('.increase-btn').addEventListener('click', () => {
            this.changeQuantity(1);
        });

        // 高级设置切换
        this.container.querySelector('.advanced-toggle').addEventListener('click', () => {
            this.toggleAdvanced();
        });

        // 号码按钮点击
        this.container.querySelectorAll('.number-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.toggleNumberState(e.target);
            });
        });

        // 清空设置
        this.container.querySelector('.clear-btn').addEventListener('click', () => {
            this.clearSettings();
        });

        // 底部操作按钮
        this.container.querySelector('.copy-btn').addEventListener('click', () => {
            this.copyPendingList();
        });

        this.container.querySelector('.save-btn').addEventListener('click', () => {
            this.savePendingList();
        });

        this.container.querySelector('.clear-pending-btn').addEventListener('click', () => {
            this.clearPendingList();
        });
    }

    changeQuantity(delta) {
        const newQuantity = this.quantity + delta;
        if (newQuantity >= this.minQuantity && newQuantity <= this.maxQuantity) {
            this.quantity = newQuantity;
            this.container.querySelector('.quantity-display').textContent = this.quantity;
            
            // 更新按钮状态
            const decreaseBtn = this.container.querySelector('.decrease-btn');
            const increaseBtn = this.container.querySelector('.increase-btn');
            
            decreaseBtn.disabled = this.quantity <= this.minQuantity;
            increaseBtn.disabled = this.quantity >= this.maxQuantity;
        }
    }

    toggleAdvanced() {
        this.advancedVisible = !this.advancedVisible;
        const panel = this.container.querySelector('.advanced-panel');
        const arrow = this.container.querySelector('.arrow');
        const toggle = this.container.querySelector('.advanced-toggle');
        
        if (this.advancedVisible) {
            panel.classList.add('show');
            toggle.classList.add('expanded');
        } else {
            panel.classList.remove('show');
            toggle.classList.remove('expanded');
        }
    }

    toggleNumberState(btn) {
        const number = parseInt(btn.dataset.number);
        const type = btn.classList.contains('red') ? 'red' : 'blue';
        const lockedKey = `${type}Locked`;
        const excludedKey = `${type}Excluded`;
        
        // 当前状态：普通 -> 锁定 -> 排除 -> 普通
        if (btn.classList.contains('locked')) {
            // 锁定 -> 排除
            btn.classList.remove('locked');
            btn.classList.add('excluded');
            this.settings[lockedKey] = this.settings[lockedKey].filter(n => n !== number);
            this.settings[excludedKey].push(number);
        } else if (btn.classList.contains('excluded')) {
            // 排除 -> 普通
            btn.classList.remove('excluded');
            this.settings[excludedKey] = this.settings[excludedKey].filter(n => n !== number);
        } else {
            // 普通 -> 锁定
            btn.classList.add('locked');
            this.settings[lockedKey].push(number);
        }
        
        console.log('当前设置:', this.settings);
    }

    clearSettings() {
        this.settings = {
            redLocked: [],
            redExcluded: [],
            blueLocked: [],
            blueExcluded: []
        };
        
        // 清除所有按钮状态
        this.container.querySelectorAll('.number-btn').forEach(btn => {
            btn.classList.remove('locked', 'excluded');
        });
        
        window.lotteryApp.showNotification('已清空所有设置', 'success');
    }

    generateNumbers() {
        try {
            const gameConfig = window.lotteryApp.getCurrentGameConfig();
            this.results = [];
            
            for (let i = 0; i < this.quantity; i++) {
                const result = this.generateSingleNumber(gameConfig);
                this.results.push(result);
            }
            
            this.renderResults();
            window.lotteryApp.showNotification(`成功生成 ${this.quantity} 注号码`, 'success');
        } catch (error) {
            window.lotteryApp.showNotification(error.message, 'error');
        }
    }

    generateSingleNumber(gameConfig) {
        // 生成红球
        let redNumbers;
        if (this.settings.redLocked.length > 0) {
            // 有胆码，先加入胆码
            redNumbers = [...this.settings.redLocked];
            const remainingCount = gameConfig.redSelect - redNumbers.length;
            
            if (remainingCount > 0) {
                const excludeList = [...this.settings.redLocked, ...this.settings.redExcluded];
                const additionalNumbers = window.lotteryApp.generateRandomNumbers(
                    gameConfig.redBalls, 
                    remainingCount, 
                    excludeList
                );
                redNumbers.push(...additionalNumbers);
            }
            redNumbers.sort((a, b) => a - b);
        } else {
            // 无胆码，正常生成
            redNumbers = window.lotteryApp.generateRandomNumbers(
                gameConfig.redBalls, 
                gameConfig.redSelect, 
                this.settings.redExcluded
            );
        }
        
        // 生成蓝球
        let blueNumbers;
        if (this.settings.blueLocked.length > 0) {
            blueNumbers = [...this.settings.blueLocked];
            const remainingCount = gameConfig.blueSelect - blueNumbers.length;
            
            if (remainingCount > 0) {
                const excludeList = [...this.settings.blueLocked, ...this.settings.blueExcluded];
                const additionalNumbers = window.lotteryApp.generateRandomNumbers(
                    gameConfig.blueBalls, 
                    remainingCount, 
                    excludeList
                );
                blueNumbers.push(...additionalNumbers);
            }
            blueNumbers.sort((a, b) => a - b);
        } else {
            blueNumbers = window.lotteryApp.generateRandomNumbers(
                gameConfig.blueBalls, 
                gameConfig.blueSelect, 
                this.settings.blueExcluded
            );
        }
        
        return {
            id: Date.now() + Math.random(),
            red: redNumbers,
            blue: blueNumbers,
            game: window.lotteryApp.currentGame,
            generated: new Date()
        };
    }

    renderResults() {
        const container = this.container.querySelector('#results-container');
        
        if (this.results.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <div class="icon">🎯</div>
                    <div class="message">点击上方按钮生成号码</div>
                    <div class="hint">从结果中挑选您心仪的组合</div>
                </div>
            `;
            return;
        }
        
        const html = this.results.map(result => `
            <div class="result-item" data-id="${result.id}">
                <div class="result-numbers">
                    ${window.lotteryApp.formatNumbers(result.red, result.blue)}
                </div>
                <button class="select-btn" title="选用这注号码">+</button>
            </div>
        `).join('');
        
        container.innerHTML = html;
        
        // 绑定选用按钮事件
        container.querySelectorAll('.select-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const resultItem = e.target.closest('.result-item');
                const resultId = resultItem.dataset.id;
                this.selectResult(resultId);
            });
        });
    }

    selectResult(resultId) {
        const result = this.results.find(r => r.id == resultId);
        if (!result) return;
        
        // 添加到待选清单
        const pendingItem = {
            ...result,
            id: Date.now() + Math.random(),
            selectedAt: new Date()
        };
        
        this.pendingList.push(pendingItem);
        this.savePendingListToStorage();
        this.renderPendingList();
        this.updateFinalActions();
        
        // 标记结果项为已选择
        const resultItem = this.container.querySelector(`[data-id="${resultId}"]`);
        if (resultItem) {
            resultItem.classList.add('selected');
            const btn = resultItem.querySelector('.select-btn');
            btn.disabled = true;
            btn.textContent = '已选用';
        }
        
        window.lotteryApp.showNotification('已添加到待选清单', 'success');
    }

    renderPendingList() {
        const container = this.container.querySelector('#pending-container');
        
        if (this.pendingList.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <div class="icon">📝</div>
                    <div class="message">请从上方结果中选用号码</div>
                    <div class="hint">选用的号码将在这里显示</div>
                </div>
            `;
            return;
        }
        
        const html = this.pendingList.map(item => `
            <div class="pending-item" data-id="${item.id}">
                <div class="pending-info">
                    <div class="pending-numbers">
                        ${window.lotteryApp.formatNumbers(item.red, item.blue)}
                    </div>
                    <div class="pending-date">
                        收藏于: ${new Date(item.selectedAt).toLocaleString()}
                    </div>
                </div>
                <button class="remove-btn" title="移除这注号码">移除</button>
            </div>
        `).join('');
        
        container.innerHTML = html;
        
        // 绑定移除按钮事件
        container.querySelectorAll('.remove-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const pendingItem = e.target.closest('.pending-item');
                const itemId = pendingItem.dataset.id;
                this.removePendingItem(itemId);
            });
        });
        
        return html;
    }

    removePendingItem(itemId) {
        this.pendingList = this.pendingList.filter(item => item.id != itemId);
        this.savePendingListToStorage();
        this.renderPendingList();
        this.updateFinalActions();
        
        window.lotteryApp.showNotification('已移除号码', 'info');
    }

    updateFinalActions() {
        const finalActions = this.container.querySelector('.final-actions');
        if (this.pendingList.length === 0) {
            finalActions.classList.add('disabled');
        } else {
            finalActions.classList.remove('disabled');
        }
    }

    copyPendingList() {
        if (this.pendingList.length === 0) return;
        
        const gameConfig = window.lotteryApp.getCurrentGameConfig();
        let text = `${gameConfig.name} 机选号码 (${this.pendingList.length}注)\n`;
        text += '━━━━━━━━━━━━━━━━━━━━━━━━━━\n';
        
        this.pendingList.forEach((item, index) => {
            const redStr = item.red.map(n => n.toString().padStart(2, '0')).join(' ');
            const blueStr = item.blue.map(n => n.toString().padStart(2, '0')).join(' ');
            text += `${(index + 1).toString().padStart(2, '0')}. 🔴 ${redStr} 🔵 ${blueStr}\n`;
        });
        
        text += '━━━━━━━━━━━━━━━━━━━━━━━━━━\n';
        text += `生成时间: ${new Date().toLocaleString()}\n`;
        text += '祝您好运！';
        
        navigator.clipboard.writeText(text).then(() => {
            window.lotteryApp.showNotification('已复制到剪贴板', 'success');
        }).catch(() => {
            window.lotteryApp.showNotification('复制失败，请手动复制', 'error');
        });
    }

    savePendingList() {
        if (this.pendingList.length === 0) return;
        
        // 保存到我的号码
        const savedNumbers = window.lotteryApp.loadFromStorage('myNumbers', []);
        
        this.pendingList.forEach(item => {
            const savedItem = {
                ...item,
                id: Date.now() + Math.random(),
                savedAt: new Date(),
                game: item.game
            };
            savedNumbers.push(savedItem);
        });
        
        window.lotteryApp.saveToStorage('myNumbers', savedNumbers);
        
        // 清空待选清单
        this.clearPendingList();
        
        window.lotteryApp.showNotification(`已收藏 ${this.pendingList.length} 注号码`, 'success');
    }

    clearPendingList() {
        window.lotteryApp.confirm('确定要清空待选清单吗？', (confirmed) => {
            if (confirmed) {
                this.pendingList = [];
                this.savePendingListToStorage();
                this.renderPendingList();
                this.updateFinalActions();
                window.lotteryApp.showNotification('已清空待选清单', 'info');
            }
        });
    }

    savePendingListToStorage() {
        window.lotteryApp.saveToStorage('pendingList', this.pendingList);
    }

    loadPendingList() {
        this.pendingList = window.lotteryApp.loadFromStorage('pendingList', []);
        this.updateFinalActions();
    }

    // 游戏类型切换时的处理
    onGameChange(game, gameConfig) {
        console.log('游戏切换到:', gameConfig.name);
        this.clearSettings();
        this.results = [];
        this.render();
        this.bindEvents();
        this.loadPendingList();
    }

    // 页面显示时的处理
    onPageShow() {
        console.log('首页显示');
    }
}

// 初始化首页组件
document.addEventListener('DOMContentLoaded', () => {
    new HomeComponent();
}); 