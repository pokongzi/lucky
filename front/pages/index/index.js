Page({
  data: {
    currentGame: 'ssq',
    numbers: []
  },
  
  onLoad: function() {
    this.generateNumbers();
  },
  
  switchGame: function(e) {
    const game = e.currentTarget.dataset.game;
    this.setData({
      currentGame: game
    });
    this.generateNumbers();
  },
  
  generateNumbers: function() {
    let numbers = [];
    if (this.data.currentGame === 'ssq') {
      // 双色球：6个红球(1-33) + 1个蓝球(1-16)
      const redBalls = this.getRandomNumbers(1, 33, 6);
      const blueBall = this.getRandomNumbers(1, 16, 1);
      numbers = [...redBalls, ...blueBall];
    } else {
      // 大乐透：5个前区(1-35) + 2个后区(1-12)
      const frontBalls = this.getRandomNumbers(1, 35, 5);
      const backBalls = this.getRandomNumbers(1, 12, 2);
      numbers = [...frontBalls, ...backBalls];
    }
    
    this.setData({
      numbers: numbers
    });
  },
  
  getRandomNumbers: function(min, max, count) {
    const numbers = [];
    while (numbers.length < count) {
      const num = Math.floor(Math.random() * (max - min + 1)) + min;
      if (!numbers.includes(num)) {
        numbers.push(num);
      }
    }
    return numbers.sort((a, b) => a - b);
  }
});
