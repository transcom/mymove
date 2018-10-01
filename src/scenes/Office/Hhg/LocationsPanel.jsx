import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues, FormSection } from 'redux-form';

import { editablePanelify } from 'shared/EditablePanel';
import { AddressElementDisplay, AddressElementEdit } from 'shared/Address';
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
        <AddressElementDisplay address={pickupAddress} title="Primary" />
        {shipment.has_secondary_pickup_address && (
          <AddressElementDisplay
            address={secondaryPickupAddress}
            title="Secondary"
          />
        )}
      </div>
      {shipment.has_delivery_address && (
        <div className="editable-panel-column">
          <span className="column-subhead">Delivery</span>
          <AddressElementDisplay address={deliveryAddress} title="Primary" />
        </div>
      )}
    </React.Fragment>
  );
};

const LocationsEdit = props => {
  let { shipment } = props;
  let pickupProps = {
    swagger: props.addressSchema,
    values: props.pickupAddress,
  };
  let secondaryPickupProps = {
    swagger: props.addressSchema,
    values: props.secondaryPickupAddress,
  };
  let deliveryProps = {
    swagger: props.addressSchema,
    values: props.deliveryAddress,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="pickupAddress">
          <span className="column-subhead">Pickup</span>
          <AddressElementEdit addressProps={pickupProps} title="Primary" />
          {shipment.has_secondary_pickup_address && (
            <AddressElementEdit
              addressProps={secondaryPickupProps}
              title="Secondary"
            />
          )}
        </FormSection>
      </div>
      <div className="editable-panel-column">
        <FormSection name="deliveryAddress">
          <span className="column-subhead">Delivery</span>
          {shipment.has_delivery_address && (
            <AddressElementEdit addressProps={deliveryProps} title="Primary" />
          )}
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_shipment_info_locations';
const editEnabled = false; // to remove the "Edit" button on panel header and disable editing

let LocationsPanel = editablePanelify(
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

    shipment,
    pickupAddress,
    secondaryPickupAddress,
    deliveryAddress,

    isUpdating: state.office.shipmentIsUpdating,

    // editablePanel
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
