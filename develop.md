# Go mirror cache
# 可缓存的软件镜像服务

1. 启动一个web服务器，监听指定的多个web目录， 如果请求的文件存在， 直接下载，不存在就向配置好的远程服务器获取
2. 可配置镜像的同步（缓存过期）时间， 超过时间的， 重新向远程服务器发送head请求检查是否需要更新
3. 文件第一次下载时，直接向客户端传送的同时， 也在本地保存一份
4. 需要考虑文件锁， 文件被下载过程中不能被覆盖写入
5. 文件的下载记录支持输出到kafka，elasticsearch, redis, amqp 等


# 与 nginx 结合
用nginx反向代理到本地的指定端口, 所以本地可以开多个cache server,
ubuntu => /ubuntu
centos => /centos
npm => npm


# 配置文件格式

listen: 服务监听地址
blocksize: 缓冲区块大小

mirrors:
    name: 镜像名
        proxytarget: 代理地址
        localdir: 本地保存路径
        expire: 缓存过期时间
        prefix: 请求的url前缀
        paths:
            /path/to/index: 360
logger：
    file：
        filepath: data/logs/

# 日志格式
{
    uri
    mirror_name:
    status_code
    hit: cached | download
    size: xx
    remote_ip
    timestamp
}
