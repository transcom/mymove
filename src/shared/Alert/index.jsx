// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';
import './index.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import faSpinner from '@fortawesome/free-solid-svg-icons/faSpinner';
import faTimes from '@fortawesome/free-solid-svg-icons/faTimes';

//this is taken from https://designsystem.digital.gov/components/alerts/
const Alert = (props) => (
  <div className={`usa-alert usa-alert--${props.type}`}>
    <div className="usa-alert__body">
      <div className="body--heading">
        {props.type === 'loading' ? (
          <div className="heading--icon">
            <FontAwesomeIcon icon={faSpinner} spin pulse size="2x" />
          </div>
        ) : null}
        <div>
          {props.heading && <h3 className="usa-alert__heading">{props.heading}</h3>}
          {props.onRemove && (
            <FontAwesomeIcon
              className="icon remove-icon actionable actionable-secondary"
              onClick={props.onRemove}
              icon={faTimes}
            />
          )}
          <div className="usa-alert__text">{props.children}</div>
        </div>
      </div>
    </div>
  </div>
);

const requiredPropsCheck = (props, propName, componentName) => {
  if (!props.heading && !props.children) {
    return new Error(`One of 'heading' or 'children' is required by '${componentName}' component.`);
  }
};

Alert.propTypes = {
  heading: requiredPropsCheck,
  children: requiredPropsCheck,
  onRemove: PropTypes.func,
  type: PropTypes.oneOf(['error', 'warning', 'info', 'success', 'loading']),
};
export default Alert;
