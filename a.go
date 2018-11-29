package main

import (
    "fmt"
    //"time"
    "net/http"
    "github.com/hpcloud/tail"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    // 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func read(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hello")
    c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c.Close()
    t, err := tail.TailFile("/tmp/tail.log", tail.Config{Follow: true})
    for line := range t.Lines {
		err = c.WriteMessage(1, []byte(line.Text + "<br/>"))
		if err != nil {
            fmt.Println(err)
			break
		}
    }
}

func view(w http.ResponseWriter, r *http.Request) {
    body := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>tail log</title>
<script src="//cdn.bootcss.com/jquery/2.1.4/jquery.js"></script>
</head>
<body>
    <div id="log-container" style="height: 450px; overflow-y: scroll; background: #333; color: #aaa; padding: 10px;">
        <div>
        </div>
    </div>
</body>
<script>
    $(document).ready(function() {
        // 指定websocket路径
        var websocket = new WebSocket('ws://localhost:8080/read');
        websocket.onmessage = function(event) {
            // 接收服务端的实时日志并添加到HTML页面中
            $("#log-container div").append(event.data);
            // 滚动条滚动到最低部
            $("#log-container").scrollTop($("#log-container div").height() - $("#log-container").height());
        };
    });
</script>
</body>
</html>`
    fmt.Fprintf(w, body)
    return
}

func main() {
    fmt.Println("tester")
    http.HandleFunc("/", view)
    http.HandleFunc("/view", view)
    http.HandleFunc("/read", read)
    http.ListenAndServe(":8080", nil)
}

