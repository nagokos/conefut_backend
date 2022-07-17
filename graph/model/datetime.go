package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nagokos/connefut_backend/logger"
)

/*
	APIサーバで受け取ったdatetimeフィールドの値をtime.Time型に変換する。
*/
func UnmarshalDateTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		if len(v) == 0 {
			return time.Time{}, nil
		}
		return time.ParseInLocation("2006/01/02 15:04", v, time.Local)
	case time.Time:
		return v, nil
	default:
		return time.Now(), fmt.Errorf("DateTime is invalid")
	}
}

/*
	APIサーバからJSONを返す際に、time.Time型をstringに変換する。
*/
func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if t.IsZero() {
			_, err := w.Write([]byte(`""`))
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}
		} else {
			_, err := w.Write([]byte(strconv.Quote(t.Format("2006/01/02 15:04"))))
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}
		}
	})
}
