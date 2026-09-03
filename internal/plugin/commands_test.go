package plugin

import (
	"testing"

	rayleabot "github.com/RayleaBot/RayleaBot/sdk/go"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name  string
		event rayleabot.Event
		kind  requestKind
		query string
	}{
		{name: "compact loot", event: commandEvent("摸航空箱"), kind: requestLoot, query: "航空箱"},
		{name: "spaced loot", event: commandEventWithArgs("摸", "航空箱"), kind: requestLoot, query: "航空箱"},
		{name: "empty loot opens list", event: commandEvent("摸"), kind: requestLoot},
		{name: "container list", event: commandEvent("可摸容器"), kind: requestContainerList},
		{name: "password primary", event: commandEvent("三角洲密码"), kind: requestPassword},
		{name: "password alias", event: commandEvent("每日密码"), kind: requestPassword},
		{name: "unrelated", event: commandEvent("天气"), kind: requestNone},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseCommand(test.event)
			if got.Kind != test.kind || got.Query != test.query {
				t.Fatalf("parseCommand() = %#v, want kind=%v query=%q", got, test.kind, test.query)
			}
		})
	}
}

func commandEvent(command string) rayleabot.Event {
	return rayleabot.Event{Payload: map[string]any{"command": command}}
}

func commandEventWithArgs(command string, args ...string) rayleabot.Event {
	values := make([]any, 0, len(args))
	for _, arg := range args {
		values = append(values, arg)
	}
	return rayleabot.Event{Payload: map[string]any{"command": command, "args": values}}
}
