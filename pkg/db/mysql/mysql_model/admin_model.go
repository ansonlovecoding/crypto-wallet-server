package mysql_model

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	db "Share-Wallet/pkg/db"
	dbModel "Share-Wallet/pkg/db/mysql"
	adminStruct "Share-Wallet/pkg/struct/admin_api"
	"Share-Wallet/pkg/utils"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	AdminUserTableName          = "w_admin_user"
	AdminUserRoleTableName      = "w_admin_role"
	AdminActionsTableName       = "w_admin_actions"
	AccountInformationTableName = "w_account_information"
	FundsLogTableName           = "w_funds_log"
	CoinCurrencyValueTableName  = "w_coin_currency_value"
)

func init() {
	//init managers

	actionCount := GetActionsCount()
	if actionCount == 0 {
		var actions []string
		for _, value := range config.Config.Manager.Actions {
			var action dbModel.AdminActions
			action.ActionName = value.Name
			action.Pid = value.Pid
			action.Status = 1
			id, err := AddAction(action)
			actionId := fmt.Sprintf("%v", id)
			actions = append(actions, actionId)
			if err != nil {
				fmt.Println("AppManager insert error", err.Error(), action, "time: ")
			}
		}
		roleActions := strings.Join(actions, ",")
		AddAdminUserRole("system", "Super Admin", "", "", roleActions, 1)
	}

	for k, v := range config.Config.Manager.AppManagerUid {
		user, err := GetRegAdminUserByUserName(v)
		if err != nil {
			fmt.Println("GetUserByUserID failed ", err.Error(), v, user)
		} else {
			continue
		}
		var appMgr dbModel.AdminUser
		appMgr.UserName = v
		appMgr.Password = config.Config.Manager.Secrets[k]
		err = AdminUserRegister(appMgr)
		if err != nil {
			fmt.Println("AppManager insert error", err.Error(), appMgr, "time: ")
		}
	}
	for _, v := range config.Config.Manager.Currencies {
		coin, _ := GetCoinCurrencyValues(v)
		if coin.Coin == "" {
			var currency dbModel.CoinCurrencyValues
			currency.Coin = v
			err := AddCurrency(currency)
			if err != nil {
				fmt.Println("Currency insert error", err.Error())
			}
		}
	}
}

func GetRegAdminUserByUserName(userName string) (*dbModel.AdminUser, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminUser
	return &r, dbConn.Debug().Table(AdminUserTableName).Where("user_name = ? and delete_time = 0",
		userName).Take(&r).Error
}
func AddAdminUser(userId int64, name, password, opUser string, role string, status int, twoFactorEnabled bool, remarks string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}

	salt := utils.RandomString(10)
	newPasswordFirst := password + salt
	passwordData := []byte(newPasswordFirst)
	has := md5.Sum(passwordData)
	newPassword := fmt.Sprintf("%x", has)
	google2fSecretKey := strings.ToUpper(utils.RandomString(16))
	roleN, err := GetRegRolesByRoleName(role)
	roleId := 0
	if err == nil {
		roleId = roleN.ID
	}

	user := map[string]interface{}{
		"user_name":                      name,
		"nick_name":                      name,
		"password":                       newPassword,
		"role_id":                        roleId,
		"google_2f_secret_key":           google2fSecretKey,
		"salt":                           salt,
		"create_user":                    opUser,
		"create_time":                    time.Now().Unix(),
		"status":                         status,
		"user_two_factor_control_status": twoFactorEnabled,
		"remarks":                        remarks,
	}

	result := dbConn.Table(AdminUserTableName).Create(&user)
	return result.Error
}
func AdminUserRegister(user dbModel.AdminUser) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	user.Salt = utils.RandomString(10)
	user.Google2fSecretKey = strings.ToUpper(utils.RandomString(16))
	newPasswordFirst := user.Password + user.Salt
	passwordData := []byte(newPasswordFirst)
	has := md5.Sum(passwordData)
	user.Password = fmt.Sprintf("%x", has)
	user.CreateTime = time.Now().Unix()
	if user.UserName == constant.SuperAdmin {
		user.RoleId = 1
	}
	// user.IPRangeStart = "172.18.0.0"
	// user.IPRangeEnd = "172.18.0.200"
	// if user.AppMangerLevel == 0 {
	// 	user.AppMangerLevel = constant.AppOrdinaryUsers
	// }
	err = dbConn.Table(AdminUserTableName).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateAdminUserPassword(user dbModel.AdminUser, userName string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(AdminUserTableName).Where("user_name=? and delete_time=0", userName).Updates(&user).Error
	return err
}
func GetAdminUsers(WhereClauseMap map[string]string, page int32, pageSize int32, orderBy string) ([]dbModel.AdminUser, error) {
	var (
		userList []dbModel.AdminUser
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbQuery := dbConn.Table(AdminUserTableName)
	if err != nil {
		return userList, err
	}
	if page != -1 {
		dbQuery = dbQuery.Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}

	sortMap := map[string]string{}
	var orderByClause string
	sortMap["create_time"] = "create_time"

	if orderBy != "" {
		direction := "DESC"
		sort := strings.Split(orderBy, ":")
		if len(sort) == 2 {
			if sort[1] == "asc" {
				direction = "ASC"
			}
			col, ok := sortMap[sort[0]]
			if ok {
				orderByClause = fmt.Sprintf("%s %s ", col, direction)
			}
		}
	}
	if orderByClause != "" {
		dbQuery = dbQuery.Order(orderByClause)
	}
	if userName, ok := WhereClauseMap["user_name"]; ok {
		if userName != "" && ok {
			dbQuery = dbQuery.Debug().Where("user_name=?", userName).Where("delete_time", 0)
		}
	}
	err = dbQuery.Debug().Find(&userList, "delete_time=0 and user_name !=?", constant.SuperAdmin).Error
	if err != nil {
		return nil, err
	}
	return userList, nil
}
func GetAdminUsersCount(WhereClauseMap map[string]string) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64

	dbQuery := dbConn.Table(AdminUserTableName)

	if userName, ok := WhereClauseMap["user_name"]; ok {
		if userName != "" && ok {
			dbQuery = dbQuery.Where("user_name=?", userName)
		}
	}
	dbError := dbQuery.Where("delete_time=0").Count(&count).Error
	if dbError != nil {
		return 0, dbError
	}
	return count, nil
}
func GetAdminUserRole(WhereClauseMap map[string]string, page int32, pageSize int32, orderBy string) ([]dbModel.AdminRole, error) {
	var (
		roleList []dbModel.AdminRole
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	dbQuery := dbConn.Table(AdminUserRoleTableName)
	if err != nil {
		return roleList, err
	}
	if page != -1 {
		dbQuery = dbQuery.Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}

	sortMap := map[string]string{}
	var orderByClause string
	sortMap["create_time"] = "create_time"

	if orderBy != "" {
		direction := "DESC"
		sort := strings.Split(orderBy, ":")
		if len(sort) == 2 {
			if sort[1] == "asc" {
				direction = "ASC"
			}
			col, ok := sortMap[sort[0]]
			if ok {
				orderByClause = fmt.Sprintf("%s %s ", col, direction)
			}
		}
	}
	if orderByClause != "" {
		dbQuery = dbQuery.Order(orderByClause)
	}
	if userName, ok := WhereClauseMap["role_name"]; ok {
		if userName != "" && ok {
			dbQuery = dbQuery.Where("role_name = ?", userName)
		}
	}
	err = dbQuery.Find(&roleList, "delete_time=0").Error
	if err != nil {
		return nil, err
	}
	return roleList, nil
}
func GetAdminUserRoleCount(WhereClauseMap map[string]string) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64

	dbQuery := dbConn.Table(AdminUserRoleTableName)

	if userName, ok := WhereClauseMap["role_name"]; ok {
		if userName != "" && ok {
			dbQuery = dbQuery.Where("role_name = ?", userName)
		}
	}
	dbError := dbQuery.Debug().Where("delete_time=0").Count(&count).Error
	if dbError != nil {
		return 0, dbError
	}
	return count, nil
}
func GetRoleCount(role int) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64

	dbQuery := dbConn.Table(AdminUserTableName)

	dbQuery = dbQuery.Where("role_id = ?", role)

	dbError := dbQuery.Debug().Where("delete_time=0").Count(&count).Error
	if dbError != nil {
		return 0, dbError
	}
	return count, nil
}

