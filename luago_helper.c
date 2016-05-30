#include <string.h>
#include <stdlib.h>
#include "./inc/lua/lua.h"
#include "./inc/lua/lauxlib.h"
#include "luago_helper.h"

int	luaGo_GetRef(struct lua_State* L, int nVariableReference)
{
	if(NULL == L)
	{
		return	0;
	}
	
	lua_rawgeti(L, LUA_REGISTRYINDEX, nVariableReference);

	if(lua_isnoneornil(L, -1))
	{
		lua_pop(L, 1);
		return	0;
	}

	return 1;
}

int	luaGo_GetGlobal(struct lua_State* L, const char* sVariableName)
{
	if(NULL == L || NULL == sVariableName || 0 == *sVariableName)
	{
		return	0;
	}

	char	tmpBuffer[256];
	// sVariableName is split by dots '.' in arrays and subarrays
	// the nextVar sVariableName contains a pointer to the next part to proceed
	const char* nextVar = sVariableName;

	do
	{
		// since we are going to modify nextVar, we store its value here
		const char* currentVar = nextVar;

		// first we extract the part between currentVar and the next dot we encounter
		nextVar = strchr(currentVar, '.');
		memset(tmpBuffer, 0, sizeof(tmpBuffer));
		if(NULL == nextVar)
		{
			strncpy(tmpBuffer, currentVar, sizeof(tmpBuffer));
		}
		else
		{
			strncpy(tmpBuffer, currentVar, nextVar - currentVar);
			// since nextVar is pointing to a dot, we have to increase
			// it first in order to find the next sVariableName
			++nextVar;
		}

		// ask lua to find the part stored in buffer
		// if currentVar == begin, this is a global sVariableName and push it on the stack
		//   otherwise we already have an array pushed on the stack by the previous loop
		if (currentVar == sVariableName)
		{
			lua_getglobal(L, tmpBuffer);
		}
		else
		{
			// if sVariableName is "a.b" and "a" is not a table (eg. it's a number or a string), this happens
			// we don't have a specific exception for this, we consider this as a sVariableName-doesn't-exist
			if (!lua_istable(L, -1))
			{
				lua_pop(L, 1);
				return	0;
			}

			// replacing the current table in the stack by its member
			lua_pushstring(L, tmpBuffer);
			lua_gettable(L, -2);
			lua_remove(L, -2);
		}

		// lua will accept anything as sVariableName name, but if the sVariableName doesn't exist
		//   it will simply push "nil" instead of a value
		// so if we have a nil on the stack, the sVariableName didn't exist and we throw
		if (lua_isnoneornil(L, -1))
		{
			lua_pop(L, 1);
			return	0;
		}
	} while (NULL != nextVar && *nextVar != 0);

	return	1;
}

int	luaGo_SetGlobal(struct lua_State* L, const char* sVariableName)
{
	// making sure there's something on the stack (ie. the value to set)
	if(NULL == L || NULL == sVariableName || 0 == *sVariableName || lua_gettop(L) < 1)
	{
		return	0;
	}

	// two possibilities: either "sVariableName" is a global sVariableName, or a member of an array
	const char* lastDot = strrchr(sVariableName, '.');
	if (lastDot == NULL)
	{
		// this is the first case, we simply call setglobal (which cleans the stack)
		lua_setglobal(L, sVariableName);
	}
	else 
	{
		size_t	tableNameLen= lastDot - sVariableName;
		char*	tableName	= (char*)malloc(tableNameLen + 1);
		memcpy(tableName, sVariableName, tableNameLen);
		tableName[tableNameLen]	= '\0';
		// in the second case, we call _getGlobal on the table name
		if(!luaGo_GetGlobal(L, tableName))
		{
			free(tableName);
			return	0;
		}

		free(tableName);
		if (!lua_istable(L, -1))
		{
			lua_pop(L, 1);
			return	0;
		}

		// now we have our value at -2 (was pushed before _setGlobal is called) and our table at -1
		lua_pushstring(L, ++lastDot);		// value at -3, table at -2, key at -1
		lua_pushvalue(L, -3);				// value at -4, table at -3, key at -2, value at -1
		lua_settable(L, -3);				// value at -2, table at -1
		lua_pop(L, 2);						// stack empty \o/
	}
	
	return	1;
}
