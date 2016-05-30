package luago

/*
#cgo	CFLAGS:-I../../inc
#cgo	LDFLAGS:-lm
#include "lua/luaconf.h"
*/
import "C"

const LUA_DIRSEP = C.LUA_DIRSEP

const LUA_ENV = C.LUA_ENV

const LUA_IDSIZE = C.LUA_IDSIZE

type LUA_INT32 C.LUA_INT32

type LUA_NUMBER C.LUA_NUMBER

const LUA_NUMBER_SCAN = C.LUA_NUMBER_SCAN
const LUA_NUMBER_FMT = C.LUA_NUMBER_FMT

// TODO: #define lua_number2str(s,n)	sprintf((s), LUA_NUMBER_FMT, (n))
// TODO: #define lua_str2number(s,p)	strtod((s), (p))
// TODO: #define lua_strx2number(s,p)	strtod((s), (p))

type LUA_INTEGER C.LUA_INTEGER

type LUA_UNSIGNED C.LUA_UNSIGNED
