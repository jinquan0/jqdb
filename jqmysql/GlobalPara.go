package jqmysql

import (
    "sync"
)

//MySQL 状态指标统计
var (
    FlagDbg bool = false
    MyConnTotal int64 = 0
    MyConnErr int64 = 0
    MyWriteTotal int64 = 0
    MyReadTotal int64 = 0
    ParaLock_my_conn_total sync.Mutex
    ParaLock_my_conn_err sync.Mutex
    ParaLock_my_write_total sync.Mutex
    ParaLock_my_read_total sync.Mutex
)