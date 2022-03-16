import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useHistory, useParams, withRouter } from 'react-router-dom';
import { generatePath } from 'react-router';
import { useMutation } from 'react-query';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import { formatSwaggerDate } from 'shared/formatters';
import scrollToTop from 'shared/scrollToTop';
import { updatePrimeMTOShipment } from 'services/primeApi';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { addressSchema } from 'utils/validation';
import { isValidWeight, isEmpty } from 'shared/utils';
import { fromPrimeAPIAddressFormat, formatAddressForPrimeAPI } from 'utils/formatters';
import PrimeUIShipmentUpdateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdateForm';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

const PrimeUIShipmentUpdate = ({ setFlashMessage }) => {
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
      setFlashMessage(`MSG_CREATE_PAYMENT_SUCCESS${shipmentId}`, 'success', `Successfully updated shipment`, '', true);
      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

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
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `Prime API: ${body.title} `,
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

  const emptyAddress = {
    streetAddress1: '',
    streetAddress2: '',
    streetAddress3: '',
    city: '',
    state: '',
    postalCode: '',
  };

  const editableWeightEstimateField = !isValidWeight(shipment.primeEstimatedWeight);
  const editableWeightActualField = true;
  const reformatPrimeApiPickupAddress = fromPrimeAPIAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeAPIAddressFormat(shipment.destinationAddress);
  const editablePickupAddress = isEmpty(reformatPrimeApiPickupAddress);
  const editableDestinationAddress = isEmpty(reformatPrimeApiDestinationAddress);

  const onSubmit = (values, { setSubmitting }) => {
    const {
      estimatedWeight,
      actualWeight,
      actualPickupDate,
      scheduledPickupDate,
      pickupAddress,
      destinationAddress,
      destinationType,
      diversion,
    } = values;

    const body = {
      primeEstimatedWeight: editableWeightEstimateField ? parseInt(estimatedWeight, 10) : null,
      primeActualWeight: parseInt(actualWeight, 10),
      scheduledPickupDate: scheduledPickupDate ? formatSwaggerDate(scheduledPickupDate) : null,
      actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
      pickupAddress: editablePickupAddress ? formatAddressForPrimeAPI(pickupAddress) : null,
      destinationAddress: editableDestinationAddress ? formatAddressForPrimeAPI(destinationAddress) : null,
      destinationType,
      diversion,
    };
    mutateMTOShipment({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag, body }).then(() => {
      setSubmitting(false);
    });
  };

  const initialValues = {
    estimatedWeight: shipment.primeEstimatedWeight?.toLocaleString(),
    actualWeight: shipment.primeActualWeight?.toLocaleString(),
    requestedPickupDate: shipment.requestedPickupDate,
    scheduledPickupDate: shipment.scheduledPickupDate,
    actualPickupDate: shipment.actualPickupDate,
    pickupAddress: editablePickupAddress ? emptyAddress : reformatPrimeApiPickupAddress,
    destinationAddress: editableDestinationAddress ? emptyAddress : reformatPrimeApiDestinationAddress,
    destinationType: shipment.destinationType,
    diversion: shipment.diversion,
  };

  const validationSchema = Yup.object().shape({
    pickupAddress: addressSchema,
    destinationAddress: addressSchema,
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
                <div className={primeStyles.errorContainer}>
                  <Alert type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
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
                      <PrimeUIShipmentUpdateForm
                        editableWeightEstimateField={editableWeightEstimateField}
                        editableWeightActualField={editableWeightActualField}
                        editablePickupAddress={editablePickupAddress}
                        editableDestinationAddress={editableDestinationAddress}
                        estimatedWeight={initialValues.estimatedWeight}
                        actualWeight={initialValues.actualWeight}
                        requestedPickupDate={initialValues.requestedPickupDate}
                        pickupAddress={initialValues.pickupAddress}
                        destinationAddress={initialValues.destinationAddress}
                        diversion={initialValues.diversion}
                      />
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

PrimeUIShipmentUpdate.propTypes = {
  setFlashMessage: func,
};

PrimeUIShipmentUpdate.defaultProps = {
  setFlashMessage: () => {},
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default withRouter(connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentUpdate));
