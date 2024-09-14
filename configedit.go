package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("用法: ./configedit <文件路径> <用户名> <密码> [端口]")
        fmt.Println("默认端口为8080")
        return
    }

    filePath := os.Args[1]
    username := os.Args[2]
    password := os.Args[3]
    port := "8080" // 默认端口

    // 如果有提供端口参数，则使用提供的端口
    if len(os.Args) == 5 {
        port = os.Args[4]
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // 检查用户是否已认证
        user, pass, ok := r.BasicAuth()
        if !ok || user != username || pass != password {
            w.Header().Set("WWW-Authenticate", `Basic realm="请提供凭证"`)
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintln(w, "需要认证")
            return
        }

        if r.Method == http.MethodPost {
            // 处理文件内容的更新
            r.ParseForm()
            newContent := r.FormValue("content")
            err := ioutil.WriteFile(filePath, []byte(newContent), 0644)
            if err != nil {
                http.Error(w, "无法写入文件", http.StatusInternalServerError)
                return
            }
            // 成功保存后，显示成功消息并设置跳转时间
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            fmt.Fprintf(w, `
                <!DOCTYPE html>
                <html lang="zh-CN">
                <head>
                    <meta charset="UTF-8">
                    <meta name="viewport" content="width=device-width, initial-scale=1.0">
                    <title>保存成功</title>
                    <script>
                        setTimeout(function() {
                            window.location.href = "/";
                        }, 3000); // 3秒后跳转
                    </script>
                </head>
                <body>
                    <h2>文件已成功更新！</h2>
                    <p>3秒后将自动跳转回编辑页面。</p>
                </body>
                </html>
            `)
            return
        }

        // 读取并显示文件内容
        content, err := ioutil.ReadFile(filePath)
        if err != nil {
            http.Error(w, "无法读取文件", http.StatusInternalServerError)
            return
        }

        // 显示文件内容和编辑表单
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprintf(w, `
            <!DOCTYPE html>
            <html lang="zh-CN">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>%s</title>
                <style>
                    body { font-family: Arial, sans-serif; margin: 20px; background-color: #f4f4f4; }
                    h2 { color: #333; }
                    form { background: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); }
                    textarea { width: 100%%; height: 300px; border: 1px solid #ccc; border-radius: 4px; padding: 10px; box-sizing: border-box; }
                    input[type="submit"] { 
                        background-color: #4CAF50; 
                        color: white; 
                        border: none; 
                        border-radius: 4px; 
                        padding: 10px 15px; 
                        cursor: pointer; 
                        margin-top: 10px; 
                    }
                    input[type="submit"]:hover { background-color: #45a049; }
                </style>
            </head>
            <body>
                <h2>编辑: %s</h2>
                <form method="post">
                    <textarea name="content">%s</textarea><br>
                    <input type="submit" value="保存修改">
                </form>
            </body>
            </html>
        `, filePath, filePath, string(content))
    })

    fmt.Printf("configedit 服务器正在运行，访问 http://localhost:%s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Println("服务器启动失败:", err)
    }
}
