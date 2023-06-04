package migration

var CreateMigrationSQLTemplate = `CREATE TABLE IF NOT EXISTS {{.TableName}} ()`
var DropMigrationSQLTemplate = `DROP TABLE IF EXISTS {{.TableName}}`
