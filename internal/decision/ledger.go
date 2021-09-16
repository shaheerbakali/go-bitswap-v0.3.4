package decision

import (
	"sync"
	"fmt"
	"encoding/csv"
	"os"
	"time"

	pb "github.com/ipfs/go-bitswap/message/pb"
	wl "github.com/ipfs/go-bitswap/wantlist"

	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

func newLedger(p peer.ID) *ledger {
	return &ledger{
		wantList: wl.New(),
		Partner:  p,
	}
}

// Keeps the wantlist for the partner. NOT threadsafe!
type ledger struct {
	// Partner is the remote Peer.
	Partner peer.ID

	// wantList is a (bounded, small) set of keys that Partner desires.
	wantList *wl.Wantlist

	lk sync.RWMutex
}

func (l *ledger) Wants(k cid.Cid, priority int32, wantType pb.Message_Wantlist_WantType) {

	//fmt.Print("---ledger.go...Wants---\n")
	//fmt.Print(" peer ", l.Partner, " wants ", k, "\n")
	//fmt.Print("\n")

	// to save in CSV file:
	//////////////////////////////////////////////////////////////////
	//Create Data Array to Write to CSV File
	t := time.Now()
	timeStamp := fmt.Sprint(t.Format("2006-01-02 15:04:05"))

	cid := fmt.Sprint(k)
	peerid := fmt.Sprint(l.Partner)

	//engine := &Engine{} // create an instance
	//foundVar := engine.MessageReceived()
	//fmt.Print("testVar: ",foundVar,"\n\n")

	data := [][]string{
		{peerid, cid, timeStamp},
	}
	headers := [][]string{
		{"peerID", "CID", "timestamp"},
	}

	// check if file exists
	// if yes, then dont write headers
	// if no, then write headers
	_, err := os.Stat("records.csv")
	if os.IsNotExist(err) { // if file doesnt exists
		// Create CSV File
		file, err := os.OpenFile("records.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			println("failed to open file", err)
		}
		w := csv.NewWriter(file)
		defer w.Flush()
		// Write headers to CSV File
		err = w.WriteAll(headers)
		if err != nil {
			fmt.Println(err)
		}
		// Write data to CSV File
		err = w.WriteAll(data)
		if err != nil {
			fmt.Println(err)
		}

	} else { // if file exists
		// Create CSV File
		file, err := os.OpenFile("records.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			println("failed to open file", err)
		}
		w := csv.NewWriter(file)
		defer w.Flush()
		// Write data to CSV File
		err = w.WriteAll(data)
		if err != nil {
			fmt.Println(err)
		}

	}
	//////////////////////////////////////////////////////////////////


	log.Debugf("peer %s wants %s", l.Partner, k)
	l.wantList.Add(k, priority, wantType)
}

func (l *ledger) CancelWant(k cid.Cid) bool {
	return l.wantList.Remove(k)
}

func (l *ledger) WantListContains(k cid.Cid) (wl.Entry, bool) {

	//println("---ledger.go...WantListContains---")
	//fmt.Print(l.wantList.Contains(k))
	//println("\n")

	return l.wantList.Contains(k)
}
