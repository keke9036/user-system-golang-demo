# Start

1. build user_bizserver
> go build -o user_bizserver ./bizsvr/user_bizserver.go
2. build user_apiserver
> go build -o user_apiserver ./apisvr/user_apiserver.go
3. copy ./conf/config_rpcserver.yaml, edit RpcServer.Addr&WorkerId 
> cp ./conf/config_rpcserver.yaml ./conf/config_rpcserver_1.yaml
4. edit RpcClient.Addrs in ./conf/config_httpserver.yaml, according to step 3
5. run MySQL & Redis
> redis-server
6. run user_bizserver 
> ./user_bizserver -config=config_rpcserver_1 -logLevel=Info -logPath=stdout
7. run a HTTP server
> ./user_apiserver -config=config_httpserver -logLevel=Info -logPath=stdout
8. view log in terminal

# Benchmark
通过./benchmark/init_db_data_py初始化测试数据
> 测试用户数据格式如下
> - username: testu_{%d}
> - password: rootpwd
