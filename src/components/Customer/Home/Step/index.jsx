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
  step,
}) => {
  const showThoughNotFunctional = false; // remove when all Edit buttons work
  const actionBtnClassName = classnames(styles['action-btn'], {
    [styles['usa-button--secondary']]: secondaryBtn,
  });

  return (
    <div className={`${styles['step-container']} ${containerClassName}`}>
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
        <Button className={actionBtnClassName} disabled={actionBtnDisabled} onClick={onActionBtnClick} type="button">
          {actionBtnLabel}
        </Button>
      )}
    </div>
  );
};

Step.propTypes = {
  actionBtnDisabled: bool,
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
  step: oneOfType([string, number]).isRequired,
};

Step.defaultProps = {
  actionBtnDisabled: false,
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
};

export default Step;
