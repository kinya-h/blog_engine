package db

func (q *Queries) GetDBTX() DBTX {
	return q.db
}
