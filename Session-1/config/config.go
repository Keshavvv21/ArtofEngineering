package config

type Config struct {
	Database        database
	DatabaseDetails DatabaseDetails
}

type database struct {
	Server string
	Port   string
}

type DatabaseDetails struct {
	BuyersDBName   string
	SellersDBName  string
	ProjectDBName  string
	CollectionName string
}
