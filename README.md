Dans Ton Cache
==============

Using [golang-LRU](github.com/hashicorp/golang-lru) as a HTTP cache.

Build it
--------

    make bin

Try it
------

Launch [imgproxy](https://github.com/imgproxy/imgproxy)

    IMGPROXY_ENABLE_WEBP_DETECTION=true ./imgproxy

Launch caching proxy

    LISTEN=:8000 BACKEND=http://localhost:8080  ./bin/dtc-proxy

Try official test url : http://localhost:8000/insecure/fill/300/400/sm/0/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk1tUTNabVk0/TnpZdFkyVm1ZaTAw/WkRSbUxUZ3lPREF0/WldZelpqaGxOemsx/TnpVMlhrRXlYa0Zx/Y0dkZVFYVnlOVGMz/TWpVek5USUAuanBn

Watch for `X-Cache` header. Modern browser accept webp, and imgproxy handle it.

Try it with curl

```
$ curl -v http://localhost:8000/insecure/fill/300/400/sm/0/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk1tUTNabVk0/TnpZdFkyVm1ZaTAw/WkRSbUxUZ3lPREF0/WldZelpqaGxOemsx/TnpVMlhrRXlYa0Zx/Y0dkZVFYVnlOVGMz/TWpVek5USUAuanBn > /dev/null

*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8000 (#0)
> GET /insecure/fill/300/400/sm/0/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk1tUTNabVk0/TnpZdFkyVm1ZaTAw/WkRSbUxUZ3lPREF0/WldZelpqaGxOemsx/TnpVMlhrRXlYa0Zx/Y0dkZVFYVnlOVGMz/TWpVek5USUAuanBn HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Cache-Control: max-age=3600, public
< Content-Disposition: inline; filename="MV5BMmQ3ZmY4NzYtY2VmYi00ZDRmLTgyODAtZWYzZjhlNzk1NzU2XkEyXkFqcGdeQXVyNTc3MjUzNTI@.jpg"
< Content-Length: 19821
< Content-Type: image/jpeg
< Date: Mon, 23 Nov 2020 22:25:19 GMT
< Expires: Tue, 24 Nov 2020 00:25:19 GMT
< Server: imgproxy
< Vary: Accept
< X-Cache: hit
< X-Request-Id: 6U2C-TTfcLYfahCWsxmb4
<
{ [16384 bytes data]
* Connection #0 to host localhost left intact
```

Watch the cache guts

```
$ file /tmp/proxy/5a1d62976a02cc35689b3321ed331def714b35d912fd3462f8cda50d9e2f257b
/tmp/proxy/5a1d62976a02cc35689b3321ed331def714b35d912fd3462f8cda50d9e2f257b: JPEG image data, baseline, precision 8, 300x400, frames 3

$ cat /tmp/proxy/5a1d62976a02cc35689b3321ed331def714b35d912fd3462f8cda50d9e2f257b.header
Cache-Control: max-age=3600, public
Content-Disposition: inline; filename="MV5BMmQ3ZmY4NzYtY2VmYi00ZDRmLTgyODAtZWYzZjhlNzk1NzU2XkEyXkFqcGdeQXVyNTc3MjUzNTI@.jpg"
Content-Length: 19821
Content-Type: image/jpeg
Date: Mon, 23 Nov 2020 22:25:19 GMT
Expires: Tue, 24 Nov 2020 00:25:19 GMT
Server: imgproxy
Vary: Accept
X-Request-Id: 6U2C-TTfcLYfahCWsxmb4

$ file /tmp/proxy/ad19087e050ee8e0cbb08c07f11ba12a1c5852727a32ab9ee5d071825c75e073
/tmp/proxy/ad19087e050ee8e0cbb08c07f11ba12a1c5852727a32ab9ee5d071825c75e073: RIFF (little-endian) data, Web/P image, VP8 encoding, 300x400, Scaling: [none]x[none], YUV color, decoders should clamp

$ cat /tmp/proxy/ad19087e050ee8e0cbb08c07f11ba12a1c5852727a32ab9ee5d071825c75e073.header
Cache-Control: max-age=3600, public
Content-Disposition: inline; filename="MV5BMmQ3ZmY4NzYtY2VmYi00ZDRmLTgyODAtZWYzZjhlNzk1NzU2XkEyXkFqcGdeQXVyNTc3MjUzNTI@.webp"
Content-Length: 14804
Content-Type: image/webp
Date: Mon, 23 Nov 2020 22:25:43 GMT
Expires: Tue, 24 Nov 2020 00:25:43 GMT
Server: imgproxy
Vary: Accept
X-Request-Id: IXQCOZ0Uzii0uBF3kKGBd

```
