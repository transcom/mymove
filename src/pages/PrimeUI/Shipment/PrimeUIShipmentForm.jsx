import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import { primeSimulatorRoutes } from '../../../constants/routes';
import { formatDate, formatWeight, formatSwaggerDate } from '../../../shared/formatters';

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
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const reformatPrimeApiPickupAddress = {
    street_address_1: shipment.pickupAddress.streetAddress1,
    street_address_2: shipment.pickupAddress.streetAddress2,
    city: shipment.pickupAddress.city,
    state: shipment.pickupAddress.state,
    postal_code: shipment.pickupAddress.postalCode,
  };

  const reformatPrimeApiDestinationAddress = {
    street_address_1: shipment.destinationAddress.streetAddress1,
    street_address_2: shipment.destinationAddress.streetAddress2,
    city: shipment.destinationAddress.city,
    state: shipment.destinationAddress.state,
    postal_code: shipment.destinationAddress.postalCode,
  };

  const onSubmit = (values) => {
    /* TODO: requestedPickupDate, make display only not available in update shipment API */
    const { estimatedWeight, actualWeight, actualPickupDate, pickupAddress, destinationAddress } = values;

    /* TODO this works, but I'd like to null out the address if there is no update, but maybe I shouldn't, right now this works and I think because the address is the same, the API/server is smart enough not to update it. However, I'm not sure if I make the address null if it would work... that might actually cause an issue. */
    console.log('building the body for the call');
    const body = {
      primeEstimatedWeight: estimatedWeight?.toInteger,
      primeActualWeight: actualWeight?.toInteger,
      actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
      pickupAddress: reformatPrimeApiPickupAddress === pickupAddress ? null : pickupAddress,
      destinationAddress: reformatPrimeApiDestinationAddress === destinationAddress ? null : destinationAddress,
    };
    /*
        const body = {
      primeEstimatedWeight: estimatedWeight?.toInteger,
      primeActualWeight: actualWeight?.toInteger,
      actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
      pickupAddress: shipment.pickupAddress === pickupAddress ? null : pickupAddress,
      destinationAddress: shipment.destinationAddress === destinationAddress ? null : destinationAddress,
    };
     */
    console.log(body);
    console.log('calling mutateMTOShipment');
    mutateMTOShipment({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag, body });
  };

  const initialValues = {
    estimatedWeight: shipment.primeEstimatedWeight?.toLocaleString(),
    actualWeight: shipment.primeActualWeight?.toLocaleString(),
    requestedPickupDate: shipment.requestedPickupDate,
    actualPickupDate: shipment.actualPickupDate,
    pickupAddress: reformatPrimeApiPickupAddress,
    destinationAddress: reformatPrimeApiDestinationAddress,
  };

  const validationSchema = Yup.object().shape({
    pickupAddress: requiredAddressSchema,
    destinationAddress: requiredAddressSchema,
    requestedPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    actualPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
  });

  const editableWeightEstimateField = shipment.primeEstimatedWeight === 0;
  const editableWeightActualField = shipment.primeActualWeight === 0;

  const emptyAddressShape = {
    street_address_1: '',
    street_address_2: '',
    city: '',
    state: '',
    postal_code: '',
  };

  const editablePickupAddress = shipment.pickupAddress === emptyAddressShape;
  const editableDestinationAddress = shipment.destinationAddress === emptyAddressShape;

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
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
