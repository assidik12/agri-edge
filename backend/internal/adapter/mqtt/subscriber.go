package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"agriedge/backend/internal/usecase"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
	client    paho.Client
	topic     string
	ingestUC  *usecase.IngestMetricUseCase
	brokerURL string
}

func NewSubscriber(brokerURL, topic string, ingestUC *usecase.IngestMetricUseCase) *Subscriber {
	return &Subscriber{
		brokerURL: brokerURL,
		topic:     topic,
		ingestUC:  ingestUC,
	}
}

// Start connects to the MQTT broker and subscribes to the telemetry topic
func (s *Subscriber) Start() error {
	opts := paho.NewClientOptions().
		AddBroker(s.brokerURL).
		SetClientID(fmt.Sprintf("agriedge-backend-%d", time.Now().UnixNano())).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(60 * time.Second)

	opts.OnConnect = func(c paho.Client) {
		log.Printf("[MQTT] Connected to broker %s. Subscribing to topic: %s", s.brokerURL, s.topic)
		if token := c.Subscribe(s.topic, 1, s.handleMessage); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] Subscribe error: %v", token.Error())
		}
	}

	opts.OnConnectionLost = func(c paho.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v. Reconnecting...", err)
	}

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	s.client = client
	return nil
}

// Stop gracefully disconnects from MQTT
func (s *Subscriber) Stop() {
	if s.client != nil && s.client.IsConnected() {
		s.client.Disconnect(250)
		log.Println("[MQTT] Disconnected gracefully.")
	}
}

// handleMessage handles incoming MQTT messages matching agriedge/<node_id>/daily_metrics
func (s *Subscriber) handleMessage(client paho.Client, msg paho.Message) {
	topic := msg.Topic()
	payloadBytes := msg.Payload()

	log.Printf("[MQTT] Received message on %s: %s", topic, string(payloadBytes))

	// Topic format: agriedge/<node_id>/daily_metrics
	nodeID, err := extractNodeID(topic)
	if err != nil {
		log.Printf("[MQTT] Failed to extract node_id from topic %s: %v", topic, err)
		return
	}

	var payload usecase.MetricPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		log.Printf("[MQTT] Malformed JSON payload from node %s: %v", nodeID, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.ingestUC.Execute(ctx, nodeID, &payload); err != nil {
		log.Printf("[MQTT] Validation or storage error for node %s: %v", nodeID, err)
		return
	}

	log.Printf("[MQTT] Successfully saved metric for node: %s", nodeID)
}

func extractNodeID(topic string) (string, error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 3 || parts[0] != "agriedge" || parts[2] != "daily_metrics" {
		return "", fmt.Errorf("unexpected topic format: %s", topic)
	}
	return parts[1], nil
}
