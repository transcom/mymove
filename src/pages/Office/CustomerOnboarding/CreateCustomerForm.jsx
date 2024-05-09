import React, { useState } from 'react';
import { GridContainer, Grid, Alert, Label, Radio, Fieldset } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router-dom';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { connect } from 'react-redux';

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
import { roleTypes } from 'constants/userRoles';

export const CreateCustomerForm = ({ roleType, setFlashMessage }) => {
  const [serverError, setServerError] = useState(null);
  const navigate = useNavigate();

  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const statesList = [
    { value: 'AL', key: 'AL' },
    { value: 'AK', key: 'AK' },
    { value: 'AR', key: 'AR' },
    { value: 'AZ', key: 'AZ' },
    { value: 'CA', key: 'CA' },
    { value: 'CO', key: 'CO' },
    { value: 'CT', key: 'CT' },
    { value: 'DC', key: 'DC' },
    { value: 'DE', key: 'DE' },
    { value: 'FL', key: 'FL' },
    { value: 'GA', key: 'GA' },
    { value: 'HI', key: 'HI' },
    { value: 'IA', key: 'IA' },
    { value: 'ID', key: 'ID' },
    { value: 'IL', key: 'IL' },
    { value: 'IN', key: 'IN' },
    { value: 'KS', key: 'KS' },
    { value: 'KY', key: 'KY' },
    { value: 'LA', key: 'LA' },
    { value: 'MA', key: 'MA' },
    { value: 'MD', key: 'MD' },
    { value: 'ME', key: 'ME' },
    { value: 'MI', key: 'MI' },
    { value: 'MN', key: 'MN' },
    { value: 'MO', key: 'MO' },
    { value: 'MS', key: 'MS' },
    { value: 'MT', key: 'MT' },
    { value: 'NC', key: 'NC' },
    { value: 'ND', key: 'ND' },
    { value: 'NE', key: 'NE' },
    { value: 'NH', key: 'NH' },
    { value: 'NJ', key: 'NJ' },
    { value: 'NM', key: 'NM' },
    { value: 'NV', key: 'NV' },
    { value: 'NY', key: 'NY' },
    { value: 'OH', key: 'OH' },
    { value: 'OK', key: 'OK' },
    { value: 'OR', key: 'OR' },
    { value: 'PA', key: 'PA' },
    { value: 'RI', key: 'RI' },
    { value: 'SC', key: 'SC' },
    { value: 'SD', key: 'SD' },
    { value: 'TN', key: 'TN' },
    { value: 'TX', key: 'TX' },
    { value: 'UT', key: 'UT' },
    { value: 'VA', key: 'VA' },
    { value: 'VT', key: 'VT' },
    { value: 'WA', key: 'WA' },
    { value: 'WI', key: 'WI' },
    { value: 'WV', key: 'WV' },
    { value: 'WY', key: 'WY' },
  ];

  const residentialAddressName = 'residential_address';
  const backupAddressName = 'backup_mailing_address';
  const backupContactName = 'backup_contact';

  const initialValues = {
    isSafetyMove: '',
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
  };

  const handleBack = () => {
    navigate(servicesCounselingRoutes.BASE_CUSTOMER_SEARCH_PATH);
  };

  const onSubmit = async (values) => {
    // Convert strings to booleans to satisfy swagger
    const createSafetyMove = values.isSafetyMove === 'true';
    const createOktaAccount = values.create_okta_account === 'true';

    const body = {
      affiliation: values.affiliation,
      edipi: values.edipi,
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
    };

    return createCustomerWithOktaOption({ body })
      .then(() => {
        setFlashMessage('CUSTOMER_CREATE_SUCCESS', 'success', `Customer created successfully.`);
        navigate(servicesCounselingRoutes.BASE_CUSTOMER_SEARCH_PATH);
      })
      .catch((e) => {
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to create service member due to server error');
        setServerError(errorMessage);
      });
  };

  const validationSchema = Yup.object().shape({
    isSafetyMove: Yup.boolean().required('Required'),
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: Yup.string().matches(/[0-9]{10}/, 'Enter a 10-digit DOD ID number'),
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
    create_okta_account: Yup.boolean().when('isSafetyMove', {
      is: false,
      then: (schema) => schema.required('Required'),
    }),
  });

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid>
          <Grid col desktop={{ col: 8 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid className={styles.nameFormContainer}>
        <Grid col desktop={{ col: 8 }} className={styles.nameForm}>
          <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
            {({ isValid, handleSubmit, setValues, values }) => {
              const handleIsSafetyMove = (e) => {
                const { checked } = e.target;
                if (checked) {
                  // clear out DoDID and OKTA fields
                  setValues({
                    ...values,
                    edipi: '',
                    create_okta_account: '',
                    isSafetyMove: 'true',
                  });
                }
              };
              return (
                <Form className={formStyles.form}>
                  <h1 className={styles.header}>Create Customer Profile</h1>
                  <SectionWrapper className={formStyles.formSection}>
                    <h3>Customer Affiliation</h3>
                    {roleType === roleTypes.SERVICES_COUNSELOR && <Alert>Role is SC</Alert>}
                    <Fieldset className={styles.trailerOwnershipFieldset}>
                      <legend className="usa-label">Is this a Safety Move?</legend>
                      <div className="grid-row grid-gap">
                        <Field
                          as={Radio}
                          id="isSafetyMoveYes"
                          label="Yes"
                          name="isSafetyMove"
                          value="true"
                          checked={values.isSafetyMove === 'true'}
                          onChange={handleIsSafetyMove}
                        />
                        <Field
                          as={Radio}
                          id="isSafetyMoveNo"
                          label="No"
                          name="isSafetyMove"
                          value="false"
                          checked={values.isSafetyMove === 'false'}
                        />
                      </div>
                    </Fieldset>
                    <DropdownInput
                      label="Branch of service"
                      name="affiliation"
                      id="affiliation"
                      required
                      options={branchOptions}
                    />
                    <TextField
                      label="DoD ID number"
                      name="edipi"
                      id="edipi"
                      labelHint="Optional"
                      maxLength="10"
                      isDisabled={values.isSafetyMove === 'true'}
                    />
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
                  {values.isSafetyMove !== 'true' && (
                    <SectionWrapper className={formStyles.formSection} disabled={values.isSafetyMove === 'true'}>
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
                          />
                          <Field
                            as={Radio}
                            id="noCreateOktaAccount"
                            label="No"
                            name="create_okta_account"
                            value="false"
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
  );
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreateCustomerForm);
