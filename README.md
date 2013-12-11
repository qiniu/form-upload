Qiniu Form Upload Demo
==========

[![Qiniu Logo](http://qiniutek.com/images/logo-2.png)](http://qiniu.com/)

这里整理了上传文件到七牛云存储的样例集。这些样例并不需要你基于任何七牛SDK来完成，方便你移植到各种技术框架体系中。


## Web (html)

这个其实是最简单的。因为七牛的上传API就是一个multipart form形式的表单。你只需要提交如下这个表单：

```html
<form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
  <input name="key" type="hidden" value="<Your file name in qiniu>">
  <input name="x:<custom_field_name>" type="hidden" value="<Value of your custom param>">
  <input name="token" type="hidden" value="<Your uptoken from server>">
  <input name="file" type="file"/>
</form>
```

关键是 token 字段。这个通常是服务端生成的。比较正常点的做法是需要服务器端配合的，在客户端需要上传或修改文件的时候服务端生成token并提供给客户端。

但如果你的需求非常简单，也可以有不需要服务端的做法：随便找一个语言的 sdk 生成一个 token（token的过期时间可以设置很长，比如100年，这样就得到了一个永不过期的token —— 除非你在 https://portal.qiniu.com/ 里面把生成这个token的accessKey/secretKey作废）。然后把这个 token 写死在以上的表单里就行。

服务端配合的常见样例（会不断补充各种语言的DEMO）：

* https://github.com/qiniu/form-upload/tree/develop/web


## Flash (action script)

能够理解上面 Web 方式的上传过程，用 Flash 只是通过 http 协议把表单发送出去而已。直接贴代码：

```ActionScript
var u :URLRequest = new URLRequest('http://up.qiniu.com');
u.method = URLRequestMethod.POST;
u.requestHeaders = [new URLRequestHeader('enctype', 'multipart/form-data')];

var ur :URLVariables = new URLVariables();
ur.key = '<Your file name in qiniu>';
ur.token = '<Your uptoken from server>';
ur['x:<custom_field_name>'] = '<Value of your custom param>';

u.data = ur;
 
f.upload(u, 'file'); // f is File or FileReference is both OK, but UploadDataFieldName must be 'file'
f.addEventListener(DataEvent.UPLOAD_COMPLETE_DATA, uploadedHandler);
 
protected function uploadedHandler(event:DataEvent):void
{
  trace(event.data);
  //{
  //  "hash":"<File etag>",
  //  "key":"<Your file name in qiniu>",
  //  "x:<custom_field_name>":"<Value of your custom param>"
  //}
}
```

这个样例由 @mani95lisa 提供。原文参见：https://gist.github.com/mani95lisa/7912530
