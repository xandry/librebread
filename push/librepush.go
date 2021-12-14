package push

import "fmt"

type Storage interface {
	AddBatchMessage(msg BatchMessage) error
	AllMessages() ([]SentMessage, error)
	ByID(id string) (SentMessage, error)
}

type BatchMessage struct {
	ID           int64             `json:"id"`
	PushService  string            `json:"push_service"`
	Title        string            `json:"title"`
	Text         string            `json:"text"`
	Data         map[string]string `json:"data,omitempty"`
	TTL          int64             `json:"ttl"`
	Tokens       []string          `json:"tokens"`
	ValidateOnly bool              `json:"validate_only"`
}

type SendResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Error     error  `json:"error"`
}

type BatchResponse struct {
	SuccessCount int             `json:"success_count"`
	FailureCount int             `json:"failure_count"`
	Responses    []*SendResponse `json:"responses"`
}

type LibrePush struct {
	storage Storage
}

func NewLibrePush(storage Storage) *LibrePush {
	return &LibrePush{
		storage: storage,
	}
}

func (p *LibrePush) Send(msg BatchMessage) (*BatchResponse, error) {
	err := p.storage.AddBatchMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("can not save batch message to storage: %w", err)
	}

	var response BatchResponse

	allSuccess(msg.Tokens, &response)

	return &response, nil
}

func (p *LibrePush) SendDryRun(msg BatchMessage) (*BatchResponse, error) {
	var response BatchResponse

	allSuccess(msg.Tokens, &response)

	return &response, nil
}

func allSuccess(tokens []string, response *BatchResponse) {
	response.SuccessCount = len(tokens)
	response.FailureCount = 0

	for _, token := range tokens {
		response.Responses = append(response.Responses, &SendResponse{
			Success:   true,
			MessageID: token,
			Error:     nil,
		})
	}
}
