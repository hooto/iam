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
    tbl_pkg := setup.NewTable("lps_pkgs")
    tbl_pkg.FieldAdd("id", "pk-string", 100, 0)
    tbl_pkg.FieldAdd("bucket", "string", 20, setup.FieldIndexIndex)
    tbl_pkg.FieldAdd("name", "string", 50, setup.FieldIndexIndex)
    tbl_pkg.FieldAdd("version", "string", 30, 0)
    tbl_pkg.FieldAdd("release", "string", 30, 0)
    tbl_pkg.FieldAdd("summary", "string-text", 0, 0)
    tbl_pkg.FieldAdd("license", "string", 100, 0)
    tbl_pkg.FieldAdd("grp_app", "string", 50, 0)
    tbl_pkg.FieldAdd("grp_dev", "string", 50, 0)
    tbl_pkg.FieldAdd("pkg_os", "string", 10, setup.FieldIndexIndex)
    tbl_pkg.FieldAdd("pkg_arch", "string", 10, setup.FieldIndexIndex)
    tbl_pkg.FieldAdd("pkg_sum", "string", 50, 0)
    tbl_pkg.FieldAdd("homepage", "string", 200, 0)
    tbl_pkg.FieldAdd("keywords", "string", 200, 0)
    tbl_pkg.FieldAdd("created", "datetime", 0, setup.FieldIndexIndex)
    tbl_pkg.FieldAdd("updated", "datetime", 0, setup.FieldIndexIndex)

    //
    tbl_buk := setup.NewTable("lps_buckets")
    tbl_buk.FieldAdd("id", "pk-string", 20, 0)
    tbl_buk.FieldAdd("type", "string", 10, 0)
    tbl_buk.FieldAdd("title", "string", 100, 0)
    tbl_buk.FieldAdd("pkg_num", "uint32", 0, 0)
    tbl_buk.FieldAdd("created", "datetime", 0, 0)
    tbl_buk.FieldAdd("updated", "datetime", 0, 0)

    //
    tbl_lpm_rpo := setup.NewTable("lpm_repos")
    tbl_lpm_rpo.FieldAdd("id", "pk-string", 20, 0)
    tbl_lpm_rpo.FieldAdd("type", "string", 10, 0)
    tbl_lpm_rpo.FieldAdd("status", "uint32", 0, setup.FieldIndexIndex)
    tbl_lpm_rpo.FieldAdd("title", "string", 100, 0)
    tbl_lpm_rpo.FieldAdd("pkg_num", "uint32", 0, 0)
    tbl_lpm_rpo.FieldAdd("srcurl", "string", 100, 0)
    tbl_lpm_rpo.FieldAdd("created", "datetime", 0, 0)
    tbl_lpm_rpo.FieldAdd("updated", "datetime", 0, setup.FieldIndexIndex)

    //
    tbl_lpm_pkg := setup.NewTable("lpm_pkgs")
    tbl_lpm_pkg.FieldAdd("id", "pk-string", 100, 0)
    tbl_lpm_pkg.FieldAdd("bucket", "string", 20, setup.FieldIndexIndex)
    tbl_lpm_pkg.FieldAdd("name", "string", 50, setup.FieldIndexIndex)
    tbl_lpm_pkg.FieldAdd("version", "string", 30, 0)
    tbl_lpm_pkg.FieldAdd("release", "string", 30, 0)
    tbl_lpm_pkg.FieldAdd("summary", "string-text", 0, 0)
    tbl_lpm_pkg.FieldAdd("license", "string", 100, 0)
    tbl_lpm_pkg.FieldAdd("grp_app", "string", 50, 0)
    tbl_lpm_pkg.FieldAdd("grp_dev", "string", 50, 0)
    tbl_lpm_pkg.FieldAdd("pkg_os", "string", 10, setup.FieldIndexIndex)
    tbl_lpm_pkg.FieldAdd("pkg_arch", "string", 10, setup.FieldIndexIndex)
    tbl_lpm_pkg.FieldAdd("pkg_sum", "string", 50, 0)
    tbl_lpm_pkg.FieldAdd("homepage", "string", 200, 0)
    tbl_lpm_pkg.FieldAdd("keywords", "string", 200, 0)
    tbl_lpm_pkg.FieldAdd("created", "datetime", 0, setup.FieldIndexIndex)
    tbl_lpm_pkg.FieldAdd("updated", "datetime", 0, setup.FieldIndexIndex)

    //
    ds := setup.NewDataSet()
    ds.Version = 11
    ds.TableAdd(tbl_pkg)
    ds.TableAdd(tbl_buk)
    ds.TableAdd(tbl_lpm_rpo)
    ds.TableAdd(tbl_lpm_pkg)

    //
    _ = cn.Setup("", ds)

    return cn, nil
}
