package db

// Conf represents a database connection configuration
type Conf struct {
	Disable      bool   `json:"disable"`
	Debug        bool   `json:"debug"`
	AppName      string `json:"app_name"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	PoolSize     int    `json:"pool_size"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// Config represents a configuration for this package
type Config struct {
	Master Conf   `json:"master"`
	Slaves []Conf `json:"slaves"`
}
