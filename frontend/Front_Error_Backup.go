package frontend

import (
    "fmt"
    "github.com/fatih/color"
    "StarRocksDevops/tools"
    "StarRocksDevops/util"
    "strings"
)

func CheckTabletErrors1(msg string) {
    type TABLET struct {
        DbName        string `bson:"DbName"`
        TableName     string `bson:"TableName"`
        PartitionName string `bson:"PartitionName"`
        IndexName     string `bson:"IndexName"`
        DbId          int    `bson:"DbId"`
        TableId       int    `bson:"TableId"`
        PartitionId   int    `bson:"PartitionId"`
        IndexId       int    `bson:"IndexId"`
        IsSync        bool   `bson:"IsSync"`
        DetailCmd     string `bson:"DetailCmd"`
    }
    var tablet TABLET
    c := color.New()

    /*判断是否副本损坏*/
    if strings.Contains(msg, "tablet_id") {

        var tabletID string
        if tools.Version() > 3.0 {
            fmt.Println(len(strings.Split(msg, "tablet_id: ")))
            fmt.Println(strings.Split(strings.Split(msg, "tablet_id: ")[0], ","))
            tid := strings.Split(strings.Split(msg, "tablet_id: ")[0], ",")[0]
            tabletID = strings.Split(tid, "=")[1]
        } else {
            tabletID = strings.Split(strings.Split(msg, "tablet_id: ")[1], ",")[0]
        }

        fmt.Println(c.Add(color.FgHiYellow).Sprint("修复判断------------------------>") + tabletID)

        r := util.Connect.Raw(fmt.Sprintf("show tablet %s", tabletID)).Scan(&tablet)
        if r.Error != nil {
            fmt.Println(r.Error.Error())
            return
        }
        _, on := Backendids1(tablet.DetailCmd)
        if on {
            sql := fmt.Sprintf(`ADMIN SET REPLICA STATUS PROPERTIES("tablet_id" = "%s", "backend_id" = "%d", "status" = "bad")`, tabletID, util.P.BackendId)
            r := util.Connect.Exec(sql)
            if r.Error != nil {
                fmt.Println(r.Error.Error())
                return
            }
            fmt.Println(c.Add(color.FgHiGreen).Sprint("修复提交------------------------>") + tabletID)
        }
    }
}

func Backendids1(sql string) (int, bool) {
    type PROCS []struct {
        ReplicaId             int    `bson:"ReplicaId"`
        BackendId             int    `bson:"BackendId"`
        Version               int    `bson:"Version"`
        VersionHash           int    `bson:"VersionHash"`
        LstSuccessVersion     int    `bson:"LstSuccessVersion"`
        LstSuccessVersionHash int    `bson:"LstSuccessVersionHash"`
        LstFailedVersion      int    `bson:"LstFailedVersion"`
        LstFailedVersionHash  int    `bson:"LstFailedVersionHash"`
        LstFailedTime         string `bson:"LstFailedTime"`
        SchemaHash            int    `bson:"SchemaHash"`
        DataSize              int    `bson:"DataSize"`
        RowCount              int    `bson:"RowCount"`
        State                 string `bson:"State"`
        IsBad                 bool   `bson:"IsBad"`
        IsSetBadForce         bool   `bson:"IsSetBadForce"`
        VersionCount          string `bson:"VersionCount"`
        PathHash              string `bson:"PathHash"`
        MetaUrl               string `bson:"MetaUrl"`
        CompactionStatus      string `bson:"CompactionStatus"`
    }

    var proc PROCS
    r := util.Connect.Raw(sql).Scan(&proc)
    if r.Error != nil {
        fmt.Println(r.Error.Error())
        return len(proc), false
    }

    for _, s := range proc {
        if s.BackendId == util.P.BackendId && !s.IsBad {
            return len(proc), true
        }
    }
    return len(proc), false
}
