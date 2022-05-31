package rules

import (
	"errors"
	"fmt"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"strings"
	"sync"
)

type requestPayload struct {
	Data model.RequestPayload `json:"data"`
	sync.RWMutex
}

var instance requestPayload

// GetRequestPayloadInstance constructs request payload instance
func GetRequestPayloadInstance() *requestPayload {
	return &instance
}

func (payload *requestPayload) SetPayload(data model.RequestPayload) {
	payload.Lock()
	defer payload.Unlock()
	payload.Data = data
}

func (payload *requestPayload) GetPayload() model.RequestPayload {
	payload.RLock()
	defer payload.RUnlock()
	return payload.Data
}

func (payload *requestPayload) ProcessRules() (bool, error) {
	ruleSets := rulesets.GetInstance().GetPayload()
	isOK := false
	var err error
	if !anyRuleExists(ruleSets) {
		return false, errors.New("please, define your rule sets first")
	}

	if isOK, err = payload.checkCardBIN(); !isOK || err != nil {
		return false, fmt.Errorf("%v check is failed!\nError Details: %v", "checkCardBIN", err)
	}

	if _, err = payload.checkThreeUniqueCardsAllowed(); !isOK || err != nil {
		return false, fmt.Errorf("%v check is failed!\nError Details: %v", "checkThreeUniqueCardsAllowed", err)
	}

	for _, ruleSet := range ruleSets {
		switch ruleSet.Key {
		case "PendingCountThreshold":
			isOK, err = payload.checkPendingCountThreshold()
		case "PendingAllowanceByTimeInterval":
			isOK, err = payload.checkPendingAllowanceByTimeInterval()
		case "ApprovedAllowanceByTimeInterval":
			isOK, err = payload.checkApprovedAllowanceByTimeInterval()
		case "MaxDailyAllowancePerUser":
			isOK, err = payload.checkMaxDailyAllowancePerUser()
		case "MinTransactionAmount":
			isOK, err = payload.checkMinTransactionAmount()
		case "MaxTransactionAmount":
			isOK, err = payload.checkMaxTransactionAmount()
			utils.CheckError(err)
		}

		if !isOK || err != nil {
			return false, fmt.Errorf("%v check is failed!\nError Details: %v", ruleSet.Key, err)
		}
	}
	return true, nil
}

func (payload *requestPayload) checkCardBIN() (bool, error) {
	binExists := false
	cardNumber := &payload.Data.Transaction.CardNumber
	*cardNumber = strings.Trim(*cardNumber, " ")
	if *cardNumber == "" {
		return binExists, errors.New("card number is empty")
	}

	cardBin := (*cardNumber)[:6]
	tx := config.MySQLDb.Raw("SELECT 1 as binExists FROM cc_binlist WHERE card_bin = ? LIMIT 1", cardBin).
		Scan(&binExists)

	if tx.Error != nil || !binExists {
		errString := "card issuer is not listed in the bin list!"
		if tx.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", tx.Error)
		}

		return binExists, errors.New(errString)
	}

	return true, nil
}

func (payload *requestPayload) getCardBINIca() (string, error) {
	bankIca := ""
	cardNumber := &payload.Data.Transaction.CardNumber
	*cardNumber = strings.Trim(*cardNumber, " ")
	if *cardNumber == "" {
		return bankIca, errors.New("card number is empty")
	}

	cardBin := (*cardNumber)[:6]
	tx := config.MySQLDb.Raw("SELECT bankIca as binExists FROM cc_binlist WHERE card_bin = ? LIMIT 1", cardBin).
		Scan(&bankIca)

	if tx.Error != nil || bankIca == "" {
		errString := "card issuer is not listed in the bin list!"
		if tx.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", tx.Error)
		}

		return bankIca, errors.New(errString)
	}

	return bankIca, nil
}

