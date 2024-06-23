package tests_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roman-kart/go-initial-project/project/managers"
)

func TestTelegramStartCommandErrorOnEmptyMessage(t *testing.T) {
	_, err := managers.TelegramStartCommandResponse(&managers.StartCommandConfig{
		Message: "",
	})

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, managers.ErrNoMessage)
	}
}

func TestTelegramStartCommandReturnSameMessageFromConfig(t *testing.T) {
	message, err := managers.TelegramStartCommandResponse(&managers.StartCommandConfig{
		Message: "test",
	})

	if assert.NoError(t, err) {
		assert.Equal(t, "test", message)
	}
}

func FuzzTelegramStartCommandResponse(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		message, err := managers.TelegramStartCommandResponse(&managers.StartCommandConfig{
			Message: s,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, s, message)
		}
	})
}

func TestTelegramHelpCommandResponseErrorOnEmptyMainMessage(t *testing.T) {
	_, err := managers.TelegramHelpCommandResponse(
		&managers.HelpCommandConfig{
			MainHelpMessage: "",
		},
		[]string{},
	)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, managers.ErrNoMessage)
	}
}

func TestTelegramHelpCommandResponseOnlyMainMessageOnZeroArgsAndNoCommandsHelpMessages(t *testing.T) {
	message, err := managers.TelegramHelpCommandResponse(
		&managers.HelpCommandConfig{
			MainHelpMessage: "test",
		},
		[]string{},
	)

	if assert.NoError(t, err) {
		assert.Equal(t, "test", message)
	}
}

func TestTelegramHelpCommandResponseMainMessage(t *testing.T) {
	message, err := managers.TelegramHelpCommandResponse(
		&managers.HelpCommandConfig{
			MainHelpMessage: "Test",
			CommandsHelpMessages: map[string]managers.HelpCommandMessages{
				"/test": {
					ShortMessage:  "testing",
					DetailMessage: "testing",
				},
				"test": {
					ShortMessage:  "testing",
					DetailMessage: "testing",
				},
			},
		},
		[]string{},
	)

	if assert.NoError(t, err) {
		expectedMessage := `Test

*Команды:*
/test - testing
test - testing`

		assert.Equal(t, expectedMessage, message)
	}
}

func TestTelegramHelpCommandResponseCommandsHelpMessages(t *testing.T) {
	cfg := &managers.HelpCommandConfig{
		MainHelpMessage: "Test",
		CommandsHelpMessages: map[string]managers.HelpCommandMessages{
			"/test": {
				ShortMessage:  "testing short /test",
				DetailMessage: "testing detail /test",
			},
			"test": {
				ShortMessage:  "testing short test",
				DetailMessage: "testing detail test",
			},
		},
	}

	message, err := managers.TelegramHelpCommandResponse(cfg, []string{"test"})
	if assert.NoError(t, err) {
		expectedMessage := `testing short /test

testing detail /test`

		assert.Equal(t, expectedMessage, message)
	}

	message, err = managers.TelegramHelpCommandResponse(cfg, []string{"/test"})
	if assert.NoError(t, err) {
		expectedMessage := `testing short /test

testing detail /test`

		assert.Equal(t, expectedMessage, message)
	}

	message, err = managers.TelegramHelpCommandResponse(cfg, []string{"wrong"})
	if assert.NoError(t, err) {
		expectedMessage := "Команда `/wrong` не найдена"

		assert.Equal(t, expectedMessage, message)
	}
}

func FuzzTelegramHelpCommandResponseOnlyMainMessageOnZeroArgsAndNoCommandsHelpMessages(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		message, err := managers.TelegramHelpCommandResponse(
			&managers.HelpCommandConfig{
				MainHelpMessage: s,
			},
			[]string{},
		)

		if assert.NoError(t, err) {
			assert.Equal(t, s, message)
		}
	})
}
