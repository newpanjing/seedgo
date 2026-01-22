package shared


// // 只需要传入 DB，自动搞定 Logic 和 Hook 的绑定
// func InitBase[T any,L any](hook HandlerHook[T], logic IBaseLogic[T]) BaseCtrl[T,L] {
// 	// 自动兜底逻辑
// 	var finalHook = hook
// 	if hook == nil {
// 		finalHook = &DefaultHandlerHook[T]{}
// 	}
// 	return BaseCtrl[T,L]{
// 		Logic: logic, // 假设 BaseLogic 结构简单
// 		Hook:  finalHook,
// 	}
// }

