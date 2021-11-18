import React from 'react';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './FinancialReviewButton.module.scss';

function FinancialReviewButton({ onClick, reviewRequested }) {
  return (
    <div className={styles.FinancialReview}>
      {reviewRequested ? (
        <>
          <Tag className={styles.financialReviewTag}>Flagged for financial review</Tag>
          <Button
            type="button"
            className={classnames(styles.financialReviewEdit, 'usa-button--unstyled')}
            onClick={onClick}
          >
            Edit
          </Button>
        </>
      ) : (
        <Button
          type="Button"
          className={classnames(styles.financialReviewButton, ['usa-button--unstyled'])}
          onClick={onClick}
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
