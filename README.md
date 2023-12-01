## Simple web ui for imapsync(Linux) ##

With this web application, you can migrate your emails from multiple accounts to new accounts simultaneously.

Installation
1. You have to install imapsync to your linux server first.
2. Change xxx.xxx.xxx to your server ip address in index.html 
3. Put template/index.html, .env, imapsync-web to same directory and run imapsync-web.
```
chmod +x imapsync-web
./imapsync-web
```

You can change the default username, password, and line display settings in the .env file.<br/>
Don't forget to open port 3000 on your firewall.<br/>
Go backend uses websocket to print stdout to web page.
You can simply modify index.html and "command" variable in go code to run any linux command from your browser.

![image](https://github.com/mozgr/imapsync-web/blob/master/screenshot.png)
