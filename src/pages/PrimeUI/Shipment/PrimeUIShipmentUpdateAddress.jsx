import React, { useState } from 'react';
import { Link, useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { useMutation } from 'react-query';
import * as Yup from 'yup';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import { primeSimulatorRoutes } from '../../../constants/routes';
import { requiredAddressSchema } from '../../../utils/validation';
import scrollToTop from '../../../shared/scrollToTop';
import { updatePrimeMTOShipmentAddress } from '../../../services/primeApi';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

import { isEmptyAddress, fromPrimeApiAddressFormat, toPrimeApiAddressFormat } from 'shared/utils';

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
  const [mutateMTOShipment] = useMutation(updatePrimeMTOShipmentAddress, {
    onSuccess: (updatedMTOShipmentAddress) => {
      const shipmentIndex = mtoShipments.findIndex((mtoShipment) => mtoShipment.id === shipmentId);
      Object.keys({ pickupAddress: '', destinationAddress: '' }).forEach((key) => {
        if (updatedMTOShipmentAddress.id === mtoShipments[shipmentIndex][key].id) {
          /*
          mtoShipments[shipmentIndex][key].streetAddress1 = updatedMTOShipmentAddress.streetAddress1;
          mtoShipments[shipmentIndex][key].streetAddress2 = updatedMTOShipmentAddress.streetAddress2;
          mtoShipments[shipmentIndex][key].streetAddress3 = updatedMTOShipmentAddress.streetAddress3;
          mtoShipments[shipmentIndex][key].city = updatedMTOShipmentAddress.city;
          mtoShipments[shipmentIndex][key].state = updatedMTOShipmentAddress.state;
          mtoShipments[shipmentIndex][key].postalCode = updatedMTOShipmentAddress.postalCode;
          mtoShipments[shipmentIndex].eTag = updatedMTOShipmentAddress.eTag;
           */
          mtoShipments[shipmentIndex][key] = updatedMTOShipmentAddress;
        }
      });
      handleClose();
    },
    onError: (error) => {
      /*
      const {
        response: { body },
      } = error;
       */
      console.log(error);

      if (false) {
        /*
        {
          "detail": "Invalid data found in input",
          "instance":"00000000-0000-0000-0000-000000000000",
          "title":"Validation Error",
          "invalidFields": {
            "primeEstimatedWeight":["the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight","Invalid Input."]
          }
        }
         */
        /*
        let invalidFieldsStr = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `${body.title} `,
          detail: `${body.detail}${invalidFieldsStr}\n\nPlease cancel and Update Shipment again`,
        });
         */
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment and values',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values, { setSubmitting }) => {
    console.log('Update address onSubmit clicked');
    console.log(values);

    const body = {
      id: values.addressID,
      streetAddress1: values.address.street_address_1,
      streetAddress2: values.address.street_address_2,
      streetAddress3: values.address.street_address_3,
      city: values.address.city,
      state: values.address.state,
      postalCode: values.address.postal_code,
    };
    console.log('etag - ');
    console.log(values.eTag);

    mutateMTOShipment({
      mtoShipmentID: shipmentId,
      addressID: values.addressID,
      ifMatchETag: values.eTag,
      // ifMatchETag: shipment.eTag,
      body,
    }).then(() => {
      setSubmitting(false);
    });
  };

  const reformatPrimeApiPickupAddress = fromPrimeApiAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeApiAddressFormat(shipment.destinationAddress);
  const editablePickupAddress = !isEmptyAddress(reformatPrimeApiPickupAddress);
  const editableDestinationAddress = !isEmptyAddress(reformatPrimeApiDestinationAddress);

  const initialValuesPickupAddress = {
    addressID: shipment.pickupAddress.id,
    address: reformatPrimeApiPickupAddress,
    eTag: shipment.pickupAddress.eTag,
  };
  const initialValuesDestinationAddress = {
    addressID: shipment.destinationAddress.id,
    address: reformatPrimeApiDestinationAddress,
    eTag: shipment.destinationAddress.eTag,
  };

  return (
    <>
      <h3>Update Shipment&apos;s Existing Addresses</h3>
      {true && (
        <PrimeUIShipmentUpdateAddressForm
          initialValues={initialValuesPickupAddress}
          onSubmit={onSubmit}
          updateShipmentAddressSchema={updateAddressSchema}
          addressLocation="Pickup address"
          address={shipment.pickupAddress}
        />
      )}
      {true && (
        <PrimeUIShipmentUpdateAddressForm
          initialValues={initialValuesDestinationAddress}
          onSubmit={onSubmit}
          updateShipmentAddressSchema={updateAddressSchema}
          addressLocation="Destination address"
          address={shipment.destinationAddress}
        />
      )}
    </>
  );
};
/*
  <Link className="usa-button usa-button--secondary" to={handleClose}>
    Cancel
  </Link>
 */

export default PrimeUIShipmentUpdateAddress;
