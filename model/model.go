package model

// OperatorType is comparison operator enum used to handle rule check
type OperatorType string

// Operator types for comparison
const (
	eq  OperatorType = "=="         // equals to
	ne  OperatorType = "!="         // not equals to
	lt  OperatorType = "<"          // less than
	gt  OperatorType = ">"          // greater than
	le  OperatorType = "<="         // less than equal to
	ge  OperatorType = ">="         // greater than equal to
	co  OperatorType = "contains"   // contains"
	sw  OperatorType = "startsWith" // starts with
	ew  OperatorType = "endsWith"   // ends with
	in  OperatorType = "in"         // in a list
	pr  OperatorType = "pr"         // present
	not OperatorType = "not"        // not of a logical expression
)

// RuleSet holds rule settings struct
type RuleSet struct {
	Name     string      `json:"name"`
	Key      string      `json:"key"` // bounded with the rule
	Priority int         `json:"priority"`
	Value    interface{} `json:"value"`
	Status   bool        `json:"status"`

	Operator OperatorType `json:"operator"`
}

// RuleSetPayload holds rule settings array
type RuleSetPayload struct {
	Data []RuleSet `json:"data"`
}

// RequestPayload holds payload struct which used for fraud detection
type RequestPayload struct {
	Client      RequestClient      `json:"client"`
	Transaction RequestTransaction `json:"data"`
	User        RequestUser        `json:"user"`
}

// RequestTransaction holds transaction struct which used for fraud detection
type RequestTransaction struct {
	Amount          string `json:"amount"`
	TransactionID   string `json:"trx"`
	CardNumber      string `json:"card_number"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	CardHoldersName string `json:"cardholders_name"`
	CVV             string `json:"cvv"`
	ReturnURL       string `json:"return_url"`
}

// RequestUser holds user struct which used for fraud detection
type RequestUser struct {
	UserID      string `json:"userID"`
	Username    string `json:"username"`
	YearOdBirth string `json:"yearofbirth"`
	FullName    string `json:"fullname"`
	Email       string `json:"email"`
	TCKN        string `json:"tckn"`
	IPAddress   string `json:"ip_address"`
}

// RequestClient holds client struct which used for fraud detection
type RequestClient struct {
	Id                       string `json:"id"`
	CCUserPermissionCheck    bool   `json:"cc_user_perm_check"`
	CCHolderAndFullNameMatch bool   `json:"fullname_cc_match"`
}

// RedisConfig holds redis config struct
type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// MySQLConfig holds mysql config struct
type MySQLConfig struct {
	// data source name
	// "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local",
	DSN string `json:"dns"`

	// default size for string fields
	// default: 256
	DefaultStringSize uint `json:"default_string_size"`

	// disable datetime precision, which not supported before MySQL 5.6
	// default: true
	DisableDatetimePrecision bool `json:"disable_datetime_precision"`

	// default datetime precision
	// default: 2
	DefaultDatetimePrecision int `json:"default_datetime_precision"`

	// drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
	// default: true
	DontSupportRenameIndex bool `json:"dont_support_rename_index"`

	// `change` when rename column, rename column not supported before MySQL 8, MariaDB
	// default: true
	DontSupportRenameColumn bool `json:"dont_support_rename_column"`

	// auto configure based on currently MySQL version
	// default: false
	SkipInitializeWithVersion bool `json:"skip_initialize_with_version"`
}

// ResponseType is response enum used to handle response status
type ResponseType string

// PreDefined response types
const (
	SuccessResponse ResponseType = "success"
	FailResponse    ResponseType = "fail"
	ErrorResponse   ResponseType = "error"
)

// ResponsePayload holds response payload used to response all HTTP requests
type ResponsePayload struct {
	Status  ResponseType
	Data    interface{}
	Code    int
	Message string
}

// CreditCardFraud holds related table data
type CreditCardFraud struct {
	ClientID              int64  `json:"client_id"`
	UserID                string `json:"user_id"`
	UserName              string `json:"username"`
	TCKN                  string `json:"tckn"`
	Breached              int64  `json:"breached"`
	PendingCount          int64  `json:"pending_count"`
	Date                  string `json:"date"`
	InitialFifteenCount   int64  `json:"initial_fifteen_count"`
	FifteenCleared        int64  `json:"fifteen_cleared"`
	FifteenNeedsClearance int64  `json:"fifteen_needs_clearance"`
}
