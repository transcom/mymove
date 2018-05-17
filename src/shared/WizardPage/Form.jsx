import React, { Component } from 'react'; // eslint-disable-line
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import { push } from 'react-router-redux';
import Alert from 'shared/Alert'; // eslint-disable-line
import generatePath from './generatePath';
import './index.css';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import { reduxForm } from 'redux-form';

import {
  getNextPagePath,
  getPreviousPagePath,
  isFirstPage,
  isLastPage,
} from './utils';

export class WizardFormPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.state = { transitionFunc: null };
  }
  componentDidUpdate() {
    if (this.props.hasSucceeded) this.onSubmitSuccessful();
    if (this.props.serverError) window.scrollTo(0, 0);
  }
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  beforeTransition(func) {
    const { dirty, pageList, pageKey, handleSubmit } = this.props;
    const path = func(pageList, pageKey);
    if (dirty && handleSubmit) {
      handleSubmit();
      this.setState({ transitionFunc: func });
    } else {
      this.goto(path);
    }
  }
  goto(path) {
    const {
      push,
      match: { params },
      additionalParams,
    } = this.props;
    const combinedParams = additionalParams
      ? Object.assign({}, additionalParams, params)
      : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, combinedParams));
  }
  onSubmitSuccessful() {
    const { transitionFunc } = this.state;
    const { pageKey, pageList } = this.props;
    if (transitionFunc) this.goto(transitionFunc(pageList, pageKey));
  }
  cancelFlow = () => {
    this.props.push(`/`);
  };
  nextPage() {
    this.beforeTransition(getNextPagePath);
  }

  previousPage() {
    this.beforeTransition(getPreviousPagePath);
  }

  render() {
    const {
      handleSubmit,
      className,
      pageKey,
      pageList,
      children,
      serverError,
      valid,
      dirty,
    } = this.props;
    const canMoveForward = valid;
    const canMoveBackward =
      (valid || !dirty) && !isFirstPage(pageList, pageKey);
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
          <div className="usa-width-one-third">
            <button className="usa-button-secondary" onClick={this.cancelFlow}>
              Cancel
            </button>
          </div>
          <div className="usa-width-two-thirds">
            <button
              className="usa-button-secondary"
              onClick={this.previousPage}
              disabled={!canMoveBackward}
            >
              Back
            </button>
            {!isLastPage(pageList, pageKey) && (
              <button
                className="usa-button-primary"
                onClick={this.nextPage}
                disabled={!canMoveForward}
              >
                Next
              </button>
            )}
            {isLastPage(pageList, pageKey) && (
              <button
                className="usa-button-primary"
                onClick={handleSubmit}
                disabled={!canMoveForward}
              >
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
  hasSucceeded: PropTypes.bool.isRequired,
  serverError: PropTypes.object,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  valid: PropTypes.bool,
  dirty: PropTypes.bool,
  push: PropTypes.func,
  match: PropTypes.object, //from withRouter
  additionalParams: PropTypes.object,
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

export const reduxifyWizardForm = (name, additionalValidations) => {
  let validations = validateRequiredFields;
  if (additionalValidations) {
    validations = composeValidations(
      validateRequiredFields,
      additionalValidations,
    );
  }
  return reduxForm({ form: name, validate: validations })(
    withRouter(connect(null, mapDispatchToProps)(WizardFormPage)),
  );
};
