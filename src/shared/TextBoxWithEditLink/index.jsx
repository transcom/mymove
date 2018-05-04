import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

class TextBoxWithEditLink extends Component {
  constructor() {
    super();
    this.state = {
      formState: 'textarea',
      isToggleOn: false,
    };
  }

  handleClick = e => {
    e.preventDefault();
    this.setState({
      formState: 'form',
      isToggleOn: true,
    });
    console.log('formState: ', this.state.formState);
  };

  render() {
    return (
      <form>
        <p>This is where informative service member move text will go.</p>
        <a href="#" onClick={this.handleClick}>
          {this.state.isToggleOn ? 'Edit' : 'Save'}
        </a>
      </form>
    );
  }
}

export default TextBoxWithEditLink;
