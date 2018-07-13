import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import { sendGexRequest } from './ducks';

const schema = {
  properties: {
    transaction_name: {
      type: 'string',
      title: 'Transaction Name',
      example: '858-1.edi',
    },
    transaction_body: {
      type: 'string',
      title: 'Transaction Body',
      format: 'textarea',
    },
  },
};

export class Gex extends Component {
  sendRequest = values => {
    this.props.sendGexRequest(values);
  };

  render() {
    return (
      <div className="usa-grid">
        <h1 className="sm-heading">Send GEX Request</h1>
        <form onSubmit={this.props.handleSubmit(this.sendRequest)}>
          <SwaggerField
            fieldName="transaction_name"
            swagger={this.props.schema}
            required
          />
          <SwaggerField
            fieldName="transaction_body"
            swagger={this.props.schema}
            required
          />
          <button type="submit">Submit</button>
        </form>
      </div>
    );
  }
}

Gex.propTypes = {
  sendGexRequest: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  return {
    ...state.gex,
    schema,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ sendGexRequest }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(
  reduxForm({ form: 'gex' })(Gex),
);
