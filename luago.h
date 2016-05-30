#if		!defined(__LUAGO_H_0X0311_20130626_)
#define	__LUAGO_H_0X0311_20130626_

//func	MyGoCFunc(L unsafe.Pointer)int;					//	32 only
//func	MyGoCFunc(L unsafe.Pointer)LuaGo_ResultSum;		//	32/64
#define	LUAGO_CFUNC_DECLARE(GF)			extern	int		(GF)(void* L);
#define	LUAGO_CFUNC_EXPORT(GF)			static	lua_CFunction	LuaCFunc_##GF()		{ LUAGO_CFUNC_DECLARE(GF); return (lua_CFunction)(GF);	}

#define	LUAGO_YIELD_CFUNC_DECLARE(GF)	static	int	Yield_##GF(lua_State* l)			{ extern int (GF)(void* L); int r = (GF)(l); return lua_yield(l, r);	}
#define	LUAGO_YIELD_CFUNC_EXPORT(GF)	static	lua_CFunction	LuaCFunc_Yield_##GF()	{ extern int (Yield_##GF)(lua_State* l); return Yield_##GF;	}

//func	MyGoReader(L unsafe.Pointer, ud unsafe.Pointer, sz *C.size_t)unsafe.Pointer;
#define	LUAGO_READER_DECLARE(GF)		extern	void*	(GF)(void* L, void *ud, size_t *sz);
#define	LUAGO_READER_EXPORT(GF)			static	lua_Reader		LuaReader_##GF()	{ LUAGO_READER_DECLARE(GF); return (lua_Reader)(GF);	}

//func	MyGoWriter(L unsafe.Pointer, p unsafe.Pointer, sz C.size_t, ud unsafe.Pointer)int;			//	32
//func	MyGoWriter(L unsafe.Pointer, p unsafe.Pointer, sz C.size_t, ud unsafe.Pointer)LuaGo_Int;	//	32/64
#define	LUAGO_WRITER_DECLARE(GF)		extern	int		(GF)(void* L, void *p, size_t sz, void* ud);
#define	LUAGO_WRITER_EXPORT(GF)			static	lua_Writer		LuaWriter_##GF()	{ LUAGO_WRITER_DECLARE(GF); return (lua_Writer)(GF);	}

//func	MyGoAlooc(ud unsafe.Pointer, ptr unsafe.Pointer, osize C.size_t, nsize C.size_t)unsafe.Pointer;
#define	LUAGO_ALLOC_DECLARE(GF)			extern	void*	(GF)(void *ud, void *ptr, size_t osize, size_t nsize);
#define	LUAGO_ALLOC_EXPORT(GF)			static	lua_Alloc		LuaAlloc_##GF()		{ LUAGO_ALLOC_DECLARE(GF); return (lua_Alloc)(GF);	}

//func	MyGoHook(L unsafe.Pointer, ar unsafe.Pointer);
#define	LUAGO_HOOK_DECLARE(GF)			extern	void	(GF)(void *L, void *ar);
#define	LUAGO_HOOK_EXPORT(GF)			static	lua_Hook		LuaHook_##GF()		{ LUAGO_HOOK_DECLARE(GF); return (lua_Hook)(GF);	}

// used by luaL_setfuncs
#define	LUAGO_REG_BEGIN(M)				\
static	luaL_Reg*	LuaReg_##M(){		\
	static	luaL_Reg	funcs[]	= {

#define	LUAGO_REG_ITEM(F)				{#F, (lua_CFunction)(F)},
#define	LUAGO_REG_ITEM_(F, N)			{(N), (lua_CFunction)(F)},

#define	LUAGO_REG_YIELD_ITEM(F)			{#F, Yield_##F},
#define	LUAGO_REG_YIELD_ITEM_(F, N)		{(N), Yield_##F},

#define	LUAGO_REG_END()					\
		{NULL, NULL},					\
	};									\
										\
	return	funcs;						\
}

/*

//
//	#include "inc/lua/lua.h"
//	#include "inc/lua/lualib.h"
//	#include "inc/lua/lauxlib.h"
//	#include "inc/lua/luago.h"
//	
//	static	void	UseGoFunc(lua_CFunction func)
//	{
//		lua_State*	l	= 0;
//		func(l);
//	}
//	static	void	UseSetFuncs(luaL_Reg* funcs)
//	{
//	}
//	LUAGO_CFUNC_DECLARE(GoFunc)
//	LUAGO_CFUNC_EXPORT(GoFunc)
//	
//	LUAGO_REG_BEGIN(XXX)
//		LUAGO_REG_ITEM(GoFunc)
//	LUAGO_REG_END()
//
import "C"

import (
	"fmt"
	"unsafe"
)

//export GoFunc
func	GoFunc(l unsafe.Pointer)int{
	fmt.Println("From GoFunc");
	return	1;
}

func	main(){
	GoFunc(nil);

	C.UseGoFunc(C.LuaCFunc_GoFunc());

	C.UseSetFuncs(C.LuaReg_XXX());
}

*/

#endif