# 基础包


## 配置文件

(/etc/chuanmoon/odoo.conf)中有db的配置 

其他的配置存储在 cy_matedata 表

``` conf
[options]
; database postgresql
db_host = 172.17.0.1
db_port = 5432
db_user = odoo
db_password = 123456
dbfilter = chuanmoon

db_name = chuanmoon
db_query = 'sslmode=disable'
```

## cy_matedata

|key| 描述|
| - | - |
|cy_nats_url|nats server url，用于rpc通讯 |
|cy_readonly_db_url|只读数据库URL |
|cy_redis_addr|redis地址，用于用户信息|
|cy_redis_auth|redis密码|
|cy_elasticsearch_url|elasticsearch url，用于商品列表与搜索|
|cy_clickhouse_addr|clickhouse地址，用于埋点|
|cy_clickhouse_database|clickhouse数据库|
|cy_clickhouse_username|clickhouse用户名|
|cy_clickhouse_password|clickhouse密码|
|cy_image_link_prefix|图片链接前缀，如 https://img.your_domain.com/|
|cy_web_link_prefix|网站链接前缀，如 https://www.your_domain.com/|
|cy_gateway_link_prefix|api网关链接前缀，如 https://gw.your_domain.com/|
|cy_gateway_internal_link_prefix|api网关内部链接前缀，如 http://192.168.1.10:8080/ ，用于内部服务器(如odoo)调用网关|
|cy_gateway_sign_key|api网关调用签名key|

