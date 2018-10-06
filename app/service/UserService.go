package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"errors"

	"dandelion/app/play"
	"math"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type UserService struct {
	dao.BaseDao
	Configuration ConfigurationService
	GiveVoucher   GiveVoucherService
	CardItem      CardItemService
	Organization OrganizationService
}

func (service UserService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := dao.Orm()

	type Result struct {
		TotalCount  uint64 `gorm:"column:TotalCount"`
		OnlineCount int
	}

	var result Result

	Orm.Table("User").Select("COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Find(&result)
	//fmt.Println(result)
	result.OnlineCount = len(gweb.Sessions.Data)
	return result
}
func (service UserService) AddUserBlockAmount(Orm *gorm.DB,UserID uint64,Menoy int64) error  {

	var user dao.User
	err:=service.Get(Orm,UserID, &user)
	if err != nil {
		return err
	}

	tm:=int64(user.BlockAmount)+Menoy
	if tm<0{
		return errors.New("冻结金额不足，无法扣款")
	}

	err = service.ChangeMap(Orm, UserID, &dao.User{}, map[string]interface{}{"BlockAmount": tm})
	return err
}
func (service UserService) FirstSettlementUserBrokerage(Orm *gorm.DB, Brokerage uint64, orders dao.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=dao.Orm()

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	var user dao.User
	service.Get(Orm, orders.UserID, &user)


	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves:=[]uint64{leve1,leve2,leve3,leve4,leve5,leve6}

	var OutBrokerageMoney int64=0
	for _,value:=range leves{
		if value<=0{
			break
		}
		var _user dao.User
		service.Get(Orm, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.AddUserBlockAmount(Orm,_user.ID,  leveMenoy)
		if err != nil {
			return err
		}
		OutBrokerageMoney=OutBrokerageMoney+leveMenoy
		user=_user
	}

	//err=service.Organization.AddOrganizationBlockAmount(Orm,orders.OID,OutBrokerageMoney)


	/*if leve1 > 0 {

		var usera dao.User
		service.Get(Orm, user.SuperiorID, &usera)
		if usera.ID <= 0 {
			return nil
		}
		leve1Menoy := int64(math.Floor(float64(leve1)/float64(100)*float64(Brokerage) + 0.5))
		err = service.AddUserBlockAmount(Orm,usera.ID,  leve1Menoy)
		if err != nil {
			return err
		}

		if leve2 > 0 {

			var userb dao.User
			service.Get(Orm, usera.SuperiorID, &userb)
			if userb.ID <= 0 {
				return nil
			}
			leve2Menoy := int64(math.Floor(float64(leve2)/float64(100)*float64(Brokerage) + 0.5))
			err = service.AddUserBlockAmount(Orm,userb.ID,  leve2Menoy)
			if err != nil {
				return err
			}

			if leve3 > 0 {

				var userc dao.User
				service.Get(Orm, userb.SuperiorID, &userc)
				if userc.ID <= 0 {
					return nil
				}
				leve3Menoy := int64(math.Floor(float64(leve3)/float64(100)*float64(Brokerage) + 0.5))
				err = service.AddUserBlockAmount(Orm,userc.ID, leve3Menoy)
				if err != nil {
					return err
				}

				if leve4 > 0 {

					var userd dao.User
					service.Get(Orm, userc.SuperiorID, &userd)
					if userd.ID <= 0 {
						return nil
					}
					leve4Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
					err = service.AddUserBlockAmount(Orm,userd.ID,leve4Menoy)
					if err != nil {
						return err
					}

					if leve5 > 0 {

						var usere dao.User
						service.Get(Orm, userd.SuperiorID, &usere)
						if usere.ID <= 0 {
							return nil
						}
						leve5Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
						err = service.AddUserBlockAmount(Orm,usere.ID, leve5Menoy)
						if err != nil {
							return err
						}

						if leve6 > 0 {

							var userf dao.User
							service.Get(Orm, usere.SuperiorID, &userf)
							if userf.ID <= 0 {
								return nil
							}
							leve6Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
							err = service.AddUserBlockAmount(Orm,userf.ID, leve6Menoy)
							if err != nil {
								return err
							}

						}

					}

				}

			}

		}

	}*/

	return err
}

//结算佣金，结算积分，结算成长值，是否送福利卷
func (service UserService) SettlementUser(Orm *gorm.DB, Brokerage uint64, orders dao.Orders) error {
	var err error
	//用户自己。下单者

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	var user dao.User
	service.Get(Orm, orders.UserID, &user)

	//fmt.Println(user.Name)

	Journal := JournalService{}
	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves:=[]uint64{leve1,leve2,leve3,leve4,leve5,leve6}

	GrowValue, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_ScoreConvertGrowValue).V, 10, 64)

	user.Score = user.Score + uint64(math.Floor(float64(orders.PayMoney)/100+0.5))
	user.Growth = user.Growth + uint64(uint64(math.Floor(float64(orders.PayMoney)/100+0.5))*GrowValue)
	err = service.ChangeModel(Orm, user.ID, &dao.User{Growth: user.Growth})
	if err != nil {
		return err
	}
	err = Journal.AddScoreJournal(Orm, user.ID, "积分", "购买商品", play.ScoreJournal_Type_GM, int64(math.Floor(float64(orders.PayMoney)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
	if err != nil {
		return err
	}

	gvs := service.GiveVoucher.FindASC()
	for _, value := range gvs {
		//主订单的金额来决定是否送卡卷
		if uint64(orders.PayMoney) >= value.ScoreMaxValue {

			err := service.CardItem.AddVoucherCardItem(Orm, orders.OrderNo, orders.UserID, value.VoucherID)
			if err != nil {
				return err
			}
			break
		}
	}

	err=Journal.AddOrganizationJournal(Orm,orders.OID,"商品交易","商品交易",play.OrganizationJournal_Goods,int64(orders.PayMoney),dao.KV{Key: "OrdersID", Value: orders.ID})

	if err != nil {
		return err
	}
	for _,value:=range leves{
		if value<=0{
			break
		}
		var _user dao.User
		service.Get(Orm, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}


		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = Journal.AddUserJournal(Orm, _user.ID, "佣金", "一级用户", play.UserJournal_Type_LEVE, leveMenoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
		if err != nil {
			return err
		}

		err = service.AddUserBlockAmount(Orm,_user.ID,  -leveMenoy)
		if err != nil {
			return err
		}

		err=Journal.AddOrganizationJournal(Orm,orders.OID,"商品交易","推广佣金"+_user.Name,play.OrganizationJournal_Brokerage,-leveMenoy,dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			return err
		}

		err = Journal.AddScoreJournal(Orm, _user.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leveMenoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			return err
		}

		user=_user
	}

	/*if leve1 > 0 {

		var usera dao.User
		service.Get(Orm, user.SuperiorID, &usera)
		if usera.ID <= 0 {
			return nil
		}
		leve1Menoy := int64(math.Floor(float64(leve1)/float64(100)*float64(Brokerage) + 0.5))
		err = Journal.AddUserJournal(Orm, usera.ID, "佣金", "一级用户", play.UserJournal_Type_LEVE, leve1Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
		if err != nil {
			return err
		}
		err = Journal.AddScoreJournal(Orm, usera.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve1Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			return err
		}

		if leve2 > 0 {

			var userb dao.User
			service.Get(Orm, usera.SuperiorID, &userb)
			if userb.ID <= 0 {
				return nil
			}
			leve2Menoy := int64(math.Floor(float64(leve2)/float64(100)*float64(Brokerage) + 0.5))
			err = Journal.AddUserJournal(Orm, userb.ID, "佣金", "二级用户", play.UserJournal_Type_LEVE, leve2Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
			if err != nil {
				return err
			}
			err = Journal.AddScoreJournal(Orm, userb.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve2Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
			if err != nil {
				return err
			}

			if leve3 > 0 {

				var userc dao.User
				service.Get(Orm, userb.SuperiorID, &userc)
				if userc.ID <= 0 {
					return nil
				}
				leve3Menoy := int64(math.Floor(float64(leve3)/float64(100)*float64(Brokerage) + 0.5))
				err = Journal.AddUserJournal(Orm, userc.ID, "佣金", "三级用户", play.UserJournal_Type_LEVE, leve3Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
				if err != nil {
					return err
				}
				err = Journal.AddScoreJournal(Orm, userc.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve3Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
				if err != nil {
					return err
				}
				if leve4 > 0 {

					var userd dao.User
					service.Get(Orm, userc.SuperiorID, &userd)
					if userd.ID <= 0 {
						return nil
					}
					leve4Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
					err = Journal.AddUserJournal(Orm, userd.ID, "佣金", "四级用户", play.UserJournal_Type_LEVE, leve4Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
					if err != nil {
						return err
					}
					err = Journal.AddScoreJournal(Orm, userd.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve4Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
					if err != nil {
						return err
					}

					if leve5 > 0 {

						var usere dao.User
						service.Get(Orm, userd.SuperiorID, &usere)
						if usere.ID <= 0 {
							return nil
						}
						leve5Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
						err = Journal.AddUserJournal(Orm, usere.ID, "佣金", "五级用户", play.UserJournal_Type_LEVE, leve5Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
						if err != nil {
							return err
						}
						err = Journal.AddScoreJournal(Orm, usere.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve5Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
						if err != nil {
							return err
						}

						if leve6 > 0 {

							var userf dao.User
							service.Get(Orm, usere.SuperiorID, &userf)
							if userf.ID <= 0 {
								return nil
							}
							leve6Menoy := int64(math.Floor(float64(leve4)/float64(100)*float64(Brokerage) + 0.5))
							err = Journal.AddUserJournal(Orm, userf.ID, "佣金", "六级用户", play.UserJournal_Type_LEVE, leve6Menoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
							if err != nil {
								return err
							}
							err = Journal.AddScoreJournal(Orm, userf.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(math.Floor(float64(leve6Menoy)/100+0.5)), dao.KV{Key: "OrdersID", Value: orders.ID})
							if err != nil {
								return err
							}

						}

					}

				}

			}

		}

	}*/

	return nil
}
func (service UserService) Leve1(UserID uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if UserID <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID=?", UserID).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve2(Leve1IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve1IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve1IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve3(Leve2IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve2IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve2IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve4(Leve3IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve3IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve3IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve5(Leve4IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve4IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve4IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve6(Leve5IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve5IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve5IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) GetUserInfo(UserID uint64) dao.UserInfo {
	Orm := dao.Orm()
	//.First(&user, 10)
	var userInfo dao.UserInfo
	Orm.Where(&dao.UserInfo{UserID: UserID}).First(&userInfo)
	if userInfo.ID == 0 && UserID != 0 {
		userInfo.UserID = UserID
		service.Add(Orm, &userInfo)
	}
	return userInfo
}
func (service UserService) ListAllTableDatas(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.User{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (service UserService) UserAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	action := context.Request.URL.Query().Get("action")
	switch action {
	case "list":
		dts := &dao.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.User{}, company.ID)
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
	}

	return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "", Data: nil}}
}
func (service UserService) FindUserByTel(Orm *gorm.DB, Tel string) *dao.User {
	user := &dao.User{}
	err := Orm.Where("Tel=?", Tel).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}

func (service UserService) FindUserByOpenID(Orm *gorm.DB, OpenID string) *dao.User {

	user := &dao.User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}
func (service UserService) AddUserByOpenID(OpenID string) *dao.User {
	Orm := dao.Orm()
	user := &dao.User{}
	user = service.FindUserByOpenID(Orm, OpenID)
	if user.ID == 0 {
		user.OpenID = OpenID
		service.Add(Orm, user)
	} else {

	}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	//err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	//tool.CheckError(err)
	return user
}
