package step

import (
	"fmt"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
	"strings"
	"testing"
	"time"
)

type DemoProcessor struct {
}

func (d *DemoProcessor) Process(_ int, in apis.StepInChan, out apis.StepOutChan) {
	recData := (<-in).(string)
	fmt.Println("rec: ", recData)
	out <- strings.ToUpper(recData)
}

func NewDemoProcessor() apis.Processor {
	return &DemoProcessor{}
}

func ExampleNewDefaultStep() {
	stepDefault := NewDefaultStep("demo step", NewDemoProcessor(), 2, 5)
	fmt.Println(stepDefault.Title())
}

func TestDefaultStep_Title(t *testing.T) {
	stepDefault := NewDefaultStep("demo step", NewDemoProcessor(), 2, 5)
	t.Logf("step name is : %s", stepDefault.Title())
}

func TestDefaultStep_Start(t *testing.T) {
	stepDefault := NewDefaultStep("demo step", NewDemoProcessor(), 2, 5)
	in := make(apis.StepChan)
	out := stepDefault.Start(in)
	for i := 0; i < 5; i++ {
		in <- fmt.Sprintf("test%d", i)
	}
	for rec := range out {
		fmt.Println("return: ", rec)
	}
}

func TestDefaultStep_Done(t *testing.T) {
	stepDefault := NewDefaultStep("demo step", NewDemoProcessor(), 2, 5)
	in := make(apis.StepChan, 2)
	out := stepDefault.Start(in)
	for i := 0; i < 5; i++ {
		in <- fmt.Sprintf("test%d", i)
	}
	close(in)
	go func() {
		for rec := range out {
			fmt.Println("return: ", rec)
			time.Sleep(2)
		}
	}()
	<-stepDefault.Done()
	fmt.Println("all step finished")
}
