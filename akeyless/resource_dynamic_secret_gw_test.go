package akeyless

import (
	"fmt"
	"testing"
)

const (
	dockerMysqlHost     = "mysql"
	dockerMysqlPort     = "3306"
	dockerMysqlUser     = "root"
	dockerMysqlPassword = "root_password"
	dockerMysqlDB       = "testdb"

	dockerPostgresHost     = "postgres"
	dockerPostgresPort     = "5432"
	dockerPostgresUser     = "postgres"
	dockerPostgresPassword = "postgres_password"
	dockerPostgresDB       = "testdb"

	dockerMongoHost     = "mongo"
	dockerMongoPort     = "27017"
	dockerMongoUser     = "admin"
	dockerMongoPassword = "mongo_password"
	dockerMongoDB       = "testdb"

	dockerMssqlHost     = "mssql"
	dockerMssqlPort     = "1433"
	dockerMssqlUser     = "sa"
	dockerMssqlPassword = "MssqlPass123!"
	dockerMssqlDB       = "master"

	dockerRedisHost     = "redis"
	dockerRedisPort     = "6379"
	dockerRedisUser     = "default"
	dockerRedisPassword = "redis_password"

	dockerCassandraHost     = "cassandra"
	dockerCassandraPort     = "9042"
	dockerCassandraUser     = "cassandra"
	dockerCassandraPassword = "cassandra"

	dockerRabbitmqURI      = "http://rabbitmq:15672"
	dockerRabbitmqUser     = "admin"
	dockerRabbitmqPassword = "rabbitmq_password"
)

