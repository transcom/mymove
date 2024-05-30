import React from 'react';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './FinancialReviewButton.module.scss';

function FinancialReviewButton({ onClick, reviewRequested, isMoveLocked }) {
  return (
    <div>
      {reviewRequested ? (
        <div className={styles.EditFinancialReviewContainer}>
          <Tag className={styles.FinancialReviewTag}>Flagged for financial review</Tag>
          <Button
            type="Button"
            className={classnames(styles.EditFinancialReviewButton, ['usa-button--unstyled'])}
            onClick={onClick}
            disabled={isMoveLocked}
          >
            Edit
          </Button>
        </div>
      ) : (
        <Button
          type="Button"
          className={classnames(styles.FinancialReviewButton, ['usa-button--unstyled'])}
          onClick={onClick}
          disabled={isMoveLocked}
        >
          Flag move for financial review
        </Button>
      )}
    </div>
  );
}

FinancialReviewButton.propTypes = {
  onClick: PropTypes.func.isRequired,
  reviewRequested: PropTypes.bool,
};

FinancialReviewButton.defaultProps = {
  reviewRequested: false,
};

export default FinancialReviewButton;
