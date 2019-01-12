import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import './index.css';

class BasicPanel extends Component {
  render() {
    const { title, children, className } = this.props;
    return (
      <div className="basic-panel">
        <div className="basic-panel-title">{title}</div>
        <div className={classnames('basic-panel-content', className)}>{children}</div>
      </div>
    );
  }
}

BasicPanel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node,
  className: PropTypes.string,
};

export default BasicPanel;
