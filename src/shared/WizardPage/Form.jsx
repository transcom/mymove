import React, { Component } from 'react'; // eslint-disable-line
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import windowSize from 'react-window-size';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import { push } from 'react-router-redux';
import Alert from 'shared/Alert'; // eslint-disable-line
import generatePath from './generatePath';
import './index.css';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import { reduxForm } from 'redux-form';
import { mobileSize } from 'shared/constants';
import scrollToTop from 'shared/scrollToTop';

import { getNextPagePath, getPreviousPagePath, isFirstPage, isLastPage, beforeTransition } from './utils';

export class WizardFormPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.cancelFlow = this.cancelFlow.bind(this);
    this.beforeTransition = beforeTransition.bind(this);
  }
  componentDidUpdate(prevProps) {
    if (this.props.additionalValues) {
      /* eslint-disable security/detect-object-injection */

      Object.keys(this.props.additionalValues).forEach(key => {
        if (this.props.additionalValues[key] !== prevProps.additionalValues[key]) {
          this.props.change(key, this.props.additionalValues[key]);
        }
      });
    }
    /* eslint-enable security/detect-object-injection */

    if (this.props.serverError) scrollToTop();
  }
  componentDidMount() {
    scrollToTop();
  }
  goto(path) {
    const {
      push,
      match: { params },
      additionalParams,
    } = this.props;
    const combinedParams = additionalParams ? Object.assign({}, additionalParams, params) : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, combinedParams));
  }

  cancelFlow() {
    this.props.push(`/`);
  }
  nextPage() {
    this.beforeTransition(getNextPagePath);
  }

  previousPage() {
    const shouldHandleSubmit = !this.props.discardOnBack;
    this.beforeTransition(getPreviousPagePath, shouldHandleSubmit);
  }

  render() {
    const isMobile = this.props.windowWidth < mobileSize;
    const { handleSubmit, className, pageKey, pageList, children, serverError, valid, dirty } = this.props;
    const canMoveForward = valid;
    const canMoveBackward = (valid || !dirty) && !isFirstPage(pageList, pageKey);
    const hideBackBtn = isFirstPage(pageList, pageKey);
    return (
      <div className="usa-grid">
        {serverError && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {serverError.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <form className={className} onSubmit={handleSubmit}>
            {children}
          </form>
        </div>
        <div className="usa-width-one-whole lower-nav-btns">
          {!isMobile && (
            <div className="left cancel">
              <button className="usa-button-secondary" onClick={this.cancelFlow}>
                Cancel
              </button>
            </div>
          )}
          <div className="prev-next">
            <button
              className={'usa-button-secondary prev ' + (hideBackBtn && 'hide-btn')}
              onClick={this.previousPage}
              disabled={!canMoveBackward}
            >
              Back
            </button>
            {!isLastPage(pageList, pageKey) && (
              <button className="usa-button-primary next" onClick={this.nextPage} disabled={!canMoveForward}>
                Next
              </button>
            )}
            {isLastPage(pageList, pageKey) && (
              <button className="usa-button-primary next" onClick={handleSubmit} disabled={!canMoveForward}>
                Complete
              </button>
            )}
          </div>
        </div>
      </div>
    );
  }
}

WizardFormPage.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  serverError: PropTypes.object,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  valid: PropTypes.bool,
  dirty: PropTypes.bool,
  push: PropTypes.func,
  match: PropTypes.object, //from withRouter
  additionalParams: PropTypes.object,
  additionalValues: PropTypes.object, // These values are passed into the form with change()
  windowWidth: PropTypes.number,
  discardOnBack: PropTypes.bool,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}

function composeValidations(initialValidations, additionalValidations) {
  return (values, form) => {
    let errors = initialValidations(values, form);
    if (errors === undefined) {
      errors = {};
    }
    const additionalErrors = additionalValidations(values, form);
    Object.assign(errors, additionalErrors);
    return errors;
  };
}

const wizardFormPageWithSize = windowSize(WizardFormPage);

export const reduxifyWizardForm = (name, additionalValidations) => {
  let validations = validateRequiredFields;
  if (additionalValidations) {
    validations = composeValidations(validateRequiredFields, additionalValidations);
  }
  return reduxForm({
    form: name,
    validate: validations,
    enableReinitialize: true,
    keepDirtyOnReinitialize: true,
  })(withRouter(connect(null, mapDispatchToProps)(wizardFormPageWithSize)));
};
