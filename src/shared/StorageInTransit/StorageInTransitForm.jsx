import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class StorageInTransitForm extends Component {
  //form submission is still to be implemented
  handleSubmit = e => {
    e.preventDefault();
  };

  render() {
    const { storageInTransitSchema, addressSchema } = this.props;
    return (
      <form onSubmit={this.handleSubmit} className="storage-in-transit-request-form">
        <SwaggerField fieldName="sit_location" swagger={storageInTransitSchema} required />
        <SwaggerField fieldName="estimated_start_date" swagger={storageInTransitSchema} required />
        <SwaggerField fieldName="notes" swagger={storageInTransitSchema} />
        <h3>Warehouse</h3>
        <SwaggerField fieldName="warehouse_id" swagger={storageInTransitSchema} required />
        <SwaggerField fieldName="warehouse_name" swagger={storageInTransitSchema} required />
        <SwaggerField fieldName="telephone" swagger={storageInTransitSchema} />
        <SwaggerField fieldName="personal_email" swagger={storageInTransitSchema} />
        <SwaggerField fieldName="street_address_1" swagger={addressSchema} required />
        <SwaggerField fieldName="street_address_2" swagger={addressSchema} required />
        <SwaggerField fieldName="city" swagger={addressSchema} required />
        <SwaggerField fieldName="state" swagger={addressSchema} required />
        <SwaggerField fieldName="postal_code" swagger={addressSchema} required />
      </form>
    );
  }
}

StorageInTransitForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  addressSchema: PropTypes.object.isRequired,
};

const formName = 'storage_in_transit_request_form';
StorageInTransitForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(StorageInTransitForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
    addressSchema: get(state, 'swaggerPublic.spec.definitions.Address', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitForm);
