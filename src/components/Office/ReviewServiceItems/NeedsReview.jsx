import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './NeedsReview.module.scss';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows if any service items have not been reviewed yet.
 * */
const NeedsReview = ({ numberOfItems, onClick }) => {
  return (
    <div data-testid="NeedsReview" className={styles.NeedsReview}>
      <strong data-testid="header">
        {numberOfItems} item{numberOfItems > 1 ? 's' : ''} still needs your review
      </strong>
      <p data-testid="content" className={styles.content}>
        Accept or reject all service items, then authorized payment.
      </p>
      <Button data-testid="finishReviewBtn" type="button" secondary onClick={onClick}>
        Finish review
      </Button>
    </div>
  );
};

NeedsReview.propTypes = {
  numberOfItems: PropTypes.number.isRequired,
  onClick: PropTypes.func,
};

NeedsReview.defaultProps = {
  onClick: null,
};

export default NeedsReview;
