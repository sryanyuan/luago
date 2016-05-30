package luago

/*
#cgo	CFLAGS:-I./inc
#cgo	LDFLAGS:-lm
#include <stdlib.h>
#include "lua/lua.h"
#include "lua/lualib.h"
#include "lua/lauxlib.h"
#include "luago.h"

#define	LUAGO_REG_LIB_NAME				script
#define	LUAGO_REG_FUNCS_LIST							\
	LUAGO_REG_YIELD_D_(script_WaitFrames,	"WaitFrames")	\
	LUAGO_REG_YIELD_D_(script_WaitSeconds,	"WaitSeconds")	\

#include "luago_define_reg_funcs.inc"

*/
import "C"

import (
	"unsafe"
)

//	LuaGo_Thread	status
const (
	THREAD_NOT_LOADED   = iota // 脚本未载入
	THREAD_LOADED              // 脚本已载入
	THREAD_RUNNING             // 脚本运行中
	THREAD_WAIT_SECONDS        // 脚本等待多少秒
	THREAD_WAIT_FRAMES         // 脚本等待多少帧
	THREAD_ERROR               // 脚本出现错误
	THREAD_DONE                // 脚本执行完毕
)

//	LuaGo_ThreadMgr
type LuaGo_ThreadMgr struct {
	LuaGo_ScriptBase
	threads map[*LuaGo_Thread]bool
}

func CreateLuaThreadMgr() *LuaGo_ThreadMgr {
	m := &LuaGo_ThreadMgr{}
	m.Create()
	return m
}

func (m *LuaGo_ThreadMgr) CreateThread(bAutoDelete bool) *LuaGo_Thread {
	t := new(LuaGo_Thread)

	t.handle = Lua_newthread(m.handle)
	t.auto_delete = bAutoDelete
	if nil == t.handle {
		t.handle = nil
		return nil
	}

	Lua_pushglobaltable(t.handle)
	Lua_pushthread(t.handle)
	Lua_pushlightuserdata(t.handle, unsafe.Pointer(t))
	Lua_settable(t.handle, -3)
	Lua_pop(t.handle, 1)

	m.threads[t] = true
	return t
}

func (m *LuaGo_ThreadMgr) DestroyThread(t *LuaGo_Thread) {
	if _, ok := m.threads[t]; ok {
		m.threads[t] = false
	}
}

func (m *LuaGo_ThreadMgr) IsValidThread(t *LuaGo_Thread) bool {
	todel, ok := m.threads[t]
	return ok && todel
}

func (m *LuaGo_ThreadMgr) GetThreadSum() int {
	return len(m.threads)
}

func (m *LuaGo_ThreadMgr) Create() {
	if nil != m.handle {
		m.Destroy()
	}
	m.handle = LuaL_newstate()
	m.threads = map[*LuaGo_Thread]bool{}
}

func (m *LuaGo_ThreadMgr) Destroy() {
	for k := range m.threads {
		delete(m.threads, k)
	}

	if nil != m.handle {
		Lua_close(m.handle)
	}
}

func (m *LuaGo_ThreadMgr) OpenScriptLib() {
	if nil != m.handle {
		LuaL_newlib(m.handle, LuaGo_RegPtr(C.LuaRegPtr_script()))
		Lua_setglobal(m.handle, "script")
	}
}

func (m *LuaGo_ThreadMgr) Update(dt float32) {
	dead := make([]*LuaGo_Thread, 0, len(m.threads))

	for t, v := range m.threads {
		if !v {
			dead = append(dead, t)
			continue
		}

		t.update(dt)
		switch t.GetStatus() {
		case THREAD_DONE, THREAD_ERROR, THREAD_NOT_LOADED:
			if t.auto_delete {
				dead = append(dead, t)
			}
		}
	}

	for i := range dead {
		delete(m.threads, dead[i])
	}
}

//	LuaGo_Thread
type LuaGo_Thread struct {
	LuaGo_ScriptBase
	status int
	mgr    *LuaGo_ThreadMgr

	auto_delete bool

	timestamp        float64 // current time
	timestamp_wakeup float64 // time to wake up
	frames_wakeup    int     // number of frames to wait
}

