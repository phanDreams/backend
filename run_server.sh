#!/bin/bash

export PG_DSN="postgres://postgres.jjniaoujdyanpybytgnq:GeX1zBdRsMn3c3P!@db.jjniaoujdyanpybytgnq.supabase.co:5432/postgres?sslmode=require"
export SERVER_ADDRESS=":3000"

go run cmd/pethelp/main.go 