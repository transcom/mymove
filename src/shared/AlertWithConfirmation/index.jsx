import React from 'react';
import PropTypes from 'prop-types';
import './index.css';

export function AlertWithConfirmation(props) {
  return (
    <div className="grid-row">
      <div className="grid-col-12">
        <div className={`usa-alert usa-alert--${props.type}`}>
          <div className="usa-alert__body">
            <div className="body--heading">
              <div>{props.heading && <h3 className="usa-alert__heading">{props.heading}</h3>}</div>
              <div className="grid-row grid-gap">
                <div className="grid-col-9">
                  <div className="usa-alert__text">{props.message}</div>
                </div>
                <div className="grid-col-3 text-right">
                  <div className="cancel-or-ok-buttons">
                    <button
                      type="button"
                      className="usa-button usa-button--secondary"
                      onClick={props.cancelActionHandler}
                    >
                      Cancel
                    </button>
                    <button type="button" className="usa-button" onClick={props.okActionHandler}>
                      OK
                    </button>
                  </div>
                </div>
              </div>
            </div>
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
  type: PropTypes.oneOf(['error', 'warning', 'info', 'success']),
};
export default AlertWithConfirmation;
