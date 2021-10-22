import React, { useState } from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import * as Yup from 'yup';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import { primeSimulatorRoutes } from '../../../constants/routes';
import { requiredAddressSchema } from '../../../utils/validation';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

const emptyAddressShape = {
  street_address_1: '',
  street_address_2: '',
  city: '',
  state: '',
  postal_code: '',
};

const updateAddressSchema = Yup.object().shape({
  address: requiredAddressSchema,
});

const PrimeUIShipmentUpdateAddress = () => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const initialValues = {
    serviceItems: [],
  };

  const onSubmit = (values, { setSubmitting }) => {
    console.log('Update address onSubmit clicked');
    /*
    const serviceItemsPayload = values.serviceItems.map((serviceItem) => {
      return { id: serviceItem };
    });
    createPaymentRequestMutation({ moveTaskOrderID: moveTaskOrder.id, serviceItems: serviceItemsPayload }).then(() => {
      setSubmitting(false);
    });
     */
  };

  return (
    <>
      <h3>Update Shipment&apos;s Existing Addresses</h3>
      {editablePickupAddress && (
        <PrimeUIShipmentUpdateAddressForm
          initialValues={initialValues}
          onSubmit={onSubmit}
          updateShipmentAddressSchema={updateAddressSchema}
          addressLocation="Pickup address"
          address={shipment.pickupAddress}
        />
      )}
      {editableDestinationAddress && (
        <PrimeUIShipmentUpdateAddressForm
          initialValues={initialValues}
          onSubmit={onSubmit}
          updateShipmentAddressSchema={updateAddressSchema}
          addressLocation="Destination address"
          address={shipment.destinationAddress}
        />
      )}
      <Link className="usa-button usa-button--secondary" to={handleClose()}>
        Cancel
      </Link>
    </>
  );
};

export default PrimeUIShipmentUpdateAddress;
