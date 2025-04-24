import { React, useEffect, useState } from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Checkbox, Radio, FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { DatePickerInput, DutyLocationInput } from 'components/form/fields';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types';
import { MoveShape, ServiceMemberShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { searchTransportationOffices } from 'services/internalApi';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { OptionalAddressSchema } from 'components/Shared/MtoShipmentForm/validationSchemas';
import { requiredAddressSchema, partialRequiredAddressSchema } from 'utils/validation';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import RequiredTag from 'components/form/RequiredTag';
import { isPreceedingAddressComplete, isPreceedingAddressPPMPrimaryDestinationComplete } from 'shared/utils';
import { handleAddressToggleChange, blankAddress } from 'utils/shipments';
import DebounceButton from 'components/DebouncedButton/DebounceButton';

let meta = '';

let validationShape = {
  useCurrentResidence: Yup.boolean(),
  hasSecondaryPickupAddress: Yup.boolean(),
  useCurrentDestinationAddress: Yup.boolean(),
  hasSecondaryDestinationAddress: Yup.boolean(),
  sitExpected: Yup.boolean().required('Required'),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  pickupAddress: Yup.object().shape({
    address: requiredAddressSchema,
  }),
  destinationAddress: Yup.object().shape({
    address: partialRequiredAddressSchema,
  }),
  secondaryPickupAddress: Yup.object().shape({
    address: OptionalAddressSchema,
  }),
  secondaryDestinationAddress: Yup.object().shape({
    address: OptionalAddressSchema,
  }),
  tertiaryPickupAddress: Yup.object().shape({
    address: OptionalAddressSchema,
  }),
  tertiaryDestinationAddress: Yup.object().shape({
    address: OptionalAddressSchema,
  }),
};

const DateAndLocationForm = ({ mtoShipment, destinationDutyLocation, serviceMember, move, onBack, onSubmit }) => {
  const initialValues = {
    useCurrentResidence: false,
    pickupAddress: {},
    secondaryPickupAddress: {},
    tertiaryPickupAddress: {},
    hasSecondaryPickupAddress: mtoShipment?.ppmShipment?.secondaryPickupAddress ? 'true' : 'false',
    hasTertiaryPickupAddress: mtoShipment?.ppmShipment?.tertiaryPickupAddress ? 'true' : 'false',
    useCurrentDestinationAddress: false,
    destinationAddress: {},
    secondaryDestinationAddress: {},
    tertiaryDestinationAddress: {},
    hasSecondaryDestinationAddress: mtoShipment?.ppmShipment?.secondaryDestinationAddress ? 'true' : 'false',
    hasTertiaryDestinationAddress: mtoShipment?.ppmShipment?.tertiaryDestinationAddress ? 'true' : 'false',
    sitExpected: mtoShipment?.ppmShipment?.sitExpected ? 'true' : 'false',
    expectedDepartureDate: mtoShipment?.ppmShipment?.expectedDepartureDate || '',
    closeoutOffice: move?.closeoutOffice || {},
  };

  if (mtoShipment?.ppmShipment?.pickupAddress) {
    initialValues.pickupAddress = { address: { ...mtoShipment.ppmShipment.pickupAddress } };
  }

  if (mtoShipment?.ppmShipment?.secondaryPickupAddress) {
    initialValues.secondaryPickupAddress = { address: { ...mtoShipment.ppmShipment.secondaryPickupAddress } };
  }

  if (mtoShipment?.ppmShipment?.tertiaryPickupAddress) {
    initialValues.tertiaryPickupAddress = { address: { ...mtoShipment.ppmShipment.tertiaryPickupAddress } };
  }

  if (mtoShipment?.ppmShipment?.destinationAddress) {
    initialValues.destinationAddress = { address: { ...mtoShipment.ppmShipment.destinationAddress } };
  }

  if (mtoShipment?.ppmShipment?.secondaryDestinationAddress) {
    initialValues.secondaryDestinationAddress = { address: { ...mtoShipment.ppmShipment.secondaryDestinationAddress } };
  }

  if (mtoShipment?.ppmShipment?.tertiaryDestinationAddress) {
    initialValues.tertiaryDestinationAddress = { address: { ...mtoShipment.ppmShipment.tertiaryDestinationAddress } };
  }

  const residentialAddress = serviceMember?.residential_address;
  const destinationDutyAddress = destinationDutyLocation?.address;

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled('third_address_available').then((enabled) => {
        setIsTertiaryAddressEnabled(enabled);
      });
    };
    fetchData();
  }, []);

  const showCloseoutOffice =
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.ARMY ||
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.AIR_FORCE ||
    serviceMember.affiliation === SERVICE_MEMBER_AGENCIES.SPACE_FORCE;
  if (showCloseoutOffice) {
    validationShape = {
      ...validationShape,
      closeoutOffice: Yup.object().shape({
        address: Yup.object().required('Required'),
      }),
    };
  } else {
    delete validationShape.closeoutOffice;
  }

  const validate = (values) => {
    if (!values.closeoutOffice) {
      meta = 'Required';
    }
    if (values.closeoutOffice) {
      meta = '';
    }
    return {};
  };

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={Yup.object().shape(validationShape)}
      onSubmit={onSubmit}
      validate={validate}
      validateOnBlur
      validateOnMount
      validateOnChange
    >
      {({ isValid, isSubmitting, handleSubmit, setValues, values, errors, ...formikProps }) => {
        const handleUseCurrentResidenceChange = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues({
              ...values,
              pickupAddress: {
                address: residentialAddress,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              pickupAddress: {
                blankAddress,
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
              destinationAddress: {
                address: destinationDutyAddress,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              destinationAddress: {
                blankAddress,
              },
            });
          }
        };

        return (
          <div className={ppmStyles.formContainer}>
            <Form className={formStyles.form}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection, 'origin')}>
                <h2>Pickup Address</h2>
                <AddressFields
                  name="pickupAddress.address"
                  labelHint="Required"
                  locationLookup
                  formikProps={formikProps}
                  render={(fields) => (
                    <>
                      <p>What address are you moving from?</p>
                      <Checkbox
                        data-testid="useCurrentResidence"
                        label="Use my current pickup address"
                        name="useCurrentResidence"
                        onChange={handleUseCurrentResidenceChange}
                        id="useCurrentResidence"
                      />
                      {fields}
                      <FormGroup>
                        <Fieldset>
                          <legend className="usa-label">Will you add items to your PPM from a second address?</legend>
                          <RequiredTag />

                          <div className={formStyles.radioGroup}>
                            <Field
                              as={Radio}
                              data-testid="yes-secondary-pickup-address"
                              id="yes-secondary-pickup-address"
                              label="Yes"
                              name="hasSecondaryPickupAddress"
                              value="true"
                              title="Yes, I have a second pickup address"
                              checked={values.hasSecondaryPickupAddress === 'true'}
                              disabled={!isPreceedingAddressComplete('true', values.pickupAddress.address)}
                              onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                            />
                            <Field
                              as={Radio}
                              data-testid="no-secondary-pickup-address"
                              id="no-secondary-pickup-address"
                              label="No"
                              name="hasSecondaryPickupAddress"
                              value="false"
                              title="No, I do not have a second pickup address"
                              checked={values.hasSecondaryPickupAddress === 'false'}
                              disabled={!isPreceedingAddressComplete('true', values.pickupAddress.address)}
                              onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                            />
                          </div>
                        </Fieldset>
                      </FormGroup>
                      {values.hasSecondaryPickupAddress === 'true' && (
                        <>
                          <h4>Second Pickup Address</h4>
                          <AddressFields
                            labelHint="Required"
                            name="secondaryPickupAddress.address"
                            locationLookup
                            formikProps={formikProps}
                          />
                          <Hint className={ppmStyles.hint}>
                            <p>
                              A second pickup address could mean that your final incentive is lower than your estimate.
                            </p>
                            <p>
                              Get separate weight tickets for each leg of the trip to show how the weight changes. Talk
                              to your move counselor for more detailed information.
                            </p>
                          </Hint>
                        </>
                      )}

                      {isTertiaryAddressEnabled && values.hasSecondaryPickupAddress === 'true' && (
                        <div>
                          <FormGroup>
                            <legend className="usa-label">Will you add items to your PPM from a third address?</legend>
                            <RequiredTag />
                            <Fieldset>
                              <div className={formStyles.radioGroup}>
                                <Field
                                  as={Radio}
                                  id="yes-tertiary-pickup-address"
                                  data-testid="yes-tertiary-pickup-address"
                                  label="Yes"
                                  name="hasTertiaryPickupAddress"
                                  value="true"
                                  title="Yes, I have a third pickup address"
                                  checked={values.hasTertiaryPickupAddress === 'true'}
                                  disabled={
                                    !isPreceedingAddressComplete(
                                      values.hasSecondaryPickupAddress,
                                      values.secondaryPickupAddress.address,
                                    )
                                  }
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                                <Field
                                  as={Radio}
                                  id="no-tertiary-pickup-address"
                                  data-testid="no-tertiary-pickup-address"
                                  label="No"
                                  name="hasTertiaryPickupAddress"
                                  value="false"
                                  title="No, I do not have a third pickup address"
                                  checked={values.hasTertiaryPickupAddress === 'false'}
                                  disabled={
                                    !isPreceedingAddressComplete(
                                      values.hasSecondaryPickupAddress,
                                      values.secondaryPickupAddress.address,
                                    )
                                  }
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                              </div>
                            </Fieldset>
                          </FormGroup>
                        </div>
                      )}
                      {isTertiaryAddressEnabled &&
                        values.hasSecondaryPickupAddress === 'true' &&
                        values.hasTertiaryPickupAddress === 'true' && (
                          <>
                            <h4>Third Pickup Address</h4>
                            <AddressFields
                              labelHint="Required"
                              name="tertiaryPickupAddress.address"
                              locationLookup
                              formikProps={formikProps}
                            />
                          </>
                        )}
                    </>
                  )}
                />
              </SectionWrapper>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Delivery Address</h2>
                <AddressFields
                  name="destinationAddress.address"
                  labelHint="Required"
                  locationLookup
                  formikProps={formikProps}
                  // White spaces are used specifically to override incoming labelHint prop
                  // not to display anything.
                  address1LabelHint=" "
                  render={(fields) => (
                    <>
                      <p>Please input your delivery address.</p>
                      <Checkbox
                        data-testid="useCurrentDestinationAddress"
                        label="Use my current delivery address"
                        name="useCurrentDestinationAddress"
                        onChange={handleUseDestinationAddress}
                        id="useCurrentDestinationAddress"
                      />
                      {fields}
                      <FormGroup>
                        <div className={formStyles.radioGroup}>
                          <Fieldset>
                            <legend className="usa-label">
                              Will you deliver part of your PPM to a second address?
                            </legend>
                            <RequiredTag />
                            <div className={formStyles.radioGroup}>
                              <Field
                                as={Radio}
                                data-testid="yes-secondary-destination-address"
                                id="hasSecondaryDestinationAddressYes"
                                label="Yes"
                                name="hasSecondaryDestinationAddress"
                                value="true"
                                title="Yes, I have a second delivery address"
                                checked={values.hasSecondaryDestinationAddress === 'true'}
                                disabled={
                                  !isPreceedingAddressPPMPrimaryDestinationComplete(values.destinationAddress.address)
                                }
                                onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                              />
                              <Field
                                as={Radio}
                                data-testid="no-secondary-destination-address"
                                id="hasSecondaryDestinationAddressNo"
                                label="No"
                                name="hasSecondaryDestinationAddress"
                                value="false"
                                title="No, I do not have a second delivery address"
                                checked={values.hasSecondaryDestinationAddress === 'false'}
                                disabled={
                                  !isPreceedingAddressPPMPrimaryDestinationComplete(values.destinationAddress.address)
                                }
                                onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                              />
                            </div>
                          </Fieldset>
                        </div>
                      </FormGroup>
                      {values.hasSecondaryDestinationAddress === 'true' && (
                        <>
                          <h4>Second Delivery Address</h4>
                          <AddressFields
                            name="secondaryDestinationAddress.address"
                            labelHint="Required"
                            locationLookup
                            formikProps={formikProps}
                          />
                          <Hint className={ppmStyles.hint}>
                            <p>
                              A second delivery address could mean that your final incentive is lower than your
                              estimate.
                            </p>
                            <p>
                              Get separate weight tickets for each leg of the trip to show how the weight changes. Talk
                              to your move counselor for more detailed information.
                            </p>
                          </Hint>
                        </>
                      )}

                      {isTertiaryAddressEnabled && values.hasSecondaryDestinationAddress === 'true' && (
                        <div>
                          <FormGroup>
                            <legend className="usa-label">Will you deliver part of your PPM to a third address?</legend>
                            <RequiredTag />
                            <Fieldset>
                              <div className={formStyles.radioGroup}>
                                <Field
                                  as={Radio}
                                  id="has-tertiary-delivery"
                                  data-testid="yes-tertiary-destination-address"
                                  label="Yes"
                                  name="hasTertiaryDestinationAddress"
                                  value="true"
                                  title="Yes, I have a third delivery address"
                                  checked={values.hasTertiaryDestinationAddress === 'true'}
                                  disabled={
                                    !isPreceedingAddressComplete(
                                      values.hasSecondaryDestinationAddress,
                                      values.secondaryDestinationAddress.address,
                                    )
                                  }
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                                <Field
                                  as={Radio}
                                  id="no-tertiary-delivery"
                                  data-testid="no-tertiary-destination-address"
                                  label="No"
                                  name="hasTertiaryDestinationAddress"
                                  value="false"
                                  title="No, I do not have a third delivery address"
                                  checked={values.hasTertiaryDestinationAddress === 'false'}
                                  disabled={
                                    !isPreceedingAddressComplete(
                                      values.hasSecondaryDestinationAddress,
                                      values.secondaryDestinationAddress.address,
                                    )
                                  }
                                  onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                />
                              </div>
                            </Fieldset>
                          </FormGroup>
                        </div>
                      )}
                      {isTertiaryAddressEnabled &&
                        values.hasSecondaryDestinationAddress === 'true' &&
                        values.hasTertiaryDestinationAddress === 'true' && (
                          <>
                            <h4>Third Delivery Address</h4>
                            <AddressFields
                              name="tertiaryDestinationAddress.address"
                              labelHint="Required"
                              locationLookup
                              formikProps={formikProps}
                            />
                          </>
                        )}
                    </>
                  )}
                />
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
                      hint="Required"
                      placeholder="Start typing a closeout office..."
                      searchLocations={searchTransportationOffices}
                      metaOverride={meta}
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
                  <RequiredTag />
                  <Field
                    as={Radio}
                    id="sitExpectedYes"
                    data-testid="storePPMYes"
                    label="Yes"
                    name="sitExpected"
                    value="true"
                    checked={values.sitExpected === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="sitExpectedNo"
                    data-testid="storePPMNo"
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
                <DatePickerInput
                  hint="Required"
                  name="expectedDepartureDate"
                  label="When do you plan to start moving your PPM?"
                />
                <Hint className={ppmStyles.hint}>
                  Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We will
                  ask for your actual departure date when you document and complete your PPM.
                </Hint>
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <DebounceButton
                  className={ppmStyles.saveButton}
                  type="button"
                  onClick={handleSubmit}
                  disabled={isSubmitting || !isValid}
                >
                  Save & Continue
                </DebounceButton>
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
};

DateAndLocationForm.defaultProps = {
  mtoShipment: undefined,
  move: undefined,
};

export default DateAndLocationForm;
