import React, { Component } from 'react';

import './index.css';

class TextBoxWithEditLink extends Component {
  constructor() {
    super();
    this.state = {
      isTextEditable: false,
      value: '',
    };
  }

  handleClick = e => {
    e.preventDefault();
    this.setState({
      isTextEditable: !this.state.isTextEditable,
    });
  };

  handleChange = e => {
    this.setState({
      value: e.target.value,
    });
  };

  handleSubmit = e => {
    e.preventDefault();
    this.setState({
      value: e.target.value,
    });
  };

  render() {
    if (this.state.isTextEditable) {
      // editable, something submitted or not
      return (
        <form onSubmit={this.onSubmit}>
          <textarea value={this.state.value} onChange={this.handleChange} />
          <a href="/save" onClick={this.handleClick}>
            Save
          </a>
        </form>
      );
    }
    if (!this.state.isTextEditable && !this.state.value) {
      // not editable, nothing submitted
      return (
        <div>
          <p>This is where informative service member move text will go.</p>
          <a href="/edit" onClick={this.handleClick}>
            Edit
          </a>
        </div>
      );
    }
    if (!this.state.isTextEditable && this.state.value) {
      // not editable, something submitted
      return (
        <div>
          <p>{this.state.value}</p>
          <a href="/edit" onClick={this.handleClick}>
            Edit
          </a>
        </div>
      );
    } // else...
  }
}

export default TextBoxWithEditLink;
