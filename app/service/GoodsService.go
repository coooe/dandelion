package service

import (
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type GoodsService struct {
	dao.BaseDao
	TimeSell TimeSellService
	Collage  CollageService
}

func (service GoodsService) GetSpecification(ID uint64, target *dao.Specification) error {
	Orm := dao.Orm()
	err := service.Get(Orm, ID, &target)

	return err
}
func (service GoodsService) OrdersStockManager(orders dao.Orders, isMinus bool) {

	go func() {

		if orders.PostType == 2 {
			//线下订单，不去维护在线商品库存
			log.Println("线下订单，不去维护在线商品库存")
			return
		}

		//管理商品库存
		Orm := dao.Orm()
		//list []dao.OrdersGoods

		list, _ := GlobalService.Orders.FindOrdersGoodsByOrdersID(Orm, orders.ID)
		for _, value := range list {
			var specification dao.Specification
			//service.Get(Orm, value.SpecificationID, &specification)
			util.JSONToStruct(value.Specification, &specification)
			var goods dao.Goods
			//service.Get(Orm, value.GoodsID, &goods)
			util.JSONToStruct(value.Goods, &goods)

			if isMinus {
				//减
				Stock := int64(specification.Stock - value.Quantity)
				if Stock < 0 {
					Stock = 0
				}
				err := service.ChangeMap(Orm, specification.ID, &dao.Specification{}, map[string]interface{}{"Stock": uint(Stock)})
				tool.CheckError(err)
				Stock = int64(goods.Stock - value.Quantity)
				if Stock < 0 {
					Stock = 0
				}
				err = service.ChangeMap(Orm, goods.ID, &dao.Goods{}, map[string]interface{}{"Stock": uint(Stock)})
				tool.CheckError(err)
			} else {
				//添加
				Stock := int64(specification.Stock + value.Quantity)
				if Stock < 0 {
					Stock = 0
				}
				err := service.ChangeMap(Orm, specification.ID, &dao.Specification{}, map[string]interface{}{"Stock": uint(Stock)})
				tool.CheckError(err)
				Stock = int64(goods.Stock + value.Quantity)
				if Stock < 0 {
					Stock = 0
				}
				err = service.ChangeMap(Orm, goods.ID, &dao.Goods{}, map[string]interface{}{"Stock": uint(Stock)})
				tool.CheckError(err)
			}

		}

	}()

}

/*func (service GoodsService) AddSpecification(context *gweb.Context) gweb.Result {
	item := &dao.Specification{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.Add(Orm, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (service GoodsService) ListSpecification(context *gweb.Context) gweb.Result {
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	var dts []dao.Specification
	service.FindWhere(Orm, &dts, dao.Specification{GoodsID: GoodsID})
	return &gweb.JsonResult{Data: &dao.ActionStatus{Success: true, Message: "OK", Data: dts}}
}*/
func (service GoodsService) DeleteSpecification(ID uint64) error {
	Orm := dao.Orm()
	err := service.Delete(Orm, &dao.Specification{}, ID)
	return err
}
func (service GoodsService) ChangeSpecification(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	GoodsID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)
	item := &dao.Specification{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	err = service.ChangeModel(Orm, GoodsID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (service GoodsService) SaveGoods(goods dao.Goods, specifications []dao.Specification) error {
	Orm := dao.Orm()
	var err error
	tx := Orm.Begin()

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if tx.NewRecord(&goods) {
		err = tx.Create(&goods).Error
	} else {
		//err = tx.Save(goods).Error
		err = tx.Model(&goods).Updates(goods).Error
	}

	if err != nil {
		return err
	}

	//添加或修改的时候不删除规格
	/*err=service.UnscopedDeleteWhere(tx,&dao.Specification{},"GoodsID=?",goods.ID)
	if err!=nil{
		return err
	}*/

	var total uint
	for _, value := range specifications {

		value.GoodsID = goods.ID

		if tx.NewRecord(&goods) {
			err = tx.Create(&value).Error
			total = total + value.Stock
		} else {
			err = tx.Save(&value).Error
			//err = tx.Model(&goods).Updates(goods).Error
			total = total + value.Stock
		}

		if err != nil {
			return err
		}

	}

	goods.Stock = total

	err = tx.Save(&goods).Error

	return err
}
func (service GoodsService) GetGoodsInfo(goods dao.Goods) dao.GoodsInfo {
	//Orm := dao.Orm()

	//user := dao.User{}
	//service.User.Get(Orm, UserID, &user)
	//brokerageProvisoConf := service.Configuration.GetConfiguration(play.ConfigurationKey_BrokerageProviso)
	//brokerageProvisoConfV, _ := strconv.ParseUint(brokerageProvisoConf.V, 10, 64)
	//vipdiscountConf := service.Configuration.GetConfiguration(play.ConfigurationKey_VIPDiscount)
	//VIPDiscount, _ := strconv.ParseUint(vipdiscountConf.V, 10, 64)
	timeSell := service.TimeSell.GetTimeSellByGoodsID(goods.ID)
	goodsInfo := dao.GoodsInfo{}
	goodsInfo.Goods = goods
	goodsInfo.Favoured = dao.Favoured{}

	if timeSell.IsEnable() {
		//Favoured:=uint64(util.Rounding45(float64(goods.Price)*(float64(timeSell.Discount)/float64(100)), 2))
		goodsInfo.Favoured = dao.Favoured{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint64(timeSell.Discount)}
	} else {
		collage := service.Collage.GetCollageByGoodsID(goods.ID)
		if collage.ID != 0 && collage.TotalNum > 0 {
			goodsInfo.Favoured = dao.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)}
		}

	}
	return goodsInfo
}
func (service GoodsService) GetGoods(DB *gorm.DB, ID uint64) dao.GoodsInfo {
	Orm := dao.Orm()
	var goods dao.Goods
	err := service.Get(Orm, ID, &goods)
	tool.Trace(err)

	var specifications []dao.Specification
	err = service.FindWhere(Orm, &specifications, dao.Specification{GoodsID: ID})
	tool.Trace(err)

	goodsInfo := service.GetGoodsInfo(goods)
	goodsInfo.Specifications = specifications

	/*var mtimeSell dao.TimeSell
	err=TimeSellService{}.Get(Orm,goods.TimeSellID,&mtimeSell)
	tool.Trace(err)
	if mtimeSell.IsEnable(){
		timeSell = mtimeSell
	}else {
		timeSell = dao.TimeSell{}
	}*/

	return goodsInfo
	//return DB.Model(target).Related(&dao.Specification{}).Where("ID=?", ID).First(target).Error
	/*Orm := dao.Orm()
	err := service.Get(Orm, ID, &goods)
	tool.Trace(err)

	err = service.FindWhere(Orm, &specifications, dao.Specification{GoodsID: ID})
	tool.Trace(err)

	var mtimeSell dao.TimeSell
	err = TimeSellService{}.Get(Orm, goods.TimeSellID, &mtimeSell)
	tool.Trace(err)
	if mtimeSell.IsEnable() {
		timeSell = mtimeSell
	} else {
		timeSell = dao.TimeSell{}
	}

	return*/
	//return DB.Model(target).Related(&dao.Specification{}).Where("ID=?", ID).First(target).Error
}

func (service GoodsService) DeleteGoods(ID uint64) *dao.ActionStatus {
	Orm := dao.Orm()
	tx := Orm.Begin()
	err := service.Delete(tx, &dao.Goods{}, ID)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where("GoodsID=?", ID).Delete(dao.Specification{}).Error
	if err != nil {
		tx.Rollback()
	}

	defer func() {
		if err == nil {
			tx.Commit()
		}
	}()

	return (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)
}
func (service GoodsService) DeleteGoodsType(ID uint64) *dao.ActionStatus {
	Orm := dao.Orm()
	tx := Orm.Begin()
	var gtcs []dao.GoodsTypeChild
	tx.Where(&dao.GoodsTypeChild{GoodsTypeID: ID}).Find(&gtcs) //Updates(map[string]interface{}{"GoodsTypeID": 0})

	var err error
	if len(gtcs) <= 0 {
		err = service.Delete(tx, &dao.GoodsType{}, ID)
		if err != nil {
			tx.Rollback()
		}
	} else {
		return (&dao.ActionStatus{}).SmartError(err, "包含子类数据，不能删除", nil)
	}

	defer func() {
		if err == nil {
			tx.Commit()
		}
	}()
	return (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)
}
func (service GoodsService) DeleteGoodsTypeChild(GoodsTypeChildID uint64) *dao.ActionStatus {
	Orm := dao.Orm()
	tx := Orm.Begin()
	tx.Model(&dao.Goods{GoodsTypeChildID: GoodsTypeChildID}).Updates(map[string]interface{}{"GoodsTypeChildID": 0})
	err := service.Delete(tx, &dao.GoodsTypeChild{}, GoodsTypeChildID)
	if err != nil {
		tx.Rollback()
	}
	defer func() {
		if err == nil {
			tx.Commit()
		}
	}()

	return (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)
}

func (service GoodsService) DeleteTimeSellGoods(DB *gorm.DB, GoodsID uint64) error {
	timesell := service.TimeSell.GetTimeSellByGoodsID(GoodsID)
	err := service.Delete(DB, &dao.TimeSell{}, timesell.ID)
	tool.CheckError(err)
	return err
}
func (service GoodsService) DeleteCollageGoods(DB *gorm.DB, GoodsID uint64) error {
	timesell := service.Collage.GetCollageByGoodsID(GoodsID)
	err := service.Delete(DB, &dao.Collage{}, timesell.ID)
	tool.CheckError(err)
	return err
}
func (service GoodsService) FindGoodsByTimeSellID(TimeSellID uint64) []dao.Goods {
	Orm := dao.Orm()

	var timesell dao.TimeSell
	err := service.Get(Orm, TimeSellID, &timesell)
	tool.CheckError(err)

	var list []dao.Goods
	err = service.FindWhere(Orm, &list, "ID=?", timesell.GoodsID)
	tool.CheckError(err)
	return list
}
func (service GoodsService) FindGoodsByTimeSellHash(Hash string) []dao.Goods {
	Orm := dao.Orm()

	var GoodsIDs []uint64
	Orm.Model(&dao.TimeSell{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []dao.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	tool.CheckError(err)
	return list
}
func (service GoodsService) FindGoodsByCollageHash(Hash string) []dao.Goods {
	Orm := dao.Orm()

	var GoodsIDs []uint64
	Orm.Model(&dao.Collage{}).Where("Hash=?", Hash).Pluck("GoodsID", &GoodsIDs)

	var list []dao.Goods
	err := service.FindWhere(Orm, &list, "ID in (?)", GoodsIDs)
	tool.CheckError(err)
	return list
}

func (service GoodsService) AllList() []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Order("CreatedAt desc") //.Limit(10)

	db.Find(&result)

	return result

}
func (service GoodsService) GetGoodsInfoList(UserID uint64, goodsList []dao.Goods) []dao.GoodsInfo {

	var results = make([]dao.GoodsInfo, 0)

	for _, value := range goodsList {
		timeSell := service.TimeSell.GetTimeSellByGoodsID(value.ID)
		goodsInfo := dao.GoodsInfo{}
		goodsInfo.Goods = value
		goodsInfo.Favoured = dao.Favoured{}
		if timeSell.IsEnable() {
			//Favoured:=uint64(util.Rounding45(float64(value.Price)*(float64(timeSell.Discount)/float64(100)), 2))
			goodsInfo.Favoured = dao.Favoured{Name: "限时抢购", Target: util.StructToJSON(timeSell), TypeName: "TimeSell", Discount: uint64(timeSell.Discount)}
		} else {
			collage := service.Collage.GetCollageByGoodsID(value.ID)
			if collage.ID != 0 && collage.TotalNum > 0 {
				goodsInfo.Favoured = dao.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)}
			}

		}
		results = append(results, goodsInfo)
	}

	return results
}

