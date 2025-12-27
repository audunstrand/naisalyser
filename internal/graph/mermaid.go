package graph

import (
	"fmt"
	"strings"
)

// ToMermaid generates a Mermaid diagram from the graph
func (g *Graph) ToMermaid() string {
	var sb strings.Builder

	sb.WriteString("graph LR\n")

	// Define node styles by type
	sb.WriteString("    %% Node definitions\n")
	for _, node := range g.Nodes {
		nodeID := sanitizeMermaidID(node.ID)
		switch node.Type {
		case NodeTypeApp:
			if node.Stateful {
				// Cylinder shape for stateful apps (with database)
				sb.WriteString(fmt.Sprintf("    %s[(%s)]\\n", nodeID, node.ID))
			} else {
				sb.WriteString(fmt.Sprintf("    %s[%s]\\n", nodeID, node.ID))
			}
		case NodeTypeKafka:
			sb.WriteString(fmt.Sprintf("    %s{{%s}}\\n", nodeID, node.ID))
		case NodeTypeExternal:
			sb.WriteString(fmt.Sprintf("    %s>%s]\\n", nodeID, node.ID))
		default:
			sb.WriteString(fmt.Sprintf("    %s[%s]\\n", nodeID, node.ID))
		}
	}

	sb.WriteString("\n    %% Edges\n")
	for _, edge := range g.Edges {
		fromID := sanitizeMermaidID(edge.From)
		toID := sanitizeMermaidID(edge.To)

		switch edge.Type {
		case EdgeTypeCall:
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromID, toID))
		case EdgeTypeKafka:
			sb.WriteString(fmt.Sprintf("    %s -.->|kafka| %s\n", fromID, toID))
		case EdgeTypeExternal:
			sb.WriteString(fmt.Sprintf("    %s -->|ext| %s\n", fromID, toID))
		default:
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromID, toID))
		}
	}

	// Add styling
	sb.WriteString("\n    %% Styling\n")
	sb.WriteString("    classDef stateless fill:#4A90D9,stroke:#2E5A8C,color:white\n")
	sb.WriteString("    classDef stateful fill:#F5A623,stroke:#C47A00,color:white\n")
	sb.WriteString("    classDef kafka fill:#7B68EE,stroke:#5A4DB2,color:white\n")
	sb.WriteString("    classDef external fill:#808080,stroke:#404040,color:white\n")

	// Apply styles to nodes
	var stateless, stateful, kafkas, externals []string
	for _, node := range g.Nodes {
		id := sanitizeMermaidID(node.ID)
		switch node.Type {
		case NodeTypeApp:
			if node.Stateful {
				stateful = append(stateful, id)
			} else {
				stateless = append(stateless, id)
			}
		case NodeTypeKafka:
			kafkas = append(kafkas, id)
		case NodeTypeExternal:
			externals = append(externals, id)
		}
	}

	if len(stateless) > 0 {
		sb.WriteString(fmt.Sprintf("    class %s stateless\n", strings.Join(stateless, ",")))
	}
	if len(stateful) > 0 {
		sb.WriteString(fmt.Sprintf("    class %s stateful\n", strings.Join(stateful, ",")))
	}
	if len(kafkas) > 0 {
		sb.WriteString(fmt.Sprintf("    class %s kafka\n", strings.Join(kafkas, ",")))
	}
	if len(externals) > 0 {
		sb.WriteString(fmt.Sprintf("    class %s external\n", strings.Join(externals, ",")))
	}

	return sb.String()
}

// sanitizeMermaidID converts a string to a valid Mermaid node ID
func sanitizeMermaidID(s string) string {
	// Replace characters that are problematic in Mermaid
	result := strings.ReplaceAll(s, "-", "_")
	result = strings.ReplaceAll(result, ".", "_")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.ReplaceAll(result, ":", "_")
	return result
}
