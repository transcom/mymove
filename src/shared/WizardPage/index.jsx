import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { push } from 'connected-react-router';

import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Alert from 'shared/Alert';
import generatePath from './generatePath';
import './index.css';

import { getNextPagePath, getPreviousPagePath, isFirstPage, isLastPage, beforeTransition } from './utils';

export class WizardPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.goHome = this.goHome.bind(this);
    this.beforeTransition = beforeTransition.bind(this);
  }

  goHome() {
    this.props.push(`/`);
  }

  goto(path) {
    const {
      push,
      match: { params },
      additionalParams,
    } = this.props;
    const combinedParams = additionalParams ? { ...additionalParams, ...params } : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, combinedParams));
  }

  nextPage() {
    this.beforeTransition(getNextPagePath);
  }

  previousPage() {
    const { pageList, pageKey } = this.props;
    // Don't submit or validate when going back
    const path = getPreviousPagePath(pageList, pageKey);
    if (path) this.goto(path);
  }

  render() {
    const {
      pageList,
      pageKey,
      children,
      error,
      pageIsValid,
      canMoveNext,
      hideBackBtn,
      showFinishLaterBtn,
      footerText,
    } = this.props;

    const canMoveForward = pageIsValid && canMoveNext;

    return (
      <div className="grid-container usa-prose">
        {error && (
          <div className="grid-row">
            <div className="desktop:grid-col-8 desktop:grid-offset-2 error-message">
              <Alert type="error" heading="An error occurred">
                {error?.message || error}
              </Alert>
            </div>
          </div>
        )}
        <div className="grid-row">
          <div className="grid-col desktop:grid-col-8 desktop:grid-offset-2">{children}</div>
        </div>
        <div className="grid-row" style={{ marginTop: '24px' }}>
          <div className="grid-col desktop:grid-col-8 desktop:grid-offset-2">
            {footerText && footerText}
            <WizardNavigation
              isFirstPage={isFirstPage(pageList, pageKey) || hideBackBtn}
              isLastPage={isLastPage(pageList, pageKey)}
              disableNext={!canMoveForward}
              showFinishLater={showFinishLaterBtn}
              onBackClick={this.previousPage}
              onNextClick={this.nextPage}
              onCancelClick={this.goHome}
            />
          </div>
        </div>
      </div>
    );
  }
}

WizardPage.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  error: PropTypes.object,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  pageIsValid: PropTypes.bool,
  canMoveNext: PropTypes.bool,
  dirty: PropTypes.bool,
  push: PropTypes.func,
  match: PropTypes.object, //from withRouter
  additionalParams: PropTypes.object,
  footerText: PropTypes.node,
};

WizardPage.defaultProps = {
  pageIsValid: true,
  canMoveNext: true,
  dirty: true,
  hideBackBtn: false,
  showFinishLaterBtn: false,
  footerText: null,
};

const mapDispatchToProps = {
  push,
};

export default withRouter(connect(null, mapDispatchToProps)(WizardPage));
