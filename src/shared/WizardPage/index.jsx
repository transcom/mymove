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
      formName,
    } = this.props;
    const CurrentForm = reduxifyForm(formName);
    return (
      <div className="usa-grid">
        <CurrentForm
          className="usa-width-one-whole"
          schema={schema}
          uiSchema={uiSchema}
          initialValues={initialValues}
          showSubmit={false}
          destroyOnUnmount={false}
        />
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
            <button
              onClick={this.nextPage}
              disabled={isLastPage(pageList, pageKey)}
            >
              Next
            </button>
          )}
          {isLastPage(pageList, pageKey) && (
            <button onClick={onSubmit}>Complete</button>
          )}
        </div>
      </div>
    );
  }
}

WizardPage.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  formName: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
  pageKey: PropTypes.string.isRequired,
};

export default WizardPage;
