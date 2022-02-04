// eslint-disable no-unused-vars
import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, FormGroup } from '@trussworks/react-uswds';

import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
// import { DutyStationShape } from 'types';
import { ZIP_CODE_REGEX } from 'utils/validation';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField } from 'components/form/fields';
import Hint from 'components/Hint/index';

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

const DatesAndLocation = ({
  mtoShipment,
  // destinationDutyStation,
  serviceMember,
  onBack,
  onSubmit,
  postalCodeValidator,
}) => {
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

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, values }) => {
        return (
          <Form>
            <TextField
              label="ZIP"
              id="pickupPostalCode"
              name="pickupPostalCode"
              maxLength={10}
              validate={(value) => postalCodeValidator(value, 'origin')}
            />
            <CheckboxField
              id="useCurrentZip"
              name="useCurrentZip"
              label={`Use my current ZIP (${serviceMember?.residentialAddress?.postalCode})`}
            />
            <FormGroup>
              <p>Will you add items to your PPM from a place in a different ZIP code?</p>
              <Field
                as={Radio}
                data-testid="yes-secondary-pickup-postal-code"
                id="yes-secondary-pickup-postal-code"
                label="Yes"
                name="hasSecondaryPickupPostalCode"
                value="yes"
                checked={values.hasSecondaryPickupPostalCode === 'yes'}
              />
              <Field
                as={Radio}
                data-testid="no-secondary-pickup-postal-code"
                id="no-secondary-pickup-postal-code"
                label="No"
                name="hasSecondaryPickupPostalCode"
                value="no"
                checked={values.hasSecondaryPickupPostalCode === 'no'}
              />
            </FormGroup>
            <TextField
              label="Second ZIP"
              id="secondaryPickupPostalCode"
              name="secondaryPickupPostalCode"
              maxLength={10}
              validate={(value) => postalCodeValidator(value, 'origin')}
            />
            <Hint>
              <p>A second origin ZIP could mean that your final incentive is lower than your estimate.</p>

              <p>
                Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to your move
                counselor for more detailed information.
              </p>
            </Hint>
            <Button type="button" unstyled onClick={onBack} data-testid="datesAndLocationBackBtn">
              Back
            </Button>
            <Button type="submit" unstyled data-testid="datesAndLocationSubmitBtn" disabled={!isValid || isSubmitting}>
              Save & Continue
            </Button>
          </Form>
        );
      }}
    </Formik>
  );
};

DatesAndLocation.propTypes = {
  mtoShipment: MtoShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  postalCodeValidator: func.isRequired,
};

DatesAndLocation.defaultProps = {
  mtoShipment: undefined,
};

export default DatesAndLocation;
