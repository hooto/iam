package conf

const databaseSchema = `
{
    "engine": "InnoDB",
    "charset": "utf8",
    "tables": [
        {
            "name": "less_dataset_version",
            "columns": [
                {
                    "name": "id",
                    "incrAble": true,
                    "type": "uint64"                    
                },
                {
                    "name": "version",
                    "type": "uint32"
                },
                {
                    "name": "action",
                    "type": "string",
                    "nullAble": true,
                    "length": "30"
                },
                {
                    "name": "created",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "version",
                    "type": 1,
                    "cols": ["version"]
                }
            ]
        },
        {
            "name": "ids_login",
            "columns": [
                {
                    "name": "uid",
                    "incrAble": true,
                    "type": "uint64"
                },
                {
                    "name": "uuid",
                    "type": "string",
                    "length": "8"
                },
                {
                    "name": "uname",
                    "type": "string",
                    "length": "20"
                },
                {
                    "name": "email",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "name",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "pass",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "group",
                    "type": "string",
                    "length": "200"
                },
                {
                    "name": "roles",
                    "type": "string",
                    "length": "200"
                },
                {
                    "name": "timezone",
                    "type": "string",
                    "length": "40"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["uid"]
                },
                {
                    "name": "uuid",
                    "type": 2,
                    "cols": ["uuid"]
                },
                {
                    "name": "uname",
                    "type": 2,
                    "cols": ["uname"]
                },
                {
                    "name": "email",
                    "type": 2,
                    "cols": ["email"]
                },
                {
                    "name": "status",
                    "type": 1,
                    "cols": ["status"]
                },
                {
                    "name": "created",
                    "type": 1,
                    "cols": ["created"]
                },
                {
                    "name": "updated",
                    "type": 1,
                    "cols": ["updated"]
                }
            ]
        },
        {
            "name": "ids_profile",
            "columns": [
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "gender",
                    "type": "uint8"
                },
                {
                    "name": "birthday",
                    "type": "date"
                },
                {
                    "name": "address",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "aboutme",
                    "type": "string-text"
                },
                {
                    "name": "photo",
                    "type": "string-text"
                },
                {
                    "name": "photosrc",
                    "type": "string-text"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["uid"]
                }
            ]
        },
        {
            "name": "ids_resetpass",
            "columns": [
                {
                    "name": "id",
                    "type": "string",
                    "length": "24"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "email",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "expired",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "expired",
                    "type": 1,
                    "cols": ["expired"]
                }
            ]
        },
        {
            "name": "ids_instance",
            "columns": [
                {
                    "name": "id",
                    "type": "string",
                    "length": "8"
                },
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "app_id",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "app_title",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "url",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "version",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "privileges",
                    "type": "uint64"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "uid",
                    "type": 1,
                    "cols": ["uid"]
                },
                {
                    "name": "status",
                    "type": 1,
                    "cols": ["status"]
                },
                {
                    "name": "app_id",
                    "type": 1,
                    "cols": ["app_id"]
                }
            ]
        },
        {
            "name": "ids_role",
            "columns": [
                {
                    "name": "rid",
                    "type": "uint64",
                    "incrAble": true
                },
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "name",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "desc",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "privileges",
                    "type": "string-text"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["rid"]
                },
                {
                    "name": "uid",
                    "type": 1,
                    "cols": ["uid"]
                },
                {
                    "name": "status",
                    "type": 1,
                    "cols": ["status"]
                }
            ]
        },
        {
            "name": "ids_privilege",
            "columns": [
                {
                    "name": "pid",
                    "type": "uint64",
                    "incrAble": true
                },
                {
                    "name": "instance",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "privilege",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "desc",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "created",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["pid"]
                },
                {
                    "name": "instance",
                    "type": 1,
                    "cols": ["instance"]
                },
                {
                    "name": "uid",
                    "type": 1,
                    "cols": ["uid"]
                }
            ]
        },
        {
            "name": "ids_sessions",
            "columns": [
                {
                    "name": "token",
                    "type": "string",
                    "length": "24"
                },
                {
                    "name": "refresh",
                    "type": "string",
                    "length": "24"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "uuid",
                    "type": "string",
                    "length": "8"
                },
                {
                    "name": "name",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "uname",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "timezone",
                    "type": "string",
                    "length": "40"
                },
                {
                    "name": "roles",
                    "type": "string",
                    "length": "200"
                },
                {
                    "name": "source",
                    "type": "string",
                    "length": "20"
                },
                {
                    "name": "data",
                    "type": "string-text"
                },
                {
                    "name": "permission",
                    "type": "uint8"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "expired",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["token"]
                },
                {
                    "name": "status",
                    "type": 1,
                    "cols": ["status"]
                },
                {
                    "name": "uid",
                    "type": 1,
                    "cols": ["uid"]
                },
                {
                    "name": "uuid",
                    "type": 1,
                    "cols": ["uuid"]
                }
            ]
        },
        {
            "name": "ids_sysconfig",
            "columns": [
                {
                    "name": "key",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "value",
                    "type": "string-text"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["key"]
                }
            ]
        },
        {
            "name": "ids_myapp",
            "columns": [
                {
                    "name": "instid",
                    "type": "string",
                    "length": "8"
                },
                {
                    "name": "uid",
                    "type": "uint64"
                },
                {
                    "name": "status",
                    "type": "uint16"
                },
                {
                    "name": "app_id",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "app_title",
                    "type": "string",
                    "length": "50"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["instid", "uid"]
                },
                {
                    "name": "status",
                    "type": 1,
                    "cols": ["status"]
                }
            ]
        }
    ]
}
`