func TestDynamicSecretMysql(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_mysql_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mysql" "%v" {
			name           = "%v"
			mysql_username = "%v"
			mysql_password = "%v"
			mysql_host     = "%v"
			mysql_port     = "%v"
			mysql_dbname   = "%v"
			user_ttl       = "30m"
		}
	`, name, itemPath, dockerMysqlUser, dockerMysqlPassword, dockerMysqlHost, dockerMysqlPort, dockerMysqlDB)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mysql" "%v" {
			name           = "%v"
			mysql_username = "%v"
			mysql_password = "%v"
			mysql_host     = "%v"
			mysql_port     = "%v"
			mysql_dbname   = "%v"
			user_ttl       = "60m"
			tags           = ["test1", "test2"]
		}
	`, name, itemPath, dockerMysqlUser, dockerMysqlPassword, dockerMysqlHost, dockerMysqlPort, dockerMysqlDB)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretPostgresql(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_postgres_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_postgresql" "%v" {
			name                = "%v"
			postgresql_username = "%v"
			postgresql_password = "%v"
			postgresql_host     = "%v"
			postgresql_port     = "%v"
			postgresql_db_name  = "%v"
			user_ttl            = "30m"
		}
	`, name, itemPath, dockerPostgresUser, dockerPostgresPassword, dockerPostgresHost, dockerPostgresPort, dockerPostgresDB)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_postgresql" "%v" {
			name                = "%v"
			postgresql_username = "%v"
			postgresql_password = "%v"
			postgresql_host     = "%v"
			postgresql_port     = "%v"
			postgresql_db_name  = "%v"
			user_ttl            = "60m"
			tags                = ["test1", "test2"]
		}
	`, name, itemPath, dockerPostgresUser, dockerPostgresPassword, dockerPostgresHost, dockerPostgresPort, dockerPostgresDB)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretMongo(t *testing.T) {
	t.Skip("SDK bug: update path is /dynamic-secret-update-mongo instead of /dynamic-secret-update-mongodb (404)")
	t.Parallel()

	name := "ds_mongo_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mongodb" "%v" {
			name                   = "%v"
			mongodb_username       = "%v"
			mongodb_password       = "%v"
			mongodb_host_port      = "%v:%v"
			mongodb_default_auth_db = "admin"
			mongodb_name           = "%v"
			user_ttl               = "30m"
		}
	`, name, itemPath, dockerMongoUser, dockerMongoPassword, dockerMongoHost, dockerMongoPort, dockerMongoDB)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mongodb" "%v" {
			name                   = "%v"
			mongodb_username       = "%v"
			mongodb_password       = "%v"
			mongodb_host_port      = "%v:%v"
			mongodb_default_auth_db = "admin"
			mongodb_name           = "%v"
			user_ttl               = "60m"
			tags                   = ["test1", "test2"]
		}
	`, name, itemPath, dockerMongoUser, dockerMongoPassword, dockerMongoHost, dockerMongoPort, dockerMongoDB)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretMssql(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_mssql_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mssql" "%v" {
			name           = "%v"
			mssql_username = "%v"
			mssql_password = "%v"
			mssql_host     = "%v"
			mssql_port     = "%v"
			mssql_dbname   = "%v"
			user_ttl       = "30m"
		}
	`, name, itemPath, dockerMssqlUser, dockerMssqlPassword, dockerMssqlHost, dockerMssqlPort, dockerMssqlDB)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_mssql" "%v" {
			name           = "%v"
			mssql_username = "%v"
			mssql_password = "%v"
			mssql_host     = "%v"
			mssql_port     = "%v"
			mssql_dbname   = "%v"
			user_ttl       = "60m"
			tags           = ["test1", "test2"]
		}
	`, name, itemPath, dockerMssqlUser, dockerMssqlPassword, dockerMssqlHost, dockerMssqlPort, dockerMssqlDB)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretRedis(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_redis_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_redis" "%v" {
			name     = "%v"
			username = "%v"
			password = "%v"
			host     = "%v"
			port     = "%v"
			user_ttl = "30m"
		}
	`, name, itemPath, dockerRedisUser, dockerRedisPassword, dockerRedisHost, dockerRedisPort)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_redis" "%v" {
			name     = "%v"
			username = "%v"
			password = "%v"
			host     = "%v"
			port     = "%v"
			user_ttl = "60m"
			tags     = ["test1", "test2"]
		}
	`, name, itemPath, dockerRedisUser, dockerRedisPassword, dockerRedisHost, dockerRedisPort)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretCassandra(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_cassandra_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_cassandra" "%v" {
			name               = "%v"
			cassandra_username = "%v"
			cassandra_password = "%v"
			cassandra_hosts    = "%v"
			cassandra_port     = "%v"
			user_ttl           = "30m"
		}
	`, name, itemPath, dockerCassandraUser, dockerCassandraPassword, dockerCassandraHost, dockerCassandraPort)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_cassandra" "%v" {
			name               = "%v"
			cassandra_username = "%v"
			cassandra_password = "%v"
			cassandra_hosts    = "%v"
			cassandra_port     = "%v"
			user_ttl           = "60m"
			tags               = ["test1", "test2"]
		}
	`, name, itemPath, dockerCassandraUser, dockerCassandraPassword, dockerCassandraHost, dockerCassandraPort)

	testItemResource(t, itemPath, config, configUpdate)
}

func TestDynamicSecretRabbitmq(t *testing.T) {
	skipIfNoGateway(t)
	t.Parallel()

	name := "ds_rabbitmq_test"
	itemPath := testPath(name)

	config := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_rabbitmq" "%v" {
			name                = "%v"
			rabbitmq_admin_user = "%v"
			rabbitmq_admin_pwd  = "%v"
			rabbitmq_server_uri = "%v"
			rabbitmq_user_tags  = "management"
			user_ttl            = "30m"
		}
	`, name, itemPath, dockerRabbitmqUser, dockerRabbitmqPassword, dockerRabbitmqURI)

	configUpdate := fmt.Sprintf(`
		resource "akeyless_dynamic_secret_rabbitmq" "%v" {
			name                = "%v"
			rabbitmq_admin_user = "%v"
			rabbitmq_admin_pwd  = "%v"
			rabbitmq_server_uri = "%v"
			rabbitmq_user_tags  = "administrator"
			user_ttl            = "60m"
			tags                = ["test1", "test2"]
		}
	`, name, itemPath, dockerRabbitmqUser, dockerRabbitmqPassword, dockerRabbitmqURI)

	testItemResource(t, itemPath, config, configUpdate)
}
