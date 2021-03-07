package constants

// RedisHostname Hostname of Redis cloud account
const RedisHostname = "redis-16136.c89.us-east-1-3.ec2.cloud.redislabs.com:16136"

// RedisUsername Username of Redis cloud account
const RedisUsername = "default"

// RedisPassword Password of Redis cloud account
const RedisPassword = "zTdtBxLplVe66hAIsgc9dDFChJakXtBf"

// RedisDialProtocol Protocol used to connect
const RedisDialProtocol = "tcp"

// DateFormat format in YYYY-MM-DD
const DateFormat = "2006-01-02"

// Delimiter chosen for the solution is a double pipe character
const Delimiter = "||"

// ControllerMapping path to Data Load controller
const ControllerMapping = "/v1/trading-rules"

// DeployedPath of the data load application
const DeployedPath = "0.0.0.0:8080"

// Currency for the solution is Dollar
const Currency = "$"

// FloatingPointPrecision for the solution is two decimal places
const FloatingPointPrecision = "%.2f"

// FileIdentifier for parsing the file
const FileIdentifier = "dataload"

// SleepTimeBetweenRedisIterations : Interval time chosen for this solution is two seconds
const SleepTimeBetweenRedisIterations = 2000

// ClientType : Key for the client type
const ClientType = "Client-Type"