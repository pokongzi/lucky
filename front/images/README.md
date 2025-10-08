# TabBar 图标文件

## 已创建的SVG图标

我已经为你创建了6个SVG格式的图标文件：

- `home.svg` / `home-active.svg` - 首页图标
- `history.svg` / `history-active.svg` - 历史开奖图标  
- `numbers.svg` / `numbers-active.svg` - 我的号码图标

## 转换为PNG的方法

### 方法1：使用在线转换工具
1. 访问 https://convertio.co/svg-png/ 或 https://cloudconvert.com/svg-to-png
2. 上传SVG文件
3. 设置输出尺寸为 81x81 像素
4. 下载PNG文件

### 方法2：使用浏览器
1. 打开 `create-icons.html` 文件
2. 点击每个图标的"下载"按钮
3. 自动下载PNG格式文件

### 方法3：使用命令行工具
```bash
# 安装ImageMagick (macOS)
brew install imagemagick

# 运行转换脚本
./convert-icons.sh
```

## 图标设计说明

- **尺寸**: 81x81 像素 (微信小程序标准)
- **风格**: 简洁的线条图标
- **颜色**: 未选中状态 #666666，选中状态 #ff4757
- **来源**: 基于Feather Icons开源图标库设计

## 文件列表

转换完成后，你需要以下6个PNG文件：
- home.png
- home-active.png  
- history.png
- history-active.png
- numbers.png
- numbers-active.png