package flow

import (
	"fmt"
	"gopkg.eicas.io/common/go-pkg/flow-util/apis"
	"gopkg.eicas.io/common/go-pkg/flow-util/step"
	"reflect"
	"strings"
	"sync"
	"testing"
)

type ProductProcessor struct {
}

func (p ProductProcessor) Process(_ int, _ apis.StepInChan, out apis.StepOutChan) {
	for i := 0; i < 2; i++ {
		outData := fmt.Sprintf("processData%d", i)
		fmt.Println("send out data: ", outData)
		out <- outData
	}
}

func NewProductProcessor() apis.Processor {
	return &ProductProcessor{}
}

type DemoProcessor struct {
}

func (d *DemoProcessor) Process(_ int, in apis.StepInChan, out apis.StepOutChan) {
	for inData := range in {
		recData := inData.(string)
		fmt.Println("rec: ", recData)
		out <- strings.ToUpper(recData)
	}
}

func NewDemoProcessor() apis.Processor {
	return &DemoProcessor{}
}

type FinialProcessor struct {
}

func (d *FinialProcessor) Process(_ int, in apis.StepInChan, _ apis.StepOutChan) {
	for inData := range in {
		recData := inData.(string)
		fmt.Println("rec: ", recData)
	}
}

func NewFinialProcessor() apis.Processor {
	return &FinialProcessor{}
}

func TestFlow_AddStep(t *testing.T) {
	type fields struct {
		stepList []apis.Step
	}
	type args struct {
		step apis.Step
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "test add",
			fields: fields{stepList: []apis.Step{step.NewDefaultStep("step1", NewDemoProcessor(), 1, 2)}},
			args:   args{step.NewDefaultStep("step2", NewDemoProcessor(), 1, 2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				stepList: tt.fields.stepList,
			}
			f.AddStep(tt.args.step)
			if f.Count() != 2 {
				t.Errorf("AddStep() = %v, want %v", f.Count(), 2)
			}
		})
	}
}

func TestFlow_Count(t *testing.T) {
	type fields struct {
		stepList []apis.Step
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "test count",
			fields: fields{stepList: []apis.Step{step.NewDefaultStep("step1", NewDemoProcessor(), 1, 2)}},
			want:   1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				stepList: tt.fields.stepList,
			}
			if got := f.Count(); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlow_Start(t *testing.T) {
	type fields struct {
		stepList []apis.Step
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test start",
			fields: fields{stepList: []apis.Step{
				step.NewDefaultStep("step1", NewProductProcessor(), 2, 3),
				step.NewDefaultStep("step2", NewDemoProcessor(), 2, 3),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				stepList: tt.fields.stepList,
			}
			out := f.Start()
			wg := sync.WaitGroup{}
			go func() {
				wg.Add(1)
				defer wg.Done()
				for rec := range out {
					fmt.Println("returned : ", rec.(string))
				}
			}()
			wg.Wait()
		})
	}
}

func TestFlow_Wait(t *testing.T) {
	type fields struct {
		stepList []apis.Step
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test wait",
			fields: fields{stepList: []apis.Step{
				step.NewDefaultStep("step1", NewProductProcessor(), 2, 3),
				step.NewDefaultStep("step2", NewFinialProcessor(), 2, 3),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				stepList: tt.fields.stepList,
			}
			f.Start()
			f.Wait()
			t.Log("Wait() success")
		})
	}
}

func TestNewFlow(t *testing.T) {
	tests := []struct {
		name string
		want *Flow
	}{
		{name: "test new", want: &Flow{make([]apis.Step, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFlow(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlow() = %v, want %v", got, tt.want)
			}
		})
	}
}
