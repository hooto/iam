package conf

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/data/rdc/setup"
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
    tbl_lgn.FieldAdd("roles", "string", 200, 0)
    tbl_lgn.FieldAdd("timezone", "string", 40, 0)
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
    /*
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
    */
    //
    /*
       tbl_app := setup.NewTable("ids_apps")
       tbl_app.FieldAdd("id", "auto", 0, 0)
       tbl_app.FieldAdd("title", "string", 50, 0)
       tbl_app.FieldAdd("created", "datetime", 0, 0)
       tbl_app.FieldAdd("updated", "datetime", 0, 0)
    */
    //
    tbl_ins := setup.NewTable("ids_instance")
    tbl_ins.FieldAdd("id", "pk-string", 8, 0)
    tbl_ins.FieldAdd("uid", "uint32", 10, setup.FieldIndexIndex)
    tbl_ins.FieldAdd("app_id", "string", 50, setup.FieldIndexIndex)
    tbl_ins.FieldAdd("app_title", "string", 50, 0)
    tbl_ins.FieldAdd("version", "string", 50, 0)
    tbl_ins.FieldAdd("data", "string-text", 0, 0)
    tbl_ins.FieldAdd("created", "datetime", 0, 0)
    tbl_ins.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_rol := setup.NewTable("ids_role")
    tbl_rol.FieldAdd("rid", "auto", 0, 0)
    tbl_rol.FieldAdd("uid", "uint32", 0, setup.FieldIndexIndex)
    tbl_rol.FieldAdd("status", "int16", 0, setup.FieldIndexIndex)
    tbl_rol.FieldAdd("name", "string", 30, 0)
    tbl_rol.FieldAdd("desc", "string", 100, 0)
    tbl_rol.FieldAdd("privileges", "string-text", 0, 0)
    tbl_rol.FieldAdd("created", "datetime", 0, 0)
    tbl_rol.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_pri := setup.NewTable("ids_privilege")
    tbl_pri.FieldAdd("pid", "auto", 0, 0)
    tbl_pri.FieldAdd("instance", "string", 30, setup.FieldIndexIndex)
    tbl_pri.FieldAdd("uid", "uint32", 0, setup.FieldIndexIndex)
    tbl_pri.FieldAdd("privilege", "string", 100, 0)
    tbl_pri.FieldAdd("desc", "string", 50, 0)
    tbl_pri.FieldAdd("created", "datetime", 0, 0)

    //
    /*
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
    */

    //
    tbl_ses := setup.NewTable("ids_sessions")
    tbl_ses.FieldAdd("token", "pk-string", 24, 0)
    tbl_ses.FieldAdd("refresh", "string", 24, 0)
    tbl_ses.FieldAdd("status", "int16", 0, setup.FieldIndexIndex)
    tbl_ses.FieldAdd("uid", "uint32", 0, setup.FieldIndexIndex)
    tbl_ses.FieldAdd("name", "string", 50, 0)
    tbl_ses.FieldAdd("uname", "string", 30, 0)
    tbl_ses.FieldAdd("timezone", "string", 40, 0)
    tbl_ses.FieldAdd("roles", "string", 200, 0)
    tbl_ses.FieldAdd("source", "string", 20, 0)
    tbl_ses.FieldAdd("data", "string-text", 0, 0)
    tbl_ses.FieldAdd("permission", "int8", 0, 0)
    tbl_ses.FieldAdd("created", "datetime", 0, 0)
    tbl_ses.FieldAdd("expired", "datetime", 0, 0)

    tbl_scf := setup.NewTable("ids_sysconfig")
    tbl_scf.FieldAdd("key", "pk-string", 50, 0)
    tbl_scf.FieldAdd("value", "string-text", 0, 0)
    tbl_scf.FieldAdd("created", "datetime", 0, 0)
    tbl_scf.FieldAdd("updated", "datetime", 0, 0)

    //
    ds := setup.NewDataSet()
    ds.Version = 2
    // accounts
    ds.TableAdd(tbl_lgn)
    ds.TableAdd(tbl_prf)
    ds.TableAdd(tbl_rst)
    // group
    //ds.TableAdd(tbl_grp)
    //ds.TableAdd(tbl_gpu)
    // applications
    //ds.TableAdd(tbl_app)
    ds.TableAdd(tbl_ins)
    // roles
    ds.TableAdd(tbl_rol)
    ds.TableAdd(tbl_pri)
    //ds.TableAdd(tbl_mes)
    // session
    ds.TableAdd(tbl_ses)

    // sysconfig
    ds.TableAdd(tbl_scf)

    //
    _ = cn.Setup("", ds)

    timenow := rdc.TimeNow("datetime")
    /* _, err := cn.ExecRaw("INSERT OR IGNORE INTO `ids_group` "+
           "(gid,pid,name,summary,status,created,updated) "+
           "VALUES (1,0,\"Administrator\",\"Root System Administrator\",1,?,?),"+
           "(100,0,\"Member\",\"Universal Member\",1,?,?)",
           timenow, timenow, timenow, timenow)
       if err != nil {
           return cn, err
       } */

    _, err := cn.ExecRaw("INSERT OR IGNORE INTO `ids_role` "+
        "(rid,uid,status,name,desc,privileges,created,updated) "+
        "VALUES (1,0,1,\"Administrator\",\"Root System Administrator\",?,?,?),"+
        "(100,0,1,\"Member\",\"Universal Member\",?,?,?),"+
        "(101,0,0,\"Anonymous\",\"Anonymous Member\",?,?,?)",
        "1", timenow, timenow, "", timenow, timenow, "", timenow, timenow)
    if err != nil {
        return cn, err
    }

    _, err = cn.ExecRaw("INSERT OR IGNORE INTO `ids_instance` "+
        "(id,uid,app_id,app_title,version,created,updated) "+
        "VALUES (\"lessids\",0,\"lessids\",\"less Identity Server\",?,?,?)",
        c.Version, timenow, timenow)
    if err != nil {
        return cn, err
    }

    _, err = cn.ExecRaw("INSERT OR IGNORE INTO `ids_privilege` "+
        "(pid,instance,uid,privilege,desc,created) "+
        "VALUES (\"1\",\"lessids\",0,\"user.admin\",\"User Management\",?),"+
        "(\"1\",\"lessids\",0,\"sys.admin\",\"System Management\",?)",
        timenow, timenow)
    if err != nil {
        return cn, err
    }

    _, err = cn.ExecRaw("INSERT OR IGNORE INTO `ids_sysconfig` "+
        "(key,value,created,updated) "+
        "VALUES (\"service_name\",\"less Identity Service\",?,?),"+
        "(\"webui_banner_title\",\"Account Center\",?,?)",
        timenow, timenow, timenow, timenow)
    if err != nil {
        return cn, err
    }

    return cn, nil
}
