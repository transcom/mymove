import React, { Component } from 'react';
import PropTypes from 'prop-types';

import {
  getNextPagePath,
  getPreviousPagePath,
  isFirstPage,
  isLastPage,
} from './utils';

class WizardPage extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
  }
  componentDidMount() {
    //  window.scrollTo(0, 0);
  }

  nextPage() {
    const { pageList, pageKey, history } = this.props;
    const path = getNextPagePath(pageList, pageKey);
    // history comes from react router: doing this moves to the route at path
    history.push(path);
  }

  previousPage() {
    const { pageList, pageKey, history } = this.props;
    const path = getPreviousPagePath(pageList, pageKey);
    // history comes from react router: doing this moves to the route at path
    history.push(path);
  }

  render() {
    const {
      handleSubmit,
      pageKey,
      pageList,
      pageIsValid,
      children,
    } = this.props;
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">{children}</div>
        <div className="usa-width-one-third">
          <button
            className="usa-button-secondary"
            onClick={this.previousPage}
            disabled={isFirstPage(pageList, pageKey)}
          >
            Prev
          </button>
        </div>
        <div className="usa-width-one-third">
          <button className="usa-button-secondary" disabled={true}>
            Save for later
          </button>
        </div>
        <div className="usa-width-one-third">
          {!isLastPage(pageList, pageKey) && (
            <button onClick={this.nextPage} disabled={!pageIsValid}>
              Next
            </button>
          )}
          {isLastPage(pageList, pageKey) && (
            <button onClick={handleSubmit}>Complete</button>
          )}
        </div>
      </div>
    );
  }
}

WizardPage.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
  pageIsValid: PropTypes.bool,
};
WizardPage.defaultProps = {
  pageIsValid: true,
};
export default WizardPage;
