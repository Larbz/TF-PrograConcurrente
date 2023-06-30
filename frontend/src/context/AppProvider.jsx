import { Dispatch, createContext,useReducer,useContext  } from "react";
import { AppReducer, InitialApp } from "./AppReducer";



export const AppContext = createContext();

const AppProvider = ({children})=>{
    const [store,dispatch]=useReducer(AppReducer,InitialApp);
    return (
        <AppContext.Provider value={{ store, dispatch }}>
            {children}
        </AppContext.Provider>
    );
}

export const useApp = () => useContext(AppContext).store;
export const useDispatch = () => useContext(AppContext).dispatch;

export default AppProvider;