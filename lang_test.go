package gotextcat


import (
    "testing"
    "fmt"
    "sort"
    "strings"
)


func TestFingerPrint(tst *testing.T) {
    testFingerPrint(tst, getFingerPrint)
}


func TestFingerPrint2(tst *testing.T) {
    gfp2 := func (src string) (fp *fingerPrint, size int) {
        return getFingerPrint2([]byte(src))
    }
    testFingerPrint(tst, gfp2)
}


func testFingerPrint(tst *testing.T, gfp func(src string) (fp *fingerPrint, size int)) {
    okdata := [][]string {
        {"a", "_ 2 _a 1 _a_ 1 a 1 a_ 1 "},
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
        //fp, _ := getFingerPrint(dat[0])
        fp, _ := gfp(dat[0])
        ret := ""
        for _, i := range fp.items {
            ret = ret + fmt.Sprintf("%s %d ", i.str, i.cnt)
        }
        if ret != dat[1] {
            tst.Error(ret, dat)
        }
    }
}


func colGram0(phase string, dst *map[string]int) {
    n := len(phase)
    for i := 0; i < n; i++ {
        m := i + gramN
        if m > n {m = n}
        for j := i + 1; j <= m; j++ {
            c := phase[i:j]
            (*dst)[c] = (*dst)[c] + 1
        }
    }
}


func colGram2(phase []byte, dst *map[string]int) {
    n := len(phase)
    for i := 0; i < n; i++ {
        for j := i; j < (i + gramN) && j < n; j++ {
            c := string(phase[i:j + 1])
            (*dst)[c] = (*dst)[c] + 1
        }
    }
}


// for banchmark compare
func xcolGram2(phase []byte, dst *map[string]int) {
    b, e := phase[0], phase[len(phase) - 1]
    phase[0], phase[len(phase) - 1] = '_', '_'
    n := len(phase)
    for i := 0; i < n; i++ {
        for j := i; j < (i + gramN) && j < n; j++ {
            c := string(phase[i:j + 1])
            (*dst)[c] = (*dst)[c] + 1
        }
    }
    phase[0], phase[len(phase) - 1] = b, e
}


func getFingerPrint2(src []byte) (fp *fingerPrint, size int) {
    var fpmp = make(map[string]int)
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    cal := func(pha []byte) {
        //fmt.Printf("phase = [%s]\n", string(pha))
        size += len(pha)
        //colGram2(pha, &fpmp)
        colGram0(string(pha), &fpmp)
    }
    sp.SplitRBM(src, cal)
    si := sortedItems{make([]fpItem, len(fpmp))}
    for k, v := range fpmp {
        si.items = append(si.items, fpItem{k, v})
    }
    // first sort by cnt desc, get top lMax items
    sort.Sort(&si)
    i := len(fpmp)
    if i > lMax { i = lMax }
    fp = &fingerPrint{sortedItems{si.items[:i]}}
    return
}


func BenchmarkFingerPrint(bm *testing.B) {
    src := "Condiciones específicas de uso de Galeon Centro de ayuda"
    for i := 0; i < bm.N; i++ {
        _, _ = getFingerPrint(src)
    }
}


func BenchmarkFingerPrint2(bm *testing.B) {
    src := []byte("Condiciones específicas de uso de Galeon Centro de ayuda")
    for i := 0; i < bm.N; i++ {
        _, _ = getFingerPrint2(src)
    }
}


////////////////////////////////////////////////////////////////////////////////////////////

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


func Benchmarkb2Ss(bm *testing.B) {
    a := ""
    b := []byte("this is a short")
    for i := 0; i < bm.N; i++ {
        a = string(b)
    }
    bm.Log(a)
}


func Benchmarkb2Sl(bm *testing.B) {
    a := ""
    b := []byte("this is a long this is a long this is a long this is a long " +
           "this is a long this is a long this is a long this is a long this " +
           "this is a long this is a long this is a long this is a long this " +
           "this is a long this is a long this is a long this is a long this " +
           "this is a long this is a long this is a long this is a long this ")
    for i := 0; i < bm.N; i++ {
        a = string(b)
    }
    bm.Log(a)
}


