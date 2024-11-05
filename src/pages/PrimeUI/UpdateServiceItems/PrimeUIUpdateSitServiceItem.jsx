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
import { formatDateForSwagger, formatDateWithUTC } from 'shared/dates';
import { SERVICE_ITEM_STATUSES } from 'constants/serviceItems';

const PrimeUIUpdateSitServiceItem = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const navigate = useNavigate();
  const { moveCodeOrID, mtoServiceItemId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  /* istanbul ignore next */
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

  const initialValues = {
    sitDepartureDate: formatDateWithUTC(serviceItem.sitDepartureDate, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitRequestedDelivery: formatDateWithUTC(serviceItem.sitRequestedDelivery, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    sitCustomerContacted: formatDateWithUTC(serviceItem.sitCustomerContacted, 'YYYY-MM-DD', 'DD MMM YYYY') || '',
    mtoServiceItemID: serviceItem.id,
    reServiceCode: serviceItem.reServiceCode,
    eTag: serviceItem.eTag,
  };

  // sending the data submitted in the form to the API
  // if any of the dates are skipped or not filled with values, we'll just make them null
  const onSubmit = (values) => {
    const {
      sitCustomerContacted,
      sitDepartureDate,
      sitRequestedDelivery,
      updateReason,
      mtoServiceItemID,
      reServiceCode,
      eTag,
    } = values;

    const body = {
      sitDepartureDate: sitDepartureDate === 'Invalid date' ? null : formatDateForSwagger(sitDepartureDate),
      sitRequestedDelivery: sitRequestedDelivery === 'Invalid date' ? null : formatDateForSwagger(sitRequestedDelivery),
      sitCustomerContacted: sitCustomerContacted === 'Invalid date' ? null : formatDateForSwagger(sitCustomerContacted),
      reServiceCode,
      modelType: 'UpdateMTOServiceItemSIT',
      updateReason,
    };

    if (serviceItem?.status === SERVICE_ITEM_STATUSES.REJECTED) {
      body.requestApprovalsRequestedStatus = true;
    }

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
                  initialValues={initialValues}
                  onSubmit={onSubmit}
                />
              ) : null}
              {modelType === 'MTOServiceItemOriginSIT' ? (
                <PrimeUIUpdateOriginSITForm
                  serviceItem={serviceItem}
                  initialValues={initialValues}
                  onSubmit={onSubmit}
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
