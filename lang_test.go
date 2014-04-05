package gotextcat


import (
    "testing"
    "fmt"
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
    Init("/usr/share/gotextcat/data/LMI")
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
    Init("/usr/share/gotextcat/data/LMI")
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
    Init("/usr/share/gotextcat/data/LMI")
    for _, dat := range okdata {
        l1, l2 := GetLanguage(dat[0])
        if l1 != nil && l1.GetName() != dat[1] {
            tst.Error(dat, l1.GetName(), l2.GetName())
        }
    }
}
