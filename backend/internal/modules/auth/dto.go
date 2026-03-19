package auth

type WeChatLoginReq struct {
	Code          string `json:"code" binding:"required"`
	EncryptedData string `json:"encrypted_data"`
	IV            string `json:"iv"`
}

type PhoneLoginReq struct {
	Phone    string `json:"phone" binding:"required,phone"`
	Code     string `json:"code" binding:"required,len=6"`
}

type PasswordLoginReq struct {
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=6"`
}

type SendSMSReq struct {
	Phone string `json:"phone" binding:"required,phone"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileReq struct {
	Nickname  *string `json:"nickname" binding:"omitempty,max=64"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,max=512"`
	Gender    *int    `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Age       *int    `json:"age" binding:"omitempty,min=0,max=150"`
	BallAge   *int    `json:"ball_age" binding:"omitempty,min=0"`
	SelfLevel *string `json:"self_level" binding:"omitempty"`
	UTRImage  *string `json:"utr_image" binding:"omitempty,max=512"`
}

type LoginResp struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserResp `json:"user"`
}

type UserResp struct {
	ID              int64   `json:"id"`
	Phone           string  `json:"phone"`
	Nickname        string  `json:"nickname"`
	AvatarURL       string  `json:"avatar_url"`
	Gender          int     `json:"gender"`
	Age             int     `json:"age"`
	UTRRating       *string `json:"utr_rating"`
	BallAge         int     `json:"ball_age"`
	SelfLevel       string  `json:"self_level"`
	MemberLevel     int     `json:"member_level"`
	MemberLevelName string  `json:"member_level_name"`
	Role            string  `json:"role"`
}

type WalletResp struct {
	Balance        string `json:"balance"`
	FrozenAmount   string `json:"frozen_amount"`
	TotalRecharged string `json:"total_recharged"`
}

func ToUserResp(u *User) UserResp {
	phone := ""
	if u.Phone != nil {
		phone = *u.Phone
	}

	var utr *string
	if u.UTRRating != nil {
		s := u.UTRRating.String()
		utr = &s
	}

	return UserResp{
		ID:              u.ID,
		Phone:           phone,
		Nickname:        u.Nickname,
		AvatarURL:       u.AvatarURL,
		Gender:          u.Gender,
		Age:             u.Age,
		UTRRating:       utr,
		BallAge:         u.BallAge,
		SelfLevel:       u.SelfLevel,
		MemberLevel:     u.MemberLevel,
		MemberLevelName: memberLevelName(u.MemberLevel),
		Role:            u.Role,
	}
}

func memberLevelName(level int) string {
	switch level {
	case 1:
		return "银卡会员"
	case 2:
		return "金卡会员"
	case 3:
		return "钻石会员"
	default:
		return "普通会员"
	}
}
