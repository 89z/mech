# Mech

> How’d you get this information?
>
> Just comes to you. This stuff just flies through the air. They send
> information out, and it’s beamed out all over the place. All you have to do
> is know how to grab it. See, I know how to grab it.
>
> [Heat (1995)](//f002.backblazeb2.com/file/ql8mlh/Heat.1995.mp4)

Anonymous API access

Some users might want to make anonymous requests, because of privacy or any
number of other reasons. This module allows people to do that. Most API these
days only offically support authenticated access. This is useful for the
company providing the API, as they can use the data for their own purposes
(analytics etc). However authentication really doesnt do anything for the end
user. Its just a pointless burden to getting the information you need for a
program you may be writing. Consider that in many cases, the same information
is available via HTML on the primary website, usually without being logged in.
So why can you do that with HTML, but not with the API? Well you can, using this
module.

https://godocs.io/github.com/89z/mech

## Sites

- Bandcamp
- GitHub
- Google Play
- Instagram
- MusicBrainz
- Reddit
- SoundCloud
- Vimeo
- YouTube

## Deezer

I have an implementation here:

https://github.com/89z/mech/tree/9dadd39c

However I have removed it for now, as I am busy with other stuff.

## Android

To setup MitmProxy, first download [1]. Then get IP address:

~~~
Get-NetIPAddress
~~~

Will look like this:

~~~
IPAddress         : 192.168.0.4
InterfaceIndex    : 11
InterfaceAlias    : Ethernet
~~~

Then start MitmProxy. Port will be in the bottom right corner, should be `8080`.
Then get Android Studio [2]. Click More Actions, AVD Manager, Create Virtual
Device. Use the default device defintion, then click Next. Click "x86 Images"
and download:

Release Name | API Level | ABI | Target
-------------|-----------|-----|------------
Nougat       | 25        | x86 | Google APIs

Launch the Android Emulator. Open Extended Controls by clicking "more". Click
settings, Proxy. Uncheck "Use Android Studio HTTP proxy settings". Click "Manual
proxy configuration", then enter IP address from above as "Host name", and port
from above as "Port number". Click Apply, you should see Proxy status Success.
Close Extended Controls. Open Chrome and browse to <https://mitm.it>. Click on
the Android certificate. Certificate name MITM.

Then if need be, you can download APKs [3]. Drag APK to device home screen to
install.

1. https://mitmproxy.org/downloads
2. https://developer.android.com/studio#downloads
3. https://apkpure.com

## Author

Steven Penny
