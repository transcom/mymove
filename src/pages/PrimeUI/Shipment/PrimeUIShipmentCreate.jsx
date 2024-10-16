import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';
import { createPrimeMTOShipmentV3 } from 'services/primeApi';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { isEmpty, isValidWeight } from 'shared/utils';
import { formatAddressForPrimeAPI, formatSwaggerDate } from 'utils/formatters';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { requiredAddressSchema } from 'utils/validation';
import PrimeUIShipmentCreateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentCreateForm';
import { OptionalAddressSchema } from 'components/Customer/MtoShipmentForm/validationSchemas';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';

const PrimeUIShipmentCreate = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };
  const { mutateAsync: mutateCreateMTOShipment } = useMutation(createPrimeMTOShipmentV3, {
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
    const isBoat = shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY || shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY;

    let body;
    if (isPPM) {
      const {
        counselorRemarks,
        ppmShipment: {
          expectedDepartureDate,
          pickupAddress,
          secondaryPickupAddress,
          destinationAddress,
          secondaryDestinationAddress,
          sitExpected,
          sitLocation,
          sitEstimatedWeight,
          sitEstimatedEntryDate,
          sitEstimatedDepartureDate,
          estimatedWeight,
          hasProGear,
          proGearWeight,
          spouseProGearWeight,
          hasSecondaryPickupAddress,
          hasSecondaryDestinationAddress,
        },
      } = values;

      body = {
        moveTaskOrderID: moveCodeOrID,
        shipmentType,
        counselorRemarks: counselorRemarks || null,
        ppmShipment: {
          expectedDepartureDate: expectedDepartureDate ? formatSwaggerDate(expectedDepartureDate) : null,
          pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
          secondaryPickupAddress: isEmpty(secondaryPickupAddress)
            ? null
            : formatAddressForPrimeAPI(secondaryPickupAddress),
          destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
          secondaryDestinationAddress: isEmpty(secondaryDestinationAddress)
            ? null
            : formatAddressForPrimeAPI(secondaryDestinationAddress),
          sitExpected,
          ...(sitExpected && {
            sitLocation: sitLocation || null,
            sitEstimatedWeight: sitEstimatedWeight ? parseInt(sitEstimatedWeight, 10) : null,
            sitEstimatedEntryDate: sitEstimatedEntryDate ? formatSwaggerDate(sitEstimatedEntryDate) : null,
            sitEstimatedDepartureDate: sitEstimatedDepartureDate ? formatSwaggerDate(sitEstimatedDepartureDate) : null,
          }),
          hasSecondaryPickupAddress: hasSecondaryPickupAddress === 'true',
          hasSecondaryDestinationAddress: hasSecondaryDestinationAddress === 'true',
          estimatedWeight: estimatedWeight ? parseInt(estimatedWeight, 10) : null,
          hasProGear,
          ...(hasProGear && {
            proGearWeight: proGearWeight ? parseInt(proGearWeight, 10) : null,
            spouseProGearWeight: spouseProGearWeight ? parseInt(spouseProGearWeight, 10) : null,
          }),
        },
      };
    } else if (isBoat) {
      const {
        counselorRemarks,
        requestedPickupDate,
        estimatedWeight,
        pickupAddress,
        destinationAddress,
        diversion,
        divertedFromShipmentId,
        boatShipment: {
          year,
          make,
          model,
          lengthInFeet,
          lengthInInches,
          widthInFeet,
          widthInInches,
          heightInFeet,
          heightInInches,
          hasTrailer,
          isRoadworthy,
        },
      } = values;

      // Sum the feet and inches fields into only inches for backend/db
      const totalLengthInInches = parseInt(lengthInFeet, 10) * 12 + parseInt(lengthInInches, 10);
      const totalWidthInInches = parseInt(widthInFeet, 10) * 12 + parseInt(widthInInches, 10);
      const totalHeightInInches = parseInt(heightInFeet, 10) * 12 + parseInt(heightInInches, 10);

      body = {
        moveTaskOrderID: moveCodeOrID,
        shipmentType,
        counselorRemarks: counselorRemarks || null,
        boatShipment: {
          year: year ? parseInt(year, 10) : null,
          make: make || null,
          model: model || null,
          lengthInInches: totalLengthInInches,
          widthInInches: totalWidthInInches,
          heightInInches: totalHeightInInches,
          hasTrailer,
          ...(hasTrailer && {
            isRoadworthy,
          }),
        },
        requestedPickupDate: requestedPickupDate ? formatSwaggerDate(requestedPickupDate) : null,
        primeEstimatedWeight: isValidWeight(estimatedWeight) ? parseInt(estimatedWeight, 10) : null,
        pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
        destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
        diversion: diversion || null,
        divertedFromShipmentId: divertedFromShipmentId || null,
      };
    } else {
      const {
        requestedPickupDate,
        estimatedWeight,
        pickupAddress,
        destinationAddress,
        diversion,
        divertedFromShipmentId,
      } = values;

      body = {
        moveTaskOrderID: moveCodeOrID,
        shipmentType,
        requestedPickupDate: requestedPickupDate ? formatSwaggerDate(requestedPickupDate) : null,
        primeEstimatedWeight: isValidWeight(estimatedWeight) ? parseInt(estimatedWeight, 10) : null,
        pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
        destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
        diversion: diversion || null,
        divertedFromShipmentId: divertedFromShipmentId || null,
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
      pickupAddress: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      secondaryPickupAddress: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      destinationAddress: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      secondaryDestinationAddress: {
        streetAddress1: '',
        streetAddress2: '',
        city: '',
        state: '',
        postalCode: '',
      },
      sitExpected: false,
      sitLocation: '',
      sitEstimatedWeight: '',
      sitEstimatedEntryDate: '',
      sitEstimatedDepartureDate: '',
      estimatedWeight: '',
      hasProGear: false,
      proGearWeight: '',
      spouseProGearWeight: '',
      hasSecondaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
    },

    // Boat Shipment
    boatShipment: {
      year: null,
      make: '',
      model: '',
      lengthInFeet: null,
      lengthInInches: null,
      widthInFeet: null,
      widthInInches: null,
      heightInFeet: null,
      heightInInches: null,
      hasTrailer: false,
      isRoadworthy: false,
    },

    // Other shipment types
    requestedPickupDate: '',
    estimatedWeight: '',
    pickupAddress: {},
    destinationAddress: {},
    diversion: '',
    divertedFromShipmentId: '',
  };

  const validationSchema = Yup.object().shape({
    shipmentType: Yup.string().required(),

    // PPM
    ppmShipment: Yup.object().when('shipmentType', {
      is: (shipmentType) => shipmentType === 'PPM',
      then: () =>
        Yup.object().shape({
          expectedDepartureDate: Yup.date()
            .required('Required')
            .typeError('Invalid date. Must be in the format: DD MMM YYYY'),
          pickupAddress: requiredAddressSchema.required('Required'),
          secondaryPickupAddress: OptionalAddressSchema,
          destinationAddress: requiredAddressSchema.required('Required'),
          secondaryDestinationAddress: OptionalAddressSchema,
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

    // Boat Shipment
    boatShipment: Yup.object().when('shipmentType', {
      is: (shipmentType) =>
        shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY || shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY,
      then: () =>
        Yup.object().shape({
          year: Yup.number().positive('Must be a postive number').required('Required'),
          make: Yup.string().min(1).trim().required('Required'),
          model: Yup.string().min(1).trim().required('Required'),
          lengthInFeet: Yup.number().moreThan(-1, 'Must be a positive number').required('Required'),
          lengthInInches: Yup.number().moreThan(-1, 'Must be a positive number').max(11).required('Required'),
          widthInFeet: Yup.number().moreThan(-1, 'Must be a positive number').required('Required'),
          widthInInches: Yup.number().moreThan(-1, 'Must be a positive number').max(11).required('Required'),
          heightInFeet: Yup.number().moreThan(-1, 'Must be a positive number').required('Required'),
          heightInInches: Yup.number().moreThan(-1, 'Must be a positive number').max(11).required('Required'),
          hasTrailer: Yup.boolean().required(),
          isRoadworthy: Yup.boolean().when('hasTrailer', {
            is: true,
            then: (schema) => schema.required('Required'),
          }),
        }),
    }),

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

export default connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentCreate);
