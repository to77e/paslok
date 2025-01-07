package main

import (
	"flag"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/atotto/clipboard"

	"github.com/to77e/paslok/internal/boot"
	"github.com/to77e/paslok/internal/models"
	"github.com/to77e/paslok/internal/printer"
	"github.com/to77e/paslok/internal/validator"
)

var version string

func main() { //nolint:funlen
	var (
		flagCreate = flag.String("create", "",
			"Create a new password entry in the locker. Usage: -create=<serviceName>. Optionally add -comment=<comment>")
		flagRead = flag.String("read", "",
			"Read the password for the specified service name. Usage: -read=<serviceName>/<id>")
		flagDelete = flag.String("delete", "",
			"Delete the password entry by service name. Usage: -delete=<serviceName>")
		flagList = flag.Bool("list", false,
			"List all stored services. Optionally filter results with -search")
		flagVersion = flag.Bool("version", false,
			"Show the current version of the application")
		flagSearch = flag.String("search", "",
			"Substring to filter results (only applicable with -list)")
		flagComment = flag.String("comment", "",
			"Comment to add to the password entry (applicable with -create)")
		flagUsername = flag.String("username", "",
			"Username to add to the password entry (applicable with -create)")
		flagLength = flag.Int("length", 18,
			"Length of the generated password (applicable with -create)")
		flagUppercase = flag.Bool("uppercase", true,
			"Include uppercase characters in the generated password (applicable with -create)")
		flagSpecial = flag.Bool("special", true,
			"Include special characters in the generated password (applicable with -create)")
		flagNumber = flag.Bool("number", true,
			"Include numbers in the generated password (applicable with -create)")
		flagDash = flag.Bool("dash", true,
			"Include dash in the generated password (applicable with -create)")
		flagPassword = flag.String("password", "",
			"Use a custom password instead of generating one (applicable with -create)")
	)

	flag.Parse()

	app, err := boot.Initialize()
	if err != nil {
		fmt.Printf("error: initializing application: %s\n", err.Error())
		return
	}
	defer app.Database.Close() //nolint:errcheck

	switch {
	case *flagVersion:
		fmt.Printf("version: %s\n", version)

	case len(*flagCreate) > 0:
		req := &models.CreatePasswordRequest{
			Service:   *flagCreate,
			Comment:   *flagComment,
			Username:  *flagUsername,
			Length:    *flagLength,
			Uppercase: *flagUppercase,
			Special:   *flagSpecial,
			Number:    *flagNumber,
			Dash:      *flagDash,
			Password:  *flagPassword,
		}
		if err = validator.GetInstance().Struct(req); err != nil {
			fmt.Printf("error: validation: %s\n", err.Error())
			return
		}
		if err = app.LockerService.Create(req); err != nil {
			fmt.Printf("error: create password: %s\n", err.Error())
			return
		}
		fmt.Printf("password for %s created\n", req.Service)

	case len(*flagRead) > 0:
		req := &models.ReadPasswordRequest{}
		req.Id, err = strconv.ParseInt(*flagRead, 10, 64)
		if err != nil {
			req.Service = *flagRead
		}
		if err = validator.GetInstance().Struct(req); err != nil {
			fmt.Printf("error: validation: %s\n", err.Error())
			return
		}
		var pass string
		pass, err = app.LockerService.Read(req)
		if err != nil {
			fmt.Printf("error: read password: %s\n", err.Error())
			return
		}
		if err = clipboard.WriteAll(pass); err != nil {
			fmt.Printf("error: copy password to clipboard: %s\n", err.Error())
		}
		fmt.Printf("password copied to clipboard\n")

	case *flagList:
		req := &models.ListPasswordsRequest{
			SearchTerm: *flagSearch,
		}
		if err = validator.GetInstance().Struct(req); err != nil {
			fmt.Printf("error: validation: %s\n", err.Error())
			return
		}
		var resources []models.Resource
		resources, err = app.LockerService.List(req)
		if err != nil {
			fmt.Printf("error: list of names: %s\n", err.Error())
			return
		}
		fmt.Println(printer.PrintResources(resources))

	case len(*flagDelete) > 0:
		req := &models.DeletePasswordRequest{}
		req.Id, err = strconv.ParseInt(*flagDelete, 10, 64)
		if err != nil {
			req.Service = *flagDelete
		}
		if err = validator.GetInstance().Struct(req); err != nil {
			fmt.Printf("error: validation: %s\n", err.Error())
			return
		}
		if err = app.LockerService.Delete(req); err != nil {
			fmt.Printf("error: delete password: %s\n", err.Error())
			return
		}
		fmt.Printf("password deleted\n")
	}
}
