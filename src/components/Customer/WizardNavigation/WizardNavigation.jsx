import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './WizardNavigation.module.scss';

const WizardNavigation = ({
  isReviewPage,
  isFirstPage,
  isLastPage,
  disableNext,
  showFinishLater,
  editMode,
  readOnly,
  onBackClick,
  onAddShipment,
  onNextClick,
  onCancelClick,
}) => {
  if (readOnly) {
    return (
      <div className={styles.WizardNavigation}>
        <Button type="button" className={styles.Button} onClick={onCancelClick} data-testid="wizardCancelButton">
          Return home
        </Button>
      </div>
    );
  }

  let submitButtonText = 'Next';
  if (isLastPage) submitButtonText = 'Complete';
  else if (editMode) submitButtonText = 'Save';

  let cancelButtonText = 'Finish Later';
  if (editMode) cancelButtonText = 'Cancel';

  return (
    <div className={styles.WizardNavigation}>
      {!isFirstPage && !editMode && (
        <Button type="button" className={styles.button} secondary onClick={onBackClick} data-testid="wizardBackButton">
          Back
        </Button>
      )}

      {(showFinishLater || editMode) && (
        <Button
          type="button"
          secondary
          className={styles.button}
          onClick={onCancelClick}
          data-testid="wizardCancelButton"
        >
          {cancelButtonText}
        </Button>
      )}

      {isReviewPage && (
        <Button type="button" onClick={onAddShipment} className={styles.button} data-testid="wizardAddShipmentButton">
          <FontAwesomeIcon icon="plus" className={styles.addShipmentIcon} />
          &nbsp;&nbsp;Add Shipment
        </Button>
      )}

      <Button
        type="button"
        onClick={onNextClick}
        className={styles.button}
        data-testid={isLastPage ? 'wizardCompleteButton' : 'wizardNextButton'}
        disabled={disableNext}
      >
        {submitButtonText}
      </Button>
    </div>
  );
};

WizardNavigation.propTypes = {
  isReviewPage: PropTypes.bool,
  isFirstPage: PropTypes.bool,
  isLastPage: PropTypes.bool,
  disableNext: PropTypes.bool,
  showFinishLater: PropTypes.bool,
  editMode: PropTypes.bool,
  readOnly: PropTypes.bool,
  onBackClick: PropTypes.func,
  onAddShipment: PropTypes.func,
  onNextClick: PropTypes.func,
  onCancelClick: PropTypes.func,
};

WizardNavigation.defaultProps = {
  isReviewPage: false,
  isFirstPage: false,
  isLastPage: false,
  disableNext: false,
  showFinishLater: false,
  editMode: false,
  readOnly: false,
  onBackClick: () => {},
  onAddShipment: () => {},
  onNextClick: () => {},
  onCancelClick: () => {},
};

export default WizardNavigation;
