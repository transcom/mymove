import React, { Fragment, Component } from 'react';
import PropTypes from 'prop-types';
import reduxifyForm from '/shared/JsonSchemaForm';
import classnames from 'classnames';

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
    const { pageList, pageKey, router } = this.props;

    //     const path = getNextPagePath(pageList, pageKey);
    //     this.onSubmit().then(() => router.push(path));
  }

  previousPage() {
    const { pageList, pageKey } = this.props;
    // const path = getPreviousPagePath(pageList, location.pathname); //see vets routing

    // this.props.router.push(path);
  }

  render() {
    const { handleSubmit, schema, uiSchema, pageKey } = this.props;
    const CurrentForm = reduxifyForm(pageKey);
    return (
      <form className="default" onSubmit={handleSubmit}>
        <CurrentForm schema={schema} uiSchema={uiSchema} showSubmit={false} />
        {!isFirstPage && (
          <button
            className={classnames({ 'usa-button-secondary': !isLastPage })}
            onClick={this.previousPage}
          >
            Prev
          </button>
        )}
        {!isLastPage && <button onClick={this.nextPage}>Next</button>}
      </form>
    );
  }
}

WizardPage.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pages: PropTypes.array.isRequired,
  pageKey: PropTypes.string.isRequired,
};
