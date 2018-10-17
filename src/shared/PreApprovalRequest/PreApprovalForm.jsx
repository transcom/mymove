import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm, Form } from 'redux-form';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './PreApprovalRequest.css';

const Codes = tariff400ngItems => props => {
  const value = props.value;
  const onChange = props.onChange;

  const localOnChange = event => {
    onChange(event.target.value);
  };
  return (
    <select onChange={localOnChange} value={value}>
      <option />
      {tariff400ngItems.map(e => (
        <option key={e.id} value={e.id}>
          {e.code} {e.item}
        </option>
      ))}
    </select>
  );
};

export class PreApprovalForm extends Component {
  render() {
    return (
      <Form onSubmit={this.props.handleSubmit(this.props.onSubmit)}>
        <div className="usa-grid">
          <div className="usa-width-one-half">
            <SwaggerField
              fieldName="accessorial"
              title="Code & Item"
              className="rounded"
              component={Codes(this.props.tariff400ngItems)}
              swagger={this.props.ship_accessorial_schema}
              required
            />
            {/* TODO andrea - set schema location enum array to accessorial selected location value */}
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

PreApprovalForm.propTypes = {
  schema: PropTypes.object,
  tariff400ngItems: PropTypes.array,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'preapproval_request_form';

PreApprovalForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PreApprovalForm);

function mapStateToProps(state, props) {
  return {
    ship_accessorial_schema: get(state, 'swaggerPublic.spec.definitions.ShipmentAccessorial', {}),
  };
}

export default connect(mapStateToProps)(PreApprovalForm);
