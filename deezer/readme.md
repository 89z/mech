# Deezer

> Listen, um, I don’t really know what’s going on here, but I’m leaving. I
> don’t know where exactly, but I’m gonna start over.
>
> Come with me. Please.
>
> [Paint it Black (2016)][1]

Download audio from Deezer

[![][3]][2]

## ARL cookie

1. http://www.deezer.com/login
2. Web Developer
3. Network
4. Log in
5. Network Settings
6. Save All As HAR

## CDN

Note that two CDN are available, but the HTTP one seems to be faster, perhaps
because HTTPS overhead is not used:

- http://e-cdn-proxy-0.deezer.com
- https://e-cdns-proxy-0.dzcdn.net

## SID cookie

After some testing, this cookie seems to expire about three hours from last use.
So if you use it once every three hours, it should stay alive, but I would need
to do more testing to know for sure.

## Repos

language   | link
-----------|-----
C          | https://github.com/yne/dzr
Go         | https://github.com/godeezer/lib
Go         | https://github.com/moon004/Go-deezer-downloader
JavaScript | https://github.com/svbnet/diezel

[1]:https://f002.backblazeb2.com/file/ql8mlh/Paint.It.Black.2016.mp4
[2]:https://pkg.go.dev/github.com/89z/deezer
[3]:https://pkg.go.dev/badge/github.com/89z/deezer.svg
