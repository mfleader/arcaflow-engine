package twostep

import (
	"context"
	"fmt"
	"go.flow.arcalot.io/engine/internal/step"
	"go.flow.arcalot.io/pluginsdk/schema"
	"sync"
)

type twostepProvider struct {
}

type StageID string

const (
	StageIDGreet StageID = "greet"
	StageIDMeet  StageID = "meet"
)

type runnableStep struct {
}

type runningStep struct {
	stageChangeHandler step.StageChangeHandler
	name               chan string
	state              step.RunningStepState
	currentStage       StageID
	lock               *sync.Mutex
	ctx                context.Context
	cancel             context.CancelFunc
	inputAvailable     bool
}

func (r *runningStep) OnStageChange(
	step_ step.RunningStep,
	previousStage *string,
	previousStageOutputID *string,
	previousStageOutput *any,
	newStage string,
	waitingForInput bool) {

}

func (r *runningStep) OnStepComplete(
	step_ step.RunningStep,
	previousStage string,
	prevoiusStageOutputID *string,
	previousStageOutput *any) {

}

func (r *runningStep) ProvideStageInput(stage string, input map[string]any) error {

	defer close(r.name)
	switch stage {
	case string(StageIDGreet):
		r.name <- fmt.Sprintf("%s", input["name"])
		r.state = step.RunningStepStateRunning
		return nil
	case string(StageIDMeet):
		return nil
	default:
		return nil
	}
	//return nil
}

func (r *runningStep) run() {
	//defer close(r.name)

	// get input on first pass
	waitingForInput := false
	r.lock.Lock()
	if !r.inputAvailable {
		r.state = step.RunningStepStateWaitingForInput
		waitingForInput = true
	} else {
		r.state = step.RunningStepStateRunning
	}
	r.lock.Unlock()

	// notify stage handler
	r.stageChangeHandler.OnStageChange(
		r,
		nil,
		nil,
		nil,
		string(StageIDGreet),
		waitingForInput)

	select {
	// wait until r.name channel has received data
	case name, ok := <-r.name:

		// break execution if the r.name channel is closed
		if !ok {
			return
		}

		// lock state, so this goroutine can modify it to running
		r.lock.Lock()
		r.state = step.RunningStepStateRunning
		r.lock.Unlock()

		// Do the thing (say hello)
		msg := fmt.Sprintf("Hello %s!", name)
		output_data := schema.PointerTo[any](map[string]any{
			"message": msg,
		})
		// that's it!

		// Assign this stage's output id
		output_id := schema.PointerTo("success")

		r.lock.Lock()
		r.state = step.RunningStepStateFinished
		r.lock.Unlock()

		r.stageChangeHandler.OnStepComplete(
			r,
			string(StageIDGreet),
			output_id,
			output_data,
		)
	}
}

func (r *runnableStep) Start(input map[string]any, handler step.StageChangeHandler) (step.RunningStep, error) {
	ctx, cancel := context.WithCancel(context.Background())
	running_step := &runningStep{
		stageChangeHandler: handler,
		name:               make(chan string, 1),
		state:              step.RunningStepStateStarting,
		currentStage:       StageIDGreet,
		lock:               &sync.Mutex{},
		ctx:                ctx,
		cancel:             cancel,
		inputAvailable:     false,
	}
	go running_step.run()
	return running_step, nil
}

func (r *runningStep) CurrentStage() string {
	r.lock.Lock()
	defer r.lock.Unlock()
	return string(r.currentStage)
}

func (r *runningStep) State() step.RunningStepState {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.state
}

func (r *runningStep) Close() error {
	r.cancel()
	return nil
}

func (r *runnableStep) RunSchema() map[string]*schema.PropertySchema {
	return nil
}

func New() step.Provider {
	return &twostepProvider{}
}

func (p *twostepProvider) LoadSchema(inputs map[string]any, workflowContext map[string][]byte) (step.RunnableStep, error) {
	return &runnableStep{}, nil
}

func (p *twostepProvider) Kind() string {
	return "twostep"
}

func (p *twostepProvider) ProviderSchema() map[string]*schema.PropertySchema {
	return map[string]*schema.PropertySchema{}
}

func (p *twostepProvider) RunProperties() map[string]struct{} {
	return map[string]struct{}{}
}

var greetingLifecycleStage = step.LifecycleStage{
	ID:           string(StageIDGreet),
	WaitingName:  "waiting for greeting",
	RunningName:  "greeting",
	FinishedName: "greeted",
	InputFields: map[string]struct{}{
		"name":     {},
		"nickname": {},
	},
	NextStages: nil,
	Fatal:      false,
}

func (p *twostepProvider) Lifecycle() step.Lifecycle[step.LifecycleStage] {
	return step.Lifecycle[step.LifecycleStage]{
		InitialStage: string(StageIDGreet),
		Stages: []step.LifecycleStage{
			greetingLifecycleStage,
		},
	}
}

var inputSchema = map[string]*schema.PropertySchema{}

var outputSchema = map[string]*schema.StepOutputSchema{}

func (r *runnableStep) Lifecycle(input map[string]any) (step.Lifecycle[step.LifecycleStageWithSchema], error) {
	return step.Lifecycle[step.LifecycleStageWithSchema]{
		InitialStage: string(StageIDGreet),
		Stages: []step.LifecycleStageWithSchema{
			{
				LifecycleStage: greetingLifecycleStage,
				InputSchema:    inputSchema,
				Outputs:        outputSchema,
			},
		},
	}, nil
}
