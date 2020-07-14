// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func SendMailInterface(w http.ResponseWriter, r *http.Request) {
	//アクセスを許可するドメインを設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//リクエストに使用可能なメソッド
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//Access-Control-Allow-Methods 及び Access-Control-Allow-Headers ヘッダーに含まれる情報をキャッシュすることができる時間の長さ(seconds)
	w.Header().Set("Access-Control-Max-Age", "86400")
	//リクエストに使用可能なHeaders
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept")
	// //レスポンスのContent-Typeを設定する
	// w.Header().Set("Content-Type", "application/json")

	// SORACOM UIより取得したJSONファイルを解読
	// タイトルの直接挿入
	var d interface{}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, fmt.Sprintf("%v", err))
		return
	}

	tomails := os.Getenv("TOMAILS")
	// ; で文字列を配列化
	mails := strings.Split(tomails, ";")

	for i := range mails {
		g := &Gmail{
			From:     os.Getenv(("SENDMAIL")),
			Username: os.Getenv("SENDMAIL"),
			Password: os.Getenv("SENDMAILPASSWORD"),
			To:       mails[i],
			Subject:  os.Getenv("MAILTITLE"),
			Message:  fmt.Sprintf("%v", d),
		}
		if err := g.Send(); err != nil {
			fmt.Fprint(w, fmt.Sprintf("%v", err))
			return
		}
	}

	return
}

type Gmail struct {
	From     string
	Username string
	// Password is password/appword
	Password string
	To       string
	Subject  string
	Message  string
}

func (g Gmail) Endpoint() string {
	return "smtp.gmail.com"
}

func (g Gmail) Send() error {
	auth := smtp.PlainAuth("", g.Username, g.Password, g.Endpoint())
	if err := smtp.SendMail(g.Endpoint()+":587", auth, g.From, []string{g.To}, g.body()); err != nil {
		return err
	}

	return nil
}

func (g Gmail) body() []byte {
	return []byte("To: " + g.To + "\r\n" +
		"Subject: " + g.Subject + "\r\n\r\n" +
		g.MaxChars() + "\r\n")
}

func (g Gmail) MaxChars() string {
	return g.Message
}
