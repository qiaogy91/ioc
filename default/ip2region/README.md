##### Ip2Region 说明
- 地址：`https://github.com/lionsoul2014/ip2region/tree/master`
- 说明：一个离线IP 地址定位框架，10微秒级别的查询效率、整个离线数据大小为11M
- 查询方式：
    - `vIndex`索引缓存：使用固定的 512KB 内存空间来作为缓存，减少从对磁盘中DB 发起的I/O操作，平均查询稳定在10-20微秒
    - `xdb` 文件缓存：将整改文件加载到内存，占用的内存大小等于DB 文件大小，从此查询DB 无磁盘I/O操作，保存10微秒内效率

##### 基于文件查询
```go
package main

func main() {
  var dbPath = "ip2region.xdb file path"
  searcher, err := xdb.NewWithFileOnly(dbPath)
  if err != nil {
    fmt.Printf("failed to create searcher: %s\n", err.Error())
    return
  }

  defer searcher.Close()
  region, err := searcher.SearchByStr("1.2.3.4")
}
```

##### 基于vIndex 查询
```go
package main

func main() {
  vIndex, err := LoadVectorIndexFromFile(dbPath)
  if err != nil {
    fmt.Printf("failed to load vector index from `%s`: %s\n", dbPath, err)
    return
  }

  searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex)
  if err != nil {
    fmt.Printf("failed to create searcher with vector index: %s\n", err)
    return
  }
  defer searcher.Close()
  region, err := searcher.SearchByStr("1.2.3.4")
}
```

##### 基于xdb 全文件缓存查询
```go
package main

func main() {
  // 1、从 dbPath 加载整个 xdb 到内存
  cBuff, err := LoadContentFromFile(dbPath)
  if err != nil {
    fmt.Printf("failed to load content from `%s`: %s\n", dbPath, err)
    return
  }

  // 2、用全局的 cBuff 创建完全基于内存的查询对象。
  searcher, err := xdb.NewWithBuffer(cBuff)
  if err != nil {
    fmt.Printf("failed to create searcher with content: %s\n", err)
    return
  }
  defer searcher.Close()
  region, err := searcher.SearchByStr("1.2.3.4")
}
```