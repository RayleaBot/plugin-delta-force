package plugin

import (
	"strings"

	rayleabot "github.com/RayleaBot/RayleaBot/sdk/go"
)

type requestKind int

const (
	requestNone requestKind = iota
	requestLoot
	requestContainerList
	requestPassword
)

type commandRequest struct {
	Kind  requestKind
	Query string
}

func parseCommand(event rayleabot.Event) commandRequest {
	command := strings.TrimSpace(event.Command())
	switch command {
	case "可摸容器":
		return commandRequest{Kind: requestContainerList}
	case "三角洲密码", "每日密码":
		return commandRequest{Kind: requestPassword}
	case "摸":
		return commandRequest{Kind: requestLoot, Query: strings.TrimSpace(strings.Join(event.Args(), ""))}
	}
	if strings.HasPrefix(command, "摸") {
		query := strings.TrimSpace(strings.TrimPrefix(command, "摸"))
		return commandRequest{Kind: requestLoot, Query: query}
	}
	return commandRequest{Kind: requestNone}
}
