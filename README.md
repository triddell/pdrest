# What is pdrest?

__pdrest__ is a utility for IBM Security Access Manager (SAM) appliances. It allows batch processing of multiple "pdadmin" commands through the appliance's REST interface. The "old" pdadmin interface still exists on the SAM appliances when connecting through SSH, however, only one command can be entered at a time (easily). __pdrest__ allows an external file of commands to be feed into the appliance. Additionally, comments and blank lines are supported in the command file.

__pdrest__ is a stand-alone binary that allows for convenient installation. It specifically deals with the pdadmin interface and does not have additional functions for calling other SAM-provided REST calls.

__pdrest__ has been successfully used to run more than 30,000 commands at a time.

Installation
------------

__pdrest__ binaries can be downloaded for specific architectures from  [releases](https://github.com/triddell/pdrest/releases) or can be built from the go source.

To clone and compile a binary from source:

```bash
$ git clone https://github.com/triddell/pdrest.git $GOPATH/src/github.com/triddell/pdrest
$ cd $GOPATH/src/github.com/triddell/pdrest
$ go install
```

The `go install` command creates the binary at `$GOBIN/pdrest`.

Alternatively, you can build a binary and copy it wherever you'd like with `go build`:

```bash
$ git clone https://github.com/triddell/pdrest.git $GOPATH/src/github.com/triddell/pdrest
$ cd $GOPATH/src/github.com/triddell/pdrest
$ go build pdrest.go
$ mv pdrest /to/some/directory
```

##Usage

__pdrest -help__

```
> pdrest -help
Usage of pdrest:
  -applianceAdmin="admin": Appliance Admin User Name
  -commands="commands.txt": Commands File Path
  -creds="": Credentials File Path
  -host="localhost": Appliance Hostname
  -runtimeAdmin="sec_master": Runtime Admin User Name
```

The command defaults to using an appliance admin user name of `admin`, a runtime admin user name of `sec_master`, and a commands file name of `commands.txt`.

The `-hostname` flag will always be needed since __pdrest__ can't run on the appliance itself.

Passwords can be provided in an external file by using the `-credentials` flag.

###Sample files for example commands:

__commands.txt__

```
server list
# The following line is indented with a tab and this line is a comment.
  server list
```

__credentials.txt__

```
# SAM Appliance Password
appPswd
# SAM Runtime Password
samPswd
```

### Examples

__Example 1:__ Prompting for passwords

```
> pdrest -host sam01.riddell.us
Appliance Admin User Password:
Runtime Admin User Password:
Status Code:  200
Result:
cmd> server list
    ivmgrd-master
    web01-webseald-sam01.riddell.us
cmd> server list
    ivmgrd-master
    web01-webseald-sam01.riddell.us
```

__Example 2:__ Passwords stored in external file

```
> pdrest -host sam01.riddell.us -creds credentials.txt
Status Code:  200
Result:
cmd> server list
    ivmgrd-master
    web01-webseald-sam01.riddell.us
cmd> server list
    ivmgrd-master
    web01-webseald-sam01.riddell.us
```
