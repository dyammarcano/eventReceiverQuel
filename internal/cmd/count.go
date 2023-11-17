package cmd

import (
	"fmt"
	"github.com/caarlos0/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	Message             struct{}
	EventWaitingCommand func() tea.Msg
	EventCountMap       chan int
	Messages            chan string

	Model struct {
		eventCountMap      EventCountMap
		messages           Messages
		eventCount         int
		spinnerModel       spinner.Model
		terminationRequest bool
	}

	CounterService struct {
		model     Model
		program   *tea.Program
		increment int
	}
)

func createWaitForEventsCommand(eventCountMap EventCountMap) EventWaitingCommand {
	return func() tea.Msg {
		return <-eventCountMap
	}
}

func (m Model) initializeModel() tea.Cmd {
	return tea.Batch(
		m.spinnerModel.Tick,
		tea.Cmd(createWaitForEventsCommand(m.eventCountMap)),
	)
}

func (m Model) updateModelBasedOnMessageType(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.terminationRequest = true
		return m, tea.Quit

	case Message:
		m.eventCount++
		return m, tea.Cmd(createWaitForEventsCommand(m.eventCountMap))

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinnerModel, cmd = m.spinnerModel.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) generateViewStateRepresentation() string {
	stateRepresentation := fmt.Sprintf("\n %s Events received: %d\n", m.spinnerModel.View(), m.eventCount)

	if m.terminationRequest {
		stateRepresentation += "\n"
	}
	return stateRepresentation
}

func startNewCounterService() *CounterService {
	newModel := Model{
		eventCountMap: make(EventCountMap),
		spinnerModel:  spinner.New(),
	}
	program := tea.NewProgram(newModel)

	return &CounterService{
		program: program,
		model:   newModel,
	}
}

func (c *CounterService) incrementEventCount() {
	c.increment++
	c.model.eventCountMap <- c.increment
}

func (c *CounterService) startCounterService() {
	if _, errorOccured := c.program.Run(); errorOccured != nil {
		log.Error(errorOccured.Error())
	}
}
