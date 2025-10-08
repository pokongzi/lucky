Page({
  data: {
    currentGame: 'ssq',
    currentTab: 'history',
    periodFilter: 30,
    historyList: [],
    redHeatmap: [],
    blueHeatmap: [],
    maxFrequency: 1,
    redMissing: [],
    blueMissing: [],
    // 分页相关状态
    pageNum: 1,
    pageSize: 10,
    hasMoreData: true
  },
  
  onLoad: function() {
    this.loadHistory();
    this.loadHeatmap();
    this.loadMissing();
  },
  
  // 切换游戏
  switchGame: function(e) {
    const game = e.currentTarget.dataset.game;
    this.setData({
      currentGame: game,
      pageNum: 1,
      hasMoreData: true,
      historyList: []
    });
    this.loadHistory();
    this.loadHeatmap();
    this.loadMissing();
  },
  
  // 切换Tab
  switchTab: function(e) {
    const tab = e.currentTarget.dataset.tab;
    this.setData({
      currentTab: tab
    });
  },
  
  // 设置期数筛选
  setPeriodFilter: function(e) {
    const period = parseInt(e.currentTarget.dataset.period);
    this.setData({
      periodFilter: period
    });
    this.loadHeatmap();
  },
  
  // 加载历史数据
  loadHistory: function() {
    // 如果没有更多数据，直接返回
    if (!this.data.hasMoreData) {
      return;
    }
    
    const app = getApp();
    const baseUrl = app.getBaseURL();
    const gameCode = this.data.currentGame;
    const pageNum = this.data.pageNum;
    const pageSize = this.data.pageSize;
    
    console.log('开始加载历史数据，游戏代码:', gameCode, '页码:', pageNum, '页大小:', pageSize, 'API URL:', `${baseUrl}/api/results/${gameCode}?page=${pageNum}&pageSize=${pageSize}`);
    
    wx.showLoading({ title: '加载中' });
    
    wx.request({
      url: `${baseUrl}/api/results/${gameCode}?page=${pageNum}&pageSize=${pageSize}`,
      method: 'GET',
      success: (res) => {
        try {
          console.log('API响应状态码:', res.statusCode);
          console.log('API响应数据:', res.data);
          
          // 确保响应是有效的JSON
          const responseData = res.data;
          
          if (res.statusCode === 200 && responseData) {
            console.log('响应数据结构检查:', {
              'code是否存在': 'code' in responseData,
              'code值': responseData.code,
              'data是否存在': 'data' in responseData,
              'data.list是否存在': responseData.data && 'list' in responseData.data,
              'list长度': responseData.data && responseData.data.list ? responseData.data.list.length : 0
            });
            
            if (responseData.data && responseData.data.list && responseData.data.list.length > 0) {
              // 处理API返回的数据，转换为页面需要的格式，根据游戏类型保留合理格式的期号数据
              const newHistoryList = responseData.data.list
                .filter(item => {
                  if (!item.period) return false;
                  const periodStr = item.period.toString();
                  // 大乐透(dlt)接受5位数字期号，双色球(ssq)接受6位及以上数字期号
                  if (gameCode === 'dlt') {
                    return /^\d{5}$/.test(periodStr);
                  } else {
                    return /^\d{6,}$/.test(periodStr);
                  }
                })
                .map(item => ({
                  id: item.id,
                  period: item.period,
                  redNumbers: item.red_balls || [],
                  blueNumbers: item.blue_balls || [],
                  date: item.draw_date ? item.draw_date.split('T')[0] : '', // 提取年月日
                  salesAmount: item.sales_amount || '',
                  prizePool: item.prize_pool || '',
                  showPrize: false,
                  prizes: item.prizes || []
                }));
              
              console.log('转换后的历史数据:', newHistoryList);
              
              // 判断是否为第一页还是加载更多
              if (this.data.pageNum === 1) {
                this.setData({
                  historyList: newHistoryList
                });
              } else {
                // 追加数据
                this.setData({
                  historyList: [...this.data.historyList, ...newHistoryList]
                });
              }
              
              console.log('数据已设置到页面，historyList长度:', this.data.historyList.length);
              
              // 如果返回的数据少于pageSize，说明没有更多数据了
              if (newHistoryList.length < this.data.pageSize) {
                this.setData({
                  hasMoreData: false
                });
              }
            } else if (responseData.data && responseData.data.list) {
              console.log('当前游戏(' + gameCode + ')暂无历史数据，使用模拟数据');
              // 当没有真实数据时，生成一些模拟数据
              const mockData = this.generateMockHistoryData(gameCode);
              
              // 判断是否为第一页还是加载更多
              if (this.data.pageNum === 1) {
                this.setData({
                  historyList: mockData
                });
              } else {
                // 追加数据
                this.setData({
                  historyList: [...this.data.historyList, ...mockData]
                });
              }
            } else {
              console.error('获取历史数据失败: 响应数据格式不正确', responseData);
              wx.showToast({ title: '获取历史数据失败', icon: 'none' });
            }
          } else {
            console.error('获取历史数据失败: 响应状态码或数据为空', res);
            wx.showToast({ title: '获取历史数据失败', icon: 'none' });
          }
        } catch (error) {
          console.error('解析历史数据失败:', error);
          wx.showToast({ title: '数据格式错误', icon: 'none' });
        }
      },
      fail: (err) => {
        console.error('请求失败:', err);
        wx.showToast({ title: '网络错误', icon: 'none' });
      },
      complete: () => {
        wx.hideLoading();
      }
    });
  },
  
  // 加载热力图数据
  loadHeatmap: function() {
    const app = getApp();
    const baseUrl = app.getBaseURL();
    const gameCode = this.data.currentGame;
    const periodCount = this.data.periodFilter;
    
    console.log('开始加载号码分布数据，游戏代码:', gameCode, '期数:', periodCount, 'API URL:', `${baseUrl}/api/results/${gameCode}/distribution?periodCount=${periodCount}`);
    
    wx.showLoading({ title: '加载中' });
    
    wx.request({
      url: `${baseUrl}/api/results/${gameCode}/distribution?periodCount=${periodCount}`,
      method: 'GET',
      success: (res) => {
        try {
          console.log('号码分布API响应状态码:', res.statusCode);
          console.log('号码分布API响应数据:', res.data);
          
          // 确保响应是有效的JSON
          const responseData = res.data;
          
          if (res.statusCode === 200 && responseData && responseData.data) {
            const redDistribution = responseData.data.red || [];
            const blueDistribution = responseData.data.blue || [];
            
            // 计算最大频率值
            let maxFreq = 0;
            
            // 检查红球数据
            if (redDistribution && redDistribution.length > 0) {
              redDistribution.forEach(item => {
                if (item.frequency > maxFreq) {
                  maxFreq = item.frequency;
                }
              });
            }
            
            // 检查蓝球数据
            if (blueDistribution && blueDistribution.length > 0) {
              blueDistribution.forEach(item => {
                if (item.frequency > maxFreq) {
                  maxFreq = item.frequency;
                }
              });
            }
            
            // 转换数据格式以适配前端显示
            const redHeatmap = redDistribution.map(item => ({
              number: item.number,
              frequency: item.frequency
            }));
            
            const blueHeatmap = blueDistribution.map(item => ({
              number: item.number,
              frequency: item.frequency
            }));
            
            this.setData({
              redHeatmap: redHeatmap,
              blueHeatmap: blueHeatmap,
              maxFrequency: maxFreq || 1
            });
          } else {
            console.error('获取号码分布数据失败: 响应数据格式不正确', responseData);
            this.generateMockHeatmapData(); // 失败时生成模拟数据
          }
        } catch (error) {
          console.error('解析号码分布数据失败:', error);
          this.generateMockHeatmapData(); // 解析失败时生成模拟数据
        }
      },
      fail: (err) => {
        console.error('获取号码分布数据请求失败:', err);
        this.generateMockHeatmapData(); // 请求失败时生成模拟数据
      },
      complete: () => {
        wx.hideLoading();
      }
    });
  },
  
  // 生成模拟热力图数据（备用）
  generateMockHeatmapData: function() {
    const redMax = this.data.currentGame === 'ssq' ? 33 : 35;
    const blueMax = this.data.currentGame === 'ssq' ? 16 : 12;
    
    // 模拟热力图数据
    const redHeatmap = [];
    const blueHeatmap = [];
    let maxFreq = 0;
    
    for (let i = 1; i <= redMax; i++) {
      const frequency = Math.floor(Math.random() * 10);
      redHeatmap.push({ number: i, frequency: frequency });
      maxFreq = Math.max(maxFreq, frequency);
    }
    
    for (let i = 1; i <= blueMax; i++) {
      const frequency = Math.floor(Math.random() * 5);
      blueHeatmap.push({ number: i, frequency: frequency });
      maxFreq = Math.max(maxFreq, frequency);
    }
    
    this.setData({
      redHeatmap: redHeatmap,
      blueHeatmap: blueHeatmap,
      maxFrequency: maxFreq || 1
    });
  },
  
  // 加载遗漏数据
  loadMissing: function() {
    const redMax = this.data.currentGame === 'ssq' ? 33 : 35;
    const blueMax = this.data.currentGame === 'ssq' ? 16 : 12;
    
    // 模拟遗漏数据
    const redMissing = [];
    const blueMissing = [];
    
    for (let i = 1; i <= redMax; i++) {
      redMissing.push({
        number: i,
        currentMissing: Math.floor(Math.random() * 20),
        maxMissing: Math.floor(Math.random() * 50) + 20,
        count: Math.floor(Math.random() * 100)
      });
    }
    
    for (let i = 1; i <= blueMax; i++) {
      blueMissing.push({
        number: i,
        currentMissing: Math.floor(Math.random() * 10),
        maxMissing: Math.floor(Math.random() * 30) + 10,
        count: Math.floor(Math.random() * 50)
      });
    }
    
    this.setData({
      redMissing: redMissing,
      blueMissing: blueMissing
    });
  },
  
  // 上滑加载更多
  onReachBottom: function() {
    if (this.data.currentTab === 'history' && this.data.hasMoreData) {
      // 增加页码
      this.setData({
        pageNum: this.data.pageNum + 1
      });
      // 加载更多数据
      this.loadHistory();
    }
  },
  
  // 切换奖品信息显示
  togglePrize: function(e) {
    const id = e.currentTarget.dataset.id;
    const historyList = this.data.historyList.map(item => {
      if (item.id === id) {
        item.showPrize = !item.showPrize;
      }
      return item;
    });
    
    this.setData({
      historyList: historyList
    });
  },
  
  // 排序数据遗漏表格
  sortMissingData: function(e) {
    const column = e.currentTarget.dataset.column;
    const currentGame = this.data.currentGame;
    
    // 获取当前游戏的遗漏数据
    let missingData = currentGame === 'ssq' ? this.data.redMissing : this.data.blueMissing;
    
    // 排序逻辑
    missingData.sort((a, b) => {
      if (a[column] < b[column]) return -1;
      if (a[column] > b[column]) return 1;
      return 0;
    });
    
    // 更新数据
    if (currentGame === 'ssq') {
      this.setData({
        redMissing: missingData
      });
    } else {
      this.setData({
        blueMissing: missingData
      });
    }
  },
  
  // 生成模拟历史数据
  generateMockHistoryData: function(gameCode) {
    const mockData = [];
    const isSSQ = gameCode === 'ssq';
    const redCount = isSSQ ? 6 : 5;
    const blueCount = isSSQ ? 1 : 2;
    const redMax = isSSQ ? 33 : 35;
    const blueMax = isSSQ ? 16 : 12;
    const pageNum = this.data.pageNum;
    const pageSize = this.data.pageSize;
    
    // 计算当前页应该生成的数据数量
    const totalMockData = 30; // 总共模拟30期数据，支持3页分页
    const startIndex = (pageNum - 1) * pageSize;
    const endIndex = Math.min(startIndex + pageSize, totalMockData);
    
    // 生成当前页的模拟数据
    for (let i = startIndex; i < endIndex; i++) {
      // 生成红球号码
      const redNumbers = [];
      while (redNumbers.length < redCount) {
        const num = Math.floor(Math.random() * redMax) + 1;
        if (!redNumbers.includes(num)) {
          redNumbers.push(num);
        }
      }
      // 排序
      redNumbers.sort((a, b) => a - b);
      
      // 生成蓝球号码
      const blueNumbers = [];
      while (blueNumbers.length < blueCount) {
        const num = Math.floor(Math.random() * blueMax) + 1;
        if (!blueNumbers.includes(num)) {
          blueNumbers.push(num);
        }
      }
      // 排序
      blueNumbers.sort((a, b) => a - b);
      
      // 生成期号，大乐透使用5位数字，双色球使用6位数字
      let period;
      if (gameCode === 'dlt') {
        // 大乐透：202加上三位数字，例如202301
        period = '202' + String(100 + i).slice(-3);
      } else {
        // 双色球：2025加上三位数，例如2025001
        period = '2025' + String(100 + i).slice(-3);
      }
      
      // 生成日期
      const now = new Date();
      now.setDate(now.getDate() - i * (isSSQ ? 2 : 3)); // 双色球每2天一期，大乐透每3天一期
      const date = now.toISOString().split('T')[0];
      
      mockData.push({
        id: 9000 + i,
        period: period,
        redNumbers: redNumbers,
        blueNumbers: blueNumbers,
        date: date,
        salesAmount: (Math.random() * 1000000000).toFixed(2),
        prizePool: (Math.random() * 5000000000).toFixed(2),
        showPrize: false,
        prizes: []
      });
    }
    
    // 检查是否还有更多数据
    this.setData({
      hasMoreData: endIndex < totalMockData
    });
    
    return mockData;
  }
});
