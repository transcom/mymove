import React from 'react';
import * as PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './DetailsPanel.module.scss';

const DetailsPanel = ({ title, tag, editButton, children, shipmentsInfoNonPpm, className }) => {
  return (
    <div className={classnames(styles.DetailsPanel, className)}>
      <div className="stackedtable-header">
        <h2>
          {title}
          {tag && (
            <Tag className={styles.tag} data-testid="detailsPanelTag">
              {tag}
            </Tag>
          )}
        </h2>
        {/* {editButton && <div>{editButton}</div>} */}
        {shipmentsInfoNonPpm.length === 0 ? null : editButton && <div>{editButton}</div>}
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
  tag: PropTypes.string,
  shipmentsInfoNonPpm: PropTypes.array,
};

DetailsPanel.defaultProps = {
  editButton: null,
  className: '',
  tag: '',
  shipmentsInfoNonPpm: [],
};

export default DetailsPanel;
