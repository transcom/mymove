import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { FormSection, reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { AddressElementEdit } from 'shared/Address';

import './StorageInTransit.css';

export class StorageInTransitForm extends Component {
  render() {
    const { storageInTransitSchema, addressSchema } = this.props;
    const warehouseAddress = get(this.props, 'formValues.warehouse_address');
    return (
      <form onSubmit={this.props.handleSubmit(this.props.onSubmit)} className="storage-in-transit-request-form">
        <fieldset key="sit-request-information">
          <div className="editable-panel-column">
            <SwaggerField
              fieldName="location"
              title="SIT location"
              swagger={storageInTransitSchema}
              className="storage-in-transit-location"
              required
            />
            <SwaggerField fieldName="estimated_start_date" swagger={storageInTransitSchema} required />
          </div>
          <div className="editable-panel-column">
            <SwaggerField fieldName="notes" swagger={storageInTransitSchema} />
          </div>
        </fieldset>
        <fieldset key="warehouse-information" className="storage-in-transit-hr-top">
          <h3>Warehouse</h3>
          <div className="editable-panel-column">
            <div className="panel-subhead" />
            <SwaggerField
              fieldName="warehouse_id"
              swagger={storageInTransitSchema}
              className="storage-in-transit-warehouse-id"
              required
            />
            <SwaggerField title="Warehouse name" fieldName="warehouse_name" swagger={storageInTransitSchema} required />
            <SwaggerField fieldName="warehouse_phone" swagger={storageInTransitSchema} />
            <SwaggerField fieldName="warehouse_email" swagger={storageInTransitSchema} />
          </div>
          <div className="editable-panel-column">
            <FormSection name="warehouse_address">
              <AddressElementEdit
                addressProps={{
                  swagger: addressSchema,
                  values: warehouseAddress,
                }}
                zipPattern="USA"
              />
            </FormSection>
          </div>
        </fieldset>
      </form>
    );
  }
}

StorageInTransitForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  addressSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'storage_in_transit_request_form';

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
