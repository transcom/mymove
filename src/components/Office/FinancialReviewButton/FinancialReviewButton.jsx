import React from 'react';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './FinancialReviewButton.module.scss';

// TODO: This button will switch states based on if the move has been flagged for financial reivew or not
// This will be covered in an up coming ticket!

function FinancialReviewButton({ onClick, reviewRequested }) {
  return (
    <div>
      {reviewRequested ? (
        <Tag className={classnames(styles.FinancialReviewTag, ['usa-tag'])}>Financial Review Requested</Tag>
      ) : (
        <Button
          type="Button"
          className={classnames(styles.FinancialReviewButton, ['usa-button usa-button--unstyled'])}
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
