mini-tun
===

超小型简易编程语言, 包含了词法分析,语法分析,语义分析以及一个解释器。可以用来当作学习编译原理的入门项目。

## Features
- [x] 支持变量声明、赋值、函数定义、函数调用
- [ ] 分支语句、循环语句

### quick start
```bash
go run ./cmd/interpreter/main.go ./example/add.tun
````
或编译后运行
```bash
go build -o tun ./cmd\interpreter/main.go
./tun ./example/add.tun
```
### example
```tun
let a = 3
let b = 2
let c = a - b
let add = function (a, b) {
    let c = a + b
    return a + c
}

let d = add(a, add(b, c))
```


### 1. 语法
TODO: 待补充

```
<statement> ::=  
```