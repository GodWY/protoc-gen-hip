# protoc-gen-hip 工具介绍

## 简介
protoc-gen-hip 是一款基于protoc的插件，通过定义proto文件可以自动生成基于gin的http接口

## 下载

```git
   // go1.17前
    go get -u github.com/GodWY/protoc-gen-hip@latest
    go install github.com/GodWY/protoc-gen-hip
    // go1.17后
    go install  github.com/GodWY/protoc-gen-hip@latest
```

## 验证

```shell
protoc-gen-hip --version
protoc-gen-hip v1.2.0
```

## 定义服务

```protobuf
syntax="proto3";
package greeter;
option go_package="examples/gen";
// test
message Request{

}
message Response{}


// @root:api/login  @middle:gin.Logger() @middle:gin.Recovery()
// @doc:this is a test
service Login{
  //@method:GET
  rpc GetUserName(Request)returns(Response){};
  //@method:POST @middle:gin.Logger() @middle:gin.Recovery() @after:gin.Recovery()
  rpc GetUserID(Request)returns(Response){};
}
```

1. 生成代码
```shell
protoc  --plugin=./protoc-gen-hip --go_out=./ --hip_out=./ examples/greeter.proto
```

2. 介绍

> 定义http服务与定义rpc服务完全一样，通过注释完成http代码的生成

3. 注解介绍

| 注解名称 | 作用   | 默认值 |  
| ----- | --------- | ----------- 
| @method | 自定义方法 | GET           
| @middle  | 中间件     | 无
| @root | 自定义组路径| api/{{组名}}
| @import | 导入包| 无
| @after| 后置中间件，函数执行完成后执行| 无

4. 结构
```shell
{
    "code": 0,
    "detail": "success",
    "data": {
       // 定义的pb返回值
    }
}
```

## 总结

使用proto的方式可以大大减少重复编写handle逻辑，同时可以很好的思考服务，结构化服务，更好的和后端交互