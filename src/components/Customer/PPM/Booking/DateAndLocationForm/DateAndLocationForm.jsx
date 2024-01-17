import React, { useState } from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Checkbox, Radio, FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { CheckboxField, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { UnsupportedZipCodePPMErrorMsg, ZIP5_CODE_REGEX, InvalidZIPTypeError } from 'utils/validation';
import { searchTransportationOffices } from 'services/internalApi';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { AddressFields } from 'components/form/AddressFields/AddressFields';

// export const residentialAddressName = 'residential_address';

const validationShape = {
  pickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  useCurrentResidence: Yup.boolean(),
  hasSecondaryPickup: Yup.boolean(),
  useCurrentDestinationAddress: Yup.boolean(),
  // secondaryPickupPostalCode: Yup.string().when('hasSecondaryPickup', {
  //   is: true,
  //   then: (schema) => schema.matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  // }),
  // useDestinationDutyLocationZIP: Yup.boolean(),
  // destinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  // hasSecondaryDestinationPostalCode: Yup.boolean().required('Required'),
  // secondaryDestinationPostalCode: Yup.string().when('hasSecondaryDestinationPostalCode', {
  //  is: true,
  //  then: (schema) => schema.matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  // }),
  sitExpected: Yup.boolean().required('Required'),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
};
const setZip = (setFieldValue, postalCodeField, postalCode, isChecked, isCheckedField) => {
  setFieldValue(isCheckedField, !isChecked);
  setFieldValue(postalCodeField, isChecked ? '' : postalCode);
};

const DateAndLocationForm = ({
  mtoShipment,
  destinationDutyLocation,
  serviceMember,
  move,
  onBack,
  onSubmit,
  postalCodeValidator,
}) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});

  const initialValues = {
    pickupPostalCode: mtoShipment?.ppmShipment?.pickupPostalCode || '',
    useCurrentResidence: false,
    secondaryPickupAddress: mtoShipment?.ppmShipment?.secondaryPickupAddress || '',
    hasSecondaryPickup: mtoShipment?.ppmShipment?.hasSecondaryPickup ? 'true' : 'false',
    secondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || '',
    useCurrentDestinationAddress: false,
    destinationPostalCode: mtoShipment?.ppmShipment?.destinationPostalCode || '',
    hasSecondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode ? 'true' : 'false',
    secondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || '',
    sitExpected: mtoShipment?.ppmShipment?.sitExpected ? 'true' : 'false',
    expectedDepartureDate: mtoShipment?.ppmShipment?.expectedDepartureDate || '',
    closeoutOffice: move?.closeout_office,
    // [residentialAddressName]: serviceMember?.residential_address,
  };

  const residentialAddress = serviceMember?.residential_address;
  const destinationAddress = destinationDutyLocation?.address;
  // const residentialAddressPostalCode = serviceMember?.residential_address?.postalCode;
  // const destinationDutyLocationPostalCode = destinationDutyLocation?.address?.postalCode;

  const postalCodeValidate = async (value, location, name) => {
    if (value?.length !== 5) {
      return undefined;
    }
    // only revalidate if the value has changed, editing other fields will re-validate unchanged ones
    if (postalCodeValid[`${name}`]?.value !== value) {
      const response = await postalCodeValidator(value, location, UnsupportedZipCodePPMErrorMsg);
      setPostalCodeValid((state) => {
        return {
          ...state,
          [name]: { value, isValid: !response },
        };
      });
      return response;
    }
    return postalCodeValid[`${name}`]?.isValid ? undefined : UnsupportedZipCodePPMErrorMsg;
  };

  const handlePrefillPostalCodeChange = (
    value,
    setFieldValue,
    postalCodeField,
    prefillValue,
    isCheckedField,
    checkedFieldValue,
  ) => {
    if (checkedFieldValue && value !== prefillValue) {
      setFieldValue(isCheckedField, false);
    }
    setFieldValue(postalCodeField, value);
  };

  const showCloseoutOffice =
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.ARMY ||
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.AIR_FORCE;
  if (showCloseoutOffice) {
    validationShape.closeoutOffice = Yup.object().required('Required');
  } else {
    delete validationShape.closeoutOffice;
  }

  return (
    <Formik initialValues={initialValues} validationSchema={Yup.object().shape(validationShape)} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, setFieldValue, setValues, values }) => {
        const handleUseCurrentResidenceChange = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues({
              ...values,
              serviceMember: {
                ...values.serviceMember,
                residential_address: residentialAddress,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              serviceMember: {
                ...values.serviceMember,
                residential_address: {
                  streetAddress1: '',
                  streetAddress2: '',
                  city: '',
                  state: '',
                  postalCode: '',
                },
              },
            });
          }
        };

        const handleUseDestinationAddress = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues({
              ...values,
              serviceMember: {
                ...values.serviceMember,
                destination_address: destinationAddress,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              serviceMember: {
                ...values.serviceMember,
                destination_address: {
                  streetAddress1: '',
                  streetAddress2: '',
                  city: '',
                  state: '',
                  postalCode: '',
                },
              },
            });
          }
        };
        return (
          <div className={ppmStyles.formContainer}>
            <Form className={(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection, 'origin')}>
                <h2>Origin</h2>
                <AddressFields
                  name="serviceMember.residential_address"
                  render={(fields) => (
                    <>
                      <p>What address are the movers picking up from?</p>
                      <Checkbox
                        data-testid="useCurrentResidence"
                        label="Use current address"
                        name="serviceMember.residential_address"
                        onChange={handleUseCurrentResidenceChange}
                        id="useCurrentResidence"
                      />
                      {fields}
                      {/* <h4>Second Origin Address</h4>
                      <FormGroup>
                        <p>Will you add items to your PPM from a different address?</p>
                        <div className={formStyles.radioGroup}>
                          <Field
                            as={Radio}
                            id="has-secondary-pickup"
                            data-testid="has-secondary-pickup"
                            label="Yes"
                            name="hasSecondaryPickup"
                            value="yes"
                            title="Yes, I have a second pickup location"
                            checked={hasSecondaryPickup === 'yes'}
                          />
                          <Field
                            as={Radio}
                            id="no-secondary-pickup"
                            data-testid="no-secondary-pickup"
                            label="No"
                            name="hasSecondaryPickup"
                            value="no"
                            title="No, I do not have a second pickup location"
                            checked={hasSecondaryPickup !== 'yes'}
                          />
                        </div>
                      </FormGroup>
                      {hasSecondaryPickup === 'yes' && <AddressFields name="secondaryPickup.address" />} */}
                    </>
                  )}
                />

                {/* <TextField
                  label="ZIP"
                  id="pickupPostalCode"
                  name="pickupPostalCode"
                  maxLength={5}
                  onChange={(e) => {
                    handlePrefillPostalCodeChange(
                      e.target.value,
                      setFieldValue,
                      'pickupPostalCode',
                      residentialAddressPostalCode,
                      'useResidentialAddressZIP',
                      values.useResidentialAddressZIP,
                    );
                  }}
                  validate={(value) => postalCodeValidate(value, 'origin', 'pickupPostalCode')}
                /> */}
                {/* <CheckboxField
                  id="useResidentialAddressZIP"
                  name="useResidentialAddressZIP"
                  label={`Use my current ZIP (${residentialAddressPostalCode})`}
                  onChange={() =>
                    setValue(
                      setFieldValue,
                      'pickupPostalCode',
                      residentialAddressPostalCode,
                      values.useResidentialAddressZIP,
                      'useResidentialAddressZIP',
                    )
                  }
                /> */}
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">
                      Will you add items to your PPM from a place in a different ZIP code?
                    </legend>
                    <Field
                      as={Radio}
                      data-testid="yes-secondary-pickup-postal-code"
                      id="yes-secondary-pickup-postal-code"
                      label="Yes"
                      name="hasSecondaryPickup"
                      value="true"
                      checked={values.hasSecondaryPickup === 'true'}
                    />
                    <Field
                      as={Radio}
                      data-testid="no-secondary-pickup-postal-code"
                      id="no-secondary-pickup-postal-code"
                      label="No"
                      name="hasSecondaryPickup"
                      value="false"
                      checked={values.hasSecondaryPickup === 'false'}
                    />
                  </Fieldset>
                </FormGroup>
                {values.hasSecondaryPickup === 'true' && (
                  <>
                    <AddressFields name="mtoShipment.secondaryPickupAddress" />
                    <Hint className={ppmStyles.hint}>
                      <p>A second origin address could mean that your final incentive is lower than your estimate.</p>
                      <p>
                        Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to
                        your move counselor for more detailed information.
                      </p>
                    </Hint>
                  </>
                )}
              </SectionWrapper>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Destination</h2>
                <AddressFields
                  name="serviceMember.destination_address"
                  render={(fields) => (
                    <>
                      <p>Please input Delivery Address</p>
                      <Checkbox
                        data-testid="useCurrentDestinationAddress"
                        label="Use my current destination address"
                        name="serviceMember.destination_address"
                        onChange={handleUseDestinationAddress}
                        id="useCurrentDestinationAddress"
                      />
                      {fields}
                      {/* <TextField
                  label="ZIP"
                  id="destinationPostalCode"
                  name="destinationPostalCode"
                  maxLength={5}
                  onChange={(e) => {
                    handlePrefillPostalCodeChange(
                      e.target.value,
                      setFieldValue,
                      'destinationPostalCode',
                      destinationDutyLocationPostalCode,
                      'useDestinationDutyLocationZIP',
                      values.useDestinationDutyLocationZIP,
                    );
                  }}
                  validate={(value) => postalCodeValidate(value, 'destination', 'destinationPostalCode')}
                /> */}
                    </>
                  )}
                />
                {/* <CheckboxField
                  id="useDestinationDutyLocationZIP"
                  name="useDestinationDutyLocationZIP"
                  label={`Use the ZIP for my new duty location (${destinationDutyLocationPostalCode})`}
                  onChange={() =>
                    setZip(
                      setFieldValue,
                      'destinationPostalCode',
                      destinationDutyLocationPostalCode,
                      values.useDestinationDutyLocationZIP,
                      'useDestinationDutyLocationZIP',
                    )
                  }
                /> */}
                {/* <Hint className={ppmStyles.hint}>
                  Use the ZIP for your new address if you know it. Use the ZIP for your new duty location if you
                  don&apos;t have a new address yet.
                </Hint> */}
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Will you add items to your PPM from a different address?</legend>
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
                  </Fieldset>
                </FormGroup>
                {values.hasSecondaryDestinationPostalCode === 'true' && (
                  <>
                    <AddressFields name="serviceMember.backup_mailing_address" />
                    {/* <TextField
                      label="Second ZIP"
                      id="secondaryDestinationPostalCode"
                      name="secondaryDestinationPostalCode"
                      maxLength={5}
                      validate={(value) => postalCodeValidate(value, 'destination', 'secondaryDestinationPostalCode')}
                    /> */}
                    <Hint className={ppmStyles.hint}>
                      <p>A second destination ZIP could mean that your final incentive is lower than your estimate.</p>
                      <p>
                        Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to
                        your move counselor for more detailed information.
                      </p>
                    </Hint>
                  </>
                )}
              </SectionWrapper>
              {showCloseoutOffice && (
                <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                  <h2>Closeout Office</h2>
                  <Fieldset>
                    <Hint className={ppmStyles.hint}>
                      <p>
                        A closeout office is where your PPM paperwork will be reviewed before you can submit it to
                        finance to receive your incentive. This will typically be your destination installation&apos;s
                        transportation office or an installation near your destination. If you are not sure what to
                        select, contact your origin transportation office.
                      </p>
                    </Hint>
                    <DutyLocationInput
                      name="closeoutOffice"
                      label="Which closeout office should review your PPM?"
                      placeholder="Start typing a closeout office..."
                      searchLocations={searchTransportationOffices}
                    />
                    <Hint className={ppmStyles.hint}>
                      If you have more than one PPM for this move, your closeout office will be the same for all your
                      PPMs.
                    </Hint>
                  </Fieldset>
                </SectionWrapper>
              )}
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Storage</h2>
                <Fieldset>
                  <legend className="usa-label">Do you plan to store items from your PPM?</legend>
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
                </Fieldset>
                {values.sitExpected === 'false' ? (
                  <Hint className={ppmStyles.hint}>
                    You can be reimbursed for up to 90 days of temporary storage (SIT).
                  </Hint>
                ) : (
                  <Hint>
                    <p>You can be reimbursed for up to 90 days of temporary storage (SIT).</p>
                    <p>
                      Your reimbursement amount is limited to the Government&apos;s Constructed Cost — what the
                      government would have paid to store your belongings.
                    </p>
                    <p>
                      You will need to pay for the storage yourself, then submit receipts and request reimbursement
                      after your PPM is complete.
                    </p>
                    <p>Your move counselor can give you more information about additional requirements.</p>
                  </Hint>
                )}
              </SectionWrapper>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Departure date</h2>
                <DatePickerInput name="expectedDepartureDate" label="When do you plan to start moving your PPM?" />
                <Hint className={ppmStyles.hint}>
                  Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We will
                  ask for your actual departure date when you document and complete your PPM.
                </Hint>
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <Button
                  className={ppmStyles.saveButton}
                  type="button"
                  onClick={handleSubmit}
                  disabled={!isValid || isSubmitting}
                >
                  Save & Continue
                </Button>
              </div>
            </Form>
          </div>
        );
      }}
    </Formik>
  );
};

DateAndLocationForm.propTypes = {
  mtoShipment: ShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  move: MoveShape,
  destinationDutyLocation: DutyLocationShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  postalCodeValidator: func.isRequired,
};

DateAndLocationForm.defaultProps = {
  mtoShipment: undefined,
  move: undefined,
};

export default DateAndLocationForm;