func (payload *requestPayload) checkThreeUniqueCardsAllowed() (bool, error) {
	cardNumber := &payload.Data.Transaction.CardNumber
	*cardNumber = strings.Trim(*cardNumber, " ")
	if *cardNumber == "" {
		return false, errors.New("card number is empty")
	}
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("tckn is empty")
	}
	var rowCount int64
	var cryptedCCs []string
	tx := config.MySQLDb.Raw(`SELECT rjr.crypted_cc FROM request_jetpay_registrations AS rjr 
			INNER JOIN request AS r ON rjr.request_id = r.ID 
			WHERE created_at >= CAST(CURDATE() AS DATETIME) 
			AND created_at <= DATE_SUB(CAST(DATE_ADD(CURDATE(), INTERVAL 1 DAY) AS DATETIME), INTERVAL 1 SECOND) 
			AND user_tckn = ? AND r.Status = 1 GROUP BY rjr.crypted_cc`, 3, *tckn).Scan(&cryptedCCs).Count(&rowCount)

	cardAllowance := rowCount <= 3
	if tx.Error != nil || !cardAllowance {
		errString := "daily card allowance limit reached!"
		if tx.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", tx.Error)
		}

		return false, errors.New(errString)
	}

	cryptedCard := utils.GetMD5Hash(utils.GetMD5Hash(*cardNumber))
	for _, cryptedCC := range cryptedCCs {
		if cryptedCard == cryptedCC {
			return true, nil
		}
	}

	return true, nil
}

