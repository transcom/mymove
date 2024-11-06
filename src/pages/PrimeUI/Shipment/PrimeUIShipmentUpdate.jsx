import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Grid, GridContainer } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';
import { updatePrimeMTOShipmentV3, updatePrimeMTOShipmentStatus } from 'services/primeApi';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { requiredAddressSchema, addressSchema, partialRequiredAddressSchema } from 'utils/validation';
import { isEmpty, isValidWeight } from 'shared/utils';
import {
  formatAddressForPrimeAPI,
  formatExtraAddressForPrimeAPI,
  formatSwaggerDate,
  fromPrimeAPIAddressFormat,
} from 'utils/formatters';
import PrimeUIShipmentUpdateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdateForm';
import PrimeUIShipmentUpdatePPMForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdatePPMForm';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { OptionalAddressSchema } from 'components/Customer/MtoShipmentForm/validationSchemas';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const PrimeUIShipmentUpdate = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const { mutateAsync: mutateMTOShipmentStatus } = useMutation(updatePrimeMTOShipmentStatus, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      setFlashMessage(`MSG_CANCELATION_SUCCESS${shipmentId}`, 'success', `Successfully canceled shipment`, '', true);
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

  const { mutateAsync: mutateMTOShipment } = useMutation(updatePrimeMTOShipmentV3, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      setFlashMessage(`MSG_CREATE_PAYMENT_SUCCESS${shipmentId}`, 'success', `Successfully updated shipment`, '', true);
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

  const isPPM = shipment.shipmentType === SHIPMENT_OPTIONS.PPM;

  const emptyAddress = {
    streetAddress1: '',
    streetAddress2: '',
    streetAddress3: '',
    city: '',
    state: '',
    postalCode: '',
  };

  const editableWeightEstimateField =
    !isValidWeight(shipment.primeEstimatedWeight) && shipment.shipmentType !== SHIPMENT_OPTIONS.NTSR;
  const editableWeightActualField = true;
  const editableProGearWeightActualField = true;
  const editableSpouseProGearWeightActualField = true;
  const reformatPrimeApiPickupAddress = fromPrimeAPIAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiSecondaryPickupAddress = fromPrimeAPIAddressFormat(shipment.secondaryPickupAddress);
  const reformatPrimeApiTertiaryPickupAddress = fromPrimeAPIAddressFormat(shipment.tertiaryPickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeAPIAddressFormat(shipment.destinationAddress);
  const reformatPrimeApiSecondaryDeliveryAddress = fromPrimeAPIAddressFormat(shipment.secondaryDeliveryAddress);
  const reformatPrimeApiTertiaryDeliveryAddress = fromPrimeAPIAddressFormat(shipment.tertiaryDeliveryAddress);
  const editablePickupAddress = isEmpty(reformatPrimeApiPickupAddress);
  const editableSecondaryPickupAddress = isEmpty(reformatPrimeApiSecondaryPickupAddress);
  const editableTertiaryPickupAddress = isEmpty(reformatPrimeApiTertiaryPickupAddress);
  const editableDestinationAddress = isEmpty(reformatPrimeApiDestinationAddress);
  const editableSecondaryDeliveryAddress = isEmpty(reformatPrimeApiSecondaryDeliveryAddress);
  const editableTertiaryDeliveryAddress = isEmpty(reformatPrimeApiTertiaryDeliveryAddress);

  const onCancelShipmentClick = () => {
    mutateMTOShipmentStatus({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag }).then(() => {
      /* console.info("It's done and canceled."); */
    });
  };

  const onSubmit = (values, { setSubmitting }) => {
    let body;
    if (isPPM) {
      const {
        ppmShipment: {
          expectedDepartureDate,
          pickupAddress,
          secondaryPickupAddress,
          tertiaryPickupAddress,
          destinationAddress,
          secondaryDestinationAddress,
          tertiaryDestinationAddress,
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
          hasTertiaryPickupAddress,
          hasSecondaryDestinationAddress,
          hasTertiaryDestinationAddress,
        },
        counselorRemarks,
      } = values;
      body = {
        ppmShipment: {
          expectedDepartureDate: expectedDepartureDate ? formatSwaggerDate(expectedDepartureDate) : null,
          pickupAddress: isEmpty(pickupAddress) ? null : formatAddressForPrimeAPI(pickupAddress),
          secondaryPickupAddress: isEmpty(secondaryPickupAddress)
            ? emptyAddress
            : formatAddressForPrimeAPI(secondaryPickupAddress),
          tertiaryPickupAddress: isEmpty(tertiaryPickupAddress)
            ? emptyAddress
            : formatAddressForPrimeAPI(tertiaryPickupAddress),
          destinationAddress: isEmpty(destinationAddress) ? null : formatAddressForPrimeAPI(destinationAddress),
          secondaryDestinationAddress: isEmpty(secondaryDestinationAddress)
            ? emptyAddress
            : formatAddressForPrimeAPI(secondaryDestinationAddress),
          tertiaryDestinationAddress: isEmpty(tertiaryDestinationAddress)
            ? emptyAddress
            : formatAddressForPrimeAPI(tertiaryDestinationAddress),
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
          hasSecondaryPickupAddress: hasSecondaryPickupAddress === 'true',
          hasTertiaryPickupAddress: hasTertiaryPickupAddress === 'true',
          hasSecondaryDestinationAddress: hasSecondaryDestinationAddress === 'true',
          hasTertiaryDestinationAddress: hasTertiaryDestinationAddress === 'true',
        },
        counselorRemarks: counselorRemarks || null,
      };
    } else {
      const {
        estimatedWeight,
        actualWeight,
        actualProGearWeight,
        actualSpouseProGearWeight,
        actualPickupDate,
        scheduledPickupDate,
        actualDeliveryDate,
        scheduledDeliveryDate,
        pickupAddress,
        secondaryPickupAddress,
        tertiaryPickupAddress,
        destinationAddress,
        secondaryDeliveryAddress,
        tertiaryDeliveryAddress,
        destinationType,
        diversion,
      } = values;

      body = {
        primeEstimatedWeight: editableWeightEstimateField ? parseInt(estimatedWeight, 10) : null,
        primeActualWeight: parseInt(actualWeight, 10),
        actualProGearWeight: parseInt(actualProGearWeight, 10),
        actualSpouseProGearWeight: parseInt(actualSpouseProGearWeight, 10),
        scheduledPickupDate: scheduledPickupDate ? formatSwaggerDate(scheduledPickupDate) : null,
        actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
        scheduledDeliveryDate: scheduledDeliveryDate ? formatSwaggerDate(scheduledDeliveryDate) : null,
        actualDeliveryDate: actualDeliveryDate ? formatSwaggerDate(actualDeliveryDate) : null,
        pickupAddress: editablePickupAddress ? formatAddressForPrimeAPI(pickupAddress) : null,
        secondaryPickupAddress: editableSecondaryPickupAddress
          ? formatExtraAddressForPrimeAPI(secondaryPickupAddress)
          : null,
        tertiaryPickupAddress: editableTertiaryPickupAddress
          ? formatExtraAddressForPrimeAPI(tertiaryPickupAddress)
          : null,
        destinationAddress: editableDestinationAddress ? formatAddressForPrimeAPI(destinationAddress) : null,
        secondaryDeliveryAddress: editableSecondaryDeliveryAddress
          ? formatExtraAddressForPrimeAPI(secondaryDeliveryAddress)
          : null,
        tertiaryDeliveryAddress: editableTertiaryDeliveryAddress
          ? formatExtraAddressForPrimeAPI(tertiaryDeliveryAddress)
          : null,
        destinationType,
        diversion,
      };
    }

    mutateMTOShipment({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag, body }).then(() => {
      setSubmitting(false);
    });
  };

  let initialValues;
  let validationSchema;
  if (isPPM) {
    initialValues = {
      ppmShipment: {
        pickupAddress: shipment.ppmShipment.pickupAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.pickupAddress)
          : emptyAddress,
        secondaryPickupAddress: shipment.ppmShipment.secondaryPickupAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.secondaryPickupAddress)
          : emptyAddress,
        tertiaryPickupAddress: shipment.ppmShipment.tertiaryPickupAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.tertiaryPickupAddress)
          : emptyAddress,
        destinationAddress: shipment.ppmShipment.destinationAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.destinationAddress)
          : emptyAddress,
        secondaryDestinationAddress: shipment.ppmShipment.secondaryDestinationAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.secondaryDestinationAddress)
          : emptyAddress,
        tertiaryDestinationAddress: shipment.ppmShipment.tertiaryDestinationAddress
          ? formatAddressForPrimeAPI(shipment.ppmShipment.tertiaryDestinationAddress)
          : emptyAddress,
        sitExpected: shipment.ppmShipment.sitExpected,
        sitLocation: shipment.ppmShipment.sitLocation,
        sitEstimatedWeight: shipment.ppmShipment.sitEstimatedWeight?.toString(),
        sitEstimatedEntryDate: shipment.ppmShipment.sitEstimatedEntryDate,
        sitEstimatedDepartureDate: shipment.ppmShipment.sitEstimatedDepartureDate,
        estimatedWeight: shipment.ppmShipment.estimatedWeight?.toString(),
        expectedDepartureDate: shipment.ppmShipment.expectedDepartureDate,
        hasProGear: shipment.ppmShipment.hasProGear,
        proGearWeight: shipment.ppmShipment.proGearWeight?.toString(),
        spouseProGearWeight: shipment.ppmShipment.spouseProGearWeight?.toString(),
        hasSecondaryPickupAddress: shipment.ppmShipment.hasSecondaryPickupAddress ? 'true' : 'false',
        hasTertiaryPickupAddress: shipment.ppmShipment.hasTertiaryPickupAddress ? 'true' : 'false',
        hasSecondaryDestinationAddress: shipment.ppmShipment.hasSecondaryDestinationAddress ? 'true' : 'false',
        hasTertiaryDestinationAddress: shipment.ppmShipment.hasTertiaryDestinationAddress ? 'true' : 'false',
      },
      counselorRemarks: shipment.counselorRemarks || '',
    };

    validationSchema = Yup.object().shape({
      ppmShipment: Yup.object().shape({
        expectedDepartureDate: Yup.date()
          .required('Required')
          .typeError('Invalid date. Must be in the format: DD MMM YYYY'),
        pickupAddress: requiredAddressSchema.required('Required'),
        secondaryPickupAddress: OptionalAddressSchema,
        tertiaryPickupAddress: OptionalAddressSchema,
        destinationAddress: partialRequiredAddressSchema.required('Required'),
        secondaryDestinationAddress: OptionalAddressSchema,
        tertiaryDestinationAddress: OptionalAddressSchema,
        sitExpected: Yup.boolean().required('Required'),
        sitLocation: Yup.string().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        sitEstimatedWeight: Yup.number().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        // TODO: Figure out how to validate this but be optional.  Right now, when you first
        //   go to the page with sitEnabled of false, the "Save" button remains disabled.
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
      // counselorRemarks is an optional string
    });
  } else {
    initialValues = {
      estimatedWeight: shipment.primeEstimatedWeight?.toLocaleString(),
      actualWeight: shipment.primeActualWeight?.toLocaleString(),
      actualProGearWeight: shipment.actualProGearWeight?.toLocaleString(),
      actualSpouseProGearWeight: shipment.actualSpouseProGearWeight?.toLocaleString(),
      requestedPickupDate: shipment.requestedPickupDate,
      scheduledPickupDate: shipment.scheduledPickupDate,
      actualPickupDate: shipment.actualPickupDate,
      scheduledDeliveryDate: shipment.scheduledDeliveryDate,
      actualDeliveryDate: shipment.actualDeliveryDate,
      pickupAddress: editablePickupAddress ? emptyAddress : reformatPrimeApiPickupAddress,
      secondaryPickupAddress: editableSecondaryPickupAddress ? emptyAddress : reformatPrimeApiSecondaryPickupAddress,
      tertiaryPickupAddress: editableTertiaryPickupAddress ? emptyAddress : reformatPrimeApiTertiaryPickupAddress,
      destinationAddress: editableDestinationAddress ? emptyAddress : reformatPrimeApiDestinationAddress,
      secondaryDeliveryAddress: editableSecondaryDeliveryAddress
        ? emptyAddress
        : reformatPrimeApiSecondaryDeliveryAddress,
      tertiaryDeliveryAddress: editableTertiaryDeliveryAddress ? emptyAddress : reformatPrimeApiTertiaryDeliveryAddress,
      destinationType: shipment.destinationType,
      diversion: shipment.diversion,
    };

    validationSchema = Yup.object().shape({
      pickupAddress: addressSchema,
      destinationAddress: addressSchema,
      scheduledPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
      actualPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    });
  }

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <Button type="button" onClick={onCancelShipmentClick} className="usa-button usa-button-secondary">
                Cancel Shipment
              </Button>
              <Formik
                initialValues={initialValues}
                onSubmit={onSubmit}
                validationSchema={validationSchema}
                validateOnMount
              >
                {({ isValid, isSubmitting, handleSubmit }) => {
                  return (
                    <Form className={formStyles.form}>
                      {isPPM ? (
                        <PrimeUIShipmentUpdatePPMForm />
                      ) : (
                        <PrimeUIShipmentUpdateForm
                          editableWeightEstimateField={editableWeightEstimateField}
                          editableWeightActualField={editableWeightActualField}
                          editableProGearWeightActualField={editableProGearWeightActualField}
                          editableSpouseProGearWeightActualField={editableSpouseProGearWeightActualField}
                          editablePickupAddress={editablePickupAddress}
                          editableSecondaryPickupAddress={editableSecondaryPickupAddress}
                          editableTertiaryPickupAddress={editableTertiaryPickupAddress}
                          editableDestinationAddress={editableDestinationAddress}
                          editableSecondaryDeliveryAddress={editableSecondaryDeliveryAddress}
                          editableTertiaryDeliveryAddress={editableTertiaryDeliveryAddress}
                          estimatedWeight={initialValues.estimatedWeight}
                          actualWeight={initialValues.actualWeight}
                          actualProGearWeight={initialValues.actualProGearWeight}
                          actualSpouseProGearWeight={initialValues.actualSpouseProGearWeight}
                          requestedPickupDate={initialValues.requestedPickupDate}
                          pickupAddress={initialValues.pickupAddress}
                          secondaryPickupAddress={initialValues.secondaryPickupAddress}
                          tertiaryPickupAddress={initialValues.tertiaryPickupAddress}
                          destinationAddress={initialValues.destinationAddress}
                          secondaryDeliveryAddress={initialValues.secondaryDeliveryAddress}
                          tertiaryDeliveryAddress={initialValues.tertiaryDeliveryAddress}
                          diversion={initialValues.diversion}
                          shipmentType={shipment.shipmentType}
                        />
                      )}
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

export default connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentUpdate);
