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
    redFrequencyData: [],
    blueFrequencyData: [],
    // 奇偶比和大小比数据
    redOddEvenRatio: { odd: 0, even: 0 },
    redBigSmallRatio: { big: 0, small: 0 },
    blueOddEvenRatio: { odd: 0, even: 0 },
    blueBigSmallRatio: { big: 0, small: 0 },
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
  
  // 切换标签页
  switchTab: function(e) {
    const tab = e.currentTarget.dataset.tab;
    this.setData({
      currentTab: tab
    });
    
    // 切换到号码分布标签时绘制频率图表和饼图
    if (tab === 'distribution') {
      // 延迟一点时间再绘制，确保DOM已经更新
      setTimeout(() => {
        this.drawFrequencyChart('redFrequencyChart', this.data.redFrequencyData, true);
        this.drawFrequencyChart('blueFrequencyChart', this.data.blueFrequencyData, false);
        
        // 绘制奇偶比和大小比饼图
        this.drawPieChart('redOddEvenChart', this.data.redOddEvenRatio, ['奇数', '偶数'], ['#e23e3e', '#3182ce']);
        this.drawPieChart('redBigSmallChart', this.data.redBigSmallRatio, ['大数', '小数'], ['#e23e3e', '#3182ce']);
        this.drawPieChart('blueOddEvenChart', this.data.blueOddEvenRatio, ['奇数', '偶数'], ['#3182ce', '#e23e3e']);
        this.drawPieChart('blueBigSmallChart', this.data.blueBigSmallRatio, ['大数', '小数'], ['#3182ce', '#e23e3e']);
      }, 100);
    }
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
    
    // 用户指定的接口URL格式
    const targetUrl = `http://47.121.26.190:8080/api/results/distribution/${gameCode}?periodCount=${periodCount}`;
    console.log('开始加载号码分布数据，游戏代码:', gameCode, '期数:', periodCount, '使用指定接口URL:', targetUrl);
    
    wx.showLoading({ title: '加载中' });
    
    wx.request({
      url: targetUrl, // 使用用户指定的接口URL
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
            
            // 计算奇偶比和大小比
            const isSSQ = gameCode === 'ssq';
            const redBigThreshold = isSSQ ? 17 : 18;
            const blueBigThreshold = isSSQ ? 8 : 6;
            
            // 计算红球奇偶比和大小比
            let redOdd = 0, redEven = 0, redBig = 0, redSmall = 0;
            redHeatmap.forEach(item => {
              const number = item.number;
              const frequency = item.frequency;
              
              // 奇偶统计（按出现次数加权）
              if (number % 2 === 1) {
                redOdd += frequency;
              } else {
                redEven += frequency;
              }
              
              // 大小统计（按出现次数加权）
              if (number > redBigThreshold) {
                redBig += frequency;
              } else {
                redSmall += frequency;
              }
            });
            
            // 计算蓝球奇偶比和大小比
            let blueOdd = 0, blueEven = 0, blueBig = 0, blueSmall = 0;
            blueHeatmap.forEach(item => {
              const number = item.number;
              const frequency = item.frequency;
              
              // 奇偶统计（按出现次数加权）
              if (number % 2 === 1) {
                blueOdd += frequency;
              } else {
                blueEven += frequency;
              }
              
              // 大小统计（按出现次数加权）
              if (number > blueBigThreshold) {
                blueBig += frequency;
              } else {
                blueSmall += frequency;
              }
            });
            
            this.setData({
              redHeatmap: redHeatmap,
              blueHeatmap: blueHeatmap,
              redFrequencyData: redHeatmap,
              blueFrequencyData: blueHeatmap,
              maxFrequency: maxFreq || 1,
              redOddEvenRatio: { odd: redOdd, even: redEven },
              redBigSmallRatio: { big: redBig, small: redSmall },
              blueOddEvenRatio: { odd: blueOdd, even: blueEven },
              blueBigSmallRatio: { big: blueBig, small: blueSmall }
            });
            
            // 数据更新后绘制频率图表和饼图
            this.drawFrequencyChart('redFrequencyChart', this.data.redFrequencyData, true);
            this.drawFrequencyChart('blueFrequencyChart', this.data.blueFrequencyData, false);
            this.drawPieChart('redOddEvenChart', this.data.redOddEvenRatio, ['奇数', '偶数'], ['#e23e3e', '#3182ce']);
            this.drawPieChart('redBigSmallChart', this.data.redBigSmallRatio, ['大数', '小数'], ['#e23e3e', '#3182ce']);
            this.drawPieChart('blueOddEvenChart', this.data.blueOddEvenRatio, ['奇数', '偶数'], ['#3182ce', '#e23e3e']);
            this.drawPieChart('blueBigSmallChart', this.data.blueBigSmallRatio, ['大数', '小数'], ['#3182ce', '#e23e3e']);
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
    
    // 计算奇偶比和大小比
    const isSSQ = this.data.currentGame === 'ssq';
    const redBigThreshold = isSSQ ? 17 : 18;
    const blueBigThreshold = isSSQ ? 8 : 6;
    
    // 计算红球奇偶比和大小比
    let redOdd = 0, redEven = 0, redBig = 0, redSmall = 0;
    redHeatmap.forEach(item => {
      const number = item.number;
      const frequency = item.frequency;
      
      // 奇偶统计（按出现次数加权）
      if (number % 2 === 1) {
        redOdd += frequency;
      } else {
        redEven += frequency;
      }
      
      // 大小统计（按出现次数加权）
      if (number > redBigThreshold) {
        redBig += frequency;
      } else {
        redSmall += frequency;
      }
    });
    
    // 计算蓝球奇偶比和大小比
    let blueOdd = 0, blueEven = 0, blueBig = 0, blueSmall = 0;
    blueHeatmap.forEach(item => {
      const number = item.number;
      const frequency = item.frequency;
      
      // 奇偶统计（按出现次数加权）
      if (number % 2 === 1) {
        blueOdd += frequency;
      } else {
        blueEven += frequency;
      }
      
      // 大小统计（按出现次数加权）
      if (number > blueBigThreshold) {
        blueBig += frequency;
      } else {
        blueSmall += frequency;
      }
    });
    
    this.setData({
      redHeatmap: redHeatmap,
      blueHeatmap: blueHeatmap,
      redFrequencyData: redHeatmap,
      blueFrequencyData: blueHeatmap,
      maxFrequency: maxFreq || 1,
      redOddEvenRatio: { odd: redOdd, even: redEven },
      redBigSmallRatio: { big: redBig, small: redSmall },
      blueOddEvenRatio: { odd: blueOdd, even: blueEven },
      blueBigSmallRatio: { big: blueBig, small: blueSmall }
    });
    
    // 生成模拟数据后绘制频率图表和饼图
    setTimeout(() => {
      this.drawFrequencyChart('redFrequencyChart', this.data.redFrequencyData, true);
      this.drawFrequencyChart('blueFrequencyChart', this.data.blueFrequencyData, false);
      this.drawPieChart('redOddEvenChart', this.data.redOddEvenRatio, ['奇数', '偶数'], ['#e23e3e', '#3182ce']);
      this.drawPieChart('redBigSmallChart', this.data.redBigSmallRatio, ['大数', '小数'], ['#e23e3e', '#3182ce']);
      this.drawPieChart('blueOddEvenChart', this.data.blueOddEvenRatio, ['奇数', '偶数'], ['#3182ce', '#e23e3e']);
      this.drawPieChart('blueBigSmallChart', this.data.blueBigSmallRatio, ['大数', '小数'], ['#3182ce', '#e23e3e']);
    }, 100);
  },
  
  // 绘制号码频率统计柱状图
  drawFrequencyChart: function(canvasId, data, isRed) {
    const ctx = wx.createCanvasContext(canvasId);
    const canvasWidth = 300; // 画布宽度
    const canvasHeight = 200; // 画布高度
    const padding = 20; // 边距
    const barWidth = 8; // 柱状图宽度
    const barGap = 4; // 柱状图间距
    
    // 清空画布
    ctx.clearRect(0, 0, canvasWidth, canvasHeight);
    
    // 计算图表区域大小
    const chartWidth = canvasWidth - 2 * padding;
    const chartHeight = canvasHeight - 40; // 预留底部空间显示号码
    
    // 找出数据中的最大值
    let maxValue = 0;
    data.forEach(item => {
      if (item.frequency > maxValue) {
        maxValue = item.frequency;
      }
    });
    
    if (maxValue === 0) maxValue = 1; // 避免除以0
    
    // 设置柱状图颜色
    const barColor = isRed ? '#e23e3e' : '#3182ce';
    
    // 绘制柱状图
    data.forEach((item, index) => {
      // 计算柱子高度
      const barHeight = (item.frequency / maxValue) * chartHeight;
      const x = padding + index * (barWidth + barGap);
      const y = canvasHeight - 30 - barHeight; // 减去底部空间
      
      // 绘制柱子
      ctx.setFillStyle(barColor);
      ctx.fillRect(x, y, barWidth, barHeight);
      
      // 绘制数值标签
      ctx.setFillStyle('#2d3748');
      ctx.setFontSize(10);
      ctx.fillText(item.frequency.toString(), x - 2, y - 5);
      
      // 绘制号码标签（每5个号码显示一个，避免拥挤）
      if (index % 5 === 0) {
        ctx.setFillStyle('#718096');
        ctx.setFontSize(9);
        ctx.fillText(item.number.toString(), x - 3, canvasHeight - 15);
      }
    });
    
    // 绘制坐标轴
    ctx.setStrokeStyle('#e2e8f0');
    ctx.setLineWidth(1);
    
    // X轴
    ctx.beginPath();
    ctx.moveTo(padding, canvasHeight - 30);
    ctx.lineTo(canvasWidth - padding, canvasHeight - 30);
    ctx.stroke();
    
    // Y轴
    ctx.beginPath();
    ctx.moveTo(padding, padding);
    ctx.lineTo(padding, canvasHeight - 30);
    ctx.stroke();
    
    // 绘制最大值标签
    ctx.setFillStyle('#718096');
    ctx.setFontSize(10);
    ctx.fillText(maxValue.toString(), 5, padding);
    
    // 绘制完成后绘制到画布上
    ctx.draw();
  },
  
  // 绘制饼图
  drawPieChart: function(canvasId, data, labels, colors) {
    const ctx = wx.createCanvasContext(canvasId);
    // 直接使用canvas的指定尺寸，不通过selectorQuery获取
    this._drawPieChartWithSize(ctx, 150, 150, data, labels, colors);
  },
  
  // 内部方法：根据指定尺寸绘制饼图
  _drawPieChartWithSize: function(ctx, canvasWidth, canvasHeight, data, labels, colors) {
    const centerX = canvasWidth / 2;
    const centerY = canvasHeight / 2;
    const radius = Math.min(centerX, centerY) - 20; // 进一步增加边距，确保图例完全显示在饼图下方
    
    // 获取数据总和
    const total = Object.values(data).reduce((sum, value) => sum + value, 0);
    
    // 如果没有数据，绘制一个灰色的饼图
    if (total === 0) {
      ctx.setFillStyle('#e2e8f0');
      ctx.beginPath();
      ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI);
      ctx.fill();
      
      // 绘制提示文字
      ctx.setFillStyle('#718096');
      ctx.setFontSize(12);
      ctx.setTextAlign('center');
      ctx.setTextBaseline('middle');
      ctx.fillText('暂无数据', centerX, centerY);
      ctx.draw();
      return;
    }
    
    let startAngle = 0;
    const keys = Object.keys(data);
    
    // 绘制饼图各部分
    for (let i = 0; i < keys.length; i++) {
      const key = keys[i];
      const value = data[key];
      const percentage = value / total;
      const endAngle = startAngle + percentage * 2 * Math.PI;
      
      // 绘制扇形
      ctx.setFillStyle(colors[i]);
      ctx.beginPath();
      ctx.moveTo(centerX, centerY);
      ctx.arc(centerX, centerY, radius, startAngle, endAngle);
      ctx.closePath();
      ctx.fill();
      
      // 计算扇形中间角度，用于放置标签
      const midAngle = startAngle + percentage * Math.PI;
      const labelX = centerX + Math.cos(midAngle) * (radius * 0.6);
      const labelY = centerY + Math.sin(midAngle) * (radius * 0.6);
      
      // 绘制百分比标签
      ctx.setFillStyle('#fff');
      ctx.setFontSize(12);
      ctx.setTextAlign('center');
      ctx.setTextBaseline('middle');
      ctx.fillText(`${Math.round(percentage * 100)}%`, labelX, labelY);
      
      // 更新起始角度
      startAngle = endAngle;
    }
    
    // 绘制图例
    const legendItemHeight = 12;
    const legendStartY = canvasHeight - 20; // 进一步调整图例起始位置，确保与饼图完全分离
    const legendItemWidth = 70;
    
    for (let i = 0; i < keys.length; i++) {
      // 计算图例项的中心位置
      const legendCenterX = centerX - (legendItemWidth * (keys.length - 1)) / 2 + i * legendItemWidth;
      
      // 绘制颜色块（位于中心偏左）
      const colorBlockX = legendCenterX - 15;
      ctx.setFillStyle(colors[i]);
      ctx.fillRect(colorBlockX, legendStartY, 10, 10);
      
      // 绘制标签文字（视觉居中）
      ctx.setFillStyle('#2d3748');
      ctx.setFontSize(9);
      ctx.setTextAlign('center'); // 设置文本居中对齐
      ctx.setTextBaseline('middle');
      ctx.fillText(`${labels[i]}`, legendCenterX + 10, legendStartY + 5);
    }
    
    // 绘制完成后绘制到画布上
    ctx.draw();
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
