import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';

import WizardPage from 'shared/WizardPage';
import { loadSchema, submitForm, resetSuccess } from 'scenes/DD1299/ducks';

import Alert from 'shared/Alert';

export class DemoWorkflow extends Component {
  componentDidMount() {
    document.title = 'Transcom My Move';
    this.props.loadSchema();
  }
  submit = values => {
    this.props.submitForm(values);
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
            Your DD1299 has been sucessfully submitted.
          </Alert>
          <button type="button" onClick={this.props.resetSuccess}>
            Do another one
          </button>
        </Fragment>
      );
    const uiSchema = Object.assign({}, this.props.uiSchema, {
      order: this.props.subsetOfUiSchema,
    });
    return (
      <Fragment>
        <WizardPage
          initialValues={this.initialValues()}
          onSubmit={this.submit}
          schema={this.props.schema}
          uiSchema={uiSchema}
          pageList={this.props.pageList}
          pageKey={this.props.path}
          history={this.props.history}
        />
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
  subsetOfUiSchema: PropTypes.arrayOf(PropTypes.string),
};

function mapStateToProps(state) {
  return state.DD1299;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadSchema, submitForm, resetSuccess }, dispatch);
}

export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(DemoWorkflow),
);
