#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import json
import urlparse
import base64
import BaseHTTPServer

# 修改下面的路径到python-sdk所在的路径
sys.path.append("path/to/qiniu/python-sdk")
from qiniu import rs
from qiniu import io
from qiniu import conf

BUCKET = "a"
DOMAIN = "aatest.qiniudn.com"
AKEY = "iN7NgwM31j4-BZacMjPrOQBs34UG1maYCAQmhdCV"
SKEY = "6QTOr2Jg1gcZEWDQXKOGZh5PziC2MCV5KsntT70j"

HOST_NAME = "localhost"
PORT_NUMBER = 8765

uploadFormFmt = """
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
"""

uploadWithKeyFormFmt = """
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
"""

returnPageFmt = """
<html>
 <body>
  <p>%s
  <p>ImageDownloadUrl: %s
  <p><a href="/upload">Back to upload</a>
  <p><a href="/upload2">Back to uploadWithKey</a>
  <p><img src="%s">
 </body>
</html>
"""

conf.ACCESS_KEY = AKEY
conf.SECRET_KEY = SKEY

class MyHandler(BaseHTTPServer.BaseHTTPRequestHandler):

	def do_GET(self):
		p = urlparse.urlparse(self.path)

		if p.path.strip("/") == "":
			self.handle_default()
		elif p.path.strip("/") == "upload":
			self.handle_upload()
		elif p.path.strip("/") == "upload2":
			self.handle_upload_with_key()
		elif p.path.strip("/") == "uploaded":
			self.handle_return()
		else:
			self.not_found()

	def not_found(self):
		self.make_response(404, "Page not found.")

	def handle_upload(self):
		policy = rs.PutPolicy(BUCKET)
		policy.endUser = "userId"
		policy.saveKey = "$(sha1)"
		policy.returnBody = '{"hash": $(etag), "key": $(key)}'
		policy.returnUrl = 'http://localhost:8765/uploaded'
		token = policy.token()
		print "token:", token
		uploadForm = uploadFormFmt % token

		self.make_response(200, uploadForm)

	def handle_upload_with_key(self):
		policy = rs.PutPolicy(BUCKET)
		policy.endUser = "userId"
		policy.returnUrl = "http://localhost:8765/uploaded"
		token = policy.token()
		print "token:", token
		uploadForm = uploadWithKeyFormFmt % token
		self.make_response(200, uploadForm)

	def handle_default(self):
		self.send_response(302)
		self.send_header("Location", "/upload")
		self.end_headers()

	def handle_return(self):
		p = urlparse.urlparse(self.path)
		qs = urlparse.parse_qs(p.query)
		if "error" in qs.keys():
			self.make_response(400, qs["error"][0])
			return
		if not "upload_ret" in qs.keys():
			self.make_response(400, "Invalid query string.")
			return
		retstr = qs["upload_ret"][0]
		if retstr == "":
			self.make_response(400, "Invalid query string.")
			return
		try:
			dec = base64.urlsafe_b64decode(retstr)
			ret = json.loads(dec)
		except Exception as e:
			self.make_response(404, "Invalid query string. Decode error.")
			return
		policy = rs.GetPolicy()
		img = policy.make_request(rs.make_base_url(DOMAIN, ret["key"]))
		returnPage = returnPageFmt % (dec, img, img)
		self.make_response(200, returnPage)

	def make_response(self, code, msg=None):
		self.send_response(code)
		self.send_header("Content-type", "text/html")
		self.end_headers()
		if msg is not None:
			self.wfile.write(msg)

if __name__ == "__main__":
	server_class = BaseHTTPServer.HTTPServer
	httpd = server_class((HOST_NAME, PORT_NUMBER), MyHandler)
	print "Server Starts - %s:%s" % (HOST_NAME, PORT_NUMBER)
	try:
		httpd.serve_forever()
	except KeyboardInterrupt:
		pass

	httpd.server_close()
	print "Server Stops - %s:%s" % (HOST_NAME, PORT_NUMBER)

