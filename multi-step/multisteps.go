package multi_step

import (
	"github.com/golang/glog"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
	"sync"
)

// multiSteps 内置step
// 每个step内仅包含一种操作
// 至于每个step的oper数量由processCount进行控制
type multiSteps struct {
	title                 string
	bufSize, processCount int
	processors            []apis.Processor
	done                  chan struct{}
}

func (d *multiSteps) Title() string {
	return d.title
}

// Start
// @Description start this step, and start all processor in step
// @Author AlpheJang
// @Date 2022-06-16 14:25:28
func (d *multiSteps) Start(in apis.StepChan) apis.StepChan {
	d.done = make(chan struct{})
	out := make(apis.StepChan, d.bufSize)

	var wg sync.WaitGroup
	wg.Add(d.processCount)
	for i := 0; i < d.processCount; i++ {
		go func(index int, processor apis.Processor, inChan <-chan apis.Item, outChan chan<- apis.Item) {
			defer wg.Done()
			glog.V(7).Infof("step %s has start processor %d", d.title, index)
			defer glog.V(7).Infof("step %s processor %d has finished", d.title, index)
			processor.Process(index, inChan, outChan)
		}(i, d.processors[i], in, out)
	}

	go func() {
		wg.Wait()
		glog.V(7).Infof("all goroutine has finished")
		close(out)
		glog.V(7).Infof("out chan closed")
		close(d.done)
		glog.V(7).Infof("done sign has closed")
	}()
	return out
}

func (d *multiSteps) Done() <-chan struct{} {
	return d.done
}

func (d *multiSteps) AddProcessor(processor apis.Processor) {
	d.processors = append(d.processors, processor)
}

func NewMultiStep(title string, processors []apis.Processor, bufSize, processCount int) apis.Step {
	if bufSize == 0 {
		bufSize = 1
	}
	if processCount == 0 {
		processCount = 1
	}
	return &multiSteps{
		bufSize:      bufSize,
		processors:   processors,
		processCount: processCount,
		title:        title,
	}
}
