package util

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"os"
)

type UserEntry struct {
	Username        string
	UplinkTraffic   string
	DownlinkTraffic string
	Enable          string
	Deleted         string
	lastSeen        string
	ActiveConn      string
	TotalConn       string
}
type UserPrinter struct {
	users []*UserEntry
	size  int
}

func (up *UserPrinter) PrintStatus() {
	log.Printf("%d users registered.\n", up.size)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "username", "UplinkTraffic", "DownlinkTraffic", "Enable", "Deleted", "lastSeen", "ActiveConn", "TotalConn"})
	for i, v := range up.users {
		t.AppendRow([]interface{}{

			i, v.Username, v.UplinkTraffic, v.DownlinkTraffic, v.Enable, v.Deleted, v.lastSeen, v.ActiveConn, v.TotalConn,
		})
	}

	t.Render()
}

func NewUserEntry(u string, ut string, dt string, e string, d string, l string, active string, total string) *UserEntry {
	return &UserEntry{
		Username:        u,
		UplinkTraffic:   ut,
		DownlinkTraffic: dt,
		Enable:          e,
		Deleted:         d,
		lastSeen:        l,
		ActiveConn:      active,
		TotalConn:       total,
	}
}

func NewUserPrinter(count int) *UserPrinter {
	return &UserPrinter{
		users: make([]*UserEntry, 0, count),
		size:  count,
	}
}
func (up *UserPrinter) AddUser(u *UserEntry) {
	up.users = append(up.users, u)
}
