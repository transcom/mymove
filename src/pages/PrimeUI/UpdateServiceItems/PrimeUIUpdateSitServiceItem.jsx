import React, { useState } from 'react';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { connect } from 'react-redux';

import PrimeUIUpdateOriginSITForm from './PrimeUIUpdateOriginSITForm';
import PrimeUIUpdateDestSITForm from './PrimeUIUpdateDestSITForm';

import { updateMTOServiceItem } from 'services/primeApi';
import { fromPrimeAPIAddressFormat } from 'utils/formatters';
import scrollToTop from 'shared/scrollToTop';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import { formatDate, formatDateForSwagger } from 'shared/dates';

const PrimeUIUpdateSitServiceItem = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const navigate = useNavigate();
  const { moveCodeOrID, mtoServiceItemId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const { mutate: createUpdateSITServiceItemRequestMutation } = useMutation(updateMTOServiceItem, {
    onSuccess: () => {
      setFlashMessage(
        `MSG_CREATE_ADDRESS_UPDATE_REQUEST_SUCCESS${moveCodeOrID}`,
        'success',
        'Successfully updated SIT service item',
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

  const serviceItem = moveTaskOrder?.mtoServiceItems.find((s) => s?.id === mtoServiceItemId);
  const { modelType } = serviceItem;
  const reformatPrimeApiSITAddress = fromPrimeAPIAddressFormat(serviceItem.sitDestinationFinalAddress);

  const destSitInitialValues = {
    address: reformatPrimeApiSITAddress,
    sitDepartureDate: formatDate(serviceItem.sitDepartureDate, 'YYYY-MM-DD', 'DD MMM YYYY'),
    sitRequestedDelivery: formatDate(serviceItem.sitRequestedDelivery, 'YYYY-MM-DD', 'DD MMM YYYY'),
    sitCustomerContacted: formatDate(serviceItem.sitCustomerContacted, 'YYYY-MM-DD', 'DD MMM YYYY'),
    mtoServiceItemID: serviceItem.id,
    eTag: serviceItem.eTag,
  };

  const originSitInitialValues = {
    sitDepartureDate: formatDate(serviceItem.sitDepartureDate, 'YYYY-MM-DD', 'DD MMM YYYY'),
    sitRequestedDelivery: formatDate(serviceItem.sitRequestedDelivery, 'YYYY-MM-DD', 'DD MMM YYYY'),
    sitCustomerContacted: formatDate(serviceItem.sitCustomerContacted, 'YYYY-MM-DD', 'DD MMM YYYY'),
    mtoServiceItemID: serviceItem.id,
    eTag: serviceItem.eTag,
  };

  // sending the data submitted in the destination SIT form to the API
  const destSitOnSubmit = (values) => {
    const { address, sitCustomerContacted, sitDepartureDate, sitRequestedDelivery, mtoServiceItemID, eTag } = values;

    const body = {
      newAddress: {
        streetAddress1: address.streetAddress1,
        streetAddress2: address.streetAddress2,
        streetAddress3: address.streetAddress3,
        city: address.city,
        state: address.state,
        postalCode: address.postalCode,
      },
      sitDepartureDate: formatDateForSwagger(sitDepartureDate),
      sitRequestedDelivery: formatDateForSwagger(sitRequestedDelivery),
      sitCustomerContacted: formatDateForSwagger(sitCustomerContacted),
      modelType: 'UpdateMTOServiceItemSIT',
    };

    createUpdateSITServiceItemRequestMutation({ mtoServiceItemID, eTag, body });
  };

  // sending the data submitted in the origin SIT form to the API
  const originSitOnSubmit = (values) => {
    const { sitCustomerContacted, sitDepartureDate, sitRequestedDelivery, mtoServiceItemID, eTag } = values;

    const body = {
      sitDepartureDate: formatDateForSwagger(sitDepartureDate),
      sitRequestedDelivery: formatDateForSwagger(sitRequestedDelivery),
      sitCustomerContacted: formatDateForSwagger(sitCustomerContacted),
      modelType: 'UpdateMTOServiceItemSIT',
    };

    createUpdateSITServiceItemRequestMutation({ mtoServiceItemID, eTag, body });
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
              {modelType === 'MTOServiceItemDestSIT' ? (
                <PrimeUIUpdateDestSITForm
                  name="address"
                  initialValues={destSitInitialValues}
                  onSubmit={destSitOnSubmit}
                />
              ) : null}
              {modelType === 'MTOServiceItemOriginSIT' ? (
                <PrimeUIUpdateOriginSITForm
                  name="address"
                  initialValues={originSitInitialValues}
                  onSubmit={originSitOnSubmit}
                />
              ) : null}
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(PrimeUIUpdateSitServiceItem);
