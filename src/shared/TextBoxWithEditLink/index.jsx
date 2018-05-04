import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

class TextBoxWithEditLink extends Component {
  constructor() {
    super();
    this.state = {
      isTextEditable: false,
    };
  }

  handleClick = e => {
    e.preventDefault();
    this.setState({
      isTextEditable: true,
    });
    console.log('isTextEditable: ', this.state.isTextEditable);
  };

  render() {
    if (this.state.isTextEditable === false) {
      return (
        <form>
          <p>This is where informative service member move text will go.</p>
          <a href="#" onClick={this.handleClick}>
            Edit
          </a>
        </form>
      );
    } else {
      return (
        <form>
          <textarea>
            This is where informative service member move text will go.
          </textarea>
          <a href="#" onClick={this.handleClick}>
            Save
          </a>
        </form>
      );
    }
  }
}

export default TextBoxWithEditLink;
