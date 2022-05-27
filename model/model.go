package model

type OperatorType string

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

type RuleSet struct {
	Name     string      `json:"name"`
	Key      string      `json:"key"` // bounded with the rule
	Priority int         `json:"priority"`
	Value    interface{} `json:"value"`
	Status   bool        `json:"status"`

	Operator OperatorType `json:"operator"`
}

type RuleSetPayload struct {
	Data []RuleSet `json:"data"`
}

type RequestPayload struct {
	Transaction RequestTransaction `json:"data"`
	User        RequestUser        `json:"user"`
}

type RequestTransaction struct {
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
	IPAddress   string `json:"ip_address"`
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type MySqlConfig struct {
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
