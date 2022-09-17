# Customize

## Websocket Handler

With the implementation of a function like `func (...args) func(pool *pool.Pool)`, you can replace the default handler with this function.

For example, if you want to use `fmt.Println()` to print all danmaku, you can write a function like the following code:

```go
func MyHandler(pool *pool.Pool) { 
    for {
        select {
        case src := <-pool.Danmaku:
            fmt.Println(string(src))
        default:
        }
    }
}
```

For further details, you can read the source code in `src/pool/handler.go`. The function `EmptyHandler() func(pool *Pool)` shows an example.

## Http Handler

You can add http handler with a format as `gin.HandlerFunc`.

For example, this is a hello-world service handler:

```go
package your_package

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	"github.com/sunist-c/bililive-danmaku/service/api"
)

func HelloWorldHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(200, "hello, world!")
	}
}

func init() {
	api.AddService(http.MethodGet, "/hello-world", HelloworldHandler())
}
```

Then import your package to `main.go`, that's all the operations. 