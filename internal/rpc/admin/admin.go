package admin

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/common/token_verify"
	db2 "Share-Wallet/pkg/db"
	db "Share-Wallet/pkg/db/mysql"
	walletdb "Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/admin"
	"Share-Wallet/pkg/utils"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	gotp "github.com/diebietse/gotp/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type adminRPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewAdminRPCServer(port int) *adminRPCServer {
	return &adminRPCServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.AdminRPC,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *adminRPCServer) Run() {
	log.NewInfo("0", "AdminCMS rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	// listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	// grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	// Service registers with etcd
	admin.RegisterAdminServer(srv, s)
	rpcRegisterIP := ""
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *adminRPCServer) TestAdminRPC(_ context.Context, req *admin.CommonReq) (*admin.CommonResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestAdminRPC!!", req.String())
	resp := &admin.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Test Success!",
	}
	return resp, nil
}

// admin login api
func (s *adminRPCServer) AdminLogin(_ context.Context, req *admin.AdminLoginReq) (*admin.AdminLoginResp, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.AdminLoginResp{}
	user, err := walletdb.GetRegAdminUserByUserName(req.AdminID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName failed!", "adminID: ", req.AdminID, err.Error())
		return resp, http.WrapError(constant.ErrUserNotExist)
	}
	if user.Status == 2 {
		return resp, http.WrapError(constant.ErrUserBanned)
	}
	password := req.Secret
	if !req.SecretHashd {
		newPasswordFirst := req.Secret + user.Salt
		passwordData := []byte(newPasswordFirst)
		has := md5.Sum(passwordData)
		password = fmt.Sprintf("%x", has)
	}

	if password == user.Password {
		if !user.User2FAuthEnable || !config.Config.AdminUser2FAuthEnable {
			req.GAuthTypeToken = false
		}
		token, expTime, err := token_verify.CreateToken(req.AdminID, constant.AdminPlatformID, req.GAuthTypeToken)
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", req.AdminID, err.Error())
			return resp, http.WrapError(constant.ErrTokenUnknown)
		}
		resp.Token = token
		updatedUser := db.AdminUser{
			UpdateTime: time.Now().Unix(),
			UpdateUser: user.UserName,
			LoginTime:  time.Now().Unix(),
			LoginIp:    req.LoginIp,
		}
		err = walletdb.UpdateUser(user.ID, updatedUser, false, true)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAdmin failed", err.Error())
			return nil, http.WrapError(constant.ErrDB)
		}
		// resp.User.UserName = user.UserName
		// 	if user.RoleId != 0 {
		// 		role, err := walletdb.GetRegRolesByRoleID(user.RoleId)
		// 		if err != nil {
		// 			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleID failed!", "adminID: ", req.AdminID, err.Error())
		// 			return resp, http.WrapError(constant.ErrDB)
		// 		}
		// 		resp.User.Role = role.RoleName
		// 	}
	}
	if resp.Token == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "failed")
		return resp, http.WrapError(constant.ErrTokenMalformed)
	}

	resp.GAuthEnabled = user.TwoFactorEnabled
	resp.GAuthSetupRequired = !user.TwoFactorEnabled && config.Config.AdminUser2FAuthEnable
	if !user.User2FAuthEnable {
		resp.GAuthSetupRequired = false
		resp.GAuthEnabled = false
	}
	if resp.GAuthSetupRequired {
		resp.GAuthSetupProvUri = genrateTOTPProvisUriQR(*user)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())

	return resp, nil
}

// genrateTOTPProvisUriQR genrating QR code URI for user
// TOTP -> Time-Based One-Time Password
func genrateTOTPProvisUriQR(user db.AdminUser) string {
	// totp := gotp.NewDefaultTOTP(user.Google2fSecretKey)
	secret, err := gotp.DecodeBase32(user.Google2fSecretKey)
	if err != nil {
		log.NewError(utils.GetSelfFuncName(), "Genrating TOTP Decoding Failed", err.Error())
	}
	totp, err := gotp.NewTOTP(secret)
	if err != nil {
		log.NewError(utils.GetSelfFuncName(), "Genrating TOTP NewTOTP Failed", err.Error())
	}
	provisingUri, err := totp.ProvisioningURI(user.UserName, config.Config.TotpIssuerName)
	if err != nil {
		log.NewError(utils.GetSelfFuncName(), "Genrating TOTP Provision URL Failed", err.Error())
	}
	return provisingUri
}

// genrateTOTPForNow genrate TOTP code for varification
func genrateTOTPForNow(user db.AdminUser) string {
	// totp := gotp.NewDefaultTOTP(user.Google2fSecretKey)
	secret, _ := gotp.DecodeBase32(user.Google2fSecretKey)
	totp, _ := gotp.NewTOTP(secret)
	totpCode, err := totp.Now()
	if err != nil {
		log.NewError(utils.GetSelfFuncName(), "Genrating TOTP Failed", err.Error())
	}
	return totpCode
}

func (s *adminRPCServer) AddAdminUser(ctx context.Context, req *admin.AddAdminUserReq) (*admin.AddAdminUserResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.AddAdminUserResp{}

	user, err := walletdb.GetRegAdminUserByUserName(req.Name)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}

	if user.ID != 0 {
		return resp, http.WrapError(constant.ErrUserIDAlreadyExsist)
	}

	err = walletdb.AddAdminUser(req.UserID, req.Name, req.Password, req.OpUserId, req.Role, int(req.Status), req.GAuthEnabled, req.Remarks)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddAdminUser failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

// admin login api
func (s *adminRPCServer) ChangeAdminUserPassword(ctx context.Context, req *admin.ChangeAdminUserPasswordReq) (*admin.ChangeAdminUserPasswordResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.ChangeAdminUserPasswordResp{}

	user, err := walletdb.GetRegAdminUserByUserName(req.UserName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName failed", "adminID: ", req.UserName, err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	// twoFacAuthFalg := user.TwoFactorEnabled && config.Config.AdminUser2FAuthEnable
	// totp := genrateTOTPForNow(*user)
	// if req.TOTP == totp || !twoFacAuthFalg {
	password := req.Secret
	newPasswordFirst := password + user.Salt
	passwordData := []byte(newPasswordFirst)
	has := md5.Sum(passwordData)
	password = fmt.Sprintf("%x", has)

	if password == user.Password {
		passwordNew := req.NewSecret
		newPasswordFirst := passwordNew + user.Salt
		passwordData := []byte(newPasswordFirst)
		has := md5.Sum(passwordData)
		passwordNew = fmt.Sprintf("%x", has)

		updateData := db.AdminUser{
			Password:   passwordNew,
			UpdateTime: time.Now().Unix(),
			UpdateUser: user.UserName,
		}
		if err := walletdb.UpdateAdminUserPassword(updateData, user.UserName); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "Update user password failed!", err.Error())
			return resp, http.WrapError(constant.ErrDB)
		}

		if err := token_verify.DeleteAdminTokenOnLogout(user.UserName, false); err != nil {
			errMsg := req.OperationID + " DeleteToken failed " + err.Error() + user.UserName + "Admin"
			log.NewError(req.OperationID, errMsg)
			return resp, http.WrapError(constant.ErrDB)
		}

		token, expTime, err := token_verify.CreateToken(req.UserName, constant.AdminPlatformID, false)
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", req.UserName, err.Error())
			return resp, http.WrapError(constant.ErrTokenUnknown)
		}
		resp.Token = token
		resp.PasswordUpdated = true
	}
	if resp.Token == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "failed")
		return resp, http.WrapError(constant.ErrWrongPassword)
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
	// }
	// return resp, errors.New("TOTP is not correct")
}

