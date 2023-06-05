<h1 align="center">
  <img width="334" alt="gig" src="https://github.com/awgst/gig/assets/60491935/2cb202e0-c2ca-48d7-94ad-a47c4292b27d">

</h1>
<p align="center">🚀 Command line interface which provides a number helpful command to assist Go project development!</p>

## ⚡️ Getting started

First, [download](https://golang.org/dl/) and install **Go**. Then install gig by using [`go install`](https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies) command
```bash
go install github.com/awgst/gig/cmd/gig@latest
```

Check if gig is installed successfully
```bash
gig
```

Create a new project
```bash
gig create <project-name>
```
and it will create a new project with the following structure
```
.
├── cmd
│   └── seeder
│   │   └── main.go
│   └── main.go
├── database
│   └── migrations
│   │   └── 20210101000000_create_users_table.go
│   └── seeder
│   │   └── seeder.go
├── pkg
│   └── common
│   └── database
│   └── env
│   └── router
│   └── ...
├── src
│   └── app
│   └── config
│   └── routes
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── gig.json
├── go.mod
├── go.sum
```

Set up the `.env` file based on the `.env.example` file, then you can run the project
```bash
go run cmd/main.go
```

or you can use the `up` command if you want to run the project with docker. It will run docker-compose up command
```bash
gig up
```
**Notes**: make sure you have docker installed on your machine

create your first module
```bash
gig make:module <name>
```
It will create a new module inside the `src/app` directory with the following structure
```
.
├── http
│   └── handler
├── model
├── repository
├── service
├── module-name.go
```

***Handler** is a place where you can put your handler function.\
***Model** is a place where you can put your model struct.\
***Repository** is a place where you can put your repository function.\
***Service** is a place where you can put your service function.\
***module-name.go** is a place where you will create instance of handler, service and repository of your module.

If you prefered to use GORM on the first installation (you can disable it later on `gig.json` configuration), please to uncomment the following line in the `./pkg/database/gorm.go` file based on your database
```go
// Dialector for open connection based on driver
// Uncomment the code based on your database driver
// Run go get gorm.io/driver/{driver_name} to install the driver
// Example: go get gorm.io/driver/mysql and import "gorm.io/driver/mysql"
var dialectors = func(driver, dsn string) gorm.Dialector {
	return map[string]gorm.Dialector{
		// "mysql":      mysql.Open(dsn),
		// "postgres": postgres.Open(dsn),
	}[driver]
}
```

## 📚 Usage

# create
Create a new project
```bash
gig create <name>
```

Arguments:
| Argument      | Description                       | Type     | Required? |
|---------------|-----------------------------------|----------|-----------|
| `name`        | The name of the project           | `string` | Yes       |


Flags:
| Option        | Description                       | Type     | Default | Required? |
|---------------|-----------------------------------|----------|---------|-----------|
| `--version`   | Set the specific Go version.      | `string` | `1.20`  | No        |

# make
Make file
```bash
gig make:<type> <module_name>
```

Arguments:
| Argument             | Description                       | Type     | Required? |
|----------------------|-----------------------------------|----------|-----------|
| `module_name`        | The name of the project           | `string` | Yes       |

Types:
| Type        | Description                       |
|-------------|-----------------------------------|
| `migration` | Create a new migration            |
| `module`    | Create a new module               |
| `handler`   | Create a new handler              |
| `model`     | Create a new model                |
| `repository`| Create a new repository           |
| `service`   | Create a new service              |
| `request`   | Create a new request              |
| `response`  | Create a new response             |

# up
Run docker-compose up command
```bash
gig up
```

# migrate
Run a simple database migration. Under the hood it used [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for managing database migration
```bash
gig migrate <command>
```
Commands:
| Command      | Description                       |
|--------------|-----------------------------------|
| `up`         | Run migrate                       |
| `down`       | Rollback migrate                  |

-----------------

To see the list of available commands, use the `--help` flag
```bash
gig --help
```

## 📝 License

This project is licensed under the [MIT License](https://github.com/awgst/gig/blob/master/LICENSE)
