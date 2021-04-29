import React from 'react';
import * as PropTypes from 'prop-types';

import styles from './DetailsPanel.module.scss';

const DetailsPanel = ({ title, editButton, children }) => {
  return (
    <div className={styles.DetailsPanel}>
      <div className="stackedtable-header">
        <h2>{title}</h2>
        {editButton && <div>{editButton}</div>}
      </div>
      {children}
    </div>
  );
};

DetailsPanel.propTypes = {
  title: PropTypes.string.isRequired,
  editButton: PropTypes.node,
  children: PropTypes.node.isRequired,
};

DetailsPanel.defaultProps = {
  editButton: null,
};

export default DetailsPanel;
