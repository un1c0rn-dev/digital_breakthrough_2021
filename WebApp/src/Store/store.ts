import thunk from 'redux-thunk';
import { applyMiddleware, compose, createStore, Store } from 'redux';

import { reducer as rootReducer } from './reducer';
import { IMainStore } from './reducer';

const composeEnhancer = (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

export default function configureStore(initialState: IMainStore): Store {
  return createStore(rootReducer, initialState, composeEnhancer(applyMiddleware(thunk)));
}
