import React from 'react';
import * as PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './DetailsTable.module.scss';

const DetailsTable = ({ title, editable, editTitle, editTestLabel, editLinkLocation, children }) => {
  return (
    <div className={styles.DetailsTable}>
      <div className="stackedtable-header">
        <h2>{title}</h2>
        {editable && (
          <div>
            <Link className="usa-button usa-button--secondary" data-testid={editTestLabel} to={editLinkLocation}>
              {editTitle}
            </Link>
          </div>
        )}
      </div>
      {children}
    </div>
  );
};

DetailsTable.propTypes = {
  title: PropTypes.string.isRequired,
  editTitle: PropTypes.string,
  editable: PropTypes.bool,
  editTestLabel: PropTypes.string,
  editLinkLocation: PropTypes.string,
  children: PropTypes.node.isRequired,
};

DetailsTable.defaultProps = {
  editable: false,
  editTitle: 'Edit',
  editTestLabel: 'edit',
  editLinkLocation: '#',
};

export default DetailsTable;
