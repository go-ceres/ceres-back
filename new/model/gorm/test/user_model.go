package test

import (
	"context"
	"github.com/go-ceres/go-ceres/store/gorm"
	"github.com/google/wire"
	"time"
)

var (
	UserSet = wire.NewSet(wire.Struct(new(UserModel), "*"))
)

type (
	UserModel struct {
		DB *gorm.DB
	}

	User struct {
		Id         int64     `gorm:"column:id;primaryKey;autoIncrement;not null" json:"id"`
		Password   string    `gorm:"column:password;type:varchar(50);not null;default:'';comment:用户密码" json:"password"`
		Mobile     string    `gorm:"column:mobile;uniqueIndex:mobile_index;type:varchar(11);not null;default:'';comment:手机号" json:"mobile"`
		Email      string    `gorm:"column:email;uniqueIndex:email_index;type:varchar(50);not null;default:'';comment:邮箱" json:"email"`
		Nickname   string    `gorm:"column:nickname;uniqueIndex:nickname_index;type:varchar(50);default:'';comment:用户昵称" json:"nickname"`
		Gender     string    `gorm:"column:gender;type:enum('0','1','2');default:'0';comment:男｜女｜未公开" json:"gender"`
		Autograph  string    `gorm:"column:autograph;type:varchar(255);default:'';comment:签名" json:"autograph"`
		Birthday   int64     `gorm:"column:birthday;type:bigint(13);default:0;comment:生日" json:"birthday"`
		Avatar     string    `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:用户头像" json:"avatar"`
		Status     int       `gorm:"column:status;type:tinyint(1);default:0;comment:用户类型" json:"status"`
		CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP" json:"create_time"`
		UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
	}

	Users []*User

	UserQueryParam struct {
		gorm.PaginationParam
		Ids      []int64 `json:"ids" form:"ids"`
		Mobile   string  `json:"mobile" form:"mobile"`
		Email    string  `json:"email" form:"email"`
		Nickname string  `json:"nickname" form:"nickname"`
	}

	UserQueryResult struct {
		List       Users
		PageResult *gorm.PaginationResult
	}
)

func (User) TableName() string {
	return "user"
}

func GetUserDb(ctx context.Context, def *gorm.DB) *gorm.DB {
	return gorm.GetDbWithModel(ctx, def, new(User))
}

func AutoMigrateGormUser(db *gorm.DB) error {
	return db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").AutoMigrate(new(User))
}

func NewUserModel(db *gorm.DB) UserModel {
	return UserModel{
		DB: db,
	}
}

func (m *UserModel) Create(ctx context.Context, param *User) error {
	result := GetUserDb(ctx, m.DB).Create(param)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *UserModel) Delete(ctx context.Context, Ids []int64) error {
	result := GetUserDb(ctx, m.DB).Where("id IN (?)", Ids).Delete(User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *UserModel) Update(ctx context.Context, Id int64, param User) error {
	result := GetUserDb(ctx, m.DB).Where("id = ?", Id).Updates(&param)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *UserModel) FindOne(ctx context.Context, params User) (*User, error) {
	var resp User
	db := GetUserDb(ctx, m.DB).Where(&params)
	_, err := gorm.FindOne(db, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *UserModel) Query(ctx context.Context, params UserQueryParam, opts ...gorm.QueryOptions) (*UserQueryResult, error) {
	opt := gorm.GetQueryOption(opts...)
	db := GetUserDb(ctx, m.DB)

	if v := params.Ids; len(v) > 0 {
		db.Where("id IN (?)", v)
	}

	if v := params.Mobile; v != "" {
		db = db.Where("mobile like ?", "%"+v+"%")
	}

	if v := params.Email; v != "" {
		db = db.Where("email like ?", "%"+v+"%")
	}

	if v := params.Nickname; v != "" {
		db = db.Where("nickname like ?", "%"+v+"%")
	}

	opt.OrderFields = append(opt.OrderFields, gorm.NewOrderField("id", gorm.OrderByDESC))
	db = db.Order(gorm.ParseOrder(opt.OrderFields))
	var list Users
	pr, err := gorm.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, err
	}
	qr := &UserQueryResult{
		List:       list,
		PageResult: pr,
	}
	return qr, nil
}
