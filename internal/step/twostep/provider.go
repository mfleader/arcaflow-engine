package twostep

import (
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
}

func (r *runningStep) ProvideStageInput(stage string, input map[string]any) error {
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

func (r *runnableStep) Lifecycle(input map[string]any) (step.Lifecycle[step.LifecycleStageWithSchema], error) {
	return step.Lifecycle[step.LifecycleStageWithSchema]{
		InitialStage: string(StageIDGreet),
		Stages:       []step.LifecycleStageWithSchema{},
	}, nil
}

func (r *runnableStep) RunSchema() map[string]*schema.PropertySchema {
	return nil
}

func (r *runnableStep) Start(input map[string]any, handler step.StageChangeHandler) (step.RunningStep, error) {

	return &runningStep{
		stageChangeHandler: handler,
	}, nil
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
