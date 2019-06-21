package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type IParseSql interface {
	SplitSql() []string
	TrimSql()
	ExtraceTableName(string) (string, []string)
}

type HQL struct {
	Hql  string
	RHql string
}

func (h *HQL) SplitSql() []string {
	var result []string
	h.TrimSql()
	res := strings.Split(h.RHql, "drop")
	//m1, _ := regexp.Compile(`drop`)
	for _, table := range res {
		if strings.Index(strings.Trim(table, " "), "table") == 0 {
			result = append(result, fmt.Sprintf("drop %s", table))
		}
	}
	return result
}

func (h *HQL) TrimSql() {
	re2, _ := regexp.Compile(`/\*[^*]*\*+(?:[^*/][^*]*\*+)*/`)
	re3, _ := regexp.Compile(`^\s*(--|#)`)
	re4, _ := regexp.Compile(`--|#`)
	rep := re2.ReplaceAllString(h.Hql, "")
	reps := strings.Split(rep, "\n")
	var result1 []string
	var result2 []string
	for _, line := range reps {
		if !re3.Match([]byte(line)) {
			result1 = append(result1, line)
		}
	}
	for _, line := range result1 {
		ls := re4.Split(line, -1)
		result2 = append(result2, ls[0])
	}
	h.RHql = strings.Join(result2, "\n")
}

func (h *HQL) ExtraceTableName(sqlStr string) (string, []string) {
	var (
		tokens   []string
		result1  []string
		objtable string
		srctable string
	)
	result := make(map[string]int)
	get_next := false
	m1, _ := regexp.Compile(`[\s)(;]+`)
	tokens = m1.Split(sqlStr, -1)
	log.Println(tokens)
	for _, token := range tokens {
		if get_next {
			if !(strings.Contains("select", strings.ToLower(token)) && strings.Contains("", strings.ToLower(token))) {
				dab := strings.Split(token, ".")
				if len(dab) == 2 {
					srctable = dab[1]
				} else {
					srctable = token
				}
				result1 = append(result1, srctable)
				result[srctable] = 0
			}
			get_next = false
		}
		get_next = strings.Contains("from", strings.ToLower(token)) || strings.Contains("join", strings.ToLower(token)) || strings.Contains("table", strings.ToLower(token))
	}
	if len(result1) > 1 {
		objtable = result1[0]
	}
	var r []string
	for k, _ := range result {
		r = append(r, k)
	}
	return objtable, r
}

func GetFilelist(InFilepath, ext string) []string {
	var fileSets []string
	err := filepath.Walk(InFilepath, func(InFilepath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if path.Ext(InFilepath) == ext {
			fileSets = append(fileSets, InFilepath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return fileSets
}

func ReadAll(filePth string) (string, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(f)
	return string(bytes), err
}

func (h *HQL) ParseDepends() ([]map[string]interface{}, []string) {
	var (
		depends_res []map[string]interface{}
		obj_names map[string]int
		depends_all []string
	)
	obj_names = make(map[string]int)
	sql_blocks := h.SplitSql()
	for _, sql_str := range sql_blocks {
		var depends []string
		obj_table_names, src_table_names := h.ExtraceTableName(sql_str)
		for _, tb := range src_table_names {
			if tb != obj_table_names {
				depends = append(depends, fmt.Sprintf("%s >> %s", tb, obj_table_names))
				depends_all = append(depends_all, fmt.Sprintf("%s >> %s", tb, obj_table_names))
			}
		}
		dep := make(map[string]interface{})
		dep["depends"] = depends
		dep["obj_table"] = obj_table_names
		dep["sql"] = sql_str
		obj_names[obj_table_names]=0
		depends_res = append(depends_res, dep)
	}
	var obj_table_sets []string
	for _, tb := range depends_all{
		tbb := strings.Split(tb,">>")
		if  _, ok := obj_names[strings.Trim(tbb[0]," ")]; ok {
			obj_table_sets = append(obj_table_sets, tb)
		}
	}
	log.Println("depends", obj_table_sets)
	return depends_res, obj_table_sets
}