package utils

import (
    "net/http"
    "net/url"
    "fmt"
    "strings"
    "sort"
    "io/ioutil"
    "github.com/roydong/gtools"
)

func CallRest(api string, query, data map[string]interface{}) *gtools.Tree {
    if query != nil {
        api = api + "?" + BuildHttpQuery(query)
    }

    var resp *http.Response
    var err error
    if data == nil {
         resp, err = http.Get(api)
    } else {
        form := url.Values{}
        for k, v := range data {
            form.Set(k, fmt.Sprintf("%v", v))
        }
        resp, err = http.PostForm(api, form)
    }
    if err != nil {
        return nil
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil
    }

    tree := gtools.NewTree()
    if err = tree.LoadJson("", body, false); err != nil {
        return nil
    }

    return tree
}

func BuildHttpQuery(data map[string]interface{}) string {
    query := make([]string, 0, len(data))
    for k, v := range data {
        query = append(query, fmt.Sprintf("%s=%v", k, v))
    }
    return strings.Join(query, "&")
}

func CreateSignature(params map[string]interface{}, skey string) string {
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }

    sort.Strings(keys)

    sigs := make([]string, 0, len(keys))
    for _, k := range keys {
        sigs = append(sigs, fmt.Sprintf("%s=%v", k, params[k]))
    }

    sigs = append(sigs, "secret_key=" + skey)

    return gtools.MD5(strings.Join(sigs, "&"))
}

