# 变更日志

## v1.2.0

1. 增加 http 中间件：
   - xss 过滤器： `middleware.XSSFilter()`
   - cors 处理器：`middleware.CORS()`
   - 有状态的 jwt 处理器：`middleware.JWTStatefulWith()`，`middleware.JWTStatefulWithout()`
   - session 管理器：`middleware.SessionWithStore()` 方便分布式部署，指定 session 存储器
   - http 缓存器：`middleware.HttpCache()`
   - 角色控制器：`middleware.WithRole()`
2. 增加 json logger 格式，为适应云日志采集(阿里云/腾讯云)
3. 增加 db 内部特地错误检查工具：
   - 唯一键冲突：`dberror.IsDuplicateEntry()`
4. 增加 动态表格结构的 restful 响应：`restful.TableWithPagination()`
5. 增加 [HTTPSQS](http://zyan.cc/httpsqs/) 客户端及相关队列：
   - 客户端：`httpsqs.NewClient()`
   - 消费者：`listener.NewHttpSQSConsumer()`
     - 其他情形：
     ```golang
     err := listener.NewHttpSQSConsumer(handler).
        WithContext(ctx).
        WithLog(global.Log).
        Consume()
     ```
   - 消费者包装器：`listener.HTTPSQS(handler)`
6. 其他组件升级：
   - DB 分页查询参数 `pager.Option{}` 支持 `Preload`

⚠️注意：本次升级存在不向下兼容部分：
   - http 中间件：
     - `middleware.HttpLogger(restful.LogOption{})` => `middleware.HttpLogger(middleware.HttpLoggerOption{})`


## v1.1.0

1. 升级 `gorm v2`，并增加 `dbutil.ConnectWith()` 方法以适应非 mysql、sqlite 的数据库连接

## v1.0.0