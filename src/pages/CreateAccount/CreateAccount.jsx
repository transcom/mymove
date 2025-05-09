import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { useNavigate, Link } from 'react-router-dom';
import { Grid, GridContainer, Alert, Button, Label } from '@trussworks/react-uswds';
import { Form, Formik } from 'formik';
import classNames from 'classnames';
import * as Yup from 'yup';

import ValidCACModal from '../../components/ValidCACModal/ValidCACModal';

import styles from './CreateAccount.module.scss';

import formStyles from 'styles/form.module.scss';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes } from 'constants/routes';
import TextField from 'components/form/fields/TextField/TextField';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import departmentIndicators from 'constants/departmentIndicators';
import StyledLine from 'components/StyledLine/StyledLine';
import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';
import RegistrationConfirmationModal from 'components/RegistrationConfirmationModal/RegistrationConfirmationModal';
import { registerUser } from 'services/internalApi';
import Hint from 'components/Hint';
import { technicalHelpDeskURL } from 'shared/constants';
import ValidationCode from 'pages/MyMove/Profile/ValidationCode';
import { isBooleanFlagEnabledUnauthenticated } from 'utils/featureFlags';

export const CreateAccount = ({ setShowLoadingSpinner }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const [showEmplid, setShowEmplid] = useState(false);
  const [isCACModalVisible, setIsCACModalVisible] = useState(false);
  const [isConfirmationModalVisible, setIsConfirmationModalVisible] = useState(false);
  const [showValidationCode, setShowValidationCode] = useState(false);
  const [validationCodeFF, setValidationCodeFF] = useState(null);
  const [showHint, setShowHint] = useState(false);
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);

  const hostname = window && window.location && window.location.hostname;
  const oktaURL =
    hostname === 'my.move.mil'
      ? 'https://milmove.okta.mil/enduser/settings'
      : 'https://test-milmove.okta.mil/enduser/settings';

  useEffect(() => {
    let timer;

    const fetchFeatureFlag = async () => {
      const flagEnabled = await isBooleanFlagEnabledUnauthenticated('validation_code_required');
      setValidationCodeFF(flagEnabled);
      if (flagEnabled) {
        setShowValidationCode(true);
      } else {
        timer = setTimeout(() => {
          setIsCACModalVisible(true);
        }, 200);
      }
    };

    fetchFeatureFlag();

    return () => {
      if (timer) {
        clearTimeout(timer);
      }
    };
  }, []);

  const handleSuccessfulValidation = () => {
    setShowValidationCode(false);
    setIsCACModalVisible(true);
  };

  const handleCACModalYes = () => {
    setIsCACModalVisible(false);
  };

  const handleCACModalNo = () => {
    navigate('/sign-in', {
      state: { noValidCAC: true },
    });
  };

  const handleConfirmationModalYes = () => {
    window.location.href = '/auth/okta';
  };

  const initialValues = {
    affiliation: '',
    edipi: '',
    edipiConfirmation: '',
    emplid: '',
    emplidConfirmation: '',
    firstName: '',
    middleInitial: '',
    lastName: '',
    email: '',
    emailConfirmation: '',
    telephone: '',
    secondaryTelephone: '',
    phoneIsPreferred: false,
    emailIsPreferred: false,
  };

  const handleCancel = () => {
    navigate(generalRoutes.SIGN_IN_PATH);
  };

  const handleSubmit = async (values) => {
    if (values.firstName && values.lastName) {
      setShowLoadingSpinner(true, `Creating MilMove Account for ${values.firstName} ${values.lastName}`);
    } else {
      setShowLoadingSpinner(true, `Creating MilMove Account`);
    }
    const payload = {
      affiliation: values.affiliation,
      edipi: values.edipi,
      emplid: values.emplid.trim() === '' ? null : values.emplid,
      firstName: values.firstName,
      middleInitial: values.middleInitial.trim() === '' ? null : values.middleInitial,
      lastName: values.lastName,
      email: values.email,
      telephone: values.telephone,
      secondaryTelephone: values.secondaryTelephone.trim() === '' ? null : values.secondaryTelephone,
      phoneIsPreferred: values.phoneIsPreferred,
      emailIsPreferred: values.emailIsPreferred,
    };
    await registerUser(payload)
      .then(() => {
        setShowLoadingSpinner(false, null);
        setIsConfirmationModalVisible(true);
      })
      .catch((e) => {
        const { response } = e;
        let errorMessage = `There was an error creating your account`;
        setShowLoadingSpinner(false, null);
        if (response.body) {
          const responseBody = response.body;
          let responseMsg = '';

          if (responseBody.detail) {
            responseMsg += `${responseBody.detail}`;
          }

          errorMessage += `\n${responseMsg}`;
        }
        setServerError(errorMessage);
        setShowHint(true);
      });
  };

  const validationSchema = Yup.object().shape({
    affiliation: Yup.mixed().oneOf(Object.keys(SERVICE_MEMBER_AGENCY_LABELS)).required('Required'),
    edipi: Yup.string()
      .matches(/^(SM[0-9]{8}|[0-9]{10})$/, 'Enter a 10-digit DoD ID number')
      .required('Required'),
    edipiConfirmation: Yup.string()
      .oneOf([Yup.ref('edipi'), null], 'DoD ID numbers must match')
      .required('Required'),
    emplid:
      showEmplid &&
      Yup.string().when('affiliation', {
        is: (affiliationValue) => affiliationValue === departmentIndicators.COAST_GUARD,
        then: () =>
          Yup.string()
            .matches(/^(SM[0-9]{5}|[0-9]{7})$/, 'Enter a 7-digit EMPLID number')
            .required(`EMPLID is required for the Coast Guard`),
        otherwise: Yup.string().notRequired(),
      }),
    emplidConfirmation:
      showEmplid &&
      Yup.string()
        .oneOf([Yup.ref('emplid'), null], 'EMPLID numbers must match')
        .required('Required'),
    firstName: Yup.string().required('Required'),
    middleName: Yup.string(),
    lastName: Yup.string().required('Required'),
    suffix: Yup.string(),
    telephone: Yup.string()
      .min(12, 'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.')
      .required('Required'),
    secondaryTelephone: Yup.string()
      .min(12, 'Please enter a valid phone number. Phone numbers must be entered as ###-###-####.')
      .nullable(),
    email: Yup.string()
      .matches(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$/, 'Must be a valid email address')
      .required('Required'),
    emailConfirmation: Yup.string()
      .oneOf([Yup.ref('email'), null], 'Emails must match')
      .required('Required'),
    phoneIsPreferred: Yup.boolean(),
    emailIsPreferred: Yup.boolean(),
  });

  return (
    <div className={classNames('usa-prose grid-container')}>
      <ValidCACModal isOpen={isCACModalVisible} onClose={handleCACModalNo} onSubmit={handleCACModalYes} />
      <RegistrationConfirmationModal isOpen={isConfirmationModalVisible} onSubmit={handleConfirmationModalYes} />
      <NotificationScrollToTop dependency={serverError} />
      <GridContainer>
        {showValidationCode && validationCodeFF ? (
          <ValidationCode onSuccess={handleSuccessfulValidation} />
        ) : (
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }} className={styles.formContainer}>
              {serverError && (
                <Grid row>
                  <Alert
                    data-testid="alert2"
                    type="error"
                    headingLevel="h4"
                    heading="An error occurred"
                    className={styles.error}
                  >
                    {serverError}
                  </Alert>
                </Grid>
              )}
              <Formik
                initialValues={initialValues}
                onSubmit={handleSubmit}
                validateOnMount
                validateOnChange
                validateOnBlur
                validationSchema={validationSchema}
              >
                {({ isSubmitting, isValid, values, setValues, handleChange }) => {
                  const handleBranchChange = (e) => {
                    if (e.target.value === departmentIndicators.COAST_GUARD) {
                      setShowEmplid(true);
                      setValues({
                        ...values,
                        affiliation: e.target.value,
                      });
                    } else if (e.target.value !== departmentIndicators.COAST_GUARD) {
                      setShowEmplid(false);
                      setValues({
                        ...values,
                        affiliation: e.target.value,
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
                    <Form className={formStyles.formSection}>
                      <SectionWrapper>
                        <div className={styles.centerColumn}>
                          <h2>MilMove Registration</h2>
                          {showHint && (
                            <Hint className={styles.hint}>
                              MilMove uses Okta for authentication. If you need to access an exsiting Okta account, you
                              can access the Okta dashboard by{' '}
                              <a className={styles.link} href={oktaURL} target="_blank" rel="noreferrer">
                                <strong> clicking this link</strong>.
                              </a>
                              <br />
                              <br />
                              If you continue to have issues with registration <br />
                              please contact the&nbsp;
                              <Link to={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                                Technical Help Desk
                              </Link>
                            </Hint>
                          )}
                        </div>
                        <div className={styles.formSection}>
                          <DropdownInput
                            label="Branch of service"
                            name="affiliation"
                            id="affiliation"
                            data-testid="affiliationInput"
                            required
                            showRequiredAsterisk
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
                            data-testid="edipiInput"
                            required
                            showRequiredAsterisk
                          />
                          <TextField
                            label="Confirm DoD ID number"
                            name="edipiConfirmation"
                            id="edipiConfirmation"
                            maxLength="10"
                            data-testid="edipiConfirmationInput"
                            disablePaste
                            required
                            showRequiredAsterisk
                          />
                          {showEmplid && (
                            <>
                              <TextField
                                label="EMPLID"
                                name="emplid"
                                id="emplid"
                                maxLength="7"
                                inputMode="numeric"
                                pattern="[0-9]{7}"
                                data-testid="emplidInput"
                                required
                                showRequiredAsterisk
                              />
                              <TextField
                                label="Confirm EMPLID"
                                name="emplidConfirmation"
                                id="emplidConfirmation"
                                maxLength="7"
                                inputMode="numeric"
                                pattern="[0-9]{7}"
                                data-testid="emplidConfirmationInput"
                                disablePaste
                                required
                                showRequiredAsterisk
                              />
                            </>
                          )}
                          <StyledLine />
                          <TextField
                            label="First Name"
                            name="firstName"
                            id="firstName"
                            data-testid="firstName"
                            required
                            showRequiredAsterisk
                          />
                          <TextField
                            label="Middle Initial"
                            name="middleInitial"
                            id="middleInitial"
                            data-testid="middleInitial"
                          />
                          <TextField
                            label="Last Name"
                            name="lastName"
                            id="lastName"
                            data-testid="lastName"
                            required
                            showRequiredAsterisk
                          />
                          <StyledLine />
                          <TextField
                            label="Email"
                            name="email"
                            id="email"
                            data-testid="email"
                            required
                            showRequiredAsterisk
                          />
                          <TextField
                            label="Confirm Email"
                            name="emailConfirmation"
                            id="emailConfirmation"
                            disablePaste
                            data-testid="emailConfirmation"
                            required
                            showRequiredAsterisk
                          />
                          <StyledLine />
                          <MaskedTextField
                            label="Telephone"
                            id="telephone"
                            name="telephone"
                            type="tel"
                            minimum="12"
                            mask="000{-}000{-}0000"
                            data-testid="telephone"
                            required
                            showRequiredAsterisk
                          />
                          <MaskedTextField
                            label="Secondary Telephone"
                            id="secondaryTelephone"
                            name="secondaryTelephone"
                            type="tel"
                            minimum="12"
                            mask="000{-}000{-}0000"
                            data-testid="secondaryTelephone"
                          />
                          <Label className={styles.checkboxLabel}>Preferred contact method</Label>
                          <div className={classNames(formStyles.radioGroup, formStyles.customerPreferredContact)}>
                            <CheckboxField
                              id="phoneIsPreferred"
                              label="Phone"
                              name="phoneIsPreferred"
                              data-testid="phoneIsPreferred"
                            />
                            <CheckboxField
                              id="emailIsPreferred"
                              label="Email"
                              name="emailIsPreferred"
                              data-testid="emailIsPreferred"
                            />
                          </div>
                        </div>
                      </SectionWrapper>

                      <div className={styles.buttonRow}>
                        <Button type="submit" disabled={!isValid || isSubmitting} data-testid="submitBtn">
                          Submit
                        </Button>
                        <Button type="button" onClick={handleCancel} secondary data-testid="cancelBtn">
                          Cancel
                        </Button>
                      </div>
                    </Form>
                  );
                }}
              </Formik>
            </Grid>
          </Grid>
        )}
      </GridContainer>
    </div>
  );
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(() => ({}), mapDispatchToProps)(CreateAccount);
