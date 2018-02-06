import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { reduxifyForm } from 'shared/JsonSchemaForm';
import { loadSchema, loadUiSchema } from './ducks';

const DD1299Form = reduxifyForm('DD1299');
export class DD1299 extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: DD1299';
    this.props.loadSchema();
    this.props.loadUiSchema();
  }
  submit = values => {
    // print the form values to the console
    console.log(values);
  };
  render() {
    return (
      <DD1299Form
        onSubmit={this.submit}
        schema={this.props.schema}
        uiSchema={this.props.uiSchema}
      />
    );
  }
}
DD1299.propTypes = {
  loadSchema: PropTypes.func.isRequired,
  schema: PropTypes.object,
  hasError: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return {
    schema: state.DD1299.schema,
    hasError: state.DD1299.hasError,
    uiSchema: state.DD1299.uiSchema,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadSchema, loadUiSchema }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DD1299);
