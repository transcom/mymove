import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import { Button, TextInput, Label, FormGroup, Fieldset, ErrorMessage, Grid, Alert } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Checkbox, FormControlLabel } from '@material-ui/core';

import styles from './SubmitMoveForm.module.scss';

import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import CertificationText from 'components/CertificationText/CertificationText';

const SubmitMoveForm = (props) => {
  const { initialValues, onPrint, onSubmit, onBack, certificationText, error, currentUser } = props;
  const [hasReadTheAgreement, setHasReadTheAgreement] = useState(false);
  const [hasAcknowledgedTerms, sethasAcknowledgedTerms] = useState(false);

  const validationSchema = Yup.object().shape({
    signature: Yup.string()
      .required('Required')
      .test('matches-user-name', 'Typed signature must match your exact user name', (signature) => {
        return signature === currentUser;
      }),
    date: Yup.date().required(),
  });

  const hasAgreedToTheTermsEvent = (event) => {
    sethasAcknowledgedTerms(event.target.checked);
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} validateOnBlur onSubmit={onSubmit}>
      {({ isValid, errors, touched, handleSubmit, isSubmitting }) => {
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
                <FormControlLabel
                  control={
                    <Checkbox
                      name="acknowledgementCheckbox"
                      color="primary"
                      disabled={!hasReadTheAgreement}
                      onChange={hasAgreedToTheTermsEvent}
                    />
                  }
                  label="I have read and understand the agreement as shown above"
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
                    <Grid tablet={{ col: 'fill' }}>
                      <FormGroup error={showSignatureError}>
                        <Label htmlFor="signature">SIGNATURE</Label>
                        {showSignatureError && (
                          <ErrorMessage id="signature-error-message">{errors.signature}</ErrorMessage>
                        )}
                        <Field
                          as={TextInput}
                          disabled={!hasAcknowledgedTerms}
                          name="signature"
                          id="signature"
                          required
                          aria-describedby={showSignatureError ? 'signature-error-message' : null}
                        />
                      </FormGroup>
                    </Grid>
                    <Grid tablet={{ col: 'auto' }}>
                      <FormGroup>
                        <Label htmlFor="date">Date</Label>
                        <Field as={TextInput} name="date" id="date" disabled />
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
                disableNext={!isValid || isSubmitting}
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
