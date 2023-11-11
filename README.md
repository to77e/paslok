# paslok
This is a useful and reliable tool for generating strong, secure passwords and storing them. 

Before using please install xclip
```shell
sudo apt install xclip
```

# Installation
```shell
    go install github.com/to77e/paslok/cmd/paslok
```

First you need to declare environment variables:   
PASLOK_CIPHER_KEY - AES 256 encryption key;   
PASLOK_DB_PATH - path to the database where the passwords will be stored (default: ~/.paslok/paslok.db).

Arguments:  
    
    -create [name comment] - generate password and save it  
    -read [name] - copy password to clipboard    
    -update [name new_password] - update password
    -delete [name] - delete password
    -list - list all names and comments
    -version - show version
    -help - show help


Example:
```shell
  paslok -create newname newcomment -read name
```

Please note that this application has only been tested on Linux and has not been tested on other platforms. Use on other platforms is not supported at this time.
