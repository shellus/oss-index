# OSS-Index

给阿里云oss储存增加列出目录功能


配置文件默认获取用户目录下的`.ossutilconfig`文件，配置文件格式为：
```ini
[Credentials]
language=CH
endpoint=endpoint_host
accessKeyID=accessKeyID
accessKeySecret=accessKeySecret
bucket=bucket_name

```

运行后会生成`.oss_index_meta`文件在每个目录下。里面包含当前路径和目录下的所有文件和文件夹

当生成完成后，你需要上传`src/templates/404.html`到oss。然后在阿里云oss控制台设置404页面。

当一切完成，你应该可以看到类似下面这样的界面

![image](https://github.com/shellus/oss-index/raw/master/thumb/64304.jpg)