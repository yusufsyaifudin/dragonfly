package model

type Connection struct {
	ID                         string
	PostgresMasterDebug        bool
	PostgresMasterHost         string
	PostgresMasterPort         string
	PostgresMasterUsername     string
	PostgresMasterPassword     string
	PostgresMasterDatabase     string
	PostgresMasterPoolSize     int
	PostgresMasterIdleTimeout  int
	PostgresMasterMaxConnAge   int
	PostgresMasterReadTimeout  int
	PostgresMasterWriteTimeout int
}

type Connections []*Connection
