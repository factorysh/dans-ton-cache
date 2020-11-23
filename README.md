Dans Ton Cache
==============

Using [golang-LRU](github.com/hashicorp/golang-lru) as a HTTP cache.

Build it
--------

    make bin

Try it
------

Launch [imgproxy](https://github.com/imgproxy/imgproxy)

Launch caching proxy

    LISTEN=:8000 BACKEND=http://localhost:8080  ./bin/dtc-proxy

Try official test url : http://localhost:8000/insecure/fill/300/400/sm/0/aHR0cHM6Ly9tLm1l/ZGlhLWFtYXpvbi5j/b20vaW1hZ2VzL00v/TVY1Qk1tUTNabVk0/TnpZdFkyVm1ZaTAw/WkRSbUxUZ3lPREF0/WldZelpqaGxOemsx/TnpVMlhrRXlYa0Zx/Y0dkZVFYVnlOVGMz/TWpVek5USUAuanBn.jpg

Watch for `X-Cache` header.

Watch the cache guts

```
$ file /tmp/4f02bf4aacb20c48f61e2c3ab4e5dc6856679b295133795165b3c0b89383656c
/tmp/4f02bf4aacb20c48f61e2c3ab4e5dc6856679b295133795165b3c0b89383656c: JPEG image data, baseline, precision 8, 300x400, frames 3

$ cat /tmp/4f02bf4aacb20c48f61e2c3ab4e5dc6856679b295133795165b3c0b89383656c.header
Cache-Control: max-age=3600, public
Content-Disposition: inline; filename="MV5BMmQ3ZmY4NzYtY2VmYi00ZDRmLTgyODAtZWYzZjhlNzk1NzU2XkEyXkFqcGdeQXVyNTc3MjUzNTI@.jpg"
Content-Length: 19821
Content-Type: image/jpeg
Date: Mon, 23 Nov 2020 21:35:56 GMT
Expires: Mon, 23 Nov 2020 23:35:56 GMT
Server: imgproxy
X-Request-Id: 7LdmdoHjpaW7TB_fyy2n5
```
