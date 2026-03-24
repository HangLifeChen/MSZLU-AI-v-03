package auths

import (
	"common/biz"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"model"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/cache"
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
	"github.com/mszlu521/thunder/tools/jwt"
	"github.com/mszlu521/thunder/tools/randoms"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	repo  repository
	cache *cache.RedisCache
}

func (s *service) register(req RegisterReq) (*RegisterResp, error) {
	//先检查用户名邮箱是否已经注册
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u, err := s.repo.findByUsername(ctx, req.Username)
	if err != nil {
		logs.Errorf("register findByUsername error: %v", err)
		return nil, errs.DBError
	}
	if u != nil {
		return nil, biz.ErrUserNameExisted
	}
	//检查邮箱是否已经注册
	u, err = s.repo.findByEmail(ctx, req.Email)
	if err != nil {
		logs.Errorf("register findByEmail error: %v", err)
		return nil, errs.DBError
	}
	if u != nil {
		return nil, biz.ErrEmailExisted
	}
	//对密码进行加密
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Errorf("register GenerateFromPassword error: %v", err)
		return nil, biz.ErrPasswordFormat
	}
	//生成邮件用的token 用于邮件激活
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errs.DBError
	}
	token := hex.EncodeToString(tokenBytes)
	userId := uuid.New()
	//存入redis 中，用于激活邮件的时候验证
	tokenKey := fmt.Sprintf("verify_token:%s", token)
	err = s.cache.Set(tokenKey, userId.String(), 24*60*60)
	if err != nil {
		logs.Errorf("register Set error: %v", err)
		return nil, errs.DBError
	}
	//存入数据库
	user := model.User{
		Id:            userId,
		Username:      req.Username,
		Password:      string(password),
		Email:         req.Email,
		EmailVerified: false,
		LastLoginTime: time.Now(),
		CurrentPlan:   model.FreePlan,
		Status:        model.UserStatusPending,
		Avatar:        "default",
	}
	err = s.repo.transaction(ctx, func(tx *gorm.DB) error {
		//创建用户
		err := s.repo.saveUser(ctx, tx, &user)
		if err != nil {
			logs.Errorf("register saveUser error: %v", err)
			return err
		}
		//发送邮件
		err = s.sendVerifyEmail(user.Email, user.Username, token)
		if err != nil {
			logs.Errorf("register sendVerifyEmail error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errs.DBError
	}
	return &RegisterResp{
		Message: "注册成功，请前往邮箱进行验证",
	}, nil
}

func (s *service) sendVerifyEmail(email string, username string, token string) error {
	//加载邮件的配置
	emailConfig := config.GetConfig().Email
	addr := fmt.Sprintf("%s:%d", emailConfig.GetHost(), emailConfig.GetPort())
	auth := smtp.PlainAuth("", emailConfig.GetUsername(), emailConfig.GetPassword(), emailConfig.GetHost())
	to := []string{email}
	subject := "请验证您的邮箱地址"
	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", emailConfig.GetBaseURL(), token)

	// 使用HTML格式，避免被反垃圾邮件系统拦截
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #333;">邮箱验证</h2>
        <p>尊敬的 <strong>%s</strong>，</p>
        <p>感谢您注册我们的服务！</p>
        <p>请点击下方按钮验证您的邮箱地址：</p>
        <div style="margin: 30px 0;">
            <a href="%s" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">验证邮箱</a>
        </div>
        <p>如果按钮无法点击，请复制以下链接到浏览器地址栏：</p>
        <p style="word-break: break-all; color: #666;">%s</p>
        <p style="color: #999; font-size: 12px; margin-top: 30px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>`, username, verifyURL, verifyURL)

	// 构建完整的邮件头
	from := emailConfig.GetFrom()
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Date: %s\r\n\r\n%s",
		from, email, subject, time.Now().Format(time.RFC1123Z), body)

	err := smtp.SendMail(addr, auth, from, to, []byte(msg))
	return err
}

func (s *service) verifyEmail(token string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//先从redis获取userId
	tokenKey := fmt.Sprintf("verify_token:%s", token)
	userIdStr, err := s.cache.Get(tokenKey)
	if err != nil {
		logs.Errorf("verifyEmail Get error: %v", err)
		return nil, biz.ErrTokenInvalid
	}
	//这个验证邮件的时候 redis的key也需要删除
	defer s.cache.Set(tokenKey, "", 1)
	//转换成uuid
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logs.Errorf("verifyEmail Parse error: %v", err)
		return nil, biz.ErrTokenInvalid
	}
	//根据用户id查找 用户
	u, err := s.repo.findById(ctx, userId)
	if err != nil {
		logs.Errorf("verifyEmail findById error: %v", err)
		return nil, errs.DBError
	}
	if u == nil {
		return nil, biz.ErrUserNotFound
	}
	//判断用户邮箱是否已经验证
	if u.EmailVerified {
		//直接返回验证成功
		return nil, nil
	}
	//更新 用户
	u.EmailVerified = true
	u.Status = model.UserStatusNormal
	err = s.repo.transaction(ctx, func(tx *gorm.DB) error {
		return s.repo.updateUser(ctx, tx, u)
	})
	if err != nil {
		logs.Errorf("verifyEmail updateUser error: %v", err)
		return nil, errs.DBError
	}
	return nil, nil
}

