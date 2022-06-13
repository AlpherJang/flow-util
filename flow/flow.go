package flow

import (
	"github.com/golang/glog"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
)

// Flow
// @Description describe whole work flow, all steps will set to flow follow there index,
// and will start one by one, data will transfer by channel
// @Author ZhangHao
// @Date 2022-06-09 14:47:45
type Flow struct {
	stepList [][]apis.Step
}

func NewFlow() *Flow {
	return &Flow{
		stepList: make([][]apis.Step, 0),
	}
}

// AddStep
// @Description add step into work flow
// @Author ZhangHao
// @Date 2022-06-09 14:49:23
func (f *Flow) AddStep(step apis.Step) {
	f.stepList = append(f.stepList, []apis.Step{step})
	glog.V(7).Infof("flow add step %s", step.Title())
}

// Count
// @Description calculates all work steps in work flow
// @Author ZhangHao
// @Date 2022-06-09 14:49:45
func (f *Flow) Count() int {
	return len(f.stepList)
}

// Start
// @Description start all steps in work flow
// @Author ZhangHao
// @Date 2022-06-09 14:50:20
func (f *Flow) Start() apis.StepChan {
	var stepChan apis.StepChan
	for _, item := range f.stepList {
		for _, itemStep := range item {
			stepChan = itemStep.Start(stepChan)
			glog.V(7).Infof("step %s has start", itemStep.Title())
		}
	}
	return stepChan
}

// Wait
// @Description wait work flow finish
// @Author ZhangHao
// @Date 2022-06-09 14:50:34
func (f *Flow) Wait() {
	for _, item := range f.stepList {
		for _, itemStep := range item {
			<-itemStep.Done()
			glog.V(7).Infof("step %s has done", itemStep.Title())
		}
	}
}
