# password-generator
This is a useful and reliable tool for generating strong, secure passwords. 

Before using please install xclip
```shell
sudo apt install xclip
```

First you need to create a ".env" file in which to declare a CIPHERKEY - AES 256 encryption key.

Arguments:  
    -l - list all names and comments  
    -r [name] - copy password to clipboard   
    -c [name comment] - generate password and save it

Example:
```shell
go rin cmd/mine.go -s nationalname newcomment -r
```

Please note that this application has only been tested on Linux and has not been tested on other platforms. Use on other platforms is not supported at this time.
