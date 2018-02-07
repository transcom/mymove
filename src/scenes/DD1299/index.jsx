import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { reduxifyForm } from 'shared/JsonSchemaForm';
import { loadSchema, createForm } from './ducks';

import Alert from 'shared/Alert';

const DD1299Form = reduxifyForm('DD1299');
export class DD1299 extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: DD1299';
    this.props.loadSchema();
  }
  submit = values => {
    this.props.createForm(values);
  };
  render() {
    if (this.props.hasError)
      return (
        <Alert type="error" heading="Server Error">
          There was a problem loading the form from the server.
        </Alert>
      );
    if (this.props.hasCreateSuccess)
      return (
        <Fragment>
          <Alert type="success" heading="Form Submitted">
            Your DD1299 has been sucessfully submitted.
          </Alert>
        </Fragment>
      );
    return (
      <Fragment>
        <DD1299Form
          onSubmit={this.submit}
          schema={this.props.schema}
          uiSchema={this.props.uiSchema}
        />
        {this.props.hasCreateError && (
          <Alert type="error" heading="Server Error">
            There was a problem saving the form on the server.
          </Alert>
        )}
      </Fragment>
    );
  }
}
DD1299.propTypes = {
  loadSchema: PropTypes.func.isRequired,
  schema: PropTypes.object,
  hasError: PropTypes.bool.isRequired,
  hasCreateError: PropTypes.bool.isRequired,
  hasCreateSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return state.DD1299;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadSchema, createForm }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DD1299);
