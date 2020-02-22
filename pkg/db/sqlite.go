package db

import (
	"strings"

	"github.com/nektro/mantle/pkg/iconst"

	"github.com/nektro/go-util/util"

	. "github.com/nektro/go-util/alias"
)

func QueryUserByUUID(uid string) (*User, bool) {
	rows := DB.Build().Se("*").Fr(iconst.TableUsers).Wh("uuid", uid).Exe()
	if !rows.Next() {
		return &User{}, false
	}
	ru := User{}.Scan(rows).(*User)
	rows.Close()
	return ru, true
}

func QueryUserBySnowflake(provider string, flake string, name string) *User {
	rows := DB.Build().Se("*").Fr(iconst.TableUsers).Wh("provider", provider).Wh("snowflake", flake).Exe()
	if rows.Next() {
		ru := User{}.Scan(rows).(*User)
		rows.Close()
		return ru
	}
	// else
	id := DB.QueryNextID(iconst.TableUsers)
	uid := newUUID()
	now := T()
	roles := ""
	if id == 1 {
		roles += "o"
		Props.Set("owner", uid)
	}
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%s', '%s', '0', '0', ?, '', '%s', '%s', '%s')", iconst.TableUsers, id, provider, flake, uid, now, now, roles), name)
	return QueryUserBySnowflake(provider, flake, name)
}

func QueryAssertUserName(uid string, name string) {
	DB.Build().Up(iconst.TableUsers, "name", name).Wh("uuid", uid).Exe()
}

func CreateRole(name string) string {
	id := DB.QueryNextID(iconst.TableRoles)
	uid := newUUID()
	util.Log("[role-create]", uid, name)
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '', 1, 1)", iconst.TableRoles, id, uid, id), name)
	return uid
}

func CreateChannel(name string) string {
	id := DB.QueryNextID(iconst.TableChannels)
	uid := newUUID()
	util.Log("[channel-create]", uid, "#"+name)
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '')", iconst.TableChannels, id, uid, id), name)
	AssertChannelMessagesTableExists(uid)
	return uid
}

func AssertChannelMessagesTableExists(uid string) {
	DB.CreateTable(F("%s%s", iconst.TableMessagesPrefix, strings.Replace(uid, "-", "_", -1)), []string{"id", "int primary key"}, [][]string{
		{"uuid", "text"},
		{"sent_at", "text"},
		{"sent_by", "text"},
		{"text", "text"},
		{"test", "text"},
	})
}