func BenchmarkSplitByByte(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " + 
           "o23293i 2;lmel324r34lrkna l"
    _ = strings.Split(src, "t")
    seps := []byte("0123456789 \t\r\n\v")
    for i := 0; i < bm.N; i++ {
        _ = splitByByte(src, seps)
    }
}


//////////////////////////////////////////////////////////////////
func TestGram5(tst *testing.T) {
    //var u uint64
    Gram5t([]byte{1}, tst)
    Gram5t([]byte{1, 2}, tst)
    Gram5t([]byte{1, 2, 3}, tst)
    Gram5t([]byte{1, 2, 3, 4}, tst)
    Gram5t([]byte{1, 2, 3, 4, 5}, tst)
    Gram5t([]byte{1, 2, 3, 4, 5, 6}, tst)
}


func Gram5t(src []byte, tst *testing.T) uint64 {
    und := uint64('_')
    u64 := und
    tst.Log("-------------------")
    var i, j int
    for j = 0; j < len(src); j++ {
        for i = 0; i < 4 && ((i + j) < len(src)); i++ {
            tst.Logf("-%x", u64)
            u64 = (u64 << 8) | uint64(src[i + j])
        }
        if i < 4 {
            tst.Logf("=%x", u64)
            u64 = (u64 << 8) | und
        }
        tst.Logf("+%x", u64)
        u64 = uint64(src[j])
    }
    tst.Logf("*%x", u64)
    tst.Logf("x%x", (u64 << 8) + und)
    tst.Logf("y%x", und)
    return u64
}


func Gram5(src []byte, dst *map[uint64]int) {
    und := uint64('_')
    u64 := und
    //tst.Log("-------------------")
    var i, j int
    for j = 0; j < len(src); j++ {
        for i = 0; i < 4 && ((i + j) < len(src)); i++ {
            //tst.Logf("-%x", u64)
            (*dst)[u64] += 1
            u64 = (u64 << 8) | uint64(src[i + j])
        }
        if i < 4 {
            //tst.Logf("=%x", u64)
            (*dst)[u64] += 1
            u64 = (u64 << 8) | und
        }
        //tst.Logf("+%x", u64)
        (*dst)[u64] += 1
        u64 = uint64(src[j])
    }
    //tst.Logf("*%x", u64)
    (*dst)[u64] += 1
    //tst.Logf("x%x", (u64 << 8) + und)
    (*dst)[(u64 << 8) + und] += 1
    //tst.Logf("y%x", und)
    (*dst)[und] += 1
}


func BenchmarkGram5(bm *testing.B) {
    src := []byte("te2t3s4i5n6g7aaaaawmd")
    dst := make(map[uint64]int)
    for i := 0; i < bm.N; i++ {
        Gram5(src, &dst)
    }
}

func BenchmarkColGram2(bm *testing.B) {
    src := []byte("te2t3s4i5n6g7aaaaawmd")
    dst := make(map[string]int)
    for i := 0; i < bm.N; i++ {
        colGram2(src, &dst)
    }
}

func BenchmarkColGram(bm *testing.B) {
    src := "te2t3s4i5n6g7aaaaawmd"
    dst := make(map[string]int)
    for i := 0; i < bm.N; i++ {
        colGram(src, &dst)
    }
}


//////////////////////////////////////////////////////////////////////////////
type Splitter struct {
    mp []byte
}


func MakeSplitter(seps []byte) *Splitter {
    sp := &Splitter{make([]byte, 256)}
    for i := range seps {
        sp.mp[seps[i]] = '1'
    }
    return sp
}

/*
func (sp *Splitter)CalGram5(src []byte, dst map[uint64]uint32) {
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i { Gram5(src[b:i], dst) }
            b = i + 1
        }
    }
    if b < i { Gram5(src[b:i], dst) }
}
*/

func (sp *Splitter)Split(src string) []string {
    ret := make([]string, 0, len(src) / 3)
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {ret = append(ret, src[b:i])}
            b = i + 1
        }
    }
    if b < i {ret = append(ret, src[b:i])}
    return ret
}


