import React, { Fragment, Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import { reduxifyForm } from '../JsonSchemaForm';
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
    window.scrollTo(0, 0);
  }

  nextPage() {
    const { pageList, pageKey, history, onSubmit } = this.props;
    const path = getNextPagePath(pageList, pageKey);
    history.push(path);
  }

  previousPage() {
    const { pageList, pageKey, history } = this.props;
    const path = getPreviousPagePath(pageList, pageKey); //see vets routing

    history.push(path);
  }

  render() {
    const {
      onSubmit,
      schema,
      uiSchema,
      initialValues,
      pageKey,
      pageList,
    } = this.props;
    const CurrentForm = reduxifyForm(pageKey);
    return (
      <Fragment>
        <CurrentForm
          schema={schema}
          uiSchema={uiSchema}
          initialValues={initialValues}
          showSubmit={false}
        />
        <div className="usa-grid">
          <div className="usa-width-one-third">
            {!isFirstPage(pageList, pageKey) && (
              <button
                className={classnames({
                  'usa-button-secondary': !isLastPage(pageList, pageKey),
                })}
                onClick={this.previousPage}
              >
                Prev
              </button>
            )}
          </div>
          <div className="usa-width-one-third" />
          <div className="usa-width-one-third">
            {!isLastPage(pageList, pageKey) && (
              <button onClick={this.nextPage}>Next</button>
            )}
            {isLastPage(pageList, pageKey) && (
              <button onClick={onSubmit}>Complete</button>
            )}
          </div>
        </div>
      </Fragment>
    );
  }
}

WizardPage.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
};

export default WizardPage;