// GetRegRolesByRoleName function fetches roles by rolename.
func GetRegRolesByRoleName(roleName string) (*dbModel.AdminRole, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminRole
	return &r, dbConn.Debug().Table(AdminUserRoleTableName).Where("role_name = ? and delete_time = 0",
		roleName).Take(&r).Error
}
func GetRegRolesByRoleID(roleID int) (*dbModel.AdminRole, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminRole
	return &r, dbConn.Debug().Table(AdminUserRoleTableName).Where("id = ? and delete_time = 0",
		roleID).Take(&r).Error
}

func GetRoleActionByActionID(actionID int) (*dbModel.AdminActions, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminActions
	return &r, dbConn.Table(AdminActionsTableName).Where("id = ? and delete_time = 0",
		actionID).Take(&r).Error
}
func GetRoleActionByActionName(action string) (*dbModel.AdminActions, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminActions
	return &r, dbConn.Table(AdminActionsTableName).Where("action_name = ? and delete_time = 0",
		action).Take(&r).Error
}

func GetRoleActionsByRoleName(roleName string) ([]string, error) {
	var (
		actionsList []string
	)
	// dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	// if err != nil {
	// 	return actionsList, err
	// }
	role, err := GetRegRolesByRoleName(roleName)
	if err != nil {
		return nil, err
	}

	roleActions := role.ActionIds

	actionIds := strings.Split(roleActions, ",")
	for _, raw := range actionIds {
		id, err := strconv.Atoi(raw)
		if err != nil {
			// log.Print(err)
			continue
		}
		action, err := GetRoleActionByActionID(id)
		if err != nil {
			return nil, err
		}
		actionsList = append(actionsList, action.ActionName)
	}
	// err = dbConn.Table(dbModel.AdminActions{}.TableName()).Where("delete_time = ?", 0).Find(&actionsList).Error
	return actionsList, err
}
func GetRoleActionIdsByRoleName(roleName string) ([]int64, error) {
	var (
		actionsList []int64
	)
	role, err := GetRegRolesByRoleName(roleName)
	if err != nil {
		return nil, err
	}

	roleActions := role.ActionIds

	actionIds := strings.Split(roleActions, ",")
	for _, raw := range actionIds {
		id, err := strconv.Atoi(raw)
		if err != nil {
			continue
		}
		action, err := GetRoleActionByActionID(id)
		if err != nil {
			return nil, err
		}
		actionsList = append(actionsList, int64(action.ID))
	}
	return actionsList, err
}

