package luago

/*
#cgo	CFLAGS:-I../../inc
#cgo	LDFLAGS:-lm
#include <stdlib.h>
#include "lua/lua.h"
*/
import "C"

import (
	"unsafe"
)

const LUA_VERSION_MAJOR = C.LUA_VERSION_MAJOR
const LUA_VERSION_MINOR = C.LUA_VERSION_MINOR
const LUA_VERSION_NUM = C.LUA_VERSION_NUM
const LUA_VERSION_RELEASE = C.LUA_VERSION_RELEASE

const LUA_VERSION = "Lua " + C.LUA_VERSION_MAJOR + "." + C.LUA_VERSION_MINOR
const LUA_RELEASE = LUA_VERSION + "." + C.LUA_VERSION_RELEASE
const LUA_COPYRIGHT = LUA_RELEASE + "  Copyright (C) 1994-2013 Lua.org, PUC-Rio"
const LUA_AUTHORS = C.LUA_AUTHORS

/* mark for precompiled code ('<esc>Lua') */
const LUA_SIGNATURE = C.LUA_SIGNATURE

/* option for multiple returns in 'lua_pcall' and 'lua_call' */
const LUA_MULTRET = C.LUA_MULTRET

/*
** pseudo-indices
 */
const LUA_REGISTRYINDEX = C.LUA_REGISTRYINDEX

func Lua_upvalueindex(i int) int {
	return LUA_REGISTRYINDEX - i
}

/* thread status */
const LUA_OK = C.LUA_OK
const LUA_YIELD = C.LUA_YIELD
const LUA_ERRRUN = C.LUA_ERRRUN
const LUA_ERRSYNTAX = C.LUA_ERRSYNTAX
const LUA_ERRMEM = C.LUA_ERRMEM
const LUA_ERRGCMM = C.LUA_ERRGCMM
const LUA_ERRERR = C.LUA_ERRERR

type Lua_Handle *C.lua_State

type Lua_CFunction C.lua_CFunction

/*
** functions that read/write blocks when loading/dumping Lua chunks
 */
type Lua_Reader C.lua_Reader

type Lua_Writer C.lua_Reader

/*
** prototype for memory-allocation functions
 */
type Lua_Alloc C.lua_Alloc

/*
** basic types
 */
const LUA_TNONE = C.LUA_TNONE

const LUA_TNIL = C.LUA_TNIL
const LUA_TBOOLEAN = C.LUA_TBOOLEAN
const LUA_TLIGHTUSERDATA = C.LUA_TLIGHTUSERDATA
const LUA_TNUMBER = C.LUA_TNUMBER
const LUA_TSTRING = C.LUA_TSTRING
const LUA_TTABLE = C.LUA_TTABLE
const LUA_TFUNCTION = C.LUA_TFUNCTION
const LUA_TUSERDATA = C.LUA_TUSERDATA
const LUA_TTHREAD = C.LUA_TTHREAD

const LUA_NUMTAGS = C.LUA_NUMTAGS

/* minimum Lua stack available to a C function */
const LUA_MINSTACK = C.LUA_MINSTACK

/* predefined values in the registry */
const LUA_RIDX_MAINTHREAD = C.LUA_RIDX_MAINTHREAD
const LUA_RIDX_GLOBALS = C.LUA_RIDX_GLOBALS
const LUA_RIDX_LAST = C.LUA_RIDX_LAST

/* type of numbers in Lua */
type Lua_Number C.lua_Number

/* type for integer functions */
type Lua_Integer C.lua_Integer

/* unsigned integer type */
type Lua_Unsigned C.lua_Unsigned

/*
** RCS ident string
 */
var Lua_ident = C.lua_ident

/*
** state manipulation
 */
func Lua_newstate(f Lua_Alloc, ud unsafe.Pointer) Lua_Handle { return C.lua_newstate(f, ud) }
func Lua_close(L Lua_Handle)                                 { C.lua_close(L) }
func Lua_newthread(L Lua_Handle) Lua_Handle                  { return C.lua_newthread(L) }

func Lua_atpanic(L Lua_Handle, panicf Lua_CFunction) Lua_CFunction {
	return Lua_CFunction(C.lua_atpanic(L, panicf))
}

func Lua_version(L Lua_Handle) Lua_Number { return Lua_Number(*C.lua_version(L)) }

/*
** basic stack manipulation
 */
