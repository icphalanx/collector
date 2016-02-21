import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { browserHistory, Router, Route } from 'react-router';
import configureStore from './configureStore';
import Root from "./components/Root";
import Hosts from "./components/Hosts";
import Host from "./components/Host";

let store = configureStore();

render(
  <Provider store={store}>
    <Router history={browserHistory}>
      <Route path="/" component={Root}>
        <Route path="hosts" component={Hosts}>
	        <Route path=":hostID" component={Host} />
        </Route>
      </Route>
    </Router>
  </Provider>,
  document.getElementById('root')
);


