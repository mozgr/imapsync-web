## Experimental web ui for imapsync(Linux) ##

1. You have to install imapsync to your linux server first.
2. Change xxx.xxx.xxx to your server ip address in index.html 
3. Put template/index.htm, .env, imapsync-web to same directory and run imapsync.
```
chmod +x imapsync-web
./imapsync-web
```

You can change the default username, password, and line display settings in the .env file.<br/>
Don't forget to open port 3000 on your firewall.

![image](https://github.com/mozgr/imapsync-web/blob/master/screenshot.png)