import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm, Field } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { AddressElementEdit } from 'shared/Address';
import RadioButton from 'shared/RadioButton';
import validator from '../JsonSchemaForm/validator';

import './StorageInTransit.css';

const RadioGroup = ({ location, change, ...input }) => {
  location = location === 'ORIGIN' ? 'ORIGIN' : 'DESTINATION';
  return (
    <div className="radio-group-wrapper normalize-margins">
      <RadioButton
        inputClassName="inline_radio"
        labelClassName="radio-label__location"
        label="Origin"
        value="origin"
        name="location"
        checked={location === 'ORIGIN'}
        onChange={() => change('location', 'ORIGIN')}
        testId="origin-radio"
      />
      <RadioButton
        inputClassName="inline_radio"
        labelClassName="radio-label__location"
        label="Destination"
        value="destination"
        name="location"
        checked={location === 'DESTINATION'}
        onChange={() => change('location', 'DESTINATION')}
        testId="destination-radio"
      />
    </div>
  );
};

export class StorageInTransitForm extends Component {
  render() {
    const { storageInTransitSchema, addressSchema, location, change } = this.props;
    return (
      <form onSubmit={this.props.handleSubmit(this.props.onSubmit)} className="storage-in-transit-form">
        <fieldset key="sit-request-information">
          <div className="editable-panel-column">
            <div className="radio-group-wrapper normalize-margins">
              <p className="radio-group-header">SIT Location</p>
              <Field
                component={RadioGroup}
                name="location"
                location={location}
                change={change}
                validate={[validator.isRequired]}
              />
            </div>
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
            <AddressElementEdit fieldName="warehouse_address" schema={addressSchema} zipPattern="USA" />
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
    location: get(state, 'form.storage_in_transit_request_form.values.location'),
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
    addressSchema: get(state, 'swaggerPublic.spec.definitions.Address', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitForm);
