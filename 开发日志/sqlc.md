sqlc 怎么处理 uuid.UUID

直接让uuid.UUID作为类型是不可以的:
```yaml
db_type: "binary"
nullable: false 
go_type: 
  import: "github.com/google/uuid" 
  type: "UUID"
```

会导致传入的是string类型
```sql
-- log日志
INSERT INTO links (
id,tenant_id,user_id,short_code,original_url,created_at,expires_at
) VALUES (
'01999a70-11da-77a9-b20a-780c29cfc2dd','01999a28-ac6d-7a32-a2b3-0da95f53b4f8','01999a2d-07aa-7db3-bbee-8b0f46061171','BmPltnZjezZ','www.baidu.com','2025-09-30 19:44:19.9305021','2025-09-30 20:44:19.9305021'
)
```

还是直接用[]byte就好
