import React from 'react';
// import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './FinancialReviewButton.module.scss';

function FinancialReviewButton({ onClick }) {
  return (
    <div>
      <Button
        type="Button"
        className={classnames(styles.FinancialReviewButton, ['usa-button usa-button--unstyled'])}
        onClick={onClick}
      >
        Flag move for financial review
      </Button>
    </div>
  );
}

FinancialReviewButton.propTypes = {
  onClick: PropTypes.func.isRequired,
};

export default FinancialReviewButton;
