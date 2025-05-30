import React, { useState, useMemo } from 'react';
import { useParams, useNavigate, generatePath } from 'react-router-dom';
import * as Yup from 'yup';
import { Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import moment from 'moment';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import { createPaymentRequest } from 'services/primeApi';
import scrollToTop from 'shared/scrollToTop';
import CreatePaymentRequestForm from 'components/PrimeUI/CreatePaymentRequestForm/CreatePaymentRequestForm';
import { primeSimulatorRoutes } from 'constants/routes';
import { PRIME_SIMULATOR_MOVE } from 'constants/queryKeys';
import { formatDateForSwagger } from 'shared/dates';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

const CreatePaymentRequest = ({ setFlashMessage }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const [errorMessage, setErrorMessage] = useState();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  const queryClient = useQueryClient();
  const { mutate: createPaymentRequestMutation } = useMutation(createPaymentRequest, {
    onSuccess: (data) => {
      if (!moveTaskOrder.paymentRequests?.length) {
        moveTaskOrder.paymentRequests = [];
      }
      moveTaskOrder.paymentRequests.push(data);

      queryClient.setQueryData([PRIME_SIMULATOR_MOVE, moveCodeOrID], moveTaskOrder);
      queryClient.invalidateQueries([PRIME_SIMULATOR_MOVE, moveCodeOrID]).then(() => {});
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
    /* Setting initial values was supposed to change formik behavior but it made no difference
    params: additionalDaySITItems.reduce(
      (acc, curr) => ({ ...acc, [curr]: { SITPaymentRequestStart: '', SITPaymentRequestEnd: '' } }),
      {},
      ),
    */
  };

  const dateValidationSchema = Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required');

  // We could ideally specify something like oneOfSchema outlined here
  // (https://gist.github.com/cb109/8eda798a4179dc21e46922a5fbb98be6) for the additional day SIT value with params
  // The behavior of Formik <FieldArray> is for dynamic lists not sparse lists as we have laid out
  const createPaymentRequestSchema = Yup.object().shape({
    serviceItems: Yup.array().of(Yup.string()).min(1),
  });

  const validateSITDate = async (id, fieldName, value, formValues, setFieldError, setFieldTouched) => {
    let error;
    // only validate service items that are being added to the payment request
    if (formValues.serviceItems?.find((serviceItem) => serviceItem === id)) {
      // The field won't be touched (and won't show the error) if the user tries to submit before editing the dates
      // even though formik says it touches all fields on submission if they are in initialValues I found this not to
      // be true.
      setFieldTouched(`params.${id}.${fieldName}`, true, false);
      await dateValidationSchema.validate(value).catch((err) => {
        error = err.message;
        // this is a way to get touched set without having to worry about other fields
        setFieldError(`params.${id}.${fieldName}`, error);
      });
    }
    return error;
  };

  const onSubmit = (values, formik) => {
    const serviceItemsPayload = values.serviceItems.map((serviceItem) => {
      // building params to send to API
      const params = [];

      // there won't always be params, so we'll check it here
      if (values.params) {
        // adding WeightBilled param to service item
        if (values.params[serviceItem]?.WeightBilled) {
          params.push({
            key: 'WeightBilled',
            value: values.params[serviceItem]?.WeightBilled,
          });
        }

        // adding add'l day SIT params to service item
        if (values.params[serviceItem]?.SITPaymentRequestStart && values.params[serviceItem]?.SITPaymentRequestEnd) {
          params.push(
            {
              key: 'SITPaymentRequestStart',
              value: formatDateForSwagger(values.params[serviceItem].SITPaymentRequestStart),
            },
            {
              key: 'SITPaymentRequestEnd',
              value: formatDateForSwagger(values.params[serviceItem].SITPaymentRequestEnd),
            },
          );
        }
      }

      if (params.length > 0) {
        return {
          id: serviceItem,
          params,
        };
      }
      return { id: serviceItem };
    });
    createPaymentRequestMutation(
      { moveTaskOrderID: moveTaskOrder.id, serviceItems: serviceItemsPayload },
      {
        onSuccess: () => {
          setFlashMessage(
            `MSG_CREATE_PAYMENT_SUCCESS${moveCodeOrID}`,
            'success',
            'Successfully created payment request',
            '',
            true,
          );
          navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
          formik.setSubmitting(false);
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
              detail:
                'An unknown error has occurred, please check the state of the shipment and service items data for this move',
            });
          }
          scrollToTop();
        },
      },
    );
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
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', primeStyles.primeContainer)}>
      <div className="grid-row">
        <div className="grid-col-12">
          {errorMessage?.detail && (
            <div className={primeStyles.errorContainer}>
              <Alert headingLevel="h4" slim type="error">
                <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
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
            handleValidateDate={validateSITDate}
            createPaymentRequestSchema={createPaymentRequestSchema}
            mtoShipments={mtoShipments}
            groupedServiceItems={groupedServiceItems}
          />
        </div>
      </div>
    </div>
  );
};

CreatePaymentRequest.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreatePaymentRequest);
