package graph

import (
	"github.com/navikt/naisalyser/internal/analyzer"
)

// NodeType represents the type of a graph node
type NodeType string

const (
	NodeTypeApp      NodeType = "app"
	NodeTypeKafka    NodeType = "kafka"
	NodeTypeExternal NodeType = "external"
)

// EdgeType represents the type of a graph edge
type EdgeType string

const (
	EdgeTypeCall     EdgeType = "call"
	EdgeTypeKafka    EdgeType = "kafka"
	EdgeTypeExternal EdgeType = "external"
)

// Node represents a node in the dependency graph
type Node struct {
	ID        string   `json:"id"`
	Type      NodeType `json:"type"`
	Namespace string   `json:"namespace,omitempty"`
	Stateful  bool     `json:"stateful,omitempty"` // true if app has database/storage
}

// Edge represents a connection between nodes
type Edge struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	Type EdgeType `json:"type"`
}

// Graph represents the complete dependency graph
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// edgeKey is used for O(1) duplicate detection
type edgeKey struct {
	From string
	To   string
	Type EdgeType
}

// Builder builds a dependency graph from multiple analysis results
type Builder struct {
	nodes   map[string]Node
	edges   []Edge
	edgeSet map[edgeKey]struct{} // For O(1) duplicate detection
}

// NewBuilder creates a new graph builder
func NewBuilder() *Builder {
	return &Builder{
		nodes:   make(map[string]Node),
		edges:   []Edge{},
		edgeSet: make(map[edgeKey]struct{}),
	}
}

// AddFromLocalAnalysis adds nodes and edges from a local analysis result
func (b *Builder) AddFromLocalAnalysis(analysis *analyzer.LocalAnalysis) {
	if analysis == nil || analysis.Repository == nil {
		return
	}

	appName := analysis.Repository.Name

	// Add app node - NaisConfig.IsStateful() handles nil checks internally
	b.addNode(Node{
		ID:       appName,
		Type:     NodeTypeApp,
		Stateful: analysis.NaisConfig.IsStateful(),
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
					Type:      NodeTypeApp,
					Namespace: rule.Namespace,
				})
				b.addEdge(Edge{
					From: rule.Application,
					To:   appName,
					Type: EdgeTypeCall,
				})
			}
		}

		// Outbound rules - apps we call
		if analysis.NaisConfig.Spec.AccessPolicy.Outbound != nil {
			for _, rule := range analysis.NaisConfig.Spec.AccessPolicy.Outbound.Rules {
				b.addNode(Node{
					ID:        rule.Application,
					Type:      NodeTypeApp,
					Namespace: rule.Namespace,
				})
				b.addEdge(Edge{
					From: appName,
					To:   rule.Application,
					Type: EdgeTypeCall,
				})
			}

			// External hosts
			for _, ext := range analysis.NaisConfig.Spec.AccessPolicy.Outbound.External {
				b.addNode(Node{
					ID:   ext.Host,
					Type: NodeTypeExternal,
				})
				b.addEdge(Edge{
					From: appName,
					To:   ext.Host,
					Type: EdgeTypeExternal,
				})
			}
		}
	}

	// Note: Databases and buckets are handled via stateful flag on app node (set above)

	// Process Kafka topics from source code analysis
	for _, topic := range analysis.KafkaTopics {
		b.addNode(Node{
			ID:   topic.Name,
			Type: NodeTypeKafka,
		})
		b.addEdge(Edge{
			From: appName,
			To:   topic.Name,
			Type: EdgeTypeKafka,
		})
	}

	// Fallback: if no topics found in code but Kafka pool is configured
	if len(analysis.KafkaTopics) == 0 && analysis.NaisConfig.Spec.Kafka != nil && analysis.NaisConfig.Spec.Kafka.Pool != "" {
		kafkaPool := analysis.NaisConfig.Spec.Kafka.Pool
		b.addNode(Node{
			ID:   kafkaPool,
			Type: NodeTypeKafka,
		})
		b.addEdge(Edge{
			From: appName,
			To:   kafkaPool,
			Type: EdgeTypeKafka,
		})
	}
}

func (b *Builder) addNode(node Node) {
	if _, exists := b.nodes[node.ID]; !exists {
		b.nodes[node.ID] = node
	}
}

func (b *Builder) addEdge(edge Edge) {
	key := edgeKey{From: edge.From, To: edge.To, Type: edge.Type}
	if _, exists := b.edgeSet[key]; exists {
		return
	}
	b.edgeSet[key] = struct{}{}
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
