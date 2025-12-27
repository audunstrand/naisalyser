package graph

import (
	"github.com/navikt/naisalyser/internal/analyzer"
)

// Node represents a node in the dependency graph
type Node struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // "app", "kafka", "external"
	Namespace string `json:"namespace,omitempty"`
	Stateful  bool   `json:"stateful,omitempty"` // true if app has database/storage
}

// Edge represents a connection between nodes
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // "inbound", "outbound", "database", "kafka", "external"
}

// Graph represents the complete dependency graph
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Builder builds a dependency graph from multiple analysis results
type Builder struct {
	nodes map[string]Node
	edges []Edge
}

// NewBuilder creates a new graph builder
func NewBuilder() *Builder {
	return &Builder{
		nodes: make(map[string]Node),
		edges: []Edge{},
	}
}

// AddFromLocalAnalysis adds nodes and edges from a local analysis result
func (b *Builder) AddFromLocalAnalysis(analysis *analyzer.LocalAnalysis) {
	if analysis == nil || analysis.Repository == nil {
		return
	}

	appName := analysis.Repository.Name

	// Determine if app is stateful (has database or storage)
	stateful := false
	if analysis.NaisConfig != nil && analysis.NaisConfig.Spec.GCP != nil {
		if len(analysis.NaisConfig.Spec.GCP.SQLInstances) > 0 || len(analysis.NaisConfig.Spec.GCP.Buckets) > 0 {
			stateful = true
		}
	}

	// Add app node
	b.addNode(Node{
		ID:       appName,
		Type:     "app",
		Stateful: stateful,
	})

	if analysis.NaisConfig == nil {
		return
	}

	// Process access policy
	if analysis.NaisConfig.Spec.AccessPolicy != nil {
		// Inbound rules - other apps that can call us
		if analysis.NaisConfig.Spec.AccessPolicy.Inbound != nil {
			for _, rule := range analysis.NaisConfig.Spec.AccessPolicy.Inbound.Rules {
				b.addNode(Node{
					ID:        rule.Application,
					Type:      "app",
					Namespace: rule.Namespace,
				})
				b.addEdge(Edge{
					From: rule.Application,
					To:   appName,
					Type: "call",
				})
			}
		}

		// Outbound rules - apps we call
		if analysis.NaisConfig.Spec.AccessPolicy.Outbound != nil {
			for _, rule := range analysis.NaisConfig.Spec.AccessPolicy.Outbound.Rules {
				b.addNode(Node{
					ID:        rule.Application,
					Type:      "app",
					Namespace: rule.Namespace,
				})
				b.addEdge(Edge{
					From: appName,
					To:   rule.Application,
					Type: "call",
				})
			}

			// External hosts
			for _, ext := range analysis.NaisConfig.Spec.AccessPolicy.Outbound.External {
				b.addNode(Node{
					ID:   ext.Host,
					Type: "external",
				})
				b.addEdge(Edge{
					From: appName,
					To:   ext.Host,
					Type: "external",
				})
			}
		}
	}

	// Note: Databases and buckets are handled via stateful flag on app node (set above)

	// Process Kafka topics from source code analysis
	for _, topic := range analysis.KafkaTopics {
		b.addNode(Node{
			ID:   topic.Name,
			Type: "kafka",
		})
		b.addEdge(Edge{
			From: appName,
			To:   topic.Name,
			Type: "kafka",
		})
	}

	// Fallback: if no topics found in code but Kafka pool is configured
	if len(analysis.KafkaTopics) == 0 && analysis.NaisConfig.Spec.Kafka != nil && analysis.NaisConfig.Spec.Kafka.Pool != "" {
		kafkaPool := analysis.NaisConfig.Spec.Kafka.Pool
		b.addNode(Node{
			ID:   kafkaPool,
			Type: "kafka",
		})
		b.addEdge(Edge{
			From: appName,
			To:   kafkaPool,
			Type: "kafka",
		})
	}
}

func (b *Builder) addNode(node Node) {
	if _, exists := b.nodes[node.ID]; !exists {
		b.nodes[node.ID] = node
	}
}

func (b *Builder) addEdge(edge Edge) {
	// Check for duplicate edges
	for _, e := range b.edges {
		if e.From == edge.From && e.To == edge.To && e.Type == edge.Type {
			return
		}
	}
	b.edges = append(b.edges, edge)
}

// Build returns the complete graph
func (b *Builder) Build() *Graph {
	nodes := make([]Node, 0, len(b.nodes))
	for _, node := range b.nodes {
		nodes = append(nodes, node)
	}
	return &Graph{
		Nodes: nodes,
		Edges: b.edges,
	}
}
