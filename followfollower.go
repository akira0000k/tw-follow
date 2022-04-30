package main

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"flag"
	"strconv"
	"io/ioutil"
	"os"
	"os/user"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

// type Cursor struct {
//  	Previous_cursor     int64
//  	Previous_cursor_str string
//  
//  	Ids []int64
//  
//  	Next_cursor     int64
//  	Next_cursor_str string
// }

var api *anaconda.TwitterApi
func main(){
	followerp := flag.Bool("followers", false, "get followers")
	friendp   := flag.Bool("friends", true, "get friends")
	flag.Parse()
	
	api = connectTwitterApi()
	var uv=url.Values{}

	//uv["screen_name"]   = []string{"??????????"} //default= auth id
	uv["cursor"]        = []string{"-1"} //default=from start
	//uv["count"]         = []string{"100"}  //default=5000

	if *followerp {
		fmt.Fprintln(os.Stderr, "get followers")
		followers(uv)
	} else if *friendp {
		fmt.Fprintln(os.Stderr, "get friends")
		friends(uv)
	}
}

func lookup(ids []int64, uv0 url.Values) {
	var uv=url.Values{}
	uv["screen_name"]   = uv0["screen_name"]
	uv["cursor"]        = []string{"-1"}  //default=from start
	uv["count"]         = nil

	total := len(ids)
	each := 99

	start := 0
	looplimit := total/each + 1
	for i:=0; i<looplimit; i++ {
		if i > 0 {
			//fmt.Printf("wait 10 second-->")
			time.Sleep(time.Second * 10)
			//fmt.Printf("wake up<---------")
		}
		fmt.Fprint(os.Stderr, ".")
		var usermap = map[string]anaconda.User{}
		
		end := start + each
		if end > total { end = total }
		partids := ids[start:end]
		//uv.Set("include_entities", "false") //skip_status:falseのときに若干効く
		uv.Set("skip_status", "true") //データ量減少

		// fmt.Println(partids)
		// fmt.Println(uv)
		u, err := api.GetUsersLookupByIds(partids, uv)
		if err != nil {
			fmt.Println("err=", err)
			return
		}
		// if i > 0 {
		//  	fmt.Printf("u=%#v\n", u)
		// }
		for _, user := range u {
			// fmt.Println(user.IdStr)
			//jsonuser, _ := json.Marshal(user) //test print
			//fmt.Println(string(jsonuser)) //test print
			usermap[user.IdStr] = user
		}
		for _, id := range partids {
			idstr := strconv.FormatInt(id, 10)
			// fmt.Println(idstr)
			user, ok := usermap[idstr]
			if !ok { break }
			//fmt.Printf("%s\t%#v\t%#v\n", user.ScreenName, user.Name, user.Description)
			username := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(user.Name,        "\n", `\n`), "\r", `\r`), "\"", `\"`)
			userdesc := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(user.Description, "\n", `\n`), "\r", `\r`), "\"", `\"`)
			fmt.Printf("%s\t\"%s\"\t\"%s\"\n", user.ScreenName, username, userdesc)
		}
		start += each
		if start >= total { break }
	}
}

func followers(uv url.Values) {
	looplimit := 20
	for i:=0; i<looplimit; i++ {
		if i > 0 {
			time.Sleep(time.Second * 5)
		}
		fmt.Fprint(os.Stderr, ".")

		cursor, err := api.GetFollowersIds(uv)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println("next cursor = ", cursor.Next_cursor_str)
		// fmt.Printf("cursor=%#v\n", cursor)
		// fmt.Printf("total ids=%d\n", len(cursor.Ids))
		lookup(cursor.Ids, uv)
		
		nexc := cursor.Next_cursor
		if nexc == 0 { break }

		uv["cursor"] = []string{cursor.Next_cursor_str}
		// fmt.Printf("cursor next=%#v\n", uv)
	}
}

func friends(uv url.Values) {
	looplimit := 20
	for i:=0; i<looplimit; i++ {
		if i > 0 {
			time.Sleep(time.Second * 5)
		}
		fmt.Fprint(os.Stderr, ".")

		cursor, err := api.GetFriendsIds(uv)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println("next cursor = ", cursor.Next_cursor_str)
		// fmt.Printf("cursor=%#v\n", cursor)
		// fmt.Printf("total ids=%d\n", len(cursor.Ids))
		lookup(cursor.Ids, uv)
		
		nexc := cursor.Next_cursor
		if nexc == 0 { break }

		uv["cursor"] = []string{cursor.Next_cursor_str}
		// fmt.Printf("cursor next=%#v\n", uv)
	}
}

func connectTwitterApi() *anaconda.TwitterApi {
	usr, _ := user.Current()
	raw, error := ioutil.ReadFile(usr.HomeDir + "/twitter/twitterAccount.json")

	if error != nil {
		fmt.Println(error.Error())
		return nil
	}

	var twitterAccount TwitterAccount
	json.Unmarshal(raw, &twitterAccount)

	return anaconda.NewTwitterApiWithCredentials(
		twitterAccount.AccessToken,
		twitterAccount.AccessTokenSecret,
		twitterAccount.ConsumerKey,
		twitterAccount.ConsumerSecret)

}

type TwitterAccount struct {
	AccessToken string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
}
