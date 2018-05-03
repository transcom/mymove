import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

class TextBoxWithEditLink extends Component {
  handleClick = e => {
    e.preventDefault();
    console.log('It worked!');
  };

  render() {
    return (
      <form>
        <textarea />
        <a href="#" onClick={this.handleClick}>
          Edit
        </a>
      </form>
    );
  }
}

export default TextBoxWithEditLink;