func Lua_absindex(L Lua_Handle, idx int) int        { return int(C.lua_absindex(L, C.int(idx))) }
func Lua_gettop(L Lua_Handle) int                   { return int(C.lua_gettop(L)) }
func Lua_settop(L Lua_Handle, idx int)              { C.lua_settop(L, C.int(idx)) }
func Lua_pushvalue(L Lua_Handle, idx int)           { C.lua_pushvalue(L, C.int(idx)) }
func Lua_remove(L Lua_Handle, idx int)              { C.lua_remove(L, C.int(idx)) }
func Lua_insert(L Lua_Handle, idx int)              { C.lua_insert(L, C.int(idx)) }
func Lua_replace(L Lua_Handle, idx int)             { C.lua_replace(L, C.int(idx)) }
func Lua_copy(L Lua_Handle, fromidx int, toidx int) { C.lua_copy(L, C.int(fromidx), C.int(toidx)) }
func Lua_checkstack(L Lua_Handle, sz int) bool      { return 0 != C.lua_checkstack(L, C.int(sz)) }

func Lua_xmove(from Lua_Handle, to Lua_Handle, n int) { C.lua_xmove(from, to, C.int(n)) }

/*
** access functions (stack -> C)
 */

func Lua_isnumber(L Lua_Handle, idx int) bool    { return 0 != C.lua_isnumber(L, C.int(idx)) }
func Lua_isstring(L Lua_Handle, idx int) bool    { return 0 != C.lua_isstring(L, C.int(idx)) }
func Lua_iscfunction(L Lua_Handle, idx int) bool { return 0 != C.lua_iscfunction(L, C.int(idx)) }
func Lua_isuserdata(L Lua_Handle, idx int) bool  { return 0 != C.lua_isuserdata(L, C.int(idx)) }
func Lua_type(L Lua_Handle, idx int) int         { return int(C.lua_type(L, C.int(idx))) }
func Lua_typename(L Lua_Handle, tp int) string   { return C.GoString(C.lua_typename(L, C.int(tp))) }

func Lua_tonumberx(L Lua_Handle, idx int, isnum *bool) Lua_Number {
	var n C.int = 0
	R := C.lua_tonumberx(L, C.int(idx), &n)
	if isnum != nil {
		*isnum = (0 != n)
	}
	return Lua_Number(R)
}
func Lua_tointegerx(L Lua_Handle, idx int, isnum *bool) Lua_Integer {
	var n C.int = 0
	R := C.lua_tointegerx(L, C.int(idx), &n)
	if isnum != nil {
		*isnum = (0 != n)
	}
	return Lua_Integer(R)
}
func Lua_tounsignedx(L Lua_Handle, idx int, isnum *bool) Lua_Unsigned {
	var n C.int = 0
	R := C.lua_tounsignedx(L, C.int(idx), &n)
	if isnum != nil {
		*isnum = (0 != n)
	}
	return Lua_Unsigned(R)
}
func Lua_toboolean(L Lua_Handle, idx int) bool { return 0 != C.lua_toboolean(L, C.int(idx)) }
func Lua_tolstring(L Lua_Handle, idx int, len *uint) string {
	var sz C.size_t
	s := C.GoString(C.lua_tolstring(L, C.int(idx), &sz))
	if len != nil {
		*len = uint(sz)
	}
	return s
}
func Lua_tolstringbytes(L Lua_Handle, idx int) []byte {
	var sz C.size_t
	lstr := unsafe.Pointer(C.lua_tolstring(L, C.int(idx), &sz))
	data := C.GoBytes(lstr, C.int(sz))

	return data
}
func Lua_rawlen(L Lua_Handle, idx int) uint { return uint(C.lua_rawlen(L, C.int(idx))) }
func Lua_tocfunction(L Lua_Handle, idx int) Lua_CFunction {
	return Lua_CFunction(C.lua_tocfunction(L, C.int(idx)))
}
func Lua_touserdata(L Lua_Handle, idx int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_touserdata(L, C.int(idx)))
}
func Lua_tothread(L Lua_Handle, idx int) Lua_Handle {
	return Lua_Handle(C.lua_tothread(L, C.int(idx)))
}
func Lua_topointer(L Lua_Handle, idx int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_topointer(L, C.int(idx)))
}

/*
** Comparison and arithmetic functions
 */

const LUA_OPADD = C.LUA_OPADD
const LUA_OPSUB = C.LUA_OPSUB
const LUA_OPMUL = C.LUA_OPMUL
const LUA_OPDIV = C.LUA_OPDIV
const LUA_OPMOD = C.LUA_OPMOD
const LUA_OPPOW = C.LUA_OPPOW
const LUA_OPUNM = C.LUA_OPUNM

