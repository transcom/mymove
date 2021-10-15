import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import { primeSimulatorRoutes } from '../../../constants/routes';
import { formatDate, formatWeight, formatSwaggerDate } from '../../../shared/formatters';
import scrollToTop from '../../../shared/scrollToTop';

import { formatAddress } from 'utils/shipmentDisplay';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { updatePrimeMTOShipment } from 'services/primeApi';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { requiredAddressSchema } from 'utils/validation';

const PrimeUIShipmentForm = () => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };
  const [mutateMTOShipment] = useMutation(updatePrimeMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      queryCache.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
      handleClose();
    },
    onError: (error) => {
      const {
        response: { body },
      } = error;

      if (body) {
        /*
        {
          "detail": "Invalid data found in input",
          "instance":"00000000-0000-0000-0000-000000000000",
          "title":"Validation Error",
          "invalidFields": {
            "primeEstimatedWeight":["the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight","Invalid Input."]
          }
        }
         */
        let invalidFieldsStr = '';
        if (Object.prototype.hasOwnProperty.call(body, 'invalidFields')) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `${body.title} `,
          detail: `${body.detail}${invalidFieldsStr}\n\nPlease cancel and Update Shipment again`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment and values',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const fromPrimeApiAddressFormat = (address) => {
    return {
      street_address_1: address.streetAddress1,
      street_address_2: address.streetAddress2,
      street_address_3: address.streetAddress3,
      city: address.city,
      state: address.state,
      postal_code: address.postalCode,
    };
  };
  const toPrimeApiAddressFormat = (address) => {
    return {
      streetAddress1: address.street_address_1,
      streetAddress2: address.street_address_2,
      streetAddress3: address.street_address_3,
      city: address.city,
      state: address.state,
      postalCode: address.postal_code,
    };
  };

  const emptyAddress = {
    street_address_1: '',
    street_address_2: '',
    street_address_3: '',
    city: '',
    state: '',
    postal_code: '',
  };

  const isValidWeight = (weight) => {
    if (weight !== 'undefined' && weight && weight > 0) {
      return true;
    }
    return false;
  };

  const editableWeightEstimateField = !isValidWeight(shipment.primeEstimatedWeight);
  const editableWeightActualField = true; // !isValidWeight(shipment.primeActualWeight);

  // Not the Prime API address format
  const isEmptyAddress = (address) => {
    if (address.street_address_1 !== 'undefined' && address.street_address_1) {
      return false;
    }
    if (address.street_address_2 !== 'undefined' && address.street_address_2) {
      return false;
    }
    if (address.street_address_3 !== 'undefined' && address.street_address_3) {
      return false;
    }
    if (address.city !== 'undefined' && address.city) {
      return false;
    }
    if (address.state !== 'undefined' && address.state) {
      return false;
    }
    if (address.postal_code !== 'undefined' && address.postal_code) {
      return false;
    }
    return true;
  };

  const reformatPrimeApiPickupAddress = fromPrimeApiAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeApiAddressFormat(shipment.destinationAddress);
  const editablePickupAddress = isEmptyAddress(reformatPrimeApiPickupAddress);
  const editableDestinationAddress = isEmptyAddress(reformatPrimeApiDestinationAddress);

  const onSubmit = (values) => {
    const { estimatedWeight, actualWeight, actualPickupDate, scheduledPickupDate, pickupAddress, destinationAddress } =
      values;

    const body = {
      primeEstimatedWeight: editableWeightEstimateField ? parseInt(estimatedWeight, 10) : null,
      primeActualWeight: parseInt(actualWeight, 10),
      scheduledPickupDate: scheduledPickupDate ? formatSwaggerDate(scheduledPickupDate) : null,
      actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
      pickupAddress: editablePickupAddress ? toPrimeApiAddressFormat(pickupAddress) : null,
      destinationAddress: editableDestinationAddress ? toPrimeApiAddressFormat(destinationAddress) : null,
    };
    mutateMTOShipment({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag, body });
  };

  const initialValues = {
    estimatedWeight: shipment.primeEstimatedWeight?.toLocaleString(),
    actualWeight: shipment.primeActualWeight?.toLocaleString(),
    requestedPickupDate: shipment.requestedPickupDate,
    scheduledPickupDate: shipment.scheduledPickupDate,
    actualPickupDate: shipment.actualPickupDate,
    pickupAddress: editablePickupAddress ? emptyAddress : reformatPrimeApiPickupAddress,
    destinationAddress: editableDestinationAddress ? emptyAddress : reformatPrimeApiDestinationAddress,
  };

  const validationSchema = Yup.object().shape({
    pickupAddress: requiredAddressSchema,
    destinationAddress: requiredAddressSchema,
    requestedPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    scheduledPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    actualPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
  });

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
              <Formik
                initialValues={initialValues}
                onSubmit={onSubmit}
                validationSchema={validationSchema}
                validateOnMount
              >
                {({ isValid, isSubmitting, handleSubmit }) => {
                  return (
                    <Form className={formStyles.form}>
                      <SectionWrapper className={`${formStyles.formSection} ${styles.formSectionHeader}`}>
                        <h2 className={styles.sectionHeader}>Shipment Dates</h2>
                        <h5 className={styles.sectionHeader}>Requested Pickup</h5>

                        <>{formatDate(shipment.requestedPickupDate)}</>
                        <DatePickerInput name="scheduledPickupDate" label="Scheduled pickup" />
                        <DatePickerInput name="actualPickupDate" label="Actual pickup" />
                        <h2 className={styles.sectionHeader}>Shipment Weights</h2>
                        {editableWeightEstimateField && (
                          <MaskedTextField
                            data-testid="estimatedWeightInput"
                            defaultValue="0"
                            name="estimatedWeight"
                            label="Estimated weight (lbs)"
                            id="estimatedWeightInput"
                            mask={Number}
                            scale={0} // digits after point, 0 for integers
                            signed={false} // disallow negative
                            thousandsSeparator=","
                            lazy={false} // immediate masking evaluation
                          />
                        )}
                        {!editableWeightEstimateField && (
                          <>
                            <dt>
                              <h5 className={styles.sectionHeader}>Estimated Weight</h5>
                            </dt>
                            <dd data-testid="authorizedWeight">{formatWeight(shipment.primeEstimatedWeight)}</dd>
                          </>
                        )}
                        {editableWeightActualField && (
                          <MaskedTextField
                            data-testid="actualWeightInput"
                            defaultValue="0"
                            name="actualWeight"
                            label="Actual weight (lbs)"
                            id="actualWeightInput"
                            mask={Number}
                            scale={0} // digits after point, 0 for integers
                            signed={false} // disallow negative
                            thousandsSeparator=","
                            lazy={false} // immediate masking evaluation
                          />
                        )}
                        {!editableWeightActualField && (
                          <>
                            <dt>
                              <h5 className={styles.sectionHeader}>Actual Weight</h5>
                            </dt>
                            <dd data-testid="authorizedWeight">{formatWeight(initialValues.actualWeight)}</dd>
                          </>
                        )}
                        <h2 className={styles.sectionHeader}>Shipment Addresses</h2>
                        <h5 className={styles.sectionHeader}>Pickup Address</h5>
                        {editablePickupAddress && <AddressFields name="pickupAddress" />}
                        {!editablePickupAddress && formatAddress(initialValues.pickupAddress)}
                        <h5 className={styles.sectionHeader}>Destination Address</h5>
                        {editableDestinationAddress && <AddressFields name="destinationAddress" />}
                        {!editableDestinationAddress && formatAddress(initialValues.destinationAddress)}
                      </SectionWrapper>
                      <div className={formStyles.formActions}>
                        <WizardNavigation
                          editMode
                          disableNext={!isValid || isSubmitting}
                          onCancelClick={handleClose}
                          onNextClick={handleSubmit}
                        />
                      </div>
                    </Form>
                  );
                }}
              </Formik>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default PrimeUIShipmentForm;
