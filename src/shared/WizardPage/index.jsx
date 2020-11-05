import React, { Component } from 'react'; // eslint-disable-line
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import { push } from 'connected-react-router';
import Alert from 'shared/Alert'; // eslint-disable-line
import generatePath from './generatePath';
import './index.css';
import scrollToTop from 'shared/scrollToTop';

import { getNextPagePath, getPreviousPagePath, isFirstPage, isLastPage, beforeTransition } from './utils';

export class WizardPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.goHome = this.goHome.bind(this);
    this.beforeTransition = beforeTransition.bind(this);
  }
  componentDidUpdate() {
    if (this.props.error) scrollToTop();
  }
  componentDidMount() {
    scrollToTop();
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
    const combinedParams = additionalParams ? Object.assign({}, additionalParams, params) : params;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, combinedParams));
  }

  nextPage() {
    this.beforeTransition(getNextPagePath);
  }

  previousPage() {
    this.beforeTransition(getPreviousPagePath);
  }

  render() {
    const {
      handleSubmit,
      pageKey,
      pageList,
      children,
      error,
      pageIsValid,
      dirty,
      canMoveNext,
      hideBackBtn,
      showFinishLaterBtn,
      footerText,
    } = this.props;
    const canMoveForward = pageIsValid && canMoveNext;
    const canMoveBackward = (pageIsValid || !dirty) && !isFirstPage(pageList, pageKey);
    return (
      <div className="grid-container usa-prose">
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
          {footerText && footerText}
          {!isFirstPage(pageList, pageKey) && !hideBackBtn && (
            <button
              type="button"
              className="usa-button usa-button--secondary margin-right-0"
              onClick={this.previousPage}
              disabled={!canMoveBackward}
              data-testid="wizardBackButton"
            >
              Back
            </button>
          )}
          {!isLastPage(pageList, pageKey) && (
            <button
              type="button"
              className="usa-button margin-right-0 next"
              onClick={this.nextPage}
              disabled={!canMoveForward}
              data-testid="wizardNextButton"
            >
              Next
            </button>
          )}
          {isLastPage(pageList, pageKey) && (
            <button
              type="button"
              className="usa-button margin-right-0"
              onClick={handleSubmit}
              disabled={!canMoveForward}
              data-testid="wizardCompleteButton"
            >
              Complete
            </button>
          )}
          {showFinishLaterBtn && (
            <button
              type="button"
              className="usa-button usa-button--unstyled finish-later margin-right-0"
              onClick={this.goHome}
              disabled={false}
              data-testid="wizardFinishLaterButton"
            >
              Finish later
            </button>
          )}
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
  footerText: '',
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}

export default withRouter(connect(null, mapDispatchToProps)(WizardPage));
