package gotextcat

/*
implement the classification technique described in
Cavnar & Trenkle, "N-Gram-Based Text Categorization".
http://www.novodynamics.com/trenkle/papers/sdair-94-bc.ps.gz
*/

/*
N-Grams:
"TEXT" -> "_TEXT_"
bi-grams:   _T, TE, EX, XT, T_
tri-grams:  _TE, TEX, EXT, XT_  // T__
quad_grams: _TEX, TEXT, EXT_    // XT__, T___
*/


import (
    "os"
    "log"
    "fmt"
    "path"
    "sort"
    "bufio"
    "errors"
    "strings"
    "strconv"
    "path/filepath"
)


//copy from libtextcat-2.2/src/constants.h
//
///* Reported matches are those fingerprints with a score less than best
// * score * THRESHOLDVALUE (i.e. a THRESHOLDVALUE of 1.03 means matches
// * must score within 3% from the best score.)  
// */
//#define THRESHOLDVALUE  1.03
//
///* If more than MAXCANDIDATES matches are found, the classifier reports
// * unknown, because the input is obviously confusing.
// */
//#define MAXCANDIDATES   5
//
///* Maximum number of n-grams in a fingerprint */
//#define MAXNGRAMS  400
//
///* Maximum size of an n-gram? */
//#define MAXNGRAMSIZE 5
//
///* Which characters are not acceptable in n-grams? */
//#define INVALID(c) (isspace((int)c) || isdigit((int)c)) 
//
///* Minimum size (in characters) for accepting a document */
//#define MINDOCSIZE  25
//
const (
    gramN int = 5
    lMax  int = 400
    tMax  int = 400 * 400 * 100
    dSize int = 25      // min doc size
)


type fpItem struct {
    str string
    cnt int
}


type sortedItems struct {
    items []fpItem
}


func (si *sortedItems)Len() int {
    return len(si.items)
}

func (si *sortedItems)Less(i, j int) bool {
    // we need a desc sorted by cnt, asc sorted by name
    if si.items[i].cnt == si.items[j].cnt {
        return si.items[i].str < si.items[j].str
    }
    return si.items[i].cnt > si.items[j].cnt
}


func (si *sortedItems)Swap(i, j int) {
    si.items[i], si.items[j] = si.items[j], si.items[i]
}


type fingerPrint struct {
    sortedItems
}


// split a string, by any byte in seps, equal
// regexp.Compile("[seps]+").Split(strings.Trim(src, seps), -1)
func splitByByte(src string, seps []byte) []string {
    mp := make([]byte, 256)
    for i := range seps {
        mp[seps[i]] = '1'
    }
    ret := make([]string, 0, len(src) / 3)
    b, i := 0, 0
    for ; i < len(src); i++ {
        if mp[src[i]] == '1' {
            if b < i {
                ret = append(ret, src[b:i])
            }
            b = i + 1
        }
    }
    if b < i {ret = append(ret, src[b:i])}
    return ret
}


func colGram(phase string, dst *map[string]int) {
    p := "_" + strings.ToLower(phase) + "_"
    n := len(p)
    for i := 0; i < n; i++ {
        m := i + gramN
        if m > n {m = n}
        for j := i + 1; j <= m; j++ {
            c := p[i:j]
            (*dst)[c] = (*dst)[c] + 1
        }
    }
}



func getFingerPrint(src string) (fp *fingerPrint, size int) {
    var fpmp = make(map[string]int)
    //for _, phase := range splitByByte(src, []byte("0123456789 \t\n\r\f")) {  // equal python [0-9\s]+
    for _, phase := range splitByByte(src, []byte("0123456789 \t\n\r\f\v")) {  // isspace||isdigit
        if len(phase) > 0 {
            size = size + len(phase) + 1    // 1 for space
            colGram(phase, &fpmp)
        }
    }
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

// Calculate gived string finger print, print top 400 items
func PrintFingerPrint(src string) {
    fp, _ := getFingerPrint(src)
    for _, it := range fp.items {
        fmt.Println(it.str, it.cnt)
    }
}


type LangInfo struct {
    id int
    name string
    fmap map[string]int
}


// fn is LMI filename, it has format like  name.id.lm
// name is language name may with codec name, id is language id
func loadLangInfo(fn string) (li *LangInfo, err error) {
    fin, err := os.Open(fn)
    if err != nil {return}
    defer fin.Close()
    bn := filepath.Base(fn)
    dat := strings.Split(bn, ".")
    if len(dat) != 3 || dat[2] != "lm" {
        return nil, errors.New("bad format filename: " + bn)
    }
    li = &LangInfo{fmap:make(map[string]int)}
    i, err := strconv.ParseInt(dat[1], 10, 16)
    if err != nil {
        return nil, err
    }
    li.id, li.name = int(i), dat[0]
    cnt := 0
    scanner := bufio.NewScanner(fin)
    for scanner.Scan() {
        cols := splitByByte(scanner.Text(), []byte(" \r\n\t"))
        li.fmap[strings.TrimSpace(cols[0])] = cnt
        cnt++
        if cnt >= lMax {break}
    }
    return
}


func (fp *fingerPrint)getDistance(dst *LangInfo, cutoff int) (acc int) {
    dis := 0
    for idx, item := range fp.items {
        i, ok := dst.fmap[item.str]
        if ok {
            dis = idx - i
            if dis < 0 { dis = -dis }
            acc = acc + dis
        } else {
            acc = acc + lMax * 100      // libtextcat using lMax
        }
        if acc > cutoff {return tMax}
    }
    return
}


func (li *LangInfo)GetId() int {
    return li.id
}


func (li *LangInfo)GetName() string {
    return li.name
}


var infoList []*LangInfo


func Init(dn string) {
    //m, err := filepath.Glob("/usr/share/gotextcat/data/LMI/*.lm")
    m, err := filepath.Glob(path.Join(dn, "*.lm"))
    if err != nil {
        log.Fatal("lang init", err)
    }
    infoList = make([]*LangInfo, len(m))
    for idx, fn := range m {
        infoList[idx], err = loadLangInfo(fn)
        if err != nil {
            log.Fatal("lang init", fn, err)
        }
    }
}


type ds struct {
    li  *LangInfo
    acc int
}


func GetLanguage(src string) (l1, l2 *LangInfo) {
    fp, size := getFingerPrint(src)
    if size < dSize {return}    // not enough data
    dists := new([5]ds)
    cutoff := tMax + 1
    for _, li := range infoList {
        acc := fp.getDistance(li, cutoff)
        //if acc == tMax {continue}
        if acc > cutoff {continue}
        for i := 0; i < 5 && li != nil; i++ {
            if dists[i].acc == 0 || acc < dists[i].acc {
                li, acc, dists[i].li, dists[i].acc = dists[i].li, dists[i].acc, li, acc
            }
        }
        cutoff = dists[0].acc * 103 / 100 + 1
    }
    if dists[4].acc == 0 || dists[4].acc > cutoff {
        // not all five in 1.03%, so has result
        l1 = dists[0].li
    }
    if dists[1].acc > 0 && dists[1].acc <= cutoff {
        // has secondary possible
        l2 = dists[1].li
    }
    return
}