// GetAdminUserList RPC returns admin users.
func (s *adminRPCServer) GetAdminUserList(ctx context.Context, req *admin.GetAdminUserListReq) (*admin.GetAdminUserListResp, error) {

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.GetAdminUserListResp{}

	whereConditionMap := map[string]string{
		"user_name": req.Name,
	}
	userCount, err := walletdb.GetAdminUsersCount(whereConditionMap)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUsersCount failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if userCount == 0 {
		return resp, nil
	}
	userList, err := walletdb.GetAdminUsers(whereConditionMap, req.Pagination.Page, req.Pagination.PageSize, req.OrderBy)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUsers failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	for _, v := range userList {
		role, err := walletdb.GetRegRolesByRoleID(v.RoleId)
		var name string
		if err != nil {
			name = ""
		} else {
			name = role.RoleName
		}
		user := &admin.AdminUser{
			Id:               int32(v.ID),
			UserName:         v.UserName,
			Role:             name,
			LastLoginIP:      v.LoginIp,
			LastLoginTime:    v.LoginTime,
			Status:           int32(v.Status),
			TwoFactorEnabled: v.User2FAuthEnable,
			Remarks:          v.Remarks,
		}
		resp.User = append(resp.User, user)
	}
	resp.TotalUsers = userCount
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

// GetAdminUserRole RPC returns admin users roles.
func (s *adminRPCServer) GetAdminUserRole(ctx context.Context, req *admin.GetAdminUserRoleReq) (*admin.GetAdminUserRoleResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.GetAdminUserRoleResp{}

	whereConditionMap := map[string]string{
		"role_name": req.Name,
	}
	roleCount, err := walletdb.GetAdminUserRoleCount(whereConditionMap)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUserRoleCount failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if roleCount == 0 {
		return resp, nil
	}
	userRoleList, err := walletdb.GetAdminUserRole(whereConditionMap, req.Pagination.Page, req.Pagination.PageSize, req.OrderBy)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUsers failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	for _, v := range userRoleList {
		totalAdmins, err := walletdb.GetRoleCount(v.ID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUserRoleCount failed", err.Error())
			return resp, http.WrapError(constant.ErrDB)
		}
		userRole := &admin.AdminUserRole{
			Id:              int32(v.ID),
			RoleName:        v.RoleName,
			RoleDescription: v.Description,
			RoleNumber:      int32(totalAdmins),
			CreateTime:      v.CreateTime,
			CreateUser:      v.CreateUser,
			UpdateUser:      v.UpdateUser,
			UpdateTime:      v.UpdateTime,
			Status:          int32(v.Status),
			Remarks:         v.Remarks,
		}
		resp.Role = append(resp.Role, userRole)
	}
	resp.TotalUserRole = roleCount
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}
func (s *adminRPCServer) AddAdminUserRole(ctx context.Context, req *admin.AddAdminUserRoleReq) (*admin.AddAdminRoleResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.AddAdminRoleResp{}

	role, err := walletdb.GetRegRolesByRoleName(req.RoleName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleName", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}

	if role.ID != 0 {
		return resp, http.WrapError(constant.ErrRoleNameAlreadyExist)
	}

	err = walletdb.AddAdminUserRole(req.UserName, req.RoleName, req.RoleDescription, req.Remarks, req.ActionIDs, int(req.Status))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddAdminUserRole failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

func (s *adminRPCServer) DeleteAdminUser(_ context.Context, req *admin.DeleteAdminReq) (*admin.CommonResp, error) {

	user, err := walletdb.GetRegAdminUserByUserName(req.UserName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	if user.ID == 0 {
		return nil, http.WrapError(constant.ErrAddInviteCodeUserNotExist) //Edit to user ID doesnt exist
	}
	deletedUser := db.AdminUser{
		DeleteTime: time.Now().Unix(),
		DeleteUser: req.DeleteUser,
		UserName:   req.UserName,
	}
	err = walletdb.DeleteUser(deletedUser)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteUser failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "DeleteAdminRPC...", req.String())
	resp := &admin.CommonResp{
		ErrCode: 0,
		ErrMsg:  "User deleted!",
	}
	if err := token_verify.DeleteAdminTokenOnLogout(user.UserName, false); err != nil {
		errMsg := req.OperationID + " DeleteToken failed " + err.Error() + user.UserName + "Admin"
		log.NewError(req.OperationID, errMsg)
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

func (s *adminRPCServer) UpdateAdminUser(_ context.Context, req *admin.UpdateAdminReq) (*admin.UpdateAdminResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	// log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TOKEN", token_verify.Claims)
	resp := &admin.UpdateAdminResp{}
	user, err := walletdb.GetRegAdminUserByUserName(req.UserName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if user.ID == 0 {
		return nil, http.WrapError(constant.ErrAddInviteCodeUserNotExist) //Edit to user ID doesnt exist
	}
	role, err := walletdb.GetRegRolesByRoleName(req.RoleName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleName", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if role.ID == 0 {
		return nil, http.WrapError(constant.ErrRoleNameDoesntExist) //Edit to user ID doesnt exist
	}

	userID := user.ID
	// log.NewInfo("123", "TWO FACTOR S IS ", req.TwoFactorEnabled)
	updatedUser := db.AdminUser{
		// Password:         passwordNew,
		RoleId:           role.ID,
		Status:           int(req.Status), //int instead of boolean
		User2FAuthEnable: req.TwoFactorEnabled,
		Remarks:          req.Remarks,
		UpdateTime:       time.Now().Unix(),
		UpdateUser:       req.UpdateUser,
	}
	if req.Password != "" {
		passwordNew := req.Password
		newPasswordFirst := passwordNew + user.Salt
		passwordData := []byte(newPasswordFirst)
		has := md5.Sum(passwordData)
		passwordNew = fmt.Sprintf("%x", has)
		updatedUser.Password = passwordNew
	}

	err = walletdb.UpdateUser(userID, updatedUser, false, false)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAdmin failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UpdateAdminRPC...", req.String())

	//log the user out when getting banned
	if req.Status == 2 {
		if err := token_verify.DeleteAdminTokenOnLogout(user.UserName, false); err != nil {
			errMsg := req.OperationID + " DeleteToken failed " + err.Error() + user.UserName + "Admin"
			log.NewError(req.OperationID, errMsg)
		}
	}
	return resp, nil
}
func (s *adminRPCServer) UpdateAdminRole(_ context.Context, req *admin.UpdateAdminRoleRequest) (*admin.UpdateAdminRoleResponse, error) {
	resp := &admin.UpdateAdminRoleResponse{}

	role, err := walletdb.GetRegRolesByRoleName(req.RoleName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleName", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if role.ID == 0 {
		return nil, http.WrapError(constant.ErrRoleNameDoesntExist)
	}

	updatedRole := db.AdminRole{
		Description: req.Description,
		ActionIds:   req.ActionIDs,
		Remarks:     req.Remarks,
		UpdateTime:  time.Now().Unix(),
		UpdateUser:  req.UpdateUser,
	}

	err = walletdb.UpdateAdminRole(role.ID, updatedRole)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateRole failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UpdateRoleRPC...", req.String())

	return resp, nil
}

func (s *adminRPCServer) DeleteRole(_ context.Context, req *admin.DeleteAdminRoleRequest) (*admin.CommonResp, error) {

	role, err := walletdb.GetRegRolesByRoleName(req.RoleName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleName", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	if role.ID == 0 {
		return nil, http.WrapError(constant.ErrRoleNameDoesntExist)
	}
	used, err := walletdb.CheckRoleUsage(req.RoleName)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CheckRoleUsage", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	if used.ID != 0 {
		return nil, http.WrapError(constant.ErrRoleIsInUse)
	}
	deletedRole := db.AdminRole{
		DeleteTime: time.Now().Unix(),
		DeleteUser: req.DeleteUser,
		RoleName:   req.RoleName,
	}
	err = walletdb.DeleteRole(deletedRole)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteRole failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "DeleteRoleRPC...", req.String())
	resp := &admin.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Role deleted!",
	}
	return resp, nil
}

// admin login api
func (s *adminRPCServer) AdminLoginV2(_ context.Context, req *admin.AdminLoginReq) (*admin.AdminLoginResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &admin.AdminLoginResp{}
	user, err := walletdb.GetRegAdminUsrByUID(req.AdminID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "Get admin By UserID failed!", "adminID: ", req.AdminID, err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	password := req.Secret
	if !req.SecretHashd {
		newPasswordFirst := req.Secret + user.Salt
		passwordData := []byte(newPasswordFirst)
		has := md5.Sum(passwordData)
		password = fmt.Sprintf("%x", has)
	}

	if password == user.Password {
		if !user.User2FAuthEnable {
			req.GAuthTypeToken = false
		}
		token, expTime, err := token_verify.CreateToken(req.AdminID, constant.AdminPlatformID, req.GAuthTypeToken)
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", req.AdminID, err.Error())
			return resp, http.WrapError(constant.ErrTokenUnknown)
		}
		resp.Token = token
	}
	if resp.Token == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "failed")
		return resp, http.WrapError(constant.ErrTokenMalformed)
	}
	// resp.GAuthEnabled = user.TwoFactorEnabled
	resp.GAuthEnabled = true
	resp.GAuthSetupRequired = !user.TwoFactorEnabled && config.Config.AdminUser2FAuthEnable
	if !user.User2FAuthEnable {
		resp.GAuthSetupRequired = false
		resp.GAuthEnabled = false
	}
	if resp.GAuthSetupRequired {
		resp.GAuthSetupProvUri = genrateTOTPProvisUriQR(*user)
	}
	updatedUser := db.AdminUser{
		TwoFactorEnabled: true,
		User2FAuthEnable: true,
		UpdateTime:       time.Now().Unix(),
		UpdateUser:       user.UserName,
	}

	err = walletdb.UpdateUser(user.ID, updatedUser, true, false)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAdmin failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func BuildActions(ids []int, relations map[int][]int) []*admin.Action {
	actions := make([]*admin.Action, len(ids))

	for i, id := range ids {
		action, err := walletdb.GetRoleActionByActionID(id)
		if err != nil {
			return nil
		}

		c := admin.Action{ID: int32(id), Name: action.ActionName}
		if childIDs, ok := relations[id]; ok { // build child's children
			c.Children = BuildActions(childIDs, relations)
		}
		actions[i] = &c
	}
	return actions
}
func (s *adminRPCServer) GetAdminRoleActions(_ context.Context, req *admin.GetAdminActionsRequest) (*admin.Action, error) {
	actions, err := walletdb.MapActionParents()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminRoleActions failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	var actionIds [][]int
	for _, action := range actions {
		var actionData []int
		actionData = append(actionData, action.ID, action.Pid)
		actionIds = append(actionIds, actionData)
	}
	x := walletdb.RelationMap(actionIds)

	//check array size
	if len(x) != 0 {
		newResult := BuildActions(x[0], x)
		return newResult[0], err
	}
	return nil, err

}

func (s *adminRPCServer) GetAdminUser(_ context.Context, req *admin.GetAdminUserRequest) (*admin.GetAdminUserResponse, error) {
	resp := &admin.GetAdminUserResponse{}
	user, err := walletdb.GetRegAdminUserByUserName(req.UserName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName failed!", "adminID: ", user.ID, err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if user.ID == 0 {
		return nil, http.WrapError(constant.ErrUserNotExist)
	}

	resp.UserName = user.UserName
	if user.RoleId != 0 {
		role, err := walletdb.GetRegRolesByRoleID(user.RoleId)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleID failed!", "Role id: ", user.RoleId, err.Error())
			// return resp, http.WrapError(constant.ErrDB)
		}
		if role.ID != 0 {
			resp.RoleName = role.RoleName
			// return nil, http.WrapError(constant.ErrRoleNameDoesntExist)

			actions, err := walletdb.GetRoleActionIdsByRoleName(role.RoleName)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRoleActionsByRoleName failed!", "Actions array ", actions, err.Error())
				// return resp, http.WrapError(constant.ErrDB)
			}
			resp.Permissions = actions
		}
	}
	return resp, err

}

