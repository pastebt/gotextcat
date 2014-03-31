package gotextcat


import (
    "testing"
    "fmt"
    "regexp"
    "strings"
)


func TestFingerPrint(tst *testing.T) {
    okdata := [][]string {
        {"ab", "_ 2 _a 1 _ab 1 _ab_ 1 a 1 ab 1 ab_ 1 b 1 b_ 1 "},
        {"a b", "_ 4 _a 1 _a_ 1 _b 1 _b_ 1 a 1 a_ 1 b 1 b_ 1 "},
        {"abb", "_ 2 b 2 _a 1 _ab 1 _abb 1 _abb_ 1 a 1 ab 1 abb " +
                "1 abb_ 1 b_ 1 bb 1 bb_ 1 "},
        {"abcd", "_ 2 _a 1 _ab 1 _abc 1 _abcd 1 a 1 ab 1 abc 1 abcd 1 " +
                 "abcd_ 1 b 1 bc 1 bcd 1 bcd_ 1 c 1 cd 1 cd_ 1 d 1 d_ 1 "},
        {"abcde", "_ 2 _a 1 _ab 1 _abc 1 _abcd 1 a 1 ab 1 abc 1 abcd 1 " +
                  "abcde 1 b 1 bc 1 bcd 1 bcde 1 bcde_ 1 c 1 cd 1 cde 1 " +
                  "cde_ 1 d 1 de 1 de_ 1 e 1 e_ 1 "},
        {"测",  "_ 2 _\xe6 1 _\xe6\xb5 1 _\xe6\xb5\x8b 1 _\xe6\xb5\x8b_ 1 " +
                "\x8b 1 \x8b_ 1 \xb5 1 \xb5\x8b 1 \xb5\x8b_ 1 \xe6 1 " +
                "\xe6\xb5 1 \xe6\xb5\x8b 1 \xe6\xb5\x8b_ 1 "},
    }
    for _, dat := range okdata {
        fp, _ := getFingerPrint(dat[0])
        ret := ""
        for _, i := range fp.items {
            ret = ret + fmt.Sprintf("%s %d ", i.str, i.cnt)
        }
        if ret != dat[1] {
            tst.Error(ret, dat)
        }
    }
}


func TestLanguage(tst *testing.T) {
    //dat := []string{"Condiciones específicas de uso de Galeon Centro de ayuda", "spanish"}
    //dat := []string{"este é um teste de sentença Inglês", "portuguese"}
    //dat := []string{"este é um teste de sentença", "portuguese"}
    dat := []string{"este é um teste de", "aa"}     // #define MINDOCSIZE  25
    //dat := []string{"este", ""}
    tst.Log(len(infoList))
    fp, _ := getFingerPrint(dat[0])
    tst.Log(fp)
    imap := make(map[int][]int)
    for _, li := range infoList {
        acc := fp.getDistance(li, tMax)
        imap[li.GetId()] = []int{acc, 0}
    }
    l1, l2 := GetLanguage(dat[0])
    if l1 != nil {
        tst.Log(l1.GetId(), l1.GetName(), imap[l1.GetId()])
    }
    if l2 != nil {
        tst.Log(l2.GetId(), l2.GetName(), imap[l2.GetId()])
    }
}


func TestGetLanguage(tst *testing.T) {
    okdata := [][]string {
        {"this is a english testing sentence", "english"},
        {"esta es una sentencia de pruebas Inglés", "spanish"},
        {"este é um teste de sentença Inglês", "portuguese"},   // if use lMax * 100, will failed
        {"il s'agit d'un test phrase anglais", "french"},
        {"dies ist ein Englisch-Tests Satz", "german"},

        {"تغییر نام و یا به طور موقت از دسترس خارج شده باشد", "arabic-iso8859_6"},
        {"Condiciones específicas de uso de Galeon Centro de ayuda", "spanish"},
        {"Lo sentimos, esta página no existe o no está disponible", "spanish"},
    }
    for _, dat := range okdata {
        l1, l2 := GetLanguage(dat[0])
        if l1 != nil && l1.GetName() != dat[1] {
            tst.Error(dat, l1.GetName(), l2.GetName())
        }
    }
}


func BenchmarkConcatString(bm *testing.B) {
    ss := "hello world"
    m := make(map[string]int)
    for i := 0; i < bm.N; i++ {
        //_ = "_" + ss + "_"
        colGram(ss, &m)
    }
}


func BenchmarkBytes2String(bm *testing.B) {
    b := []byte(" hello world ")
    m := make(map[string]int)
    for i := 0; i < bm.N; i++ {
        //_ = string(b)
        colGram2(b, &m)
    }
}


func strSplit(src string) []string {
    ret := make([]string, 0, 10)
    b, i := 0, 0
    for ; i < len(src); i++ {
        c := src[i]
        if ('0' <= c && c <= '9') || c == ' ' || c == '\t' ||
           c == '\n' || c == '\r' || c == '\v' {
            if b < i {
                ret = append(ret, src[b:i])
            }
            b = i + 1
        }
    }
    if b < i {ret = append(ret, src[b:i])}
    return ret
}


func TestSplit(tst *testing.T) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaa"
    re := regexp.MustCompile(`[0-9\s]+`)
    s1 := strSplit(src)
    tst.Log(s1)
    s2 := splitByByte(src, []byte("0123456789 \t\r\n\v\f"))
    tst.Log(s2)
    s3 := re.Split(src, -1)
    tst.Log(s3)
}


func BenchmarkStringSplit(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaa"
    _ = strings.Split(src, "t")
    for i := 0; i < bm.N; i++ {
        //_ = strings.Split(src, "0123456789 \t\r\n\v\f")
        _ = strSplit(src)
        //_ = strings.Split(src, "t")
    }
}


func BenchmarkSplitByByte(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaa"
    _ = strings.Split(src, "t")
    //seps := []byte("0123456789 \t\r\n\v")
    seps := []byte("0123456789 \t\r\n\v")
    for i := 0; i < bm.N; i++ {
        _ = util.SplitByByte(src, seps)
        //_ = strSplitString(src, " \t\r\n\v")
    }
}


func BenchmarkStringSplitter(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaa"
    _ = strings.Split(src, "t")
    seps := []byte("0123456789 \t\r\n\v")
    sp := util.MakeSplitter(seps)
    for i := 0; i < bm.N; i++ {
        _ = sp.Split(src)
        //_ = util.SplitByByte(src, seps)
        //_ = strSplitString(src, " \t\r\n\v")
    }
}


func BenchmarkReSplit(bm *testing.B) {
    //re := regexp.MustCompile(`\s+`)
    re := regexp.MustCompile(`[0-9\s]+`)
    //re := regexp.MustCompile(`t`)
    src := "test th1is a2   te2t3s4i5n6g7aaaaa"
    for i := 0; i < bm.N; i++ {
        _ = re.Split(src, -1)
    }
}
