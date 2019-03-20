import { Component } from 'react';
import React from 'react';
import PropTypes from 'prop-types';

class PopUp extends Component {
  handleClick = e => {
    // Prevent this from checking the box after opening the alert.
    e.preventDefault();
    alert(this.props.alertMessage);
  };

  render() {
    return <a onClick={this.handleClick}>{this.props.children}</a>;
  }
}

PopUp.propTypes = {
  alertMessage: PropTypes.string.isRequired,
};

export default PopUp;