func (s *adminRPCServer) GetAccountInformation(_ context.Context, req *admin.GetAccountInformationReq) (*admin.GetAccountInformationResp, error) {
	resp := &admin.GetAccountInformationResp{}
	filters := make(map[string]string)
	if req.Sort != "" {
		filters["sort"] = req.Sort
	}
	if req.OrderBy != "" {
		filters["order_by"] = req.OrderBy
	}
	if req.From != "" {
		filters["to"] = req.To
		filters["from"] = req.From
	}
	if req.MerchantUid != "" {
		filters["merchant_uid"] = req.MerchantUid
	}
	if req.AccountAddress != "" {
		filters["account_address"] = req.AccountAddress
	}
	if req.CoinsType != "" {
		filters["coins_type"] = req.CoinsType
	} else {
		filters["coins_type"] = "all"
	}
	if req.AccountSource != "" {
		filters["account_source"] = req.AccountSource
	}

	accountInformation, count, coinsTotal, err := walletdb.GetAccountInformation(filters, req.Pagination.Page, req.Pagination.PageSize, req.Uid)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformation", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}

	coinsType := []string{"BTC", "ETH", "USDT-ERC20", "TRX", "USDT-TRC20"}
	var totalAssetUsd, totalAssetYuan, totalAssetEuro float64
	totalAssetUsd, totalAssetYuan, totalAssetEuro = 0, 0, 0
	currencies, _ := walletdb.GetCoinStatuses()
	btcIndex := constant.BTCCoin - 1
	ethIndex := constant.ETHCoin - 1
	trxIndex := constant.TRX - 1
	usdtIndex := constant.USDTERC20 - 1

	var unique []db.AccountInformation
	if len(accountInformation) > 0 {
		for _, v := range accountInformation {
			skip := false
			for _, u := range unique {
				if v.UUID == u.UUID {
					skip = true
					break
				}
			}
			if !skip {
				unique = append(unique, v)
			}
			// creationLogin := time.Unix(v.CreationTime, 0)
			// lastLogin := time.Unix(v.LastLoginTime, 0)
			accountAddresses := &admin.AccountAddresses{
				BtcPublicAddress: v.BtcPublicAddress,
				EthPublicAddress: v.EthPublicAddress,
				ErcPublicAddress: v.ErcPublicAddress,
				TrcPublicAddress: v.TrcPublicAddress,
				TrxPublicAddress: v.TrxPublicAddress,
			}
			lastLoginInformation := &admin.LoginInformation{
				LoginIp:       v.LastLoginIp,
				LoginRegion:   v.LastLoginRegion,
				LoginTerminal: v.LastLoginTerminal,
				LoginTime:     v.LastLoginTime,
			}
			creationLoginInformation := &admin.LoginInformation{
				LoginIp:       v.CreationLoginIp,
				LoginRegion:   v.CreationLoginRegion,
				LoginTerminal: v.CreationLoginTerminal,
				LoginTime:     v.CreationTime,
			}
			totalBalance := v.BtcBalance + v.EthBalance + v.TrxBalance + v.ErcBalance + v.TrcBalance
			var accountAssetUsd, accountAssetYuan, accountAssetEuro float64
			accountAssetUsd = (currencies[btcIndex].Usd * v.BtcBalance) + (currencies[ethIndex].Usd * v.EthBalance) + (currencies[usdtIndex].Usd * v.ErcBalance) + (currencies[usdtIndex].Usd * v.TrcBalance) + (currencies[trxIndex].Usd * v.TrxBalance)
			accountAssetYuan = (currencies[btcIndex].Yuan * v.BtcBalance) + (currencies[ethIndex].Yuan * v.EthBalance) + (currencies[usdtIndex].Yuan * v.ErcBalance) + (currencies[usdtIndex].Yuan * v.TrcBalance) + (currencies[trxIndex].Yuan * v.TrxBalance)
			accountAssetEuro = (currencies[btcIndex].Euro * v.BtcBalance) + (currencies[ethIndex].Euro * v.EthBalance) + (currencies[usdtIndex].Euro * v.ErcBalance) + (currencies[usdtIndex].Euro * v.TrcBalance) + (currencies[trxIndex].Euro * v.TrxBalance)
			accountAssets := &admin.AccountAsset{
				UsdAmount:  utils.RoundFloat(accountAssetUsd, 2),
				YuanAmount: utils.RoundFloat(accountAssetYuan, 2),
				EuroAmount: utils.RoundFloat(accountAssetEuro, 2),
			}
			Btc := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.BtcBalance),
				UsdBalance:  fmt.Sprintf("%.8f", currencies[btcIndex].Usd*v.BtcBalance),
				YuanBalance: fmt.Sprintf("%.8f", currencies[btcIndex].Yuan*v.BtcBalance),
				EuroBalance: fmt.Sprintf("%.8f", currencies[btcIndex].Euro*v.BtcBalance),
			}
			Eth := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.EthBalance),
				UsdBalance:  fmt.Sprintf("%.8f", currencies[ethIndex].Usd*v.EthBalance),
				YuanBalance: fmt.Sprintf("%.8f", currencies[ethIndex].Yuan*v.EthBalance),
				EuroBalance: fmt.Sprintf("%.8f", currencies[ethIndex].Euro*v.EthBalance),
			}
			Erc := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.ErcBalance),
				UsdBalance:  fmt.Sprintf("%.8f", currencies[usdtIndex].Usd*v.ErcBalance),
				YuanBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Yuan*v.ErcBalance),
				EuroBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Euro*v.ErcBalance),
			}
			Trc := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrcBalance),
				UsdBalance:  fmt.Sprintf("%.8f", currencies[usdtIndex].Usd*v.TrcBalance),
				YuanBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Yuan*v.TrcBalance),
				EuroBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Euro*v.TrcBalance),
			}
			Trx := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrxBalance),
				UsdBalance:  fmt.Sprintf("%.8f", currencies[trxIndex].Usd*v.TrxBalance),
				YuanBalance: fmt.Sprintf("%.8f", currencies[trxIndex].Yuan*v.TrxBalance),
				EuroBalance: fmt.Sprintf("%.8f", currencies[trxIndex].Euro*v.TrxBalance),
			}
			account := &admin.AccountInformation{
				ID:                       v.ID,
				Uid:                      v.UUID,
				MerchantUid:              v.MerchantUid,
				CoinsType:                coinsType,
				Addresses:                accountAddresses,
				AccountAssets:            accountAssets,
				Btc:                      Btc,
				Eth:                      Eth,
				Trx:                      Trx,
				Erc:                      Erc,
				Trc:                      Trc,
				TotalBalance:             utils.RoundFloat(totalBalance, 8),
				AccountSource:            v.AccountSource,
				CreationLoginInformation: creationLoginInformation,
				LastLoginInformation:     lastLoginInformation,
			}
			resp.Account = append(resp.Account, account)
		}
	}
	totalAssetUsd = (coinsTotal.Btc * currencies[btcIndex].Usd) + (coinsTotal.Eth * currencies[ethIndex].Usd) + (coinsTotal.Erc * currencies[usdtIndex].Usd) + (coinsTotal.Trx * currencies[trxIndex].Usd) + (coinsTotal.Trc * currencies[usdtIndex].Usd)
	totalAssetEuro = (coinsTotal.Btc * currencies[btcIndex].Euro) + (coinsTotal.Eth * currencies[ethIndex].Euro) + (coinsTotal.Erc * currencies[usdtIndex].Euro) + (coinsTotal.Trx * currencies[trxIndex].Euro) + (coinsTotal.Trc * currencies[usdtIndex].Euro)
	totalAssetYuan = (coinsTotal.Btc * currencies[btcIndex].Yuan) + (coinsTotal.Eth * currencies[ethIndex].Yuan) + (coinsTotal.Erc * currencies[usdtIndex].Yuan) + (coinsTotal.Trx * currencies[trxIndex].Yuan) + (coinsTotal.Trc * currencies[usdtIndex].Yuan)
	totalAssets := &admin.AccountAsset{
		UsdAmount:  utils.RoundFloat(totalAssetUsd, 2),
		YuanAmount: utils.RoundFloat(totalAssetYuan, 2),
		EuroAmount: utils.RoundFloat(totalAssetEuro, 2),
	}
	TotalBtc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", coinsTotal.Btc),
		UsdBalance:  fmt.Sprintf("%.8f", coinsTotal.Btc*currencies[btcIndex].Usd),
		YuanBalance: fmt.Sprintf("%.8f", coinsTotal.Btc*currencies[btcIndex].Yuan),
		EuroBalance: fmt.Sprintf("%.8f", coinsTotal.Btc*currencies[btcIndex].Euro),
	}
	TotalEth := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", coinsTotal.Eth),
		UsdBalance:  fmt.Sprintf("%.8f", coinsTotal.Eth*currencies[ethIndex].Usd),
		YuanBalance: fmt.Sprintf("%.8f", coinsTotal.Eth*currencies[ethIndex].Yuan),
		EuroBalance: fmt.Sprintf("%.8f", coinsTotal.Eth*currencies[ethIndex].Euro),
	}
	TotalErc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", coinsTotal.Erc),
		UsdBalance:  fmt.Sprintf("%.8f", coinsTotal.Erc*currencies[usdtIndex].Usd),
		YuanBalance: fmt.Sprintf("%.8f", coinsTotal.Erc*currencies[usdtIndex].Yuan),
		EuroBalance: fmt.Sprintf("%.8f", coinsTotal.Erc*currencies[usdtIndex].Euro),
	}
	TotalTrc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", coinsTotal.Trc),
		UsdBalance:  fmt.Sprintf("%.8f", coinsTotal.Trc*currencies[usdtIndex].Usd),
		YuanBalance: fmt.Sprintf("%.8f", coinsTotal.Trc*currencies[usdtIndex].Yuan),
		EuroBalance: fmt.Sprintf("%.8f", coinsTotal.Trc*currencies[usdtIndex].Euro),
	}
	TotalTrx := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", coinsTotal.Trx),
		UsdBalance:  fmt.Sprintf("%.8f", coinsTotal.Trx*currencies[trxIndex].Usd),
		YuanBalance: fmt.Sprintf("%.8f", coinsTotal.Trx*currencies[trxIndex].Yuan),
		EuroBalance: fmt.Sprintf("%.8f", coinsTotal.Trx*currencies[trxIndex].Euro),
	}
	resp.TotalAccounts = count
	resp.TotalAssets = totalAssets
	resp.BtcTotal = TotalBtc
	resp.EthTotal = TotalEth
	resp.ErcTotal = TotalErc
	resp.TrcTotal = TotalTrc
	resp.TrxTotal = TotalTrx
	return resp, nil
}

