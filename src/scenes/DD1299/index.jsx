import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { reduxifyForm } from 'shared/JsonSchemaForm';
import { submitForm, resetSuccess } from './ducks';

import Alert from 'shared/Alert';

const DD1299Form = reduxifyForm('DD1299');
export class DD1299 extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: DD1299';
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
            Your DD1299 has been successfully submitted.
          </Alert>
          <button type="button" onClick={this.props.resetSuccess}>
            Do another one
          </button>
        </Fragment>
      );
    return (
      <Fragment>
        <DD1299Form
          initialValues={this.initialValues()}
          onSubmit={this.submit}
          schema={this.props.schema}
          uiSchema={this.props.uiSchema}
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

DD1299.propTypes = {
  submitForm: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  hasSchemaError: PropTypes.bool.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = { ...state.DD1299, schema: {} };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateForm1299Payload;
  }
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ submitForm, resetSuccess }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DD1299);
