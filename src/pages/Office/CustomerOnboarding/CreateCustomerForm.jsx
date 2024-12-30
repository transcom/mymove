import classnames from 'classnames';
import React, { useEffect, useState } from 'react';
import { GridContainer, Grid, Alert, Label, Radio, Fieldset } from '@trussworks/react-uswds';
import { generatePath, useNavigate } from 'react-router-dom';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { connect } from 'react-redux';

import styles from './CreateCustomerForm.module.scss';

import { Form } from 'components/form/Form';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
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
import { generateUniqueDodid, generateUniqueEmplid } from 'utils/customer';
import Hint from 'components/Hint';
import { setCanAddOrders as setCanAddOrdersAction } from 'store/general/actions';

export const CreateCustomerForm = ({ userPrivileges, setFlashMessage, setCanAddOrders }) => {
  const [serverError, setServerError] = useState(null);
  const [showEmplid, setShowEmplid] = useState(false);
  const [isSafetyMove, setIsSafetyMove] = useState(false);
  const [showSafetyMoveHint, setShowSafetyMoveHint] = useState(false);
  const navigate = useNavigate();

  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);

  const residentialAddressName = 'residential_address';
  const backupAddressName = 'backup_mailing_address';
  const backupContactName = 'backup_contact';

  const [isSafetyMoveFF, setSafetyMoveFF] = useState(false);

  const uniqueDodid = generateUniqueDodid();
  const uniqueEmplid = generateUniqueEmplid();

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
    emplid: null,
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
    is_safety_move: 'false',
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
      emplid: values.emplid,
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
        setCanAddOrders(true);
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
    // All branches require an EDIPI unless it is a safety move
    // where a fake DoD ID may be used
    edipi:
      !isSafetyMove &&
      Yup.string()
        .matches(/^(SM[0-9]{8}|[0-9]{10})$/, 'Enter a 10-digit DoD ID number')
        .required('Required'),
    // Only the coast guard requires both EDIPI and EMPLID
    // unless it is a safety move
    emplid:
      !isSafetyMove &&
      showEmplid &&
      Yup.string().when('affiliation', {
        is: (affiliationValue) => affiliationValue === departmentIndicators.COAST_GUARD,
        then: () =>
          Yup.string()
            .matches(/^(SM[0-9]{5}|[0-9]{7})$/, 'Enter a 7-digit EMPLID number')
            .required(`EMPLID is required for the Coast Guard`),
        otherwise: Yup.string().notRequired(),
      }),
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

  const sectionStyles = classnames(styles.noTopMargin, formStyles.formSection);

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
            {({ isValid, handleSubmit, setValues, values, handleChange, ...formikProps }) => {
              const handleIsSafetyMove = (e) => {
                const { value } = e.target;
                if (value === 'true') {
                  setIsSafetyMove(true);
                  setShowSafetyMoveHint(true);
                  setValues({
                    ...values,
                    affiliation: '',
                    create_okta_account: '',
                    cac_user: 'true',
                    is_safety_move: 'true',
                  });
                } else if (value === 'false') {
                  setIsSafetyMove(false);
                  setValues({
                    ...values,
                    affiliation: '',
                    edipi: '',
                    emplid: null,
                    is_safety_move: 'false',
                  });
                }
              };
              const handleBranchChange = (e) => {
                setShowSafetyMoveHint(false);
                if (e.target.value === departmentIndicators.COAST_GUARD && isSafetyMove) {
                  setShowEmplid(true);
                  setValues({
                    ...values,
                    affiliation: e.target.value,
                    edipi: uniqueDodid,
                    emplid: uniqueEmplid,
                  });
                } else if (e.target.value === departmentIndicators.COAST_GUARD && !isSafetyMove) {
                  setShowEmplid(true);
                  setValues({
                    ...values,
                    affiliation: e.target.value,
                    edipi: '',
                    emplid: null,
                  });
                } else if (e.target.value !== departmentIndicators.COAST_GUARD && isSafetyMove) {
                  setShowEmplid(false);
                  setValues({
                    ...values,
                    affiliation: e.target.value,
                    edipi: uniqueDodid,
                    emplid: null,
                  });
                } else {
                  setShowEmplid(false);
                  setValues({
                    ...values,
                    affiliation: e.target.value,
                    edipi: '',
                    emplid: null,
                  });
                }
              };
              return (
                <Form className={classnames(formStyles.form, styles.form)}>
                  <h1 className={styles.header}>Create Customer Profile</h1>
                  <SectionWrapper className={sectionStyles}>
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
                            checked={values.is_safety_move === 'true'}
                          />
                          <Field
                            as={Radio}
                            id="isSafetyMoveNo"
                            label="No"
                            name="is_safety_move"
                            value="false"
                            data-testid="is-safety-move-no"
                            onChange={handleIsSafetyMove}
                            checked={values.is_safety_move === 'false'}
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
                      maxLength="10"
                      isDisabled={isSafetyMove}
                      data-testid="edipiInput"
                    />
                    {showEmplid && (
                      <TextField
                        label="EMPLID"
                        name="emplid"
                        id="emplid"
                        maxLength="7"
                        inputMode="numeric"
                        pattern="[0-9]{7}"
                        isDisabled={isSafetyMove}
                        data-testid="emplidInput"
                      />
                    )}
                    {isSafetyMove && showSafetyMoveHint && (
                      <Hint data-testid="safetyMoveHint">
                        Once a branch is selected, this will generate a random safety move identifier
                      </Hint>
                    )}
                  </SectionWrapper>
                  <SectionWrapper className={sectionStyles}>
                    <h3>Customer Name</h3>
                    <TextField label="First name" name="first_name" id="firstName" required />
                    <TextField label="Middle name" name="middle_name" id="middleName" labelHint="Optional" />
                    <TextField label="Last name" name="last_name" id="lastName" required />
                    <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
                  </SectionWrapper>
                  <SectionWrapper className={sectionStyles}>
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
                    <h3>Pickup Address</h3>
                    <AddressFields
                      name={residentialAddressName}
                      labelHint="Required"
                      locationLookup
                      formikProps={formikProps}
                    />
                  </SectionWrapper>
                  <SectionWrapper className={sectionStyles}>
                    <h3>Backup Address</h3>
                    <AddressFields
                      name={backupAddressName}
                      labelHint="Required"
                      locationLookup
                      formikProps={formikProps}
                    />
                  </SectionWrapper>
                  <SectionWrapper className={sectionStyles}>
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
                    <SectionWrapper className={sectionStyles}>
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
                    <SectionWrapper className={sectionStyles}>
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
  setCanAddOrders: setCanAddOrdersAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreateCustomerForm);
