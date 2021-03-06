package rules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type requestPayload struct {
	Data model.RequestPayload `json:"data"`
	sync.RWMutex
}

const MethodID = 5

var requestPayloadInstance requestPayload
var activeRules = rulesets.GetInstance().GetPayloadKeyMapping()

// GetRequestPayloadInstance constructs request payload instance
func GetRequestPayloadInstance() *requestPayload {
	return &requestPayloadInstance
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
	tx := config.MySQLDb.Raw("SELECT bank_ica as bankIca FROM cc_binlist WHERE card_bin = ? LIMIT 1", cardBin).
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
	var cryptedCCs []string
	tx := config.MySQLDb.Raw(`SELECT rjr.crypted_cc 
		FROM request_jetpay_registrations AS rjr 
			INNER JOIN request AS r ON rjr.request_id = r.ID 
		WHERE created_at >= CAST(CURDATE() AS DATETIME) 
			AND created_at <= DATE_SUB(CAST(DATE_ADD(CURDATE(), INTERVAL 1 DAY) AS DATETIME), INTERVAL 1 SECOND) 
			AND user_tckn = ? AND r.Status = 1 GROUP BY rjr.crypted_cc`, *tckn).Scan(&cryptedCCs)

	cardAllowance := len(cryptedCCs) <= 3
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
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *userID == "" || *clientID == "" {
		return false, errors.New("userID and/or clientID is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT fifteen_needs_clearance AS fifteenNeedsClearance FROM cc_fraud WHERE user_id = ? AND client_id = ?`, *userID, *clientID).
		Scan(&fifteenNeedsClearance)

	if tx.Error != nil || fifteenNeedsClearance {
		errString := "kredi Kart?? harcama g??venli??inizin sa??lanmas?? kapsam??nda, sitemizin ??a??r?? merkezi ile ileti??ime ge??erek i??lemlerin sizin taraf??n??zdan yap??ld??????n?? do??rulaman??z gerekmektedir!"
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
	clientID := &payload.Data.Client.Id
	fullName := &payload.Data.User.FullName
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("tckn is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT
      COUNT(1) = 0 AS allowance
	  FROM request r
		  INNER JOIN request_jetpay_registrations rjr ON rjr.request_id = r.ID
	  WHERE Status = 1 AND payment_method = 5 AND (StartDate > DATE_SUB(NOW(), INTERVAL 30 MINUTE)) 
		AND ((r.SID = ? AND r.UserID = ?) OR (r.FullName = ? AND rjr.user_tckn = ?))`, clientID, *userID, fullName, *tckn).
		Scan(&allowance)

	if tx.Error != nil || !allowance {
		errString := "her kullan??c?? 30 dakikada bir adet ba??ar??l?? i??lem ger??ekle??tirebilir!"
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
	clientID := &payload.Data.Client.Id
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
	clientID := &payload.Data.Client.Id
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
				AND r.Status = 1 
		GROUP BY ccb.bank_ica, rjr.crypted_cc LIMIT 1`,
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
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	recTCKN := ""
	tx := config.MySQLDb.Raw(`SELECT tckn FROM cc_fraud WHERE user_id = ? AND client_id = ?  LIMIT 1`,
		*userID, *clientID).Scan(&recTCKN)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	if recTCKN == *tckn {
		return false, errors.New("user is only allowed to perform transactions with a single tckn")
	}

	return true, nil
}

func (payload *requestPayload) checkLastTenTransactions() (bool, error) {
	var txPaymentMethods []int64
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT CAST(IFNULL(payment_method, 0) AS UNSIGNED) AS payment_method 
		FROM request WHERE UserID = ? AND SID = ? ORDER BY StartDate DESC LIMIT 10`,
		*userID, *clientID).Scan(&txPaymentMethods)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}
	ccTxCount := 0
	for _, v := range txPaymentMethods {
		if v == MethodID {
			ccTxCount++
		}
	}

	if ccTxCount == 10 {
		if err := payload.changeUserPerm("0"); err != nil {
			return false, err
		}

		if err := clearPendingCount(*clientID, *userID); err != nil {
			return false, err
		}

		return false, fmt.Errorf("user's deposit privilege is revoked. Please contact live support")
	}

	return true, nil
}

func clearPendingCount(clientId string, userId string) error {
	tx := config.MySQLDb.Exec(`UPDATE cc_fraud SET pending_count = 0 WHERE user_id = ? AND client_id = ?`,
		clientId, userId)
	if tx.Error != nil {
		return fmt.Errorf("\nError Details: %v", tx.Error)
	}

	return nil
}

func (payload *requestPayload) createBlacklistRecord() error {
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return errors.New("card number is empty")
	}
	userName := strings.Trim(payload.Data.User.Username, " ")
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	fullName := utils.SanitizeName(payload.Data.User.FullName)
	if *clientID == "" || *userID == "" {
		return errors.New("clientID and/or userID is empty")
	}
	if tx := config.MySQLDb.Exec(`INSERT INTO blacklist (client_id, tckn, username, fullname, notes) VALUES (?, ?, ?, ?, ?)`,
		*clientID, *tckn, userName, fullName, "Ba??ar??s??z i??lem denemesi!"); tx.Error != nil {
		return fmt.Errorf("\nError Details: %v", tx.Error)
	}

	return nil
}

func (payload *requestPayload) checkPendingCountThreshold() (bool, error) {

	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("card number is empty")
	}
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	txCount := int64(0)
	if tx := config.MySQLDb.Raw(`SELECT COUNT(1) AS txCount
		FROM cc_fraud WHERE user_id = ? AND client_id = ? AND tckn = ?`, *userID, *clientID, *tckn).Scan(&txCount); tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}
	if txCount > config.PendingCountThreshold {
		clearPendingCount(*clientID, *userID)
		payload.createBlacklistRecord()
		return false, fmt.Errorf("user blacklisted")
	}
	return true, nil
}

func (payload *requestPayload) checkPendingAllowanceByTimeInterval() (bool, error) {
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("card number is empty")
	}
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	fullName := utils.SanitizeName(payload.Data.User.FullName)
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	txCount := int64(0)
	if tx := config.MySQLDb.Raw(`SELECT COUNT(1) AS txCount
		FROM request r
			INNER JOIN request_jetpay_registrations rjr ON rjr.request_id = r.ID
      	WHERE Status = 0 AND payment_method = 5 AND StartDate > DATE_SUB(NOW(), INTERVAL ? MINUTE)
           	AND ((SID = ? AND UserID = ?) OR (r.FullName = ? AND rjr.user_tckn = ?))
      	ORDER BY StartDate DESC;`, config.PendingAllowanceByTimeInterval, *clientID, *userID, fullName, *tckn).Scan(&txCount); tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}
	if txCount > config.PendingAllowanceByTimeInterval {
		return false, fmt.Errorf("user already has a pending transaction")
	}
	return true, nil
}

// TODO
func (payload *requestPayload) checkApprovedAllowanceByTimeInterval() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMaxDailyAllowancePerUser() (bool, error) {
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("card number is empty")
	}
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	fullName := utils.SanitizeName(payload.Data.User.FullName)
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	txCount := int64(0)
	if tx := config.MySQLDb.Raw(`SELECT COUNT(1) AS txCount
		FROM request r
              	INNER JOIN request_jetpay_registrations rjr ON rjr.request_id = r.ID
      	WHERE Status = 1 AND payment_method = 5 AND StartDate >= CAST(CURDATE() AS DATETIME) 
			AND StartDate <= DATE_SUB(CAST(DATE_ADD(CURDATE(), INTERVAL 1 DAY) AS DATETIME), INTERVAL 1 SECOND)
           	AND ((SID = ? AND UserID = ?) OR (r.FullName = ? AND rjr.user_tckn = ?))
      	ORDER BY StartDate DESC;`, *clientID, *userID, fullName, *tckn).Scan(&txCount); tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}
	if txCount >= config.MaxDailyAllowancePerUser {
		return false, fmt.Errorf("user exceeds daily approved transaction count limits")
	}
	return true, nil
}

func (payload *requestPayload) checkMinTransactionAmount() (bool, error) {
	txAmount, err := strconv.ParseFloat(payload.Data.Transaction.Amount, 64)
	if err != nil {
		return false, fmt.Errorf("unexpacted amount value")
	}
	if txAmount < config.MinTransactionAmount {
		return false, fmt.Errorf("amount below minimum limits")
	}
	return true, nil
}

func (payload *requestPayload) checkMaxTransactionAmount() (bool, error) {
	fullName := utils.LocalizeToEnglish(utils.SanitizeName(payload.Data.User.FullName))
	cardHoldersName := utils.LocalizeToEnglish(payload.Data.Transaction.CardHoldersName)
	if fullName != cardHoldersName {
		return false, fmt.Errorf("card holder's name must match with fullname")
	}
	return true, nil
}

func (payload *requestPayload) checkCardholdersNameMatch() (bool, error) {
	txAmount, err := strconv.ParseFloat(payload.Data.Transaction.Amount, 64)
	if err != nil {
		return false, fmt.Errorf("unexpacted amount value")
	}
	if txAmount > config.MaxTransactionAmount {
		return false, fmt.Errorf("amount above maximum limits")
	}
	return true, nil
}

func (payload *requestPayload) checkUserPerm() (response bool, err error) {
	userID := &payload.Data.User.UserID
	clientID := &payload.Data.Client.Id
	*userID = strings.Trim(*userID, " ")
	*clientID = strings.Trim(*clientID, " ")
	userName := strings.Trim(payload.Data.User.Username, " ")
	fullName := utils.SanitizeName(payload.Data.User.FullName)
	if *clientID == "" || *userID == "" {
		return false, errors.New("clientID and/or userID is empty")
	}

	tx := config.MySQLDb.Raw(`SELECT privilege
		FROM cc_client_users WHERE client_id = ? AND user_id = ? LIMIT 1;`, *clientID, *userID).Scan(&response)

	//fmt.Printf("SELECT privilege AS txCount FROM cc_client_users WHERE client_id = %v AND user_id = %v LIMIT 1;\n", *clientID, *userID)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}
	if tx.RowsAffected > 0 {
		if response {
			return response, nil
		}
		return response, fmt.Errorf("user is not allowed to use the method. Please contact live support")
	}

	tx = config.MySQLDb.Exec(`INSERT INTO cc_client_users(client_id,username,user_id,fullname,privilege) VALUES (?,?,?,?,?);`,
		*clientID, userName, *userID, fullName, 0)
	if tx.Error != nil {
		return false, fmt.Errorf("\nError Details: %v", tx.Error)
	}

	return true, nil
}

func checkTcknViaApi(requestParams model.TcknCheckRequestParams, printResponse bool) (response bool, err error) {
	jsonBody, err := json.Marshal(requestParams)
	if err != nil {
		fmt.Printf("JSON marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", os.Getenv("TCKN_CHECK_URL"), bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TCKN_CHECK_SERVICE_KEY"))
	req.Header.Set("Content-Type", "application/json")

	if printResponse {
		fmt.Println("request Headers:", req.Header)
		fmt.Println("request Body:", string(jsonBody))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if printResponse {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}

	apiRes := make(map[string]string)
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		return false, err
	}
	response, err = strconv.ParseBool(apiRes["approved"])

	return response, err
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
