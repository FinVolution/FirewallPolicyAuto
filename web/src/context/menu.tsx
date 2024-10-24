import { BreadcrumbItemType } from 'antd/es/breadcrumb/Breadcrumb';
import { createContext, useContext, useReducer } from 'react'

export const MenuContext = createContext<any>(null)
export const MenuDispatchContext = createContext<any>(null)

interface MenuType {
  selectedKeys: string
  items: BreadcrumbItemType[]
}


interface ActionType extends MenuType {
  type: 'changed' | string
}

function menuReducer(menu: MenuType, action: ActionType) {
  switch (action.type) {
    case 'changed': {
      return {
        ...menu,
        items: action.items,
        selectedKeys: action.selectedKeys
      };
    }
    default: {
      throw Error('Unknown action: ' + action.type);
    }
  }
}


const initialMenu: MenuType = {
  selectedKeys: '',
  items: []
}
export function TasksProvider({ children }: any) {
  const [tasks, dispatch] = useReducer(menuReducer, initialMenu);

  return (
    <MenuContext.Provider value={tasks}>
      <MenuDispatchContext.Provider value={dispatch}>
        {children}
      </MenuDispatchContext.Provider>
    </MenuContext.Provider>
  );
}


export function useMenus() {
  return useContext<MenuType>(MenuContext);
}

export function useMenusDispatch() {
  return useContext(MenuDispatchContext);
}