# nomad-demo

## Build 🔨 and run 🏃
### Nomad
#### 1. Install Nomad https://developer.hashicorp.com/nomad/tutorials/get-started/gs-install
#### 2. Run Nomad cluster
```
sudo nomad agent -dev \
  -bind 0.0.0.0 \
  -network-interface='{{ GetDefaultInterfaces | attr "name" }}'
```
#### 3. Run consul
```
consul agent -dev
```
### Golang API
Run Golang REST API with an endpoint PUT /services/\<name>, which allows creating or updating a service with the given name. It accepts the following parameters:
- script (boolean)
- url (string)

You set server port through PORT env variable (default is 3000)

```
make run
```

## Tests
Example with script set to false:
- Create a new entry on pastebin.com with the content "Hello world, this is the
content of my webpage"
- Create a new service using your API: curl localhost:3000/services/mypage -X PUT -d
'{"url": "https://pastebin.com/raw/hEFbnx33", "script": false}'
- The API returns:
```
{
"url": "http://192.168.1.104:20676"
}
```
- The content of the pastebin is downloaded and served with Nginx. A request to
http://192.168.1.104:20676 returns "Hello world, this is the content of my webpage"

Example with script set to true:
- Create a new entry on pastebin.com with the content:
```
#!/bin/sh
echo "Hello world!"
```
- Create a new service using your API: curl localhost:3000/services/mypage -X PUT -d
'{"url": "https://pastebin.com/raw/abcde123", "script": true}'
- The API returns:
```
{
"url": "http://192.168.1.104:22733"
}
```
- The content of the pastebin is executed as a script, and accessing the URL at
http://192.168.1.104:22733 returns "Hello world!"

## Unit Tests
```
make unit_test
```