func (s *service) login(req LoginReq) (*LoginResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//根据用户名或邮箱查询用户
	u, err := s.repo.findByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		logs.Errorf("login findByUsernameOrEmail error: %v", err)
		return nil, errs.DBError
	}
	if u == nil {
		return nil, biz.ErrUserNotFound
	}
	if !u.EmailVerified {
		return nil, biz.ErrEmailNotVerified
	}
	//比对密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, biz.ErrPasswordInvalid
	}
	return s.token(u)
}

func (s *service) token(u *model.User) (*LoginResp, error) {
	//生成token和refreshToken
	expire := config.GetConfig().Jwt.GetExpire()
	refreshExpire := config.GetConfig().Jwt.GetRefresh()
	token, err := jwt.GenToken(u.Id.String(), u.Username, expire)
	if err != nil {
		logs.Errorf("token GenToken error: %v", err)
		return nil, biz.ErrTokenGen
	}
	refreshToken, err := jwt.GenToken(u.Id.String(), u.Username, refreshExpire)
	if err != nil {
		logs.Errorf("token GenToken error: %v", err)
		return nil, biz.ErrTokenGen
	}
	return &LoginResp{
		Expire:        time.Now().Add(expire).UnixMilli(),
		Token:         token,
		RefreshExpire: time.Now().Add(refreshExpire).UnixMilli(),
		RefreshToken:  refreshToken,
		UserInfo: &model.UserDTO{
			Id:          u.Id,
			Username:    u.Username,
			Avatar:      u.Avatar,
			Status:      u.Status,
			CurrentPlan: u.CurrentPlan,
		},
	}, nil
}

func (s *service) refreshToken(refreshToken string) (*LoginResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//解析refreshToken
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, biz.ErrTokenInvalid
	}
	userIdStr := claims.UserId
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil, biz.ErrTokenInvalid
	}
	//根据用户id查询用户
	u, err := s.repo.findById(ctx, userId)
	if err != nil {
		logs.Errorf("refreshToken findById error: %v", err)
		return nil, errs.DBError
	}
	//重新生成token和refreshToken
	return s.token(u)
}

func (s *service) forgotPassword(forgetReq ForgetPasswordReq) (any, error) {
	//先检查邮件是否存在
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u, err := s.repo.findByEmail(ctx, forgetReq.Email)
	if err != nil {
		logs.Errorf("forgetPassword findByEmail error: %v", err)
		return nil, errs.DBError
	}
	if u == nil {
		return nil, biz.ErrUserNotFound
	}
	//生成验证码
	code, err := randoms.Gen6Code()
	if err != nil {
		logs.Errorf("forgetPassword Gen6Code error: %v", err)
		return nil, biz.ErrCodeGen
	}
	//保存验证码到redis
	codeKey := fmt.Sprintf("forget_password_code:%s", u.Email)
	err = s.cache.Set(codeKey, code, 5*60)
	if err != nil {
		logs.Errorf("forgetPassword Set error: %v", err)
		return nil, errs.DBError
	}
	//发送邮件
	err = s.sendForgetPasswordEmail(u.Email, u.Username, code)
	if err != nil {
		logs.Errorf("forgetPassword sendForgetPasswordEmail error: %v", err)
		return nil, errs.DBError
	}
	return map[string]any{
		"message": "已发送验证码，请查收邮件",
	}, nil
}

