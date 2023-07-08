package messenger

import (
      "database/sql"
       _ "github.com/lib/pq"
)

type MessageProc struct {
	Connection *sql.DB
}

func New(connection *sql.DB) *MessageProc {
    return &MessageProc{
        Connection: connection,
    }
}
