package model

type Connection struct {
	ID                         string
	PostgresMasterDebug        bool
	PostgresMasterHost         string
	PostgresMasterPort         int
	PostgresMasterUsername     string
	PostgresMasterPassword     string
	PostgresMasterDatabase     string
	PostgresMasterPoolSize     int
	PostgresMasterReadTimeout  int
	PostgresMasterWriteTimeout int
}

type Connections []*Connection
