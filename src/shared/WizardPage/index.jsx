import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import { push } from 'react-router-redux';
import Alert from 'shared/Alert';
import generatePath from './generatePath';
import './index.css';

import {
  getNextPagePath,
  getPreviousPagePath,
  isFirstPage,
  isLastPage,
} from './utils';

export class WizardPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.state = { transitionTo: null };
  }
  componentDidUpdate() {
    if (this.props.hasSucceeded) this.onSubmitSuccessful();
    if (this.props.error) window.scrollTo(0, 0);
  }
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  beforeTransition(func) {
    const {
      isAsync,
      pageIsDirty,
      pageList,
      pageKey,
      handleSubmit,
    } = this.props;
    const path = func(pageList, pageKey);
    if (pageIsDirty && handleSubmit) {
      handleSubmit();
      if (isAsync) {
        this.setState({ transitionTo: path });
      } else {
        this.goto(path);
      }
    } else {
      this.goto(path);
    }
  }
  goto(path) {
    const {
      push,
      match: { params },
    } = this.props;
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, params));
  }
  onSubmitSuccessful() {
    if (this.state.transitionTo) this.goto(this.state.transitionTo);
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
      pageIsDirty,
    } = this.props;
    const canMoveForward = pageIsValid;
    const canMoveBackward =
      (pageIsValid || !pageIsDirty) && !isFirstPage(pageList, pageKey);
    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">{children}</div>
        <div className="usa-width-one-whole lower-nav-btns">
          <div className="usa-width-one-third">
            <button
              className="usa-button-secondary"
              onClick={this.previousPage}
              disabled={!canMoveBackward}
            >
              Prev
            </button>
          </div>
          <div className="usa-width-one-third center">
            <button className="usa-button-secondary" disabled={true}>
              Save for later
            </button>
          </div>
          <div className="usa-width-one-third right-align">
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

WizardPage.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  isAsync: PropTypes.bool.isRequired,
  hasSucceeded: (props, propName, componentName) => {
    if (props['isAsync'] && typeof props[propName] !== 'boolean') {
      return new Error('Async WizardPages must have hasSucceeded boolean prop');
    }
  },
  error: PropTypes.object,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  pageIsValid: PropTypes.bool,
  pageIsDirty: PropTypes.bool,
  push: PropTypes.func,
  match: PropTypes.object, //from withRouter
};

WizardPage.defaultProps = {
  isAsync: false,
  pageIsValid: true,
  pageIsDirty: true,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}

export default withRouter(connect(null, mapDispatchToProps)(WizardPage));
