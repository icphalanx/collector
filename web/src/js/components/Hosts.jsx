import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { fetchHostsStateIfNeeded } from '../actions';

class Hosts extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(fetchHostsStateIfNeeded());
  }

  render() {
    const { hosts, isFetching } = this.props;

    const actualHosts = [];
    for (let hostID in hosts) {
      if (!hosts.hasOwnProperty(hostID)) continue;
      actualHosts.push(hosts[hostID]);
    }
    console.log(this.props, hosts, actualHosts);

    return (
      <div className="row">
        <h1>Hosts</h1>
        <ul>
          {actualHosts.filter((host) => host.humanName).map((host) => (
            <li key={host.id}><Link to={`/hosts/${host.id}`}>{host.humanName}</Link></li>
          ))}
        </ul>
        {this.props.children}
      </div>
    );
  }
}

Hosts.propTypes = {
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

export default connect(mapStateToProps)(Hosts);
