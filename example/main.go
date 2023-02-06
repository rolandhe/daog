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
	}
	var err error
	datasource, err = daog.NewDatasource(conf)
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {
	defer datasource.Shutdown()

	//createUser()
	//create()
	//query()
	//queryUser()
	queryRawSQLForCount()
	//queryByIds()
	//queryByIdsUsingDao()
	//queryByMatcher()
	//queryByMatcherOrder()
	//countByMatcher()
	//update()

	//deleteById()
}

func query() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer daog.CompleteTransContext(tc, err)

	g, err := daog.GetById(tc,9, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, err := json.Marshal(g)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("query", string(j))
	var rg dal.GroupInfo
	json.Unmarshal(j,&rg)
	fmt.Println(g.CreateAt)
	fmt.Println(string(g.BinData))
}

func queryUser() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer daog.CompleteTransContext(tc, err)

	g, err := daog.GetById(tc,1, dal.UserInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, err := json.Marshal(g)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("queryUser", string(j))
	var rg dal.UserInfo
	json.Unmarshal(j,&rg)
	fmt.Println(rg)
}

func deleteById() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.DeleteById(tc,2, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("delete", g)
}

func queryByIds() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	gs, err := daog.GetByIds(tc,[]int64{1, 2}, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIds", string(j))
	fmt.Println(gs)
}

func queryByIdsUsingDao() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	gs, err := dal.GroupInfoDao.GetByIds(tc,[]int64{1, 2})
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByIdsUsingDao", string(j))
	fmt.Println(gs)
}

func queryByMatcher() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer func() {
		tc.Complete(err)
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	gs, err := daog.QueryListMatcher(tc,matcher, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(gs)
	fmt.Println("queryByMatcher", string(j))
	fmt.Println(gs)
}

func queryByMatcherOrder() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestNone, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer func() {
		tc.Complete(err)
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	gs, err := daog.QueryListMatcher(tc,matcher, dal.GroupInfoMeta,  daog.NewDescOrder(dal.GroupInfoFields.Id))
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
	// 无事务情况下也需要加上这段代码，用于释放底层链接
	defer func() {
		tc.Complete(err)
	}()
	matcher := daog.NewMatcher().Like(dal.GroupInfoFields.Name, "roland", daog.LikeStyleLeft).Lt(dal.GroupInfoFields.Id, 4)
	c, err := daog.Count(tc,matcher, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("countByMatcher",c)
	fmt.Println(c)
}


func create() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	amount, err := decimal.NewFromString("128.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &dal.GroupInfo{
		Name:        "roland",
		MainData:    `{"a":102}`,
		Content:     "hello world!!",
		BinData:     []byte("byte data"),
		CreateAt:    ttypes.NormalDatetime(time.Now()),
		TotalAmount: amount,
		Remark:      *ttypes.FromString("haha"),
	}
	affect, err := daog.Insert(tc,t, dal.GroupInfoMeta)
	fmt.Println(affect, t.Id, err)

	t.Name = "roland he"
	af, err := daog.Update(tc,t, dal.GroupInfoMeta)
	fmt.Println(af, err)
}

func createUser() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-1001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &dal.UserInfo{
		Name:        "roland",
		CreateAt:    ttypes.NormalDatetime(time.Now()),
		ModifyAt: *ttypes.FromDatetime(time.Now()),
	}
	affect, err := daog.Insert(tc,t, dal.UserInfoMeta)
	fmt.Println(affect, t.Id, err)
}

func update() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestWrite, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()
	g, err := daog.GetById(tc,5, dal.GroupInfoMeta)
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(g)
	fmt.Println("query", string(j))

	g.Name = "Eric"
	af, err := daog.Update(tc,g, dal.GroupInfoMeta)
	fmt.Println(af, err)

}

type GroupInfoCounter struct {
	Name string
	Count int64
}
func queryRawSQLForCount() {
	tc, err := daog.NewTransContext(datasource, txrequest.RequestReadonly, "trace-100099")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		tc.Complete(err)
	}()

	list, err := daog.QueryRawSQL(tc, func (ins *GroupInfoCounter) []any{
		ret := make([]any,2)
		ret[0] = &ins.Name
		ret[1] = &ins.Count
		return ret
	},"select name,count(id) from group_info group by name")
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(list)
	fmt.Println("queryRawSQLForCount", string(j))
}