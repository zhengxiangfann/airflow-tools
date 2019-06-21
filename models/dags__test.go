package models

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	data := `
	--上周日所处周数，如：20190609为201923
set yearweeknum_upperlimit=year(date_add(from_unixtime(unix_timestamp()),1 - dayofweek(from_unixtime(unix_timestamp()))))*100+weekofyear(date_add(from_unixtime(unix_timestamp()),1 - dayofweek(from_unixtime(unix_timestamp()))));
-- 2019年款保留2018年最后4周及以后的数据
set product_year_name=2019;
-- 上周日
set invsunday = regexp_replace(date_add(from_unixtime(unix_timestamp()),1 - dayofweek(from_unixtime(unix_timestamp()))),'-','');
-- 上上周日
set last_invsunday = regexp_replace(date_add(from_unixtime(unix_timestamp()),1 - dayofweek(from_unixtime(unix_timestamp()))-7),'-','');
set version_tag = v5;
--权重，随时间调整
set wgt_store_spr=0.3; 
set wgt_sku_spr=0.4;
set wgt_store_smr=0.4;
set wgt_sku_smr=0.3;
	`
	re2, _ := regexp.Compile(`/\*[^*]*\*+(?:[^*/][^*]*\*+)*/`);
	re3, _ := regexp.Compile(`^\s*(--|#)`)
	rep := re2.ReplaceAllString(data, "");
	reps := strings.Split(rep,"\n")
	var result1 []string
	var result2 []string
	for _, line := range reps{
		if !re3.Match([]byte(line)){
			result1 = append(result1, line)
		}
	}
	log.Print(result1)
	log.Print("------------------")
	for _, line := range result1{
			l1 := strings.Split(line,"--")
			var l string
			if len(l1)>1{
				l = l1[0]
			} else {
				l = line
			}
			l2 := strings.Split(l, "#")
			if len(l2) > 1 {
				l =  l2[0]
			}
			result2 = append(result2, l)
	}
	log.Print(result2)
	log.Print(strings.Join(result2, "\n"))
}