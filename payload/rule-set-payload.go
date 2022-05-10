package payload

import (
	model "fraud-service/model"
	"sync"
)

type ruleSetPayload struct {
	Data []model.RuleSet `json:"data"`
	sync.RWMutex
}

var instance ruleSetPayload

func GetInstance() *ruleSetPayload {
	return &instance
}

func (p *ruleSetPayload) SetPayload(data []model.RuleSet) {
	p.Lock()
	defer p.Unlock()
	p.Data = data
}

func (p *ruleSetPayload) GetPayload() []model.RuleSet {
	p.RLock()
	defer p.RUnlock()
	return p.Data
}
