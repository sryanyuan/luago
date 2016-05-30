
/*

#define	LUAGO_REG_LIB_NAME				XXXX
#define	LUAGO_REG_FUNCS_LIST			\
	LUAGO_REG(GoFunc1)					\
	LUAGO_REG_(GoFunc1, "GoFunc1_1")	\
	LUAGO_REG_D(GoFunc2)				\
	LUAGO_REG_D_(GoFunc2, "GoFunc2_1")	\

#include	"inc/lua/luago_define_reg_funcs.inc"

*/

#include	"luago.h"

/* predefine check */
#if	!defined(LUAGO_REG_LIB_NAME)	||	!defined(LUAGO_REG_FUNCS_LIST)
#	error please define LUAGO_REG_LIB_NAME and LUAGO_REG_FUNCS_LIST
#endif

#if	defined(LUAGO_REG) || defined(LUAGO_REG_) || defined(LUAGO_REG_D) || defined(LUAGO_REG_D_)
#	error LUAGO_REG or LUAGO_REG_ or LUAGO_REG_D or LUAGO_REG_D_ has been defined
#endif


/* declare */
#define	LUAGO_REG(X)			
#define	LUAGO_REG_(X,N)			
#define	LUAGO_REG_YIELD(X)		
#define	LUAGO_REG_YIELD_(X,N)	
#define	LUAGO_REG_D(X)			LUAGO_CFUNC_DECLARE(X)
#define	LUAGO_REG_D_(X,N)		LUAGO_CFUNC_DECLARE(X)
#define	LUAGO_REG_YIELD_D(X)	LUAGO_YIELD_CFUNC_DECLARE(X)
#define	LUAGO_REG_YIELD_D_(X,N)	LUAGO_YIELD_CFUNC_DECLARE(X)

LUAGO_REG_FUNCS_LIST

/* reg funcs wrraper */
#undef	LUAGO_REG
#undef	LUAGO_REG_
#undef	LUAGO_REG_YIELD
#undef	LUAGO_REG_YIELD_
#undef	LUAGO_REG_D
#undef	LUAGO_REG_D_
#undef	LUAGO_REG_YIELD_D
#undef	LUAGO_REG_YIELD_D_

#define	LUAGO_REG(X)			LUAGO_REG_ITEM(X)
#define	LUAGO_REG_(X,N)			LUAGO_REG_ITEM_(X, N)
#define	LUAGO_REG_YIELD(X)		LUAGO_REG_YIELD_ITEM(X)
#define	LUAGO_REG_YIELD_(X,N)	LUAGO_REG_YIELD_ITEM_(X, N)
#define	LUAGO_REG_D(X)			LUAGO_REG_ITEM(X)
#define	LUAGO_REG_D_(X,N)		LUAGO_REG_ITEM_(X, N)
#define	LUAGO_REG_YIELD_D(X)	LUAGO_REG_YIELD_ITEM(X)
#define	LUAGO_REG_YIELD_D_(X,N)	LUAGO_REG_YIELD_ITEM_(X, N)

#define	LUAGO_I_REG_WRAPPER_NAME(N)			LUAGO_I_REG_WRAPPER_NAME_VALUE(N)
#define	LUAGO_I_REG_WRAPPER_NAME_VALUE(N)	LuaReg_##N

#define	LUAGO_I_REGPTR_WRAPPER_NAME(N)			LUAGO_I_REGPTR_WRAPPER_NAME_VALUE(N)
#define	LUAGO_I_REGPTR_WRAPPER_NAME_VALUE(N)	LuaRegPtr_##N

static	luaL_Reg*	LUAGO_I_REG_WRAPPER_NAME(LUAGO_REG_LIB_NAME) (){
	static	luaL_Reg	funcs[]	= {
		LUAGO_REG_FUNCS_LIST
		{NULL, NULL},
	};

	return	funcs;
}

static	void*		LUAGO_I_REGPTR_WRAPPER_NAME(LUAGO_REG_LIB_NAME) (){
	return	LUAGO_I_REG_WRAPPER_NAME(LUAGO_REG_LIB_NAME) ();
}

#undef	LUAGO_I_REGPTR_WRAPPER_NAME
#undef	LUAGO_I_REGPTR_WRAPPER_NAME_VALUE

#undef	LUAGO_I_REG_WRAPPER_NAME
#undef	LUAGO_I_REG_WRAPPER_NAME_VALUE

/* clean */
#undef	LUAGO_REG_LIB_NAME
#undef	LUAGO_REG_FUNCS_LIST
#undef	LUAGO_REG
#undef	LUAGO_REG_
#undef	LUAGO_REG_YIELD
#undef	LUAGO_REG_YIELD_
#undef	LUAGO_REG_D
#undef	LUAGO_REG_D_
#undef	LUAGO_REG_YIELD_D
#undef	LUAGO_REG_YIELD_D_
