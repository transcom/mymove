import React, { useEffect, useState } from 'react';
import { GridContainer, Grid, Alert, Label, Radio, Fieldset } from '@trussworks/react-uswds';
import { generatePath, useNavigate } from 'react-router-dom';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { connect } from 'react-redux';

import { statesList } from '../../../constants/states';

import styles from './CreateCustomerForm.module.scss';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { servicesCounselingRoutes } from 'constants/routes';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { CheckboxField, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { backupContactInfoSchema, requiredAddressSchema } from 'utils/validation';
import { createCustomerWithOktaOption } from 'services/ghcApi';
import { getResponseError } from 'services/internalApi';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import departmentIndicators from 'constants/departmentIndicators';

export const CreateCustomerForm = ({ userPrivileges, setFlashMessage }) => {
  const [serverError, setServerError] = useState(null);
  const [showEmplid, setShowEmplid] = useState(false);
  const [isSafetyMove, setIsSafetyMove] = useState(false);
  const navigate = useNavigate();

  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);

  const residentialAddressName = 'residential_address';
  const backupAddressName = 'backup_mailing_address';
  const backupContactName = 'backup_contact';

  const [isSafetyMoveFF, setSafetyMoveFF] = useState(false);
  const [secondaryTelephoneNum, setSecondaryTelephoneNum] = useState('');

  useEffect(() => {
    isBooleanFlagEnabled('safety_move')?.then((enabled) => {
      setSafetyMoveFF(enabled);
    });
  }, []);

  const isSafetyPrivileged = isSafetyMoveFF
    ? userPrivileges?.some((privilege) => privilege.privilegeType === elevatedPrivilegeTypes.SAFETY)
    : false;

  const initialValues = {
    affiliation: '',
    edipi: '',
    first_name: '',
    middle_name: '',
    last_name: '',
    suffix: '',
    telephone: '',
    secondary_telephone: null,
    personal_email: '',
    phone_is_preferred: false,
    email_is_preferred: false,
    [residentialAddressName]: {
      streetAddress1: '',
      streetAddress2: '',
      streetAddress3: '',
      city: '',
      state: '',
      postalCode: '',
    },
    [backupAddressName]: {
      streetAddress1: '',
      streetAddress2: '',
      streetAddress3: '',
      city: '',
      state: '',
      postalCode: '',
    },
    [backupContactName]: {
      name: '',
      telephone: '',
      email: '',
    },
    create_okta_account: '',
    cac_user: '',
    is_safety_move: false,
  };

  const handleBack = () => {
    navigate(servicesCounselingRoutes.BASE_CUSTOMER_SEARCH_PATH);
  };

  const onSubmit = async (values) => {
    // Convert strings to booleans to satisfy swagger
    const createOktaAccount = values.create_okta_account === 'true';
    const cacUser = values.cac_user === 'true';

    const body = {
      affiliation: values.affiliation,
      edipi: values.edipi,
      emplid: values.emplid || '',
      firstName: values.first_name,
      middleName: values.middle_name,
      lastName: values.last_name,
      suffix: values.suffix,
      telephone: values.telephone,
      secondaryTelephone: values.secondary_telephone,
      personalEmail: values.personal_email,
      phoneIsPreferred: values.phone_is_preferred,
      emailIsPreferred: values.email_is_preferred,
      residentialAddress: values[residentialAddressName],
      backupMailingAddress: values[backupAddressName],
      backupContact: {
        name: values[backupContactName].name,
        email: values[backupContactName].email,
        phone: values[backupContactName].telephone,
      },
      createOktaAccount,
      cacUser,
    };

    return createCustomerWithOktaOption({ body })
      .then((res) => {
        const customerId = Object.keys(res.createdCustomer)[0];
        setFlashMessage('CUSTOMER_CREATE_SUCCESS', 'success', `Customer created successfully.`);
        navigate(
          generatePath(servicesCounselingRoutes.BASE_CUSTOMERS_ORDERS_ADD_PATH, {
            customerId,
          }),
          { state: { isSafetyMoveSelected: isSafetyMove } },
        );
      })
      .catch((e) => {
        let errorMessage;
        if (e.status === 409) errorMessage = 'This EMPLID is already in use';
        else errorMessage = getResponseError(e?.response, 'failed to create service member due to server error');
        setServerError(errorMessage);
      });
  };

  const validationSchema = Yup.object().shape({
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: Yup.string().matches(/[0-9]{10}/, 'Enter a 10-digit DOD ID number'),
    emplid: Yup.string()
      .notRequired()
      .matches(/[0-9]{7}/, 'Enter a 7-digit EMPLID number'),
    first_name: Yup.string().required('Required'),
    middle_name: Yup.string(),
    last_name: Yup.string().required('Required'),
    suffix: Yup.string(),
    telephone: Yup.string()
      .min(12, 'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.')
      .required('Required'),
    secondary_telephone: Yup.string()
      .min(12, 'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.')
      .nullable(),
    personal_email: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    phoneIsPreferred: Yup.boolean(),
    emailIsPreferred: Yup.boolean(),
    [residentialAddressName]: requiredAddressSchema.required(),
    [backupAddressName]: requiredAddressSchema.required(),
    [backupContactName]: backupContactInfoSchema.required(),
    create_okta_account: isSafetyMove ? '' : Yup.boolean().required('Required'),
    cac_user: isSafetyMove ? '' : Yup.boolean().required('Required'),
    is_safety_move: isSafetyMoveFF ? Yup.boolean().required('Required') : '',
  });

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid className={styles.nameFormContainer}>
          <Grid col desktop={{ col: 8 }} className={styles.nameForm}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid className={styles.nameFormContainer}>
        <Grid col desktop={{ col: 8 }} className={styles.nameForm}>
          <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
            {({ isValid, handleSubmit, setValues, values, handleChange }) => {
              const handleSubmitNext = () => {
                setValues({
                  ...values,
                  secondary_telephone: secondaryTelephoneNum,
                });
                handleSubmit();
              };
              const handlePhoneNumChange = (value) => {
                setSecondaryTelephoneNum(value);
              };
              const handleIsSafetyMove = (e) => {
                const { value } = e.target;
                if (value === 'true') {
                  setIsSafetyMove(true);
                  // clear out DoDID, emplid, and OKTA fields
                  setValues({
                    ...values,
                    edipi: '',
                    emplid: '',
                    create_okta_account: '',
                    cac_user: 'true',
                    is_safety_move: 'true',
                  });
                } else if (value === 'false') {
                  setIsSafetyMove(false);
                  setValues({
                    ...values,
                    is_safety_move: 'false',
                  });
                }
              };
              const handleBranchChange = (e) => {
                if (e.target.value === departmentIndicators.COAST_GUARD) {
                  setShowEmplid(true);
                } else {
                  setShowEmplid(false);
                }
              };
              return (
                <Form className={formStyles.form}>
                  <h1 className={styles.header}>Create Customer Profile</h1>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Customer Affiliation</h3>
                    {isSafetyPrivileged && (
                      <Fieldset className={styles.trailerOwnershipFieldset}>
                        <legend className="usa-label">Is this a Safety move?</legend>
                        <div className="grid-row grid-gap">
                          <Field
                            as={Radio}
                            id="isSafetyMoveYes"
                            label="Yes"
                            name="is_safety_move"
                            value="true"
                            data-testid="is-safety-move-yes"
                            onChange={handleIsSafetyMove}
                          />
                          <Field
                            as={Radio}
                            id="isSafetyMoveNo"
                            label="No"
                            name="is_safety_move"
                            value="false"
                            data-testid="is-safety-move-no"
                            onChange={handleIsSafetyMove}
                          />
                        </div>
                      </Fieldset>
                    )}
                    <DropdownInput
                      label="Branch of service"
                      name="affiliation"
                      id="affiliation"
                      required
                      onChange={(e) => {
                        handleChange(e);
                        handleBranchChange(e);
                      }}
                      options={branchOptions}
                    />
                    <TextField
                      label="DoD ID number"
                      name="edipi"
                      id="edipi"
                      labelHint="Optional"
                      maxLength="10"
                      isDisabled={isSafetyMove}
                    />
                    {showEmplid && (
                      <TextField
                        label="EMPLID"
                        name="emplid"
                        id="emplid"
                        maxLength="7"
                        labelHint="Optional"
                        inputMode="numeric"
                        pattern="[0-9]{7}"
                        isDisabled={isSafetyMove}
                      />
                    )}
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Customer Name</h3>
                    <TextField label="First name" name="first_name" id="firstName" required />
                    <TextField label="Middle name" name="middle_name" id="middleName" labelHint="Optional" />
                    <TextField label="Last name" name="last_name" id="lastName" required />
                    <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Contact Info</h3>
                    <MaskedTextField
                      label="Best contact phone"
                      id="telephone"
                      name="telephone"
                      type="tel"
                      minimum="12"
                      mask="000{-}000{-}0000"
                      required
                    />
                    <MaskedTextField
                      label="Alt. phone"
                      labelHint="Optional"
                      id="altTelephone"
                      name="secondary_telephone"
                      type="tel"
                      minimum="12"
                      mask="000{-}000{-}0000"
                      onChange={handlePhoneNumChange}
                    />
                    <TextField label="Personal email" id="personalEmail" name="personal_email" required />
                    <Label>Preferred contact method (optional)</Label>
                    <div className={formStyles.radioGroup}>
                      <CheckboxField id="phoneIsPreferred" label="Phone" name="phone_is_preferred" />
                      <CheckboxField id="emailIsPreferred" label="Email" name="email_is_preferred" />
                    </div>
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Current Address</h3>
                    <TextField
                      label="Address 1"
                      id="mailingAddress1"
                      name="residential_address.streetAddress1"
                      data-testid="res-add-street1"
                    />
                    <TextField
                      label="Address 2"
                      labelHint="Optional"
                      id="mailingAddress2"
                      name="residential_address.streetAddress2"
                      data-testid="res-add-street2"
                    />
                    <TextField
                      label="Address 3"
                      labelHint="Optional"
                      id="mailingAddress3"
                      name="residential_address.streetAddress3"
                      data-testid="res-add-street3"
                    />
                    <TextField label="City" id="city" name="residential_address.city" data-testid="res-add-city" />

                    <div className="grid-row grid-gap">
                      <div className="mobile-lg:grid-col-6">
                        <DropdownInput
                          name="residential_address.state"
                          id="state"
                          label="State"
                          options={statesList}
                          data-testid="res-add-state"
                        />
                      </div>
                      <div className="mobile-lg:grid-col-6">
                        <TextField
                          label="ZIP"
                          id="zip"
                          name="residential_address.postalCode"
                          maxLength={10}
                          data-testid="res-add-zip"
                        />
                      </div>
                    </div>
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Backup Address</h3>
                    <TextField
                      label="Address 1"
                      id="backupMailingAddress1"
                      name="backup_mailing_address.streetAddress1"
                      data-testid="backup-add-street1"
                    />
                    <TextField
                      label="Address 2"
                      labelHint="Optional"
                      id="backupMailingAddress2"
                      name="backup_mailing_address.streetAddress2"
                      data-testid="backup-add-street2"
                    />
                    <TextField
                      label="Address 3"
                      labelHint="Optional"
                      id="backupMailingAddress3"
                      name="backup_mailing_address.streetAddress3"
                      data-testid="backup-add-street3"
                    />
                    <TextField
                      label="City"
                      id="backupCity"
                      name="backup_mailing_address.city"
                      data-testid="backup-add-city"
                    />

                    <div className="grid-row grid-gap">
                      <div className="mobile-lg:grid-col-6">
                        <DropdownInput
                          name="backup_mailing_address.state"
                          id="backupState"
                          label="State"
                          options={statesList}
                          data-testid="backup-add-state"
                        />
                      </div>
                      <div className="mobile-lg:grid-col-6">
                        <TextField
                          label="ZIP"
                          id="backupZip"
                          name="backup_mailing_address.postalCode"
                          maxLength={10}
                          data-testid="backup-add-zip"
                        />
                      </div>
                    </div>
                  </SectionWrapper>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Backup Contact</h3>
                    <TextField label="Name" id="backupContactName" name="backup_contact.name" required />
                    <TextField label="Email" id="backupContactEmail" name="backup_contact.email" required />
                    <MaskedTextField
                      label="Phone"
                      id="backupContactTelephone"
                      name="backup_contact.telephone"
                      type="tel"
                      minimum="12"
                      mask="000{-}000{-}0000"
                      required
                    />
                  </SectionWrapper>
                  {values.is_safety_move !== 'true' && (
                    <SectionWrapper className={formStyles.formSection}>
                      <h3>Okta Account</h3>
                      <Fieldset className={styles.trailerOwnershipFieldset}>
                        <legend className="usa-label">Do you want to create an Okta account for this customer?</legend>
                        <div className="grid-row grid-gap">
                          <Field
                            as={Radio}
                            id="yesCreateOktaAccount"
                            label="Yes"
                            name="create_okta_account"
                            value="true"
                            data-testid="create-okta-account-yes"
                          />
                          <Field
                            as={Radio}
                            id="noCreateOktaAccount"
                            label="No"
                            name="create_okta_account"
                            value="false"
                            data-testid="create-okta-account-no"
                          />
                        </div>
                      </Fieldset>
                    </SectionWrapper>
                  )}
                  {values.is_safety_move !== 'true' && (
                    <SectionWrapper className={formStyles.formSection}>
                      <h3>Non-CAC Users</h3>
                      <Fieldset className={styles.trailerOwnershipFieldset}>
                        <legend className="usa-label">Does the customer have a CAC?</legend>
                        <div className="grid-row grid-gap">
                          <Field
                            as={Radio}
                            id="yesCacUser"
                            label="Yes"
                            name="cac_user"
                            value="true"
                            data-testid="cac-user-yes"
                          />
                          <Field
                            as={Radio}
                            id="NonCacUser"
                            label="No"
                            name="cac_user"
                            value="false"
                            data-testid="cac-user-no"
                          />
                        </div>
                      </Fieldset>
                    </SectionWrapper>
                  )}
                  <div className={formStyles.formActions}>
                    <WizardNavigation
                      editMode
                      onCancelClick={handleBack}
                      disableNext={!isValid}
                      onNextClick={handleSubmitNext}
                    />
                  </div>
                </Form>
              );
            }}
          </Formik>
        </Grid>
      </Grid>
    </GridContainer>
  );
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreateCustomerForm);
