BingGo
=======

Download the [bing](http://www.bing.com/) wallpaper and set it wallpaper.
This application currently supports macOS and Gnome

Installation
-------------

`go get -u github.com/kmtr/binggo`

## Gnome

`gsettings` is required.

Usage
------

`binggo --pictdir /path/to/downloads`

If you want to change wallpaper of 2nd display, you will add "--display 2" (macOS).

`binggo --pictdir /path/to/downloads --display 2`

crontab
--------

```sh
(crontab -l; echo "0 0 * * * /path/to/binggo --pictdir /path/to/downloads >/dev/null 2>&1") | crontab -
```

