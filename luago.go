package luago

//
import "C"

import (
	"unsafe"
)

type LuaGo_Int C.int
type LuaGo_Ref int
type LuaGo_ResultSum C.int

func LuaGo_Handle(l unsafe.Pointer) Lua_Handle {
	return Lua_Handle(l)
}

func LuaGo_RegPtr(r unsafe.Pointer) *LuaL_Reg {
	return (*LuaL_Reg)(r)
}
