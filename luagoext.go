package luago

/*
#cgo	CFLAGS:-I../../inc
#cgo	LDFLAGS:-lm
#include <stdlib.h>
#include "lua/lua.h"
#include "lua/lauxlib.h"
#include "luagoext.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//	提供lua与go的交互，go函数注册

type LuaGo_State struct {
	//	base struct
	LuaGo_Script

	//	for gofunction export
	gofunctions map[int]interface{}

	//	base seed
	baseSeq int
}

type LuaGo_Function func(L Lua_Handle) int

//	开启gofunction注册支持
func (this *LuaGo_State) LuaGo_OpenInvokeCompenent() {
	//	初始化导出map
	C.luago_open(this.handle, unsafe.Pointer(this))
	this.gofunctions = make(map[int]interface{})
	this.baseSeq = 1
}

//	获得gofunction注册序号
func (this *LuaGo_State) luago_getGoFunctionSeq() int {
	this.baseSeq = this.baseSeq + 1
	return this.baseSeq
}

//	根据seq获得对应的GoFunction
func (this *LuaGo_State) luago_getGoFunctionBySeq(_seq int) LuaGo_Function {
	gfuncinterface := this.gofunctions[_seq]
	if nil == gfuncinterface {
		return nil
	}
	var gfunc LuaGo_Function = gfuncinterface.(LuaGo_Function)
	return gfunc
}

//	压入gofunction
func (this *LuaGo_State) LuaGo_PushGoFunction(_name string, _func LuaGo_Function) int {
	funcSeq := this.luago_getGoFunctionSeq()

	this.gofunctions[funcSeq] = _func

	//	注册到lua中
	cfuncname := C.CString(_name)
	C.luago_pushGoFunction(this.handle, cfuncname, C.int(funcSeq))
	C.free(unsafe.Pointer(cfuncname))

	return funcSeq
}

func (this *LuaGo_State) LuaGo_PopGoFunction(_name string) {
	cfuncname := C.CString(_name)
	C.luago_popGoFunction(this.handle, cfuncname)
	C.free(unsafe.Pointer(cfuncname))
}

func (this *LuaGo_State) LuaGo_SafeDoFile(_name string) bool {
	if R := LuaL_loadfile(this.handle, _name); LUA_OK != R {
		err_msg := Lua_tostring(this.handle, -1)
		Lua_pop(this.handle, 1)
		fmt.Print(err_msg)
		return false
	}

	//this.luago_pushErrorCallback()
	//return this.luago_safeDoFile()
	this.LuaGo_SafeCall(0, true)
	return true
}

func (this *LuaGo_State) luago_pushErrorCallback() {
	handler_content := `function _TRACEBACK_(errmsg)
		print("-------------------")
		print("LUA ERROR:"..tostring(errmsg))
		print("")
		print(debug.traceback("", 2))
		print("-------------------") end`
	LuaL_dostring(this.handle, handler_content)
}

//	只支持返回单返回值:number or boolean , 其余自己操作lua_state获取
//	removeret : 自动pop返回值，假设返回值个数为0，传入false
func (this *LuaGo_State) LuaGo_SafeCall(argnum int, removeret bool) int {
	var funcIndex int = -(argnum + 1)
	if !Lua_isfunction(this.handle, funcIndex) {
		//	pop all arguments and the calling function
		Lua_pop(this.handle, argnum+1)
		fmt.Print("trying to call a non-callable object")
		return 0
	}

	var traceback int = 0
	Lua_getglobal(this.handle, "_TRACEBACK_")
	if !Lua_isfunction(this.handle, -1) {
		Lua_pop(this.handle, 1)
	} else {
		Lua_insert(this.handle, funcIndex-1)
		traceback = funcIndex - 1
	}

	var err int = 0
	err = Lua_pcall(this.handle, argnum, 1, traceback)
	if err != 0 {
		err_msg := ""

		if 0 == traceback {
			//	没有指定错误回调
			err_msg = Lua_tostring(this.handle, -1)
			Lua_pop(this.handle, 1)
			fmt.Print("[LUA ERROR] ", err_msg)
		} else {
			//	指定了错误回调
			//	before call stack: traceback function
			//	after call with error stack: trackback errmsg
			Lua_pop(this.handle, 2)
		}
		return 0
	}

	var ret int = 0
	//	成功调用 stack: traceback result
	//	remove the result
	if removeret {
		if Lua_isnumber(this.handle, -1) {
			ret = (int)(Lua_tonumber(this.handle, -1))
		} else if Lua_isboolean(this.handle, -1) {
			bret := Lua_toboolean(this.handle, -1)
			if bret {
				ret = 1
			}
		}
		Lua_pop(this.handle, 1)
	}

	//	remove the traceback callback
	if 0 != traceback {
		var removeStackIndex = -2
		if removeret {
			removeStackIndex = -1
		}
		Lua_remove(this.handle, removeStackIndex)
	}

	return ret
}

func (this *LuaGo_State) luago_safeDoFile() (bool, string) {
	var funcIndex int = -(0 + 1)
	if !Lua_isfunction(this.handle, funcIndex) {
		Lua_pop(this.handle, 0+1)
		return false, "calling a non-callable object"
	}

	var traceback int = 0
	Lua_getglobal(this.handle, "_TRACEBACK_")
	if !Lua_isfunction(this.handle, -1) {
		Lua_pop(this.handle, 1)
	} else {
		Lua_insert(this.handle, funcIndex-1)
		traceback = funcIndex - 1
	}

	var err int = 0
	err = Lua_pcall(this.handle, 0, 1, traceback)
	if err != 0 {
		err_msg := ""

		if 0 == traceback {
			//	没有指定错误回调
			err_msg = Lua_tostring(this.handle, -1)
			Lua_pop(this.handle, 1)
		} else {
			//	指定了错误回调
			//	before call stack: traceback function
			//	after call with error stack: trackback errmsg
			Lua_pop(this.handle, 2)
		}
		return false, err_msg
	}

	//	成功调用 stack: traceback result
	//	remove the result
	Lua_pop(this.handle, 1)
	//	remove the traceback callback
	if 0 != traceback {
		Lua_remove(this.handle, -1)
	}

	return true, ""
}

//	增加lua搜索路径
func (this *LuaGo_State) LuaGo_AddSearchPath(_path string, _tip int) {
	cstr := C.CString(_path)
	C.luago_addSearchPath(this.handle, cstr, C.int(_tip))
	C.free(unsafe.Pointer(cstr))
}

//	内部GoFunction dispatcher
func (this *LuaGo_State) luago_internal_call(_seq int) int {
	gfunc := this.luago_getGoFunctionBySeq(_seq)
	if nil == gfunc {
		errmsg := C.CString("trying to call a unregisted go function")
		C.luago_error(this.handle, errmsg)
		C.free(unsafe.Pointer(errmsg))
		return 0
	}

	return gfunc(this.handle)
}

func LuaGo_newState() *LuaGo_State {
	luago_state := &LuaGo_State{}
	luago_state.Create()
	luago_state.LuaGo_OpenInvokeCompenent()
	luago_state.luago_pushErrorCallback()
	return luago_state
}

//	lua调用go函数的回调
//export luagoc_call
func luagoc_call(_luaGo_State unsafe.Pointer, _nFuncIndex int) int {
	//	get the dest function
	luago_state := (*LuaGo_State)(_luaGo_State)
	if nil == luago_state {
		return 0
	}

	return luago_state.luago_internal_call(_nFuncIndex)
}

//export luagoc_output
func luagoc_output(_text *C.char) {
	gostr := C.GoString(_text)
	fmt.Println(gostr)
}
