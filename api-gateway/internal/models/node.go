package models

import (
	"time"

	"github.com/google/uuid"
)

// ChainType represents the type of blockchain
type ChainType string

const (
	ChainTypeEthereum ChainType = "ethereum"
	ChainTypePolygon  ChainType = "polygon"
	ChainTypeArbitrum ChainType = "arbitrum"
	ChainTypeBSC      ChainType = "bsc"
	ChainTypeCustom   ChainType = "custom"
)

// NodeStatus represents the status of a blockchain node
type NodeStatus string

const (
	NodeStatusRunning     NodeStatus = "running"
	NodeStatusStopped     NodeStatus = "stopped"
	NodeStatusStarting    NodeStatus = "starting"
	NodeStatusSyncing     NodeStatus = "syncing"
	NodeStatusError       NodeStatus = "error"
	NodeStatusMaintenance NodeStatus = "maintenance"
)

// CloudProvider represents the cloud provider where the node is hosted
type CloudProvider string

const (
	CloudProviderAWS          CloudProvider = "aws"
	CloudProviderGCP          CloudProvider = "gcp"
	CloudProviderAzure        CloudProvider = "azure"
	CloudProviderDigitalOcean CloudProvider = "digitalocean"
	CloudProviderOnPremise    CloudProvider = "onpremise"
)

// SyncStatus represents the synchronization status of a blockchain node
type SyncStatus struct {
	IsSyncing          bool    `json:"is_syncing"`
	CurrentBlock       uint64  `json:"current_block"`
	HighestBlock       uint64  `json:"highest_block"`
	StartingBlock      uint64  `json:"starting_block"`
	ProgressPercentage float64 `json:"progress_percentage"`
}

// BlockchainNode represents a blockchain node in the system
type BlockchainNode struct {
	ID                 uuid.UUID              `json:"id"`
	Name               string                 `json:"name"`
	ChainType          ChainType              `json:"chain_type"`
	EndpointURL        string                 `json:"endpoint_url"`
	Status             NodeStatus             `json:"status"`
	Version            string                 `json:"version"`
	SyncStatus         SyncStatus             `json:"sync_status"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	Region             string                 `json:"region"`
	Provider           CloudProvider          `json:"provider"`
	PerformanceMetrics map[string]float64     `json:"performance_metrics,omitempty"`
	Config             map[string]interface{} `json:"config,omitempty"`
}

// CreateNodeRequest is used to create a new blockchain node
type CreateNodeRequest struct {
	Name        string                 `json:"name" binding:"required"`
	ChainType   ChainType              `json:"chain_type" binding:"required"`
	EndpointURL string                 `json:"endpoint_url" binding:"required"`
	Region      string                 `json:"region" binding:"required"`
	Provider    CloudProvider          `json:"provider" binding:"required"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// UpdateNodeRequest is used to update an existing blockchain node
type UpdateNodeRequest struct {
	Name        *string                `json:"name,omitempty"`
	EndpointURL *string                `json:"endpoint_url,omitempty"`
	Status      *NodeStatus            `json:"status,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// ListNodesResponse is the response for listing nodes
type ListNodesResponse struct {
	Items      []BlockchainNode `json:"items"`
	Total      uint64           `json:"total"`
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"page_size"`
	TotalPages uint64           `json:"total_pages"`
}

// NodeResponse is the standard response for node operations
type NodeResponse struct {
	ID uuid.UUID `json:"id"`
}
