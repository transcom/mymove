import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { addressElementDisplay, addressElementEdit } from '../AddressElement';
import { no_op } from 'shared/utils';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelField } from 'shared/EditablePanel';

const LocationsDisplay = props => {
  const pickupAddress = {
    street_address_1: '123 4th St.',
    street_address_2: 'Flat 5',
    city: 'Sixto',
    state: 'LA',
    postal_code: '89101',
  };

  const deliveryAddress = {
    street_address_1: '234 5th St.',
    street_address_2: 'Flat 6',
    city: 'Sevento',
    state: 'BE',
    postal_code: '91011',
  };

  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  // const pickupAddress = props.shipment.pickup_address;
  // const deliveryAddress = props.shipment.delivery_address;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Pickup</span>
        {addressElementDisplay(pickupAddress, 'Primary')}
      </div>
      <div className="editable-panel-column">
        <span className="column-subhead">Delivery</span>
        {addressElementDisplay(deliveryAddress, 'Primary')}
      </div>
    </React.Fragment>
  );
};

const LocationsEdit = props => {
  const pickupAddress = {
    street_address_1: '123 4th St.',
    street_address_2: 'Flat 5',
    city: 'Sixto',
    state: 'LA',
    postal_code: '89101',
  };

  const deliveryAddress = {
    street_address_1: '234 5th St.',
    street_address_2: 'Flat 6',
    city: 'Sevento',
    state: 'BE',
    postal_code: '91011',
  };
  let pickupAddressProps = {
    swagger: props.addressSchema,
    values: props.backupMailingAddress,
  };
  const { shipmentSchema } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="pickupAddress">
          <span className="column-subhead">Pickup</span>
          {addressElementEdit(pickupAddress, 'Primary')}
        </FormSection>
      </div>
      <div className="editable-panel-column">
        <FormSection name="deliveryAddress">
          <span className="column-subhead">Delivery</span>
          {addressElementEdit(deliveryAddress, 'Primary')}
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_shipment_info_locations';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let LocationsPanel = editablePanel(
  LocationsDisplay,
  LocationsEdit,
  editEnabled,
);
LocationsPanel = reduxForm({ form: formName })(LocationsPanel);

function mapStateToProps(state) {
  let shipment = get(state, 'office.officeMove.shipments.0', {});
  return {
    initialValues: {
      pickupAddress: get(shipment, 'pickup_address', {}),
      deliveryAddress: get(shipment, 'delivery_address', {}),
    },
    // Wrapper
    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),
    addressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    hasError:
      state.office.shipmentHasLoadError || state.office.shipmentHasUpdateError,
    errorMessage: state.office.error,

    shipment: shipment,
    isUpdating: state.office.shipmentIsUpdating,

    // editablePanel
    formIsValid: isValid(formName)(state),
    getUpdateArgs: function() {
      let values = getFormValues(formName)(state);
      return [shipment.id, values];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(LocationsPanel);
