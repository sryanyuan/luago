#ifndef _INC_LUAGOEXT_
#define _INC_LUAGOEXT_
//////////////////////////////////////////////////////////
struct	lua_State;

//	创建元方法的元表
int luago_open(struct lua_State* _pLuaState, void* _pLuaGoState);
int luago_close(struct lua_State* L);

//	push GoFunction
int luago_pushGoFunction(struct lua_State* L, const char* _pszName, int _nSeq);
int luago_popGoFunction(struct lua_State* L, const char* _pszName);

//	helpers
void luago_error(struct lua_State* L, const char* _pszErrMsg);
void luago_addSearchPath(struct lua_State* L, const char* _pszPath, int _tip);
//////////////////////////////////////////////////////////
#endif