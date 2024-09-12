import React, { useState } from 'react';
import { useNavigate, useParams, generatePath, useLocation } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PrimeUIShipmentUpdateAddressForm from './PrimeUIShipmentUpdateAddressForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import { ZIP_CODE_REGEX } from 'utils/validation';
import scrollToTop from 'shared/scrollToTop';
import { updatePrimeMTOShipmentAddress } from 'services/primeApi';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { isEmpty } from 'shared/utils';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import { PRIME_SIMULATOR_MOVE } from 'constants/queryKeys';
import { getAddressLabel } from 'shared/constants';

const updateAddressSchema = Yup.object().shape({
  addressID: Yup.string(),
  address: Yup.object().shape({
    id: Yup.string(),
    streetAddress1: Yup.string().required('Required'),
    streetAddress2: Yup.string(),
    city: Yup.string().required('Required'),
    state: Yup.string().required('Required').length(2, 'Must use state abbreviation'),
    postalCode: Yup.string().required('Required').matches(ZIP_CODE_REGEX, 'Must be valid zip code'),
  }),
  eTag: Yup.string(),
});

const PrimeUIShipmentUpdateAddress = () => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const navigate = useNavigate();
  const location = useLocation();

  const addressType = location?.state?.addressType;
  const addressLabel = getAddressLabel(addressType);
  const addressData = shipment ? shipment[addressType] : null;

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const queryClient = useQueryClient();
  const { mutateAsync: mutateMTOShipment } = useMutation(updatePrimeMTOShipmentAddress, {
    onSuccess: (updatedMTOShipmentAddress) => {
      const shipmentIndex = mtoShipments.findIndex((mtoShipment) => mtoShipment.id === shipmentId);
      let updateQuery = false;
      if (updatedMTOShipmentAddress.id === mtoShipments[shipmentIndex][addressType].id) {
        mtoShipments[shipmentIndex][addressType] = updatedMTOShipmentAddress;
        updateQuery = true;
      }
      if (updateQuery) {
        moveTaskOrder.mtoShipments = mtoShipments;
        queryClient.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
        queryClient.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]);
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
    const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode } = values.address;

    const body = {
      id: values.addressID,
      streetAddress1,
      streetAddress2,
      streetAddress3,
      city,
      state,
      postalCode,
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

  const reformatPriApiAddress = fromPrimeAPIAddressFormat(addressData);
  const editableAddress = !isEmpty(reformatPriApiAddress);

  const initialValues = {
    addressID: addressData?.id,
    address: reformatPriApiAddress,
    eTag: addressData?.eTag,
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
              <h1>Update Existing {`${addressLabel}`}</h1>
              {editableAddress && (
                <PrimeUIShipmentUpdateAddressForm
                  initialValues={initialValues}
                  onSubmit={onSubmit}
                  updateShipmentAddressSchema={updateAddressSchema}
                  addressLocation={addressLabel}
                  name="address"
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
