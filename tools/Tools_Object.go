package tools

import (
    "bufio"
    "crypto/rand"
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    "net/http"
    "os"
    "regexp"
    "StarRocksDevops/util"
    "strconv"
    "strings"
    "time"
)

const (
    KB = 1024
    MB = 1024 * KB
    GB = 1024 * MB
    TB = 1024 * GB
    PB = 1024 * TB
)

func Post(method, u string, body io.Reader) []byte {
    request, err := http.NewRequest(method, u, body)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }
    request.Header.Set("Content-Type", "application/json;charset=utf-8")
    client := &http.Client{
        Timeout:   time.Second * 30,
        Transport: &http.Transport{},
    }
    respone, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }
    defer respone.Body.Close()
    b, err := ioutil.ReadAll(respone.Body)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }
    return b
}

func Writefile(fname, msg string) {
    fileHandle, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer fileHandle.Close()
    // NewWriter 默认缓冲区大小是 4096
    // 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
    buf := bufio.NewWriterSize(fileHandle, len(msg))

    buf.WriteString(msg)

    err = buf.Flush()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
}

func GetHour(second int) (formatString string) {
    // GetHour 秒格式化
    hours := second / 3600
    minutes := (second % 3600) / 60
    secs := second % 60

    if hours >= 1 {
        return fmt.Sprintf("%02dh:%02dmin:%02ds", hours, minutes, secs)
    }
    if minutes >= 1 {
        return fmt.Sprintf("%02dmin:%02ds", minutes, secs)
    }
    return fmt.Sprintf("%02ds", secs)
}

// DuplicatesAndCount /*数组去重*/
func DuplicatesAndCount(strs []string) ([]string, map[string]int) {
    unique := make([]string, 0)
    counts := make(map[string]int)

    for _, str := range strs {
        // 计算每个字符串出现的次数
        counts[str]++
        // 如果是第一次出现，则添加到unique数组中
        if counts[str] == 1 {
            unique = append(unique, str)
        }
    }
    return unique, counts
}

func Version() float64 {
    sql := fmt.Sprintf("select current_version() as version")
    var m map[string]interface{}
    util.Connect.Raw(sql).Scan(&m)

    arr := strings.Split(m["version"].(string), " ")[0]
    if len(arr) < 2 {
        return 0
    }
    if !strings.Contains(arr, ".") {
        return 0
    }

    version, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(arr, ".")[0], strings.Split(arr, ".")[1]), 64)
    if err != nil {
        return 0
    }
    return version
}

// ByteSizeToString 将字节数转换为人类可读的字符串表示形式（如 KB、MB、GB）
func ByteSizeToString(s int64) string {
    size, _ := strconv.ParseFloat(strconv.FormatInt(s, 10), 64)
    if size < KB {
        return fmt.Sprintf("%.2fB", size)
    } else if size < MB {
        return fmt.Sprintf("%.2fKB", size/KB)
    } else if size < GB {
        return fmt.Sprintf("%.2fMB", size/MB)
    } else if size < TB {
        return fmt.Sprintf("%.2fGB", size/GB)
    } else if size < PB {
        return fmt.Sprintf("%.2fTB", size/TB)
    } else {
        return fmt.Sprintf("%.2fPB", size/PB)
    }
}

// PrintProgress 用于在一行内打印进度条
func PrintProgress(current, total int) {
    percent := current * 100 / total                                      // 计算进度百分比
    length := 20                                                          // 设定进度条的长度
    s := strings.Repeat("■", current*length/total)                        // 根据完成度填充进度条
    e := strings.Repeat("□", length-len(s)/3)                             // 填充剩余的空格(这里一个□占用了3个字符，所以除以3)
    bar := s + e                                                          // 拼接
    fmt.Printf("\033[2K\r%d%% [%s](%d/%d)", percent, bar, current, total) // 使用ANSI转义序列将光标移动到行首
}

// RmDupSlice /*数组去重*/
func RmDupSlice(strs []string) []string {
    result := []string{}
    tempMap := map[string]byte{} // 存放不重复字符串
    for _, e := range strs {
        l := len(tempMap)
        tempMap[e] = 0
        if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
            result = append(result, e)
        }
    }
    return result
}

func SchemaRegexp(input string) ([]string, error) {
    // 创建正则表达式对象，匹配三个点分隔的模式
    //re := regexp.MustCompile(`([^\s.]+)\.([^\s.]+)\.([^\s.]+)`)
    re := regexp.MustCompile(`([a-zA-Z][^\s=,'.]+)\.([^\s=,'.]+)`)
    // 使用FindAllString查找所有匹配项
    //matches := re.FindAllString(input, -1)
    var result []string
    for _, s := range re.FindAllString(input, -1) {
        b := regexp.MustCompile(`[\\/\(\),:|+><~!@#%^&*='";?-]`).FindString(s) != ""
        if !b {
            result = append(result, s)
        }
    }
    return result, nil
}

// RandomPassWord 根据指定长度的生产高敏感度字符串
func RandomPassWord(n int) string {
    allowedChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

    b := make([]byte, n)
    for i := range b {
        // 生成一个随机索引
        ri, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
        if err != nil {
            return ""
        }
        // 使用随机索引获取一个字符
        b[i] = allowedChars[ri.Int64()]
    }

    return string(b)
}

func VToFloat(version string) float64 {
    // Split the string by '-' and take the first part
    parts := strings.Split(version, "-")
    mainPart := parts[0]
    // Split the main part by '.'
    subParts := strings.Split(mainPart, ".")
    // Construct the new version format as float
    var newVersionFloat float64
    var err error
    if len(subParts) >= 2 {
        // Combine the first two parts and convert to float
        combined := subParts[0] + "." + subParts[1]
        if len(subParts[2]) > 1 {
            combined += subParts[2][:2] // Append the second character of the third part if it exists
        } else {
            combined += subParts[2]
        }
        newVersionFloat, err = strconv.ParseFloat(combined, 64)
        if err != nil {
            fmt.Println(err.Error())
            return 0
        }
    }
    return newVersionFloat
}
