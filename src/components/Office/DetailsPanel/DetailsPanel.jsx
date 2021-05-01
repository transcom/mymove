import React from 'react';
import * as PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './DetailsPanel.module.scss';

const DetailsPanel = ({ title, editButton, children, className }) => {
  return (
    <div className={classnames(styles.DetailsPanel, className)}>
      <div className="stackedtable-header">
        <h2>{title}</h2>
        {editButton && <div>{editButton}</div>}
      </div>
      {children}
    </div>
  );
};

DetailsPanel.propTypes = {
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  editButton: PropTypes.node,
  title: PropTypes.string.isRequired,
};

DetailsPanel.defaultProps = {
  editButton: null,
  className: '',
};

export default DetailsPanel;
