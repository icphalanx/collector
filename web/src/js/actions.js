export const RECEIVE_LOGIN_STATE = 'RECEIVE_LOGIN_STATE';

export function receiveLoginState(userInfo) {
  return {
    type: RECEIVE_LOGIN_STATE,
    userInfo
  };
};

export const REQUEST_LOGIN_STATE = 'REQUEST_LOGIN_STATE';

export function requestLoginState() {
  return {
    type: REQUEST_LOGIN_STATE,
  };
};

export function fetchLoginState() {
  return dispatch => {
    dispatch(requestLoginState());
    return fetch('/api/user/me', {credentials: 'same-origin'})
             .then(req => req.json())
             .then(json => dispatch(receiveLoginState(json)));
  };
};

export function shouldFetchLoginState(state) {
  const { authState } = state;

  if (authState.isFetching) {
    return false;
  }

  return authState.didInvalidate;
};

export function fetchLoginStateIfNeeded() {
  return ((dispatch, getState) => {
    if (shouldFetchLoginState(getState())) {
      return dispatch(fetchLoginState());
    }
  });
};

export const RECEIVE_HOSTS_STATE = 'RECEIVE_HOSTS_STATE';
export const REQUEST_HOSTS_STATE = 'REQUEST_HOSTS_STATE';
export const RECEIVE_HOST_STATE = 'RECEIVE_HOST_STATE';
export const REQUEST_HOST_STATE = 'REQUEST_HOST_STATE';

export function receiveHostsState(hosts) {
  return {
    type: RECEIVE_HOSTS_STATE,
    hosts: hosts,
  };
};

export function requestHostsState() {
  return {
    type: REQUEST_HOSTS_STATE
  };
};

export function receiveHostState(host) {
  return {
    type: RECEIVE_HOST_STATE,
    host: host,
  };
};

export function requestHostState(hostID) {
  return {
    type: REQUEST_HOST_STATE,
    hostID: hostID,
  };
};

export function shouldFetchHostsState(state) {
  const { hostsState } = state;
  if (!hostsState.length) {
    return true;
  }

  if (hostsState.didInvalidate) {
    return true;
  }

  return false;
};

export function fetchHostsStateIfNeeded() {
  return ((dispatch, getState) => {
    if (shouldFetchHostsState(getState())) {
      return dispatch(fetchHostsState());
    }
  });
};

export function shouldFetchHostState(state, hostID) {
  if (shouldFetchHostsState(state)) {
    return true;
  }

  const { hostsState } = state;

  if (!hostsState[hostID]) {
    return true;
  }

  if (hostsState[hostID].didInvalidate) {
    return true;
  }

  return false;
};

export function fetchHostStateIfNeeded(hostID) {
  return ((dispatch, getState) => {
    if (shouldFetchHostState(getState(), hostID)) {
      return dispatch(fetchHostState(hostID));
    }
  });
};

export function fetchHostsStateIfNeeded() {
  return ((dispatch, getState) => {
    if (shouldFetchHostsState(getState())) {
      return dispatch(fetchHostsState());
    }
  });
};

export function fetchHostState(hostID) {
  return (dispatch) => {
    dispatch(requestHostState(hostID));
    fetch(`/api/hosts/${hostID}`, {credentials: 'same-origin'})
      .then(req => req.json())
      .then(json => dispatch(receiveHostState(json)));
  };
};

export function fetchHostsState() {
  return (dispatch) => {
    dispatch(requestHostsState());
    fetch(`/api/hosts`, {credentials: 'same-origin'})
      .then(req => req.json())
      .then(json => dispatch(receiveHostsState(json)));
  };
};

export const START_LOGIN = 'START_LOGIN';
export const LOGIN_FAILED = 'LOGIN_FAILED';
export const LOGIN_SUCCEEDED = 'LOGIN_SUCCEEDED';

export function startLogin(username, password) {
  return {
    type: START_LOGIN,
    credentials: {
      username: username,
      password: password,
    }
  };
};

export function loginFailed(username, password) {
  return {
    type: LOGIN_FAILED,
    credentials: {
      username: username,
      password: password,
    }
  };
};

export function loginSucceeded(username, password) {
  return {
    type: LOGIN_SUCCEEDED,
    credentials: {
      username: username,
      password: password,
    }
  };
};

export function performLogin(username, password) {
  return dispatch => {
    dispatch(startLogin(username, password));
    let headers = new Headers();
    headers.set('Content-type', 'application/x-www-form-urlencoded');
    let fs = `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`;
    return fetch('/api/user/login', {credentials: 'same-origin', method: 'post', body: fs, headers: headers})
      .then(req => {
        dispatch(loginSucceeded(username, password));
        return fetchLoginState()(dispatch)
      }, req => {
        dispatch(loginFailed(username, password));
      });
  };
};

export function performLogout() {
  return dispatch => {
    return fetch('/api/user/logout', {credentials: 'same-origin', method: 'post'})
      .then(req => {
        return fetchLoginState()(dispatch)
      });
  };
};

export const REQUEST_LOGS_FOR_HOST = 'REQUEST_LOGS_FOR_HOST';
export const RECEIVE_LOGS_FOR_HOST = 'RECEIVE_LOGS_FOR_HOST';

export function requestLogsForHost(hostID) {
  return {
    type: REQUEST_LOGS_FOR_HOST,
    hostID: hostID,
  }
}

export function receiveLogsForHost(hostID, logs) {
  return {
    type: RECEIVE_LOGS_FOR_HOST,
    hostID: hostID,
    logs: logs,
  };
}

export function fetchLogsForHost(hostID) {
  return ((dispatch) => {
    dispatch(requestLogsForHost(hostID));
    fetch(`/api/hosts/${hostID}/logs`, {credentials: 'same-origin'})
      .then((resp) => resp.json())
      .then((json) => dispatch(receiveLogsForHost(hostID, json)))
  });
}


export function fetchLogsForHostIfNeeded(hostID) {
  return ((dispatch, getState) => {
    const s = getState();
    if (!s[hostID] || (!s[hostID].isFetching && !s[hostID].logs)) {
      return dispatch(fetchLogsForHost(hostID));
    }
  });
};
