import React, { useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import { addressSchema } from 'utils/validation';
import scrollToTop from 'shared/scrollToTop';
import { updatePrimeMTOShipmentAddress } from 'services/primeApi';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { isEmpty } from 'shared/utils';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { PRIME_SIMULATOR_MOVE } from 'constants/queryKeys';

const updatePickupAddressSchema = Yup.object().shape({
  addressID: Yup.string(),
  pickupAddress: Yup.object().shape({
    address: addressSchema,
  }),
  eTag: Yup.string(),
});

const updateDestinationAddressSchema = Yup.object().shape({
  addressID: Yup.string(),
  destinationAddress: Yup.object().shape({
    address: addressSchema,
  }),
  eTag: Yup.string(),
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
      let updateQuery = false;
      ['pickupAddress', 'destinationAddress'].forEach((key) => {
        if (updatedMTOShipmentAddress.id === mtoShipments[shipmentIndex][key].id) {
          mtoShipments[shipmentIndex][key] = updatedMTOShipmentAddress;
          updateQuery = true;
        }
      });
      if (updateQuery) {
        moveTaskOrder.mtoShipments = mtoShipments;
        queryCache.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
        queryCache.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]);
      }
      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the address values used',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values, { setSubmitting }) => {
    // Choose pickupAddress or destinationAddress by the presence of the object
    // by the same name. It's possible that these values are blank and set to
    // `undefined` or an empty string `""`.
    const address = values.pickupAddress ? values.pickupAddress.address : values.destinationAddress.address;

    const body = {
      id: values.addressID,
      streetAddress1: address.streetAddress1,
      streetAddress2: address.streetAddress2,
      streetAddress3: address.streetAddress3,
      city: address.city,
      state: address.state,
      postalCode: address.postalCode,
    };

    // Check if the address payload contains any blank properties and remove
    // them. This will allow the backend to send the proper error messages
    // since the properties won't exist in the payload that is sent.
    Object.keys(body).forEach((k) => {
      if (!body[k]) {
        delete body[k];
      }
    });

    mutateMTOShipment({
      mtoShipmentID: shipmentId,
      addressID: values.addressID,
      ifMatchETag: values.eTag,
      body,
    }).then(() => {
      setSubmitting(false);
    });
  };

  const reformatPrimeApiPickupAddress = fromPrimeAPIAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeAPIAddressFormat(shipment.destinationAddress);
  const editablePickupAddress = !isEmpty(reformatPrimeApiPickupAddress);
  const editableDestinationAddress = !isEmpty(reformatPrimeApiDestinationAddress);

  const initialValuesPickupAddress = {
    addressID: shipment.pickupAddress?.id,
    pickupAddress: {
      address: reformatPrimeApiPickupAddress,
    },
    eTag: shipment.pickupAddress?.eTag,
  };
  const initialValuesDestinationAddress = {
    addressID: shipment.destinationAddress?.id,
    destinationAddress: {
      address: reformatPrimeApiDestinationAddress,
    },
    eTag: shipment.destinationAddress?.eTag,
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <h1>Update Existing Pickup & Destination Address</h1>
              {editablePickupAddress && (
                <PrimeUIShipmentUpdateAddressForm
                  initialValues={initialValuesPickupAddress}
                  onSubmit={onSubmit}
                  updateShipmentAddressSchema={updatePickupAddressSchema}
                  addressLocation="Pickup address"
                  name="pickupAddress.address"
                />
              )}
              {editableDestinationAddress && (
                <PrimeUIShipmentUpdateAddressForm
                  initialValues={initialValuesDestinationAddress}
                  onSubmit={onSubmit}
                  updateShipmentAddressSchema={updateDestinationAddressSchema}
                  addressLocation="Destination address"
                  name="destinationAddress.address"
                />
              )}
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default PrimeUIShipmentUpdateAddress;
