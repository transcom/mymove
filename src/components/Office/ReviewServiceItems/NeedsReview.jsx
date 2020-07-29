import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './NeedsReview.module.scss';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows if any service items have not been reviewed yet.
 * */
const NeedsReview = ({ numberOfItems, handleFinishReviewBtn }) => {
  return (
    <div data-testid="NeedsReview" className={styles.NeedsReview}>
      <div className={styles.header}>{`${numberOfItems} item still needs your review`}</div>
      <div>Accept or reject all service items, then authorized payment.</div>
      <Button type="button" secondary onClick={handleFinishReviewBtn}>
        Finish review
      </Button>
    </div>
  );
};

NeedsReview.propTypes = {
  numberOfItems: PropTypes.number.isRequired,
  handleFinishReviewBtn: PropTypes.func,
};

NeedsReview.defaultProps = {
  handleFinishReviewBtn: null,
};

export default NeedsReview;
