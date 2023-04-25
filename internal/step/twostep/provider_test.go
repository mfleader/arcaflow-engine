package twostep_test

import (
	"fmt"
	"go.arcalot.io/assert"
	"go.flow.arcalot.io/engine/internal/step"
	"go.flow.arcalot.io/engine/internal/step/twostep"
	"testing"
)

type stageChangeHandler struct {
	message chan string
}

func (s *stageChangeHandler) OnStageChange(_ step.RunningStep, _ *string, _ *string, _ *any, _ string, _ bool) {

}

func (s *stageChangeHandler) OnStepComplete(
	_ step.RunningStep,
	previousStage string,
	previousStageOutputID *string,
	previousStageOutput *any) {

	if previousStage != "meet" {
		panic(fmt.Errorf("invalid previous stage: %s",
			previousStage))
	}
	if previousStageOutputID == nil {
		panic(fmt.Errorf("no previous stage output ID"))
	}
	if *previousStageOutputID != "success" {
		panic(fmt.Errorf("invalid previous stage output ID: %s",
			*previousStageOutputID))
	}
	if previousStageOutput == nil {
		panic(fmt.Errorf("no previous stage output"))
	}
	message := (*previousStageOutput).(map[string]any)["message"].(string)
	s.message <- message
}

func TestProviderKind(t *testing.T) {
	provider := twostep.New()
	assert.Equals(t, provider.Kind(), "twostep")
}

func TestProvider(t *testing.T) {
	provider := twostep.New()
	runnable, err := provider.LoadSchema(
		map[string]any{},
		map[string][]byte{})
	assert.NoError(t, err)

	handler := &stageChangeHandler{
		message: make(chan string),
	}

	running, err := runnable.Start(map[string]any{}, handler)
	assert.NoError(t, err)
	assert.NoError(t, running.ProvideStageInput(
		"greet", map[string]any{
			"name": "Arca Lot",
		}))
	message := <-handler.message
	//assert.Equals(t, message, "Hello Arca Lot!")

	assert.Equals(t, message, "Arca Lot meet Anon.")

}
