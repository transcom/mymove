import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Grid, GridContainer, Alert, Button } from '@trussworks/react-uswds';
import { Form, Formik } from 'formik';
import classNames from 'classnames';
import * as Yup from 'yup';

import styles from './SignUp.module.scss';
import ValidCACModal from './ValidCACModal';

import formStyles from 'styles/form.module.scss';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes } from 'constants/routes';
import TextField from 'components/form/fields/TextField/TextField';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import departmentIndicators from 'constants/departmentIndicators';
import StyledLine from 'components/StyledLine/StyledLine';
import LoadingSpinnerModal from 'components/LoadingSpinnerModal/LoadingSpinnerModal';

export const SignUp = () => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const [showEmplid, setShowEmplid] = useState(false);
  const [isDisabled, setIsDisabled] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [loadingMessage, setLoadingMessage] = useState('Loading');
  const branchOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);

  // Timer to show the modal after 5 seconds
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsModalVisible(true);
    }, 500);

    return () => clearTimeout(timer);
  }, []);

  const handleModalYes = () => {
    setIsModalVisible(false);
  };

  const handleModalNo = () => {
    navigate('/sign-in', {
      state: { noValidCAC: true },
    });
  };

  const initialValues = {
    affiliation: '',
    dodid: '',
    dodidConfirmation: '',
    firstName: '',
    middleInitial: '',
    lastName: '',
    email: '',
    telephone: '',
    secondaryTelephone: '',
    phoneIsPreferred: false,
    emailIsPreferred: false,
  };

  const handleCancel = () => {
    navigate(generalRoutes.SIGN_IN_PATH);
  };

  const delay = (ms) => {
    return new Promise((resolve) => {
      setTimeout(() => resolve(), ms);
    });
  };

  const handleStartLoading = async () => {
    setIsDisabled(true);
    setIsLoading(true);
    setLoadingMessage('Creating MilMove Profile');

    await delay(4000); // Wait 3 seconds
    setLoadingMessage('Creating Okta Profile');

    await delay(4000); // Wait 3 more seconds

    // Navigate to the desired URL
    window.location.href = '/auth/okta';
    setIsLoading(false);
  };

  const handleSubmit = async (values) => {
    const body = {
      affiliation: values.affiliation,
      edipi: values.edipi,
      firstName: values.firstName,
      middleInitial: values.middleInitial,
      lastName: values.lastName,
      email: values.email,
      telephone: values.telephone,
      secondaryTelephone: values.secondaryTelephone,
      phoneIsPreferred: values.phoneIsPreferred,
      emailIsPreferred: values.emailIsPreferred,
    };
    await handleStartLoading();
    console.log('submitted!', body);
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
    phoneIsPreferred: Yup.boolean(),
    emailIsPreferred: Yup.boolean(),
  });

  return (
    <div className={classNames('usa-prose grid-container padding-top-3')}>
      <LoadingSpinnerModal isOpen={isLoading} message={loadingMessage} />
      <ValidCACModal isOpen={isModalVisible} onClose={handleModalNo} onSubmit={handleModalYes} />
      <GridContainer>
        <NotificationScrollToTop dependency={serverError} />

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

        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }} className={styles.formContainer}>
            <Formik
              initialValues={initialValues}
              onSubmit={handleSubmit}
              validateOnMount
              validateOnChange
              validateOnBlur
              validationSchema={validationSchema}
            >
              {({ isValid, values, setValues, handleChange }) => {
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
                  <Form>
                    <SectionWrapper className={formStyles.formSection}>
                      <h2>MilMove Registration</h2>
                      <DropdownInput
                        label="Branch of service"
                        name="affiliation"
                        id="affiliation"
                        data-testid="affiliationInput"
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
                        data-testid="edipiInput"
                      />
                      <TextField
                        label="Confirm DoD ID number"
                        name="edipiConfirmation"
                        id="edipiConfirmation"
                        maxLength="10"
                        data-testid="edipiConfirmationInput"
                      />
                      {showEmplid && (
                        <TextField
                          label="EMPLID"
                          name="emplid"
                          id="emplid"
                          maxLength="7"
                          inputMode="numeric"
                          pattern="[0-9]{7}"
                          data-testid="emplidInput"
                        />
                      )}
                      <StyledLine />
                      <TextField label="First Name" name="firstName" id="firstName" />
                      <TextField label="Middle Initial" name="middleInitial" id="middleInitial" />
                      <TextField label="Last Name" name="lastName" id="lastName" />
                      <StyledLine />
                      <TextField label="Email" name="email" id="email" />
                      <MaskedTextField
                        label="Telephone"
                        id="telephone"
                        name="telephone"
                        type="tel"
                        minimum="12"
                        mask="000{-}000{-}0000"
                      />
                      <MaskedTextField
                        label="Secondary Telephone"
                        id="secondaryTelephone"
                        name="secondaryTelephone"
                        type="tel"
                        minimum="12"
                        mask="000{-}000{-}0000"
                      />
                      <div className={styles.radioGroup}>
                        <CheckboxField id="phoneIsPreferred" label="Phone" name="phone_is_preferred" />
                        <CheckboxField id="emailIsPreferred" label="Email" name="email_is_preferred" />
                      </div>
                    </SectionWrapper>

                    <div className={styles.buttonRow}>
                      <Button type="submit" disabled={!isValid || isDisabled} onClick={handleSubmit}>
                        Submit
                      </Button>
                      <Button type="button" disabled={isDisabled} onClick={handleCancel} secondary>
                        Cancel
                      </Button>
                    </div>
                  </Form>
                );
              }}
            </Formik>
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(SignUp);
