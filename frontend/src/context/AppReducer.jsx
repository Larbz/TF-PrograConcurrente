export const InitialApp = {
    tableInfo: [],
};

export const types = {
    addTableInfo: "addTableInfo",
};

export const AppReducer = (state, action) => {
    switch (action.type) {
        case types.addTableInfo:
            return {
                ...state,
                tableInfo: action.payload,
            };
    }
};
