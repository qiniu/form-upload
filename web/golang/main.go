package main

import (
	"log"
	"fmt"
	"os"
	"time"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
)

const (
	BUCKET = "a"
	DOMAIN = "aatest.qiniudn.com"
)

// --------------------------------------------------------------------------------

type PutPolicy struct {
	Scope            string `json:"scope,omitempty"`
	CallbackUrl      string `json:"callbackUrl,omitempty"`
	CallbackBody     string `json:"callbackBody,omitempty"`
	ReturnUrl        string `json:"returnUrl,omitempty"`
	ReturnBody       string `json:"returnBody,omitempty"`
	AsyncOps         string `json:"asyncOps,omitempty"`
	EndUser          string `json:"endUser,omitempty"`
	Expires          uint32 `json:"deadline"` 			// 截止时间（以秒为单位）
}

func (r PutPolicy) Token() string {
	if r.Expires == 0 {
		r.Expires = 3600
	}
	r.Expires += uint32(time.Now().Unix())
	return digest.SignJson(ACCESS_KEY, []byte(SECRET_KEY), &r)
}

// --------------------------------------------------------------------------------

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
}

// --------------------------------------------------------------------------------

var uploadFormFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   Image to upload: <input name="file" type="file"/>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`

var returnPageFmt = `
<html>
 <body>
  <p>%s
  <p><a href="/upload">Back to upload</a>
  <p><img src="%s">
 </body>
</html>
`

type UploadRet struct {
	Key string `json:"key"`
}

func handleReturn(w http.ResponseWriter, req *http.Request) {

	ret := req.FormValue("upload_ret")
	b, err := base64.URLEncoding.DecodeString(ret)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	var uploadRet UploadRet
	err = json.Unmarshal(b, &uploadRet)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	policy := rs.GetPolicy{Scope: "*/" + uploadRet.Key}
	token := policy.Token()
	img := "http://" + DOMAIN + "/" + uploadRet.Key + "?token=" + token
	returnPage := fmt.Sprintf(returnPageFmt, string(b), img)
	w.Write([]byte(returnPage))
}

func handleUpload(w http.ResponseWriter, req *http.Request) {

	policy := PutPolicy{Scope:BUCKET, ReturnUrl: "http://localhost:8765/uploaded"}
	token := policy.Token()
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadFormFmt, token)
	w.Write([]byte(uploadForm))	
}

func handleDefault(w http.ResponseWriter, req *http.Request) {

	http.Redirect(w, req, "/upload", 302)
}

func main() {

	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/uploaded", handleReturn)
	log.Fatal(http.ListenAndServe(":8765", nil))
}

// --------------------------------------------------------------------------------

