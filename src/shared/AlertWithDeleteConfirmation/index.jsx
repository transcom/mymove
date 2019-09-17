// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';
import styles from './index.module.scss';

export function AlertWithDeleteConfirmation(props) {
  return (
    <div
      className={`usa-alert usa-alert-warning usa-width-two-thirds ${styles['delete-alert']} ${
        styles[`${props.type}`]
      }`}
    >
      <div className={`usa-alert-body ${styles['delete-alert-body']} usa-width-one-whole`}>
        <div className={styles['delete-body--heading']}>
          <div>{props.heading && <h3 className={styles['delete-alert-heading']}>{props.heading}</h3>}</div>
          <div className="usa-alert-text">{props.message}</div>
          <div className={styles['delete-or-cancel-buttons']}>
            <button
              type="button"
              className={`usa-button ${styles['delete-button']}`}
              onClick={props.deleteActionHandler}
            >
              Delete
            </button>
            <button type="button" className="usa-button-secondary" onClick={props.cancelActionHandler}>
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

const requiredPropsCheck = (props, propName, componentName) => {
  if (!props.heading || !props.message) {
    return new Error(`A heading or message is required by '${componentName}' component.`);
  }
};

AlertWithDeleteConfirmation.propTypes = {
  heading: requiredPropsCheck,
  message: requiredPropsCheck,
  type: PropTypes.oneOf(['weight-ticket-list-alert', 'expense-ticket-list-alert']),
};
export default AlertWithDeleteConfirmation;
