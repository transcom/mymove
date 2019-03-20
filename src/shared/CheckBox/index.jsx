import React, { Component } from 'react';
import PropTypes from 'prop-types';

class CheckBox extends Component {
  handleOnChange = e => {
    this.props.onChangeHandler(e.target.checked);
  };

  render() {
    return (
      <>
        <input id="agree-checkbox" type="checkbox" checked={this.props.checked} onChange={this.handleOnChange} />
        <label htmlFor="agree-checkbox"> {this.props.children}</label>
      </>
    );
  }
}

CheckBox.propTypes = {
  checked: PropTypes.bool,
  onChangeHandler: PropTypes.func,
};

export default CheckBox;