func GetAllActions() ([]*dbModel.AdminActions, error) {
	var (
		actionsList []*dbModel.AdminActions
		// root        *dbModel.AdminActions
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return actionsList, err
	}

	err = dbConn.Table(dbModel.AdminActions{}.TableName()).Where("delete_time = ?", 0).Find(&actionsList).Error
	// err = dbConn.Table(dbModel.AdminActions{}.TableName()).Where("pid=? and delete_time=0", 0).Take(&root).Error
	return actionsList, err
}

// AddAdminUserRole function creates new user roles.
func AddAdminUserRole(userName, roleName, description, remarks, action_ids string, status int) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	userRole := dbModel.AdminRole{
		RoleName:    roleName,
		Description: description,
		Remarks:     remarks,
		Status:      status,
		// Have to increment memberNum when new user is added.
		MemberNum:  0,
		ActionIds:  action_ids,
		CreateUser: userName,
		CreateTime: time.Now().Unix(),
	}
	result := dbConn.Table(AdminUserRoleTableName).Create(&userRole)
	return result.Error
}

func DeleteUser(user dbModel.AdminUser) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	// return dbConn.Table(AdminUserTableName).Where("user_name=?", user.UserName).Delete(&user).Error
	return dbConn.Table(AdminUserTableName).Where("user_name=? and delete_time=0", user.UserName).Updates(&user).Error

}

func UpdateUser(userID int, user dbModel.AdminUser, verify bool, login bool) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if verify {
		return dbConn.Debug().Table(AdminUserTableName).Where("id=? and delete_time=0", userID).Updates(map[string]interface{}{"user_two_factor_control_status": true, "two_factor_enabled": true, "update_time": user.UpdateTime, "update_user": user.UpdateUser}).Error
	}

	updatedUser := map[string]interface{}{
		"update_time": user.UpdateTime,
		"update_user": user.UpdateUser,
	}
	if user.RoleId != 0 {
		updatedUser["role_id"] = user.RoleId
	}
	if user.Password != "" {
		updatedUser["password"] = user.Password
	}
	if !login {
		updatedUser["user_two_factor_control_status"] = user.User2FAuthEnable
	}
	if user.Remarks != "" {
		updatedUser["remarks"] = user.Remarks
	}
	if user.Status != 0 {
		updatedUser["status"] = user.Status
	}
	if user.LoginTime != 0 {
		updatedUser["login_time"] = user.LoginTime
		updatedUser["login_ip"] = user.LoginIp
	}
	return dbConn.Debug().Table(AdminUserTableName).Where("id=? and delete_time=0", userID).Updates(updatedUser).Error
}

func UpdateAdminRole(roleID int, role dbModel.AdminRole) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	result := dbConn.Table(AdminUserRoleTableName).Where("id=?", roleID).Updates(&role)
	return result.Error
}

func DeleteRole(role dbModel.AdminRole) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}

	return dbConn.Table(AdminUserRoleTableName).Where("role_name=? and delete_time=0", role.RoleName).Updates(&role).Error

}

func CheckRoleUsage(role string) (*dbModel.AdminUser, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	tempRole, err := GetRegRolesByRoleName(role)
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminUser
	return &r, dbConn.Debug().Table(AdminUserTableName).Where("role_id=? and delete_time=0", tempRole.ID).Take(&r).Error
}

func GetRegAdminUsrByUID(userID string) (*dbModel.AdminUser, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var r dbModel.AdminUser
	return &r, dbConn.Table(AdminUserTableName).Where("user_name = ? and delete_time = 0",
		userID).Take(&r).Error
}

type Action struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Child []Action `json:"child,omitempty"`
}

func RelationMap(dataset [][]int) map[int][]int {
	relations := make(map[int][]int)
	for _, relation := range dataset {
		child, parent := relation[0], relation[1]
		relations[parent] = append(relations[parent], child)
	}
	return relations
}

func BuildActions(ids []int, relations map[int][]int) []Action {
	actions := make([]Action, len(ids))
	for i, id := range ids {
		action, err := GetRoleActionByActionID(id)
		if err != nil {
			return nil
		}
		c := Action{ID: id, Name: action.ActionName}
		if childIDs, ok := relations[id]; ok { // build child's children
			c.Child = BuildActions(childIDs, relations)
		}
		actions[i] = c
	}
	return actions
}

func MapActionParents() ([]*dbModel.AdminActions, error) {
	var (
		actionsList []*dbModel.AdminActions
	)
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, nil
	}
	err = dbConn.Table(dbModel.AdminActions{}.TableName()).Select([]string{"id", "pid"}).Where("delete_time=0").Order("pid").Find(&actionsList).Error
	return actionsList, err
}

func CheckRolePermission(roleName string, permissionName string) bool {
	role, err := GetRegRolesByRoleName(roleName)
	if err != nil {
		return false
	}
	actions, err := GetRoleActionsByRoleName(role.RoleName)
	if err != nil {
		return false
	}
	for i := 0; i < len(actions); i++ {
		// check
		if actions[i] == permissionName {
			return true
		}
	}
	return false
}

