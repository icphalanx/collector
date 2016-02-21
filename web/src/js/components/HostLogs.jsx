import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { fetchLogsForHostIfNeeded } from '../actions';

class HostLogs extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { dispatch, host } = this.props;
    dispatch(fetchLogsForHostIfNeeded(host.id));
  }

  render() {
    const { logs, host, isFetching } = this.props;

    const hostLogs = logs[host.id];

    if (!hostLogs || hostLogs.isFetching) {
      return (<div></div>);
    }

    return (
      <div className="row">
        <h3>Logs</h3>
        <ul className="logs-list">
          {hostLogs.logs.map((logLine) => (
            <li>[{logLine.Timestamp}] {logLine.LogLine}</li>
          ))}
        </ul>
      </div>
    );
  }
}

HostLogs.propTypes = {
  dispatch: PropTypes.func.isRequired,
  logs: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const { logsState } = state;
  return {
    logs: logsState,
  };
};

export default connect(mapStateToProps)(HostLogs);
