import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'react-router-dom';
import { push } from 'react-router-redux';
import generatePath from './generatePath';

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
  }
  componentDidMount() {
    window.scrollTo(0, 0);
  }

  nextPage() {
    const { pageList, pageKey, push, match: { params } } = this.props;
    const path = getNextPagePath(pageList, pageKey);
    // comes from react router redux: doing this moves to the route at path  (might consider going back to history since we need withRouter)
    push(generatePath(path, params));
  }

  previousPage() {
    const { pageList, pageKey, push, match: { params } } = this.props;
    const path = getPreviousPagePath(pageList, pageKey);
    // push comes from react router redux : doing this moves to the route at path
    push(generatePath(path, params));
  }

  render() {
    const { handleSubmit, pageKey, pageList, children } = this.props;
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
            <button onClick={this.nextPage}>Next</button>
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
  push: PropTypes.func,
  match: PropTypes.object, //from withRouter
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push }, dispatch);
}

export default withRouter(connect(null, mapDispatchToProps)(WizardPage));
