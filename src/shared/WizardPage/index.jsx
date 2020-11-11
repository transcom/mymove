import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { push } from 'connected-react-router';

import ScrollToTop from 'components/ScrollToTop';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Alert from 'shared/Alert';
import generatePath from './generatePath';
import './index.css';

import { getNextPagePath, getPreviousPagePath, isFirstPage, isLastPage } from './utils';

/**
 * TODO:
 * - verify next/previous actions
 * - style buttons - WIP need design input
 */

export const WizardPage = (props) => {
  const {
    push,
    match,
    additionalParams,
    pageList,
    pageKey,
    dirty,
    handleSubmit,
    children,
    error,
    pageIsValid,
    canMoveNext,
    hideBackBtn,
    showFinishLaterBtn,
    footerText,
  } = props;

  const goHome = () => {
    push(`/`);
  };

  const goto = (path) => {
    const { params } = match;
    const combinedParams = additionalParams ? { ...additionalParams, ...params } : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, combinedParams));
  };

  const nextPage = async () => {
    if (isLastPage(pageList, pageKey)) return handleSubmit();

    if (dirty && handleSubmit) {
      const awaitSubmit = await handleSubmit(); // wait for API save
      if (awaitSubmit?.error) {
        console.error('Wizard submit error', awaitSubmit.error);
        return;
      }
    }

    const path = getNextPagePath(pageList, pageKey);
    if (path) goto(path);
  };

  const previousPage = () => {
    // Don't submit or validate when going back
    const path = getPreviousPagePath(pageList, pageKey);
    if (path) goto(path);
  };

  const canMoveForward = pageIsValid && canMoveNext;

  return (
    <div className="grid-container usa-prose">
      <ScrollToTop />
      {error && (
        <div className="grid-row">
          <div className="grid-col-12 error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        </div>
      )}
      <div className="grid-row">
        <div className="grid-col">{children}</div>
      </div>
      <div className="grid-row" style={{ marginTop: '24px' }}>
        <div className="grid-col">
          {footerText && footerText}
          <WizardNavigation
            isFirstPage={isFirstPage(pageList, pageKey) || hideBackBtn}
            isLastPage={isLastPage(pageList, pageKey)}
            disableNext={!canMoveForward}
            showFinishLater={showFinishLaterBtn}
            onBackClick={previousPage}
            onNextClick={nextPage}
            onCancelClick={goHome}
          />
        </div>
      </div>
    </div>
  );
};

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
