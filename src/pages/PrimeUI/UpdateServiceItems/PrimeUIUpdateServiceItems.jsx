import React, { useState } from 'react';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import PrimeUIRequestSITDestAddressChangeForm from './PrimeUIRequestSITDestAddressChangeForm';

import { createSITAddressUpdateRequest } from 'services/primeApi';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import scrollToTop from 'shared/scrollToTop';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';

const PrimeUIUpdateServiceItems = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const navigate = useNavigate();
  const { moveCodeOrID } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const destinationServiceItem = moveTaskOrder?.mtoServiceItems.find(
    (serviceItem) => serviceItem?.reServiceCode === 'DDDSIT',
  );

  const reformatPrimeApiSITDestinationAddress = fromPrimeAPIAddressFormat(
    destinationServiceItem.sitDestinationFinalAddress,
  );

  const initialValues = {
    address: reformatPrimeApiSITDestinationAddress,
    contractorRemarks: '',
    mtoServiceItemID: destinationServiceItem.id,
  };

  const { mutate: createAdressUpdateRequestMutation } = useMutation(createSITAddressUpdateRequest, {
    onSuccess: () => {
      setFlashMessage(
        `MSG_CREATE_ADDRESS_UPDATE_REQUEST_SUCCESS${moveCodeOrID}`,
        'success',
        'Successfully created SIT address update request',
        '',
        true,
      );

      navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
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

  const onSubmit = (values) => {
    const { address, contractorRemarks, mtoServiceItemID } = values;

    const body = {
      newAddress: {
        streetAddress1: address.streetAddress1,
        streetAddress2: address.streetAddress2,
        streetAddress3: address.streetAddress3,
        city: address.city,
        state: address.state,
        postalCode: address.postalCode,
      },
      contractorRemarks,
      mtoServiceItemID,
    };

    createAdressUpdateRequestMutation({ body });
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 9, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <h1 className={styles.sectionHeader}>Update Service Items</h1>
              <PrimeUIRequestSITDestAddressChangeForm
                name="address"
                initialValues={initialValues}
                onSubmit={onSubmit}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

PrimeUIUpdateServiceItems.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(PrimeUIUpdateServiceItems);
