## CHANGE LOG

### v5.0.0

2013-06-12 issue [#5](https://github.com/qiniu/form-upload/pull/5)

- Golang 版的 form upload 样例，演示内容：
  - 两种上传模式
    - upload（自动生成 key，可自动消重）
    - uploadWithKey（手工指定key，同名会导致上传失败，要想支持 overwrite 请修改 PutPolicy.Scope）
  - 用户自定义域（x:开头的字段）
  - 演示如何生成私有资源的临时下载链接（privateUrl）
- 遵循 [sdkspec v1.0.2](https://github.com/qiniu/sdkspec/tree/v1.0.2)
  - rs.GetPolicy 删除 Scope，也就是不再支持批量下载的授权。
  - rs.New, PutPolicy.Token, GetPolicy.MakeRequest 增加 mac *digest.Mac 参数。

