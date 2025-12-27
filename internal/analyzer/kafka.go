package analyzer

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/navikt/naisalyser/internal/local"
)

// KafkaTopic represents a Kafka topic found in source code
type KafkaTopic struct {
	Name       string
	ConstName  string
	SourceFile string
}

// ExtractKafkaTopics searches for Kafka topic definitions in Kotlin/Java source files
func ExtractKafkaTopics(reader *local.Reader, repoPath string) []KafkaTopic {
	var topics []KafkaTopic

	// Try known paths first
	knownPaths := []string{
		"src/main/kotlin/no/nav/helse/kafka/Topics.kt",
	}

	for _, path := range knownPaths {
		content, err := reader.GetFileContent(repoPath, path)
		if err != nil {
			continue
		}
		topics = append(topics, parseKotlinTopics(content, path)...)
	}

	// Also search in any file containing "Topic" in the kafka directory
	kafkaDir := filepath.Join(repoPath, "src/main/kotlin")
	entries, _ := filepath.Glob(filepath.Join(kafkaDir, "**/kafka/*Topic*.kt"))
	for _, fullPath := range entries {
		relPath, _ := filepath.Rel(repoPath, fullPath)
		content, err := reader.GetFileContent(repoPath, relPath)
		if err != nil {
			continue
		}
		topics = append(topics, parseKotlinTopics(content, relPath)...)
	}

	return topics
}

// parseKotlinTopics extracts topic constants from Kotlin source code
func parseKotlinTopics(content, sourceFile string) []KafkaTopic {
	var topics []KafkaTopic

	// Pattern: const val TOPIC_NAME = "actual.topic.name"
	constPattern := regexp.MustCompile(`const\s+val\s+(\w+)\s*=\s*"([^"]+)"`)
	matches := constPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			constName := match[1]
			topicName := match[2]

			// Filter for likely Kafka topics (contain common Kafka naming patterns)
			if isLikelyKafkaTopic(topicName, constName) {
				topics = append(topics, KafkaTopic{
					Name:       topicName,
					ConstName:  constName,
					SourceFile: sourceFile,
				})
			}
		}
	}

	return topics
}

// isLikelyKafkaTopic checks if a string constant is likely a Kafka topic name
func isLikelyKafkaTopic(topicName, constName string) bool {
	topicNameLower := strings.ToLower(topicName)
	constNameLower := strings.ToLower(constName)

	// Check topic name patterns
	topicPatterns := []string{"privat-", "public-", "aapen-", "-mottatt", "-topic", "dusseldorf."}
	for _, pattern := range topicPatterns {
		if strings.Contains(topicNameLower, pattern) {
			return true
		}
	}

	// Check constant name patterns
	constPatterns := []string{"topic", "kafka"}
	for _, pattern := range constPatterns {
		if strings.Contains(constNameLower, pattern) {
			return true
		}
	}

	return false
}
