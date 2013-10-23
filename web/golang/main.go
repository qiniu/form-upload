package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

const (
	BUCKET = "a"
	DOMAIN = "aatest.qiniudn.com"
	AKEY = "iN7NgwM31j4-BZacMjPrOQBs34UG1maYCAQmhdCV"
	SKEY = "6QTOr2Jg1gcZEWDQXKOGZh5PziC2MCV5KsntT70j"
)

// --------------------------------------------------------------------------------

func init() {
	ACCESS_KEY = AKEY
	SECRET_KEY = SKEY
}

// --------------------------------------------------------------------------------

var uploadFormFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   Album belonged to: <input name="x:album" value="albumId"><br>
   Image to upload: <input name="file" type="file"/><br>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`

var uploadWithKeyFormFmt = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   Album belonged to: <input name="x:album" value="albumId"><br>
   Image key in qiniu cloud storage: <input name="key" value="foo bar.jpg"><br>
   Image to upload: <input name="file" type="file"/><br>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`

var returnPageFmt = `
<html>
 <body>
  <p>%s
  <p>ImageDownloadUrl: %s
  <p><a href="/upload">Back to upload</a>
  <p><a href="/upload2">Back to uploadWithKey</a>
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

	policy := rs.GetPolicy{}
	img := policy.MakeRequest(rs.MakeBaseUrl(DOMAIN, uploadRet.Key), nil)
	returnPage := fmt.Sprintf(returnPageFmt, string(b), img, img)
	w.Write([]byte(returnPage))
}

func handleUpload(w http.ResponseWriter, req *http.Request) {

	policy := rs.PutPolicy{
		Scope: BUCKET,
		ReturnUrl: "http://localhost:8765/uploaded",
		EndUser: "userId",
		SaveKey: "$(sha1)",
		ReturnBody: `{"hash": $(etag), "key": $(key)}`,
	}
	token := policy.Token(nil)
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadFormFmt, token)
	w.Write([]byte(uploadForm))
}

func handleUploadWithKey(w http.ResponseWriter, req *http.Request) {

	policy := rs.PutPolicy{Scope: BUCKET, ReturnUrl: "http://localhost:8765/uploaded", EndUser: "userId"}
	token := policy.Token(nil)
	log.Println("token:", token)
	uploadForm := fmt.Sprintf(uploadWithKeyFormFmt, token)
	w.Write([]byte(uploadForm))
}

func handleDefault(w http.ResponseWriter, req *http.Request) {

	http.Redirect(w, req, "/upload", 302)
}

func main() {

	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/upload2", handleUploadWithKey)
	http.HandleFunc("/uploaded", handleReturn)
	log.Fatal(http.ListenAndServe(":8765", nil))
}

// --------------------------------------------------------------------------------