func (sp *Splitter)SplitRight(src string) []string {
    ret := make([]string, 0, len(src) / 3)
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {ret = append(ret, "_" + src[b:i] + "_")}
            b = i + 1
        }
    }
    if b < i {ret = append(ret, "_" + src[b:i] + "_")}
    return ret
}
func (sp *Splitter)SplitRightB(src []byte) [][]byte {
    ret := make([][]byte, 0, len(src) / 3)
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {
                src[i] = '_'
                if b == 0 {
                    ret = append(ret, src[b:i + 1])
                } else {
                    src[b - 1] = '_'
                    ret = append(ret, src[b -1:i + 1])
                }
            }
            b = i + 1
        }
    }
    if b < i {
        if b == 0 {
            ret = append(ret, src[b:i])
        } else {
            src[b - 1] = '_'
            ret = append(ret, src[b - 1:i])
        }
    }
    return ret
}
func (sp *Splitter)SplitRBM(src []byte, dst func(phase []byte)) {
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {
                src[i] = '_'
                if b == 0 {
                    //ret = append(ret, src[b:i + 1])
                    //dst(src[b:i + 1])
                    dst(append([]byte("_"), src[b:i + 1]...))
                } else {
                    src[b - 1] = '_'
                    //ret = append(ret, src[b -1:i + 1])
                    dst(src[b - 1:i + 1])
                }
            }
            b = i + 1
        }
    }
    if b < i {
        s := src[b:i]
        if b == 0 {
            s = append([]byte("_"), s...)
        } else {
            s = src[b - 1:i]
            s[0] = '_'
        }
        if i == len(src) {
            s = append(s, '_')
        }
        dst(s)
    }
}


func TestSplitRBM(tst *testing.T) {
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    dst := func(src []byte) {
        tst.Log(string(src))
    }
    sp.SplitRBM([]byte("a"), dst)
    sp.SplitRBM([]byte("a "), dst)
    sp.SplitRBM([]byte(" a"), dst)
    sp.SplitRBM([]byte("ab"), dst)
    sp.SplitRBM([]byte(" ab"), dst)
    sp.SplitRBM([]byte("ab "), dst)
    sp.SplitRBM([]byte("a b"), dst)
    sp.SplitRBM([]byte("abc"), dst)
    sp.SplitRBM([]byte("abcd"), dst)
    sp.SplitRBM([]byte("ab cd"), dst)
    sp.SplitRBM([]byte("abcde"), dst)
    sp.SplitRBM([]byte("abc  de"), dst)
}

func (sp *Splitter)fakeSplit(src string) (cnt int) {
    b, i, mp := 0, 0, sp.mp
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {cnt += 1}
            b = i + 1
        }
    }
    if b < i {cnt += 1}
    return cnt
}


func BenchmarkSplitterRight(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l"
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    //bm.Log(sp.SplitRight2(src))
    for i := 0; i < bm.N; i++ {
        _ = sp.SplitRight(src)
    }
}
func BenchmarkSplitterRightB(bm *testing.B) {
    src := []byte("test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l")
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    //bm.Log(sp.SplitRightB(src))
    for i := 0; i < bm.N; i++ {
        _ = sp.SplitRightB(src)
    }
}
func BenchmarkSplitterRightBs(bm *testing.B) {
    src := []byte("test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l")
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    //bm.Log(sp.SplitRightB(src))
    for i := 0; i < bm.N; i++ {
        bs := sp.SplitRightB(src)
        ret := make([]string, len(bs))
        for i, b := range bs {
            ret[i] = string(b)
        }
    }
}
func BenchmarkSplitterRightBl(bm *testing.B) {
    src := []byte("test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l")
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    //bm.Log(sp.SplitRightB(src))
    for i := 0; i < bm.N; i++ {
        bs := sp.SplitRightB(src)
        ret := make([]uint32, len(bs))
        for i, b := range bs {
            for _, bt := range b {
                ret[i] += (ret[i] << 8) + uint32(bt)
            }
        }
    }
}
func BenchmarkSplitter(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l"
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    for i := 0; i < bm.N; i++ {
        _ = sp.Split(src)
    }
}
func BenchmarkSplitterFake(bm *testing.B) {
    src := "test th1is a2   te2t3s4i5n6g7aaaaawmd lmwlmd " +
           "o23293i 2;lmel324r34lrkna l"
    seps := []byte("0123456789 \t\r\n\v")
    sp := MakeSplitter(seps)
    for i := 0; i < bm.N; i++ {
        _ = sp.fakeSplit(src)
    }
}

