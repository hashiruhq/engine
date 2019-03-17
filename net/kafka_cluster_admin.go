package net

import (
	"github.com/Shopify/sarama"
)

// ClusterAdmin manages the cluster configuration
type ClusterAdmin struct {
	ca sarama.ClusterAdmin
}

// NewClusterAdmin creates a new cluster admin manager
func NewClusterAdmin(addrs []string, conf *sarama.Config) (*ClusterAdmin, error) {
	ca, err := sarama.NewClusterAdmin(addrs, conf)
	return &ClusterAdmin{
		ca: ca,
	}, err
}

// EnsureTopicsExist checks if the topics are available
func (admin *ClusterAdmin) EnsureTopicsExist(topics []string) error {
	topicDetails, err := admin.ca.ListTopics()
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if _, ok := topicDetails[topic]; !ok {
			// Setup the Topic details in CreateTopicRequest struct
			topicDetail := &sarama.TopicDetail{}
			topicDetail.NumPartitions = int32(1)
			topicDetail.ReplicationFactor = int16(1)
			topicDetail.ConfigEntries = make(map[string]*string)
			err := admin.ca.CreateTopic(topic, topicDetail, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
