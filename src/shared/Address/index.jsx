import React, { Component } from 'react';
import { FormSection } from 'redux-form';
import PropTypes from 'prop-types';

import { PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const AddressElementDisplay = ({ address, title }) => (
  <PanelField title={title}>
    {address.street_address_1 && (
      <span>
        {address.street_address_1}
        <br />
      </span>
    )}
    {address.street_address_2 && (
      <span>
        {address.street_address_2}
        <br />
      </span>
    )}
    {address.street_address_3 && (
      <span>
        {address.street_address_3}
        <br />
      </span>
    )}
    {address.city}, {address.state} {address.postal_code}
  </PanelField>
);

AddressElementDisplay.defaultProps = {
  address: {},
};

AddressElementDisplay.propTypes = {
  address: PropTypes.shape({
    street_address_1: PropTypes.string,
    street_address_2: PropTypes.string,
    street_address_3: PropTypes.string,
    city: PropTypes.string.isRequired,
    state: PropTypes.string.isRequired,
    postal_code: PropTypes.string.isRequired,
  }),
  title: PropTypes.string.isRequired,
};

export class AddressElementEdit extends Component {
  render() {
    return (
      <FormSection name={this.props.fieldName}>
        <div className="panel-subhead">{this.props.title}</div>
        <SwaggerField fieldName="street_address_1" swagger={this.props.schema} required />
        <SwaggerField fieldName="street_address_2" swagger={this.props.schema} />
        <SwaggerField fieldName="street_address_3" swagger={this.props.schema} />
        <SwaggerField fieldName="city" swagger={this.props.schema} required />
        <SwaggerField fieldName="state" swagger={this.props.schema} required />
        <SwaggerField fieldName="postal_code" swagger={this.props.schema} zipPattern={this.props.zipPattern} required />
      </FormSection>
    );
  }
}

AddressElementEdit.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  zipPattern: PropTypes.string,
  title: PropTypes.string,
};
