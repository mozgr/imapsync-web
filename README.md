## Experimental web ui for imapsync(Linux) ##

1. You have to install imapsync to your linux server first.
2. Change xxx.xxx.xxx to your server ip address in index.html 
3. Put template/index.htm, .env, imapsync to same directory and run imapsync.

 chmod +x imapsync
 ./imapsync

You can change default username, password and lines display from .env
Don't forget to open port 3000 on your firewall.

![image](https://github.com/mozgr/imapsync-web/blob/master/screenshot.png)