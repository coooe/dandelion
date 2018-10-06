package admin

import (
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/util"
	"net/url"
	"strings"

	"dandelion/app/service/dao"

	"strconv"

	"errors"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type InterceptorManager struct {
}

//Execute(Session *Session,Request *http.Request)(bool,Result)
func (this InterceptorManager) Execute(context *gweb.Context) (bool, gweb.Result) {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionAdmin) == nil {
		//http.SetCookie(context.Response, &http.Cookie{Name: "UID", MaxAge:-1, Path: "/"})
		//fmt.Println(context.Request.URL.Path)
		//fmt.Println(context.Request.URL.Query().Encode())
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}

		//fmt.Println(url.QueryEscape(redirect))
		//http.Redirect(context.Response, context.Request, "/account/loginAdminPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false, &gweb.RedirectToUrlResult{Url: "/account/loginAdminPage?redirect=" + url.QueryEscape(redirect)}
	} else {
		return true, nil
	}
}

type Controller struct {
	gweb.BaseController
	Admin           service.AdminService
	Goods           service.GoodsService
	ScoreGoods      service.ScoreGoodsService
	Voucher         service.VoucherService
	Store           service.StoreService
	StoreStock      service.StoreStockService
	ExpressTemplate service.ExpressTemplateService
	FullCut         service.FullCutService
	TimeSell        service.TimeSellService
	Orders          service.OrdersService
	Rank            service.RankService
	Configuration   service.ConfigurationService
	GiveVoucher     service.GiveVoucherService
	User            service.UserService
	CardItem        service.CardItemService
	Content         service.ContentService
	Article         service.ArticleService
}

