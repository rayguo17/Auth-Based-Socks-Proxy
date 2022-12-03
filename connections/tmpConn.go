package connections

import "fmt"

type TmpConn struct {
	AuthChan    chan bool //authenticated
	id          string    // ip:port
	TermChan    chan bool
	addedChan   chan bool
	deletedChan chan bool
}

type TmpManager struct {
	conns         map[string]*TmpConn
	addChan       chan *TmpConn
	deleteChan    chan string //id
	TerminateChan chan bool
}

func (tm *TmpManager) MainRoutine() {
	for {
		select {
		case conn := <-tm.addChan:
			tm.conns[conn.id] = conn
			conn.addedChan <- true
		case id := <-tm.deleteChan:
			fmt.Println("haha find it meaning less...", id)
		}
	}
}
