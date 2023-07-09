# paslok
This is a useful and reliable tool for generating strong, secure passwords and storing them. 

Before using please install xclip
```shell
sudo apt install xclip
```

# Installation
```go
go install github.com/to77e/paslok/cmd/paslok
```

First you need to declare environment variables:   
PASLOK_CIPHER_KEY - AES 256 encryption key;   
PASLOK_FILE_PATH - path to the file where the passwords will be stored (default: ~/.paslok/.paslok).

Arguments:  
    -l - list all names and comments  
    -r [name] - copy password to clipboard   
    -c [name comment] - generate password and save it

Example:
```shell
paslok -c newname newcomment -r name
```

Please note that this application has only been tested on Linux and has not been tested on other platforms. Use on other platforms is not supported at this time.
