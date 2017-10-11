# ginlogger

This package contains a simple Go http.Handler middleware that mimics the HTTP Request logger used in [gin](https://github.com/gin-gonic/gin).

Usage:

```go
myHandler := myCustomHandler() // satisfied http.Handler interface
myHandlerWithLogging := ginlogger.Handler(myHandler, os.Stdout) //os.Stdout or any io.Writer
```

Requests will be logged with colors, status code, latency etc.:
```
[HTTP] 2017/10/11 - 12:00:06 | 200 |      28.096µs |       127.0.0.1 | POST     /
```

Thanks to `github.com/gin-gonic/gin` for the inspiration (and most of the code)!
