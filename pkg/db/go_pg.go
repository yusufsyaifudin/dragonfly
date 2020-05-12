package db

// NewConnectionGoPG will create new connection for selected package.
func NewConnectionGoPG(config Config) (sql SQL, err error) {
	return &goPgSQL{
		masterConf: config.Master,
		slaveConf:  config.Slaves,
	}, nil
}

// goPgSQL is a struct implements SQL interface
type goPgSQL struct {
	masterConf Conf
	slaveConf  []Conf
}

// Writer always use the master.
func (g *goPgSQL) Writer() SQLWriter {
	return connectorGoPgWriter(g.masterConf)
}

// Reader using slaves instance.
// TODO: selecting slave
func (g *goPgSQL) Reader() SQLReader {
	return connectorGoPgWriter(g.masterConf)
}
