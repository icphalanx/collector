import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import rootReducer from './reducers';
import { syncHistory } from 'react-router-redux';
import { browserHistory } from 'react-router';

const reduxRouterMiddleware = syncHistory(browserHistory);

export default function configureStore(initialState) {
  const store = createStore(
    rootReducer,
    initialState,
    applyMiddleware(
      thunkMiddleware,
      reduxRouterMiddleware
    )
  )

  reduxRouterMiddleware.listenForReplays(store);

  return store;
};
