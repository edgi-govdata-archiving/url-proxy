# url-proxy
--
Url proxy, uh, proxies a url, maybe so you don't have to expose an API key or
something. To get up & running, make sure you have [go](http://golang.org)
installed, `cd` to the directory, run `go build`, then `ENDPOINT=[url to proxy]
./url-proxy`. It reads from the environemnt variable `ENDPOINT`, so you can set
that however you like.
