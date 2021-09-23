package apis

type Item interface{}

type StepChan chan Item
type StepInChan <-chan Item
type StepOutChan chan<- Item

// Step 描述每个步骤要做的事情
type Step interface {
	// Title 获取步骤名称
	Title() string
	// Start 启动步骤协程
	Start(in StepChan) StepChan
	// Done 等待协程退出
	Done() <-chan struct{}
	// todo Pause 暂停某一step
	//Pause()
	// todo Resume 继续执行某一step
	//Resume()
}

// Processor 执行器
type Processor interface {
	// Process 执行操作
	Process(id int, in StepInChan, out StepOutChan)
}
