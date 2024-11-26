package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/example/dal"
	"github.com/rolandhe/daog/ttypes"
	txrequest "github.com/rolandhe/daog/tx"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

var datasource daog.Datasource

func init() {
	conf := &daog.DbConf{
		DbUrl:  "root:12345678@tcp(localhost:3306)/daog?parseTime=true&timeout=1s&readTimeout=2s&writeTimeout=2s",
		LogSQL: true,
		Size:   1,
	}
	var err error
	datasource, err = daog.NewDatasource(conf)
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {
	defer datasource.Shutdown()

	//testMapConn()
	//createUserUseAutoTrans()
	//createUser()z
	//create()
	//query()
	//queryUser()
	//queryUserPageForUpdate()
	//queryRawSQLForCount()
	//queryByIds()
	//queryByIdsUsingDao()
	//queryByMatcher()
	//queryAll()
	//queryByMatcherOrder()
	//countByMatcher()
	update()

	//deleteById()
}

//func testMapConn() {
//	_, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
//	//defer tc.Complete(nil)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	//go func() {
//	//	time.Sleep(time.Second * 5)
//	//	tc.Complete(nil)
//	//}()
//
//	tc2, err2 := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
//	if err2 != nil {
//		fmt.Println(err)
//		return
//	}
//	tc2.Complete(nil)
//}

func query() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	}
	g, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) (*dal.GroupInfo, error) {
		return daog.GetById(tc, 9, dal.GroupInfoMeta)
	})

	if err != nil {
		fmt.Println(err)
	}
	j, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("query", string(j))
	var rg dal.GroupInfo
	json.Unmarshal(j, &rg)
	fmt.Println(g.CreateAt)
	fmt.Println(string(g.BinData))
}

func queryUser() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	}
	g, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) (*dal.UserInfo, error) {
		return daog.GetById(tc, 1, dal.UserInfoMeta)
	})
	if err != nil {
		fmt.Println(err)
	}
	j, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("queryUser", string(j))
	var rg dal.UserInfo
	json.Unmarshal(j, &rg)
	fmt.Println(rg)
}

func queryUserPageForUpdate() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	}
	mat := daog.NewMatcher()
	mat.In(dal.UserInfoFields.Id, []any{1, 3, 5})
	pager := &daog.Pager{
		PageNumber: 1,
		PageSize:   2,
	}
	viewColumns := []string{
		dal.UserInfoFields.Id,
		dal.UserInfoFields.Name,
		dal.UserInfoFields.CreateAt,
	}
	userInfos, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.UserInfo, error) {
		return daog.QueryPageListMatcherWithViewColumnsForUpdate(tc, mat, dal.UserInfoMeta, viewColumns, pager, true)
	})
	if err != nil {
		fmt.Println(err)
	}
	j, err := json.Marshal(userInfos)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("queryUser", string(j))

}

func deleteById() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}
	daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		g, err := daog.DeleteById(tc, 2, dal.GroupInfoMeta)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("delete", g)
		return err
	})
}

func queryByIds() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	}
	gs, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return daog.GetByIds(tc, []int64{1, 2}, dal.GroupInfoMeta)
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIds", string(j))
	fmt.Println(gs)
}

func queryByIdsUsingDao() {
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	}
	gs, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return dal.GroupInfoDao.GetByIds(tc, []int64{1, 2})
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIdsUsingDao", string(j))
	fmt.Println(gs)
}

func queryByMatcher() {
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleRight).Lt(dal.GroupInfoFields.Id, 4)
	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	}
	gs, err := daog.AutoTransWithResult(tcCreate, func(tc *daog.TransContext) ([]*dal.GroupInfo, error) {
		return daog.QueryListMatcher(tc, matcher, dal.GroupInfoMeta)
	})

	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcher", string(j))
	fmt.Println(gs)
}

func queryAll() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	// 必须使用匿名函数，不能使用 tc.CompleteWithPanic(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	gs, err := daog.GetAll(tc, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryAll", string(j))
	fmt.Println(gs)
}

func queryByMatcherOrder() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	// 必须使用匿名函数，不能使用 tc.Complete(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	gs, err := daog.QueryListMatcher(tc, matcher, dal.GroupInfoMeta, daog.NewDescOrder(dal.GroupInfoFields.Id))
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcherOrder", string(j))
	fmt.Println(gs)
}

func countByMatcher() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	// 必须使用匿名函数，不能使用 tc.CompleteWithPanic(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	c, err := daog.Count(tc, matcher, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("countByMatcher", c)
	fmt.Println(c)
}

func create() {
	amount, err := decimal.NewFromString("128.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &dal.GroupInfo{
		Name:        "roland-one",
		MainData:    `{"a":102}`,
		Content:     "hello world!!",
		BinData:     []byte("byte data"),
		CreateAt:    ttypes.NormalDatetime(time.Now()),
		TotalAmount: amount,
		Remark:      *ttypes.FromString("haha"),
	}

	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}
	daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		affect, err := daog.Insert(tc, t, dal.GroupInfoMeta)
		fmt.Println(affect, t.Id, err)
		if err != nil {
			return err
		}
		t.Name = "rolandx"
		af, err := daog.Update(tc, t, dal.GroupInfoMeta)
		fmt.Println(af, err)
		if err != nil {
			return err
		}
		return nil
	})
}

func createUserUseAutoTrans() {
	t := &dal.UserInfo{
		Name: "roland",
		//CreateAt: ttypes.NormalDatetime(time.Now()),
		//ModifyAt: *ttypes.FromDatetime(time.Now()),
	}
	affect, err := daog.AutoTransWithResult[int64](func() (*daog.TransContext, error) {
		return daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	}, func(tc *daog.TransContext) (int64, error) {
		return daog.Insert(tc, t, dal.UserInfoMeta)
	})
	fmt.Println(affect, t.Id, err)
}

func update() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 必须使用匿名函数，不能使用 tc.CompleteWithPanic(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()
	g, err := dal.UserInfoDao.GetById(tc, 1)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))

	g.Name = "Eric1"

	af, err := dal.UserInfoDao.Update(tc, g)
	fmt.Println(af, err)

}

type GroupInfoCounter struct {
	Name  string
	Count int64
}

func queryRawSQLForCount() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	// tc 创建后必须马上跟上 defer func， 如果这之间有return或者panic，连接将被泄露
	// 必须使用匿名函数，不能使用 tc.CompleteWithPanic(err)， 因为defer 后面函数的参数在执行defer语句是就会被确定
	defer func() {
		// 注意：后面代码的error都要使用err变量来接收，否则在发生错误的情况下，事务不会被回滚
		tc.CompleteWithPanic(err, recover())
	}()

	list, err := daog.QueryRawSQL(tc, func(ins *GroupInfoCounter) []any {
		ret := make([]any, 2)
		ret[0] = &ins.Name
		ret[1] = &ins.Count
		return ret
	}, "select name,count(id) from group_info group by name")
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(list)
	fmt.Println("queryRawSQLForCount", string(j))
}
