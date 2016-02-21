import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { fetchHostStateIfNeeded } from '../actions';
import HostLogs from './HostLogs';

class Host extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { dispatch, params } = this.props;
    dispatch(fetchHostStateIfNeeded(params.hostID));
  }

  render() {
    const { hosts, params, isFetching } = this.props;

    const host = hosts[params.hostID];

    if (isFetching || !host) {
      return (<div></div>);
    }

    return (
      <div className="row">
        <h2>{host.humanName}</h2>
        <h3>Metrics</h3>
        <dl>
          <dt>Total packages installed</dt>
            <dd>1566</dd>
          <dt>Packages needing updates</dt>
            <dd>46</dd>
        </dl>
        <HostLogs host={host} />
        {this.props.children}
      </div>
    );
  }
}

Host.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  hosts: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const { hostsState } = state;
  const {
    isFetching,
    hosts,
  } = hostsState;
  return {
    isFetching,
    hosts,
  };
};

export default connect(mapStateToProps)(Host);
