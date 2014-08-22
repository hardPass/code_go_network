package main

import (
	"bufio"
	"bytes"
	// "encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

const tpl_insert_user = `INSERT INTO user (id,org_id,passwd,cellphone,email,reg_date,gender,position,nickname,valid) 
VALUES ({{.Id}},{{.Org_id}},'{{.Passwd}}','{{.Cellphone}}','{{.Email}}','{{.Reg_date}}','{{.Gender}}','{{.Position}}','{{.Nickname}}','{{.Valid}}');`

const tpl_insert_org = `INSERT INTO org (id,pid,pname,name,representative,phone,fax,reg_date,address,valid)
VALUES ({{.Id}},{{.Pid}},'{{.Pname}}','{{.Name}}','{{.Representative}}','{{.Phone}}','{{.Fax}}','{{.Reg_date}}','{{.Address}}','{{.Valid}}');`

var (
	orglist             []*Org
	userlist            []*User
	org_AUTO_INCREMENT  = 80000000
	user_AUTO_INCREMENT = 80000000
	mapCellphone        = make(map[string]interface{})
)

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// ##单  位,姓  名,性别,职  务, 电  话,传  真, 手  机,MSN/QQ,业务电话,地  址,
type Org struct {
	Id             int
	Pid            int
	Pname          string
	Name           string // 单  位
	Representative string // 姓  名
	Reg_date       string
	Valid          string
	Phone          string // 电  话
	Fax            string // 传  真
	Address        string // 地  址
}

func (this *Org) InsertSql() string {
	t := template.New("any-template")
	t, err := t.Parse(tpl_insert_org)
	checkError(err)

	buffer := bytes.NewBuffer(make([]byte, 0))

	err = t.Execute(buffer, this)
	checkError(err)

	sql := (string)(buffer.Bytes())
	// log.Println(sql)

	return sql
}

type User struct {
	Id       int
	Org_id   int
	Passwd   string
	Email    string
	Reg_date string
	Valid    string

	Nickname  string // 姓  名
	Gender    string // 性别
	Position  string // 职  务
	Cellphone string // 手  机
}

func (this *User) InsertSql() string {
	t := template.New("any-template")
	t, err := t.Parse(tpl_insert_user)
	checkError(err)

	buffer := bytes.NewBuffer(make([]byte, 0))

	err = t.Execute(buffer, this)
	checkError(err)

	sql := (string)(buffer.Bytes())
	// log.Println(sql)

	return sql
}

func handleLine(line string) {
	fields := strings.Split(line, ",")

	org_AUTO_INCREMENT++
	org := new(Org)
	org.Id = org_AUTO_INCREMENT
	org.Pid = 0
	org.Pname = fields[0]
	org.Name = fields[0]
	org.Representative = fields[1]
	org.Reg_date = "1407312335276"
	org.Valid = "0"
	org.Phone = fields[4]
	org.Fax = fields[5]
	org.Address = fields[9]
	orglist = append(orglist, org)

	user_AUTO_INCREMENT++
	user := new(User)
	user.Id = user_AUTO_INCREMENT
	user.Org_id = org.Id
	user.Passwd = "d1285815febb9781f709d12d6c821230"
	user.Reg_date = "1407312335276"
	user.Valid = "0"
	user.Nickname = fields[1]

	user.Gender = "0"
	if fields[2] == "男" {
		user.Gender = "1"
	}

	user.Position = fields[3]
	user.Cellphone = fields[6]
	user.Email = user.Cellphone

	if mapCellphone[user.Cellphone] == nil {
		userlist = append(userlist, user)
	} else {
		log.Printf("[WARN]this cellpone is duplicated: %s.", user.Cellphone)
		userlist = append(userlist, nil)
	}
	mapCellphone[user.Cellphone] = &Holder{}
}

type Holder struct{}

func gen(r *bufio.Reader) {
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if strings.HasPrefix(line, "##") {
			continue
		}

		handleLine(line)

	}
	return
}

func main() {
	from, _ := os.OpenFile("sp.txt", os.O_RDONLY, 0660)
	defer func() {
		from.Close()
	}()

	gen(bufio.NewReader(from))

	to, err := os.OpenFile("insert.sql", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		from.Close()
	}()

	buffer := bytes.NewBuffer(make([]byte, 0, 100000))

	for i := 0; i < len(orglist); i++ {
		org := orglist[i]
		user := userlist[i]

		// bOrg, _ := json.Marshal(org)
		// bUser, _ := json.Marshal(user)
		// log.Printf("org:%s\n user:%s\n", bOrg, bUser)

		buffer.WriteString(org.InsertSql() + "\n")

		if user != nil {
			buffer.WriteString(user.InsertSql() + "\n")
		}

	}

	to.WriteString(buffer.String())
}
