package luago

/*
#cgo	CFLAGS:-I../../inc
#cgo	LDFLAGS:-lm
#include "lua/lualib.h"
*/
import "C"

func Luaopen_base(L Lua_Handle) int { return int(C.luaopen_base(L)) }

const LUA_COLIBNAME = C.LUA_COLIBNAME

func Luaopen_coroutine(L Lua_Handle) int { return int(C.luaopen_coroutine(L)) }

const LUA_TABLIBNAME = C.LUA_TABLIBNAME

func Luaopen_table(L Lua_Handle) int { return int(C.luaopen_table(L)) }

const LUA_IOLIBNAME = C.LUA_IOLIBNAME

func Luaopen_io(L Lua_Handle) int { return int(C.luaopen_io(L)) }

const LUA_OSLIBNAME = C.LUA_OSLIBNAME

func Luaopen_os(L Lua_Handle) int { return int(C.luaopen_os(L)) }

const LUA_STRLIBNAME = C.LUA_STRLIBNAME

func Luaopen_string(L Lua_Handle) int { return int(C.luaopen_string(L)) }

const LUA_BITLIBNAME = C.LUA_BITLIBNAME

func Luaopen_bit32(L Lua_Handle) int { return int(C.luaopen_bit32(L)) }

const LUA_MATHLIBNAME = C.LUA_MATHLIBNAME

func Luaopen_math(L Lua_Handle) int { return int(C.luaopen_math(L)) }

const LUA_DBLIBNAME = C.LUA_DBLIBNAME

func Luaopen_debug(L Lua_Handle) int { return int(C.luaopen_debug(L)) }

const LUA_LOADLIBNAME = C.LUA_LOADLIBNAME

func Luaopen_package(L Lua_Handle) int { return int(C.luaopen_package(L)) }

/* open all previous libraries */
func LuaL_openlibs(L Lua_Handle) { C.luaL_openlibs(L) }
