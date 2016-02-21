import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { fetchLoginStateIfNeeded } from '../actions';
import LoginPage from './LoginPage';
import Dashboard from './Dashboard';
import { performLogout } from '../actions';
import { Link } from 'react-router';


class Root extends Component {
  constructor(props) {
    super(props);
    this.handleLogout = this.handleLogout.bind(this);
  }

  handleLogout() {
    const { dispatch } = this.props;
    dispatch(performLogout());
  }

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch(fetchLoginStateIfNeeded());
  }

  render() {
    const { userInfo, isFetching } = this.props;

    if (!userInfo.loggedIn) {
      return (
        <div className={isFetching ? "fetching-wrap fetching" : "fetching-wrap"}>
          <LoginPage />
        </div>
      );
    }

    return (
      <div className={isFetching ? "fetching-wrap fetching" : "fetching-wrap"}>
        <div className="container">
          <div className="header clearfix">
            <a className="text-muted pull-xs-right" onClick={this.handleLogout}>Log out</a>
            <h3 className="text-muted"><Link to="/">Phalanx</Link></h3>
          </div>
          <div>{this.props.children ? this.props.children : (<Dashboard />)}</div>
        </div>
      </div>
    );
  }
}

Root.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  userInfo: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const { authState } = state;
  const {
    isFetching,
    userInfo,
  } = authState;
  return {
    isFetching,
    userInfo,
  };
};

export default connect(mapStateToProps)(Root);