func Lua_arith(L Lua_Handle, op int) { C.lua_arith(L, C.int(op)) }

const LUA_OPEQ = C.LUA_OPEQ
const LUA_OPLT = C.LUA_OPLT
const LUA_OPLE = C.LUA_OPLE

func Lua_rawequal(L Lua_Handle, idx1 int, idx2 int) bool {
	return 0 != C.lua_rawequal(L, C.int(idx1), C.int(idx2))
}
func Lua_compare(L Lua_Handle, idx1 int, idx2 int, op int) int {
	R := C.lua_compare(L, C.int(idx1), C.int(idx2), C.int(op))
	return int(R)
}

/*
** push functions (C -> stack)
 */
func Lua_pushnil(L Lua_Handle)                      { C.lua_pushnil(L) }
func Lua_pushnumber(L Lua_Handle, n Lua_Number)     { C.lua_pushnumber(L, C.lua_Number(n)) }
func Lua_pushinteger(L Lua_Handle, n Lua_Integer)   { C.lua_pushinteger(L, C.lua_Integer(n)) }
func Lua_pushunsigned(L Lua_Handle, n Lua_Unsigned) { C.lua_pushunsigned(L, C.lua_Unsigned(n)) }
func Lua_pushlstring(L Lua_Handle, s unsafe.Pointer, l uint) {
	C.lua_pushlstring(L, (*C.char)(s), C.size_t(l))
}
func Lua_pushstring(L Lua_Handle, s string) string {
	S1 := C.CString(s)
	defer C.free(unsafe.Pointer(S1))
	C.lua_pushstring(L, S1)
	return s
}

// TODO: LUA_API const char *(lua_pushvfstring)	(L Lua_Handle, const char *fmt, va_list argp);
// TODO: LUA_API const char *(lua_pushfstring) (L Lua_Handle, const char *fmt, ...);
func Lua_pushcclosure(L Lua_Handle, fn Lua_CFunction, n int) {
	C.lua_pushcclosure(L, C.lua_CFunction(fn), C.int(n))
}
func Lua_pushboolean(L Lua_Handle, b bool) {
	var n C.int
	if b {
		n = 1
	} else {
		n = 0
	}
	C.lua_pushboolean(L, n)
}
func Lua_pushlightuserdata(L Lua_Handle, p unsafe.Pointer) { C.lua_pushlightuserdata(L, p) }
func Lua_pushthread(L Lua_Handle) bool                     { return 0 != C.lua_pushthread(L) }

/*
** get functions (Lua -> stack)
 */
func Lua_getglobal(L Lua_Handle, varname string) {
	S1 := C.CString(varname)
	defer C.free(unsafe.Pointer(S1))
	C.lua_getglobal(L, S1)
}
func Lua_gettable(L Lua_Handle, idx int) { C.lua_gettable(L, C.int(idx)) }
func Lua_getfield(L Lua_Handle, idx int, k string) {
	S1 := C.CString(k)
	defer C.free(unsafe.Pointer(S1))
	C.lua_getfield(L, C.int(idx), S1)
}
func Lua_rawget(L Lua_Handle, idx int)                    { C.lua_rawget(L, C.int(idx)) }
func Lua_rawgeti(L Lua_Handle, idx int, n int)            { C.lua_rawgeti(L, C.int(idx), C.int(n)) }
func Lua_rawgetp(L Lua_Handle, idx int, p unsafe.Pointer) { C.lua_rawgetp(L, C.int(idx), p) }
func Lua_createtable(L Lua_Handle, narr int, nrec int) {
	C.lua_createtable(L, C.int(narr), C.int(nrec))
}
func Lua_newuserdata(L Lua_Handle, sz uint) unsafe.Pointer {
	return C.lua_newuserdata(L, C.size_t(sz))
}
func Lua_getmetatable(L Lua_Handle, objindex int) bool {
	return 0 != C.lua_getmetatable(L, C.int(objindex))
}
func Lua_getuservalue(L Lua_Handle, idx int) { C.lua_getuservalue(L, C.int(idx)) }

/*
** set functions (stack -> Lua)
 */
