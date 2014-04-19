package conf

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/data/rdc/setup"
    "time"
)

func (c *Config) DatabaseInstance() (*rdc.Conn, error) {

    rdccfg := rdc.NewConfig()
    rdccfg.DbPath = c.DatabasePath
    rdccfg.Driver = "sqlite3"

    //
    cn, _ := rdccfg.Instance()

    //
    tbl_lgn := setup.NewTable("ids_login")
    tbl_lgn.FieldAdd("uid", "auto", 0, 0)
    tbl_lgn.FieldAdd("uname", "string", 20, setup.FieldIndexUnique)
    tbl_lgn.FieldAdd("email", "string", 50, setup.FieldIndexUnique)
    tbl_lgn.FieldAdd("name", "string", 50, 0)
    tbl_lgn.FieldAdd("pass", "string", 100, 0)
    tbl_lgn.FieldAdd("group", "string", 200, 0)
    tbl_lgn.FieldAdd("status", "int16", 0, setup.FieldIndexIndex)
    tbl_lgn.FieldAdd("created", "datetime", 0, setup.FieldIndexIndex)
    tbl_lgn.FieldAdd("updated", "datetime", 0, setup.FieldIndexIndex)

    //
    tbl_prf := setup.NewTable("ids_profile")
    tbl_prf.FieldAdd("uid", "pk", 0, 0)
    tbl_prf.FieldAdd("gender", "int8", 0, 0)
    tbl_prf.FieldAdd("birthday", "date", 0, 0)
    tbl_prf.FieldAdd("address", "string", 100, 0)
    tbl_prf.FieldAdd("url_personal", "string", 100, 0)
    tbl_prf.FieldAdd("aboutme", "string-text", 0, 0)
    tbl_prf.FieldAdd("photo", "string-text", 0, 0)
    tbl_prf.FieldAdd("created", "datetime", 0, 0)
    tbl_prf.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_rst := setup.NewTable("ids_resetpass")
    tbl_rst.FieldAdd("id", "pk-string", 24, 0)
    tbl_rst.FieldAdd("status", "int16", 0, 0)
    tbl_rst.FieldAdd("email", "string", 100, 0)
    tbl_rst.FieldAdd("expired", "datetime", 0, setup.FieldIndexIndex)

    // group
    tbl_grp := setup.NewTable("ids_group")
    tbl_grp.FieldAdd("gid", "auto", 0, 0)
    tbl_grp.FieldAdd("pid", "int32", 0, setup.FieldIndexIndex)
    tbl_grp.FieldAdd("name", "string", 50, 0)
    tbl_grp.FieldAdd("summary", "string", 100, 0)
    tbl_grp.FieldAdd("status", "int16", 0, setup.FieldIndexIndex)
    tbl_grp.FieldAdd("created", "datetime", 0, 0)
    tbl_grp.FieldAdd("updated", "datetime", 0, 0)

    // group_users
    tbl_gpu := setup.NewTable("ids_group_users")
    tbl_gpu.FieldAdd("gukey", "pk-string", 20, 0)
    tbl_gpu.FieldAdd("created", "datetime", 0, 0)

    //
    tbl_app := setup.NewTable("ids_apps")
    tbl_app.FieldAdd("id", "auto", 0, 0)
    tbl_app.FieldAdd("title", "string", 50, 0)
    tbl_app.FieldAdd("created", "datetime", 0, 0)
    tbl_app.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_ins := setup.NewTable("ids_instances")
    tbl_ins.FieldAdd("id", "pk-string", 8, 0)
    tbl_ins.FieldAdd("appid", "string", 50, setup.FieldIndexIndex)
    tbl_ins.FieldAdd("version", "string", 50, 0)
    tbl_ins.FieldAdd("uid", "uint32", 10, 0)
    tbl_ins.FieldAdd("created", "datetime", 0, 0)
    tbl_ins.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_rol := setup.NewTable("ids_roles")
    tbl_rol.FieldAdd("id", "auto", 0, 0)
    tbl_rol.FieldAdd("name", "string", 30, 0)
    tbl_rol.FieldAdd("weight", "int32", 0, 0)

    //
    tbl_pem := setup.NewTable("ids_perms")
    tbl_pem.FieldAdd("id", "auto", 0, 0)
    tbl_pem.FieldAdd("rid", "uint32", 0, setup.FieldIndexIndex)
    tbl_pem.FieldAdd("instance", "string", 30, setup.FieldIndexIndex)
    tbl_pem.FieldAdd("permission", "string", 100, 0)

    //
    tbl_mes := setup.NewTable("ids_menus")
    tbl_mes.FieldAdd("id", "auto", 0, 0)
    tbl_mes.FieldAdd("pid", "uint32", 0, 0)
    tbl_mes.FieldAdd("type", "uint16", 0, setup.FieldIndexIndex)
    tbl_mes.FieldAdd("status", "int16", 0, 0)
    tbl_mes.FieldAdd("instance", "string", 50, setup.FieldIndexIndex)
    tbl_mes.FieldAdd("uid", "uint32", 0, setup.FieldIndexIndex)
    tbl_mes.FieldAdd("title", "string", 100, 0)
    tbl_mes.FieldAdd("link", "string", 100, 0)
    tbl_mes.FieldAdd("weight", "int8", 0, 0)
    tbl_mes.FieldAdd("permission", "string", 50, 0)
    tbl_mes.FieldAdd("created", "datetime", 0, 0)
    tbl_mes.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_ses := setup.NewTable("ids_sessions")
    tbl_ses.FieldAdd("token", "pk-string", 24, 0)
    tbl_ses.FieldAdd("refresh", "string", 24, 0)
    tbl_ses.FieldAdd("status", "int16", 0, setup.FieldIndexIndex)
    tbl_ses.FieldAdd("uid", "uint32", 0, setup.FieldIndexIndex)
    tbl_ses.FieldAdd("name", "string", 50, 0)
    tbl_ses.FieldAdd("uname", "string", 30, 0)
    tbl_ses.FieldAdd("source", "string", 20, 0)
    tbl_ses.FieldAdd("data", "string-text", 0, 0)
    tbl_ses.FieldAdd("permission", "int8", 0, 0)
    tbl_ses.FieldAdd("created", "datetime", 0, 0)
    tbl_ses.FieldAdd("expired", "datetime", 0, 0)

    //
    ds := setup.NewDataSet()
    ds.Version = 7
    // accounts
    ds.TableAdd(tbl_lgn)
    ds.TableAdd(tbl_prf)
    ds.TableAdd(tbl_rst)
    // group
    ds.TableAdd(tbl_grp)
    ds.TableAdd(tbl_gpu)
    // applications
    ds.TableAdd(tbl_app)
    ds.TableAdd(tbl_ins)
    // roles
    ds.TableAdd(tbl_rol)
    ds.TableAdd(tbl_pem)
    ds.TableAdd(tbl_mes)
    // session
    ds.TableAdd(tbl_ses)

    //
    _ = cn.Setup("", ds)

    timenow := time.Now().Format(time.RFC3339)
    _, err := cn.ExecRaw("INSERT OR IGNORE INTO `ids_group` "+
        "(gid,pid,name,summary,status,created,updated) "+
        "VALUES (1,0,\"Administrator\",\"Root System Administrator\",1,?,?),"+
        "(100,0,\"Member\",\"Universal Member\",1,?,?)",
        timenow, timenow, timenow, timenow)
    if err != nil {
        return cn, err
    }

    return cn, nil
}
