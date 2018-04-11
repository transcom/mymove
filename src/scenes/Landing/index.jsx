import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { bindActionCreators } from 'redux';

import { createMove } from 'scenes/Moves/ducks';
import Alert from 'shared/Alert';
import LoginButton from 'shared/User/LoginButton';
import { loadUserAndToken } from 'shared/User/ducks';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
    this.props.loadUserAndToken();
  }
  componentDidUpdate() {
    if (this.props.hasSubmitSuccess)
      this.props.push(`moves/${this.props.currentMove.id}`);
  }
  startMove = values => {
    this.props.createMove({});
  };
  render() {
    return (
      <div className="usa-grid">
        <h1>Welcome! </h1>
        <div>
          {this.props.hasSubmitError && (
            <Alert type="error" heading="An error occurred">
              There was an error starting your move.
            </Alert>
          )}
          {!this.props.isLoggedIn && <LoginButton />}
          <div />
        </div>
        {this.props.isLoggedIn && (
          <button onClick={this.startMove}>Start a move</button>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
  currentMove: state.submittedMoves.currentMove,
  hasSubmitError: state.submittedMoves.hasSubmitError,
  hasSubmitSuccess: state.submittedMoves.hasSubmitSuccess,
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, loadUserAndToken, createMove }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