func (t *LuaGo_Thread) update(dt float32) {
	t.timestamp += float64(dt)

	switch t.status {
	case THREAD_WAIT_SECONDS:
		{ // 脚本等待多少秒
			if t.timestamp >= t.timestamp_wakeup {
				t.resume(false)
			}
		}
	case THREAD_WAIT_FRAMES:
		{ // 脚本等待多少帧
			t.frames_wakeup--
			if t.frames_wakeup <= 0 {
				t.resume(false)
			}
		}
	}
}

func (t *LuaGo_Thread) resume(bAbortWait bool) bool {
	switch t.status {
	case THREAD_NOT_LOADED:
		{
			return false
		}
	case THREAD_ERROR:
		{
			return false
		}
	}

	// we're about to run/resume the thread, so set the global
	t.status = THREAD_RUNNING

	// param is treated as a return value from the function that yielded
	Lua_pushboolean(t.handle, bAbortWait)

	switch Lua_resume(t.handle, nil, 1) {
	case LUA_OK:
		{
			t.status = THREAD_DONE
			return true
		}
	case LUA_YIELD:
		{
			return true
		}
		break
	default:
		{
			t.status = THREAD_ERROR
			t.err_msg = Lua_tostring(t.handle, -1)
			Lua_pop(t.handle, -1)
		}
	}

	return false
}

func (t *LuaGo_Thread) GetMgr() *LuaGo_ThreadMgr {
	return t.mgr
}

func (t *LuaGo_Thread) GetStatus() int {
	return t.status
}

func (t *LuaGo_Thread) GetAutoDelete() bool {
	return t.auto_delete
}

func (t *LuaGo_Thread) SetAutoDelete(bAutoDelete bool) {
	t.auto_delete = bAutoDelete
}

func (t *LuaGo_Thread) RunFile(file string) bool {
	t.status = THREAD_NOT_LOADED

	if LUA_OK == LuaL_loadfile(t.handle, file) {
		t.status = THREAD_LOADED
		return t.resume(false)
	}

	t.status = THREAD_NOT_LOADED
	t.err_msg = Lua_tostring(t.handle, -1)
	Lua_pop(t.handle, 1)

	return false
}

func (t *LuaGo_Thread) RunString(code string) bool {
	t.status = THREAD_NOT_LOADED

	if LUA_OK != LuaL_loadstring(t.handle, code) {
		t.err_msg = Lua_tostring(t.handle, -1)
		Lua_pop(t.handle, 1)
		return false
	}

	t.status = THREAD_LOADED

	return t.resume(false)
}

func (t *LuaGo_Thread) RunBuffer(buffer unsafe.Pointer, size uint) bool {
	t.status = THREAD_NOT_LOADED

	if LUA_OK != LuaL_loadbuffer(t.handle, buffer, size, "LuaGo_Thread.RunBuffer") {
		t.err_msg = Lua_tostring(t.handle, -1)
		Lua_pop(t.handle, 1)
		return false
	}

	t.status = THREAD_LOADED

	return t.resume(false)
}

func (t *LuaGo_Thread) AbortWait() {
	t.resume(true)
}

//
//	script Lib
//
func script_GetScriptObject(L Lua_Handle) *LuaGo_Thread {
	Lua_pushglobaltable(L)
	Lua_pushthread(L)
	Lua_gettable(L, -2)

	R := (*LuaGo_Thread)(Lua_touserdata(L, -1))

	Lua_pop(L, 2)
	return R
}

//export script_WaitFrames
func script_WaitFrames(l unsafe.Pointer) LuaGo_ResultSum {
	L := LuaGo_Handle(l)
	s := script_GetScriptObject(L)

	s.frames_wakeup = int(LuaL_optinteger(L, 1, 1))
	s.status = THREAD_WAIT_FRAMES

	return 0
}

//export script_WaitSeconds
func script_WaitSeconds(l unsafe.Pointer) LuaGo_ResultSum {
	L := LuaGo_Handle(l)
	s := script_GetScriptObject(L)

	s.timestamp_wakeup = s.timestamp + float64(LuaL_optnumber(L, 1, 1.0))
	s.status = THREAD_WAIT_SECONDS

	return 0
}
