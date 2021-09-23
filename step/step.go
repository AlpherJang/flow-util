package step

import (
	"github.com/golang/glog"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
	"sync"
)

// defaultStep 内置step
// 每个step内仅包含一种操作
// 至于每个step的oper数量由processCount进行控制
type defaultStep struct {
	title                 string
	bufSize, processCount int
	processor             apis.Processor
	done                  chan struct{}
}

func (d *defaultStep) Title() string {
	return d.title
}

func (d *defaultStep) Start(in apis.StepChan) apis.StepChan {
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
		}(i, d.processor, in, out)
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

func (d *defaultStep) Done() <-chan struct{} {
	return d.done
}

func NewDefaultStep(title string, processor apis.Processor, bufSize, processCount int) apis.Step {
	if bufSize == 0 {
		bufSize = 1
	}
	if processCount == 0 {
		processCount = 1
	}
	return &defaultStep{
		bufSize:      bufSize,
		processor:    processor,
		processCount: processCount,
		title:        title,
	}
}