func Lua_setglobal(L Lua_Handle, varname string) {
	S1 := C.CString(varname)
	defer C.free(unsafe.Pointer(S1))
	C.lua_setglobal(L, S1)
}
func Lua_settable(L Lua_Handle, idx int) { C.lua_settable(L, C.int(idx)) }
func Lua_setfield(L Lua_Handle, idx int, k string) {
	S1 := C.CString(k)
	defer C.free(unsafe.Pointer(S1))
	C.lua_setfield(L, C.int(idx), S1)
}
func Lua_rawset(L Lua_Handle, idx int)                    { C.lua_rawset(L, C.int(idx)) }
func Lua_rawseti(L Lua_Handle, idx int, n int)            { C.lua_rawseti(L, C.int(idx), C.int(n)) }
func Lua_rawsetp(L Lua_Handle, idx int, p unsafe.Pointer) { C.lua_rawsetp(L, C.int(idx), p) }
func Lua_setmetatable(L Lua_Handle, objindex int) bool {
	return 0 != C.lua_setmetatable(L, C.int(objindex))
}
func Lua_setuservalue(L Lua_Handle, idx int) { C.lua_setuservalue(L, C.int(idx)) }

/*
** 'load' and 'call' functions (load and run Lua code)
 */
func Lua_callk(L Lua_Handle, nargs int, nresults int, ctx int, k Lua_CFunction) {
	C.lua_callk(L, C.int(nargs), C.int(nresults), C.int(ctx), C.lua_CFunction(k))
}
func Lua_call(L Lua_Handle, nargs int, nresults int) {
	C.lua_callk(L, C.int(nargs), C.int(nresults), 0, C.lua_CFunction(nil))
}
func Lua_getctx(L Lua_Handle, ctx *int) int {
	var c C.int
	R := C.lua_getctx(L, &c)
	if ctx != nil {
		*ctx = int(c)
	}
	return int(R)
}
func Lua_pcallk(L Lua_Handle, nargs int, nresults int, errfunc int, ctx int, k Lua_CFunction) int {
	R := C.lua_pcallk(L, C.int(nargs), C.int(nresults), C.int(errfunc), C.int(ctx), C.lua_CFunction(k))
	return int(R)
}
func Lua_pcall(L Lua_Handle, nargs int, nresults int, errfunc int) int {
	R := C.lua_pcallk(L, C.int(nargs), C.int(nresults), C.int(errfunc), C.int(0), C.lua_CFunction(nil))
	return int(R)
}
func Lua_load(L Lua_Handle, reader Lua_Reader, dt unsafe.Pointer, chunkname string, mode string) int {
	S1 := C.CString(chunkname)
	defer C.free(unsafe.Pointer(S1))
	S2 := C.CString(mode)
	defer C.free(unsafe.Pointer(S2))
	R := C.lua_load(L, C.lua_Reader(reader), dt, S1, S2)
	return int(R)
}
func Lua_dump(L Lua_Handle, writer Lua_Writer, data unsafe.Pointer) int {
	R := C.lua_dump(L, C.lua_Writer(writer), data)
	return int(R)
}

/*
** coroutine functions
 */
func Lua_yieldk(L Lua_Handle, nresults int, ctx int, k Lua_CFunction) int {
	R := C.lua_yieldk(L, C.int(nresults), C.int(ctx), C.lua_CFunction(k))
	return int(R)
}
func Lua_yield(L Lua_Handle, nresults int) int {
	R := C.lua_yieldk(L, C.int(nresults), C.int(0), C.lua_CFunction(nil))
	return int(R)
}
func Lua_resume(L Lua_Handle, from Lua_Handle, narg int) int {
	R := C.lua_resume(L, from, C.int(narg))
	return int(R)
}
func Lua_status(L Lua_Handle) int { R := C.lua_status(L); return int(R) }

/*
** garbage-collection function and options
 */

const LUA_GCSTOP = C.LUA_GCSTOP
const LUA_GCRESTART = C.LUA_GCRESTART
const LUA_GCCOLLECT = C.LUA_GCCOLLECT
const LUA_GCCOUNT = C.LUA_GCCOUNT
const LUA_GCCOUNTB = C.LUA_GCCOUNTB
const LUA_GCSTEP = C.LUA_GCSTEP
const LUA_GCSETPAUSE = C.LUA_GCSETPAUSE
const LUA_GCSETSTEPMUL = C.LUA_GCSETSTEPMUL
const LUA_GCSETMAJORINC = C.LUA_GCSETMAJORINC
const LUA_GCISRUNNING = C.LUA_GCISRUNNING
const LUA_GCGEN = C.LUA_GCGEN
const LUA_GCINC = C.LUA_GCINC

func Lua_gc(L Lua_Handle, what int, data int) int {
	R := C.lua_gc(L, C.int(what), C.int(data))
	return int(R)
}

/*
** miscellaneous functions
 */

func Lua_error(L Lua_Handle) int { R := C.lua_error(L); return int(R) }

