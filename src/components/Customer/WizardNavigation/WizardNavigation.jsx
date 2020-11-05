import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

const WizardNavigation = ({
  isFirstPage,
  isLastPage,
  disableNext,
  showFinishLater,
  onBackClick,
  onNextClick,
  onCancelClick,
}) => {
  return (
    <div>
      {!isFirstPage && (
        <Button type="button" secondary onClick={onBackClick} data-testid="wizardBackButton">
          Back
        </Button>
      )}
      <Button
        type="button"
        onClick={onNextClick}
        data-testid={isLastPage ? 'wizardCompleteButton' : 'wizardNextButton'}
        disabled={disableNext}
      >
        {isLastPage ? 'Complete' : 'Next'}
      </Button>

      {showFinishLater && (
        <Button type="button" unstyled onClick={onCancelClick} data-testid="wizardFinishLaterButton">
          Finish later
        </Button>
      )}
    </div>
  );
};

WizardNavigation.propTypes = {
  isFirstPage: PropTypes.bool,
  isLastPage: PropTypes.bool,
  disableNext: PropTypes.bool,
  showFinishLater: PropTypes.bool,
  onBackClick: PropTypes.func,
  onNextClick: PropTypes.func,
  onCancelClick: PropTypes.func,
};

WizardNavigation.defaultProps = {
  isFirstPage: false,
  isLastPage: false,
  disableNext: false,
  showFinishLater: false,
  onBackClick: () => {},
  onNextClick: () => {},
  onCancelClick: () => {},
};

export default WizardNavigation;
