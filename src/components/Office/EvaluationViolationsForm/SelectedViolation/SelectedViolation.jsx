import React from 'react';
import { Button } from '@trussworks/react-uswds';

import styles from './SelectedViolation.module.scss';

const SelectedViolation = ({ violation, unselectViolation, isReadOnly }) => {
  if (!violation) {
    return null;
  }

  return (
    <div key={`${violation.id}-violation`} className={styles.violation}>
      <div className={styles.grow}>
        <h5>{`${violation.paragraphNumber} ${violation.title}`}</h5>
        <p>
          <small> {violation.requirementSummary}</small>
        </p>
      </div>
      {!isReadOnly && (
        <Button type="button" unstyled onClick={() => unselectViolation(violation.id)} role="button">
          Remove
        </Button>
      )}
    </div>
  );
};

export default SelectedViolation;
