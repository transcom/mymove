import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik } from 'formik';

import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { DutyStationShape } from 'types';
import { ZIP_CODE_REGEX } from 'utils/validation';

// TODO: conditional validation for optional ZIPs
const validationSchema = Yup.object().shape({
  pickupPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  useResidentialAddressZIP: Yup.boolean(),
  hasSecondaryPickupPostalCode: Yup.boolean().required('Required'),
  secondaryPickupPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code'),
  useDestinationDutyLocationZIP: Yup.boolean(),
  destinationPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code'),
  hasSecondaryDestinationPostalCode: Yup.boolean().required('Required'),
  secondaryDestinationPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code'),
  sitExpected: Yup.boolean().required('Required'),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
});

// eslint-disable-next-line no-unused-vars
const DatesAndLocation = ({ mtoShipment, destinationDutyStation, serviceMember, onBack, onSubmit }) => {
  const initialValues = {
    pickupPostalCode: mtoShipment?.ppmShipment?.pickupPostalCode || '',
    useResidentialAddressZIP: '',
    hasSecondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || 'no',
    secondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || '',
    useDestinationDutyLocationZIP: '',
    destinationPostalCode: mtoShipment?.ppmShipment?.destinationPostalCode || '',
    hasSecondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || 'no',
    secondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || '',
    sitExpected: mtoShipment?.ppmShipment?.sitExpected || 'no',
    expectedDepartureDate: mtoShipment?.ppmShipment?.expectedDepartureDate || '',
  };

  // TODO: async validation call to validate postal codes are valid for rate engine

  return <Formik initialValues={initialValues} validationSchema={validationSchema} />;
};

DatesAndLocation.propTypes = {
  mtoShipment: MtoShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyStation: DutyStationShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

DatesAndLocation.defaultProps = {
  mtoShipment: undefined,
};