func GetAccountInformation(filters map[string]string, page int32, pageSize int32, uid string) ([]dbModel.AccountInformation, int64, *adminStruct.TotalCoins, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, nil, err
	}
	var accInfo []dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName)

	if accountSource, ok := filters["account_source"]; ok {
		if accountSource != "all" {
			dbQuery = dbQuery.Debug().Where("account_source=?", accountSource)
		}
	}
	if accountAddress, ok := filters["account_address"]; ok {
		switch filters["coins_type"] {
		case "all":
			dbQuery = dbQuery.Debug().Where("(btc_public_address=? OR erc_public_address=? OR trx_public_address=? OR trc_public_address=? OR eth_public_address=?) ",
				accountAddress, accountAddress, accountAddress, accountAddress, accountAddress)
		case "btc":
			dbQuery = dbQuery.Debug().Where("btc_public_address=?", accountAddress)
		case "eth":
			dbQuery = dbQuery.Debug().Where("eth_public_address=?", accountAddress)
		case "usdt-erc20":
			dbQuery = dbQuery.Debug().Where("erc_public_address=?", accountAddress)
		case "usdt-trc20":
			dbQuery = dbQuery.Debug().Where("trc_public_address=?", accountAddress)
		case "trx":
			dbQuery = dbQuery.Debug().Where("trx_public_address=?", accountAddress)
		default:
			dbQuery = dbQuery.Debug().Where("(btc_public_address=? OR erc_public_address=? OR trx_public_address=? OR trc_public_address=? OR eth_public_address=?) ",
				accountAddress, accountAddress, accountAddress, accountAddress, accountAddress)
		}

	}
	if merchantUid, ok := filters["merchant_uid"]; ok {
		dbQuery = dbQuery.Debug().Where("merchant_uid=?", merchantUid)
	}
	if uid != "" {
		dbQuery = dbQuery.Debug().Where("uuid=?", uid)
	}

	from, fok := filters["from"]
	to, tok := filters["to"]
	if orderBy, ok := filters["order_by"]; ok {
		dir := "desc"
		if sort, ok := filters["sort"]; ok {
			dir = sort
		}
		orderByClause := fmt.Sprintf("%s %s ", orderBy, dir)
		dbQuery = dbQuery.Debug().Order(orderByClause)
		if tok && fok {
			if orderBy == "creation_time" {
				dbQuery = dbQuery.Debug().Where("creation_time between ? and ?", from, to)

			} else if orderBy == "last_login_time" {
				dbQuery = dbQuery.Debug().Where("last_login_time between ? and ?", from, to)

			}
		}
	}
	var totals *adminStruct.TotalCoins
	dbQuery.Debug().Select(`
	SUM(btc_balance) as btc,
	SUM(eth_balance) as eth,
	SUM(erc_balance) as erc,
	SUM(trx_balance) as trx,
	SUM(trc_balance) as trc
	`).Scan(&totals)
	var count int64
	dbQuery.Debug().Count(&count)
	if page != -1 {
		dbQuery = dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}
	dbQuery = dbQuery.Debug().Find(&accInfo)
	err = dbQuery.Error
	return accInfo, count, totals, err
}

func GetFundsLog(filters map[string]string, page int32, pageSize int32, uid string) ([]dbModel.FundsLog, int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, err
	}
	var fundsLog []dbModel.FundsLog
	dbQuery := dbConn.Table(FundsLogTableName)
	from, fok := filters["from"]
	to, tok := filters["to"]
	if tok && fok {
		dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", from, to)
	}
	if transactionType, ok := filters["transaction_type"]; ok {
		if transactionType != "all" {
			dbQuery = dbQuery.Debug().Where("transaction_type=?", transactionType)
		}
	}
	if userAddress, ok := filters["user_address"]; ok {
		dbQuery = dbQuery.Debug().Where("user_address=?", userAddress)
	}
	if oppositeAddress, ok := filters["opposite_address"]; ok {
		dbQuery = dbQuery.Debug().Where("opposite_address=?", oppositeAddress)
	}

	// use coins type to decide which table to retreive the data from
	if coinsType, ok := filters["coins_type"]; ok {
		if coinsType != "all" {
			dbQuery = dbQuery.Debug().Where("coin_type=?", coinsType)
		}
	}
	if state, ok := filters["state"]; ok {
		if state != "all" {
			switch state {
			case "fail":
				dbQuery = dbQuery.Debug().Where("state=?", 0)
			case "success":
				dbQuery = dbQuery.Debug().Where("state=?", 1)
			}
		} else {
			dbQuery = dbQuery.Debug().Where("state != ?", 2)
		}
	}
	if txid, ok := filters["txid"]; ok {
		dbQuery = dbQuery.Debug().Where("txid=?", txid)
	}
	if merchantUid, ok := filters["merchant_uid"]; ok {
		dbQuery = dbQuery.Debug().Where("merchant_uid=?", merchantUid)
	}
	if uid != "" {
		dbQuery = dbQuery.Debug().Where("uid=?", uid)
	}
	dbQuery = dbQuery.Debug().Where("merchant_uid != ?", "")
	var count int64
	dbQuery.Count(&count)
	if page > 0 && pageSize > 0 {
		dbQuery = dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}
	dbQuery = dbQuery.Debug().Order("confirmation_time DESC,id DESC")
	result := dbQuery.Debug().Find(&fundsLog)
	err = result.Error
	return fundsLog, count, err
}

func GetFundLogByCoinAndHash(coinType int, txHash string) (*dbModel.FundsLog, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var fundLog dbModel.FundsLog
	if err = dbConn.Table(FundsLogTableName).Debug().Where("coin_type = ? and txid = ?", utils.GetCoinName(uint8(coinType)), txHash).Find(&fundLog).Error; err != nil {
		return nil, err
	}
	return &fundLog, err
}

func GetCoinCurrencyValues(coin string) (*dbModel.CoinCurrencyValues, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var coinCurrencyValue *dbModel.CoinCurrencyValues
	err = dbConn.Table(dbModel.CoinCurrencyValues{}.TableName()).Debug().Where("coin=?", coin).Find(&coinCurrencyValue).Error
	return coinCurrencyValue, err
}