func (service GoodsService) GoodsList(UserID uint64, SqlOrder string, Index int, where interface{}, args ...interface{}) []dao.GoodsInfo {
	Orm := dao.Orm()
	var goodsList []dao.Goods
	//db := Orm.Model(&dao.Goods{}).Order("CountSale desc").Limit(10)
	//db.Find(&result)
	err := service.FindWherePaging(Orm, SqlOrder, &goodsList, Index, where, args)
	tool.CheckError(err)

	return service.GetGoodsInfoList(UserID, goodsList)
}
func (service GoodsService) HotList() []dao.Goods {

	Orm := dao.Orm()

	var result []dao.Goods

	db := Orm.Model(&dao.Goods{}).Order("CountSale desc").Limit(10)

	db.Find(&result)

	return result

}
func (service GoodsService) ListGoodsType() []dao.GoodsType {
	Orm := dao.Orm()
	var gts []dao.GoodsType
	service.FindAll(Orm, &gts)
	return gts
}
func (service GoodsService) ListGoodsTypeChild(GoodsTypeID uint64) []dao.GoodsTypeChild {
	Orm := dao.Orm()
	var gts []dao.GoodsTypeChild
	service.FindWhere(Orm, &gts, dao.GoodsTypeChild{GoodsTypeID: GoodsTypeID})
	return gts
}
func (service GoodsService) ListGoodsChildByGoodsTypeID(GoodsTypeID, GoodsTypeChildID uint64) []dao.Goods {
	Orm := dao.Orm()
	var gts []dao.Goods
	service.FindWhere(Orm, &gts, dao.Goods{GoodsTypeID: GoodsTypeID, GoodsTypeChildID: GoodsTypeChildID})
	return gts
}
