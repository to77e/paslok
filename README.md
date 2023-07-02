# paslok
This is a useful and reliable tool for generating strong, secure passwords and storing them. 

Before using please install xclip
```shell
sudo apt install xclip
```

First you need to create a "internal/config/config.yaml" file in which to declare:   
cipherKey - AES 256 encryption key;   
filePath - path to the file where the passwords will be stored.   

Arguments:  
    -l - list all names and comments  
    -r [name] - copy password to clipboard   
    -c [name comment] - generate password and save it

Example:
```shell
go run cmd/main.go -c newname newcomment -r name
```

Please note that this application has only been tested on Linux and has not been tested on other platforms. Use on other platforms is not supported at this time.