func GetReceiveDetails(filters map[string]string, page int32, pageSize int32, uid string) ([]dbModel.FundsLog, int64, *adminStruct.TotalCoins, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, nil, err
	}
	var receiveDetails []dbModel.FundsLog
	dbQuery := dbConn.Table(FundsLogTableName)
	from, fok := filters["from"]
	to, tok := filters["to"]
	if tok && fok {
		dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", from, to)
	}
	// else {
	// 	lastWeek := time.Now().AddDate(0, 0, -7).Unix()
	// 	now := time.Now().Unix()
	// 	dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", lastWeek, now)
	// }
	if receivingAddress, ok := filters["receiving_address"]; ok {
		dbQuery = dbQuery.Debug().Where("user_address=?", receivingAddress)
	}
	if depositAddress, ok := filters["deposit_address"]; ok {
		dbQuery = dbQuery.Debug().Where("opposite_address=?", depositAddress)
	}

	// use coins type to decide which table to retreive the data from
	if coinsType, ok := filters["coins_type"]; ok {
		if filters["coins_type"] != "all" {
			dbQuery = dbQuery.Debug().Where("coin_type=?", coinsType)
		}
	}
	if merchantUid, ok := filters["merchant_uid"]; ok {
		dbQuery = dbQuery.Debug().Where("merchant_uid=?", merchantUid)
	}
	if txid, ok := filters["txid"]; ok {
		dbQuery = dbQuery.Debug().Where("txid=?", txid)
	}
	if uid != "" {
		dbQuery = dbQuery.Debug().Where("uid=?", uid)
	}
	dbQuery = dbQuery.Debug().Where("merchant_uid != ?", "")
	var totals *adminStruct.TotalCoins
	dbQuery.Debug().Select(`
	SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'btc' AND state = '1' THEN amount_of_coins ELSE 0 END) as btc,
	SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'eth' AND state = '1' THEN amount_of_coins ELSE 0 END) as eth,
	SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'usdt-erc20' AND state = '1' THEN amount_of_coins ELSE 0 END) as erc,
	SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'trx' AND state = '1' THEN amount_of_coins  ELSE 0 END) as trx,
	SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'usdt-trc20' AND state = '1' THEN amount_of_coins ELSE 0 END) as trc
	`).Scan(&totals)
	dbQuery = dbQuery.Debug().Where("transaction_type = ?", "received").Order("confirmation_time DESC")
	var count int64
	dbQuery.Debug().Count(&count)
	if page != -1 {
		dbQuery = dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}
	dbQuery = dbQuery.Debug().Find(&receiveDetails)
	err = dbQuery.Error
	return receiveDetails, count, totals, err
}
func GetTransferDetails(filters map[string]string, page int32, pageSize int32, uid string, state int32) ([]dbModel.FundsLog, int64, *adminStruct.TotalCoins, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, nil, err
	}
	var transferDetails []dbModel.FundsLog
	dbQuery := dbConn.Table(FundsLogTableName)
	dbQuery = dbQuery.Debug().Where("transaction_type=?", "transfer")
	from, fok := filters["from"]
	to, tok := filters["to"]
	if filters["order_by"] != "confirmation_time" || filters["order_by"] == "" {
		if tok && fok {
			dbQuery = dbQuery.Debug().Where("creation_time between ? and ?", from, to).Order("creation_time DESC")
		}
		// else {
		// 	lastWeek := time.Now().AddDate(0, 0, -7).Unix()
		// 	now := time.Now().Unix()
		// 	dbQuery = dbQuery.Debug().Where("creation_time between ? and ?", lastWeek, now).Order("creation_time DESC")
		// }
	} else {
		if tok && fok {
			dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", from, to).Order("confirmation_time DESC")
		}
		// else {
		// 	lastWeek := time.Now().AddDate(0, 0, -7).Unix()
		// 	now := time.Now().Unix()
		// 	dbQuery = dbQuery.Debug().Where("confirmation_time between ? and ?", lastWeek, now).Order("confirmation_time DESC")
		// }
	}
	if transferAddress, ok := filters["transfer_address"]; ok {
		dbQuery = dbQuery.Debug().Where("user_address=?", transferAddress)
	}
	if receivingAddress, ok := filters["receiving_address"]; ok {
		dbQuery = dbQuery.Debug().Where("opposite_address=?", receivingAddress)
	}
	if uid != "" {
		dbQuery = dbQuery.Debug().Where("uid=?", uid)
	}
	if merchantUid, ok := filters["merchant_uid"]; ok {
		dbQuery = dbQuery.Debug().Where("merchant_uid=?", merchantUid)
	}

	// use coins type to decide which table to retreive the data from
	if coinsType, ok := filters["coins_type"]; ok {
		if filters["coins_type"] != "all" {
			dbQuery = dbQuery.Debug().Where("coin_type=?", coinsType)
		}
	}
	if state != -1 {
		dbQuery = dbQuery.Debug().Where("state=?", state)
	}
	if txid, ok := filters["txid"]; ok {
		dbQuery = dbQuery.Debug().Where("txid=?", txid)
	}
	var totals *adminStruct.TotalCoins
	dbQuery.Debug().Select(`
	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'btc' AND state = '1' THEN amount_of_coins ELSE 0 END) as btc,
	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'eth' AND state = '1' THEN amount_of_coins ELSE 0 END) as eth,
	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'usdt-erc20' AND state = '1' THEN amount_of_coins ELSE 0 END) as erc,
	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'trx' AND state = '1' THEN amount_of_coins ELSE 0 END) as trx,
	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'usdt-trc20' AND state = '1' THEN amount_of_coins ELSE 0 END) as trc,
	SUM(CASE WHEN coin_type = 'btc' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'btc_fee',
	SUM(CASE WHEN coin_type = 'eth' and transaction_type='transfer' THEN gas_price*POWER(0.1,18)*gas_used  ELSE 0 END) as 'eth_fee',
	SUM(CASE WHEN coin_type = 'usdt-erc20' and transaction_type='transfer' THEN gas_price*POWER(0.1,18)*gas_used  ELSE 0 END) as 'erc_fee',
	SUM(CASE WHEN coin_type = 'trx' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'trx_fee',
	SUM(CASE WHEN coin_type = 'usdt-trc20' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'trc_fee'
	`).Scan(&totals)
	var count int64
	dbQuery.Debug().Count(&count)
	if page != -1 {
		dbQuery = dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1)))
	}
	dbQuery = dbQuery.Debug().Find(&transferDetails)
	err = dbQuery.Error
	return transferDetails, count, totals, err
}
func ResetGooogleKey(userName string, updateUser string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	google2fSecretKey := strings.ToUpper(utils.RandomString(16))
	result := dbConn.Table(AdminUserTableName).Debug().Table(AdminUserTableName).Where("user_name=? and delete_time=0", userName).Updates(map[string]interface{}{
		"google_2f_secret_key":           google2fSecretKey,
		"two_factor_enabled":             false,
		"user_two_factor_control_status": true,
		"update_time":                    time.Now().Unix(),
		"update_user":                    updateUser})
	return result.Error
}

