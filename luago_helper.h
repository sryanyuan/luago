#if		!defined(__LUAGO_HELPER_H_0X0311_20130626_)
#define	__LUAGO_HELPER_H_0X0311_20130626_

struct	lua_State;

int		luaGo_GetGlobal(struct lua_State* L, const char* sVariableName);
int		luaGo_SetGlobal(struct lua_State* L, const char* sVariableName);

int		luaGo_GetRef(struct lua_State* L, int nVariableReference);

#endif
