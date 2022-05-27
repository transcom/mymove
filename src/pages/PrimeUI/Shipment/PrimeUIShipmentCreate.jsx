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
import { primeSimulatorRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';
import { createPrimeMTOShipment } from 'services/primeApi';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { addressSchema } from 'utils/validation';
import { isValidWeight, isEmpty } from 'shared/utils';
import { formatAddressForPrimeAPI, formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import PrimeUIShipmentCreateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreateForm';

const PrimeUIShipmentCreate = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID } = useParams();
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };
  const [mutateCreateMTOShipment] = useMutation(createPrimeMTOShipment, {
    onSuccess: (createdMTOShipment) => {
      setFlashMessage(
        `MSG_CREATE_PAYMENT_SUCCESS${createdMTOShipment.id}`,
        'success',
        `Successfully created shipment ${createdMTOShipment.id}`,
        '',
        true,
      );
      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        let invalidFieldsStr = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${invalidFieldsStr}\n\nPlease try again`,
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

  const onSubmit = (values, { setSubmitting }) => {
    const { shipmentType, requestedPickupDate, estimatedWeight, pickupAddress, destinationAddress, diversion } = values;

    const body = {
      moveTaskOrderID: moveCodeOrID,
      shipmentType,
      requestedPickupDate: requestedPickupDate ? formatSwaggerDate(requestedPickupDate) : null,
      primeEstimatedWeight: isValidWeight(estimatedWeight) ? parseInt(estimatedWeight, 10) : null,
      pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
      destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
      diversion: diversion || null,
    };
    mutateCreateMTOShipment({ body }).then(() => {
      setSubmitting(false);
    });
  };

  const initialValues = {
    shipmentType: '',
    requestedPickupDate: '',
    estimatedWeight: '',
    pickupAddress: {},
    destinationAddress: {},
    diversion: '',
  };

  const validationSchema = Yup.object().shape({
    shipmentType: Yup.string(),
    requestedPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    pickupAddress: addressSchema.optional(),
    destinationAddress: addressSchema.optional(),
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
                      <PrimeUIShipmentCreateForm />
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

PrimeUIShipmentCreate.propTypes = {
  setFlashMessage: func,
};

PrimeUIShipmentCreate.defaultProps = {
  setFlashMessage: () => {},
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default withRouter(connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentCreate));
