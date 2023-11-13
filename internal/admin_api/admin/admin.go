package admin

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	http2 "Share-Wallet/pkg/common/http"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/common/token_verify"
	db "Share-Wallet/pkg/db/mysql"
	sql "Share-Wallet/pkg/db/mysql/mysql_model"
	"Share-Wallet/pkg/grpc-etcdv3/getcdv3"
	"Share-Wallet/pkg/proto/admin"
	"Share-Wallet/pkg/proto/eth"
	"Share-Wallet/pkg/proto/tron"
	adminStruct "Share-Wallet/pkg/struct/admin_api"
	"Share-Wallet/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	gotp "github.com/diebietse/gotp/v2"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// TestAdminApi godoc
// @Summary      Testing admin server
// @Description  Testing admin server
// @Tags         Test
// @Accept       json
// @Produce      json
// @Param        req body admin_api.TestRequest true "operationID is only for tracking"
// @Success      200  {object}  admin_api.TestResponse
// @Failure      400  {object}  admin_api.TestResponse
// @Router       /admin/test [post]
func TestAdminApi(c *gin.Context) {
	var (
		req   adminStruct.TestRequest
		resp  adminStruct.TestResponse
		reqPb admin.CommonReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "TestAdminApi!!")

	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	_, err := client.TestAdminRPC(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Name = req.Name
	secret, _ := gotp.DecodeBase32("MRAJWWHTHCTCUAXH")
	totp, _ := gotp.NewTOTP(secret)
	totpCode, _ := totp.Now()
	log.NewInfo("0", utils.GetSelfFuncName(), "THE CODE IS : ", totpCode)

	// roles, err := sql.GetAllActions()
	// // role, err := sql.GetRoleActionByActionID(1)
	// log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Role information", roles[0].ActionName, "err", err)

	http2.RespHttp200(c, constant.OK, resp)
}

// Admin Login godoc
// @Summary      Admin Login
// @Description  Admin Login
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.AdminLoginRequest true "admin_name and secret are required"
// @Success      200  {object}  admin_api.AdminLoginResponse
// @Failure      400  {object}  admin_api.AdminLoginResponse
// @Router       /admin/login [post]
func AdminLogin(c *gin.Context) {
	var (
		req   adminStruct.AdminLoginRequest
		resp  adminStruct.AdminLoginResponse
		reqPb admin.AdminLoginReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.LoginIp = c.ClientIP()
	reqPb.Secret = req.Secret
	reqPb.AdminID = req.AdminName
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.GAuthTypeToken = true
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.AdminLogin(c, &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Token = respPb.Token
	resp.GAuthEnabled = respPb.GAuthEnabled
	resp.GAuthSetupRequired = respPb.GAuthSetupRequired
	resp.GAuthSetupProvUri = respPb.GAuthSetupProvUri
	resp.User.UserName = req.AdminName
	// resp.User.Permissions = respPb.User.Permissions
	// resp.User.Role = respPb.User.Role

	http2.RespHttp200(c, constant.OK, resp)
}

// ChangePassword godoc
// @Summary      ChangePassword
// @Description  ChangePassword
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.AdminPasswordChangeRequest true "secret and new_secret are required"
// @Success      200  {object}  admin_api.AdminPasswordChangeResponse
// @Failure      400  {object}  admin_api.AdminPasswordChangeResponse
// @Router       /admin/reset-password [post]
func ChangePassword(c *gin.Context) {
	var (
		resp     adminStruct.AdminPasswordChangeResponse
		reqPb    admin.ChangeAdminUserPasswordReq
		userName string
	)
	params := adminStruct.AdminPasswordChangeRequest{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	OperationID := utils.OperationIDGenerator()
	//get the userId from middleware
	userIDInter, existed := c.Get("userID")
	if existed {
		userName = userIDInter.(string)
	}
	if userName != "" {
		reqPb.Secret = params.Secret
		reqPb.NewSecret = params.NewSecret
		// reqPb.TOTP = params.TOTP
		reqPb.UserName = userName
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
		if etcdConn == nil {
			errMsg := OperationID + "getcdv3.GetConn == nil"
			log.NewError(OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := admin.NewAdminClient(etcdConn)
		respPb, err := client.ChangeAdminUserPassword(context.Background(), &reqPb)
		if err != nil {
			log.NewError(OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
			http2.RespHttp200(c, err, nil)
			return
		}
		resp.Token = respPb.Token
		resp.PasswordUpdated = respPb.PasswordUpdated
		http2.RespHttp200(c, constant.OK, resp)
		return

	}
	c.JSON(http.StatusOK, gin.H{"code": constant.ErrTokenInvalid, "err_msg": constant.TokenInvalidMsg, "data": nil})

}

// AdminUserList godoc
// @Summary      AdminUserList
// @Description  AdminUserList
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetAdminUsersRequest true "operationID is required"
// @Success      200  {object}  admin_api.GetAdminUsersResponse
// @Failure      400  {object}  admin_api.GetAdminUsersResponse
// @Router       /admin/users [get]
func AdminUserList(c *gin.Context) {
	var (
		req   adminStruct.GetAdminUsersRequest
		reqPb admin.GetAdminUserListReq
		resp  adminStruct.GetAdminUsersResponse
	)

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "AdminUserList")
	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	reqPb.Name = req.Name
	reqPb.OrderBy = req.OrderBy
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetAdminUserList(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, resp)
		return
	}
	utils.CopyStructFields(&resp.Users, respPb.User)
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.UserNums = respPb.TotalUsers
	http2.RespHttp200(c, err, resp)
}

// AdminUserRole godoc
// @Summary      AdminUserRole
// @Description  AdminUserRole
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetAdminUserRoleRequest true "operationID is required"
// @Success      200  {object}  admin_api.GetAdminUserRoleResponse
// @Failure      400  {object}  admin_api.GetAdminUserRoleResponse
// @Router       /admin/roles [get]
func AdminUserRole(c *gin.Context) {
	var (
		req   adminStruct.GetAdminUserRoleRequest
		reqPb admin.GetAdminUserRoleReq
		resp  adminStruct.GetAdminUserRoleResponse
	)

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "AdminUserRole")
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	reqPb.Name = req.Name
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetAdminUserRole(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, resp)
		return
	}
	utils.CopyStructFields(&resp.Roles, respPb.Role)
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.RoleNums = respPb.TotalUserRole
	http2.RespHttp200(c, err, resp)
}

// AddAdminUser godoc
// @Summary      AddAdminUser
// @Description  AddAdminUser
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.PostAdminUserRequest true "user_name,secret,role,remarks,two_factor_enabled and status are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/users [post]
func AddAdminUser(c *gin.Context) {
	var (
		req   adminStruct.PostAdminUserRequest
		reqPb admin.AddAdminUserReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "AddAdminUser")
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.Name = req.UserName
	reqPb.Password = req.Secret
	reqPb.Role = req.Role
	reqPb.Status = req.Status
	reqPb.GAuthEnabled = req.TwoFactorEnabled
	reqPb.Remarks = req.Remarks
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}
	user, er := sql.GetRegAdminUsrByUID(userID)
	if er != nil {
		log.NewError(req.OperationID, "Admin user have not register", user.ID, er.Error())
		c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
		return
	}
	hasPermission := CheckPermission(c, user, "Add/Edit administrator")
	if !hasPermission {
		http2.RespHttp200(c, constant.ErrUnauthorized, nil)
		return
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	_, err := client.AddAdminUser(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, nil)
}

// AddAdminUserRole godoc
// @Summary      AddAdminUserRole
// @Description  AddAdminUserRole
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.PostAdminRoleRequest true "role_name,description,actionIDs,remarks,status and username are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/roles [post]
func AddAdminUserRole(c *gin.Context) {
	var (
		req   adminStruct.PostAdminRoleRequest
		reqPb admin.AddAdminUserRoleReq
	)

	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "AddAdminUserRole")
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.RoleName = req.RoleName
	reqPb.RoleDescription = req.Description
	reqPb.Remarks = req.Remarks
	reqPb.Status = req.Status
	reqPb.ActionIDs = req.ActionIDs
	// Todo: GetUserID from context, define middleware to inject UserID/UserName in requestContext.
	// For the time being getting userName from request directly.
	reqPb.UserName = req.UserName
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}
	user, er := sql.GetRegAdminUsrByUID(userID)
	if er != nil {
		log.NewError(req.OperationID, "Admin user have not register", user.ID, er.Error())
		c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
		return
	}
	hasPermission := CheckPermission(c, user, "Add/Edit role")
	if !hasPermission {
		http2.RespHttp200(c, constant.ErrUnauthorized, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	_, err := client.AddAdminUserRole(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, nil)

}

// DeleteUserAPI godoc
// @Summary      DeleteUserAPI
// @Description  DeleteUserAPI
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.DeleteAdminRequest true "user_name and operationID are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/user-delete [post]
func DeleteUserAPI(c *gin.Context) {
	var (
		req   adminStruct.DeleteAdminRequest
		reqPb admin.DeleteAdminReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.UserName = req.UserName
	var ok bool
	var errInfo string
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)
		reqPb.DeleteUser = user.UserName
		if err != nil {
			log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		hasPermission := CheckPermission(c, user, "Delete admin")
		if !hasPermission {
			http2.RespHttp200(c, constant.ErrUnauthorized, nil)
			return
		}
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}

	client := admin.NewAdminClient(etcdConn)
	_, err := client.DeleteAdminUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, "Deleted")
}

// UpdateAdmin godoc
// @Summary      UpdateAdmin
// @Description  UpdateAdmin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.UpdateAdminReq true "old_name,new_name,password,role,status, google_verification and operationID are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/users-update [post]
func UpdateAdmin(c *gin.Context) {
	var (
		req   adminStruct.UpdateAdminReq
		reqPb admin.UpdateAdminReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.UserName = req.UserName
	reqPb.Password = req.Password
	reqPb.RoleName = req.RoleName
	reqPb.Remarks = req.Remarks
	reqPb.Status = req.Status
	reqPb.TwoFactorEnabled = req.TwoFactorEnabled

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}

	var ok bool
	var errInfo string
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)
		reqPb.UpdateUser = user.UserName
		if err != nil {
			log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		hasPermission := CheckPermission(c, user, "Add/Edit administrator")
		if !hasPermission {
			http2.RespHttp200(c, constant.ErrUnauthorized, nil)
			return
		}
	}
	client := admin.NewAdminClient(etcdConn)
	_, err := client.UpdateAdminUser(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	http2.RespHttp200(c, constant.OK, "Updated")
}

// UpdateAdminRole godoc
// @Summary      UpdateAdminRole
// @Description  UpdateAdminRole
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.UpdateAdminRoleRequest true "role_name,description,actionIDs,remarks and operationID are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/roles-update [post]
func UpdateAdminRole(c *gin.Context) {
	var (
		req   adminStruct.UpdateAdminRoleRequest
		reqPb admin.UpdateAdminRoleRequest
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.RoleName = req.RoleName
	reqPb.Description = req.Description
	reqPb.ActionIDs = req.ActionIDs
	reqPb.Remarks = req.Remarks

	var ok bool
	var errInfo string
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)
		reqPb.UpdateUser = user.UserName
		if err != nil {
			log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		hasPermission := CheckPermission(c, user, "Add/Edit role")
		if !hasPermission {
			http2.RespHttp200(c, constant.ErrUnauthorized, nil)
			return
		}
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	_, err := client.UpdateAdminRole(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}

	http2.RespHttp200(c, constant.OK, "Updated")
}

// DeleteRole godoc
// @Summary      DeleteRole
// @Description  DeleteRole
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.DeleteRoleRequest true "role_name and operationID are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/role-delete [post]
func DeleteRole(c *gin.Context) {
	var (
		req   adminStruct.DeleteRoleRequest
		reqPb admin.DeleteAdminRoleRequest
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.RoleName = req.RoleName
	var ok bool
	var errInfo string
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)
		reqPb.DeleteUser = user.UserName
		if err != nil {
			log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		//Check if the user has permission to delete
		hasPermission := CheckPermission(c, user, "Delete role")
		if !hasPermission {
			http2.RespHttp200(c, constant.ErrUnauthorized, nil)
			return
		}
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}

	client := admin.NewAdminClient(etcdConn)
	_, err := client.DeleteRole(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, "Deleted")
}

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

// VerifyTOTPAdminUser godoc
// @Summary      VerifyTOTPAdminUser
// @Description  VerifyTOTPAdminUser
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.ParamsTOTPVerify true "totp and operationID are required"
// @Success      200  {object}  admin_api.AdminLoginResponse
// @Failure      400  {object}  admin_api.AdminLoginResponse
// @Router       /admin/admin-verify-totp [post]
func VerifyTOTPAdminUser(c *gin.Context) {
	var (
		resp  adminStruct.AdminLoginResponse
		reqPb admin.AdminLoginReq
	)
	params := adminStruct.ParamsTOTPVerify{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	var ok bool
	var errInfo string
	// userID, ok := c.Get("userID")
	ok, userID, errInfo := token_verify.GetAdminUserIDFromToken(c.Request.Header.Get("token"), params.OperationID, true)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)

		if err != nil {
			log.NewError(params.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		totp := genrateTOTPForNow(*user)

		if params.TOTP == totp {
			reqPb.Secret = user.Password
			reqPb.SecretHashd = true
			reqPb.AdminID = user.UserName
			reqPb.OperationID = utils.OperationIDGenerator()
			reqPb.GAuthTypeToken = false
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
			if etcdConn == nil {
				errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
				log.NewError(reqPb.OperationID, errMsg)
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
				return
			}

			client := admin.NewAdminClient(etcdConn)
			respPb, err := client.AdminLoginV2(context.Background(), &reqPb)
			if err != nil {
				log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
				http2.RespHttp200(c, err, nil)
				return
			}

			resp.Token = respPb.Token
			http2.RespHttp200(c, constant.OK, resp)
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": constant.PasswordErr, "err_msg": "TOTP is not correct", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": constant.ErrTokenInvalid, "err_msg": constant.TokenInvalidMsg, "data": nil})
}

// AdminUserActions godoc
// @Summary      AdminUserActions
// @Description  AdminUserActions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  admin_api.GetAdminUserActionsResponse
// @Failure      400  {object}  admin_api.GetAdminUserActionsResponse
// @Router       /admin/actions [get]
func AdminUserActions(c *gin.Context) {
	var (
		req   adminStruct.GetAdminUserActionsRequest
		reqPb admin.GetAdminActionsRequest
		// resp  adminStruct.GetAdminUserActionsResponse
	)

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "AdminUserAction")
	reqPb.OperationID = utils.OperationIDGenerator()

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetAdminRoleActions(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	http2.RespHttp200(c, err, respPb)
}

// GetAdminUser godoc
// @Summary      GetAdminUser
// @Description  GetAdminUser
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.GetAdminUserRequest true "operationID is required"
// @Success      200  {object}  admin_api.GetAdminUserResponse
// @Failure      400  {object}  admin_api.GetAdminUserResponse
// @Router       /admin/user-info [get]
func GetAdminUser(c *gin.Context) {
	var (
		req   adminStruct.GetAdminUserRequest
		reqPb admin.GetAdminUserRequest
		resp  adminStruct.GetAdminUserResponse
	)

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetAdminUser")
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}
	user, err := sql.GetRegAdminUsrByUID(userID)
	if err != nil {
		log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
		c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.UserName = user.UserName

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetAdminUser(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	resp.UserName = respPb.UserName
	resp.RoleName = respPb.RoleName
	resp.Permissions = respPb.Permissions
	http2.RespHttp200(c, err, resp)
}

func CheckPermission(c *gin.Context, user *db.AdminUser, permission string) bool {
	if user.RoleId != 0 {
		role, err := sql.GetRegRolesByRoleID(user.RoleId)
		if err != nil {
			log.NewError("0", "Role was not found ", role.ID, err.Error())
		}
		hasPermission := sql.CheckRolePermission(role.RoleName, permission)
		if !hasPermission {
			log.NewInfo("0", "Doesnt have permission ")
			return false
		}
		return true
	}
	return false
}

// GetAccountInformation godoc
// @Summary      GetAccountInformation
// @Description  GetAccountInformation
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetAccountInformationRequest true "role_name,description,actionIDs,remarks,status and username are required"
// @Success      200  {object}  admin_api.GetAccountInformationResponse
// @Failure      400  {object}  admin_api.GetAccountInformationResponse
// @Router       /admin/account-management [get]
func GetAccountInformation(c *gin.Context) {
	var (
		req   adminStruct.GetAccountInformationRequest
		reqPb admin.GetAccountInformationReq
		resp  adminStruct.GetAccountInformationResponse
	)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetAccountInformation")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	if req.Sort != "" {
		reqPb.Sort = req.Sort
	} else {
		reqPb.Sort = "desc"
	}
	if req.OrderBy != "" {
		reqPb.OrderBy = req.OrderBy
	} else {
		reqPb.OrderBy = "creation_time"
	}
	if req.From != "" {
		reqPb.From = req.From
		reqPb.To = req.To
	}
	if req.MerchantUid != "" {
		reqPb.MerchantUid = req.MerchantUid
	}
	if req.Uid != "" {
		reqPb.Uid = req.Uid
	}
	if req.AccountAddress != "" {
		reqPb.AccountAddress = req.AccountAddress
	}
	if req.CoinsType != "" {
		reqPb.CoinsType = req.CoinsType
	}
	if req.AccountSource != "" {
		reqPb.AccountSource = req.AccountSource
	}

	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	j, _ := json.Marshal(req)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "request ", string(j))
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetAccountInformation(context.Background(), &reqPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	if len(respPb.Account) > 0 {
		utils.CopyStructFields(&resp.Accounts, respPb.Account)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}

	resp.TotalNum = respPb.TotalAccounts
	resp.TotalAssets.UsdAmount = respPb.TotalAssets.UsdAmount
	resp.TotalAssets.YuanAmount = respPb.TotalAssets.YuanAmount
	resp.TotalAssets.EuroAmount = respPb.TotalAssets.EuroAmount
	utils.CopyStructFields(&resp.BtcTotal, respPb.BtcTotal)
	utils.CopyStructFields(&resp.EthTotal, respPb.EthTotal)
	utils.CopyStructFields(&resp.ErcTotal, respPb.ErcTotal)
	utils.CopyStructFields(&resp.TrcTotal, respPb.TrcTotal)
	utils.CopyStructFields(&resp.TrxTotal, respPb.TrxTotal)
	http2.RespHttp200(c, err, resp)
}

// GetFundsLog godoc
// @Summary      GetFundsLog
// @Description  GetFundsLog
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetFundsLogRequest true "coins type is required;transaction_type(transfer/received)"
// @Success      200  {object}  admin_api.GetFundsLogResponse
// @Failure      400  {object}  admin_api.GetFundsLogResponse
// @Router       /admin/funds-log [get]
func GetFundsLog(c *gin.Context) {
	var (
		req   adminStruct.GetFundsLogRequest
		reqPb admin.GetFundsLogReq
		resp  adminStruct.GetFundsLogResponse
	)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	if req.From != "" {
		reqPb.From = req.From
		reqPb.To = req.To
	}
	if req.TransactionType != "" {
		reqPb.TransactionType = req.TransactionType
	}
	if req.UserAddress != "" {
		reqPb.UserAddress = req.UserAddress
	}
	if req.OppositeAddress != "" {
		reqPb.OppositeAddress = req.OppositeAddress
	}
	if req.CoinsType != "" {
		reqPb.CoinsType = req.CoinsType
	}
	if req.State != "" {
		reqPb.State = req.State
	}
	if req.Txid != "" {
		reqPb.Txid = req.Txid
	}
	if req.Uid != "" {
		reqPb.Uid = req.Uid
	}
	if req.MerchantUid != "" {
		reqPb.MerchantUid = req.MerchantUid

	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetFundsLog(context.Background(), &reqPb)
	log.NewError(req.OperationID, respPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	if len(respPb.FundLog) > 0 {
		utils.CopyStructFields(&resp.Funds, respPb.FundLog)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TotalNum = respPb.TotalFundLogs
	f := excelize.NewFile()
	index := f.NewSheet("FundLog")
	f.SetActiveSheet(index)

	for idx, row := range resp.Funds {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "EXCEL", err.Error())
		}
		f.SetSheetRow("FundLog", cell, &row)
		// log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Row values are ", row, err)
		// log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Cell values are ", cell)
	}
	if err := f.SaveAs("FundLogs.xlsx"); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "Sheet saving failed")
	}

	// Set the headers necessary to get browsers to interpret the downloadable file
	// c.Writer.Header().Set("Content-Type", "application/octet-stream")
	// c.Writer.Header().Set("Content-Disposition", "attachment; filename=FundLogs.xlsx")
	// c.Writer.Header().Set("File-Name", "FundLogs.xlsx")
	// c.Writer.Header().Set("Content-Transfer-Encoding", "binary")
	// c.Writer.Header().Set("Expires", "0")
	// fmt.Fprint(c.Writer, nil)
	http2.RespHttp200(c, err, resp)

}

// GetReceiveDetails godoc
// @Summary      GetReceiveDetails
// @Description  GetReceiveDetails
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetReceiveDetailsRequest true "role_name,description,actionIDs,remarks,status and username are required"
// @Success      200  {object}  admin_api.GetRecieveDetailsResponse
// @Failure      400  {object}  admin_api.GetRecieveDetailsResponse
// @Router       /admin/receive-details [get]
func GetReceiveDetails(c *gin.Context) {
	var (
		req   adminStruct.GetReceiveDetailsRequest
		reqPb admin.GetReceiveDetailsReq
		resp  adminStruct.GetRecieveDetailsResponse
	)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	if req.From != "" {
		reqPb.From = req.From
		reqPb.To = req.To
	}
	if req.Uid != "" {
		reqPb.Uid = req.Uid
	}
	if req.MerchantUid != "" {
		reqPb.MerchantUid = req.MerchantUid

	}
	if req.ReceivingAddress != "" {
		reqPb.ReceivingAddress = req.ReceivingAddress
	}
	if req.DepositAddress != "" {
		reqPb.DepositAddress = req.DepositAddress
	}
	if req.CoinsType != "" {
		reqPb.CoinsType = req.CoinsType
	}
	if req.Txid != "" {
		reqPb.Txid = req.Txid
	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetReceiveDetails(context.Background(), &reqPb)
	log.NewError(req.OperationID, respPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	if len(respPb.ReceiveDetail) > 0 {
		utils.CopyStructFields(&resp.ReceiveDetails, respPb.ReceiveDetail)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TotalNum = respPb.TotalDetails
	resp.TotalAmountReceivedUsd = respPb.TotalAmountReceivedUsd
	resp.TotalAmountReceivedYuan = respPb.TotalAmountReceivedYuan
	resp.TotalAmountReceivedEuro = respPb.TotalAmountReceivedEuro
	resp.GrandTotalUsd = respPb.GrandTotalUsd
	resp.GrandTotalYuan = respPb.GrandTotalYuan
	resp.GrandTotalEuro = respPb.GrandTotalEuro
	http2.RespHttp200(c, err, resp)
}

// GetTransferDetails godoc
// @Summary      GetTransferDetails
// @Description  GetTransferDetails
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetTransferDetailsRequest true "username are required"
// @Success      200  {object}  admin_api.GetTransferDetailsResponse
// @Failure      400  {object}  admin_api.GetTransferDetailsResponse
// @Router       /admin/transfer-details [get]
func GetTransferDetails(c *gin.Context) {
	var (
		req   adminStruct.GetTransferDetailsRequest
		reqPb admin.GetTransferDetailsReq
		resp  adminStruct.GetTransferDetailsResponse
	)

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetFundsLog")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constant.FormattingError, "err_msg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	if req.OrderBy != "" {
		reqPb.OrderBy = req.OrderBy
	} else {
		reqPb.OrderBy = "creation_time"
	}
	if req.From != "" {
		reqPb.From = req.From
		reqPb.To = req.To
	}
	if req.MerchantUid != "" {
		reqPb.MerchantUid = req.MerchantUid

	}
	if req.CoinsType != "" {
		reqPb.CoinsType = req.CoinsType
	}
	if req.State != "" {
		reqPb.State = req.State
	}
	if req.TransferAddress != "" {
		reqPb.TransferAddress = req.TransferAddress
	}
	if req.ReceivingAddress != "" {
		reqPb.ReceivingAddress = req.ReceivingAddress
	}
	if req.Txid != "" {
		reqPb.Txid = req.Txid
	}
	if req.Uid != "" {
		reqPb.Uid = req.Uid
	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetTransferDetails(context.Background(), &reqPb)
	log.NewError(req.OperationID, respPb)
	if err != nil {
		http2.RespHttp200(c, err, "error")
		return
	}
	if len(respPb.TransferDetail) > 0 {
		utils.CopyStructFields(&resp.TransferDetails, respPb.TransferDetail)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TotalNum = respPb.TotalDetails
	resp.TotalAmountTransferedUsd = respPb.TotalAmountTransferedUsd
	resp.TotalAmountTransferedYuan = respPb.TotalAmountTransferedYuan
	resp.TotalAmountTransferedEuro = respPb.TotalAmountTransferedEuro
	resp.TotalFeeAmountUsd = respPb.TotalFeeAmountUsd
	resp.TotalFeeAmountYuan = respPb.TotalFeeAmountYuan
	resp.TotalFeeAmountEuro = respPb.TotalFeeAmountEuro
	resp.GrandTotalUsd = respPb.GrandTotalUsd
	resp.GrandTotalYuan = respPb.GrandTotalYuan
	resp.GrandTotalEuro = respPb.GrandTotalEuro
	resp.TotalTransferUsd = respPb.TotalTransferUsd
	resp.TotalTransferYuan = respPb.TotalTransferYuan
	resp.TotalTransferEuro = respPb.TotalTransferEuro
	http2.RespHttp200(c, err, resp)
}

// ResetGoogleKey godoc
// @Summary      ResetGoogleKey
// @Description  ResetGoogleKey
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.ResetGoogleKeyRequest true "username are required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/reset-google-key [post]
func ResetGoogleKey(c *gin.Context) {
	var (
		req   adminStruct.ResetGoogleKeyRequest
		reqPb admin.ResetGoogleKeyReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.UserName = req.UserName

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}

	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo //+ " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		return
	}

	if userID != "" {
		var user *db.AdminUser
		var err error
		userID := userID
		user, err = sql.GetRegAdminUsrByUID(userID)
		reqPb.UpdateUser = user.UserName
		if err != nil {
			log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
			c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
			return
		}
		hasPermission := CheckPermission(c, user, "Set/Reset Google Verification Code")
		if !hasPermission {
			http2.RespHttp200(c, constant.ErrUnauthorized, nil)
			return
		}
	}
	client := admin.NewAdminClient(etcdConn)
	resp, err := client.ResetGoogleKey(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, err, resp)
}

// GetRoleActions godoc
// @Summary      GetRoleActions
// @Description  GetRoleActions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetRoleActionsRequest true "rolename are required"
// @Success      200  {object}  admin_api.GetRoleActionsResponse
// @Failure      400  {object}  admin_api.GetRoleActionsResponse
// @Router       /admin/role-actions [get]
func GetRoleActions(c *gin.Context) {
	var (
		req   adminStruct.GetRoleActionsRequest
		reqPb admin.GetRoleActionsReq
		resp  adminStruct.GetRoleActionsResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	reqPb.RoleName = req.RoleName
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetRoleActions(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	resp.Actions = respPb.Actions
	http2.RespHttp200(c, constant.OK, resp)
}

// GetCurrencies godoc
// @Summary      GetCurrencies
// @Description  GetCurrencies
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetCurrenciesRequest true "opid is required"
// @Success      200  {object}  admin_api.GetCurrenciesResponse
// @Failure      400  {object}  admin_api.GetCurrenciesResponse
// @Router       /admin/currencies [get]
func GetCurrencies(c *gin.Context) {
	var (
		req   adminStruct.GetCurrenciesRequest
		reqPb admin.GetCurrenciesReq
		resp  adminStruct.GetCurrenciesResponse
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetCurrencies(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	if len(respPb.Currencies) > 0 {
		utils.CopyStructFields(&resp.Currencies, respPb.Currencies)
	}
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TotalNum = respPb.TotalCurrencies
	http2.RespHttp200(c, constant.OK, resp)
}

// UpdateCurrency godoc
// @Summary      UpdateCurrency
// @Description  UpdateCurrency
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.UpdateCurrencyRequest true "currency_id is required"
// @Success      200  {string}  string    "ok"
// @Failure      400  {string}  string    "error"
// @Router       /admin/update-currency [post]
func UpdateCurrency(c *gin.Context) {
	var (
		req   adminStruct.UpdateCurrencyRequest
		reqPb admin.UpdateCurrencyReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo
		log.NewError(req.OperationID, errMsg)
		return
	}
	user, err := sql.GetRegAdminUsrByUID(userID)
	if err != nil {
		log.NewError(req.OperationID, "Admin user have not register", userID, err.Error())
		c.JSON(http.StatusOK, gin.H{"code": constant.NotRegistered, "err_msg": "No user found! Kindly Register"})
		return
	}
	reqPb.UpdateUser = user.UserName
	reqPb.CurrencyId = req.CurrencyId
	reqPb.State = req.State
	client := admin.NewAdminClient(etcdConn)
	_, error := client.UpdateCurrency(context.Background(), &reqPb)
	if error != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	http2.RespHttp200(c, constant.OK, "Updated")
}

// GetOperationalReport godoc
// @Summary      GetOperationalReport
// @Description  GetOperationalReport
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.GetOperationalReportRequest true "currency_id is required"
// @Success      200  {object}  admin_api.GetOperationalReportResponse
// @Failure      400  {object}  admin_api.GetOperationalReportResponse
// @Router       /admin/operational-report [get]
func GetOperationalReport(c *gin.Context) {
	var (
		req   adminStruct.GetOperationalReportRequest
		reqPb admin.GetOperationalReportReq
		resp  adminStruct.GetOperationalReportResponse
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	if req.Page == 0 {
		req.Page = constant.DefaultPageNumber
	}
	if req.PageSize == 0 && req.Page != -1 {
		req.PageSize = constant.DefaultPageSize
	}
	reqPb.Pagination = &admin.RequestPagination{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	reqPb.From = req.From
	reqPb.To = req.To

	client := admin.NewAdminClient(etcdConn)
	respPb, err := client.GetOperationalReport(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		http2.RespHttp200(c, err, nil)
		return
	}
	if len(respPb.OperationalReports) > 0 {
		utils.CopyStructFields(&resp.OperationalReports, respPb.OperationalReports)
	}
	utils.CopyStructFields(&resp.GrandTotals, respPb.GrandTotals)
	utils.CopyStructFields(&resp.TotalAssets, respPb.TotalAssets)
	if req.Page == -1 {
		resp.Page = constant.DefaultPageNumber
		resp.PageSize = 0
	} else {
		resp.Page = int(req.Page)
		resp.PageSize = int(req.PageSize)
	}
	resp.TotalNum = respPb.TotalNum
	resp.TotalUsers = respPb.TotalUsers
	resp.NewUsersToday = respPb.NewUsersToday
	http2.RespHttp200(c, constant.OK, resp)
}

// ConfirmTransaction godoc
// @Summary      ConfirmTransaction
// @Description  ConfirmTransaction
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req query admin_api.ConfirmTransactionRequest true "tx_hash_id is required"
// @Success      200  {object}  admin_api.ConfirmTransactionResponse
// @Failure      400  {object}  admin_api.ConfirmTransactionResponse
// @Router       /admin/confirm_tx [post]
func ConfirmTransaction(c *gin.Context) {
	var (
		req  adminStruct.ConfirmTransactionRequest
		resp adminStruct.ConfirmTransactionResponse
	)
	if err := c.ShouldBind(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.CoinsType, req.TxHashId)
	if req.CoinsType == 0 || req.TxHashId == "" {
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	switch req.CoinsType {
	case constant.BTCCoin:
	case constant.ETHCoin, constant.USDTERC20:
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto eth rpc!")
		var reqPb eth.GetEthConfirmationReq
		reqPb.OperationID = req.OperationID
		reqPb.CoinType = uint32(req.CoinsType)
		reqPb.TransactionHash = req.TxHashId
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := eth.NewEthClient(etcdConn)
		_, err := client.GetConfirmationRPC(c, &reqPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

	case constant.TRX, constant.USDTTRC20:
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto tron rpc!")
		var reqPb tron.GetTronConfirmationReq
		reqPb.OperationID = req.OperationID
		reqPb.CoinType = uint32(req.CoinsType)
		reqPb.TransactionHash = req.TxHashId
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := tron.NewTronClient(etcdConn)
		_, err := client.GetConfirmationRPC(c, &reqPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

	}

	fundLog, err := sql.GetFundLogByCoinAndHash(int(req.CoinsType), req.TxHashId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFundLogByCoinAndHash error", err.Error())
		http2.RespHttp200(c, constant.ErrDB, nil)
		return
	}
	if fundLog.Txid == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "fundLog.Txid is nil", err.Error())
		http2.RespHttp200(c, constant.ErrDB, nil)
		return
	}

	log.Info(req.OperationID, utils.GetSelfFuncName(), "fundLog", fundLog.Txid, fundLog.State, fundLog.ConfirmationTime)
	var fundLogResponse adminStruct.FundsLog
	err = utils.CopyStructFields(&fundLogResponse, fundLog)
	if err != nil {
		http2.RespHttp200(c, constant.NewErrInfo(202, err.Error()), nil)
		return
	}
	fundLogResponse.State = utils.GetFundLogStateToString(fundLog.State)

	resp.TransferDetail = &fundLogResponse
	http2.RespHttp200(c, constant.OK, resp)
}

// UpdateAccountBalance godoc
// @Summary      UpdateAccountBalance
// @Description  UpdateAccountBalance
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        req body admin_api.UpdateAccountBalanceRequest true "merchant_uid and uuid are required"
// @Success      200  {object}  admin_api.AccountInformation
// @Failure      400  {object}  admin_api.AccountInformation
// @Router       /admin/update_account_balance [post]
func UpdateAccountBalance(c *gin.Context) {
	var (
		req         adminStruct.UpdateAccountBalanceRequest
		reqETHPb    eth.GetBalanceReq
		reqTronPb   tron.GetBalanceReq
		updateReqPb admin.UpdateAccountBalanceReq
		resp        adminStruct.AccountInformation
		// btcBalance float64
		ethBalance float64
		ercBalance float64
		trxBalance float64
		trcBalance float64
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrArgs, nil)
		return
	}

	account, err := sql.GetAccountInformationByMerchantUidAndUid(req.MerchantUid, req.Uuid)
	if err != nil {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), err.Error())
		http2.RespHttp200(c, constant.ErrDB, nil)
		return
	}
	if account.BtcPublicAddress != "" {
		//get btc balance
	}
	if account.EthPublicAddress != "" {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto eth rpc!")
		reqETHPb.CoinType = uint32(constant.ETHCoin)
		reqETHPb.Address = account.EthPublicAddress
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqETHPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := eth.NewEthClient(etcdConn)
		ethResp, err := client.GetEthBalanceRPC(c, &reqETHPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}
		balanceStr := ethResp.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		balanceDecimal, err := decimal.NewFromString(balanceStr)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

		bal := utils.Wei2Eth_str(balanceDecimal.BigInt())
		balFloat, err := strconv.ParseFloat(bal, 64)
		if err != nil {
			return
		}
		ethBalance = utils.RoundFloat(balFloat, 8)
	}
	if account.ErcPublicAddress != "" {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto eth rpc!")
		reqETHPb.CoinType = uint32(constant.USDTERC20)
		reqETHPb.Address = account.ErcPublicAddress
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.EthRPC, reqETHPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := eth.NewEthClient(etcdConn)
		ercResp, err := client.GetEthBalanceRPC(c, &reqETHPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}
		balanceStr := ercResp.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		balanceDecimal, err := decimal.NewFromString(balanceStr)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

		tmpBalance, _ := balanceDecimal.Float64()
		tmpBalance = (tmpBalance / 1000000)
		usdtBalance, _ := utils.FormatFloat(tmpBalance, 6)
		ercBalance = usdtBalance

	}
	if account.TrxPublicAddress != "" {
		//get trx balance
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto tron rpc!")
		reqTronPb.CoinType = uint32(constant.TRX)
		reqTronPb.Address = account.TrxPublicAddress
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqTronPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := tron.NewTronClient(etcdConn)
		tronResp, err := client.GetTronBalanceRPC(c, &reqTronPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}
		balanceStr := tronResp.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		balanceDecimal, err := decimal.NewFromString(balanceStr)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

		bal := utils.SunToTrx(balanceDecimal.BigInt())
		balFloat, _ := bal.Float64()
		trxBalance = utils.RoundFloat(balFloat, 6)
	}
	if account.TrcPublicAddress != "" {
		//get trc balance
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "goto tron rpc!")
		reqTronPb.CoinType = uint32(constant.USDTTRC20)
		reqTronPb.Address = account.TrcPublicAddress
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.TronRPC, reqTronPb.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
			return
		}
		client := tron.NewTronClient(etcdConn)
		tronResp, err := client.GetTronBalanceRPC(c, &reqTronPb)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}
		balanceStr := tronResp.Balance
		if balanceStr == "" {
			balanceStr = "0"
		}
		balanceDecimal, err := decimal.NewFromString(balanceStr)
		if err != nil {
			log.NewError(req.OperationID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
			return
		}

		tmpBalance, _ := balanceDecimal.Float64()
		tmpBalance = (tmpBalance / 1000000)
		usdtBalance, _ := utils.FormatFloat(tmpBalance, 6)
		trcBalance = usdtBalance
	}

	//update balances
	updateReqPb.OperationID = req.OperationID
	updateReqPb.MerchantUid = req.MerchantUid
	updateReqPb.Uuid = req.Uuid
	updateReqPb.EthBalance = ethBalance
	updateReqPb.Erc20Balance = ercBalance
	updateReqPb.TrxBalance = trxBalance
	updateReqPb.Trc20Balance = trcBalance
	log.NewInfo(req.OperationID, "Balances: ", ethBalance, ercBalance, trxBalance, trcBalance)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.AdminRPC, updateReqPb.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "err_msg": errMsg})
		return
	}
	client := admin.NewAdminClient(etcdConn)
	respPb, er := client.UpdateAccountBalance(c, &updateReqPb)
	if er != nil {
		log.NewError(req.OperationID, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 501, "err_msg": err.Error()})
		return
	}
	err = utils.CopyStructFields(&resp, respPb.Account)
	if err != nil {
		log.NewError(req.OperationID, "Failed to copy struct fields", err.Error())
		http2.RespHttp200(c, constant.NewErrInfo(202, err.Error()), nil)
		return
	}
	http2.RespHttp200(c, constant.OK, resp)
}
