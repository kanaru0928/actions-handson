package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
    dbPath := os.Getenv("SQLITE_DB_PATH")
    if dbPath == "" {
        dbPath = "./data/app.db"
    }
    os.MkdirAll("./data", 0755)

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, content TEXT)`)
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        fmt.Fprintln(w, `<!DOCTYPE html><html lang='ja'><head><meta charset='utf-8'><title>掲示板</title><style>
        :root { 
          --primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          --accent: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
          --surface: rgba(255, 255, 255, 0.95);
          --glass: rgba(255, 255, 255, 0.1);
          --text: #2d3748;
          --border: rgba(255, 255, 255, 0.2);
          --shadow: 0 8px 32px rgba(31, 38, 135, 0.37);
        }
        body { 
          font-family: 'Hiragino Kaku Gothic ProN', 'Noto Sans JP', sans-serif; 
          background: var(--primary);
          min-height: 100vh;
          margin: 0;
          padding: 20px;
        }
        .container { 
          max-width: 600px; 
          margin: 0 auto; 
          background: var(--surface);
          border-radius: 20px; 
          box-shadow: var(--shadow);
          padding: 32px;
          backdrop-filter: blur(10px);
          border: 1px solid var(--border);
        }
        h1 { 
          text-align: center; 
          background: var(--accent);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          font-size: 2.5rem;
          font-weight: 700;
          margin-bottom: 32px;
          text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
        }
        ul { padding: 0; margin: 0; }
        li { 
          list-style: none; 
          background: var(--glass);
          backdrop-filter: blur(10px);
          border-radius: 12px;
          padding: 16px;
          margin-bottom: 12px;
          border: 1px solid var(--border);
          transition: all 0.3s ease;
        }
        li:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 16px rgba(0,0,0,0.1);
        }
        form { 
          display: flex; 
          gap: 12px; 
          margin-bottom: 24px; 
        }
        input[type=text] { 
          flex: 1; 
          padding: 16px; 
          border: 1px solid var(--border);
          border-radius: 12px; 
          font-size: 16px;
          background: var(--glass);
          backdrop-filter: blur(10px);
          transition: all 0.3s ease;
        }
        input[type=text]:focus {
          outline: none;
          border-color: #667eea;
          box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        button { 
          padding: 16px 24px; 
          border: none; 
          background: var(--accent);
          color: #fff; 
          border-radius: 12px; 
          cursor: pointer;
          font-weight: 600;
          font-size: 16px;
          transition: all 0.3s ease;
          box-shadow: 0 4px 16px rgba(245, 87, 108, 0.3);
        }
        button:hover { 
          transform: translateY(-2px);
          box-shadow: 0 8px 24px rgba(245, 87, 108, 0.4);
        }
        .message-id {
          background: var(--primary);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          font-weight: 600;
        }
        </style></head><body><div class='container'>`)
        fmt.Fprintln(w, `<h1>🌸 掲示板 v2.0 🌸</h1>`)
        fmt.Fprintln(w, `<form method='POST' action='/add'><input type='text' name='msg' placeholder='メッセージを入力'><button type='submit'>投稿</button></form>`)
        fmt.Fprintln(w, `<ul>`)
        rows, err := db.Query("SELECT id, content FROM messages ORDER BY id DESC")
        if err != nil {
            fmt.Fprintf(w, "<li style='color:red;'>%s</li>", err.Error())
        } else {
            defer rows.Close()
            for rows.Next() {
                var id int
                var content string
                rows.Scan(&id, &content)
                fmt.Fprintf(w, "<li><span class='message-id'>#%d</span>: %s</li>", id, content)
            }
        }
        fmt.Fprintln(w, `</ul></div></body></html>`)
    })

    http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        var msg string
        if r.Method == "POST" {
            r.ParseForm()
            msg = r.FormValue("msg")
        } else {
            msg = r.URL.Query().Get("msg")
        }
        if msg == "" {
            http.Error(w, "msg required", 400)
            return
        }
        _, err := db.Exec("INSERT INTO messages(content) VALUES(?)", msg)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("Listening on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
