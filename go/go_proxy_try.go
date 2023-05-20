package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest(goproxy.UrlIs("/api/v1")).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusForbidden,
				"{\"data\":\"ok\",\"msg\":\"\"}")
		})
	proxy.OnRequest(goproxy.Not(goproxy.SrcIpIs("192.168.1.1"))).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusForbidden,
				"DENY")
		})
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
