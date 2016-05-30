# luago
A lua wrapper for golang.

## 简介

一个采用cgo模式集成入lua 5.2.2的小模块。基本导出了所有的lua函数，在原来的基础上加上了导入go函数的功能，你可以在lua中很简单的调用注册过的go函数。

## Details

This project implement a wrapper for lua 5.2.2 using cgo.It just export almost all functions from lua, and support pushing your own go functions to the lua VM.So you can invoke your go functions in your lua script easily.

## Usage

	package main
	
	import (
		"strconv"
	
		"github.com/sryanyuan/luago"
	)
	
	func export(L luago.Lua_Handle) int {
		//	get args from lua
		num := int(luago.Lua_tonumber(L, -2))
		str := luago.Lua_tostring(L, -1)
	
		//	push value to lua
		val := str + strconv.Itoa(num)
		luago.Lua_pushstring(L, val)
		return 1
	}
	
	func main() {
		L := luago.LuaGo_newState()
		L.OpenStdLibs()
	
		L.LuaGo_PushGoFunction("export", export)
	
		//	invoke
		luago.LuaL_dostring(L.GetHandle(), ` 
			local val = export(1, "hello luago")
			print(val)
		`)
	}
