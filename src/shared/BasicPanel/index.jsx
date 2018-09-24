import React, { Component } from 'react';
import './index.css';

class BasicPanel extends Component {
  render() {
    const { title, children } = this.props;
    return (
      <div className="basic-panel">
        <div className="basic-panel-title">{title}</div>
        <div className="basic-panel-content">{children}</div>
      </div>
    );
  }
}

export default BasicPanel;
