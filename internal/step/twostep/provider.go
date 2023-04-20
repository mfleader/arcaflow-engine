package twostep

import (
	"fmt"
	"go.flow.arcalot.io/engine/internal/step"
	"go.flow.arcalot.io/pluginsdk/schema"
)

type twostepProvider struct {
}

type StageID string

const (
	StageIDGreet StageID = "greet"
)

type runnableStep struct {
}

type runningStep struct {
	stageChangeHandler step.StageChangeHandler
	name               chan string
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

type outPut struct {
	msg string
}

func (r *runningStep) ProvideStageInput(stage string, input map[string]any) error {
	//r.stageChangeHandler
	msg := fmt.Sprintf("Hello %s!", input["name"])
	//out := &outPut{msg: msg}
	//s := &string
	//*s = "success"
	prev_stage := "greet"
	prev_stage_out_id := "success"
	// this is weird
	outputData := schema.PointerTo[any](map[string]any{
		//"message": fmt.Sprintf("Hello %s!", name),
		"message": msg,
	})
	r.name <- msg
	r.stageChangeHandler.OnStepComplete(
		nil,
		prev_stage,
		&prev_stage_out_id,
		outputData,
	)

	defer close(r.name)
	return nil
}

func (r *runningStep) CurrentStage() string {
	return "derp"
}

func (r *runningStep) State() step.RunningStepState {
	return step.RunningStepStateStarting
}

func (r *runningStep) Close() error {

	return nil
}

func (r *runningStep) run() {
	//r.
}

func (r *runnableStep) Start(input map[string]any, handler step.StageChangeHandler) (step.RunningStep, error) {
	running_step := &runningStep{
		stageChangeHandler: handler,
		name:               make(chan string, 1),
	}
	//running_step.run()
	return running_step, nil
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
