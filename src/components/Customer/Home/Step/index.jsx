/* eslint-disable react/prop-types */
import React from 'react';
import { bool, node, string, oneOfType, number, func } from 'prop-types';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import styles from './Step.module.scss';

import { ReactComponent as AcceptIcon } from 'shared/icon/accept-inversed.svg';

const NumberCircle = ({ num }) => <div className={styles['number-circle']}>{num}</div>;

NumberCircle.propTypes = {
  num: string.isRequired,
};

const Step = ({
  actionBtnDisabled,
  actionBtnLabel,
  actionBtnId,
  children,
  complete,
  completedHeaderText,
  containerClassName,
  editBtnDisabled,
  editBtnLabel,
  headerText,
  onActionBtnClick,
  onEditBtnClick,
  secondaryBtn,
  secondaryBtnClassName,
  step,
}) => {
  const showThoughNotFunctional = false; // remove when all Edit buttons work
  const actionBtnClassName = classnames(
    styles['action-btn'],
    {
      [styles['action-button--secondary']]: secondaryBtn,
    },
    secondaryBtnClassName,
  );

  return (
    <div data-testid={`stepContainer${step}`} className={`${styles['step-container']} ${containerClassName}`}>
      <div className={styles['step-header-container']}>
        {complete ? <AcceptIcon aria-hidden className={styles.accept} /> : <NumberCircle num={step} />}
        <strong>{complete ? completedHeaderText : headerText}</strong>
        {showThoughNotFunctional && editBtnLabel && (
          <Button className={styles['edit-btn']} disabled={editBtnDisabled} onClick={onEditBtnClick} type="button">
            {editBtnLabel}
          </Button>
        )}
      </div>

      {children}
      {actionBtnLabel && (
        <Button
          className={actionBtnClassName}
          disabled={actionBtnDisabled}
          data-testid={actionBtnId}
          onClick={onActionBtnClick}
          type="button"
          secondary={secondaryBtn}
        >
          {actionBtnLabel}
        </Button>
      )}
    </div>
  );
};

Step.propTypes = {
  actionBtnDisabled: bool,
  actionBtnId: string,
  actionBtnLabel: string,
  children: node,
  complete: bool,
  completedHeaderText: string,
  containerClassName: string,
  editBtnDisabled: bool,
  editBtnLabel: string,
  headerText: string.isRequired,
  onActionBtnClick: func,
  onEditBtnClick: func,
  secondaryBtn: bool,
  secondaryBtnClassName: string,
  step: oneOfType([string, number]).isRequired,
};

Step.defaultProps = {
  actionBtnDisabled: false,
  actionBtnId: 'button',
  actionBtnLabel: '',
  children: null,
  complete: false,
  completedHeaderText: '',
  containerClassName: '',
  editBtnDisabled: false,
  editBtnLabel: '',
  onActionBtnClick: () => {},
  onEditBtnClick: () => {},
  secondaryBtn: false,
  secondaryBtnClassName: '',
};

export default Step;
