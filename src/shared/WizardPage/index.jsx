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
    this.state = {
      currentPageIndex: 0,
    };
  }
  nextPage() {
    const { pageList, pageKey, router, onSubmit } = this.props;

    const path = getNextPagePath(pageList, pageKey);
    router.push(path);
  }

  previousPage() {
    const { pageList, pageKey, router } = this.props;
    const path = getPreviousPagePath(pageList, pageKey); //see vets routing

    router.push(path);
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
        {!isFirstPage && (
          <button
            className={classnames({
              'usa-button-secondary': !isLastPage(pageList, pageKey),
            })}
            onClick={this.previousPage}
          >
            Prev
          </button>
        )}
        {!isLastPage(pageList, pageKey) && (
          <button onClick={this.nextPage}>Next</button>
        )}
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
