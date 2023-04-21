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
	default:
		return nil
	}
	//return nil
}

func (r *runningStep) run() {
	//switch r.State() {
	//case step.RunningStepStateWaitingForInput:
	//case step.RunningStepStateRunning:
	//
	//}
	select {
	case name, ok := <-r.name:
		// wait until r.name channel has received data
		if !ok {
			// break execution if the r.name channel is closed
			return
		}
		r.state = step.RunningStepStateRunning
		msg := fmt.Sprintf("Hello %s!", name)
		prev_stage := "greet"
		prev_stage_out_id := "success"
		outputData := schema.PointerTo[any](map[string]any{
			"message": msg,
		})
		r.state = step.RunningStepStateFinished
		r.stageChangeHandler.OnStepComplete(
			nil,
			prev_stage,
			&prev_stage_out_id,
			outputData,
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

func (r *runnableStep) Lifecycle(input map[string]any) (step.Lifecycle[step.LifecycleStageWithSchema], error) {
	return step.Lifecycle[step.LifecycleStageWithSchema]{
		InitialStage: string(StageIDGreet),
		Stages:       []step.LifecycleStageWithSchema{},
	}, nil
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
