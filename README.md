# System Information Tool

[![Build](https://github.com/junler/sysinfo/actions/workflows/build.yml/badge.svg)](https://github.com/junler/sysinfo/actions/workflows/build.yml)
[![CI](https://github.com/junler/sysinfo/actions/workflows/ci.yml/badge.svg)](https://github.com/junler/sysinfo/actions/workflows/ci.yml)
[![Release](https://github.com/junler/sysinfo/actions/workflows/release.yml/badge.svg)](https://github.com/junler/sysinfo/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/junler/sysinfo)](https://goreportcard.com/report/github.com/junler/sysinfo)

一个功能强大的系统信息展示工具，支持命令行和Web界面两种展示方式。

## 功能特性

- **全面的系统信息收集**：CPU、内存、磁盘、网络、进程等
- **增强的监控指标**：负载平均值、I/O统计、温度传感器、交换内存等
- **进程监控**：显示CPU和内存占用最高的进程，包括详细的进程信息
- **网络统计**：详细的网络接口信息、流量统计、错误和丢包统计
- **开放端口监控**：显示所有TCP/UDP端口及其对应进程
- **系统服务状态**：显示关键系统服务的运行状态
- **用户会话监控**：显示当前登录的用户及其会话信息
- **磁盘I/O性能**：读写操作统计和平均响应时间
- **命令行界面**：快速查看系统信息和详细监控数据
- **Web界面**：使用TailwindCSS + HyperUI 的现代化Web界面
- **实时刷新**：Web界面支持自动刷新
- **嵌入式资源**：Web资源以embed方式打包，单文件部署
- **丰富的API接口**：提供RESTful API用于集成和自动化

## 安装

### 从 GitHub Releases 下载

前往 [Releases 页面](https://github.com/junler/sysinfo/releases) 下载适合你系统的预编译二进制文件：

- Linux (x64): `sysinfo-linux-amd64.tar.gz`
- Linux (ARM64): `sysinfo-linux-arm64.tar.gz`
- macOS (x64): `sysinfo-darwin-amd64.tar.gz`
- macOS (ARM64): `sysinfo-darwin-arm64.tar.gz`
- Windows (x64): `sysinfo-windows-amd64.zip`

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/junler/sysinfo.git
cd sysinfo

# 构建应用
make build

# 或者直接使用 go build
go build -o sysinfo .
```

## 使用方法

### 命令行界面

```bash
# 显示帮助信息
./sysinfo

# 显示系统信息
./sysinfo info

# 显示详细监控数据 (新功能)
./sysinfo monitor

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
   - 负载平均值（1分钟、5分钟、15分钟）

2. **CPU信息**
   - CPU型号和规格
   - 物理核心和逻辑核心数
   - 频率信息
   - 每核心使用率（实时图表）
   - CPU温度（如果可用）
   - 缓存大小

3. **内存信息**
   - 总内存、已用内存、可用内存
   - 缓存和缓冲区信息
   - 活跃和非活跃内存
   - 共享内存使用量
   - 内存使用率图表

4. **交换内存**
   - 交换分区总大小
   - 已用交换空间
   - 交换使用率

5. **磁盘信息**
   - 所有磁盘分区信息
   - 使用量和可用空间
   - 文件系统类型
   - inode使用情况
   - 使用率图表

6. **网络信息**
   - 网络接口列表和详细配置
   - 网络流量统计（发送/接收字节数和包数）
   - 网络错误和丢包统计
   - IP地址和硬件地址信息

7. **进程监控**
   - CPU使用率最高的进程
   - 内存使用情况
   - 进程状态和用户信息
   - 命令行参数

8. **I/O统计**
   - 磁盘读写字节数
   - 读写操作次数
   - 平均响应时间

9. **开放端口**
   - TCP/UDP端口列表
   - 对应进程信息
   - PID和绑定地址

10. **用户会话**
    - 当前登录用户
    - 终端和主机信息
    - 登录时间

11. **系统服务**
    - 关键系统服务状态
    - 服务进程ID
    - 运行状态

## API接口

Web服务器提供以下API接口：

### 基础接口
- `GET /api/info` - 获取完整系统信息（JSON格式）
- `GET /api/ports` - 获取开放端口信息（JSON格式）
- `GET /api/health` - 健康检查

### 增强监控接口 (新增)
- `GET /api/monitoring` - 获取核心监控数据（包含CPU、内存、I/O、网络等）
- `GET /api/processes` - 获取top进程列表
- `GET /api/iostats` - 获取磁盘I/O统计
- `GET /api/temperature` - 获取温度传感器数据
- `GET /api/users` - 获取当前登录用户信息
- `GET /api/services` - 获取系统服务状态

### API响应示例

#### /api/monitoring
```json
{
  "load_average": {
    "load1": 3.85,
    "load5": 3.25,
    "load15": 3.24
  },
  "memory": {
    "total": 17179869184,
    "used": 11357822976,
    "used_percent": 66.11,
    "available": 5822046208
  },
  "cpu": {
    "model_name": "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
    "cores": 6,
    "usage": [37.0, 0, 36.9, 0.99, 24.2, 0]
  },
  "io_stats": {
    "disk_read_bytes": 2355452559360,
    "disk_write_bytes": 2100280946688
  }
}
```

#### /api/processes
```json
{
  "processes": [
    {
      "pid": 98143,
      "name": "Code Helper",
      "cpu_percent": 13.7,
      "mem_percent": 3.8,
      "status": "sleep",
      "username": "sunjun"
    }
  ]
}
```

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

### 示例输出

#### 基础命令行输出

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

#### 增强监控输出 (新功能)

```text
=== SYSTEM MONITORING DASHBOARD ===
Host: junler-mac-pro-2.local | OS: darwin 15.5 | Uptime: 312 hours
Architecture: x86_64 | Kernel: 24.5.0

=== LOAD AVERAGE ===
1min:   3.57 | 5min:   3.09 | 15min:   3.20

=== CPU METRICS ===
Model: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
Cores: 6 Physical / 1 Logical | Frequency: 2600 MHz
Per-Core Usage:  52.5%   2.0%  45.0%   0.0%  46.0%   0.0%

=== MEMORY METRICS ===
RAM:  Total:    17.2 GB | Used:    11.3 GB ( 66.0%) | Available:     5.8 GB
Swap: Total:     4.3 GB | Used:     2.6 GB ( 59.7%) | Free:       1.7 GB

=== TOP PROCESSES BY CPU ===
PID      NAME                 STATUS         CPU%     MEM%    RSS(MB) USER           
------------------------------------------------------------------------------------------
63067    sysinfo              sleep         25.3%     0.1%      15.0 sunjun         
98143    Code Helper (Rend... sleep         13.6%     4.0%     659.8 sunjun         
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

- [x] 添加更多系统监控指标
  - [x] 负载平均值监控
  - [x] 交换内存统计
  - [x] 进程详细信息（CPU、内存使用率排序）
  - [x] 增强的网络统计（包括错误和丢包）
  - [x] 磁盘I/O性能统计
  - [x] 用户会话监控
  - [x] 系统服务状态监控
  - [x] 详细的磁盘使用信息（包括inode）
  - [x] 温度监控框架（待平台实现）
  - [x] 新的`monitor`命令用于详细监控展示
  - [x] 增强的API接口（/api/monitoring等）
- [ ] 支持历史数据存储和图表
- [ ] 添加警报和通知功能
- [ ] 支持多节点监控
- [ ] Docker容器化部署

## 开发

### 本地开发

```bash
# 安装依赖
make install

# 运行测试
make test

# 构建
make build

# 启动 Web 服务器
make run-web
```

### 发布新版本

1. 更新 `Makefile` 中的 `VERSION`
2. 提交代码并推送到 main 分支
3. 创建并推送标签：

   ```bash
   make tag
   git push origin v1.0.0
   ```

4. GitHub Actions 将自动构建并创建 Release

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！
