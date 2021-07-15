![globetrot logo red](./docs/assets/globetrot_logo_r.png)
![globetrot logo white](./docs/assets/globetrot_logo_w.png)
![globetrot logo blue](./docs/assets/globetrot_logo_b.png)
---

Globetrot is a cross-platform CLI tool used for managing database changes/migrations. It currently works with MySQL, PostgreSQL and SQL Server with plans of other database support in the future.

## Getting globetrot

You can download Globetrot releases from https://github.com/dustyhoppe/globetrot/releases


## Getting started

From the command line run the following initialization command to create your globetrot project directory.

```bash
./globetrot init --path "./{directory_path}"

# ./globetrot init --path "./globetrot_playground"
```

You'll now see your globetrot project directory structure setup at the path you specified. The directory consists of the following sub-directories:

* up
* procs


The `up` directory consists of update scripts that will be run one-time. `up` scripts are typically used for managing database schema changes as well as one-time data modifications.

***Note: making modifications to up scripts after execution will result in an error***

The `procs` directory consists of scripts for creating database procedures (functions, stored procedures, triggers, etc). These scripts will run anytime the script has been modified since the last migration.

You can now add scripts to the appropriate directories and run the `migrate` command to apply a migration.

```bash
./globetrot migrate --filePath "./globetrot_playground" \
    --database "sample_db" \
    --databaseType "mysql" \
    --server "localhost" \
    --username "root" \
    --password "password123!" \
    --port 3306
```


## Commands

Below is a reference of the commands supported by globetrot

### init

The `init` sub-command is used for creating a globetrot project directory structure at whatever directory path provided as an argument.

```
Usage:
  globetrot init [flags]

Flags:
  -h, --help          help for init
  -p, --path string   The path to the directory that should be initialized as a globetrot project directory.
```

A sample execution of the `init` command looks as follows:

```bash
./globetrot init --path "./playground" 
```

### migrate

Executing the `migrate` command performs a database migration. The migration first executes any `up` scripts that have yet to be run against the target database. Once `up` scripts have been run, any new or modified `procs` scripts will then be applied.

***Note: upon first execution, this command will setup it's migration tracking schema which is currently a single table named `scripts_run`***

```
Usage:
  globetrot migrate [flags]

Flags:
  -c, --configPath string     The path to the configuration to use in-place of passing command line arguments
  -d, --database string       The name of the database
  -t, --databaseType string   Indicates the database type/platform (mysql, postgres, sqlserver) (default "mysql")
      --dryRun                Indicates whether the command should perform a dry run
  -e, --env string            The environment the migration is targeting
  -f, --filePath string       The directory where your SQL scripts are located (default ".\\")
  -h, --help                  help for migrate
  -p, --password string       The password to use when connecting to the database
  -P, --port int              The port the database server the database is listening at (default 3306)
  -s, --server string         The host server the database is located at (default "localhost")
  -u, --username string       The username to use when connecting to the database (default "root")

required flag(s) "database", "databaseType", "filePath", "password", "port", "server", "username" not set
```

A sample execution of the `migrate` command looks as follows:

```bash
./globetrot migrate --filePath "./playground" \
    --database "sample_db" \
    --databaseType "mysql" \
    --server "localhost" \
    --username "root" \
    --password "password123!" \
    --port 3306 
```

## Environment Scripts

Globetrot provides support for running scripts that target a specific environment. To flag a script as being an environment specific script, make sure the script file extension is `.env.sql`. The full name of the script should follow the format `{SCRIPT_NAME}{ENVIRONMENT_NAME}.env.sql`

If I wanted to create an environment script that targets the `PROD` environment, I could name the script `LoadUsers.PROD.env.sql`. When running the migration, I'd also need to specify the `--env` argument and give it a value of `PROD`.

```bash
./globetrot migrate --filePath "./playground" \
    --database "sample_db" \
    --databaseType "mysql" \
    --server "localhost" \
    --username "root" \
    --password "password123!" \
    --port 3306 \
    --env "PROD" # This tells globetrot I'm running the script against the PROD environment
```

## Configuration

In addition to passing configuration via command line arguments, you can use config files or environment variables as an alternative or supplemental approach for passing configuration.

If configuration is passed using multiple sources, the priority order for configuration sources are:

1. Command line arguments
2. Environment variables
3. Configuration files

### Environment Variables

To pass configuration via environment variables, you'll need to use the `GLOBETROT` prefix and make sure variables are set in all-caps.

For example, to set the `password` to use when connecting to the database:

```bash
export GLOBETROT_PASSWORD=password123!
```

### Configuration Files

You may load configuration via a file by specifying the directory in which the file is located using the `--configPath` argument. The supported formats are:
* JSON
* YAML
* envfile
* HCL
* TOML

Below are examples of providing configuration using each of these formats.

#### **`globetrot.json`**
```json
{
    "username": "sa",
    "password": "password123!",
    "server": "127.0.0.1",
    "database": "globetrot",
    "port": 1433,
    "databaseType": "sqlserver",
    "filePath": "./playground",
    "environment": "production"
}
```

#### **`globetrot.yaml`**
```yaml
username: sa
password: password123!
port: 1433
databaseType: sqlserver
server: 127.0.0.1
database: Globetrot
filePath: ./playground
env: production
```

#### **`globetrot.env`**
```env
USERNAME=sa
PASSWORD=password123!
PORT=1433
DATABASETYPE=sqlserver
SERVER=127.0.0.1
DATABASE=Globetrot
FILEPATH=./playground
ENV=production
```

#### **`globetrot.hcl`**
```hcl
username = "sa"
password = "password123!"
port = 1433
databaseType = "sqlserver"
server = "127.0.0.1"
database = "Globetrot"
filePath = "./playground"
env = "production"
```

#### **`globetrot.toml`**
```toml
username = "sa"
password = "password123!"
port = 1433
databaseType = "sqlserver"
server = "127.0.0.1"
database = "Globetrot"
filePath = "../playground"
env = "production"
```