package dbtime

import (
	"database/sql/driver"
	"strings"
	"time"
)

type NormalDatetime time.Time

func (ndt NormalDatetime) Value() (driver.Value, error) {
	return time.Time(ndt), nil
}
func (c *NormalDatetime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", value) //parse time
	if err != nil {
		return err
	}
	*c = NormalDatetime(t) //set result using the pointer
	return nil
}

func (c NormalDatetime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02 15:04:05") + `"`), nil
}
