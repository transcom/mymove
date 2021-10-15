import React, { useState } from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, FormGroup, Checkbox, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useMutation } from 'react-query';

import { formatDateFromIso } from '../../../shared/formatters';
import { shipmentTypeLabels } from '../../../content/shipments';
import { ErrorMessage } from '../../../components/form';
import { createPaymentRequest } from '../../../services/primeApi';
import scrollToTop from '../../../shared/scrollToTop';
import Hint from '../../../components/Hint';

import styles from './CreatePaymentRequest.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { AgentShape } from 'types/agent';
import { AddressShape } from 'types/address';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

const ServiceItem = ({ serviceItem }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{`${serviceItem.reServiceName}`}</h3>
      <div className={descriptionListStyles.row}>
        <dt>Status:</dt>
        <dd>{serviceItem.status}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>ID:</dt>
        <dd>{serviceItem.id}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Service Code:</dt>
        <dd>{serviceItem.reServiceCode}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Service Name:</dt>
        <dd>{serviceItem.reServiceName}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>eTag:</dt>
        <dd>{serviceItem.eTag}</dd>
      </div>
    </dl>
  );
};

ServiceItem.propTypes = {
  serviceItem: PropTypes.shape({
    id: PropTypes.string,
    reServiceCode: PropTypes.string,
    reServiceName: PropTypes.string,
    eTag: PropTypes.string,
    status: PropTypes.string,
  }).isRequired,
};

const Shipment = ({ shipment }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
      <div className={descriptionListStyles.row}>
        <dt>Status:</dt>
        <dd>{shipment.status}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment ID:</dt>
        <dd>{shipment.id}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment eTag:</dt>
        <dd>{shipment.eTag}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Requested Pickup Date:</dt>
        <dd>{shipment.requestedPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Pickup Date:</dt>
        <dd>{shipment.actualPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Estimated Weight:</dt>
        <dd>{shipment.primeEstimatedWeight}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Weight:</dt>
        <dd>{shipment.primeActualWeight}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Pickup Address:</dt>
        <dd>
          {shipment.pickupAddress.streetAddress1} {shipment.pickupAddress.streetAddress2} {shipment.pickupAddress.city}{' '}
          {shipment.pickupAddress.state} {shipment.pickupAddress.postalCode}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>
          {shipment.destinationAddress.streetAddress1} {shipment.destinationAddress.streetAddress2}{' '}
          {shipment.destinationAddress.city} {shipment.destinationAddress.state}{' '}
          {shipment.destinationAddress.postalCode}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Created at:</dt>
        <dd>{formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD')}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Approved at:</dt>
        <dd>{shipment.approvedDate}</dd>
      </div>
    </dl>
  );
};

Shipment.propTypes = {
  shipment: PropTypes.shape({
    id: PropTypes.string,
    eTag: PropTypes.string,
    shipmentType: ShipmentOptionsOneOf,
    requestedPickupDate: PropTypes.string,
    scheduledPickupDate: PropTypes.string,
    actualPickupDate: PropTypes.string,
    pickupAddress: AddressShape,
    secondaryPickupAddress: AddressShape,
    destinationAddress: AddressShape,
    secondaryDeliveryAddress: AddressShape,
    agents: PropTypes.arrayOf(AgentShape),
    primeEstimatedWeight: PropTypes.number,
    primeActualWeight: PropTypes.number,
    diversion: PropTypes.bool,
    counselorRemarks: PropTypes.string,
    customerRemarks: PropTypes.string,
    status: PropTypes.string,
    reweigh: PropTypes.shape({
      id: PropTypes.string,
    }),
    createdAt: PropTypes.string,
    approvedDate: PropTypes.string,
  }).isRequired,
};

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
      const {
        response: { body },
      } = error;

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

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { mtoShipments, mtoServiceItems } = moveTaskOrder;
  const MoveServiceCodes = ['MS', 'CS'];

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
            <dl className={descriptionListStyles.descriptionList}>
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
          <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={createPaymentRequestSchema}>
            {({ isValid, dirty, isSubmitting, errors }) => (
              <Form className={formStyles.form}>
                <FormGroup error={errors != null && Object.keys(errors).length > 0}>
                  {errors != null && errors.serviceItems && (
                    <ErrorMessage display>
                      At least 1 service item must be added when creating a payment request
                    </ErrorMessage>
                  )}
                  <SectionWrapper className={formStyles.formSection}>
                    <dl className={descriptionListStyles.descriptionList}>
                      <h2>Move Service Items</h2>
                      {mtoServiceItems?.map((mtoServiceItem, mtoServiceItemIndex) => {
                        return (
                          MoveServiceCodes.includes(mtoServiceItem.reServiceCode) && (
                            <SectionWrapper
                              key={`moveServiceItems${mtoServiceItem.id}`}
                              className={formStyles.formSection}
                            >
                              <div className={styles.serviceItemInputGroup}>
                                <Label htmlFor={mtoServiceItem.id}>Add to payment request</Label>
                                <Field
                                  as={Checkbox}
                                  type="checkbox"
                                  name="serviceItems"
                                  value={mtoServiceItem.id}
                                  id={mtoServiceItem.id}
                                />
                              </div>
                              <ServiceItem
                                serviceItem={mtoServiceItem}
                                shipmentServiceItemNumber={mtoServiceItemIndex}
                              />
                            </SectionWrapper>
                          )
                        );
                      })}
                    </dl>
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <dl className={descriptionListStyles.descriptionList}>
                      <h2>Shipments</h2>
                      {mtoShipments?.map((mtoShipment) => {
                        return (
                          <div key={mtoShipment.id}>
                            <Shipment shipment={mtoShipment} />
                            <h2>Shipment Service Items</h2>
                            {mtoServiceItems?.map((mtoServiceItem, mtoServiceItemIndex) => {
                              return (
                                mtoServiceItem.mtoShipmentID === mtoShipment.id && (
                                  <SectionWrapper
                                    key={`shipmentServiceItems${mtoServiceItem.id}`}
                                    className={formStyles.formSection}
                                  >
                                    <div className={styles.serviceItemInputGroup}>
                                      <Label htmlFor={mtoServiceItem.id}>Add to payment request</Label>
                                      <Field
                                        as={Checkbox}
                                        type="checkbox"
                                        name="serviceItems"
                                        value={mtoServiceItem.id}
                                        id={mtoServiceItem.id}
                                      />
                                    </div>
                                    <ServiceItem
                                      serviceItem={mtoServiceItem}
                                      shipmentServiceItemNumber={mtoServiceItemIndex}
                                    />
                                  </SectionWrapper>
                                )
                              );
                            })}
                          </div>
                        );
                      })}
                    </dl>
                    <Button
                      aria-label="Submit Payment Request"
                      type="submit"
                      disabled={!dirty || isSubmitting || !isValid}
                    >
                      Submit Payment Request
                    </Button>
                    <Hint>
                      At least one basic service item or shipment service item is required to create a payment request
                    </Hint>
                  </SectionWrapper>
                </FormGroup>
              </Form>
            )}
          </Formik>
        </div>
      </div>
    </div>
  );
};

export default CreatePaymentRequest;
