Contributing to MumbleDJ
========================

Contributions are always welcome to MumbleDJ. This document will give you some tips and guidelines to follow while implementing your contribution.

## Table of Contents
* [Implementing a new command](#implementing-a-new-command)
  * [Create files for your new command](#create-files-for-your-new-command)
  * [Copy templates into your new files](#copy-templates-into-your-new-files)
    * [Command implementation template (`command.go`)](#command-implementation-template-commandgo)
    * [Command test suite template (`command_test.go`)](#command-test-suite-template-command_testgo)
  * [Implement your new command](#implement-your-new-command)
  * [Add command to `commands/pkg_init.go`](#add-command-to-commandspkg_initgo)
  * [Add necessary configuration values to `config.yaml` and `config.go`](#add-necessary-configuration-values-to-configyaml-and-configgo)
  * [Regenerate `bindata.go`](#regenerate-bindatago)
  * [Document your new command](#document-your-new-command)
* [Implementing support for a new service](#implementing-support-for-a-new-service)
  * [Create file for your new service](#create-file-for-your-new-service)
  * [Copy template into your new file](#copy-template-into-your-new-file)
  * [Implement your new service](#implement-your-new-service)
  * [Add service to `services/pkg_init.go`](#add-command-to-servicespkg_initgo)
  * [Add API key configuration value to `config.yaml` and `config.go` if necessary](#add-api-key-configuration-value-to-configyaml-and-configgo-if-necessary)
  * [Document your new service](#document-your-new-service)

## Implementing a new command
Commands are the portion of MumbleDJ that allows users to interact with the bot. Here is a step-by-step guide on how to implement a new command:

### Create files for your new command
All commands possess their own `command.go` file and `command_test.go` file. These files must reside in the `commands` directory.

### Copy templates into your new files
Templates for both command implementations and command tests have been created to ensure consistency across the codebase. Please use these templates, they will make your implementation easier and will make the codebase much cleaner.

#### Command implementation template (`command.go`)
```go
/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/yournewcommand.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"

	"layeh.com/gumble/gumble"
	"reik.pl/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// YourNewCommand is a command... (put a short description of the command here)
type YourNewCommand struct{}

// Aliases returns the current aliases for the command.
func (c *YourNewCommand) Aliases() []string {
	return viper.GetStringSlice("commands.yournewcommand.aliases")
}

// Description returns the description for the command.
func (c *YourNewCommand) Description() string {
	return viper.GetString("commands.yournewcommand.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *YourNewCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.yournewcommand.is_admin")
}

// Execute executes the command with the given user and arguments.
// Return value descriptions:
//    string: A message to be returned to the user upon successful execution.
//    bool:   Whether the message should be private or not. true = private,
//            false = public (sent to whole channel).
//    error:  An error message to be returned upon unsuccessful execution.
//            If no error has occurred, pass nil instead.
// Example return statement:
//    return "This is a private message!", true, nil
func (c *YourNewCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	
}
```

#### Command test suite template (`command_test.go`)
```go
/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/yournewcommand_test.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"layeh.com/gumble/gumbleffmpeg"
	"reik.pl/mumbledj/bot"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type YourNewCommandTestSuite struct {
	Command YourNewCommand
	suite.Suite
}

func (suite *YourNewCommandTestSuite) SetupSuite() {
	DJ = bot.NewMumbleDJ()
	bot.DJ = DJ

	viper.Set("commands.yournewcommand.aliases", []string{"yournewcommand", "c"})
	viper.Set("commands.yournewcommand.description", "yournewcommand")
	viper.Set("commands.yournewcommand.is_admin", false)
}

func (suite *YourNewCommandTestSuite) SetupTest() {
	DJ.Queue = bot.NewQueue()
}

func (suite *YourNewCommandTestSuite) TestAliases() {
	suite.Equal([]string{"yournewcommand", "c"}, suite.Command.Aliases())
}

func (suite *YourNewCommandTestSuite) TestDescription() {
	suite.Equal("yournewcommand", suite.Command.Description())
}

func (suite *YourNewCommandTestSuite) TestIsAdminCommand() {
	suite.False(suite.Command.IsAdminCommand())
}

// Implement more tests here as necessary! It may be helpful to take a look
// at the stretchr/testify documentation:
// https://github.com/stretchr/testify
// Remove this comment before sending a pull request.

func TestYourNewCommandTestSuite(t *testing.T) {
	suite.Run(t, new(YourNewCommandTestSuite))
}
```

### Implement your new command
Now the fun starts! Write the implementation for your command in the `Execute()` method. Then, write tests for your new command, making sure to test each possible execution flow of your command.

For writing the implementation and unit tests for your new command, it may be helpful to [look at previously created commands](https://reik.pl/mumbledj/blob/master/commands).

**Make sure to rename the example names to represent your new command!**

### Add command to `commands/pkg_init.go`
`commands/pkg_init.go` contains a slice of enabled commands. If you do not put your command in this slice, your command will not be enabled.

**Please keep the commands in alphabetical order!**

### Add necessary configuration values to `config.yaml` and `config.go`
Go to `config.yaml` and `bot/config.go` and add the necessary configuration values for your new command. 

**Please keep the commands in alphabetical order!**

### Regenerate `bindata.go`
This step is very easy, but is very important. This allows the bot to store a copy of the new `config.yaml` internally and use it to write to disk.

Simply execute `make bindata` and this step will be taken care of!

### Document your new command
Make sure to put information in `README.md` about your new command. It would be a shame for your new command to go unnoticed!

**Please keep the commands in alphabetical order!**

## Implementing support for a new service
Services are the portion of MumbleDJ that allows the bot to interact with various media services. Here is a step-by-step guide on how to implement support for a new service:

### Create file for your new service
All services possess their own `service.go` file. This file must reside in the `services` directory.

### Copy template into your new file
A template for service implementations has been created to ensure consistency across the codebase. Please use this template, it will make your implementation easier and will make the codebase much cleaner.

```go
/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/yournewservice.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"layeh.com/gumble/gumble"
	"reik.pl/mumbledj/bot"
	"reik.pl/mumbledj/interfaces"
)

// YourNewService is a... (description here)
type YourNewService struct {
	*GenericService
}

// NewYourNewServiceService returns an initialized YourNewService service object.
func NewYourNewServiceService() *YourNewService {
	return &YourNewService{
		&GenericService{
			ReadableName: "Your new service",
			Format:       "bestaudio",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`regex for track URLs in your new service`),
			},
			PlaylistRegex: []*regexp.Regexp{
                regexp.MustCompile(`regex for playlist URLs in your new service`),   
            },
		},
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (yn *YourNewService) CheckAPIKey() error {
	
}

// GetTracks uses the passed URL to find and return
// tracks associated with the URL. An error is returned
// if any error occurs during the API call.
func (yn *YourNewService) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	
}
```

### Implement your new service
Now the fun starts! Implement `CheckAPIKey()` and `GetTracks()`.

For writing the implementation for your new service, it may be helpful to [look at previously created service wrappers](https://reik.pl/mumbledj/blob/master/services).

**Make sure to rename the example names to represent your new service!**

### Add service to `services/pkg_init.go`
`services/pkg_init.go` contains a slice of enabled services. If you do not put your service in this slice, your service will not be enabled.

**Please keep the services in alphabetical order!**

### Add API key configuration value to `config.yaml` and `config.go` if necessary
Some services will require an API key for users to interact with their service. If an API key is required for your service, add it to the configuration file and `config.go` and run `make bindata` to regenerate the `bindata.go` file.

### Document your new service
In sections of `README.md` that describe which services are supported, add your new service to the list. It would be a shame for your new service to go unnoticed!

Also, if your service requires an API key, make sure to document the steps to retrieve an API key in the "Requirements" section of the `README`.
