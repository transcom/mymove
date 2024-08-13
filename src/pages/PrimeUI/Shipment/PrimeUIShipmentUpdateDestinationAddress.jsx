import React, { useState } from 'react';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PrimeUIShipmentUpdateDestinationAddressForm from './PrimeUIShipmentUpdateDestinationAddressForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import { addressSchema } from 'utils/validation';
import scrollToTop from 'shared/scrollToTop';
import { updateShipmentDestinationAddressV3 } from 'services/primeApi';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { isEmpty } from 'shared/utils';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { setFlashMessage } from 'store/flash/actions';

const updateDestinationAddressSchema = Yup.object().shape({
  mtoShipmentID: Yup.string(),
  newAddress: Yup.object().shape({
    address: addressSchema,
  }),
  newSecondaryAddress: Yup.object().shape({
    address: addressSchema,
  }),
  newTertiaryAddress: Yup.object().shape({
    address: addressSchema,
  }),
  contractorRemarks: Yup.string().required('Contractor remarks are required to make these changes'),
  eTag: Yup.string(),
});

const PrimeUIShipmentUpdateDestinationAddress = () => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  /* istanbul ignore next */
  const { mutate: updateShipmentDestinationAddressAPI } = useMutation(updateShipmentDestinationAddressV3, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      setFlashMessage(`MSG_UPDATE_SUCCESS${shipmentId}`, 'success', `Successfully updated shipment`, '', true);
      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        let additionalDetails = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            additionalDetails += `:\n${key} - ${body.invalidFields[key]}`;
          });
        }

        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${additionalDetails}`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail:
            'An unknown error has occurred, please check the state of the shipment and service items data for this move',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = (values, { setSubmitting }) => {
    const { mtoShipmentID, newAddress, newSecondaryAddress, newTertiaryAddress } = values;

    const body = {
      newAddress: {
        id: newAddress.address.id,
        streetAddress1: newAddress.address.streetAddress1,
        streetAddress2: newAddress.address.streetAddress2,
        streetAddress3: newAddress.address.streetAddress3,
        city: newAddress.address.city,
        state: newAddress.address.state,
        postalCode: newAddress.address.postalCode,
      },
      newSecondaryAddress: {
        id: newSecondaryAddress.address.id,
        streetAddress1: newSecondaryAddress.address.streetAddress1,
        streetAddress2: newSecondaryAddress.address.streetAddress2,
        streetAddress3: newSecondaryAddress.address.streetAddress3,
        city: newSecondaryAddress.address.city,
        state: newSecondaryAddress.address.state,
        postalCode: newSecondaryAddress.address.postalCode,
      },
      newTertiaryAddress: {
        id: newTertiaryAddress.address.id,
        streetAddress1: newTertiaryAddress.address.streetAddress1,
        streetAddress2: newTertiaryAddress.address.streetAddress2,
        streetAddress3: newTertiaryAddress.address.streetAddress3,
        city: newTertiaryAddress.address.city,
        state: newTertiaryAddress.address.state,
        postalCode: newTertiaryAddress.address.postalCode,
      },
      contractorRemarks: values.contractorRemarks,
    };

    updateShipmentDestinationAddressAPI({
      mtoShipmentID,
      ifMatchETag: values.eTag,
      body,
    }).then(() => {
      setSubmitting(false);
    });
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <PrimeUIShipmentUpdateDestinationAddressForm
                onSubmit={onSubmit}
                updateDestinationAddressSchema={updateDestinationAddressSchema}
                shipment={shipment}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default PrimeUIShipmentUpdateDestinationAddress;
