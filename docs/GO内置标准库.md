一、Go 内置标准库（不用安装，直接 import）
1. 基础工具
fmt：打印、格式化输出
os：操作文件、环境变量、命令行
io / ioutil：读写流、文件
strconv：字符串与数字互转
strings：字符串处理
time：时间、日期、定时器
sync：并发锁、等待组
2. 网络 & HTTP
net/http：写 HTTP 服务 / 客户端
net/url：URL 解析
net：TCP/UDP 编程
3. 数据格式
encoding/json：JSON 序列化
encoding/xml：XML 处理
encoding/base64：base64
4. 数据库
database/sql：数据库接口
5. 日志
log：简单日志