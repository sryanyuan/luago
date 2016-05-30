#include <string.h>
#include <stdlib.h>
#include "./inc/lua/lua.h"
#include "./inc/lua/lauxlib.h"
#include "luagoext.h"
#include "lpack.h"
//////////////////////////////////////////////////////////
//	go export
extern int luagoc_call(void*, int);
//////////////////////////////////////////////////////////
//	该表存着GoFunction Index -> LuaGo_State的表
const char* g_szGoStateTable = "gostate_table";
//	元方法的元表
const char* g_szGoMetaTable = "gometatable";
//	全局表中的golua_state key
const char* g_szGoLuaStateKey = "golua_state_key";
//////////////////////////////////////////////////////////
//	得到lua_state对应的golua_state
static void* luago_getGoLuaState(struct lua_State* L)
{
	lua_getglobal(L, g_szGoLuaStateKey);
	void* pRet = 0;
	
	if(lua_isnil(L, -1))
	{
		//	nothing
		luagoc_output("invalid GoLua_State");
	}
	else
	{
		pRet = lua_touserdata(L, -1);
	}
	lua_pop(L, 1);

	return pRet;
}

//	设置lua_state对应的golua_state
static void luago_setGoLuaState(struct lua_State* L, void* _pLuaGoState)
{
	lua_pushlightuserdata(L, _pLuaGoState);
	lua_setglobal(L, g_szGoLuaStateKey);
}

//	元方法
static int luago_metamethod_call(struct lua_State* L)
{
	//	stack userdata, parameters...
	if(lua_isuserdata(L, 1))
	{
		//	check valid
		int* pUserData = (int*)luaL_checkudata(L, 1, g_szGoMetaTable);
		if(0 != pUserData)
		{
			int nGoFunctionSeq = *pUserData;

			// get the golua_state
			void* pLuaGoState = luago_getGoLuaState(L);
			//	call go callback
			return (int)luagoc_call(pLuaGoState, nGoFunctionSeq);
		}
		else
		{
			luaL_argcheck(L, pUserData != 0, 1, "GoFunction expected");
		}
	}
	
	luago_error(L, "trying to call a non-callable object");
	return 0;
}

//	初始化metatable内的原方法
static void luago_buildMetatable(struct lua_State* L)
{
	//	stack: mt
	lua_pushstring(L, "__call");
	lua_pushcfunction(L, luago_metamethod_call);
	lua_rawset(L, -3);
}

//	创建新的userdata
static int luago_registerGoFunction(struct lua_State* _pLuaState, int _nSeq)
{
	struct lua_State* L = _pLuaState;
	
	void* pUserData = lua_newuserdata(_pLuaState, sizeof(int));
	int* pIntUserData = (int*)pUserData;
	*pIntUserData = _nSeq;

	//	set metatable
	//	stack: ud
	//lua_pushvalue(L, -1);	//	stack: ud ud
	luaL_getmetatable(L, g_szGoMetaTable);
	if(lua_isnil(L, -1))
	{
		//	??? metatable not found
		lua_pop(L, 2);		//	stack: ud
		return 0;
	}

	//	stack: ud mt
	lua_setmetatable(L, -2);	//	stack: ud
	
	return 1;
}
//////////////////////////////////////////////////////////
//	创建GoFunction的metatable
int luago_open(struct lua_State* _pLuaState, void* _pLuaGoState)
{
	int r = luaL_newmetatable(_pLuaState, g_szGoMetaTable);
	
	//	lua 5.1
#ifdef LUA_VERSION_NUM
	if(r)
	{
		//	stack: mt
		lua_pushvalue(_pLuaState, -1);
		//	stack: mt mt
		lua_pushstring(_pLuaState, g_szGoMetaTable);
		//	stack: mt mt value
		lua_settable(_pLuaState, LUA_REGISTRYINDEX);
		//	reg[mt] = value
	}
#endif
	//	stack: mt
	if(r)
	{
		//	lpack support
		luaopen_pack(_pLuaState);
		luago_buildMetatable(_pLuaState);	
		luago_setGoLuaState(_pLuaState, _pLuaGoState);
	}
	
	return r;
}

int luago_close(struct lua_State* L)
{
	//
	return 1;
}

int luago_registerLuaState(struct lua_State* _pLuaState, void* _pLuaGoState, int _nSeq)
{
	lua_getglobal(_pLuaState, g_szGoStateTable);
	if(lua_istable(_pLuaState, -1))
	{
		//	already exists
		//	stack: gostatetable
	}
	else
	{
		//	create it
		//	stack: nil
		lua_pop(_pLuaState, 1);
		lua_newtable(_pLuaState);
		//	stack: gostatetable
		lua_setglobal(_pLuaState, g_szGoStateTable);

		//	push the global table into stack
		lua_getglobal(_pLuaState, g_szGoStateTable);
	}
	
	//	stack: gostatetable
	if(!lua_istable(_pLuaState, -1))
	{
		return 0;
	}

	//	push key
	lua_pushnumber(_pLuaState, _nSeq);
	
	//	push value
	if(0 == _pLuaGoState)
	{
		lua_pushnil(_pLuaState);
	}
	else
	{
		lua_pushlightuserdata(_pLuaState, _pLuaGoState);	
	}
	
	//	stack: gostatetable seq luagostate
	lua_rawset(_pLuaState, -3);
	
	//	gostatetable: [seq] = luagostate
	
	return 1;
}

int luago_unregisterLuaState(struct lua_State* _pLuaState, int _nSeq)
{
	return luago_registerLuaState(_pLuaState, 0, _nSeq);
}


//	push GoFunction
int luago_pushGoFunction(struct lua_State* L, const char* _pszName, int _nSeq)
{
	//	new userdata	stack: ud
	luago_registerGoFunction(L, _nSeq);
	//	全局表保持着对GoFunction userdata的引用，确保不被GC
	lua_setglobal(L, _pszName);
	return 1;
}

int luago_popGoFunction(struct lua_State* L, const char* _pszName)
{
	lua_pushnil(L);
	lua_setglobal(L, _pszName);
	return 1;
}

//	输出错误信息
void luago_error(struct lua_State* L, const char* _pszErrMsg)
{
	luaL_where(L, 1);
	lua_pushstring(L, _pszErrMsg);
	lua_concat(L, 2);
	lua_error(L);
}

//	添加lua模块搜寻目录
void luago_addSearchPath(struct lua_State* L, const char* _pszPath, int _tip)
{
	lua_getglobal(L, "package");	//	stack:package
	lua_getfield(L, -1, "path");		//	stack:package package[path]
	
	const char* pszCurrentPath = lua_tostring(L, -1);
	char szPath[512] = {0};
	strcpy(szPath, pszCurrentPath);
	strcat(szPath, ";");
	strcat(szPath, _pszPath);
	strcat(szPath, "\\?.lua");
	
	if(0 != _tip)
	{
		luagoc_output("lua file search path:");
		luagoc_output(szPath);	
	}
	
	lua_pushstring(L, szPath);	//	stack:package package[path] newpath
	lua_setfield(L, -3, "path");		//	stack:package package[path]
	
	lua_pop(L, 2);						//	stack: -
}