func (payload *requestPayload) checkFifteenCountClearance() (bool, error) {
	fifteenNeedsClearance := false
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.ClientID
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *userID == "" || *clientID == "" {
		return false, errors.New("userID and/or clientID is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT fifteenNeedsClearance FROM cc_fraud WHERE user_id = ? AND client_id = ?`, *userID, *clientID).
		Scan(&fifteenNeedsClearance)

	if tx.Error != nil || fifteenNeedsClearance {
		errString := "kredi Kartı harcama güvenliğinizin sağlanması kapsamında, sitemizin çağrı merkezi ile iletişime geçerek işlemlerin sizin tarafınızdan yapıldığını doğrulamanız gerekmektedir!"
		if tx.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", tx.Error)
		}

		return false, errors.New(errString)
	}

	return true, nil
}

func (payload *requestPayload) checkOneApprovedAllowedByThirtyMinuteInterval() (bool, error) {
	allowance := false
	tckn := &payload.Data.User.TCKN
	userID := &payload.Data.User.UserID
	sid := 0
	fullName := &payload.Data.User.FullName
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("tckn is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT
      COUNT(1) == 0 AS allowance
	  FROM
		  request r
		  INNER JOIN request_jetpay_registrations rjr ON rjr.request_id = r.ID
	  WHERE Status = 1 AND payment_method = 5 AND (StartDate > DATE_SUB(NOW(), INTERVAL 30 MINUTE)) AND 
		((r.SID = ? AND r.UserID = ?) OR (r.FullName = ? AND rjr.user_tckn = ?))`, sid, *userID, fullName, *tckn).
		Scan(&allowance)

	if tx.Error != nil || !allowance {
		errString := "her kullanıcı 30 dakikada bir adet başarılı işlem gerçekleştirebilir!"
		if tx.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", tx.Error)
		}

		return false, errors.New(errString)
	}

	return true, nil
}

func getUserFraudRecordExternal(clientID *string, userID *string) (model.CreditCardFraud, error) {
	creditCardFraud := model.CreditCardFraud{}
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return creditCardFraud, errors.New("clientID and/or userID is empty")
	}

	tx := config.MySQLDb.Raw(`"SELECT * FROM cc_fraud WHERE client_id = ? AND user_id = ?`, *clientID, *userID).
		Scan(&creditCardFraud)

	if tx.Error != nil {
		return creditCardFraud, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	return creditCardFraud, nil
}

func incrementFifteenCount(clientID *string, userID *string) (bool, error) {
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	fraudRecord, err := getUserFraudRecordExternal(clientID, userID)
	if err != nil {
		return false, fmt.Errorf("\nError Details: %v", err.Error())
	}

	if fraudRecord.FifteenCleared == 1 {
		return true, nil
	}

	tx := config.MySQLDb.Exec(`"UPDATE cc_fraud SET initial_fifteen_count += 1 WHERE client_id = ? AND user_id = ?`, *clientID, *userID)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	if fraudRecord.InitialFifteenCount == 14 {
		if err := changeUserPermExternal(clientID, userID, 0); err != nil {
			return false, fmt.Errorf("error Details: %v", err)
		}
		tx := config.MySQLDb.Exec(`"UPDATE cc_fraud SET fifteen_needs_clearance = 1 WHERE client_id = ? AND user_id = ?`, *clientID, *userID)
		if tx.Error != nil {
			return false, fmt.Errorf("error Details: %v", tx.Error)
		}
	}

	return true, nil
}

func changeUserPermExternal(clientID *string, userID *string, privillage int64) error {
	tx := config.MySQLDb.Exec(`"UPDATE cc_client_users SET privilege = ? WHERE client_id = ? AND user_id = ?`, privillage, *clientID, *userID)
	if tx.Error != nil {
		return fmt.Errorf("\nError Details: %v", tx.Error)
	}
	return nil
}

func (payload *requestPayload) changeUserPerm(privillage string) error {
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.ClientID
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	tx := config.MySQLDb.Exec(`"UPDATE cc_client_users SET privilege = ? WHERE client_id = ? AND user_id = ?`, privillage, *clientID, *userID)
	if tx.Error != nil {
		return fmt.Errorf("\nError Details: %v", tx.Error)
	}
	return nil
}

func (payload *requestPayload) getUserFraudRecord() (model.CreditCardFraud, error) {
	creditCardFraud := model.CreditCardFraud{}
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.ClientID
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return creditCardFraud, errors.New("clientID and/or userID is empty")
	}

	tx := config.MySQLDb.Raw(`"SELECT * FROM cc_fraud WHERE client_id = ? AND user_id = ?`, *clientID, *userID).
		Scan(&creditCardFraud)

	if tx.Error != nil {
		return creditCardFraud, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	return creditCardFraud, nil
}

func (payload *requestPayload) checkOneCardPerBank() (bool, error) {
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	cardNumber := &payload.Data.Transaction.CardNumber
	*cardNumber = strings.Trim(*cardNumber, " ")
	if *tckn == "" {
		return false, errors.New("card number is empty")
	}
	bankIca, err := payload.getCardBINIca()
	if err != nil {
		return false, err
	}

	cryptedCC := ""
	tx := config.MySQLDb.Raw(`SELECT ccb.bank_ica , rjr.crypted_cc AS crypted_cc
        FROM request_jetpay_registrations AS rjr 
			INNER JOIN request AS r ON rjr.request_id = r.ID 
			INNER JOIN cc_binlist AS ccb ON ccb.card_bin = rjr.card_bin 
        AND rjr.created_at >= ? 
		AND user_tckn = ? 
		AND ccb.bank_ica = ? 
		AND r.Status = 1 GROUP BY ccb.bank_ica, rjr.crypted_cc LIMIT 1`,
		"\"2022-04-05 08:00:00\"", *tckn, bankIca).Scan(&cryptedCC)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	cryptedCard := utils.GetMD5Hash(utils.GetMD5Hash(*cardNumber))
	if cryptedCC == cryptedCard {
		return false, errors.New("only one unique card per bank is permitted")
	}

	return true, nil
}

func (payload *requestPayload) checkOneTcknPerUser() (bool, error) {
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("card number is empty")
	}
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.ClientID
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	recTCKN := ""
	tx := config.MySQLDb.Raw(`SELECT tckn FROM cc_fraud WHERE user_id = ? AND client_id =  LIMIT 1`,
		*userID, clientID).Scan(&recTCKN)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	if recTCKN == *tckn {
		return false, errors.New("user is only allowed to perform transactions with a single tckn")
	}

	return true, nil
}

func (payload *requestPayload) checkPendingCountThreshold() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkPendingAllowanceByTimeInterval() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkApprovedAllowanceByTimeInterval() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMaxDailyAllowancePerUser() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMinTransactionAmount() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMaxTransactionAmount() (bool, error) {

	return true, nil
}

func anyRuleExists(ruleSets []model.RuleSet) bool {
	preDefinedRuleKeys := []string{
		"PendingCountThreshold",
		"PendingAllowanceByTimeInterval",
		"ApprovedAllowanceByTimeInterval",
		"MaxDailyAllowancePerUser",
		"MinTransactionAmount",
		"MaxTransactionAmount",
	}
	var ruleKeys []string
	for _, ruleSet := range ruleSets {
		ruleKeys = append(ruleKeys, ruleSet.Key)
	}

	return len(utils.Intersection(ruleKeys, preDefinedRuleKeys)) > 0
}
