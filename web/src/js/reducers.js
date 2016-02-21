import { combineReducers } from 'redux';
import { routeReducer } from 'react-router-redux';
import { REQUEST_LOGIN_STATE, RECEIVE_LOGIN_STATE } from './actions';
import { START_LOGIN, LOGIN_FAILED, LOGIN_SUCCEEDED } from './actions';
import { RECEIVE_HOSTS_STATE, REQUEST_HOSTS_STATE, RECEIVE_HOST_STATE, REQUEST_HOST_STATE } from './actions';
import { RECEIVE_LOGS_FOR_HOST, REQUEST_LOGS_FOR_HOST } from './actions';


function authState(state = {isFetching: false, didInvalidate: true, userInfo: {loggedIn: false}}, action) {
  switch (action.type) {
    case RECEIVE_LOGIN_STATE:
      return Object.assign({}, state, {
        isFetching: false,
        didInvalidate: false,
        userInfo: action.userInfo,
      });
    case REQUEST_LOGIN_STATE:
      return Object.assign({}, state, {
        isFetching: true,
        didInvalidate: false,
      });
    default:
      return state;
  }
}

function loggingInState(state = {isFetching: false, credentials: {username: '', password: ''}}, action) {
  switch (action.type) {
    case START_LOGIN:
      return Object.assign({}, state, {
        isFetching: true,
        credentials: action.credentials,
      });
    case LOGIN_FAILED:
      return Object.assign({}, state, {
        isFetching: false,
        credentials: action.credentials,
      });
    case LOGIN_SUCCEEDED:
      return Object.assign({}, state, {
        isFetching: false,
        credentials: action.credentials,
      });
    default:
      return state;
  }
}

function hostsState(state = {isFetching: false, didInvalidate: false, hosts: {}}, action) {
  switch (action.type) {
    case RECEIVE_HOSTS_STATE:
      return Object.assign({}, state, {
        isFetching: false,
        didInvalidate: false,
        hosts: action.hosts
      });
    case REQUEST_HOSTS_STATE:
      return Object.assign({}, state, {
        isFetching: true,
        didInvalidate: false,
      });
    case RECEIVE_HOST_STATE:
    {
      let o = Object.assign({}, {
        isFetching: false,
        didInvalidate: false,
      }, action.host);

      let x = {};
      x[o.id] = o;

      let y = {hosts: Object.assign({}, state.hosts, x)};

      return Object.assign({}, state, y);
    }
    case REQUEST_HOST_STATE:
    {
      let o = {
        isFetching: true,
        didInvalidate: false
      };

      let x = {};
      x[action.hostID] = Object.assign({}, state.hosts[action.hostID], o);

      let y = {hosts: Object.assign({}, state.hosts, x)};

      return Object.assign({}, state, y);
    }
    default:
      return state;
  }
}

function logsState(state={}, action) {
  switch (action.type) {
  case REQUEST_LOGS_FOR_HOST:
  {
    let newO = Object.assign({}, state[action.hostID], {
      isFetching: true,
      didInvalidate: false,
    });
    let x = {};
    x[action.hostID] = newO;
    return Object.assign({}, state, x);
  }
  case RECEIVE_LOGS_FOR_HOST:
  {
    let newO = Object.assign({}, state[action.hostID], {
      isFetching: false,
      didInvalidate: false,
      logs: action.logs,
    });
    let x = {};
    x[action.hostID] = newO;
    return Object.assign({}, state, x);
  }
  default:
    return state;
  }
}

const phalanxApp = combineReducers({
  routing: routeReducer,
  authState,
  loggingInState,
  hostsState,
  logsState
});

export default phalanxApp;
