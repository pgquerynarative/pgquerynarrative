package story

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ParseNarrative extracts JSON from LLM response and parses it
func ParseNarrative(response string) (*NarrativeContent, error) {
	// Try to extract JSON from the response
	jsonStr := extractJSON(response)

	var narrative NarrativeContent
	if err := json.Unmarshal([]byte(jsonStr), &narrative); err != nil {
		return nil, fmt.Errorf("failed to parse narrative JSON: %w", err)
	}

	// Validate required fields
	if narrative.Headline == "" {
		return nil, fmt.Errorf("narrative missing required field: headline")
	}

	if len(narrative.Takeaways) == 0 {
		return nil, fmt.Errorf("narrative missing required field: takeaways")
	}

	return &narrative, nil
}

// extractJSON extracts JSON from a response that might contain markdown or other text
func extractJSON(response string) string {
	// Remove markdown code blocks
	response = strings.TrimSpace(response)

	// Try to find JSON object
	jsonPattern := regexp.MustCompile(`(?s)\{.*\}`)
	matches := jsonPattern.FindString(response)
	if matches != "" {
		return matches
	}

	// If no JSON found, try the whole response
	return response
}
