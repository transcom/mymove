import React, { Component } from 'react';
import { connect } from 'react-redux';

import LoginButton from 'shared/User/LoginButton';
import { bindActionCreators } from 'redux';
import { loadUserAndToken } from 'shared/User/ducks';
import { push } from 'react-router-redux';
export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
    this.props.loadUserAndToken();
  }
  gotoLegalese = values => {
    this.props.push(getLegaleseRoute());
  };
  render() {
    return (
      <div className="usa-grid">
        <h1>Welcome!</h1>
        <div>
          <LoginButton />
        </div>
        {this.props.isLoggedIn && (
          <button onClick={this.gotoLegalese}>Start a move</button>
        )}
      </div>
    );
  }
}

function getLegaleseRoute() {
  //this is a horrible hack until we can get move id from server
  const moveId = process.env.REACT_APP_MOVE_UUID;
  return `moves/${moveId}/legalese`;
}

const mapStateToProps = state => ({
  isLoggedIn: state.user.isLoggedIn,
});
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, loadUserAndToken }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Landing);