func (s *adminRPCServer) GetFundsLog(c context.Context, req *admin.GetFundsLogReq) (*admin.GetFundsLogResp, error) {
	resp := &admin.GetFundsLogResp{}
	filters := make(map[string]string)
	if req.From != "" {
		filters["to"] = req.To
		filters["from"] = req.From
	}
	if req.TransactionType != "" {
		filters["transaction_type"] = req.TransactionType
	}
	if req.UserAddress != "" {
		filters["user_address"] = req.UserAddress
	}
	if req.OppositeAddress != "" {
		filters["opposite_address"] = req.OppositeAddress
	}
	if req.CoinsType != "" {
		filters["coins_type"] = req.CoinsType
	} else {
		filters["coins_type"] = "btc"
	}
	if req.State != "" {
		filters["state"] = req.State
	}
	if req.Txid != "" {
		filters["txid"] = req.Txid
	}
	if req.MerchantUid != "" {
		filters["merchant_uid"] = req.MerchantUid
	}
	fundsLog, count, err := walletdb.GetFundsLog(filters, req.Pagination.Page, req.Pagination.PageSize, req.Uid)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	//fetch all coins
	currencies, _ := walletdb.GetCoinStatuses()
	if len(fundsLog) > 0 {
		for _, v := range fundsLog {
			coinsAmount, _ := v.AmountOfCoins.Float64()
			// networkFee, _ := v.NetworkFee.Float64()
			gasPrice := utils.Wei2Eth_str(v.GasPrice.BigInt())
			gasPriceFloat, _ := strconv.ParseFloat(gasPrice, 64)
			networkFee := gasPriceFloat * float64(v.GasUsed)

			//get the coin id
			coinId := utils.GetCoinType(v.CoinType) - 1

			if v.CoinType == utils.GetCoinName(constant.USDTTRC20) || v.CoinType == utils.GetCoinName(constant.TRX) {
				networkFee = v.NetworkFee.InexactFloat64()
			}
			totalCoins := v.AmountOfCoins.Add(v.NetworkFee)
			totalCoinsAmount, _ := totalCoins.Float64()

			//get coin rate based on the coin id
			totalUsdCoins := currencies[coinId].Usd * coinsAmount
			totalYuanCoins := currencies[coinId].Yuan * coinsAmount
			totalEuroCoins := currencies[coinId].Euro * coinsAmount
			usdFee := currencies[coinId].Usd * networkFee
			yuanFee := currencies[coinId].Yuan * networkFee
			euroFee := currencies[coinId].Euro * networkFee
			if v.CoinType == utils.GetCoinName(constant.USDTERC20) {
				coinId = utils.GetCoinType("ETH") - 1
				ethRate := currencies[coinId]

				usdFee = ethRate.Usd * networkFee
				yuanFee = ethRate.Yuan * networkFee
				euroFee = ethRate.Euro * networkFee
			} else if v.CoinType == utils.GetCoinName(constant.USDTTRC20) {
				coinId = utils.GetCoinType("TRX") - 1
				trxRate := currencies[coinId]
				usdFee = trxRate.Usd * networkFee
				yuanFee = trxRate.Yuan * networkFee
				euroFee = trxRate.Euro * networkFee
			}
			totalUsdTransfered := totalUsdCoins + usdFee
			totalYuanTransfered := totalYuanCoins + yuanFee
			totalEuroTransfered := totalEuroCoins + euroFee
			// creationTime := time.Unix(v.CreationTime, 0)
			// confirmationTime := time.Unix(v.ConfirmationTime, 0)
			//state 0 failed, 1 success, 2 pending
			state := "failed"
			if v.State == 1 {
				state = "success"
			}
			coins := fmt.Sprintf("%.8f", utils.RoundFloat(coinsAmount, 8))
			totalTransferedCoins := fmt.Sprintf("%.8f", utils.RoundFloat(totalCoinsAmount, 8))
			totalNetworkFee := fmt.Sprintf("%.8f", utils.RoundFloat(networkFee, 8))
			UsdAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalUsdCoins, 8))
			YuanAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalYuanCoins, 8))
			EuroAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalEuroCoins, 8))
			TotalUsdTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalUsdTransfered, 8))
			TotalYuanTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalYuanTransfered, 8))
			TotalEuroTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalEuroTransfered, 8))
			UsdFee := fmt.Sprintf("%.2f", utils.RoundFloat(usdFee, 8))
			YuanFee := fmt.Sprintf("%.2f", utils.RoundFloat(yuanFee, 8))
			EuroFee := fmt.Sprintf("%.2f", utils.RoundFloat(euroFee, 8))
			log.NewInfo(req.OperationID, "amount of coins", coins)
			fundLog := &admin.FundsLog{
				ID:                   v.ID,
				Txid:                 v.Txid,
				Uid:                  v.UID,
				MerchantUid:          v.MerchantUid,
				TransactionType:      v.TransactionType,
				UserAddress:          v.UserAddress,
				OppositeAddress:      v.OppositeAddress,
				CoinType:             v.CoinType,
				AmountOfCoins:        coins,
				UsdAmount:            UsdAmount,
				YuanAmount:           YuanAmount,
				EuroAmount:           EuroAmount,
				NetworkFee:           "0",
				UsdNetworkFee:        "0",
				EuroNetworkFee:       "0",
				YuanNetworkFee:       "0",
				TotalCoinsTransfered: coins,
				TotalUsdTransfered:   UsdAmount,
				TotalYuanTransfered:  YuanAmount,
				TotalEuroTransfered:  EuroAmount,
				CreationTime:         v.CreationTime,
				State:                state,
				ConfirmationTime:     v.ConfirmationTime,
			}
			if v.TransactionType != constant.TransactionTypeReceiveString {
				fundLog.NetworkFee = totalNetworkFee
				fundLog.UsdNetworkFee = UsdFee
				fundLog.YuanNetworkFee = YuanFee
				fundLog.EuroNetworkFee = EuroFee
				fundLog.TotalCoinsTransfered = totalTransferedCoins
				fundLog.TotalUsdTransfered = TotalUsdTransfered
				fundLog.TotalYuanTransfered = TotalYuanTransfered
				fundLog.TotalEuroTransfered = TotalEuroTransfered
			}
			log.NewInfo(req.OperationID, "fund log amount of coins", fundLog.AmountOfCoins)
			resp.FundLog = append(resp.FundLog, fundLog)
		}
	}
	resp.TotalFundLogs = count

	return resp, nil
}
func (s *adminRPCServer) GetReceiveDetails(c context.Context, req *admin.GetReceiveDetailsReq) (*admin.GetReceiveDetailsResp, error) {
	resp := &admin.GetReceiveDetailsResp{}
	filters := make(map[string]string)
	if req.From != "" {
		filters["to"] = req.To
		filters["from"] = req.From
	}
	if req.MerchantUid != "" {
		filters["merchant_uid"] = req.MerchantUid
	}
	if req.ReceivingAddress != "" {
		filters["receiving_address"] = req.ReceivingAddress
	}
	if req.DepositAddress != "" {
		filters["deposit_address"] = req.DepositAddress
	}
	if req.CoinsType != "" {
		filters["coins_type"] = req.CoinsType
	}
	if req.Txid != "" {
		filters["txid"] = req.Txid
	}
	receiveDetails, count, coinsTotal, err := walletdb.GetReceiveDetails(filters, req.Pagination.Page, req.Pagination.PageSize, req.Uid)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	fundsTotal, err := walletdb.GetTotalFunds(constant.TransactionTypeReceiveString)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	var totalUsd, totalYuan, totalEuro, grandTotalUsd, grandTotalYuan, grandTotalEuro float64
	totalUsd, totalYuan, totalEuro, grandTotalUsd, grandTotalYuan, grandTotalEuro = 0, 0, 0, 0, 0, 0
	currencies, _ := walletdb.GetCoinStatuses()
	if len(receiveDetails) > 0 {
		for _, v := range receiveDetails {
			coinIndex := utils.GetCoinType(v.CoinType) - 1
			coinsAmount, _ := v.AmountOfCoins.Float64()
			totalUsdCoins := currencies[coinIndex].Usd * coinsAmount
			totalYuanCoins := currencies[coinIndex].Yuan * coinsAmount
			totalEuroCoins := currencies[coinIndex].Euro * coinsAmount
			coins := fmt.Sprintf("%.8f", utils.RoundFloat(coinsAmount, 8))
			UsdAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalUsdCoins, 8))
			YuanAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalYuanCoins, 8))
			EuroAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalEuroCoins, 8))
			receiveDetail := &admin.ReceiveDetails{
				ID:               v.ID,
				Uid:              v.UID,
				MerchantUid:      v.MerchantUid,
				ReceivingAddress: v.UserAddress,
				CoinType:         v.CoinType,
				AmountOfReceived: coins,
				UsdAmount:        UsdAmount,
				YuanAmount:       YuanAmount,
				EuroAmount:       EuroAmount,
				DepositAddress:   v.OppositeAddress,
				Txid:             v.Txid,
				CreationTime:     v.CreationTime,
			}
			resp.ReceiveDetail = append(resp.ReceiveDetail, receiveDetail)
		}
		for _, value := range fundsTotal {
			if value.Name != "" {
				coinIndex := utils.GetCoinType(value.Name) - 1
				grandTotalUsd += value.Total * currencies[coinIndex].Usd
				grandTotalYuan += value.Total * currencies[coinIndex].Yuan
				grandTotalEuro += value.Total * currencies[coinIndex].Euro
			}
		}

		//getting the totals for assets based on the given filters, using "coinsTotal"
		coin := strings.ToUpper(req.CoinsType)
		switch coin {
		case utils.GetCoinName(constant.BTCCoin):
			coinIndex := constant.BTCCoin - 1
			totalUsd = (coinsTotal.Btc * currencies[coinIndex].Usd)
			totalEuro = (coinsTotal.Btc * currencies[coinIndex].Euro)
			totalYuan = (coinsTotal.Btc * currencies[coinIndex].Yuan)

		case utils.GetCoinName(constant.ETHCoin):
			coinIndex := constant.ETHCoin - 1
			totalUsd = (coinsTotal.Eth * currencies[coinIndex].Usd)
			totalEuro = (coinsTotal.Eth * currencies[coinIndex].Euro)
			totalYuan = (coinsTotal.Eth * currencies[coinIndex].Yuan)

		case utils.GetCoinName(constant.USDTERC20):
			coinIndex := constant.USDTERC20 - 1
			totalUsd = (coinsTotal.Erc * currencies[coinIndex].Usd)
			totalEuro = (coinsTotal.Erc * currencies[coinIndex].Euro)
			totalYuan = (coinsTotal.Erc * currencies[coinIndex].Yuan)
		case utils.GetCoinName(constant.TRX):
			coinIndex := constant.TRX - 1
			totalUsd = (coinsTotal.Trx * currencies[coinIndex].Usd)
			totalEuro = (coinsTotal.Trx * currencies[coinIndex].Euro)
			totalYuan = (coinsTotal.Trx * currencies[coinIndex].Yuan)
		case utils.GetCoinName(constant.USDTTRC20):
			coinIndex := constant.USDTTRC20 - 1
			totalUsd = (coinsTotal.Trc * currencies[coinIndex].Usd)
			totalEuro = (coinsTotal.Trc * currencies[coinIndex].Euro)
			totalYuan = (coinsTotal.Trc * currencies[coinIndex].Yuan)

		default:
			btcIndex := constant.BTCCoin - 1
			ethIndex := constant.ETHCoin - 1
			trxIndex := constant.TRX - 1
			usdtIndex := constant.USDTERC20 - 1
			totalUsd = (coinsTotal.Btc * currencies[btcIndex].Usd) + (coinsTotal.Eth * currencies[ethIndex].Usd) + (coinsTotal.Erc * currencies[usdtIndex].Usd) + (coinsTotal.Trx * currencies[trxIndex].Usd) + (coinsTotal.Trc * currencies[usdtIndex].Usd)
			totalYuan = (coinsTotal.Btc * currencies[btcIndex].Yuan) + (coinsTotal.Eth * currencies[ethIndex].Yuan) + (coinsTotal.Erc * currencies[usdtIndex].Yuan) + (coinsTotal.Trx * currencies[trxIndex].Yuan) + (coinsTotal.Trc * currencies[usdtIndex].Yuan)
			totalEuro = (coinsTotal.Btc * currencies[btcIndex].Euro) + (coinsTotal.Eth * currencies[ethIndex].Euro) + (coinsTotal.Erc * currencies[usdtIndex].Euro) + (coinsTotal.Trx * currencies[trxIndex].Euro) + (coinsTotal.Trc * currencies[usdtIndex].Euro)

		}

	}
	resp.TotalDetails = count
	resp.TotalAmountReceivedUsd = utils.RoundFloat(totalUsd, 2)
	resp.TotalAmountReceivedYuan = utils.RoundFloat(totalYuan, 2)
	resp.TotalAmountReceivedEuro = utils.RoundFloat(totalEuro, 2)
	resp.GrandTotalUsd = utils.RoundFloat(grandTotalUsd, 2)
	resp.GrandTotalYuan = utils.RoundFloat(grandTotalYuan, 2)
	resp.GrandTotalEuro = utils.RoundFloat(grandTotalEuro, 2)
	return resp, nil
}

