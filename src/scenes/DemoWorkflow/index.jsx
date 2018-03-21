import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { reduxifyForm } from 'shared/JsonSchemaForm';

import WizardPage from 'shared/WizardPage';
import { loadSchema, submitForm, resetSuccess } from 'scenes/DD1299/ducks';

import Alert from 'shared/Alert';

export class DemoWorkflow extends Component {
  componentDidMount() {
    document.title = 'Transcom My Move';
    this.props.loadSchema();
  }
  submit = () => {
    this.props.submitForm();
  };
  initialValues() {
    return {
      mobile_home_contents_packed_requested: false,
      mobile_home_blocked_requested: false,
      mobile_home_unblocked_requested: false,
      mobile_home_stored_at_origin_requested: false,
      mobile_home_stored_at_destination_requested: false,
    };
  }
  componentDidUpdate() {
    if (this.props.hasSubmitSuccess) window.scrollTo(0, 0);
  }
  render() {
    if (this.props.hasSchemaError)
      return (
        <Alert type="error" heading="Server Error">
          There was a problem loading the form from the server.
        </Alert>
      );
    if (this.props.hasSubmitSuccess)
      return (
        <Fragment>
          <Alert type="success" heading="Form Submitted">
            Your DD1299 has been successfully submitted.
          </Alert>
          <button type="button" onClick={this.props.resetSuccess}>
            Do another one
          </button>
        </Fragment>
      );

    const CurrentForm = reduxifyForm('DD1299');
    const isValid = this.refs.currentForm && this.refs.currentForm.valid;
    return (
      <Fragment>
        <WizardPage
          handleSubmit={this.submit}
          pageList={this.props.pageList}
          pageKey={this.props.path}
          history={this.props.history}
          pageIsValid={isValid}
        >
          <CurrentForm
            ref="currentForm"
            className="usa-width-one-whole"
            schema={this.props.schema}
            uiSchema={this.props.uiSchema}
            initialValues={this.initialValues()}
            showSubmit={false}
            destroyOnUnmount={false}
            subsetOfUiSchema={this.props.subsetOfUiSchema}
          />
        </WizardPage>
        {this.props.hasSubmitError && (
          <Alert type="error" heading="Server Error">
            There was a problem saving the form on the server.
          </Alert>
        )}
      </Fragment>
    );
  }
}

DemoWorkflow.propTypes = {
  loadSchema: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  hasSchemaError: PropTypes.bool.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.DD1299;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadSchema, submitForm, resetSuccess }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DemoWorkflow);
