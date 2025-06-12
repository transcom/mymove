import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import { Button, TextInput, Label, FormGroup, Fieldset, ErrorMessage, Grid, Alert } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Checkbox, FormControlLabel } from '@material-ui/core';

import styles from './SubmitMoveForm.module.scss';

import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import CertificationText from 'components/CertificationText/CertificationText';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const SubmitMoveForm = (props) => {
  const { initialValues, onPrint, onSubmit, onBack, certificationText, error, currentUser } = props;
  const [hasReadTheAgreement, setHasReadTheAgreement] = useState(false);
  const [hasAcknowledgedTerms, sethasAcknowledgedTerms] = useState(false);

  const normalizeString = (str) => {
    return str
      .toLowerCase()
      .trim()
      .replace(/\s+/g, ' ')
      .replace(/[^\w\s]/gi, '');
  };

  const compareSignature = (signature, fullName) => {
    const normalizedSignature = normalizeString(signature);
    const normalizedFullName = normalizeString(fullName);

    return normalizedSignature === normalizedFullName;
  };

  const validationSchema = Yup.object().shape({
    signature: Yup.string()
      .required('Required')
      .test('matches-user-name', 'Typed signature must match your exact user name', (signature) => {
        return compareSignature(signature, currentUser);
      }),
    date: Yup.date().required(),
  });

  const hasAgreedToTheTermsEvent = (event) => {
    sethasAcknowledgedTerms(event.target.checked);
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} validateOnBlur onSubmit={onSubmit}>
      {({ isValid, errors, touched, handleSubmit, isSubmitting, dirty }) => {
        const showSignatureError = !!(errors.signature && touched.signature);

        return (
          <Form className={`${formStyles.form} ${styles.SubmitMoveForm}`}>
            <h1>Now for the official part&hellip;</h1>
            <p>
              Please read this agreement, type your name in the <strong>SIGNATURE</strong> field to sign it, then click
              the <strong>Complete</strong> button.
            </p>
            <p>This agreement covers the shipment of your personal property.</p>

            <SectionWrapper>
              <Button type="button" unstyled onClick={onPrint} className={styles.hideForPrint}>
                Print
              </Button>

              <CertificationText certificationText={certificationText} onScrollToBottom={setHasReadTheAgreement} />

              <FormGroup>
                {requiredAsteriskMessage}
                <FormControlLabel
                  className={!hasReadTheAgreement ? styles.disabledCheckbox : ''}
                  control={
                    <Checkbox
                      data-testid="acknowledgementCheckbox"
                      name="acknowledgementCheckbox"
                      color={!hasReadTheAgreement ? '#6c757d' : 'primary'}
                      checked={hasAcknowledgedTerms}
                      readOnly={!hasReadTheAgreement}
                      onChange={hasAgreedToTheTermsEvent}
                      inputProps={{
                        'aria-disabled': !hasReadTheAgreement,
                      }}
                    />
                  }
                  label={
                    <>
                      <RequiredAsterisk /> I have read and understand the agreement as shown above
                    </>
                  }
                />
              </FormGroup>

              <div className={styles.signatureBox}>
                <h3>SIGNATURE</h3>
                <p>
                  In consideration of said household goods or mobile homes being shipped at Government expense, I hereby
                  agree to the certifications stated above.
                </p>

                <Fieldset>
                  <Grid row gap>
                    <Grid tablet={{ col: 'fill' }} className={styles.dateGrid}>
                      <FormGroup error={showSignatureError}>
                        <Label htmlFor="signature">
                          <span>
                            SIGNATURE <RequiredAsterisk />
                          </span>
                        </Label>
                        {showSignatureError && (
                          <ErrorMessage id="signature-error-message">{errors.signature}</ErrorMessage>
                        )}
                        <Field
                          as={TextInput}
                          id="signature"
                          name="signature"
                          readOnly={!hasAcknowledgedTerms}
                          aria-readonly={!hasAcknowledgedTerms}
                          aria-required="true"
                          aria-describedby={errors.signature && touched.signature ? 'signature-error' : undefined}
                          onFocus={(e) => !hasAcknowledgedTerms && e.target.blur()}
                          className={!hasAcknowledgedTerms ? styles.readOnlyInput : ''}
                        />
                      </FormGroup>
                    </Grid>
                    <Grid tablet={{ col: 'auto' }} className={styles.dateGrid}>
                      <FormGroup gap>
                        <Label htmlFor="date" className="dateGrid">
                          Date
                        </Label>
                        <Field
                          as={TextInput}
                          id="date"
                          name="date"
                          readOnly
                          aria-readonly="true"
                          className={styles.readOnlyInput}
                        />
                      </FormGroup>
                    </Grid>
                  </Grid>
                  <Grid row>
                    <Grid tablet={{ col: 'fill' }}>
                      <p>{currentUser}</p>
                    </Grid>
                  </Grid>
                  <Grid row gap>
                    <Grid tablet={{ col: 'fill' }}>
                      <p>Typed signature must match displayed name</p>
                    </Grid>
                  </Grid>
                </Fieldset>
              </div>

              {error && (
                <Alert type="error" headingLevel="h4" heading="Server Error">
                  There was a problem saving your signature.
                </Alert>
              )}
            </SectionWrapper>
            <div className={formStyles.formActions}>
              <WizardNavigation
                isLastPage
                onBackClick={onBack}
                disableNext={!isValid || isSubmitting || !dirty || !hasAcknowledgedTerms}
                onNextClick={handleSubmit}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

SubmitMoveForm.propTypes = {
  certificationText: PropTypes.string,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
  onPrint: PropTypes.func,
  error: PropTypes.oneOfType([PropTypes.bool, PropTypes.object]),
  initialValues: PropTypes.shape({
    signature: PropTypes.string.isRequired,
    date: PropTypes.string.isRequired,
  }).isRequired,
};

SubmitMoveForm.defaultProps = {
  certificationText: null,
  onPrint: () => window.print(),
  error: false,
};

export default SubmitMoveForm;
