/* eslint-disable react/prop-types */
import React from 'react';
import { bool, element, string, oneOfType, number, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import { ReactComponent as AcceptIcon } from 'shared/icon/accept-inversed.svg';

const NumberCircle = ({ num }) => <div className={styles['number-circle']}>{num}</div>;

NumberCircle.propTypes = {
  num: string.isRequired,
};

const Step = ({
  actionBtnDisabled,
  actionBtnLabel,
  children,
  complete,
  completedHeaderText,
  containerClassName,
  description,
  editBtnDisabled,
  editBtnLabel,
  headerText,
  onActionBtnClick,
  onEditBtnClick,
  secondary,
  step,
}) => {
  const secondaryClassName = styles['usa-button--secondary'];
  const disabledClassName = styles['btn--disabled'];
  return (
    <div className={`${containerClassName} margin-bottom-6`}>
      <div className={`${styles['step-header-container']} margin-bottom-2`}>
        {complete ? <AcceptIcon aria-hidden className={styles.accept} /> : <NumberCircle num={step} />}
        <strong>{complete ? completedHeaderText : headerText}</strong>
        {editBtnLabel && (
          <Button
            disabled={editBtnDisabled}
            className={`${styles['edit-button']} ${editBtnDisabled ? disabledClassName : ''}`}
            onClick={onEditBtnClick}
          >
            {editBtnLabel}
          </Button>
        )}
      </div>

      {children || <p>{description}</p>}
      {actionBtnLabel && (
        <Button
          className={`margin-top-3 ${styles['action-btn']} ${secondary ? secondaryClassName : ''} ${
            actionBtnDisabled ? disabledClassName : ''
          }`}
          disabled={actionBtnDisabled}
          onClick={onActionBtnClick}
        >
          {actionBtnLabel}
        </Button>
      )}
    </div>
  );
};

Step.propTypes = {
  actionBtnDisabled: bool,
  actionBtnLabel: string,
  children: element,
  complete: bool,
  completedHeaderText: string,
  containerClassName: string,
  description: string,
  editBtnDisabled: bool,
  editBtnLabel: string,
  headerText: string.isRequired,
  onActionBtnClick: func,
  onEditBtnClick: func,
  secondary: bool,
  step: oneOfType([string, number]).isRequired,
};

Step.defaultProps = {
  actionBtnDisabled: false,
  actionBtnLabel: '',
  children: null,
  complete: false,
  completedHeaderText: '',
  containerClassName: '',
  description: '',
  editBtnDisabled: false,
  editBtnLabel: '',
  onActionBtnClick: () => {},
  onEditBtnClick: () => {},
  secondary: false,
};

export default Step;