func (controller *Controller) Apply() {
	controller.Interceptors.Add(&InterceptorManager{})
	//Index.RequestMapping = make(map[string]mvc.Function)
	controller.AddHandler(gweb.ALLMethod("", controller.rootPage))
	controller.AddHandler(gweb.ALLMethod("*", controller.defaultPage))
	controller.AddHandler(gweb.ALLMethod("index", controller.indexPage))
	controller.AddHandler(gweb.ALLMethod("goods", controller.Goods.Action))

	controller.AddHandler(gweb.POSMethod("score_goods", controller.ScoreGoods.AddScoreGoods))
	controller.AddHandler(gweb.GETMethod("score_goods/:ID", controller.ScoreGoods.GetScoreGoods))
	controller.AddHandler(gweb.POSMethod("score_goods/list", controller.ScoreGoods.DatatablesScoreGoods))
	controller.AddHandler(gweb.DELMethod("score_goods/:ID", controller.ScoreGoods.DeleteScoreGoods))
	controller.AddHandler(gweb.PUTMethod("score_goods/:ID", controller.ScoreGoods.ChangeScoreGoods))

	controller.AddHandler(gweb.POSMethod("voucher", controller.Voucher.AddItem))
	controller.AddHandler(gweb.GETMethod("voucher/:ID", controller.Voucher.GetItem))
	controller.AddHandler(gweb.POSMethod("voucher/list", controller.Voucher.ListItem))
	controller.AddHandler(gweb.DELMethod("voucher/:ID", controller.Voucher.DeleteItem))
	controller.AddHandler(gweb.PUTMethod("voucher/:ID", controller.Voucher.ChangeItem))

	controller.AddHandler(gweb.POSMethod("fullcut/save", controller.FullCut.SaveItem))
	controller.AddHandler(gweb.GETMethod("fullcut/:ID", controller.FullCut.GetItem))
	controller.AddHandler(gweb.POSMethod("fullcut/datatables/list", controller.FullCut.DataTablesItem))
	controller.AddHandler(gweb.DELMethod("fullcut/:ID", controller.FullCut.DeleteItem))

	controller.AddHandler(gweb.POSMethod("timesell/save", controller.TimeSell.SaveItem))
	controller.AddHandler(gweb.GETMethod("timesell/:ID", controller.TimeSell.GetItem))
	controller.AddHandler(gweb.POSMethod("timesell/datatables/list", controller.TimeSell.DataTablesItem))
	controller.AddHandler(gweb.GETMethod("timesell/goods/:TimeSellID/list", controller.ListTimeSellGoods))
	controller.AddHandler(gweb.DELMethod("timesell/goods/:GoodsID", controller.DeleteTimeSellGoods))
	controller.AddHandler(gweb.DELMethod("timesell/:ID", controller.TimeSell.DeleteItem))

	controller.AddHandler(gweb.POSMethod("store", controller.Store.AddItem))
	controller.AddHandler(gweb.GETMethod("store/:ID", controller.Store.GetItem))
	controller.AddHandler(gweb.POSMethod("store/list", controller.Store.ListItem))
	controller.AddHandler(gweb.DELMethod("store/:ID", controller.Store.DeleteItem))
	controller.AddHandler(gweb.PUTMethod("store/:ID", controller.Store.ChangeItem))

	controller.AddHandler(gweb.POSMethod("store_stock", controller.StoreStock.SaveItem))
	controller.AddHandler(gweb.PUTMethod("store_stock", controller.StoreStock.SaveItem))

	controller.AddHandler(gweb.GETMethod("store_stock/:ID", controller.StoreStock.GetItem))
	controller.AddHandler(gweb.GETMethod("store_stock/exist/goods/:StoreID", controller.StoreStock.ListExistGoodsIDS))
	controller.AddHandler(gweb.POSMethod("store_stock/list/:StoreID/:GoodsID", controller.StoreStock.ListByGoods))
	controller.AddHandler(gweb.POSMethod("store_stock/list", controller.StoreStock.ListItem))
	controller.AddHandler(gweb.DELMethod("store_stock/:ID", controller.StoreStock.DeleteItem))

	controller.AddHandler(gweb.POSMethod("express_template/save", controller.saveExpressTemplate))
	controller.AddHandler(gweb.DELMethod("express_template/:ID", controller.deleteExpressTemplate))
	controller.AddHandler(gweb.GETMethod("express_template/:ID", controller.getExpressTemplate))
	controller.AddHandler(gweb.POSMethod("express_template/datatables/list", controller.datatablesExpressTemplate))
	controller.AddHandler(gweb.GETMethod("express_template/list", controller.listExpressTemplate))

	controller.AddHandler(gweb.POSMethod("store_journal/list", controller.storeJournalListAction))

	controller.AddHandler(gweb.POSMethod("order/list", controller.listOrderAction))
	controller.AddHandler(gweb.PUTMethod("order/change", controller.orderChangeAction))

	controller.AddHandler(gweb.POSMethod("configuration/list", controller.configurationListAction))
	controller.AddHandler(gweb.POSMethod("configuration/change", controller.configurationChangeAction))

	controller.AddHandler(gweb.POSMethod("rank/add", controller.rankAddAction))
	controller.AddHandler(gweb.POSMethod("rank/list", controller.rankListAction))
	controller.AddHandler(gweb.DELMethod("rank/:RankID", controller.rankDeleteAction))

	//carditem_list.html
	controller.AddHandler(gweb.POSMethod("carditem/list", controller.carditemListAction))

	controller.AddHandler(gweb.POSMethod("give_voucher/save", controller.giveVoucherSaveAction))
	controller.AddHandler(gweb.POSMethod("give_voucher/list", controller.giveVoucherListAction))
	controller.AddHandler(gweb.DELMethod("give_voucher/:GiveVoucherID", controller.giveVoucherDeleteAction))

	controller.AddHandler(gweb.POSMethod("situation", controller.situationAction))

	controller.AddHandler(gweb.POSMethod("admin", controller.Admin.AddItem))
	controller.AddHandler(gweb.GETMethod("admin/:ID", controller.Admin.GetItem))
	controller.AddHandler(gweb.POSMethod("admin/list", controller.Admin.ListItem))
	controller.AddHandler(gweb.DELMethod("admin/:ID", controller.Admin.DeleteItem))
	controller.AddHandler(gweb.PUTMethod("admin/:ID", controller.Admin.ChangePassWork))
	controller.AddHandler(gweb.PUTMethod("admin/authority/:ID", controller.Admin.ChangeAuthority))

	controller.AddHandler(gweb.POSMethod("user/all/list", controller.User.ListAllTableDatas))

	controller.AddHandler(gweb.ALLMethod("loginOut", controller.loginOutAction))

	//--------------content------------------
	controller.AddHandler(gweb.GETMethod("content_type/list", controller.Content.ListContentTypeAction))
	controller.AddHandler(gweb.POSMethod("content", controller.Content.AddContentAction))
	controller.AddHandler(gweb.GETMethod("content/:ID", controller.Content.GetContentAction))
	controller.AddHandler(gweb.GETMethod("content/list", controller.Content.ListContentsAction))
	controller.AddHandler(gweb.DELMethod("content/:ID", controller.Content.DeleteContentAction))
	controller.AddHandler(gweb.PUTMethod("content/:ID", controller.Content.ChangeContentAction))
	controller.AddHandler(gweb.PUTMethod("content/index/:ID", controller.Content.ChangeContentIndexAction))
	controller.AddHandler(gweb.PUTMethod("content/hide/:ID", controller.Content.ChangeHideContentAction))

	controller.AddHandler(gweb.POSMethod("content_sub_type", controller.Content.AddClassify))
	controller.AddHandler(gweb.GETMethod("content_sub_type/list/:ContentID", controller.Content.ListClassify))
	controller.AddHandler(gweb.GETMethod("content_sub_type/child/list/:ContentID/:ParentContentSubTypeID", controller.Content.ListChildClassify))
	controller.AddHandler(gweb.DELMethod("content_sub_type/:ID", controller.Content.DeleteClassify))
	controller.AddHandler(gweb.PUTMethod("content_sub_type/:ID", controller.Content.ChangeClassify))
	controller.AddHandler(gweb.GETMethod("content_sub_type/:ID", controller.Content.GetContentSubTypeAction))
	//------------------ArticleService.go-datatables------------------------
	controller.AddHandler(gweb.POSMethod("article/datatables/list", controller.Article.DataTablesAction))
	controller.AddHandler(gweb.POSMethod("article/save", controller.Article.SaveArticleAction))
	controller.AddHandler(gweb.POSMethod("article/delete", controller.Article.DeleteArticleAction))
	controller.AddHandler(gweb.GETMethod("article/get/:ID", controller.Article.GetArticleAction))

}

