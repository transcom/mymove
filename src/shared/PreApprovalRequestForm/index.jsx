import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm, Form } from 'redux-form';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './index.css';

const Codes = accessorials => props => {
  let value, onChange;
  if (props.input) {
    value = props.input.value;
    onChange = props.input.onChange;
  } else {
    value = props.value;
    onChange = props.onChange;
  }

  const localOnChange = event => {
    onChange(event.target.value);
  };
  return (
    <select onChange={localOnChange} value={value}>
      <option />
      {accessorials.map(e => (
        <option key={e.id} value={e.id}>
          {e.code} {e.item}
        </option>
      ))}
    </select>
  );
};

class PreApprovalRequestForm extends Component {
  render() {
    return (
      <Form onSubmit={this.props.handleSubmit(this.props.onSubmit)}>
        <div className="usa-grid">
          <div className="usa-width-one-half">
            <SwaggerField
              fieldName="accessorial"
              title="Code & Item"
              className="three-quarter-width rounded"
              component={Codes(this.props.accessorials)}
              swagger={this.props.ship_accessorial_schema}
              required
            />
            <SwaggerField
              fieldName="location"
              className="one-third-width rounded"
              swagger={this.props.ship_accessorial_schema}
              required
            />
            <SwaggerField
              fieldName="quantity_1"
              className="half-width"
              swagger={this.props.ship_accessorial_schema}
              required
            />
          </div>
          <div className="usa-width-one-half">
            <SwaggerField
              fieldName="notes"
              className="three-quarter-width"
              swagger={this.props.ship_accessorial_schema}
            />
          </div>
        </div>
      </Form>
    );
  }
}

PreApprovalRequestForm.propTypes = {
  schema: PropTypes.object,
  accessorials: PropTypes.array,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'preapproval_request_form';

PreApprovalRequestForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PreApprovalRequestForm);

function mapStateToProps(state, props) {
  return {
    ship_accessorial_schema: get(
      state,
      'swaggerPublic.spec.definitions.ShipmentAccessorial',
      {},
    ),
  };
}

export default connect(mapStateToProps)(PreApprovalRequestForm);