func (s *adminRPCServer) GetTransferDetails(c context.Context, req *admin.GetTransferDetailsReq) (*admin.GetTransferDetailsResp, error) {
	resp := &admin.GetTransferDetailsResp{}
	filters := make(map[string]string)
	if req.OrderBy != "" {
		filters["order_by"] = req.OrderBy
	}
	if req.From != "" {
		filters["from"] = req.From
		filters["to"] = req.To
	}
	if req.TransferAddress != "" {
		filters["transfer_address"] = req.TransferAddress
	}
	if req.ReceivingAddress != "" {
		filters["receiving_address"] = req.ReceivingAddress
	}
	if req.CoinsType != "" {
		filters["coins_type"] = req.CoinsType
	} else {
		filters["coins_type"] = "all"
	}
	var state int32
	state = -1
	if req.State != "" {
		switch req.State {
		case "fail":
			state = 0
		case "success":
			state = 1
		case "pending":
			state = 2
		case "all":
			state = -1
		}
	}
	if req.Txid != "" {
		filters["txid"] = req.Txid
	}
	if req.MerchantUid != "" {
		filters["merchant_uid"] = req.MerchantUid
	}
	transferDetails, count, coinsTotal, err := walletdb.GetTransferDetails(filters, req.Pagination.Page, req.Pagination.PageSize, req.Uid, state)
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "transferDetails", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	var totalUsd, totalYuan, totalEuro, totalUsdFee, totalYuanFee, totalEuroFee, totalTransferUsd, totalTransferYuan, totalTransferEuro float64
	totalUsd, totalYuan, totalEuro, totalUsdFee, totalYuanFee, totalEuroFee = 0, 0, 0, 0, 0, 0
	currencies, _ := walletdb.GetCoinStatuses()
	if len(transferDetails) > 0 {
		for _, v := range transferDetails {
			coinIndex := utils.GetCoinType(v.CoinType) - 1

			txState := ""
			switch v.State {
			case 0:
				txState = "fail"
			case 1:
				txState = "success"
			case 2:
				txState = "pending"
			}

			coinsAmount, _ := v.AmountOfCoins.Float64()
			// networkFee, _ := v.NetworkFee.Float64()
			gasPrice := utils.Wei2Eth_str(v.GasPrice.BigInt())
			gasPriceFloat, _ := strconv.ParseFloat(gasPrice, 64)
			networkFee := gasPriceFloat * float64(v.GasUsed)
			totalUsdCoins := currencies[coinIndex].Usd * coinsAmount
			totalYuanCoins := currencies[coinIndex].Yuan * coinsAmount
			totalEuroCoins := currencies[coinIndex].Euro * coinsAmount
			if v.CoinType == utils.GetCoinName(constant.USDTTRC20) || v.CoinType == utils.GetCoinName(constant.TRX) {
				networkFee = v.NetworkFee.InexactFloat64()
			}
			usdFee := currencies[coinIndex].Usd * networkFee
			yuanFee := currencies[coinIndex].Yuan * networkFee
			euroFee := currencies[coinIndex].Euro * networkFee
			if v.CoinType == utils.GetCoinName(constant.USDTERC20) {
				coinIndex := utils.GetCoinType("ETH") - 1
				usdFee = currencies[coinIndex].Usd * networkFee
				yuanFee = currencies[coinIndex].Yuan * networkFee
				euroFee = currencies[coinIndex].Euro * networkFee
			} else if v.CoinType == utils.GetCoinName(constant.USDTTRC20) {
				coinIndex := utils.GetCoinType("TRX") - 1
				usdFee = currencies[coinIndex].Usd * networkFee
				yuanFee = currencies[coinIndex].Yuan * networkFee
				euroFee = currencies[coinIndex].Euro * networkFee
			}
			// ConfirmationTime := time.Unix(v.ConfirmationTime, 0)
			// CreationTime := time.Unix(v.CreationTime, 0)
			totalcoins := v.AmountOfCoins.Add(v.NetworkFee)
			totalCoinsAmount, _ := totalcoins.Float64()
			coins := fmt.Sprintf("%.8f", utils.RoundFloat(coinsAmount, 8))
			totalTransferedCoins := fmt.Sprintf("%.8f", utils.RoundFloat(totalCoinsAmount, 8))
			totalNetworkFee := fmt.Sprintf("%.8f", utils.RoundFloat(networkFee, 8))
			UsdAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalUsdCoins, 8))
			YuanAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalYuanCoins, 8))
			EuroAmount := fmt.Sprintf("%.2f", utils.RoundFloat(totalEuroCoins, 8))
			TotalUsdTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalUsdCoins+usdFee, 8))
			TotalYuanTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalYuanCoins+yuanFee, 8))
			TotalEuroTransfered := fmt.Sprintf("%.2f", utils.RoundFloat(totalEuroCoins+euroFee, 8))
			UsdFee := fmt.Sprintf("%.2f", utils.RoundFloat(usdFee, 8))
			YuanFee := fmt.Sprintf("%.2f", utils.RoundFloat(yuanFee, 8))
			EuroFee := fmt.Sprintf("%.2f", utils.RoundFloat(euroFee, 8))

			transferDetail := &admin.FundsLog{
				ID:                   v.ID,
				Uid:                  v.UID,
				MerchantUid:          v.MerchantUid,
				Txid:                 v.Txid,
				UserAddress:          v.UserAddress,
				CoinType:             v.CoinType,
				TransactionType:      v.TransactionType,
				AmountOfCoins:        coins,
				UsdAmount:            UsdAmount,
				YuanAmount:           YuanAmount,
				EuroAmount:           EuroAmount,
				OppositeAddress:      v.OppositeAddress,
				NetworkFee:           totalNetworkFee,
				UsdNetworkFee:        UsdFee,
				YuanNetworkFee:       YuanFee,
				EuroNetworkFee:       EuroFee,
				TotalUsdTransfered:   TotalUsdTransfered,
				TotalYuanTransfered:  TotalYuanTransfered,
				TotalEuroTransfered:  TotalEuroTransfered,
				TotalCoinsTransfered: totalTransferedCoins,
				CreationTime:         v.CreationTime,
				ConfirmationTime:     v.ConfirmationTime,
				State:                txState,
			}
			resp.TransferDetail = append(resp.TransferDetail, transferDetail)
		}
		fundsTotal, err := walletdb.GetTotalFunds(constant.TransactionTypeSendString)
		if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog", err.Error())
			return resp, http.WrapError(constant.ErrDB)
		}
		for _, value := range fundsTotal {
			if value.Name != "" {
				coinIndex := utils.GetCoinType(value.Name) - 1
				totalTransferUsd += value.Total * currencies[coinIndex].Usd
				totalTransferYuan += value.Total * currencies[coinIndex].Yuan
				totalTransferEuro += value.Total * currencies[coinIndex].Euro
			}
		}

		btcIndex := constant.BTCCoin - 1
		ethIndex := constant.ETHCoin - 1
		trxIndex := constant.TRX - 1
		usdtIndex := constant.USDTERC20 - 1
		coin := strings.ToUpper(req.CoinsType)
		switch coin {
		case utils.GetCoinName(constant.BTCCoin):
			totalUsd = (coinsTotal.Btc * currencies[btcIndex].Usd)
			totalYuan = (coinsTotal.Btc * currencies[btcIndex].Yuan)
			totalEuro = (coinsTotal.Btc * currencies[btcIndex].Euro)
			totalUsdFee = (coinsTotal.BtcFee * currencies[btcIndex].Usd)
			totalYuanFee = (coinsTotal.BtcFee * currencies[btcIndex].Yuan)
			totalEuroFee = (coinsTotal.BtcFee * currencies[btcIndex].Euro)
		case utils.GetCoinName(constant.ETHCoin):
			totalUsd = (coinsTotal.Eth * currencies[ethIndex].Usd)
			totalEuro = (coinsTotal.Eth * currencies[ethIndex].Euro)
			totalYuan = (coinsTotal.Eth * currencies[ethIndex].Yuan)
			totalUsdFee = (coinsTotal.EthFee * currencies[ethIndex].Usd)
			totalEuroFee = (coinsTotal.EthFee * currencies[ethIndex].Euro)
			totalYuanFee = (coinsTotal.EthFee * currencies[ethIndex].Yuan)
		case utils.GetCoinName(constant.USDTERC20):
			totalUsd = (coinsTotal.Erc * currencies[usdtIndex].Usd)
			totalYuan = (coinsTotal.Erc * currencies[usdtIndex].Yuan)
			totalEuro = (coinsTotal.Erc * currencies[usdtIndex].Euro)
			totalUsdFee = (coinsTotal.ErcFee * currencies[ethIndex].Usd)
			totalEuroFee = (coinsTotal.ErcFee * currencies[ethIndex].Euro)
			totalYuanFee = (coinsTotal.ErcFee * currencies[ethIndex].Yuan)
		case utils.GetCoinName(constant.TRX):
			totalUsd = (coinsTotal.Trx * currencies[trxIndex].Usd)
			totalEuro = (coinsTotal.Trx * currencies[trxIndex].Euro)
			totalYuan = (coinsTotal.Trx * currencies[trxIndex].Yuan)
			totalUsdFee = (coinsTotal.TrxFee * currencies[trxIndex].Usd)
			totalEuroFee = (coinsTotal.TrxFee * currencies[trxIndex].Euro)
			totalYuanFee = (coinsTotal.TrxFee * currencies[trxIndex].Yuan)
		case utils.GetCoinName(constant.USDTTRC20):
			totalUsd = (coinsTotal.Trc * currencies[usdtIndex].Usd)
			totalEuro = (coinsTotal.Trc * currencies[usdtIndex].Euro)
			totalYuan = (coinsTotal.Trc * currencies[usdtIndex].Yuan)
			totalUsdFee = (coinsTotal.TrcFee * currencies[trxIndex].Usd)
			totalEuroFee = (coinsTotal.TrcFee * currencies[trxIndex].Euro)
			totalYuanFee = (coinsTotal.TrcFee * currencies[trxIndex].Yuan)
		default:
			//getting the totals for assets based on the given filters, using "coinsTotal"
			totalUsd = (coinsTotal.Btc * currencies[btcIndex].Usd) + (coinsTotal.Eth * currencies[ethIndex].Usd) + (coinsTotal.Erc * currencies[usdtIndex].Usd) + (coinsTotal.Trx * currencies[trxIndex].Usd) + (coinsTotal.Trc * currencies[usdtIndex].Usd)
			totalYuan = (coinsTotal.Btc * currencies[btcIndex].Yuan) + (coinsTotal.Eth * currencies[ethIndex].Yuan) + (coinsTotal.Erc * currencies[usdtIndex].Yuan) + (coinsTotal.Trx * currencies[trxIndex].Yuan) + (coinsTotal.Trc * currencies[usdtIndex].Yuan)
			totalEuro = (coinsTotal.Btc * currencies[btcIndex].Euro) + (coinsTotal.Eth * currencies[ethIndex].Euro) + (coinsTotal.Erc * currencies[usdtIndex].Euro) + (coinsTotal.Trx * currencies[trxIndex].Euro) + (coinsTotal.Trc * currencies[usdtIndex].Euro)

			totalUsdFee = (coinsTotal.BtcFee * currencies[btcIndex].Usd) + (coinsTotal.EthFee * currencies[ethIndex].Usd) + (coinsTotal.ErcFee * currencies[ethIndex].Usd) + (coinsTotal.TrxFee * currencies[trxIndex].Usd) + (coinsTotal.TrcFee * currencies[trxIndex].Usd)
			totalYuanFee = (coinsTotal.BtcFee * currencies[btcIndex].Yuan) + (coinsTotal.EthFee * currencies[ethIndex].Yuan) + (coinsTotal.ErcFee * currencies[ethIndex].Yuan) + (coinsTotal.TrxFee * currencies[trxIndex].Yuan) + (coinsTotal.TrcFee * currencies[trxIndex].Yuan)
			totalEuroFee = (coinsTotal.BtcFee * currencies[btcIndex].Euro) + (coinsTotal.EthFee * currencies[ethIndex].Euro) + (coinsTotal.ErcFee * currencies[ethIndex].Euro) + (coinsTotal.TrxFee * currencies[trxIndex].Euro) + (coinsTotal.TrcFee * currencies[trxIndex].Euro)
		}

		//total currency without fees
		resp.TotalAmountTransferedUsd = utils.RoundFloat(totalUsd, 2)
		resp.TotalAmountTransferedYuan = utils.RoundFloat(totalYuan, 2)
		resp.TotalAmountTransferedEuro = utils.RoundFloat(totalEuro, 2)

		//total fees
		resp.TotalFeeAmountUsd = utils.RoundFloat(totalUsdFee, 2)
		resp.TotalFeeAmountYuan = utils.RoundFloat(totalYuanFee, 2)
		resp.TotalFeeAmountEuro = utils.RoundFloat(totalEuroFee, 2)

		//total of currency + fee
		resp.GrandTotalUsd = utils.RoundFloat(totalUsd+totalUsdFee, 2)
		resp.GrandTotalYuan = utils.RoundFloat(totalYuan+totalYuanFee, 2)
		resp.GrandTotalEuro = utils.RoundFloat(totalEuro+totalEuroFee, 2)

		//total for all the page
		resp.TotalTransferUsd = utils.RoundFloat(totalTransferUsd, 2)
		resp.TotalTransferYuan = utils.RoundFloat(totalTransferYuan, 2)
		resp.TotalTransferEuro = utils.RoundFloat(totalTransferEuro, 2)
	}
	resp.TotalDetails = count
	return resp, nil
}
func (s *adminRPCServer) ResetGoogleKey(c context.Context, req *admin.ResetGoogleKeyReq) (*admin.CommonResp, error) {
	user, err := walletdb.GetRegAdminUserByUserName(req.UserName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegAdminUserByUserName failed!", "adminID: ", user.ID, err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	err = walletdb.ResetGooogleKey(user.UserName, req.UpdateUser)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddAdminUser failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	resp := &admin.CommonResp{
		ErrCode: 0,
		ErrMsg:  "Key reset!",
	}
	return resp, err
}
func (s *adminRPCServer) GetRoleActions(_ context.Context, req *admin.GetRoleActionsReq) (*admin.GetRoleActionsResp, error) {
	resp := &admin.GetRoleActionsResp{}

	role, err := walletdb.GetRegRolesByRoleName(req.RoleName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRegRolesByRoleName failed!", err.Error())
	}
	if role.ID != 0 {
		actions, err := walletdb.GetRoleActionIdsByRoleName(role.RoleName)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetRoleActionIdsByRoleName failed!", "Actions array ", actions, err.Error())
		}
		resp.Actions = actions
	}
	return resp, err
}
func (s *adminRPCServer) GetCurrencies(c context.Context, req *admin.GetCurrenciesReq) (*admin.GetCurrenciesResp, error) {
	resp := &admin.GetCurrenciesResp{}

	currencies, count, err := walletdb.GetCurrencies(req.Pagination.Page, req.Pagination.PageSize)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCurrencies failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}

	for _, v := range currencies {
		// updateTime := time.Unix(v.UpdateTime, 0)
		currency := &admin.Currency{
			ID:             int32(v.ID),
			CoinType:       v.Coin,
			LastEditedTime: v.UpdateTime,
			Editor:         v.UpdateUser,
			State:          int32(v.State),
		}
		resp.Currencies = append(resp.Currencies, currency)
	}
	resp.TotalCurrencies = count
	return resp, err
}
func (s *adminRPCServer) UpdateCurrency(c context.Context, req *admin.UpdateCurrencyReq) (*admin.CommonResp, error) {
	resp := &admin.CommonResp{}
	updatedCurrency := db.CoinCurrencyValues{
		State:      int8(req.State),
		UpdateUser: req.UpdateUser,
		UpdateTime: time.Now().Unix(),
	}
	err := walletdb.UpdateCurrency(req.CurrencyId, updatedCurrency)

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateCurrency failed", err.Error())
		return nil, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "UpdateCurrencyRPC...", req.String())

	return resp, nil
}

