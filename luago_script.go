package luago

import (
	"reflect"
	"unsafe"
)

//
//	LuaGo_ScriptBase
//
type LuaGo_ScriptBase struct {
	handle  Lua_Handle
	err_msg string
}

func (s *LuaGo_ScriptBase) GetHandle() Lua_Handle {
	return s.handle
}

func (s *LuaGo_ScriptBase) SetHandle(l Lua_Handle) {
	s.handle = l
}

// called must after dostring / dofile
func (s *LuaGo_ScriptBase) GetErrorMessage() string {
	return s.err_msg
}

func (s *LuaGo_ScriptBase) OpenStdLibs() {
	if nil != s.handle {
		LuaL_openlibs(s.handle)
	}
}

func (s *LuaGo_ScriptBase) HasRef(ref LuaGo_Ref) bool {
	if !LuaGoH_GetRef(s.handle, ref) {
		return false
	}

	R := !Lua_isnoneornil(s.handle, -1)
	Lua_pop(s.handle, 1)
	return R
}

func (s *LuaGo_ScriptBase) GetRef(var_name string, pref *LuaGo_Ref) bool {
	if !LuaGoH_GetGlobal(s.handle, var_name) {
		return false
	}

	bExist := !Lua_isnoneornil(s.handle, -1)
	if bExist && nil != pref {
		*pref = LuaGo_Ref(LuaL_ref(s.handle, LUA_REGISTRYINDEX))
	} else {
		Lua_pop(s.handle, 1)
	}

	return bExist
}

func (s *LuaGo_ScriptBase) RemoveRef(ref LuaGo_Ref) bool {
	LuaL_unref(s.handle, LUA_REGISTRYINDEX, int(ref))

	return true
}

func (s *LuaGo_ScriptBase) LoadRef(ref LuaGo_Ref) bool {
	return LuaGoH_GetRef(s.handle, ref)
}

func (s *LuaGo_ScriptBase) HasVar(var_name string) bool {
	if !LuaGoH_GetGlobal(s.handle, var_name) {
		return false
	}

	R := !Lua_isnoneornil(s.handle, -1)
	Lua_pop(s.handle, 1)
	return R
}

func (s *LuaGo_ScriptBase) RemoveVar(var_name string) {
	Lua_pushnil(s.handle)
	if !LuaGoH_SetGlobal(s.handle, var_name) {
		Lua_pop(s.handle, 1)
	}
}

func (s *LuaGo_ScriptBase) GetVar(var_name string, value interface{}) bool {
	return s.GetObject(var_name, value, true)
}

func (s *LuaGo_ScriptBase) SetVar(var_name string, value interface{}) bool {
	return s.SetObject(var_name, value, false)
}

func (s *LuaGo_ScriptBase) GetObject(var_name string, value interface{}, ignore_nonexistent_field bool) bool {
	r := reflect.ValueOf(value)
	if r.Kind() != reflect.Ptr {
		return false
	}

	v := r.Elem()
	if !v.CanSet() {
		return false
	}

	if !LuaGoH_GetGlobal(s.handle, var_name) {
		return false
	}

	return LuaGoH_FetchVariable(s.handle, value, ignore_nonexistent_field)
}

func (s *LuaGo_ScriptBase) SetObject(var_name string, value interface{}, keep_nonexistent_field bool) bool {
	if !keep_nonexistent_field {
		s.RemoveVar(var_name)
	}

	if !LuaGoH_PushVariable(s.handle, value) {
		return false
	}

	if !LuaGoH_SetGlobal(s.handle, var_name) {
		Lua_pop(s.handle, 1)
		return false
	}

	return true
}

func (s *LuaGo_ScriptBase) Call(func_name string, args ...interface{}) bool {
	if !LuaGoH_GetGlobal(s.handle, func_name) {
		return false
	}

	for i := 0; i < len(args); i++ {
		if !LuaGoH_PushVariable(s.handle, args[i]) {
			Lua_pop(s.handle, i+1) // 0, 1, .. i - 1, + LuaGoH_GetGlobal
			return false
		}
	}

	return LuaGoH_InvokeFunction(s.handle, len(args), 0, nil, &s.err_msg)
}

func (s *LuaGo_ScriptBase) Invoke(ret_value interface{}, func_name string, args ...interface{}) bool {
	if !LuaGoH_GetGlobal(s.handle, func_name) {
		return false
	}

	for i := 0; i < len(args); i++ {
		if !LuaGoH_PushVariable(s.handle, args[i]) {
			Lua_pop(s.handle, i+1) // 0, 1, .. i - 1, + LuaGoH_GetGlobal
			return false
		}
	}

	retsum := 0
	if nil != ret_value {
		retsum = 1
	}

	if !LuaGoH_InvokeFunction(s.handle, len(args), retsum, nil, &s.err_msg) {
		return false
	}

	if nil != ret_value {
		return LuaGoH_FetchVariable(s.handle, ret_value, true)
	}

	return true
}

func (s *LuaGo_ScriptBase) RunFile(file string) bool {
	if R := LuaL_loadfile(s.handle, file); LUA_OK != R {
		s.err_msg = Lua_tostring(s.handle, -1)
		Lua_pop(s.handle, 1)
		return false
	}
	return LuaGoH_InvokeFunction(s.handle, 0, LUA_MULTRET, nil, &s.err_msg)
}

func (s *LuaGo_ScriptBase) RunString(code string) bool {
	if R := LuaL_loadstring(s.handle, code); LUA_OK != R {
		s.err_msg = Lua_tostring(s.handle, -1)
		Lua_pop(s.handle, 1)
		return false
	}
	return LuaGoH_InvokeFunction(s.handle, 0, LUA_MULTRET, nil, &s.err_msg)
}

func (s *LuaGo_ScriptBase) RunBuffer(buffer unsafe.Pointer, size uint) bool {
	if LUA_OK == LuaL_loadbuffer(s.handle, buffer, size, "LuaGo_ScriptBase.RunBuffer") {
		return LuaGoH_InvokeFunction(s.handle, 0, LUA_MULTRET, nil, &s.err_msg)
	}
	s.err_msg = Lua_tostring(s.handle, -1)
	Lua_pop(s.handle, 1)
	return false
}

//	simple impls.
type LuaGo_Script struct {
	LuaGo_ScriptBase
}

func CreateLuaScript() *LuaGo_Script {
	s := &LuaGo_Script{}
	s.Create()
	return s
}

func (s *LuaGo_Script) Create() {
	if nil != s.handle {
		s.Destroy()
	}
	s.handle = LuaL_newstate()
}

func (s *LuaGo_Script) Destroy() {
	if nil != s.handle {
		Lua_close(s.handle)
	}
}
