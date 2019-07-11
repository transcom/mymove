// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './index.css';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';
import { withRouter } from 'react-router-dom';

//this is taken from https://designsystem.digital.gov/components/alerts/
class AlertWithConfirmation extends Component {
  cancelActionHandler = () => {
    return this.props.cancelActionHandler;
  };

  confirmActionHandler = () => {
    return this.props.okActionHandler;
  };

  render() {
    return (
      <div className="usa-width-one-whole">
        <div className={`usa-alert usa-alert-${this.props.type}`}>
          <div className="usa-alert-body">
            <div className="body--heading">
              <div>
                <div>
                  {this.props.heading && <h3 className="usa-alert-heading">{this.props.heading}</h3>}
                  {this.props.onRemove && (
                    <FontAwesomeIcon
                      className="icon remove-icon actionable actionable-secondary"
                      onClick={this.props.onRemove}
                      icon={faTimes}
                    />
                  )}
                </div>
                <div className="usa-alert-text">{this.props.message}</div>
                <div className="cancel-or-ok-buttons">
                  <button type="button" className="usa-button-secondary" onClick={this.cancelActionHandler()}>
                    Cancel
                  </button>
                  <button type="button" className="usa-button" onClick={this.confirmActionHandler()}>
                    OK
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

const requiredPropsCheck = (props, propName, componentName) => {
  if (!props.message) {
    return new Error(`Message is required by '${componentName}' component.`);
  }
};

AlertWithConfirmation.propTypes = {
  message: requiredPropsCheck,
  onRemove: PropTypes.func,
  type: PropTypes.oneOf(['error', 'warning', 'info', 'success', 'loading']),
};
export default withRouter(AlertWithConfirmation);
