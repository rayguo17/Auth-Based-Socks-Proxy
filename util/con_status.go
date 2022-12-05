package util

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"os"
)

type Connection struct {
	id            string
	username      string
	DstAddr       string
	status        string
	executeStatus string
	cmdType       string
}

type ConnPrinter struct {
	conns []*Connection
	size  int
}

func NewConnPrinter(count int) *ConnPrinter {
	return &ConnPrinter{
		conns: make([]*Connection, 0, count),
		size:  count,
	}
}

func (cp *ConnPrinter) AddCon(connection *Connection) {
	cp.conns = append(cp.conns, connection)
}
func (cp *ConnPrinter) PrintStatus() {

	log.Printf("%d connections accepted\n", cp.size)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "id", "username", "dstAddr", "status", "executeStatus", "cmdType"})
	for i, v := range cp.conns {
		t.AppendRow([]interface{}{
			i, v.id, v.username, v.DstAddr, v.status, v.executeStatus, v.cmdType,
		})
	}
	t.Render()
}

func NewConnection(id string, username string, DstAddr string, status string, executeStatus string, cmdType string) *Connection {
	return &Connection{
		id:            id,
		username:      username,
		DstAddr:       DstAddr,
		status:        status,
		executeStatus: executeStatus,
		cmdType:       cmdType,
	}
}
