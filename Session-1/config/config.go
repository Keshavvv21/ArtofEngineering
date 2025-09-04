package config

// Config is the top-level configuration struct for the application.
// It contains both general database connection info (host/port)
// and logical database details (specific DBs and collection names).
type Config struct {
	Database        database        // Basic DB connection settings (host, port)
	DatabaseDetails DatabaseDetails // Names of logical DBs and collections
}

// database holds the raw connection details for the database server.
type database struct {
	Server string // Database server hostname or IP address
	Port   string // Port on which the database server is listening
}

// DatabaseDetails holds logical names of databases and collections
// used by the application for Buyers, Sellers, and Projects.
type DatabaseDetails struct {
	BuyersDBName   string // Name of the database that stores Buyers
	SellersDBName  string // Name of the database that stores Sellers
	ProjectDBName  string // Name of the database that stores Projects
	CollectionName string // Shared or default collection name for inserts/queries
}
