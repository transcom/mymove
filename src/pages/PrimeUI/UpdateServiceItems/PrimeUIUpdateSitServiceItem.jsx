import React, { useState } from 'react';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { connect } from 'react-redux';

import PrimeUIUpdateOriginSITForm from './PrimeUIUpdateOriginSITForm';
import PrimeUIUpdateDestSITForm from './PrimeUIUpdateDestSITForm';

import { updateMTOServiceItem } from 'services/primeApi';
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
        `UPDATE_SIT_SERVICE_ITEM_REQUEST_SUCCESS${moveCodeOrID}`,
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

  const destSitInitialValues = {
    sitDestinationFinalAddress: {
      id: serviceItem?.sitDestinationFinalAddress?.id || '',
      streetAddress1: serviceItem?.sitDestinationFinalAddress?.streetAddress1 || '',
      streetAddress2: serviceItem?.sitDestinationFinalAddress?.streetAddress2 || '',
      streetAddress3: serviceItem?.sitDestinationFinalAddress?.streetAddress3 || '',
      city: serviceItem?.sitDestinationFinalAddress?.city || '',
      state: serviceItem?.sitDestinationFinalAddress?.state || '',
      postalCode: serviceItem?.sitDestinationFinalAddress?.postalCode || '',
    },
    sitDepartureDate: formatDate(serviceItem.sitDepartureDate, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitRequestedDelivery: formatDate(serviceItem.sitRequestedDelivery, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitCustomerContacted: formatDate(serviceItem.sitCustomerContacted, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    mtoServiceItemID: serviceItem.id,
    reServiceCode: serviceItem.reServiceCode,
    eTag: serviceItem.eTag,
  };

  const originSitInitialValues = {
    sitDepartureDate: formatDate(serviceItem.sitDepartureDate, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitRequestedDelivery: formatDate(serviceItem.sitRequestedDelivery, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitCustomerContacted: formatDate(serviceItem.sitCustomerContacted, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    mtoServiceItemID: serviceItem.id,
    reServiceCode: serviceItem.reServiceCode,
    eTag: serviceItem.eTag,
  };

  // sending the data submitted in the destination SIT form to the API
  // if any of the dates are skipped or not filled with values, we'll just make them null
  const destSitOnSubmit = (values) => {
    const {
      sitDestinationFinalAddress,
      sitCustomerContacted,
      sitDepartureDate,
      sitRequestedDelivery,
      mtoServiceItemID,
      reServiceCode,
      eTag,
    } = values;

    const body = {
      sitDestinationFinalAddress: {
        id: sitDestinationFinalAddress.id,
        streetAddress1: sitDestinationFinalAddress.streetAddress1,
        streetAddress2: sitDestinationFinalAddress.streetAddress2,
        streetAddress3: sitDestinationFinalAddress.streetAddress3,
        city: sitDestinationFinalAddress.city,
        state: sitDestinationFinalAddress.state,
        postalCode: sitDestinationFinalAddress.postalCode,
      },
      sitDepartureDate: sitDepartureDate === 'Invalid date' ? null : formatDateForSwagger(sitDepartureDate),
      sitRequestedDelivery: sitRequestedDelivery === 'Invalid date' ? null : formatDateForSwagger(sitRequestedDelivery),
      sitCustomerContacted: sitCustomerContacted === 'Invalid date' ? null : formatDateForSwagger(sitCustomerContacted),
      reServiceCode,
      modelType: 'UpdateMTOServiceItemSIT',
    };

    createUpdateSITServiceItemRequestMutation({ mtoServiceItemID, eTag, body });
  };

  // sending the data submitted in the origin SIT form to the API
  // if any of the dates are skipped or not filled with values, we'll just make them null
  const originSitOnSubmit = (values) => {
    const { sitCustomerContacted, sitDepartureDate, sitRequestedDelivery, mtoServiceItemID, reServiceCode, eTag } =
      values;

    const body = {
      sitDepartureDate: sitDepartureDate === 'Invalid date' ? null : formatDateForSwagger(sitDepartureDate),
      sitRequestedDelivery: sitRequestedDelivery === 'Invalid date' ? null : formatDateForSwagger(sitRequestedDelivery),
      sitCustomerContacted: sitCustomerContacted === 'Invalid date' ? null : formatDateForSwagger(sitCustomerContacted),
      reServiceCode,
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
                  name="sitDestinationFinalAddress"
                  serviceItem={serviceItem}
                  initialValues={destSitInitialValues}
                  onSubmit={destSitOnSubmit}
                />
              ) : null}
              {modelType === 'MTOServiceItemOriginSIT' ? (
                <PrimeUIUpdateOriginSITForm
                  serviceItem={serviceItem}
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
