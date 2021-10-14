import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';
import { primeSimulatorRoutes } from '../../../constants/routes';

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
  const history = useHistory();

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const onSubmit = () => {};
  const onBack = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const { mtoShipments } = moveTaskOrder;
  const shipment = mtoShipments.find((mtoShipment) => mtoShipment.id === shipmentId);

  const initialValues = {
    estimatedWeight: shipment.primeEstimatedWeight.toLocaleString(),
    actualWeight: shipment.primeActualWeight.toLocaleString(),
    requestedPickupDate: shipment.requestedPickupDate,
    actualPickupDate: shipment.actualPickupDate,
    pickupAddress: shipment.pickupAddress,
    destinationAddress: shipment.destinationAddress,
  };

  const validationSchema = Yup.object().shape({
    pickupAddress: requiredAddressSchema,
    destinationAddress: requiredAddressSchema,
    requestedPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    actualPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
  });

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
                        <DatePickerInput name="requestedPickupDate" label="Date issued" />
                        <DatePickerInput name="actualPickupDate" label="Date issued" />
                        <h2 className={styles.sectionHeader}>Shipment Weights</h2>
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
                        <h2 className={styles.sectionHeader}>Shipment Addresses</h2>
                        <h3 className={styles.sectionHeader}>Pickup Address</h3>
                        <AddressFields name="pickupAddress" />
                        <h3 className={styles.sectionHeader}>Destination Address</h3>
                        <AddressFields name="destinationAddress" />
                      </SectionWrapper>
                      <div className={formStyles.formActions}>
                        <WizardNavigation
                          editMode
                          disableNext={!isValid || isSubmitting}
                          onCancelClick={onBack}
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

/*
PrimeUIShipmentForm.propTypes = {
  initialValues: PropTypes.shape({
    estimatedWeight: PropTypes.number,
    actualWeight: PropTypes.number,
    requestedPickupDate: PropTypes.string,
    actualPickupDate: PropTypes.string,
    pickupAddress: ResidentialAddressShape,
    destinationAddress: ResidentialAddressShape,
  }),
};

 */

export default PrimeUIShipmentForm;
