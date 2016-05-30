package luago

/*
#cgo	CFLAGS:-I../../inc
#cgo	LDFLAGS:-lm
#include <stdlib.h>
#include "lua/lua.h"
#include "./luago_helper.h"
*/
import "C"

import (
	"reflect"
	"unsafe"
)

// helpers.

// varname -> a.b.c
func LuaGoH_GetGlobal(l Lua_Handle, varname string) bool {
	S1 := C.CString(varname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaGo_GetGlobal(l, S1)
	return (0 != R)
}
func LuaGoH_SetGlobal(l Lua_Handle, varname string) bool {
	S1 := C.CString(varname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaGo_SetGlobal(l, S1)
	return (0 != R)
}

func LuaGoH_GetRef(l Lua_Handle, var_ref LuaGo_Ref) bool {
	R := C.luaGo_GetRef(l, C.int(var_ref))
	return (0 != R)
}

func LuaGoH_InvokeFunction(l Lua_Handle, nargs int, nresults int, err_code *int, err_msg *string) bool {
	var ret int
	if nil != err_msg {
		*err_msg = ""
	}
	if nil == err_code {
		err_code = &ret
	}

	switch *err_code = Lua_pcall(l, nargs, nresults, 0); *err_code {
	case LUA_OK, LUA_YIELD:
		{
			return true
		}
	default:
		{
			if nil != err_msg {
				*err_msg = Lua_tostring(l, -1)
			}
			Lua_pop(l, 1)
		}
	}
	return false
}

func luagoh_fetch_Int(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isnumber(l, -1); *bValid {
		v.SetInt(int64(Lua_tointeger(l, -1)))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_Uint(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isnumber(l, -1); *bValid {
		v.SetUint(uint64(Lua_tounsigned(l, -1)))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_Float(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isnumber(l, -1); *bValid {
		v.SetFloat(float64(Lua_tonumber(l, -1)))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_Bool(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isboolean(l, -1); *bValid {
		v.SetBool(Lua_toboolean(l, -1))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_String(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isstring(l, -1); *bValid {
		v.SetString(Lua_tostring(l, -1))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_Pointer(l Lua_Handle, v *reflect.Value, bValid *bool) {
	if *bValid = Lua_isthread(l, -1); *bValid {
		v.SetPointer(unsafe.Pointer(Lua_tothread(l, -1)))
	}
	Lua_pop(l, 1)
}

func luagoh_fetch_ERROR(l Lua_Handle, v *reflect.Value, bValid *bool) {
	*bValid = false
	Lua_pop(l, 1)
}

func LuaGoH_FetchParams(L Lua_Handle, ignore_nonexistent_field bool, args ...interface{}) bool {
	for i := len(args) - 1; i >= 0; i-- {
		if !LuaGoH_FetchVariable(L, args[i], ignore_nonexistent_field) {
			return false
		}
	}
	return true
}

func LuaGoH_FetchValue(L Lua_Handle, value *reflect.Value, ignore_nonexistent_field bool) bool {
	if !value.CanSet() {
		Lua_pop(L, 1)
		return false
	}

	if Lua_isnil(L, -1) {
		Lua_pop(L, 1)
		return false
	}

	bValid := false
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		luagoh_fetch_Int(L, value, &bValid)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		luagoh_fetch_Uint(L, value, &bValid)
	case reflect.Float32, reflect.Float64:
		luagoh_fetch_Float(L, value, &bValid)
	case reflect.Bool:
		luagoh_fetch_Bool(L, value, &bValid)
	case reflect.String:
		luagoh_fetch_String(L, value, &bValid)
	case reflect.Ptr:
		luagoh_fetch_Pointer(L, value, &bValid)
	case reflect.Struct:
		if !Lua_istable(L, -1) {
			Lua_pop(L, 1)
			break
		}
		bValid = true
		t := value.Type()
		for i := 0; i < value.NumField(); i++ {
			item := value.Field(i)
			if !item.CanSet() {
				continue
			}
			Lua_pushstring(L, t.Field(i).Name)
			Lua_gettable(L, -2)
			if !LuaGoH_FetchValue(L, &item, ignore_nonexistent_field) {
				if !ignore_nonexistent_field {
					bValid = false
					break //for
				}
			}
		}
		Lua_pop(L, 1)
	case reflect.Slice:
		if !Lua_istable(L, -1) {
			Lua_pop(L, 1)
			break
		}
		bValid = true
		l := int(Lua_rawlen(L, -1))
		tempRfV := reflect.MakeSlice(value.Type(), 0, l) //*value
		for i := 0; i < l; i++ {
			item := reflect.New(tempRfV.Type().Elem()).Elem()
			Lua_rawgeti(L, -1, i+1)
			if !LuaGoH_FetchValue(L, &item, ignore_nonexistent_field) {
				if !ignore_nonexistent_field {
					bValid = false
					break //for
				}
			}
			tempRfV = reflect.Append(tempRfV, item)
		}
		Lua_pop(L, 1)
		value.Set(tempRfV)
	case reflect.Array:
		if !Lua_istable(L, -1) {
			Lua_pop(L, 1)
			break
		}
		bValid = true
		l := value.Len()
		for i := 0; i < l; i++ {
			Lua_rawgeti(L, -1, i+1)
			item := value.Index(i)
			if !LuaGoH_FetchValue(L, &item, ignore_nonexistent_field) {
				if !ignore_nonexistent_field {
					bValid = false
					break //for
				}
			}
		}
		Lua_pop(L, 1)
	case reflect.Map:
		if !Lua_istable(L, -1) {
			Lua_pop(L, 1)
			break
		}
		bValid = true

		if value.IsNil() {
			value.Set(reflect.MakeMap(value.Type()))
		}

		Lua_pushnil(L)
		for 0 != Lua_next(L, -2) {
			Lua_pushvalue(L, -2) // table key value key

			k := reflect.New(value.Type().Key()).Elem()
			if !LuaGoH_FetchValue(L, &k, ignore_nonexistent_field) {
				if !ignore_nonexistent_field {
					Lua_pop(L, 2)
					bValid = false
					break //for
				}
				Lua_pop(L, 1)
				continue
			}
			v := reflect.New(value.Type().Elem()).Elem()
			if !LuaGoH_FetchValue(L, &v, ignore_nonexistent_field) {
				if !ignore_nonexistent_field {
					Lua_pop(L, 1)
					bValid = false
					break //for
				}
				continue
			}
			value.SetMapIndex(k, v)
		}
		Lua_pop(L, 1)
	default:
		luagoh_fetch_ERROR(L, value, &bValid)
	}
	return bValid
}

func LuaGoH_FetchVariable(L Lua_Handle, value interface{}, ignore_nonexistent_field bool) bool {
	r := reflect.ValueOf(value)
	if r.Kind() != reflect.Ptr {
		Lua_pop(L, 1)
		return false
	}

	v := r.Elem()
	if !v.CanSet() {
		Lua_pop(L, 1)
		return false
	}

	return LuaGoH_FetchValue(L, &v, ignore_nonexistent_field)
}

func LuaGoH_PushValue(L Lua_Handle, value *reflect.Value) bool {
	if !value.IsValid() {
		return false
	}

	ok := true
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if r, ok := value.Interface().(LuaGo_Ref); ok {
			LuaGoH_GetRef(L, r)
		} else {
			Lua_pushinteger(L, Lua_Integer(value.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		Lua_pushunsigned(L, Lua_Unsigned(value.Uint()))
	case reflect.Float32, reflect.Float64:
		Lua_pushnumber(L, Lua_Number(value.Float()))
	case reflect.Bool:
		Lua_pushboolean(L, value.Bool())
	case reflect.String:
		Lua_pushstring(L, value.String())
	case reflect.Struct:
		Lua_newtable(L)
		t := value.Type()
		for i := 0; i < value.NumField(); i++ {
			item := value.Field(i)
			if !item.CanSet() {
				continue
			}
			Lua_pushstring(L, t.Field(i).Name)
			if !LuaGoH_PushValue(L, &item) {
				Lua_pop(L, 2)
				ok = false
				break
			}
			Lua_settable(L, -3)
		}
	case reflect.Slice, reflect.Array:
		Lua_newtable(L)
		l := value.Len()
		for i := 0; i < l; i++ {
			v := value.Index(i)
			if !LuaGoH_PushValue(L, &v) {
				Lua_pop(L, 1)
				ok = false
				break
			}
			Lua_rawseti(L, -2, i+1) // go begin with 0, lua table is 1
		}
	case reflect.Map:
		Lua_newtable(L)
		keys := value.MapKeys()
		l := len(keys)
		for i := 0; i < l; i++ {
			k := keys[i]
			v := value.MapIndex(k)
			if !LuaGoH_PushValue(L, &k) {
				Lua_pop(L, 1)
				ok = false
				break
			}
			if !LuaGoH_PushValue(L, &v) {
				Lua_pop(L, 2)
				ok = false
				break
			}
			Lua_settable(L, -3)
		}
	default:
		{
			v := value.Interface()
			switch v.(type) {
			case Lua_Handle:
				Lua_pushthread(L)
			case Lua_CFunction:
				Lua_pushcfunction(L, v.(Lua_CFunction))
			default:
				ok = false
			}
		}
	}
	return ok
}

func LuaGoH_PushVariable(L Lua_Handle, value interface{}) bool {
	r := reflect.ValueOf(value)
	if r.Kind() == reflect.Ptr {
		v := r.Elem()
		if !v.CanSet() {
			return false
		}

		return LuaGoH_PushValue(L, &v)
	} else {
		return LuaGoH_PushValue(L, &r)
	}
	return false
}
