import React, { useState, useMemo } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import * as Yup from 'yup';
import { Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useMutation } from 'react-query';
import moment from 'moment';

import { createPaymentRequest } from '../../../services/primeApi';
import scrollToTop from '../../../shared/scrollToTop';
import CreatePaymentRequestForm from '../../../components/PrimeUI/CreatePaymentRequestForm/CreatePaymentRequestForm';

import styles from './CreatePaymentRequest.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

const createPaymentRequestSchema = Yup.object().shape({
  serviceItems: Yup.array().of(Yup.string()).min(1),
});

const CreatePaymentRequest = () => {
  const { moveCodeOrID } = useParams();
  const history = useHistory();

  const [errorMessage, setErrorMessage] = useState();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const [createPaymentRequestMutation] = useMutation(createPaymentRequest, {
    onSuccess: () => {
      history.push(`/simulator/moves/${moveCodeOrID}/details`);
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        setErrorMessage({ title: body.title, detail: body.detail });
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

  const { mtoShipments, mtoServiceItems } = moveTaskOrder || {};

  const groupedServiceItems = useMemo(() => {
    const serviceItems = { basic: [] };
    mtoServiceItems?.forEach((mtoServiceItem) => {
      if (mtoServiceItem.mtoShipmentID == null) {
        serviceItems.basic.push(mtoServiceItem);
      } else if (!serviceItems[mtoServiceItem.mtoShipmentID]) {
        serviceItems[mtoServiceItem.mtoShipmentID] = [mtoServiceItem];
      } else {
        serviceItems[mtoServiceItem.mtoShipmentID].push(mtoServiceItem);
      }
    });
    return serviceItems;
  }, [mtoServiceItems]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  // always display the shipments in order of creation date to not disorient the user
  mtoShipments.sort((firstShipment, secondShipment) => {
    return moment(firstShipment.createdAt) - moment(secondShipment.createdAt);
  });

  const initialValues = {
    serviceItems: [],
  };

  const onSubmit = (values, formik) => {
    const serviceItemsPayload = values.serviceItems.map((serviceItem) => {
      return { id: serviceItem };
    });
    createPaymentRequestMutation({ moveTaskOrderID: moveTaskOrder.id, serviceItems: serviceItemsPayload }).then(() => {
      formik.setSubmitting(false);
    });
  };

  const handleShipmentSelectAll = (shipmentID, values, setValues, event) => {
    const shipmentServiceItems = groupedServiceItems[shipmentID];
    const existingServiceItems = values.serviceItems;

    if (!event.target.checked) {
      // unselected the select all
      shipmentServiceItems.forEach((serviceItem) => {
        // remove the single element in place
        existingServiceItems.splice(existingServiceItems.indexOf(serviceItem.id), 1);
      });
    } else {
      shipmentServiceItems.forEach((serviceItem) => {
        // don't add duplicates if one is already selected prior to clicking select all
        if (!existingServiceItems.includes(serviceItem.id)) {
          existingServiceItems.push(serviceItem.id);
        }
      });
    }
    setValues({ serviceItems: existingServiceItems });
  };

  return (
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', styles.CreatePaymentRequest)}>
      <div className="grid-row">
        <div className="grid-col-12">
          {errorMessage?.detail && (
            <div className={styles.errorContainer}>
              <Alert slim type="error">
                <span className={styles.errorTitle}>{errorMessage.title}</span>
                <span className={styles.errorDetail}>{errorMessage.detail}</span>
              </Alert>
            </div>
          )}
          <SectionWrapper className={formStyles.formSection}>
            <dl className={descriptionListStyles.descriptionList} data-testid="moveDetails">
              <h2>Move</h2>
              <div className={descriptionListStyles.row}>
                <dt>Move Code:</dt>
                <dd>{moveTaskOrder.moveCode}</dd>
              </div>
              <div className={descriptionListStyles.row}>
                <dt>Move Id:</dt>
                <dd>{moveTaskOrder.id}</dd>
              </div>
            </dl>
          </SectionWrapper>
          <CreatePaymentRequestForm
            initialValues={initialValues}
            onSubmit={onSubmit}
            handleSelectAll={handleShipmentSelectAll}
            createPaymentRequestSchema={createPaymentRequestSchema}
            mtoShipments={mtoShipments}
            groupedServiceItems={groupedServiceItems}
          />
        </div>
      </div>
    </div>
  );
};

export default CreatePaymentRequest;
