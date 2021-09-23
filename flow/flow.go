package flow

import (
	"github.com/golang/glog"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
)

// Flow 描述完成流程，每个step都是按顺序写入到flow中，依次启动，数据并行处理串行流转
type Flow struct {
	stepList []apis.Step
}

func NewFlow() *Flow {
	return &Flow{
		stepList: make([]apis.Step, 0),
	}
}

// AddStep 添加步骤
func (f *Flow) AddStep(step apis.Step) {
	f.stepList = append(f.stepList, step)
	glog.V(7).Infof("flow add step %s", step.Title())
}

// Count 统计总步骤
func (f *Flow) Count() int {
	return len(f.stepList)
}

// Start 启动flow
func (f *Flow) Start() apis.StepChan {
	var stepChan apis.StepChan
	for _, item := range f.stepList {
		stepChan = item.Start(stepChan)
		glog.V(7).Infof("step %s has start", item.Title())
	}
	return stepChan
}

// Wait 等待flow执行完成
func (f *Flow) Wait() {
	for _, item := range f.stepList {
		<-item.Done()
		glog.V(7).Infof("step %s has done", item.Title())
	}
}
