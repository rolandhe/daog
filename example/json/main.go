package main

import (
	"encoding/json"
	"fmt"
	"github.com/rolandhe/daog/ttypes"
	"time"
)

type Data struct {
	//Id int64
	//Name string
	CreateAt ttypes.NormalDatetime
	Modify ttypes.NilableDatetime
}

func main()  {
	d := &Data{
		//2,
		//"roland",
		ttypes.NormalDatetime(time.Now()),
		*ttypes.FromDatetime(time.Now()),
	}
	j,_ := json.Marshal(d)

	var t Data
	json.Unmarshal(j,&t)
	fmt.Println(t)
}
