package models

import (
	"time"
)

// normally we would use gorm.Model but we need to filter out unnecessary json fields
type DBCustomModel struct {
	// Unique identifier for the record.
	ID uint `gorm:"primary_key" json:"-"`
	// Timestamp indicating when the record was created.
	CreatedAt time.Time `json:"-"`
	// Timestamp indicating when the record was last updated.
	UpdatedAt time.Time `json:"-"`
	// Optional timestamp indicating when the record was deleted. Indexed for efficient querying.
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type DBContainerMetricsModel struct {
	DBCustomModel
	// Name of the container.
	Name string `json:"name" example:"nginx"`
	// CPU usage of the container in millicores (1/1000th of a core).
	CPUUsage int64 `json:"cpu_usage" example:"20"`
	// Memory usage of the container in bytes.
	MemoryUsage int64 `json:"memory_usage" example:"10485760"`
	// Foreign key that references DBPodMetricsModel's ID field to make the relationship between containers and pods.
	PodID uint `json:"-"`
}

type DBPodMetricsModel struct {
	DBCustomModel
	// Name of the pod.
	Name string `json:"name" example:"mypod"`
	// Metrics records grouped by containers.
	Containers []DBContainerMetricsModel `gorm:"foreignKey:PodID" json:"containers"`
	// Foreign key that references ClusterMetricsModel's ID field to make the relationship between pods and clusters.
	CRecID uint `json:"-"`
}

type DBClusterMetricsModel struct {
	DBCustomModel
	// Timestamp indicating when the record was created. This is the time when the metrics were collected.
	Pods []DBPodMetricsModel `gorm:"foreignKey:CRecID" json:"pods"`
}
