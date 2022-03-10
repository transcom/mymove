import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Office/WeightDisplay/WeightDisplay.module.scss';
import { formatWeight } from 'utils/formatters';

const WeightDisplay = ({ heading, weightValue, onEdit, children }) => {
  return (
    <div className={classnames('maxw-tablet', styles.WeightDisplay)}>
      <div className={styles.heading}>
        <div>{heading}</div>
        {onEdit && (
          <Button unstyled type="button" className={styles.editButton} onClick={onEdit} data-testid="weightDisplayEdit">
            <FontAwesomeIcon icon="pen" title="edit" alt="" />
          </Button>
        )}
      </div>
      <div data-testid="weight-display" className={styles.value}>
        {Number.isFinite(weightValue) ? formatWeight(weightValue) : '—'}
      </div>
      {children && <div className={styles.details}>{children}</div>}
    </div>
  );
};

WeightDisplay.propTypes = {
  heading: PropTypes.string.isRequired,
  weightValue: PropTypes.number,
  children: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  onEdit: PropTypes.func,
};

WeightDisplay.defaultProps = {
  weightValue: null,
  children: null,
  onEdit: null,
};

export default WeightDisplay;
