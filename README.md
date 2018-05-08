## gocheck2

Gocheck2 is my first prgoram in the Go language. It was orginally a bash script, however after reading a article about why bash scripts are horrible (not really) I decided to rewrite it. I started on python but actually did Go because of how fast and easy it is to deploy.

## What does it do?

Gocheck2 reads the .XST file that a Xerox printer leaves behind during a network scan. You're able to set custom entrys on scan screen but not able to rename based on entrys alone. This will auto detect when XST is uploaded and will move the .tif/.pdf to a processed folder.

```
./check2.exe
```
It can also be ran as a daemon.
```
./check2.exe d=true
```