func (s *adminRPCServer) GetOperationalReport(c context.Context, req *admin.GetOperationalReportReq) (*admin.GetOperationalReportResp, error) {
	resp := &admin.GetOperationalReportResp{}
	operation, count, err := walletdb.GetOperationalReport(req.From, req.To, req.Pagination.Page, req.Pagination.PageSize)
	currencies, _ := walletdb.GetCoinStatuses()
	btcIndex := constant.BTCCoin - 1
	ethIndex := constant.ETHCoin - 1
	trxIndex := constant.TRX - 1
	usdtIndex := constant.USDTERC20 - 1

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCoinStatuses", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if len(operation) > 0 {
		var grandTotalTransferedUsd, grandTotalReceivedUsd, grandTotalFeeUsd, grandTotalTransferedEuro, grandTotalReceivedEuro, grandTotalFeeEuro, grandTotalTransferedYuan, grandTotalReceivedYuan, grandTotalFeeYuan float64
		var grandTotalTransferedBtc, grandTotalTransferedEth, grandTotalTransferedErc, grandTotalTransferedTrc, grandTotalTransferedTrx float64
		var grandTotalReceivedBtc, grandTotalReceivedEth, grandTotalReceivedErc, grandTotalReceivedTrc, grandTotalReceivedTrx, grandTotalFeeBtc, grandTotalFeeEth, grandTotalFeeErc, grandTotalFeeTrc, grandTotalFeeTrx float64
		var totalAssetUsd, totalAssetYuan, totalAssetEuro float64
		var totalUsers int64
		for _, v := range operation {
			layout := "2006-01-02"
			t, err := time.Parse(layout, v.Date)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(t.Unix())
			newUsers, err := walletdb.GetCreatedUsersCountPerDay(t.Unix())
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCreatedUsersCountPerDay", err.Error())
				return resp, http.WrapError(constant.ErrDB)
			}
			totalTransferUsd := (v.BtcTransfer * currencies[btcIndex].Usd) + (v.EthTransfer * currencies[ethIndex].Usd) + (v.ErcTransfer * currencies[usdtIndex].Usd) + (v.TrcTransfer * currencies[usdtIndex].Usd) + (v.TrxTransfer * currencies[trxIndex].Usd)
			totalTransferYuan := (v.BtcTransfer * currencies[btcIndex].Yuan) + (v.EthTransfer * currencies[ethIndex].Yuan) + (v.ErcTransfer * currencies[usdtIndex].Yuan) + (v.TrcTransfer * currencies[usdtIndex].Yuan) + (v.TrxTransfer * currencies[trxIndex].Yuan)
			totalTransferEuro := (v.BtcTransfer * currencies[btcIndex].Euro) + (v.EthTransfer * currencies[ethIndex].Euro) + (v.ErcTransfer * currencies[usdtIndex].Euro) + (v.TrcTransfer * currencies[usdtIndex].Euro) + (v.TrxTransfer * currencies[trxIndex].Euro)
			totalReceivedUsd := (v.BtcReceived * currencies[btcIndex].Usd) + (v.EthReceived * currencies[ethIndex].Usd) + (v.ErcReceived * currencies[usdtIndex].Usd) + (v.TrcReceived * currencies[usdtIndex].Usd) + (v.TrxReceived * currencies[trxIndex].Usd)
			totalReceivedYuan := (v.BtcReceived * currencies[btcIndex].Yuan) + (v.EthReceived * currencies[ethIndex].Yuan) + (v.ErcReceived * currencies[usdtIndex].Yuan) + (v.TrcReceived * currencies[usdtIndex].Yuan) + (v.TrxReceived * currencies[trxIndex].Yuan)
			totalReceivedEuro := (v.BtcReceived * currencies[btcIndex].Euro) + (v.EthReceived * currencies[ethIndex].Euro) + (v.ErcReceived * currencies[usdtIndex].Euro) + (v.TrcReceived * currencies[usdtIndex].Euro) + (v.TrxReceived * currencies[trxIndex].Euro)

			totalFeeUsd := (v.BtcFee * currencies[btcIndex].Usd) + (v.EthFee * currencies[ethIndex].Usd) + (v.ErcFee * currencies[ethIndex].Usd) + (v.TrcFee * currencies[trxIndex].Usd) + (v.TrxFee * currencies[trxIndex].Usd)
			totalFeeYuan := (v.BtcFee * currencies[btcIndex].Yuan) + (v.EthFee * currencies[ethIndex].Yuan) + (v.ErcFee * currencies[ethIndex].Yuan) + (v.TrcFee * currencies[trxIndex].Yuan) + (v.TrxFee * currencies[trxIndex].Yuan)
			totalFeeEuro := (v.BtcFee * currencies[btcIndex].Euro) + (v.EthFee * currencies[ethIndex].Euro) + (v.ErcFee * currencies[ethIndex].Euro) + (v.TrcFee * currencies[trxIndex].Euro) + (v.TrxFee * currencies[trxIndex].Euro)
			grandTotalTransferedUsd += totalTransferUsd
			grandTotalTransferedYuan += totalTransferYuan
			grandTotalTransferedEuro += totalTransferEuro
			grandTotalReceivedUsd += totalReceivedUsd
			grandTotalReceivedYuan += totalReceivedYuan
			grandTotalReceivedEuro += totalReceivedEuro
			grandTotalFeeUsd += totalFeeUsd
			grandTotalFeeYuan += totalFeeYuan
			grandTotalFeeEuro += totalFeeEuro
			grandTotalTransferedBtc += v.BtcTransfer
			grandTotalTransferedEth += v.EthTransfer
			grandTotalTransferedErc += v.ErcTransfer
			grandTotalTransferedTrc += v.TrcTransfer
			grandTotalTransferedTrx += v.TrxTransfer
			grandTotalReceivedBtc += v.BtcReceived
			grandTotalReceivedEth += v.EthReceived
			grandTotalReceivedErc += v.ErcReceived
			grandTotalReceivedTrc += v.TrcReceived
			grandTotalReceivedTrx += v.TrxReceived
			grandTotalFeeBtc += v.BtcFee
			grandTotalFeeEth += v.EthFee
			grandTotalFeeErc += v.ErcFee
			grandTotalFeeTrc += v.TrcFee
			grandTotalFeeTrx += v.TrxFee
			totalUsers += newUsers
			BtcTransfer := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.BtcTransfer),
				UsdBalance:  fmt.Sprintf("%.8f", v.BtcTransfer*currencies[btcIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.BtcTransfer*currencies[btcIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.BtcTransfer*currencies[btcIndex].Euro),
			}
			EthTransfer := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.EthTransfer),
				UsdBalance:  fmt.Sprintf("%.8f", v.EthTransfer*currencies[ethIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.EthTransfer*currencies[ethIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.EthTransfer*currencies[ethIndex].Euro),
			}
			ErcTransfer := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.ErcTransfer),
				UsdBalance:  fmt.Sprintf("%.8f", v.ErcTransfer*currencies[usdtIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.ErcTransfer*currencies[usdtIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.ErcTransfer*currencies[usdtIndex].Euro),
			}
			TrcTransfer := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrcTransfer),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrcTransfer*currencies[usdtIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrcTransfer*currencies[usdtIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrcTransfer*currencies[usdtIndex].Euro),
			}
			TrxTransfer := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrxTransfer),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrxTransfer*currencies[trxIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrxTransfer*currencies[trxIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrxTransfer*currencies[trxIndex].Euro),
			}
			TransferAssets := &admin.AccountAsset{
				UsdAmount:  utils.RoundFloat(totalTransferUsd, 2),
				YuanAmount: utils.RoundFloat(totalTransferYuan, 2),
				EuroAmount: utils.RoundFloat(totalTransferEuro, 2),
			}
			transferStats := &admin.Statistics{
				AccountAssets: TransferAssets,
				Btc:           BtcTransfer,
				Eth:           EthTransfer,
				Erc:           ErcTransfer,
				Trc:           TrcTransfer,
				Trx:           TrxTransfer,
			}

			BtcReceived := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.BtcReceived),
				UsdBalance:  fmt.Sprintf("%.8f", v.BtcReceived*currencies[btcIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.BtcReceived*currencies[btcIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.BtcReceived*currencies[btcIndex].Euro),
			}
			EthReceived := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.EthReceived),
				UsdBalance:  fmt.Sprintf("%.8f", v.EthReceived*currencies[ethIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.EthReceived*currencies[ethIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.EthReceived*currencies[ethIndex].Euro),
			}
			ErcReceived := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.ErcReceived),
				UsdBalance:  fmt.Sprintf("%.8f", v.ErcReceived*currencies[usdtIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.ErcReceived*currencies[usdtIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.ErcReceived*currencies[usdtIndex].Euro),
			}
			TrcReceived := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrcReceived),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrcReceived*currencies[usdtIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrcReceived*currencies[usdtIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrcReceived*currencies[usdtIndex].Euro),
			}
			TrxReceived := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrxReceived),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrxReceived*currencies[trxIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrxReceived*currencies[trxIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrxReceived*currencies[trxIndex].Euro),
			}
			ReceivedAssets := &admin.AccountAsset{
				UsdAmount:  utils.RoundFloat(totalReceivedUsd, 2),
				YuanAmount: utils.RoundFloat(totalReceivedYuan, 2),
				EuroAmount: utils.RoundFloat(totalReceivedEuro, 2),
			}
			receivedStats := &admin.Statistics{
				AccountAssets: ReceivedAssets,
				Btc:           BtcReceived,
				Eth:           EthReceived,
				Erc:           ErcReceived,
				Trc:           TrcReceived,
				Trx:           TrxReceived,
			}
			BtcFee := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.BtcFee),
				UsdBalance:  fmt.Sprintf("%.8f", v.BtcFee*currencies[btcIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.BtcFee*currencies[btcIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.BtcFee*currencies[btcIndex].Euro),
			}
			EthFee := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.EthFee),
				UsdBalance:  fmt.Sprintf("%.8f", v.EthFee*currencies[ethIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.EthFee*currencies[ethIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.EthFee*currencies[ethIndex].Euro),
			}
			ErcFee := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.ErcFee),
				UsdBalance:  fmt.Sprintf("%.8f", v.ErcFee*currencies[ethIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.ErcFee*currencies[ethIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.ErcFee*currencies[ethIndex].Euro),
			}
			TrcFee := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrcFee),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrcFee*currencies[trxIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrcFee*currencies[trxIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrcFee*currencies[trxIndex].Euro),
			}
			TrxFee := &admin.Coin{
				Balance:     fmt.Sprintf("%.8f", v.TrxFee),
				UsdBalance:  fmt.Sprintf("%.8f", v.TrxFee*currencies[trxIndex].Usd),
				YuanBalance: fmt.Sprintf("%.8f", v.TrxFee*currencies[trxIndex].Yuan),
				EuroBalance: fmt.Sprintf("%.8f", v.TrxFee*currencies[trxIndex].Euro),
			}
			NetworkFeeAssets := &admin.AccountAsset{
				UsdAmount:  utils.RoundFloat(totalFeeUsd, 2),
				YuanAmount: utils.RoundFloat(totalFeeYuan, 2),
				EuroAmount: utils.RoundFloat(totalFeeEuro, 2),
			}
			networkFeeStats := &admin.Statistics{
				AccountAssets: NetworkFeeAssets,
				Btc:           BtcFee,
				Eth:           EthFee,
				Erc:           ErcFee,
				Trc:           TrcFee,
				Trx:           TrxFee,
			}
			operationalReport := &admin.OperationalReport{
				Date:            v.Date,
				NewUsers:        newUsers,
				TotalTransfered: transferStats,
				TotalReceived:   receivedStats,
				NetworkFee:      networkFeeStats,
			}

			resp.OperationalReports = append(resp.OperationalReports, operationalReport)
		}
		GrandTotalBtcTransfered := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalTransferedBtc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalTransferedBtc*currencies[btcIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalTransferedBtc*currencies[btcIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalTransferedBtc*currencies[btcIndex].Euro),
		}
		GrandTotalEthTransfered := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalTransferedEth),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalTransferedEth*currencies[ethIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalTransferedEth*currencies[ethIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalTransferedEth*currencies[ethIndex].Euro),
		}
		GrandTotalErcTransfered := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalTransferedErc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalTransferedErc*currencies[usdtIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalTransferedErc*currencies[usdtIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalTransferedErc*currencies[usdtIndex].Euro),
		}
		GrandTotalTrcTransfered := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalTransferedTrc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalTransferedTrc*currencies[usdtIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalTransferedTrc*currencies[usdtIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalTransferedTrc*currencies[usdtIndex].Euro),
		}
		GrandTotalTrxTransfered := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalTransferedTrx),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalTransferedTrx*currencies[trxIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalTransferedTrx*currencies[trxIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalTransferedTrx*currencies[trxIndex].Euro),
		}
		GrandTotalAssets := &admin.AccountAsset{
			UsdAmount:  utils.RoundFloat(grandTotalTransferedUsd, 2),
			YuanAmount: utils.RoundFloat(grandTotalTransferedYuan, 2),
			EuroAmount: utils.RoundFloat(grandTotalTransferedEuro, 2),
		}
		grandTotalTransfer := &admin.Statistics{
			AccountAssets: GrandTotalAssets,
			Btc:           GrandTotalBtcTransfered,
			Eth:           GrandTotalEthTransfered,
			Erc:           GrandTotalErcTransfered,
			Trc:           GrandTotalTrcTransfered,
			Trx:           GrandTotalTrxTransfered,
		}
		GrandTotalBtcReceived := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalReceivedBtc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalReceivedBtc*currencies[btcIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalReceivedBtc*currencies[btcIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalReceivedBtc*currencies[btcIndex].Euro),
		}
		GrandTotalEthReceived := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalReceivedEth),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalReceivedEth*currencies[ethIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalReceivedEth*currencies[ethIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalReceivedEth*currencies[ethIndex].Euro),
		}
		GrandTotalErcReceived := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalReceivedErc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalReceivedErc*currencies[usdtIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalReceivedErc*currencies[usdtIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalReceivedErc*currencies[usdtIndex].Euro),
		}
		GrandTotalTrcReceived := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalReceivedTrc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalReceivedTrc*currencies[usdtIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalReceivedTrc*currencies[usdtIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalReceivedTrc*currencies[usdtIndex].Euro),
		}
		GrandTotalTrxReceived := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalReceivedTrx),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalReceivedTrx*currencies[trxIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalReceivedTrx*currencies[trxIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalReceivedTrx*currencies[trxIndex].Euro),
		}
		GrandTotalReceivedAssets := &admin.AccountAsset{
			UsdAmount:  utils.RoundFloat(grandTotalReceivedUsd, 2),
			YuanAmount: utils.RoundFloat(grandTotalReceivedYuan, 2),
			EuroAmount: utils.RoundFloat(grandTotalReceivedEuro, 2),
		}
		grandTotalReceived := &admin.Statistics{
			AccountAssets: GrandTotalReceivedAssets,
			Btc:           GrandTotalBtcReceived,
			Eth:           GrandTotalEthReceived,
			Erc:           GrandTotalErcReceived,
			Trc:           GrandTotalTrcReceived,
			Trx:           GrandTotalTrxReceived,
		}

		GrandTotalBtcFee := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalFeeBtc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalFeeBtc*currencies[btcIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalFeeBtc*currencies[btcIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalFeeBtc*currencies[btcIndex].Euro),
		}
		GrandTotalEthFee := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalFeeEth),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalFeeEth*currencies[ethIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalFeeEth*currencies[ethIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalFeeEth*currencies[ethIndex].Euro),
		}
		GrandTotalErcFee := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalFeeErc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalFeeErc*currencies[ethIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalFeeErc*currencies[ethIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalFeeErc*currencies[ethIndex].Euro),
		}
		GrandTotalTrcFee := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalFeeTrc),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalFeeTrc*currencies[trxIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalFeeTrc*currencies[trxIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalFeeTrc*currencies[trxIndex].Euro),
		}
		GrandTotalTrxFee := &admin.Coin{
			Balance:     fmt.Sprintf("%.8f", grandTotalFeeTrx),
			UsdBalance:  fmt.Sprintf("%.8f", grandTotalFeeTrx*currencies[trxIndex].Usd),
			YuanBalance: fmt.Sprintf("%.8f", grandTotalFeeTrx*currencies[trxIndex].Yuan),
			EuroBalance: fmt.Sprintf("%.8f", grandTotalFeeTrx*currencies[trxIndex].Euro),
		}
		GrandFeeTotalAssets := &admin.AccountAsset{
			UsdAmount:  utils.RoundFloat(grandTotalFeeUsd, 2),
			YuanAmount: utils.RoundFloat(grandTotalFeeYuan, 2),
			EuroAmount: utils.RoundFloat(grandTotalFeeEuro, 2),
		}
		grandTotalFee := &admin.Statistics{
			AccountAssets: GrandFeeTotalAssets,
			Btc:           GrandTotalBtcFee,
			Eth:           GrandTotalEthFee,
			Erc:           GrandTotalErcFee,
			Trc:           GrandTotalTrcFee,
			Trx:           GrandTotalTrxFee,
		}

		grandTotals := &admin.OperationalReport{
			NewUsers:        totalUsers,
			TotalTransfered: grandTotalTransfer,
			TotalReceived:   grandTotalReceived,
			NetworkFee:      grandTotalFee,
		}

		wallets, err := walletdb.GetAllWallets()
		var unique []db.AccountInformation

		for _, v := range wallets {
			skip := false
			for _, u := range unique {
				if v.UUID == u.UUID {
					skip = true
					break
				}
			}
			if !skip {
				totalAssetUsd += (v.BtcBalance * currencies[btcIndex].Usd) + (v.EthBalance * currencies[ethIndex].Usd) + (v.ErcBalance * currencies[usdtIndex].Usd) + (v.TrxBalance * currencies[trxIndex].Usd) + (v.TrcBalance * currencies[usdtIndex].Usd)
				totalAssetEuro += (v.BtcBalance * currencies[btcIndex].Euro) + (v.EthBalance * currencies[ethIndex].Euro) + (v.ErcBalance * currencies[usdtIndex].Euro) + (v.TrxBalance * currencies[trxIndex].Euro) + (v.TrcBalance * currencies[usdtIndex].Euro)
				totalAssetYuan += (v.BtcBalance * currencies[btcIndex].Yuan) + (v.EthBalance * currencies[ethIndex].Yuan) + (v.ErcBalance * currencies[usdtIndex].Yuan) + (v.TrxBalance * currencies[trxIndex].Yuan) + (v.TrcBalance * currencies[usdtIndex].Yuan)
				unique = append(unique, v)
			}
		}

		totalAssets := &admin.AccountAsset{
			UsdAmount:  utils.RoundFloat(totalAssetUsd, 2),
			YuanAmount: utils.RoundFloat(totalAssetYuan, 2),
			EuroAmount: utils.RoundFloat(totalAssetEuro, 2),
		}
		resp.GrandTotals = grandTotals
		resp.TotalNum = int32(count)
		resp.TotalAssets = totalAssets
		users, dErr := walletdb.GetWalletsCount()
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAdminUsersCount", dErr.Error())
			return resp, http.WrapError(constant.ErrDB)
		}
		resp.TotalUsers = users
	}
	newUsersToday, err := walletdb.GetCreatedUsersCountPerDay(time.Now().Unix())
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetCreatedUsersCountPerDay", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	resp.NewUsersToday = newUsersToday
	return resp, nil
}
func (s *adminRPCServer) UpdateAccountBalance(c context.Context, req *admin.UpdateAccountBalanceReq) (*admin.UpdateAccountBalanceResp, error) {
	resp := &admin.UpdateAccountBalanceResp{}
	account := &db.AccountInformation{
		UUID:        req.Uuid,
		MerchantUid: req.MerchantUid,
		BtcBalance:  req.BtcBalance,
		EthBalance:  req.EthBalance,
		ErcBalance:  req.Erc20Balance,
		TrcBalance:  req.Trc20Balance,
		TrxBalance:  req.TrxBalance,
	}
	err := walletdb.UpdateAccountInformation(account)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAccountInformation", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if req.MessageID != "" {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Removing key after updating", req.MessageID)
		db2.DB.RedisDB.RemoveAddressFromList(req.MessageID)
	}
	currencies, _ := walletdb.GetCoinStatuses()
	btcIndex := constant.BTCCoin - 1
	ethIndex := constant.ETHCoin - 1
	trxIndex := constant.TRX - 1
	usdtIndex := constant.USDTERC20 - 1

	updatedAccount, er := walletdb.GetAccountInformationByMerchantUidAndUid(req.MerchantUid, req.Uuid)
	if er != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformationByMerchantUidAndUid", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	accountAddresses := &admin.AccountAddresses{
		BtcPublicAddress: updatedAccount.BtcPublicAddress,
		EthPublicAddress: updatedAccount.EthPublicAddress,
		ErcPublicAddress: updatedAccount.ErcPublicAddress,
		TrcPublicAddress: updatedAccount.TrcPublicAddress,
		TrxPublicAddress: updatedAccount.TrxPublicAddress,
	}
	lastLoginInformation := &admin.LoginInformation{
		LoginIp:       updatedAccount.LastLoginIp,
		LoginRegion:   updatedAccount.LastLoginRegion,
		LoginTerminal: updatedAccount.LastLoginTerminal,
		LoginTime:     updatedAccount.LastLoginTime,
	}
	creationLoginInformation := &admin.LoginInformation{
		LoginIp:       updatedAccount.CreationLoginIp,
		LoginRegion:   updatedAccount.CreationLoginRegion,
		LoginTerminal: updatedAccount.CreationLoginTerminal,
		LoginTime:     updatedAccount.CreationTime,
	}
	var accountAssetUsd, accountAssetYuan, accountAssetEuro float64
	accountAssetUsd = (currencies[btcIndex].Usd * updatedAccount.BtcBalance) + (currencies[ethIndex].Usd * updatedAccount.EthBalance) + (currencies[usdtIndex].Usd * updatedAccount.ErcBalance) + (currencies[usdtIndex].Usd * updatedAccount.TrcBalance) + (currencies[trxIndex].Usd * updatedAccount.TrxBalance)
	accountAssetYuan = (currencies[btcIndex].Yuan * updatedAccount.BtcBalance) + (currencies[ethIndex].Yuan * updatedAccount.EthBalance) + (currencies[usdtIndex].Yuan * updatedAccount.ErcBalance) + (currencies[usdtIndex].Yuan * updatedAccount.TrcBalance) + (currencies[trxIndex].Yuan * updatedAccount.TrxBalance)
	accountAssetEuro = (currencies[btcIndex].Euro * updatedAccount.BtcBalance) + (currencies[ethIndex].Euro * updatedAccount.EthBalance) + (currencies[usdtIndex].Euro * updatedAccount.ErcBalance) + (currencies[usdtIndex].Euro * updatedAccount.TrcBalance) + (currencies[trxIndex].Euro * updatedAccount.TrxBalance)
	accountAssets := &admin.AccountAsset{
		UsdAmount:  utils.RoundFloat(accountAssetUsd, 2),
		YuanAmount: utils.RoundFloat(accountAssetYuan, 2),
		EuroAmount: utils.RoundFloat(accountAssetEuro, 2),
	}
	Btc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", updatedAccount.BtcBalance),
		UsdBalance:  fmt.Sprintf("%.8f", currencies[btcIndex].Usd*updatedAccount.BtcBalance),
		YuanBalance: fmt.Sprintf("%.8f", currencies[btcIndex].Yuan*updatedAccount.BtcBalance),
		EuroBalance: fmt.Sprintf("%.8f", currencies[btcIndex].Euro*updatedAccount.BtcBalance),
	}
	Eth := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", updatedAccount.EthBalance),
		UsdBalance:  fmt.Sprintf("%.8f", currencies[ethIndex].Usd*updatedAccount.EthBalance),
		YuanBalance: fmt.Sprintf("%.8f", currencies[ethIndex].Yuan*updatedAccount.EthBalance),
		EuroBalance: fmt.Sprintf("%.8f", currencies[ethIndex].Euro*updatedAccount.EthBalance),
	}
	Erc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", updatedAccount.ErcBalance),
		UsdBalance:  fmt.Sprintf("%.8f", currencies[usdtIndex].Usd*updatedAccount.ErcBalance),
		YuanBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Yuan*updatedAccount.ErcBalance),
		EuroBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Euro*updatedAccount.ErcBalance),
	}
	Trc := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", updatedAccount.TrcBalance),
		UsdBalance:  fmt.Sprintf("%.8f", currencies[usdtIndex].Usd*updatedAccount.TrcBalance),
		YuanBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Yuan*updatedAccount.TrcBalance),
		EuroBalance: fmt.Sprintf("%.8f", currencies[usdtIndex].Euro*updatedAccount.TrcBalance),
	}
	Trx := &admin.Coin{
		Balance:     fmt.Sprintf("%.8f", updatedAccount.TrxBalance),
		UsdBalance:  fmt.Sprintf("%.8f", currencies[trxIndex].Usd*updatedAccount.TrxBalance),
		YuanBalance: fmt.Sprintf("%.8f", currencies[trxIndex].Yuan*updatedAccount.TrxBalance),
		EuroBalance: fmt.Sprintf("%.8f", currencies[trxIndex].Euro*updatedAccount.TrxBalance),
	}
	coinsType := []string{utils.GetCoinName(constant.BTCCoin),
		utils.GetCoinName(constant.ETHCoin),
		utils.GetCoinName(constant.USDTERC20),
		utils.GetCoinName(constant.TRX),
		utils.GetCoinName(constant.USDTTRC20)}

	accountInformation := &admin.AccountInformation{
		ID:                       updatedAccount.ID,
		Uid:                      updatedAccount.UUID,
		MerchantUid:              updatedAccount.MerchantUid,
		CoinsType:                coinsType,
		Addresses:                accountAddresses,
		AccountAssets:            accountAssets,
		Btc:                      Btc,
		Eth:                      Eth,
		Trx:                      Trx,
		Erc:                      Erc,
		Trc:                      Trc,
		AccountSource:            updatedAccount.AccountSource,
		CreationLoginInformation: creationLoginInformation,
		LastLoginInformation:     lastLoginInformation,
		TotalBalance:             0,
	}
	resp.Account = accountInformation
	return resp, nil
}
