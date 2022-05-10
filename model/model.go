package model

type RuleSet struct {
	Name     string      `json:"name"`
	Key      string      `json:"key"`
	Priority int         `json:"priority"`
	Value    interface{} `json:"value"`
	Status   bool        `json:"status"`

	// TODO - make Enum
	Operator string `json:"operator"`
}

type RuleSetPayload struct {
	Data []RuleSet `json:"data"`
}

type RequestPayload struct {
	Data RequestData `json:"data"`
	User RequestUser `json:"user"`
}

type RequestData struct {
	Amount          string `json:"amount"`
	TransactionId   string `json:"trx"`
	CardNumber      string `json:"card_number"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	CardHoldersName string `json:"cardholders_name"`
	CVV             string `json:"cvv"`
	ReturnURL       string `json:"return_url"`
}

type RequestUser struct {
	Username    string `json:"username"`
	UserId      string `json:"userID"`
	YearOdBirth string `json:"yearofbirth"`
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	TCKN        string `json:"tckn"`
	CVV         string `json:"cvv"`
	IPAddress   string `json:"ip_address"`
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type ResponseType string

const (
	SuccessResponse ResponseType = "success"
	FailResponse    ResponseType = "fail"
	ErrorResponse   ResponseType = "error"
)

type ResponsePayload struct {
	Status  ResponseType
	Data    interface{}
	Code    int
	Message string
}