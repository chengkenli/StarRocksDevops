/*
 *@author  chengkenli
 *@project StarRocksDDL
 *@package behavioral
 *@file    behavioral_ddl_drop
 *@date    2024/12/5 17:07
 */

package recovers

import (
	"fmt"
	"regexp"
	"starrocks/tools"
	"strings"
)

func ddldrop(m2 map[string]interface{}) (bool, string) {
	sql := strings.ToLower(m2["stmt"].(string))
	if !strings.Contains(strings.ToLower(sql), "drop") {
		return false, ""
	}
	var tablename string
	//匹配drop [action] [schema].[tablename]
	regix1 := regexp.MustCompile(fmt.Sprintf(`\s*drop\s+%s\s+([^\s;]+)\.([^\s;]+\S+)`, "table"))
	matched := regix1.MatchString(sql)
	if matched {
		schema, _ := tools.SchemaRegexp(strings.ReplaceAll(sql, "`", ""))
		if schema == nil {
			return false, ""
		}
		tablename = schema[0]
		//命中
		return true, schema[0]
	}

	//匹配drop [action] [tablename]
	regix2 := regexp.MustCompile(fmt.Sprintf(`(?i)drop\s+%s\s*(?:if\s+exists\s*)?([^;\s]+\S+)`, "table"))
	matched2 := regix2.FindStringSubmatch(sql)
	if matched2 != nil && len(matched2) > 1 {
		if !strings.Contains(matched2[1], ".") {
			tablename = fmt.Sprintf("%s.%s", m2["db"].(string), strings.NewReplacer("`", "", "#", "").Replace(matched2[1]))
			//命中
			return true, tablename
		} else {
			return true, matched2[1]
		}
	}

	re := regexp.MustCompile(`(?i)drop table if exists\s+(\w+\.\w+.\w+.\w+.\w+\S+)`)
	matches := re.FindStringSubmatch(sql)
	if matches != nil && len(matches) > 1 {
		if !strings.Contains(matches[1], ".") {
			tablename = fmt.Sprintf("%s.%s", m2["db"].(string), strings.NewReplacer("`", "", "#", "").Replace(matches[1]))
			//命中
			return true, tablename
		} else {
			return true, matches[1]
		}
	}

	re = regexp.MustCompile(`(?i)alter table\s+(\S+)\s+drop partition if exists\s+(\S+)`)
	matches = re.FindStringSubmatch(sql)
	if matches != nil && len(matches) > 1 {
		if !strings.Contains(matches[1], ".") {
			tablename = fmt.Sprintf("%s.%s^%s", m2["db"].(string), strings.NewReplacer("`", "", "#", "").Replace(matches[1]), matches[2])
			//命中
			return true, tablename
		} else {
			return true, fmt.Sprintf("%s^%s", matches[1], matches[2])
		}
	}
	return false, ""
}