func (s *service) sendForgetPasswordEmail(email string, username string, code string) error {
	//加载邮件的配置
	emailConfig := config.GetConfig().Email
	addr := fmt.Sprintf("%s:%d", emailConfig.GetHost(), emailConfig.GetPort())
	auth := smtp.PlainAuth("", emailConfig.GetUsername(), emailConfig.GetPassword(), emailConfig.GetHost())
	to := []string{email}
	subject := "您的验证码"

	// 使用HTML格式，避免被反垃圾邮件系统拦截
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #333;">密码重置验证码</h2>
        <p>尊敬的 <strong>%s</strong>，</p>
        <p>您正在重置密码，验证码是：</p>
        <div style="margin: 30px 0; padding: 20px; background-color: #f5f5f5; text-align: center; border-radius: 5px;">
            <span style="font-size: 32px; font-weight: bold; color: #007bff; letter-spacing: 5px;">%s</span>
        </div>
        <p>验证码5分钟内有效，如非本人操作请忽略。</p>
        <p style="color: #999; font-size: 12px; margin-top: 30px;">此邮件由系统自动发送，请勿回复。</p>
    </div>
</body>
</html>`, username, code)

	// 构建完整的邮件头
	from := emailConfig.GetFrom()
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Date: %s\r\n\r\n%s",
		from, email, subject, time.Now().Format(time.RFC1123Z), body)

	err := smtp.SendMail(addr, auth, from, to, []byte(msg))
	return err
}

func (s *service) verifyCode(req VerifyCodeReq) (*VerifyCodeResp, error) {
	//先从redis获取验证码
	codeKey := fmt.Sprintf("forget_password_code:%s", req.Email)
	code, err := s.cache.Get(codeKey)
	if err != nil {
		logs.Errorf("verifyCode Get error: %v", err)
		return nil, biz.ErrCodeInvalid
	}
	//检查验证码是否正确
	if code != req.Code {
		return nil, biz.ErrCodeInvalid
	}
	//生成一个用于重置密码的临时令牌
	token, err := s.generateResetPasswordToken(req.Email)
	if err != nil {
		logs.Errorf("verifyCode generrateResetPasswordToken error: %v", err)
		return nil, biz.ErrTokenGen
	}
	//删除redis中的验证码 设置1秒过期
	defer s.cache.Set(codeKey, "", 1)
	return &VerifyCodeResp{
		Message: "验证成功",
		Token:   token,
	}, nil
}

func (s *service) generateResetPasswordToken(email string) (string, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", errs.DBError
	}
	token := hex.EncodeToString(tokenBytes)
	//放入redis中
	key := fmt.Sprintf("reset_password_token:%s", token)
	//1小时 这个重置密码后需要及时删除 或者设置的短一些 15分钟
	err := s.cache.Set(key, email, 15*60)
	if err != nil {
		logs.Errorf("generateResetPasswordToken Set error: %v", err)
		return "", errs.DBError
	}
	return token, nil
}

func (s *service) resetPassword(c context.Context, resetReq ResetPasswordReq) (any, error) {
	//验证重置密码的令牌
	tokenKey := fmt.Sprintf("reset_password_token:%s", resetReq.Token)
	email, err := s.cache.Get(tokenKey)
	if err != nil {
		return nil, biz.ErrTokenInvalid
	}
	defer s.cache.Set(tokenKey, "", 1)
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	//检查邮箱是否匹配
	if email != resetReq.Email {
		return nil, biz.ErrEmailNotMatch
	}
	//根据邮箱查询用户
	u, err := s.repo.findByEmail(ctx, resetReq.Email)
	if err != nil {
		logs.Errorf("resetPassword findByEmail error: %v", err)
		return nil, errs.DBError
	}
	if u == nil {
		return nil, biz.ErrUserNotFound
	}
	//新的密码加密
	newPassword, err := bcrypt.GenerateFromPassword([]byte(resetReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logs.Errorf("resetPassword GenerateFromPassword error: %v", err)
		return nil, biz.ErrPasswordFormat
	}
	//更新密码
	u.Password = string(newPassword)
	err = s.repo.transaction(ctx, func(tx *gorm.DB) error {
		return s.repo.updateUser(ctx, tx, u)
	})
	if err != nil {
		logs.Errorf("resetPassword updateUser error: %v", err)
		return nil, errs.DBError
	}
	return map[string]any{
		"message": "密码重置成功",
	}, nil
}

func newService() *service {
	return &service{
		repo:  newModel(database.GetPostgresDB().GormDB),
		cache: cache.NewRedisCache(),
	}
}
