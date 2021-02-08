package errmsg

type Code int

const (
	SUCCESS     = 200
	ERROR       = 500
	PARAM_ERROR = 501
	// code = 1000 用户模块错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003
	ERROR_TOKEN_NOT_EXIST  = 1004
	ERROR_TOKEN_OVERTIME   = 1005
	ERROR_TOKEN_WRONG      = 1006
	ERROR_TOKEN_ILLIGAL    = 1007
	ERROR_USER_SELECT_FAIL = 1008
	ERROR_NOT_ADMIN        = 1009

	// code = 2000 分类模块错误
	ERROR_CATEGORY_EXIST = 2001
	// code = 3000 文章模块错误
	ERROR_POST_NOT_FOUND = 3001
)

var codemsg = map[Code]string{
	ERROR_NOT_ADMIN:        "Not administartor!",
	PARAM_ERROR:            "Wrong parameters!",
	ERROR_POST_NOT_FOUND:   "Post not found!",
	ERROR_CATEGORY_EXIST:   "Category exist!",
	ERROR_USER_SELECT_FAIL: "Select users fails!",
	ERROR_TOKEN_ILLIGAL:    "Illigal token!",
	ERROR_TOKEN_WRONG:      "Wrong Token!",
	ERROR_TOKEN_OVERTIME:   "Token over time!",
	ERROR_TOKEN_NOT_EXIST:  "Token not exists!",
	ERROR_USER_NOT_EXIST:   "User not exists!",
	ERROR_PASSWORD_WRONG:   "Wrong passord!",
	ERROR_USERNAME_USED:    "Username already exists!",
	ERROR:                  "FAIL",
	SUCCESS:                "OK",
}

func GetErrMsg(code Code) string {
	return codemsg[code]
}
