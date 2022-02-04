// eslint-disable no-unused-vars
import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, FormGroup } from '@trussworks/react-uswds';

import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { ZIP_CODE_REGEX } from 'utils/validation';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { DutyStationShape } from 'types';

const validationSchema = Yup.object().shape({
  pickupPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  useResidentialAddressZIP: Yup.boolean(),
  hasSecondaryPickupPostalCode: Yup.boolean().required('Required'),
  secondaryPickupPostalCode: Yup.string().when('hasSecondaryPickupPostalCode', {
    is: true,
    then: (schema) => schema.matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  }),
  useDestinationDutyLocationZIP: Yup.boolean(),
  destinationPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code'),
  hasSecondaryDestinationPostalCode: Yup.boolean().required('Required'),
  secondaryDestinationPostalCode: Yup.string().when('hasSecondaryDestinationPostalCode', {
    is: true,
    then: (schema) => schema.matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  }),
  sitExpected: Yup.boolean().required('Required'),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
});

const DatesAndLocation = ({
  mtoShipment,
  destinationDutyStation,
  serviceMember,
  onBack,
  onSubmit,
  postalCodeValidator,
}) => {
  const initialValues = {
    pickupPostalCode: mtoShipment?.ppmShipment?.pickupPostalCode || '',
    useResidentialAddressZIP: '',
    hasSecondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || 'false',
    secondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || '',
    useDestinationDutyLocationZIP: '',
    destinationPostalCode: mtoShipment?.ppmShipment?.destinationPostalCode || '',
    hasSecondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || 'false',
    secondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || '',
    sitExpected: mtoShipment?.ppmShipment?.sitExpected || 'false',
    expectedDepartureDate: mtoShipment?.ppmShipment?.expectedDepartureDate || '',
  };

  // TODO: async validation call to validate postal codes are valid for rate engine

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <Form>
            <TextField
              label="ZIP"
              id="pickupPostalCode"
              name="pickupPostalCode"
              maxLength={10}
              validate={(value) => postalCodeValidator(value, 'origin')}
            />
            {/* TODO: call setFieldValue when this checkbox is selected to populate pickupPostalCode */}
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
                value="true"
                checked={values.hasSecondaryPickupPostalCode === 'true'}
              />
              <Field
                as={Radio}
                data-testid="no-secondary-pickup-postal-code"
                id="no-secondary-pickup-postal-code"
                label="No"
                name="hasSecondaryPickupPostalCode"
                value="false"
                checked={values.hasSecondaryPickupPostalCode === 'false'}
              />
            </FormGroup>
            {values.hasSecondaryPickupPostalCode === 'true' && (
              <>
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
                    Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to your
                    move counselor for more detailed information.
                  </p>
                </Hint>
              </>
            )}
            <TextField
              label="ZIP"
              id="destinationPostalCode"
              name="destinationPostalCode"
              maxLength={10}
              validate={(value) => postalCodeValidator(value, 'destination')}
            />
            <CheckboxField
              id="useDestinationDutyLocationZIP"
              name="useDestinationDutyLocationZIP"
              label={`Use the ZIP for my new duty location (${destinationDutyStation?.address?.postalCode})`}
            />
            <Hint>
              Use the ZIP for your new address if you know it. Use the ZIP for your new duty location if you don&apos;t
              have a new address yet.
            </Hint>
            <FormGroup>
              <p>Will you deliver part of your PPM to another place in a different ZIP code?</p>
              <Field
                as={Radio}
                id="hasSecondaryDestinationPostalCodeYes"
                label="Yes"
                name="hasSecondaryDestinationPostalCode"
                value="true"
                checked={values.hasSecondaryDestinationPostalCode === 'true'}
              />
              <Field
                as={Radio}
                id="hasSecondaryDestinationPostalCodeNo"
                label="No"
                name="hasSecondaryDestinationPostalCode"
                value="false"
                checked={values.hasSecondaryDestinationPostalCode === 'false'}
              />
            </FormGroup>
            {values.hasSecondaryDestinationPostalCode === 'true' && (
              <>
                <TextField
                  label="Second ZIP"
                  id="secondaryDestinationPostalCode"
                  name="secondaryDestinationPostalCode"
                  maxLength={10}
                  validate={(value) => postalCodeValidator(value, 'destination')}
                />
                <Hint>
                  <p>A second destination ZIP could mean that your final incentive is lower than your estimate.</p>
                  <p>
                    Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to your
                    move counselor for more detailed information.
                  </p>
                </Hint>
              </>
            )}
            <FormGroup>
              <p>Do you plan to store items from your PPM?</p>
              <Field
                as={Radio}
                id="sitExpectedYes"
                label="Yes"
                name="sitExpected"
                value="true"
                checked={values.sitExpected === 'true'}
              />
              <Field
                as={Radio}
                id="sitExpectedNo"
                label="No"
                name="sitExpected"
                value="false"
                checked={values.sitExpected === 'false'}
              />
            </FormGroup>
            {values.sitExpected === 'false' ? (
              <Hint>You can be reimbursed for up to 90 days of temporary storage (SIT).</Hint>
            ) : (
              <Hint>
                <p>You can be reimbursed for up to 90 days of temporary storage (SIT).</p>
                <p>
                  Your reimbursement amount is limited to the Government&apos;s Constructed Cost — what the government
                  would have paid to store your belongings.
                </p>
                <p>
                  You will need to pay for the storage yourself, then submit receipts and request reimbursement after
                  your PPM is complete.
                </p>
                <p>Your move counselor can give you more information about additional requirements.</p>
              </Hint>
            )}
            <DatePickerInput name="expectedDepartureDate" label="When do you plan to start moving your PPM?" />
            <Hint>
              Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We will ask
              for your actual departure date when you document and complete your PPM.
            </Hint>
            <Button type="button" unstyled onClick={onBack} data-testid="datesAndLocationBackBtn">
              Back
            </Button>
            <Button
              type="button"
              onClick={handleSubmit}
              unstyled
              data-testid="datesAndLocationSubmitBtn"
              disabled={!isValid || isSubmitting}
            >
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
  destinationDutyStation: DutyStationShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  postalCodeValidator: func.isRequired,
};

DatesAndLocation.defaultProps = {
  mtoShipment: undefined,
};

export default DatesAndLocation;
