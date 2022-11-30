import React from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './SelectedViolation.module.scss';

import { PWSViolationShape } from 'types';

const SelectedViolation = ({ violation, unselectViolation }) => {
  if (!violation) {
    return null;
  }

  return (
    <div key={`${violation.id}-violation`} className={styles.violation}>
      <div className={styles.grow}>
        <h5>{`${violation.paragraphNumber} ${violation.title}`}</h5>
        <p>
          <small>{violation.requirementSummary}</small>
        </p>
      </div>
      <Button type="button" unstyled onClick={() => unselectViolation(violation.id)} role="button">
        Remove
      </Button>
    </div>
  );
};

SelectedViolation.propTypes = {
  violation: PWSViolationShape,
  unselectViolation: PropTypes.func.isRequired,
};

SelectedViolation.defaultProps = {
  violation: null,
};

export default SelectedViolation;
