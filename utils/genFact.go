package utils

import (
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateFact() map[string]string {
	fact := make(map[string]string)
	startDate := gofakeit.DateRange(time.Now().AddDate(0, -3, 0), time.Now())
	fact["period_start"] = startDate.Format("2006-01-02")
	fact["period_end"] = startDate.AddDate(0, 1, 0).Format("2006-01-02")
	fact["period_key"] = "month"
	fact["indicator_to_mo_id"] = "227373"
	fact["indicator_to_mo_fact_id"] = strconv.Itoa(gofakeit.Number(0, 1))
	fact["value"] = strconv.Itoa(gofakeit.Number(0, 1))
	fact["fact_time"] = startDate.AddDate(0, 0, 15).Format("2006-01-02")
	fact["is_plan"] = strconv.Itoa(gofakeit.Number(0, 1))
	fact["auth_user_id"] = "40"
	fact["comment"] = gofakeit.City()
	return fact
}
