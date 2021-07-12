package flags

import "github.com/urfave/cli"

var (

	// ServerNameFlag ...
	ServerNameFlag = cli.StringFlag{
		Name:   "server_name",
		Usage:  "server name",
		EnvVar: "SERVER_NAME",
		Value:  "Hotel Management",
	}

	// ServerHostFlag ...
	ServerHostFlag = cli.StringFlag{
		Name:   "server_host",
		Usage:  "Server Host",
		EnvVar: "SERVER_HOST",
		Value:  "",
	}

	// ServerPortFlag ...
	ServerPortFlag = cli.StringFlag{
		Name:   "server_port",
		Usage:  "Server Port",
		EnvVar: "SERVER_PORT",
		Value:  "8080",
	}

	// MongoDatabaseNameFlag ...
	MongoDatabaseNameFlag = cli.StringFlag{
		Name:   "database_name",
		Usage:  "Database name",
		EnvVar: "DATABASE_NAME",
		Value:  "hotel",
	}

	// MongoHostFlag ...
	MongoHostFlag = cli.StringFlag{
		Name:   "database_host",
		Usage:  "Database host",
		EnvVar: "DATABASE_HOST",
		Value:  "mongodb://my_first_mongodb",
	}

	// MongoPortFlag ...
	MongoPortFlag = cli.StringFlag{
		Name:   "database_port",
		Usage:  "Database port",
		EnvVar: "DATABASE_PORT",
		Value:  "27017",
	}

	// StorageAccessKeyFlag ...
	StorageAccessKeyFlag = cli.StringFlag{
		Name:   "storage_access_key",
		Usage:  "Storage access key",
		EnvVar: "STORAGE_ACCESS_KEY",
		Value:  "",
	}

	// StorageSecretKeyFlag ...
	StorageSecretKeyFlag = cli.StringFlag{
		Name:   "storage_secret_key",
		Usage:  "Storage secret key",
		EnvVar: "STORAGE_SECRET_KEY",
		Value:  "",
	}

	// StorageRegionFlag ...
	StorageRegionFlag = cli.StringFlag{
		Name:   "storage_region",
		Usage:  "Storage region",
		EnvVar: "STORAGE_REGION",
		Value:  "",
	}

	// StorageName ...
	StorageNameFlag = cli.StringFlag{
		Name:   "storage_name",
		Usage:  "Storage name",
		EnvVar: "STORAGE_NAME",
		Value:  "room-image",
	}

	// JaegerHostFlag ...
	JaegerHostFlag = cli.StringFlag{
		Name:   "Jaeger_Host",
		Usage:  "Jaeger Host",
		EnvVar: "JAEGER_HOST",
		Value:  "tracer",
	}

	// JaegerPortFlag ...
	JaegerPortFlag = cli.StringFlag{
		Name:   "Jaeger_Port",
		Usage:  "Jaeger Port",
		EnvVar: "JAEGER_PORT",
		Value:  "6831",
	}
)

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			ctx.GlobalSet(name, ctx.String(name))
		}
		return action(ctx)
	}
}