type FundsTotal struct {
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

func GetTotalFunds(filter string) ([]FundsTotal, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var fundsTotal []FundsTotal
	dbQuery := dbConn.Table(FundsLogTableName)
	dbQuery = dbQuery.Select("coin_type as name, sum(amount_of_coins) as total").Where("transaction_type = ?", filter).Group("coin_type").Scan(&fundsTotal)
	err = dbQuery.Error
	return fundsTotal, err
}

func GetCurrencies(page int32, pageSize int32) ([]dbModel.CoinCurrencyValues, int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, err
	}
	var currencies []dbModel.CoinCurrencyValues
	dbQuery := dbConn.Table(CoinCurrencyValueTableName).Debug().Scan(&currencies)
	count := dbQuery.RowsAffected
	if page != -1 {
		dbQuery.Debug().Limit(int(pageSize)).Offset(int(pageSize * (page - 1))).Scan(&currencies)
	}
	err = dbQuery.Error
	return currencies, count, err
}

func UpdateCurrency(currencyID int32, currency dbModel.CoinCurrencyValues) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	updatedCurrency := map[string]interface{}{
		"update_user": currency.UpdateUser,
		"update_time": currency.UpdateTime,
		"state":       currency.State,
	}
	err = dbConn.Table(CoinCurrencyValueTableName).Debug().Where("id=?", currencyID).Updates(updatedCurrency).Error
	return err
}

func AddAction(action dbModel.AdminActions) (int, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	err = dbConn.Table(AdminActionsTableName).Debug().Create(&action).Error
	if err != nil {
		return 0, err
	}
	return action.ID, err
}
func GetActionsCount() int64 {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0
	}
	var count int64
	dbConn.Table(AdminActionsTableName).Count(&count)
	if err != nil {
		return 0
	}
	return count
}
func CreateAccountInformation(account dbModel.AccountInformation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(AccountInformationTableName).Debug().Create(&account).Error
	if err != nil {
		return err
	}
	return err
}

func CreateFundLog(fundLog dbModel.FundsLog) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(FundsLogTableName).Debug().Create(&fundLog).Error
	if err != nil {
		return err
	}
	return err
}

func UpdateFundLog(fundLog dbModel.FundsLog) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	updatedLog := map[string]interface{}{
		"state":             fundLog.State,
		"confirmation_time": fundLog.ConfirmationTime,
	}
	err = dbConn.Table(FundsLogTableName).Debug().Where("txid=?", fundLog.Txid).Updates(updatedLog).Error
	if err != nil {
		return err
	}
	return err
}
func UpdateAccountInformation(account *dbModel.AccountInformation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	updatedAccount := map[string]interface{}{
		"btc_balance": account.BtcBalance,
		"eth_balance": account.EthBalance,
		"erc_balance": account.ErcBalance,
		"trx_balance": account.TrxBalance,
		"trc_balance": account.TrcBalance,
	}

	dbQuery := dbConn.Table(AccountInformationTableName)
	err = dbQuery.Debug().Where("uuid=? and merchant_uid=?", account.UUID, account.MerchantUid).Updates(&updatedAccount).Error
	if err != nil {
		return err
	}
	return err
}
func UpdateAccountInfoWithAccount(account *dbModel.AccountInformation) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(AccountInformationTableName).Debug().Updates(account).Error
	if err != nil {
		return err
	}
	return err
}

func GetTransactionByTxid(txid string) (*dbModel.FundsLog, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var fundLog dbModel.FundsLog
	err = dbConn.Table(FundsLogTableName).Where("txid", txid).Take(&fundLog).Error
	return &fundLog, err
}

func GetTotalFee() ([]FundsTotal, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var fundsTotal []FundsTotal
	dbQuery := dbConn.Table("w_funds_log").Debug().Select("coin_type as name, sum(network_fee) as total").Group("coin_type").Scan(&fundsTotal)
	err = dbQuery.Error
	return fundsTotal, err

}

