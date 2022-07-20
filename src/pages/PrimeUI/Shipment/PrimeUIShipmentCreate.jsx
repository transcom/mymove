import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useHistory, useParams, withRouter } from 'react-router-dom';
import { generatePath } from 'react-router';
import { useMutation } from 'react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
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
import { isEmpty, isValidWeight } from 'shared/utils';
import { formatAddressForPrimeAPI, formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import PrimeUIShipmentCreateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreateForm';
import { OptionalAddressSchema } from 'components/Customer/MtoShipmentForm/validationSchemas';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { InvalidZIPTypeError, ZIP5_CODE_REGEX } from 'utils/validation';

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
    const { shipmentType } = values;
    const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;

    let body;
    if (isPPM) {
      const {
        counselorRemarks,
        ppmShipment: {
          expectedDepartureDate,
          pickupPostalCode,
          secondaryPickupPostalCode,
          destinationPostalCode,
          secondaryDestinationPostalCode,
          sitExpected,
          sitLocation,
          sitEstimatedWeight,
          sitEstimatedEntryDate,
          sitEstimatedDepartureDate,
          estimatedWeight,
          hasProGear,
          proGearWeight,
          spouseProGearWeight,
        },
      } = values;

      body = {
        moveTaskOrderID: moveCodeOrID,
        shipmentType,
        counselorRemarks: counselorRemarks || null,
        ppmShipment: {
          expectedDepartureDate: expectedDepartureDate ? formatSwaggerDate(expectedDepartureDate) : null,
          pickupPostalCode,
          secondaryPickupPostalCode: secondaryPickupPostalCode || null,
          destinationPostalCode,
          secondaryDestinationPostalCode: secondaryDestinationPostalCode || null,
          sitExpected,
          ...(sitExpected && {
            sitLocation: sitLocation || null,
            sitEstimatedWeight: sitEstimatedWeight ? parseInt(sitEstimatedWeight, 10) : null,
            sitEstimatedEntryDate: sitEstimatedEntryDate ? formatSwaggerDate(sitEstimatedEntryDate) : null,
            sitEstimatedDepartureDate: sitEstimatedDepartureDate ? formatSwaggerDate(sitEstimatedDepartureDate) : null,
          }),
          estimatedWeight: estimatedWeight ? parseInt(estimatedWeight, 10) : null,
          hasProGear,
          ...(hasProGear && {
            proGearWeight: proGearWeight ? parseInt(proGearWeight, 10) : null,
            spouseProGearWeight: spouseProGearWeight ? parseInt(spouseProGearWeight, 10) : null,
          }),
        },
      };
    } else {
      const { requestedPickupDate, estimatedWeight, pickupAddress, destinationAddress, diversion } = values;

      body = {
        moveTaskOrderID: moveCodeOrID,
        shipmentType,
        requestedPickupDate: requestedPickupDate ? formatSwaggerDate(requestedPickupDate) : null,
        primeEstimatedWeight: isValidWeight(estimatedWeight) ? parseInt(estimatedWeight, 10) : null,
        pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
        destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
        diversion: diversion || null,
      };
    }

    mutateCreateMTOShipment({ body }).then(() => {
      setSubmitting(false);
    });
  };

  const initialValues = {
    shipmentType: '',

    // PPM
    counselorRemarks: '',
    ppmShipment: {
      expectedDepartureDate: '',
      pickupPostalCode: '',
      secondaryPickupPostalCode: '',
      destinationPostalCode: '',
      secondaryDestinationPostalCode: '',
      sitExpected: false,
      sitLocation: '',
      sitEstimatedWeight: '',
      sitEstimatedEntryDate: '',
      sitEstimatedDepartureDate: '',
      estimatedWeight: '',
      hasProGear: false,
      proGearWeight: '',
      spouseProGearWeight: '',
    },

    // Other shipment types
    requestedPickupDate: '',
    estimatedWeight: '',
    pickupAddress: {},
    destinationAddress: {},
    diversion: '',
  };

  const validationSchema = Yup.object().shape({
    shipmentType: Yup.string().required(),

    // PPM
    ppmShipment: Yup.object().when('shipmentType', {
      is: (shipmentType) => shipmentType === 'PPM',
      then: Yup.object().shape({
        expectedDepartureDate: Yup.date()
          .required('Required')
          .typeError('Invalid date. Must be in the format: DD MMM YYYY'),
        pickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
        secondaryPickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).nullable(),
        destinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
        secondaryDestinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).nullable(),
        sitExpected: Yup.boolean().required('Required'),
        sitLocation: Yup.string().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        sitEstimatedWeight: Yup.number().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        // TODO: Figure out how to validate this but be optional.  Right now, when you uncheck
        //  sitEnabled, the "Save" button remains disabled in certain situations.
        // sitEstimatedEntryDate: Yup.date().when('sitExpected', {
        //   is: true,
        //   then: (schema) =>
        //     schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
        // }),
        // sitEstimatedDepartureDate: Yup.date().when('sitExpected', {
        //   is: true,
        //   then: (schema) =>
        //     schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
        // }),
        estimatedWeight: Yup.number().required('Required'),
        hasProGear: Yup.boolean().required('Required'),
        proGearWeight: Yup.number().when(['hasProGear', 'spouseProGearWeight'], {
          is: (hasProGear, spouseProGearWeight) => hasProGear && !spouseProGearWeight,
          then: (schema) =>
            schema.required(
              `Enter a weight into at least one pro-gear field. If you won't have pro-gear, uncheck above.`,
            ),
        }),
        spouseProGearWeight: Yup.number(),
      }),
    }),
    // counselorRemarks is an optional string

    // Other shipment types
    requestedPickupDate: Yup.date().when('shipmentType', {
      is: (shipmentType) => shipmentType !== 'PPM',
      then: (schema) => schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).'),
    }),
    pickupAddress: Yup.object().when('shipmentType', {
      is: (shipmentType) => shipmentType !== 'PPM',
      then: () => OptionalAddressSchema,
    }),
    destinationAddress: Yup.object().when('shipmentType', {
      is: (shipmentType) => shipmentType !== 'PPM',
      then: () => OptionalAddressSchema,
    }),
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
