import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './WizardNavigation.module.scss';

import { ButtonUsa as Button } from 'shared/standardUI/Buttons/ButtonUsa';

export const wizardActionButtonStyle = styles['wizard-action-button'];
export const wizardMainButtonStyle = styles['wizard-main-button'];

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
        <Button
          type="button"
          className={wizardActionButtonStyle}
          onClick={onCancelClick}
          data-testid="wizardCancelButton"
        >
          Return home
        </Button>
      </div>
    );
  }

  let submitButtonText = 'Next';
  if (isLastPage) submitButtonText = 'Complete';
  else if (editMode) submitButtonText = 'Save';

  let cancelButtonText = 'Finish later';
  if (editMode) cancelButtonText = 'Cancel';

  return (
    <div className={styles.WizardNavigation}>
      {!isFirstPage && !editMode && (
        <Button
          type="button"
          className={wizardActionButtonStyle}
          secondary
          onClick={onBackClick}
          data-testid="wizardBackButton"
        >
          Back
        </Button>
      )}

      {(showFinishLater || editMode) && (
        <Button
          type="button"
          secondary
          className={wizardActionButtonStyle}
          onClick={onCancelClick}
          data-testid="wizardCancelButton"
        >
          {cancelButtonText}
        </Button>
      )}

      {isReviewPage && (
        <Button
          type="button"
          onClick={onAddShipment}
          className={wizardMainButtonStyle}
          data-testid="wizardAddShipmentButton"
        >
          <FontAwesomeIcon icon="plus" className={styles.addShipmentIcon} />
          <span>Add shipment</span>
        </Button>
      )}

      <Button
        type="button"
        onClick={onNextClick}
        className={wizardMainButtonStyle}
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