func Lua_next(L Lua_Handle, idx int) int { R := C.lua_next(L, C.int(idx)); return int(R) }

func Lua_concat(L Lua_Handle, n int) { C.lua_concat(L, C.int(n)) }
func Lua_len(L Lua_Handle, idx int)  { C.lua_len(L, C.int(idx)) }

func Lua_getallocf(L Lua_Handle, ud *unsafe.Pointer) Lua_Alloc {
	R := C.lua_getallocf(L, ud)
	return Lua_Alloc(R)
}
func Lua_setallocf(L Lua_Handle, f Lua_Alloc, ud unsafe.Pointer) {
	C.lua_setallocf(L, C.lua_Alloc(f), ud)
}

/*
** ===============================================================
** some useful macros
** ===============================================================
 */

func Lua_tonumber(L Lua_Handle, idx int) Lua_Number { return Lua_tonumberx(L, idx, nil) }
func Lua_tointeger(L Lua_Handle, idx int) Lua_Integer {
	R := C.lua_tointegerx(L, C.int(idx), nil)
	return Lua_Integer(R)
}
func Lua_tounsigned(L Lua_Handle, idx int) Lua_Unsigned {
	R := C.lua_tounsignedx(L, C.int(idx), nil)
	return Lua_Unsigned(R)
}

func Lua_pop(L Lua_Handle, n int) { C.lua_settop(L, C.int(-(n)-1)) }

func Lua_newtable(L Lua_Handle) { C.lua_createtable(L, C.int(0), C.int(0)) }

func Lua_register(L Lua_Handle, varname string, fn Lua_CFunction) {
	S1 := C.CString(varname)
	defer C.free(unsafe.Pointer(S1))
	Lua_pushcfunction(L, fn)
	C.lua_setglobal(L, S1)
}

func Lua_pushcfunction(L Lua_Handle, fn Lua_CFunction) {
	C.lua_pushcclosure(L, C.lua_CFunction(fn), C.int(0))
}

func Lua_isfunction(L Lua_Handle, n int) bool { return (C.lua_type(L, C.int(n)) == LUA_TFUNCTION) }
func Lua_istable(L Lua_Handle, n int) bool    { return (C.lua_type(L, C.int(n)) == LUA_TTABLE) }
func Lua_islightuserdata(L Lua_Handle, n int) bool {
	return (C.lua_type(L, C.int(n)) == LUA_TLIGHTUSERDATA)
}
func Lua_isnil(L Lua_Handle, n int) bool       { return (C.lua_type(L, C.int(n)) == LUA_TNIL) }
func Lua_isboolean(L Lua_Handle, n int) bool   { return (C.lua_type(L, C.int(n)) == LUA_TBOOLEAN) }
func Lua_isthread(L Lua_Handle, n int) bool    { return (C.lua_type(L, C.int(n)) == LUA_TTHREAD) }
func Lua_isnone(L Lua_Handle, n int) bool      { return (C.lua_type(L, C.int(n)) == LUA_TNONE) }
func Lua_isnoneornil(L Lua_Handle, n int) bool { return (C.lua_type(L, C.int(n)) <= 0) }

// TODO: #define lua_pushliteral			(L Lua_Handle, s string)		lua_pushlstring(L, "" s, (sizeof(s)/sizeof(char))-1)
func Lua_pushglobaltable(L Lua_Handle) {
	C.lua_rawgeti(L, C.int(LUA_REGISTRYINDEX), C.int(LUA_RIDX_GLOBALS))
}

func Lua_tostring(L Lua_Handle, i int) string {
	s := C.GoString(C.lua_tolstring(L, C.int(i), nil))
	return s
}

/*
** {======================================================================
** Debug API
** =======================================================================
 */

/*
** Event codes
 */
const LUA_HOOKCALL = C.LUA_HOOKCALL
const LUA_HOOKRET = C.LUA_HOOKRET
const LUA_HOOKLINE = C.LUA_HOOKLINE
const LUA_HOOKCOUNT = C.LUA_HOOKCOUNT
const LUA_HOOKTAILCALL = C.LUA_HOOKTAILCALL

/*
** Event masks
 */
const LUA_MASKCALL = C.LUA_MASKCALL
const LUA_MASKRET = C.LUA_MASKRET
const LUA_MASKLINE = C.LUA_MASKLINE
const LUA_MASKCOUNT = C.LUA_MASKCOUNT

// TODO: debug not support.

type Lua_Debug C.lua_Debug

type Lua_Hook C.lua_Hook
