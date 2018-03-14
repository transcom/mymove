import React, { Component } from 'react';
import { connect } from 'react-redux';

import LoginButton from 'shared/User/LoginButton';
import { bindActionCreators } from 'redux';
import { loadUserAndToken } from 'shared/User/ducks';
import { push } from 'react-router-redux';
import { createMove } from 'scenes/Moves/ducks';
export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
    this.props.loadUserAndToken();
  }
  startMove = values => {
    this.props.createMove({ selected_move_type: 'COMBO' });
  };
  render() {
    return (
      <div className="usa-grid">
        <h1>Welcome!</h1>
        <div>
          <LoginButton />
        </div>
        {this.props.isLoggedIn && (
          <button onClick={this.startMove}>Start a move</button>
        )}
        {this.props.hasSubmitSuccess &&
          this.props.push(getLegaleseRoute(this.prop))}
      </div>
    );
  }
}

function getLegaleseRoute(props) {
  //this is a horrible hack until we can get move id from server
  var moveID = props.currentMove;
  return `moves/${moveID}/legalese`;
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
