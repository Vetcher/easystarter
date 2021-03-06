# Easy Starter
Tool to manage your services

## Install
```
$ go get github.com/vetcher/easystarter
```

## Command line arguments
|Argument                         |Description                                                                                                     |
|:-------------------------------:|----------------------------------------------------------------------------------------------------------------|
|`-env <path-of-env.ini>`         | Path to file with environment variables. Default `env.ini`                                                     |
|`-config <path-of-serices.json>` | File with services configuration. Default `services.json`                                                      |
|`-s={true/false}`                | This flag means start all services after startup. Same as enter `start -all` after run program. Default `false`|

## Commands

| Title              | Command   | Description                                                                                                         | Parameters             | Other |
|:-------------------|:----------|:--------------------------------------------------------------------------------------------------------------------|:-----------------------|-------|
| Exit               | `exit`    | Exit program                                                                                                        |                        |       |
| Start service      | `start`   | Start specified services or start all                                                                               | `-all` or service name |       |
| Stop service       | `stop`    | Stop specified services or stop all (send `SIGTERM` signal)                                                         | `-all` or service name |       |
| Kill service       | `kill`    | Kill specified services or kill all                                                                                 | `-all` or service name |       |
| Restart service    | `restart` | Stop and start services.                                                                                            | `-all` or service name |       |
| List services      | `ps`      | Print all services, their args and status.                                                                          | `-all`                 |       |
| List environment   | `env`     | Print environment variables from `env.ini` file or all. With flag `-reload` reloads environment from `env.ini` file | `-all` or `-reload`    |       |

## Usage
1. Add `GOPATH` to your `PATH`.
2. Create in `HOME` folder file `services.json` with [configuration](#service-configuration).
3. Create in `HOME` folder file `env.ini` with [environment variables](#environment-configuration).
4. Type `easystarter` in terminal and press Enter.

Program creates `logs` folder in current directory if it does not exist yet.    
Logs for each service writes to `./logs/<service-name>.log` file.    
For services names app use some sort of auto-completion, so you can specify only beginning of service name.    

## Service configuration
You can specify services in file `services.json`, where you may set name, target Makefile with `install` _rule_, custom directory (absolute or relative) to service folder and command line arguments for service.    
If `services.json` not in current directory, program use file from `$HOME` folder.    
For file structure refer at `services.json` file in repository. __Fields `name` and `target` are required__.

#### Example
Example/template for `services.json` file.
```
[
  {
    "name": "testing",
    "target": "Makefile",
    "dir": "",
    "args": ["-duration", "3", "-x", "2"]
  },
  {
    "name": "testing1",
    "target": "path/to/makefile/inside/testing1/dir/Makefile",
    "dir": "/home/vetcher/bin",
    "args": ["-duration", "5", "-x", "10"]
  }
]
```

## Environment configuration
You can add custom environment variables to `env.ini`, which will be added before service execution.
It looks up for `env.ini` file in current directory. If there is no file, it open `env.ini` in `HOME` folder.

#### Example
Example/template for `env.ini` file.
```
GOPATH=/home/vetcher/go
PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```
