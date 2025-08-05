# System Information Tool

一个功能强大的系统信息展示工具，支持命令行和Web界面两种展示方式。

## 功能特性

- **全面的系统信息收集**：CPU、内存、磁盘、网络、进程等
- **开放端口监控**：显示所有TCP/UDP端口及其对应进程
- **命令行界面**：快速查看系统信息
- **Web界面**：使用TailwindCSS + HyperUI 的现代化Web界面
- **实时刷新**：Web界面支持自动刷新
- **嵌入式资源**：Web资源以embed方式打包，单文件部署

## 安装

```bash
# 克隆仓库
git clone https://github.com/junler/sysinfo.git
cd sysinfo

# 构建应用
go build -o sysinfo .
```

## 使用方法

### 命令行界面

```bash
# 显示帮助信息
./sysinfo

# 显示系统信息
./sysinfo info

# 显示开放端口
./sysinfo ports

# 显示详细帮助
./sysinfo --help
```

### Web界面

```bash
# 启动Web服务器（默认端口8080）
./sysinfo serve

# 指定端口启动Web服务器
./sysinfo serve --port 9090

# 或使用Makefile
make run-web
```

然后在浏览器中访问 `http://localhost:8080`

## Web界面功能

Web界面提供以下信息的实时展示：

1. **系统概览**
   - 操作系统版本
   - 主机名
   - 系统运行时间
   - 系统架构

2. **CPU信息**
   - CPU型号和规格
   - 物理核心和逻辑核心数
   - 频率信息
   - 每核心使用率（实时图表）

3. **内存信息**
   - 总内存、已用内存、可用内存
   - 缓存和缓冲区信息
   - 内存使用率图表

4. **磁盘信息**
   - 所有磁盘分区信息
   - 使用量和可用空间
   - 文件系统类型
   - 使用率图表

5. **网络信息**
   - 网络接口列表
   - 网络流量统计
   - IP地址信息

6. **开放端口**
   - TCP/UDP端口列表
   - 对应进程信息
   - PID和绑定地址

## API接口

Web服务器提供以下API接口：

- `GET /api/info` - 获取系统信息（JSON格式）
- `GET /api/ports` - 获取开放端口信息（JSON格式）
- `GET /api/health` - 健康检查

## 系统要求

- Go 1.24.1 或更高版本
- 支持 Linux、macOS、Windows

## 依赖项

- `github.com/shirou/gopsutil/v3` - 系统信息收集
- `github.com/spf13/cobra` - 命令行界面
- `github.com/gin-gonic/gin` - Web服务器

## 技术特点

1. **现代化UI**：使用TailwindCSS和HyperUI组件库
2. **响应式设计**：支持桌面和移动端
3. **实时更新**：Web界面每30秒自动刷新
4. **嵌入式资源**：使用Go 1.16+ embed特性，单文件部署
5. **RESTful API**：标准化的API接口
6. **跨平台支持**：支持主流操作系统

## 示例输出

### 命令行输出

```text
=== System Information ===
OS: darwin 15.5
Hostname: junler-mac-pro-2.local
Architecture: x86_64
Kernel Version: 24.5.0
Uptime: 311 hours
Last Boot: 2025-07-23 17:18:11
Process Count: 833

=== CPU Information ===
Model: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
Physical Cores: 6
Logical Cores: 12
Frequency: 2600.00 MHz
CPU Usage per core: 58.0%, 2.0%, 41.0%, 0.0%, 36.0%, 0.0%...
```

### API JSON输出

```json
{
  "os": "darwin 15.5",
  "hostname": "junler-mac-pro-2.local",
  "uptime": "311 hours",
  "cpu": {
    "model_name": "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
    "cores": 6,
    "logical_cores": 12,
    "usage": [58.0, 2.0, 41.0, 0.0, 36.0, 0.0],
    "frequency": 2600.0
  },
  "memory": {
    "total": 18446744073709551616,
    "used": 11928395776,
    "used_percent": 64.5,
    "free": 1213595648,
    "available": 6545543168
  }
}
```

## 开发计划

- [ ] 添加更多系统监控指标
- [ ] 支持历史数据存储和图表
- [ ] 添加警报和通知功能
- [ ] 支持多节点监控
- [ ] Docker容器化部署

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！
