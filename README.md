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

    -comment string
        Comment to add to the password entry (applicable with -create)
    -create string
        Create a new password entry in the locker. Usage: -create=<serviceName>. Optionally add -comment=<comment>
    -dash
        Include dash in the generated password (applicable with -create) (default true)
    -delete string
        Delete the password entry by service name. Usage: -delete=<serviceName>
    -length int
        Length of the generated password (applicable with -create) (default 18)
    -list
        List all stored services. Optionally filter results with -search
    -number
        Include numbers in the generated password (applicable with -create) (default true)
    -password string
        Use a custom password instead of generating one (applicable with -create)
    -read string
        Read the password for the specified service name. Usage: -read=<serviceName>/<id>
    -search string
        Substring to filter results (only applicable with -list)
    -special
        Include special characters in the generated password (applicable with -create) (default true)
    -uppercase
        Include uppercase characters in the generated password (applicable with -create) (default true)
    -username string
        Username to add to the password entry (applicable with -create)
    -version
        Show the current version of the application



Example:
```shell
  paslok -create service -username user -comment comment
```

Please note that this application has only been tested on Linux and has not been tested on other platforms. Use on other platforms is not supported at this time.
