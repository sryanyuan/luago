package luago

/*
#cgo	CFLAGS:-I./inc
#cgo	LDFLAGS:-lm
#include <stdlib.h>
#include "lua/lauxlib.h"
void	xluaL_addchar(luaL_Buffer* B, char c)
{
  ((void)((B)->n < (B)->size || luaL_prepbuffsize((B), 1)), \
   ((B)->b[(B)->n++] = (c)));
}
*/
import "C"

import (
	"unsafe"
)

/* extra error code for `luaL_load' */
const LUA_ERRFILE = C.LUA_ERRFILE

type LuaL_Reg C.luaL_Reg

func LuaL_checkversion_(L Lua_Handle, ver Lua_Number) { C.luaL_checkversion_(L, C.lua_Number(ver)) }
func LuaL_checkversion(L Lua_Handle)                  { C.luaL_checkversion_(L, LUA_VERSION_NUM) }

func LuaL_getmetafield(L Lua_Handle, obj int, e string) int {
	S1 := C.CString(e)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_getmetafield(L, C.int(obj), S1)
	return int(R)
}
func LuaL_callmeta(L Lua_Handle, obj int, e string) int {
	S1 := C.CString(e)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_callmeta(L, C.int(obj), S1)
	return int(R)
}
func LuaL_tolstring(L Lua_Handle, idx int, len *uint) string {
	var sz C.size_t
	s := C.GoString(C.luaL_tolstring(L, C.int(idx), &sz))
	if len != nil {
		*len = uint(sz)
	}
	return s
}
func LuaL_argerror(L Lua_Handle, numarg int, extramsg string) int {
	S1 := C.CString(extramsg)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_argerror(L, C.int(numarg), S1)
	return int(R)
}
func LuaL_checklstring(L Lua_Handle, numArg int, l *uint) string {
	var sz C.size_t
	s := C.GoString(C.luaL_checklstring(L, C.int(numArg), &sz))
	if l != nil {
		*l = uint(sz)
	}
	return s
}
func LuaL_optlstring(L Lua_Handle, numArg int, def string, l *uint) string {
	S1 := C.CString(def)
	defer C.free(unsafe.Pointer(S1))
	var sz C.size_t
	s := C.GoString(C.luaL_optlstring(L, C.int(numArg), S1, &sz))
	if l != nil {
		*l = uint(sz)
	}
	return s
}
func LuaL_checknumber(L Lua_Handle, numArg int) Lua_Number {
	R := C.luaL_checknumber(L, C.int(numArg))
	return Lua_Number(R)
}
func LuaL_optnumber(L Lua_Handle, nArg int, def Lua_Number) Lua_Number {
	R := C.luaL_optnumber(L, C.int(nArg), C.lua_Number(def))
	return Lua_Number(R)
}

func LuaL_checkinteger(L Lua_Handle, numArg int) Lua_Integer {
	R := C.luaL_checkinteger(L, C.int(numArg))
	return Lua_Integer(R)
}
func LuaL_optinteger(L Lua_Handle, nArg int, def Lua_Integer) Lua_Integer {
	R := C.luaL_optinteger(L, C.int(nArg), C.lua_Integer(def))
	return Lua_Integer(R)
}
func LuaL_checkunsigned(L Lua_Handle, numArg int) Lua_Unsigned {
	R := C.luaL_checkunsigned(L, C.int(numArg))
	return Lua_Unsigned(R)
}
func LuaL_optunsigned(L Lua_Handle, numArg int, def Lua_Unsigned) Lua_Unsigned {
	R := C.luaL_optunsigned(L, C.int(numArg), C.lua_Unsigned(def))
	return Lua_Unsigned(R)
}

func LuaL_checkstack(L Lua_Handle, sz int, msg string) {
	S1 := C.CString(msg)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_checkstack(L, C.int(sz), S1)
}
func LuaL_checktype(L Lua_Handle, narg int, t int) { C.luaL_checktype(L, C.int(narg), C.int(t)) }
func LuaL_checkany(L Lua_Handle, narg int)         { C.luaL_checkany(L, C.int(narg)) }

