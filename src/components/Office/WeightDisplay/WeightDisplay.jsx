import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Office/WeightDisplay/WeightDisplay.module.scss';
import { formatWeight } from 'shared/formatters';

const WeightDisplay = ({ heading, value, showEditBtn, onEdit }) => {
  return (
    <div className={classnames('maxw-tablet', styles.WeightDisplay)}>
      <div className={styles.heading}>
        <div>{heading}</div>
        {showEditBtn && (
          <Button unstyled type="button" className={styles.editButton} onClick={onEdit}>
            <FontAwesomeIcon icon="pen" title="edit" alt="" />
          </Button>
        )}
      </div>
      <div className={styles.value}>{value ? formatWeight(value) : null}</div>
    </div>
  );
};

WeightDisplay.propTypes = {
  heading: PropTypes.string.isRequired,
  value: PropTypes.number,
  showEditBtn: PropTypes.bool,
  onEdit: PropTypes.func,
};

WeightDisplay.defaultProps = {
  value: null,
  showEditBtn: false,
  onEdit: null,
};

export default WeightDisplay;