func GetOperationalReport(fromDate, toDate string, page int32, pageSize int32) ([]adminStruct.Operation, int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, 0, err
	}
	var operation []adminStruct.Operation
	limit := int64(10)
	offset := 0

	if page != -1 {
		limit = int64(pageSize)
		offset = int(pageSize * (page - 1))
	}
	if fromDate == "" && toDate == "" {
		if page == -1 {
			dbConn.Table(FundsLogTableName).Debug().Count(&limit)
		}
		dbConn.Table(FundsLogTableName).Debug().Select("min(creation_time)").Scan(&fromDate)
		dbConn.Debug().Raw(`SELECT UNIX_TIMESTAMP() AS to_date`).Scan(&toDate)
	}
	dbQuery := dbConn.Debug().Raw(`SELECT *
	FROM
	  (
		SELECT a.Date AS date
		FROM (
			   SELECT FROM_UNIXTIME(?,'%Y-%m-%d') - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date
			   FROM (SELECT 0 AS a
					 UNION ALL SELECT 1
					 UNION ALL SELECT 2
					 UNION ALL SELECT 3
					 UNION ALL SELECT 4
					 UNION ALL SELECT 5
					 UNION ALL SELECT 6
					 UNION ALL SELECT 7
					 UNION ALL SELECT 8
					 UNION ALL SELECT 9) AS a
				 CROSS JOIN (SELECT 0 AS a
							 UNION ALL SELECT 1
							 UNION ALL SELECT 2
							 UNION ALL SELECT 3
							 UNION ALL SELECT 4
							 UNION ALL SELECT 5
							 UNION ALL SELECT 6
							 UNION ALL SELECT 7
							 UNION ALL SELECT 8
							 UNION ALL SELECT 9) AS b
				 CROSS JOIN (SELECT 0 AS a
							 UNION ALL SELECT 1
							 UNION ALL SELECT 2
							 UNION ALL SELECT 3
							 UNION ALL SELECT 4
							 UNION ALL SELECT 5
							 UNION ALL SELECT 6
							 UNION ALL SELECT 7
							 UNION ALL SELECT 8
							 UNION ALL SELECT 9) AS c
			 ) a
		WHERE a.Date BETWEEN FROM_UNIXTIME(?,'%Y-%m-%d') AND FROM_UNIXTIME(?,'%Y-%m-%d')
	  ) dates
	  LEFT JOIN
	  (
		SELECT FROM_UNIXTIME(confirmation_time,'%Y-%m-%d') AS week_day,	SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'btc' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'btc_transfer',
		SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'eth' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'eth_transfer',
		SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'usdt-erc20' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'erc_transfer',
		SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'usdt-trc20' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'trc_transfer',
		SUM(CASE WHEN transaction_type = 'transfer' AND coin_type = 'trx' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'trx_transfer',
		SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'btc' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'btc_received',
		SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'eth' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'eth_received',
		SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'usdt-erc20' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'erc_received',
		SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'usdt-trc20' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'trc_received',
		SUM(CASE WHEN transaction_type = 'received' AND coin_type = 'trx' AND state = 1 THEN amount_of_coins ELSE 0 END) as 'trx_received',
		SUM(CASE WHEN coin_type = 'btc' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'btc_fee',
		SUM(CASE WHEN coin_type = 'eth' and transaction_type='transfer' THEN gas_price*POWER(0.1,18)*gas_used  ELSE 0 END) as 'eth_fee',
		SUM(CASE WHEN coin_type = 'usdt-erc20' and transaction_type='transfer' THEN gas_price*POWER(0.1,18)*gas_used  ELSE 0 END) as 'erc_fee',
		SUM(CASE WHEN coin_type = 'trx' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'trx_fee',
		SUM(CASE WHEN coin_type = 'usdt-trc20' and transaction_type='transfer' THEN network_fee ELSE 0 END) as 'trc_fee'
			FROM
		  w_funds_log
				WHERE merchant_uid !=? AND confirmation_time between ? AND ?
				GROUP BY week_day
	  ) data
		ON DATE_FORMAT(dates.date, '%Y-%m-%d') = week_day ORDER BY date DESC
	LIMIT ? OFFSET ?`, toDate, fromDate, toDate, "", fromDate, toDate, limit, offset).Scan(&operation)
	var count int64
	dbConn.Debug().Raw("SELECT DATEDIFF(FROM_UNIXTIME(?,'%Y-%m-%d'), FROM_UNIXTIME(?,'%Y-%m-%d')) as count", toDate, fromDate).Scan(&count)

	err = dbQuery.Error
	count += 1
	return operation, count, err
}

func GetCreatedUsersCountPerDay(day int64) (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	timeDay := time.Unix(day, 0).Day()
	timeMonth := time.Unix(day, 0).Month()
	timeYear := time.Unix(day, 0).Year()
	var count int64
	result := dbConn.Debug().Raw(`SELECT COUNT(uuid) FROM w_account_information WHERE FROM_UNIXTIME(creation_time,'%d') = ? AND FROM_UNIXTIME(creation_time,'%m') = ? AND FROM_UNIXTIME(creation_time,'%Y') = ?
		AND uuid NOT IN (SELECT DISTINCT uuid FROM w_account_information WHERE FROM_UNIXTIME(creation_time,'%Y-%m-%d') < FROM_UNIXTIME(?,'%Y-%m-%d') )`,
		timeDay, timeMonth, timeYear, day).Scan(&count)
	err = result.Error
	return count, err
}

func GetAllWallets() ([]dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var wallets []dbModel.AccountInformation
	dbQuery := dbConn.Debug().Table(AccountInformationTableName).Find(&wallets)
	err = dbQuery.Error
	return wallets, err
}

