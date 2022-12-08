package common

type Communicator struct {
	checkRulesetChan  chan *CheckRulesetWrap
	uploadTrafficChan chan *UploadTrafficWrap
	delCon            chan *DCWrap
}

func NewCommunicator(checkChan chan *CheckRulesetWrap, upChan chan *UploadTrafficWrap, delChan chan *DCWrap) *Communicator {
	return &Communicator{
		checkRulesetChan:  checkChan,
		uploadTrafficChan: upChan,
		delCon:            delChan,
	}
}

func (c *Communicator) CheckRuleset(wrap *CheckRulesetWrap) {
	go func() {
		c.checkRulesetChan <- wrap
	}()
}
func (c *Communicator) UploadTrrafic(wrap *UploadTrafficWrap) {
	go func() {
		c.uploadTrafficChan <- wrap
	}()
}

func (c *Communicator) DelCon(wrap *DCWrap) {
	go func() {
		c.delCon <- wrap
	}()
}
