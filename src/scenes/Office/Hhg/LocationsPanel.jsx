import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, isValid, FormSection } from 'redux-form';

import editablePanel from '../editablePanel';
import { addressElementDisplay, addressElementEdit } from '../AddressElement';
import { no_op } from 'shared/utils';

const LocationsDisplay = props => {
  const {
    shipment,
    pickupAddress,
    secondaryPickupAddress,
    deliveryAddress,
  } = props;

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">Pickup</span>
        {addressElementDisplay(pickupAddress, 'Primary')}
        {shipment.has_secondary_pickup_address &&
          addressElementDisplay(secondaryPickupAddress, 'Secondary')}
      </div>
      {shipment.has_delivery_address && (
        <div className="editable-panel-column">
          <span className="column-subhead">Delivery</span>
          {addressElementDisplay(deliveryAddress, 'Primary')}
        </div>
      )}
    </React.Fragment>
  );
};

const LocationsEdit = props => {
  const { shipment } = props;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="pickupAddress">
          <span className="column-subhead">Pickup</span>
          {addressElementEdit(shipment.pickup_address, 'Primary')}
        </FormSection>
      </div>
      <div className="editable-panel-column">
        <FormSection name="deliveryAddress">
          <span className="column-subhead">Delivery</span>
          {addressElementEdit(shipment.deliveryAddress, 'Primary')}
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
  let pickupAddress = get(shipment, 'pickup_address', {});
  let secondaryPickupAddress = get(shipment, 'secondary_pickup_address', {});
  let deliveryAddress = get(shipment, 'delivery_address', {});

  return {
    initialValues: {
      pickupAddress,
      secondaryPickupAddress,
      deliveryAddress,
    },
    // Wrapper
    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),
    addressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    hasError:
      state.office.shipmentHasLoadError || state.office.shipmentHasUpdateError,
    errorMessage: state.office.error,

    shipment: shipment,
    pickupAddress,
    secondaryPickupAddress,
    deliveryAddress,

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
