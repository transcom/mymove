import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import './StorageInTransitForm.css';

export class StorageInTransitForm extends Component {
  //form submission is still to be implemented
  handleSubmit = e => {
    e.preventDefault();
  };

  render() {
    const { storageInTransitSchema, addressSchema } = this.props;
    return (
      <form onSubmit={this.handleSubmit} className="storage-in-transit-request-form">
        <fieldset key="sit-request-information">
          <h3>Request SIT</h3>
          <div className="storage-in-transit-location">
            <SwaggerField fieldName="sit_location" swagger={storageInTransitSchema} required />
            <SwaggerField fieldName="estimated_start_date" swagger={storageInTransitSchema} required />
          </div>
          <div className="storage-in-transit-notes">
            <SwaggerField fieldName="notes" className="three-quarter-width" swagger={storageInTransitSchema} />
          </div>
        </fieldset>
        <fieldset key="warehouse-information" className="storage-in-transit-hr-top">
          <h3>Warehouse</h3>
          <div className="storage-in-transit-warehouse-information">
            <SwaggerField
              className="storage-in-transit-warehouse-id"
              fieldName="warehouse_id"
              swagger={storageInTransitSchema}
              required
            />
            <SwaggerField fieldName="warehouse_name" swagger={storageInTransitSchema} required />
            <SwaggerField fieldName="telephone" swagger={storageInTransitSchema} />
            <SwaggerField fieldName="personal_email" swagger={storageInTransitSchema} />
          </div>
          <div className="storage-in-transit-warehouse-address">
            <SwaggerField fieldName="street_address_1" swagger={addressSchema} required />
            <SwaggerField fieldName="street_address_2" swagger={addressSchema} required />
            <SwaggerField fieldName="city" swagger={addressSchema} required />
            <SwaggerField fieldName="state" swagger={addressSchema} required />
            <SwaggerField fieldName="postal_code" swagger={addressSchema} required />
          </div>
        </fieldset>
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
