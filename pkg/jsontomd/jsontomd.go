package jsontomd

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types/events"
)

func ConvertJsonToMD(jsonData string) (string, error) {
	result := "# Docker Event \n\n ## Event Details \n\n"
	var dockerMsg events.Message
	err := json.Unmarshal([]byte(jsonData), &dockerMsg)
	if err != nil {
		slog.Error("Error in unmarshaling the message in the redis", "error", err)
		return "", err
	}

	result += fmt.Sprintf("- **Type**: `%s`\n", dockerMsg.Type)
	result += fmt.Sprintf("- **Action**: `%s`\n", dockerMsg.Action)
	result += fmt.Sprintf("- **Scope**: `%s`\n", dockerMsg.Scope)
	result += fmt.Sprintf("- **Time**: `%d`\n", dockerMsg.Time)
	result += fmt.Sprintf("- **TimeNano**: `%d`\n", dockerMsg.TimeNano)
	result += "## Actor \n"
	result += fmt.Sprintf("- **Actor.ID**: `%s`\n", dockerMsg.Actor.ID)

	for k, v := range dockerMsg.Actor.Attributes {
		result += fmt.Sprintf("  - **%s**: `%s`\n", k, v)
	}
	return result, nil
}