func LuaL_newmetatable(L Lua_Handle, tname string) int {
	S1 := C.CString(tname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_newmetatable(L, S1)
	return int(R)
}
func LuaL_setmetatable(L Lua_Handle, tname string) {
	S1 := C.CString(tname)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_setmetatable(L, S1)
}
func LuaL_testudata(L Lua_Handle, ud int, tname string) unsafe.Pointer {
	S1 := C.CString(tname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_testudata(L, C.int(ud), S1)
	return (R)
}
func LuaL_checkudata(L Lua_Handle, ud int, tname string) unsafe.Pointer {
	S1 := C.CString(tname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_checkudata(L, C.int(ud), S1)
	return (R)
}

func LuaL_where(L Lua_Handle, lvl int) { C.luaL_where(L, C.int(lvl)) }

// TODO: func	int (luaL_error) (L Lua_Handle, const char *fmt, ...);

func LuaL_checkoption(L Lua_Handle, narg int, def string, lst []string) int {
	l := make([](*C.char), 0, len(lst)+1)
	for i := 0; i < len(lst); i++ {
		S1 := C.CString(lst[i])
		defer C.free(unsafe.Pointer(S1))
		l = append(l, S1)
	}
	l = append(l, nil)

	S2 := C.CString(def)
	defer C.free(unsafe.Pointer(S2))

	R := C.luaL_checkoption(L, C.int(narg), S2, &l[0])
	return int(R)
}

func LuaL_fileresult(L Lua_Handle, stat int, fname string) int {
	S1 := C.CString(fname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_fileresult(L, C.int(stat), S1)
	return int(R)
}
func LuaL_execresult(L Lua_Handle, stat int) int {
	R := C.luaL_execresult(L, C.int(stat))
	return int(R)
}

/* pre-defined references */
const LUA_NOREF = C.LUA_NOREF
const LUA_REFNIL = C.LUA_REFNIL

func LuaL_ref(L Lua_Handle, t int) int        { R := C.luaL_ref(L, C.int(t)); return int(R) }
func LuaL_unref(L Lua_Handle, t int, ref int) { C.luaL_unref(L, C.int(t), C.int(ref)) }

func LuaL_loadfilex(L Lua_Handle, filename string, mode string) int {
	S1 := C.CString(filename)
	defer C.free(unsafe.Pointer(S1))
	S2 := C.CString(mode)
	defer C.free(unsafe.Pointer(S2))
	R := C.luaL_loadfilex(L, S1, S2)
	return int(R)
}

func LuaL_loadfile(L Lua_Handle, filename string) int {
	S1 := C.CString(filename)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_loadfilex(L, S1, nil)
	return int(R)
}

func LuaL_loadbufferx(L Lua_Handle, buff unsafe.Pointer, sz uint, name string, mode string) int {
	S1 := C.CString(name)
	defer C.free(unsafe.Pointer(S1))
	S2 := C.CString(mode)
	defer C.free(unsafe.Pointer(S2))
	R := C.luaL_loadbufferx(L, (*C.char)(buff), C.size_t(sz), S1, S2)
	return int(R)
}
func LuaL_loadstring(L Lua_Handle, s string) int {
	S1 := C.CString(s)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_loadstring(L, S1)
	return int(R)
}

func LuaL_newstate() Lua_Handle { R := C.luaL_newstate(); return Lua_Handle(R) }

func LuaL_len(L Lua_Handle, idx int) int { R := C.luaL_len(L, C.int(idx)); return int(R) }

func LuaL_gsub(L Lua_Handle, s string, p string, r string) string {
	S1 := C.CString(s)
	defer C.free(unsafe.Pointer(S1))
	S2 := C.CString(p)
	defer C.free(unsafe.Pointer(S2))
	S3 := C.CString(r)
	defer C.free(unsafe.Pointer(S3))
	R := C.luaL_gsub(L, S1, S2, S3)
	return C.GoString(R)
}

func LuaL_setfuncs(L Lua_Handle, lst *LuaL_Reg, nup int) {
	C.luaL_setfuncs(L, (*C.luaL_Reg)(lst), C.int(nup))
}

func LuaL_getsubtable(L Lua_Handle, idx int, fname string) int {
	S1 := C.CString(fname)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_getsubtable(L, C.int(idx), S1)
	return int(R)
}

func LuaL_traceback(L Lua_Handle, L1 Lua_Handle, msg string, level int) {
	S1 := C.CString(msg)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_traceback(L, L1, S1, C.int(level))
}

func LuaL_requiref(L Lua_Handle, modname string, openf Lua_CFunction, glb int) {
	S1 := C.CString(modname)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_requiref(L, S1, C.lua_CFunction(openf), C.int(glb))
}

/*
** ===============================================================
** some useful macros
** ===============================================================
 */

// NOTE!!!	LuaL_newlibtable second param -> len(lst) => 8

func LuaL_newlibtable(L Lua_Handle, lst *LuaL_Reg) { Lua_createtable(L, 0, 8) }

func LuaL_newlib(L Lua_Handle, lst *LuaL_Reg) { LuaL_newlibtable(L, lst); LuaL_setfuncs(L, lst, 0) }

func LuaL_argcheck(L Lua_Handle, cond bool, numarg int, extramsg string) {
	if !cond {
		S1 := C.CString(extramsg)
		defer C.free(unsafe.Pointer(S1))
		C.luaL_argerror(L, C.int(numarg), S1)
	}
}
func LuaL_checkstring(L Lua_Handle, numArg int) string { return LuaL_checklstring(L, numArg, nil) }
func LuaL_optstring(L Lua_Handle, numArg int, def string) string {
	return LuaL_optlstring(L, numArg, def, nil)
}
func LuaL_checkint(L Lua_Handle, numArg int) int { R := LuaL_checkinteger(L, numArg); return int(R) }
func LuaL_optint(L Lua_Handle, numArg int, def int) int {
	R := LuaL_optinteger(L, numArg, Lua_Integer(def))
	return int(R)
}
func LuaL_checklong(L Lua_Handle, numArg int) int32 {
	R := LuaL_checkinteger(L, numArg)
	return int32(R)
}
func LuaL_optlong(L Lua_Handle, numArg int, def int32) int32 {
	R := LuaL_optinteger(L, numArg, Lua_Integer(def))
	return int32(R)
}

func LuaL_typename(L Lua_Handle, idx int) string { return Lua_typename(L, Lua_type(L, (idx))) }

func LuaL_dofile(L Lua_Handle, fn string) int {
	R := LuaL_loadfile(L, fn)
	if LUA_OK != R {
		return R
	}
	return Lua_pcall(L, 0, LUA_MULTRET, 0)
}

func LuaL_dostring(L Lua_Handle, s string) int {
	R := LuaL_loadstring(L, s)
	if LUA_OK != R {
		return R
	}
	return Lua_pcall(L, 0, LUA_MULTRET, 0)
}

func LuaL_getmetatable(L Lua_Handle, name string) { Lua_getfield(L, LUA_REGISTRYINDEX, name) }

// TODO: func	LuaL_opt(L Lua_Handle, f Lua_CFunction, n int, d int)int				{ if Lua_isnoneornil(L,(n)) { return (d); } else { return f(L,(n)); }; }

func LuaL_loadbuffer(L Lua_Handle, buff unsafe.Pointer, sz uint, name string) int {
	S1 := C.CString(name)
	defer C.free(unsafe.Pointer(S1))
	R := C.luaL_loadbufferx(L, (*C.char)(buff), C.size_t(sz), S1, nil)
	return int(R)
}

/*
** {======================================================
** Generic Buffer manipulation
** =======================================================
 */

type LuaL_Buffer C.luaL_Buffer

func LuaL_addchar(B *LuaL_Buffer, c byte) { C.xluaL_addchar((*C.luaL_Buffer)(B), C.char(c)) }

func LuaL_addsize(B *LuaL_Buffer, s uint) { B.n += C.size_t(s) }

func LuaL_buffinit(L Lua_Handle, B *LuaL_Buffer) { C.luaL_buffinit(L, (*C.luaL_Buffer)(B)) }
func LuaL_prepbuffsize(B *LuaL_Buffer, sz uint) unsafe.Pointer {
	R := C.luaL_prepbuffsize((*C.luaL_Buffer)(B), C.size_t(sz))
	return unsafe.Pointer(R)
}
func LuaL_addlstring(B *LuaL_Buffer, s string, l uint) {
	S1 := C.CString(s)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_addlstring((*C.luaL_Buffer)(B), S1, C.size_t(l))
}
func LuaL_addstring(B *LuaL_Buffer, s string) {
	S1 := C.CString(s)
	defer C.free(unsafe.Pointer(S1))
	C.luaL_addstring((*C.luaL_Buffer)(B), S1)
}
func LuaL_addvalue(B *LuaL_Buffer)   { C.luaL_addvalue((*C.luaL_Buffer)(B)) }
func LuaL_pushresult(B *LuaL_Buffer) { C.luaL_pushresult((*C.luaL_Buffer)(B)) }
func LuaL_pushresultsize(B *LuaL_Buffer, sz uint) {
	C.luaL_pushresultsize((*C.luaL_Buffer)(B), C.size_t(sz))
}
func LuaL_buffinitsize(L Lua_Handle, B *LuaL_Buffer, sz uint) unsafe.Pointer {
	R := C.luaL_buffinitsize(L, (*C.luaL_Buffer)(B), C.size_t(sz))
	return unsafe.Pointer(R)
}

func LuaL_prepbuffer(B *LuaL_Buffer) { C.luaL_prepbuffsize((*C.luaL_Buffer)(B), C.LUAL_BUFFERSIZE) }

/* }====================================================== */

/*
** {======================================================
** File handles for IO library
** =======================================================
 */

/*
** A file handle is a userdata with metatable 'LUA_FILEHANDLE' and
** initial structure 'luaL_Stream' (it may contain other fields
** after that initial structure).
 */

const LUA_FILEHANDLE = C.LUA_FILEHANDLE

type LuaL_Stream C.luaL_Stream

/* }====================================================== */

// TODO: compat functions.
