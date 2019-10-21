[![Build Status](https://travis-ci.org/mhewedy/onesignal-client.svg?branch=master)](https://travis-ci.org/mhewedy/onesignal-client)

### Usage

1. Create file with name `playerids.txt` contains player IDs separated by new lines.

2. Issue the command:
```shell script
APP_ID=<onesignal api key> ./onesignal-client<.exe|.mac|.linux> <tag list>
```
example to run on mac:
```shell script
APP_ID=3e64a763-0557-4094-9bbb-345fded09cd9 ./onesignal-client.mac '{"isApproved": true}'
```

example to run on linux:
```shell script
APP_ID=3e64a763-0557-4094-9bbb-345fded09cd9 ./onesignal-client.linux '{"isApproved": "", "isVerfied": true}'
```

see [onesignal docs](https://documentation.onesignal.com/reference#edit-device) for more about tags

### Downloads:
[Windows](https://github.com/mhewedy/onesignal-client/releases/download/v3.0/onesignal-client.exe)
|[Mac](https://github.com/mhewedy/onesignal-client/releases/download/v3.0/onesignal-client.mac)
|[Linux](https://github.com/mhewedy/onesignal-client/releases/download/v3.0/onesignal-client.linux)
