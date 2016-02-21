import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { performLogin } from '../actions';

class LoginPage extends Component {
  constructor(props) {
    super(props);
    this.state = {username: '', password: ''};
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
  }

  handleSubmit(e) {
    e.preventDefault();

    const { username, password } = this.state;
    const { dispatch } = this.props;
    dispatch(performLogin(username, password));
  }

  handleUsernameChange(e) {
    this.setState({username: e.target.value});
  }

  handlePasswordChange(e) {
    this.setState({password: e.target.value});
  }

  render() {
    const { credentials, isFetching } = this.props;

    return (
      <div className="container">
        <div className="header clearfix">
          <h3 className="text-muted">Phalanx</h3>
        </div>

        <div className="jumbotron">
          <h1 className="display-3">Log in</h1>
          <p className="lead">A valid user account is required to access Phalanx.</p>
        </div>

        <div className="row">
          <form onSubmit={this.handleSubmit}>
            <fieldset className="form-group">
              <label htmlFor="login-username">Username:</label>
              <input disabled={isFetching} id="login-username" name="username" value={this.state.username} type="text" className="form-control" onChange={this.handleUsernameChange} />
            </fieldset>
            <fieldset className="form-group">
              <label htmlFor="login-password">Password:</label>
              <input disabled={isFetching} id="login-password" name="password" value={this.state.password} type="password" className="form-control" onChange={this.handlePasswordChange} />
            </fieldset>
            <button type="submit" className="btn btn-primary">Log in</button>
          </form>
        </div>
      </div>
    );
  }
}

LoginPage.propTypes = {
  dispatch: PropTypes.func.isRequired,
  isFetching: PropTypes.bool.isRequired,
  credentials: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const { loggingInState } = state;
  const {
    isFetching,
    credentials,
  } = loggingInState;
  return {
    isFetching,
    credentials,
  };
};

export default connect(mapStateToProps)(LoginPage);
