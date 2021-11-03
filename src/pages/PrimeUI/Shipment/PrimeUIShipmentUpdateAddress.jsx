import React, { useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import { primeSimulatorRoutes } from '../../../constants/routes';
import { requiredAddressSchema } from '../../../utils/validation';
import scrollToTop from '../../../shared/scrollToTop';
import { updatePrimeMTOShipmentAddress } from '../../../services/primeApi';
import styles from '../../../components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { PRIME_SIMULATOR_MOVE } from '../../../constants/queryKeys';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

import { isEmpty } from 'shared/utils';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';

const updateAddressSchema = Yup.object().shape({
  addressID: Yup.string(),
  address: requiredAddressSchema,
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
          title: `${body.title} `,
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
    const body = {
      id: values.addressID,
      streetAddress1: values.address.street_address_1,
      streetAddress2: values.address.street_address_2,
      streetAddress3: values.address.street_address_3,
      city: values.address.city,
      state: values.address.state,
      postalCode: values.address.postal_code,
    };

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
    address: reformatPrimeApiPickupAddress,
    eTag: shipment.pickupAddress?.eTag,
  };
  const initialValuesDestinationAddress = {
    addressID: shipment.destinationAddress?.id,
    address: reformatPrimeApiDestinationAddress,
    eTag: shipment.destinationAddress?.eTag,
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={styles.errorContainer}>
                  <Alert type="error">
                    <span className={styles.errorTitle}>{errorMessage.title}</span>
                    <span className={styles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <h1>Update Existing Pickup & Destination Address</h1>
              {editablePickupAddress && (
                <PrimeUIShipmentUpdateAddressForm
                  initialValues={initialValuesPickupAddress}
                  onSubmit={onSubmit}
                  updateShipmentAddressSchema={updateAddressSchema}
                  addressLocation="Pickup address"
                />
              )}
              {editableDestinationAddress && (
                <PrimeUIShipmentUpdateAddressForm
                  initialValues={initialValuesDestinationAddress}
                  onSubmit={onSubmit}
                  updateShipmentAddressSchema={updateAddressSchema}
                  addressLocation="Destination address"
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
