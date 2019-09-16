// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import './index.css';

export function AlertWithConfirmation(props) {
  return (
    <div className="usa-alert usa-alert-warning usa-width-two-thirds">
      <div className="usa-alert-body usa-width-one-whole">
        <div className="body--heading">
          <div>{props.heading && <h3 className="usa-alert-heading">{props.heading}</h3>}</div>
          <div className="usa-alert-text">{props.message}</div>
          <div className="delete-or-cancel-buttons">
            <button type="button" className="usa-button" onClick={props.deleteActionHandler}>
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

AlertWithConfirmation.propTypes = {
  heading: requiredPropsCheck,
  message: requiredPropsCheck,
};
export default AlertWithConfirmation;