func GetAccountInformationByPublicAddress(accountAddress string, coinType uint32) (*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account *dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName)
	switch coinType {
	case constant.BTCCoin:
		dbQuery = dbQuery.Where("btc_public_address=?", accountAddress)
	case constant.ETHCoin:
		dbQuery = dbQuery.Where("eth_public_address=?", accountAddress)
	case constant.USDTERC20:
		dbQuery = dbQuery.Where("erc_public_address=?", accountAddress)
	case constant.USDTTRC20:
		dbQuery = dbQuery.Where("trc_public_address=?", accountAddress)
	case constant.TRX:
		dbQuery = dbQuery.Where("trx_public_address=?", accountAddress)
	}
	dbQuery = dbQuery.Find(&account)
	err = dbQuery.Error
	return account, err
}
func GetAccountInformationByMerchantUid(merchantUid string) (*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account *dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName).Debug().Where("merchant_uid", merchantUid).Take(&account)
	err = dbQuery.Error
	return account, err
}
func GetAccountInformationByMerchantUidAndUid(merchantUid, uid string) (*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account *dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName).Debug().Where("merchant_uid = ? and uuid = ?", merchantUid, uid).Take(&account)
	err = dbQuery.Error
	return account, err
}
func GetAccountInformationByUUID(uuid string) (*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account *dbModel.AccountInformation
	err = dbConn.Table(AccountInformationTableName).Debug().Where("uuid = ?", uuid).Find(&account).Error
	return account, err
}
func GetAccountInformationListByUid(uid string) ([]*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var accounts []*dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName).Debug().Where("uuid = ?", uid).Find(&accounts)
	err = dbQuery.Error
	return accounts, err
}
func UpdateAccountLoginInfo(account dbModel.AccountInformation, merchantUid, uid string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(AccountInformationTableName).Debug().Where("merchant_uid=? and uuid = ?", merchantUid, uid).Updates(&account).Error
	if err != nil {
		return err
	}
	return err
}
func GetAccountByAddress(accountAddress string) (*dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account *dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName)
	dbQuery = dbQuery.Debug().Where("btc_public_address=?", accountAddress).
		Or("eth_public_address=?", accountAddress).
		Or("erc_public_address=?", accountAddress).
		Or("trc_public_address=?", accountAddress).
		Or("trx_public_address=?", accountAddress)
	dbQuery = dbQuery.Take(&account)
	err = dbQuery.Error
	return account, err
}

func GetCoinStatuses() ([]dbModel.CoinCurrencyValues, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var coins []dbModel.CoinCurrencyValues
	dbQuery := dbConn.Table(CoinCurrencyValueTableName).Debug().Find(&coins)
	err = dbQuery.Error
	return coins, err
}

func UpdateCoinRates(coin *dbModel.CoinCurrencyValues) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(CoinCurrencyValueTableName).Debug().Where("coin=?", coin.Coin).Updates(&coin).Error
	if err != nil {
		return err
	}
	return err
}
func BulkUpdateCoinRates(coins map[string]*dbModel.CoinCurrencyValues, coinTypes []string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	var updatedCoins []dbModel.CoinCurrencyValues
	dbConn.Debug().Raw(`
	UPDATE w_coin_currency_value 
	SET usd=(CASE WHEN coin = 'BTC' THEN ? WHEN coin = 'ETH' THEN ? WHEN coin = 'USDT-ERC20' THEN ? WHEN coin = 'TRX' THEN ? WHEN coin = 'USDT-TRC20' THEN ? END),
	euro=(CASE WHEN coin = 'BTC' THEN ? WHEN coin = 'ETH' THEN ? WHEN coin = 'USDT-ERC20' THEN ? WHEN coin = 'TRX' THEN ? WHEN coin = 'USDT-TRC20' THEN ? END),
	yuan=(CASE WHEN coin = 'BTC' THEN ? WHEN coin = 'ETH' THEN ? WHEN coin = 'USDT-ERC20' THEN ? WHEN coin = 'TRX' THEN ? WHEN coin = 'USDT-TRC20' THEN ? END)
	WHERE coin IN (?)
	`,
		coins["btc"].Usd, coins["eth"].Usd, coins["erc"].Usd, coins["trx"].Usd, coins["trc"].Usd,
		coins["btc"].Euro, coins["eth"].Euro, coins["erc"].Euro, coins["trx"].Euro, coins["trc"].Euro,
		coins["btc"].Yuan, coins["eth"].Yuan, coins["erc"].Yuan, coins["trx"].Yuan, coins["trc"].Yuan, coinTypes,
	).Scan(&updatedCoins)
	if err != nil {
		return err
	}
	return err
}

func GetWalletsCount() (int64, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return 0, err
	}
	var count int64
	dbQuery := dbConn.Table(AccountInformationTableName).Debug().Distinct("uuid").Count(&count)
	return count, dbQuery.Error
}

func GetAccountInformationByPublicAddressTemp(accountAddress string) ([]dbModel.AccountInformation, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var account []dbModel.AccountInformation
	dbQuery := dbConn.Table(AccountInformationTableName)
	dbQuery = dbQuery.Debug().Where("(btc_public_address=? OR erc_public_address=? OR trx_public_address=? OR trc_public_address=? OR eth_public_address=?) ",
		accountAddress, accountAddress, accountAddress, accountAddress, accountAddress)
	dbQuery = dbQuery.Find(&account)
	err = dbQuery.Error
	return account, err
}
func AddCurrency(currency dbModel.CoinCurrencyValues) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	err = dbConn.Table(CoinCurrencyValueTableName).Debug().Create(&currency).Error
	if err != nil {
		return err
	}
	return err
}