func (controller *Controller) carditemListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.CardItem.DatatablesListOrder(Orm, dts, &[]dao.CardItem{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) situationAction(context *gweb.Context) gweb.Result {
	context.Request.ParseForm()
	StartTime, _ := strconv.ParseInt(context.Request.FormValue("StartTime"), 10, 64)
	EndTime, _ := strconv.ParseInt(context.Request.FormValue("EndTime"), 10, 64)

	results := make(map[string]interface{})

	results["Orders"] = controller.Orders.Situation(StartTime, EndTime)
	results["Voucher"] = controller.Voucher.Situation(StartTime, EndTime)
	results["ScoreGoods"] = controller.ScoreGoods.Situation(StartTime, EndTime)
	results["User"] = controller.User.Situation(StartTime, EndTime)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", results)}
}
func (controller *Controller) giveVoucherSaveAction(context *gweb.Context) gweb.Result {

	item := dao.GiveVoucher{}
	util.RequestBodyToJSON(context.Request.Body, &item)
	err := controller.GiveVoucher.SaveItem(item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (controller *Controller) giveVoucherListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.ExpressTemplate.DatatablesListOrder(Orm, dts, &[]dao.GiveVoucher{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) giveVoucherDeleteAction(context *gweb.Context) gweb.Result {

	GiveVoucherID, _ := strconv.ParseUint(context.PathParams["GiveVoucherID"], 10, 64)

	err := controller.Rank.Delete(dao.Orm(), &dao.GiveVoucher{}, GiveVoucherID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}

func (controller *Controller) configurationChangeAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := dao.Configuration{}
	util.RequestBodyToJSON(context.Request.Body, &item)
	err := controller.Configuration.ChangeConfiguration(company.ID, item.K, item.V)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (controller *Controller) configurationListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var ks []uint64
	util.RequestBodyToJSON(context.Request.Body, &ks)
	list := controller.Configuration.GetConfigurations(company.ID, ks)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "OK", list)}
}
func (controller *Controller) rankAddAction(context *gweb.Context) gweb.Result {

	rank := dao.Rank{}
	util.RequestBodyToJSON(context.Request.Body, &rank)
	err := controller.Rank.AddRank(rank)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}
func (controller *Controller) rankListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.ExpressTemplate.DatatablesListOrder(Orm, dts, &[]dao.Rank{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (controller *Controller) rankDeleteAction(context *gweb.Context) gweb.Result {

	RankID, _ := strconv.ParseUint(context.PathParams["RankID"], 10, 64)

	err := controller.Rank.Delete(dao.Orm(), &dao.Rank{}, RankID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *Controller) orderChangeAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	context.Request.ParseForm()
	Action := context.Request.FormValue("Action")
	switch Action {
	case "RefundComplete":
		OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		RefundType, _ := strconv.ParseUint(context.Request.FormValue("RefundType"), 10, 64)
		err, info := controller.Orders.RefundComplete(OrdersGoodsID, RefundType)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "RefundOk":
		OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		err, info := controller.Orders.RefundOk(OrdersGoodsID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "RefundNo":
		OrdersGoodsID, _ := strconv.ParseUint(context.Request.FormValue("OrdersGoodsID"), 10, 64)
		err, info := controller.Orders.RefundNo(OrdersGoodsID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "Cancel":
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		err, info := controller.Orders.Cancel(ID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "CancelOk":
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		RefundType, _ := strconv.ParseUint(context.Request.FormValue("RefundType"), 10, 64) //退款资金来源	 0=未结算资金退款,1=可用余额退款
		err, info := controller.Orders.CancelOk(ID, RefundType)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, info, nil)}
	case "Deliver":
		ShipName := context.Request.FormValue("ShipName")
		ShipNo := context.Request.FormValue("ShipNo")
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)

		err := controller.Orders.Deliver(ShipName, ShipNo, ID)
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "发货成功", nil)}
	case "PayMoney":
		PayMoney, _ := strconv.ParseFloat(context.Request.FormValue("PayMoney"), 64)
		ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
		err := controller.Orders.ChangeMap(Orm, ID, &dao.Orders{}, map[string]interface{}{"PayMoney": uint64(PayMoney * 100)})
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
		success, message := controller.Orders.ChangeOrdersPayMoney(PayMoney, ID)
		return &gweb.JsonResult{Data: &dao.ActionStatus{Success: success, Message: message, Data: nil}}

	}

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("999"), "OK", nil)}

}
func (controller *Controller) storeJournalListAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.CardItem.DatatablesListOrder(Orm, dts, &[]dao.StoreJournal{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}

}
func (controller *Controller) listOrderAction(context *gweb.Context) gweb.Result {

	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	//Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	UserID, _ := strconv.ParseUint(dts.Columns[0].Search.Value, 10, 64)
	PostType, _ := strconv.ParseInt(dts.Columns[1].Search.Value, 10, 64)
	Status := dts.Columns[2].Search.Value

	var StatusList []string
	if !strings.EqualFold(Status, "") {
		StatusList = strings.Split(Status, ",")
	}
	//fmt.Println(dts)
	list, recordsTotal := controller.Orders.ListOrders(UserID, company.ID, int(PostType), StatusList, dts.Length, dts.Start)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": dts.Draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsTotal}}
}
func (controller *Controller) DeleteTimeSellGoods(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["GoodsID"], 10, 64)

	list := controller.Goods.DeleteTimeSellGoods(Orm, ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "删除成功", list)}

}
func (controller *Controller) ListTimeSellGoods(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["TimeSellID"], 10, 64)

	list := controller.Goods.FindByTimeSellID(ID)
	//var item dao.ExpressTemplate
	//err := controller.ExpressTemplate.Get(service.Orm, ID, &item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(nil, "", list)}
	//2002
}
func (controller *Controller) getExpressTemplate(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	var item dao.ExpressTemplate
	err := controller.ExpressTemplate.Get(Orm, ID, &item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", item)}
	//2002
}
func (controller *Controller) deleteExpressTemplate(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)

	err := controller.ExpressTemplate.Delete(Orm, &dao.ExpressTemplate{}, ID)

	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (controller *Controller) saveExpressTemplate(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	item := &dao.ExpressTemplate{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	tool.Trace(err)
	item.OID = company.ID
	err = controller.ExpressTemplate.SaveExpressTemplate(item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "保存成功", nil)}
}
func (controller *Controller) listExpressTemplate(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	var list []dao.ExpressTemplate
	err := controller.ExpressTemplate.FindAllByOID(Orm, &list, company.ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", list)}
}
func (controller *Controller) datatablesExpressTemplate(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := controller.ExpressTemplate.DatatablesListOrder(Orm, dts, &[]dao.ExpressTemplate{}, company.ID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}

}

func (controller *Controller) loginOutAction(context *gweb.Context) gweb.Result {
	context.Session.Attributes.Delete(play.SessionAdmin)
	return &gweb.RedirectToUrlResult{Url: "/admin"}
}

func (controller *Controller) indexPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}

func (controller *Controller) defaultPage(context *gweb.Context) gweb.Result {

	return &gweb.HTMLResult{}
}
func (controller *Controller) rootPage(context *gweb.Context) gweb.Result {

	return &gweb.RedirectToUrlResult{"index"}
}
