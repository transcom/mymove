import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { Button } from '@trussworks/react-uswds';

import Alert from 'shared/Alert';
import generatePath from './generatePath';
import './index.css';
import { validateRequiredFields } from 'shared/JsonSchemaForm';
import styles from 'components/Customer/WizardNavigation/WizardNavigation.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';

import { beforeTransition, getNextPagePath, getPreviousPagePath, isFirstPage, isLastPage } from './utils';
import withRouter from 'utils/routing';

export class WizardFormPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.beforeTransition = beforeTransition.bind(this);
    this.submit = this.submit.bind(this);
  }

  static defaultProps = {
    readyToSubmit: true,
  };

  componentDidUpdate(prevProps) {
    if (this.props.additionalValues) {
      Object.keys(this.props.additionalValues).forEach((key) => {
        if (this.props.additionalValues[`${key}`] !== prevProps.additionalValues[`${key}`]) {
          this.props.change(key, this.props.additionalValues[`${key}`]);
        }
      });
    }
  }

  goto(path) {
    const {
      router: { navigate, params },
      additionalParams,
    } = this.props;

    const combinedParams = additionalParams ? { ...additionalParams, ...params } : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    navigate(generatePath(path, combinedParams));
  }

  nextPage() {
    if (this.props.reduxFormSubmit) {
      return this.props.reduxFormSubmit().then(() => this.beforeTransition(getNextPagePath, false));
    }
    return this.beforeTransition(getNextPagePath);
  }

  previousPage() {
    const shouldHandleSubmit = !this.props.discardOnBack;
    if (this.props.reduxFormSubmit && shouldHandleSubmit) {
      return this.props.reduxFormSubmit().then(() => this.beforeTransition(getPreviousPagePath, false));
    }
    this.beforeTransition(getPreviousPagePath, shouldHandleSubmit);
    return undefined;
  }

  submit() {
    if (this.props.reduxFormSubmit) {
      return this.props.reduxFormSubmit();
    }
    return this.props.handleSubmit();
  }

  render() {
    // when reduxFormSubmit is supplied it's expected that the form will use redux-form's handlesubmit prop
    // and accompanying submit validation https://redux-form.com/8.2.0/examples/submitvalidation/
    // while forms that provide their own handlesubmit prop are expected to not be using redux-form's submit validation
    const hasReduxFormSubmitHandler = !!this.props.reduxFormSubmit;
    const { handleSubmit, className, pageKey, pageList, children, serverError, valid, dirty, readyToSubmit } =
      this.props;
    const canMoveForward = valid && readyToSubmit;
    const canMoveBackward = (valid || !dirty) && !isFirstPage(pageList, pageKey);
    const hideBackBtn = isFirstPage(pageList, pageKey);
    return (
      <div className="grid-container usa-prose">
        <NotificationScrollToTop dependency={serverError} />
        {serverError && (
          <div className="grid-row">
            <div className="desktop:grid-col-8 desktop:grid-offset-2 error-message">
              <Alert type="error" heading="An error occurred">
                {serverError.message}
              </Alert>
            </div>
          </div>
        )}
        <div className="grid-row">
          <div className="grid-col desktop:grid-col-8 desktop:grid-offset-2">
            <form className={className}>{children}</form>
          </div>
        </div>
        <div className="grid-row" style={{ marginTop: '24px' }}>
          <div className="grid-col desktop:grid-col-8 desktop:grid-offset-2">
            <div className={styles.WizardNavigation}>
              {!hideBackBtn && (
                <Button
                  type="button"
                  secondary
                  className={styles.button}
                  onClick={hasReduxFormSubmitHandler ? handleSubmit(this.previousPage) : this.previousPage}
                  disabled={!canMoveBackward}
                  data-testid="wizardBackButton"
                >
                  Back
                </Button>
              )}

              {isLastPage(pageList, pageKey) ? (
                <Button
                  type="button"
                  className={styles.button}
                  onClick={hasReduxFormSubmitHandler ? handleSubmit(this.submit) : this.submit}
                  disabled={!canMoveForward}
                  data-testid="wizardCompleteButton"
                >
                  Complete
                </Button>
              ) : (
                <Button
                  type="button"
                  className={styles.button}
                  onClick={hasReduxFormSubmitHandler ? handleSubmit(this.nextPage) : this.nextPage}
                  disabled={!canMoveForward}
                  data-testid="wizardNextButton"
                >
                  Next
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

WizardFormPage.propTypes = {
  handleSubmit: PropTypes.func,
  reduxFormSubmit: PropTypes.func, // function supplied to use w/ redux-form's submit validation
  serverError: PropTypes.object,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  valid: PropTypes.bool,
  dirty: PropTypes.bool,
  router: PropTypes.object, //from withRouter
  additionalParams: PropTypes.object,
  additionalValues: PropTypes.object, // These values are passed into the form with change()
  discardOnBack: PropTypes.bool,
};

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

export const reduxifyWizardForm = (name, additionalValidations, asyncValidate, asyncBlurFields) => {
  let validations = validateRequiredFields;
  if (additionalValidations) {
    validations = composeValidations(validateRequiredFields, additionalValidations);
  }
  return reduxForm({
    form: name,
    validate: validations,
    asyncValidate,
    asyncBlurFields,
    enableReinitialize: true,
    keepDirtyOnReinitialize: true,
  })(withRouter(WizardFormPage));
};
