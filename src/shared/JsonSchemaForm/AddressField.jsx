import { FormSection } from 'redux-form';
import React, { Component } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import PropTypes from 'prop-types';

export class AddressField extends Component {
  render() {
    return (
      //Todo: Add css formatting to the fields.
      <FormSection name={this.props.fieldName}>
        <div className="address-segment usa-grid">
          <SwaggerField fieldName="street_address_1" swagger={this.props.schema} required />
          <SwaggerField fieldName="street_address_2" swagger={this.props.schema} />
          <SwaggerField fieldName="city" swagger={this.props.schema} required />
          <SwaggerField fieldName="state" swagger={this.props.schema} required />
          <SwaggerField fieldName="postal_code" swagger={this.props.schema} required />
        </div>
      </FormSection>
    );
  }
}

AddressField.propTypes = {
  schema: PropTypes.object.isRequired,
  fieldName: PropTypes.string.isRequired,
};
