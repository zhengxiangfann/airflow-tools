package utils

import (
	"log"
	"testing"
)

func TestHQL_SplitSql(t *testing.T) {
	sql := `
	drop table belle_sh.hg_shanghai_19_sale_inv_last_week_temp;
create table belle_sh.hg_shanghai_19_sale_inv_last_week_temp as
select *
from belle_sh.hg_shanghai_19_sale_inv;

drop table belle_sh.lrs_st_open_store_on_UI;
create table belle_sh.lrs_st_open_store_on_UI 
(store_no string)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY ','
STORED AS TEXTFILE
;
insert into belle_sh.lrs_st_open_store_on_UI
values
('I0ZHST')
,('IAJ013')
,('I038ST')
,('I022ST')
,('I032ST')
,('I023ST')
,('I045ST')
,('I01IST')
,('I010ST')
,('I044ST')
,('I051ST')
,('IAJ002')
,('I01MST')
,('I01SST')
,('IAJ009')
,('I01BST')
,('IAJ004')
,('IAJ001')
,('I01JST')
,('IAJ003')
,('IAJ011')
,('I014ST')
,('I046ST')
,('I072ST')
,('I028ST')
,('IAJ014')
,('I078ST')
,('I01FST')
,('IAJ007')
,('I064ST')
,('I063ST')
,('I070ST')
,('I067ST')
,('I035ST')
,('IAJ006')
,('I029ST')
,('I033ST')
,('I01HST')
,('I074ST')
,('I066ST')
,('IAJ008')
,('I056ST')
,('IAJ012')
,('I048ST')
,('I01WST')
,('IA02MA')
,('IAJ016') -- map to I01MST
,('IAJA04') -- map to I046ST
;
`
	hsql := HQL{
		Hql:sql,
	}
	//hsql.TrimSql()
	//log.Print(hsql.RHql)
	res := hsql.SplitSql()
	log.Print("res length = ", len(res))
	log.Print(res[0])
	if len(res) > 1 {
		ctable, tb := hsql.ExtraceTableName(res[1])
		log.Print(ctable, tb)
	}
}

func TestGetFilelist(t *testing.T) {
	inpath := "/Users/zhengxiangfan/work/周补通/"
	fileSets := GetFilelist(inpath, ".sql")
	log.Println(fileSets)
}


func TestHQL_ParseDepends(t *testing.T) {
	sql := `
	drop table belle_sh.hg_shanghai_19_sale_inv_last_week_temp;
create table belle_sh.hg_shanghai_19_sale_inv_last_week_temp as
select *
from belle_sh.hg_shanghai_19_sale_inv;

drop table belle_sh.lrs_st_open_store_on_UI;
create table belle_sh.lrs_st_open_store_on_UI as
select *
from  belle_sh.hg_shanghai_19_sale_inv_last_week_temp a inner join belle_sh.hg_shanghai_19_sale_inv_last_week b
on a.id = b.id
`
	hsql := HQL{
		Hql:sql,
	}
	hsql.TrimSql()
	hsql.ParseDepends()
}
