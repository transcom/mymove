import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import { Button, TextInput, Label, FormGroup, Fieldset, ErrorMessage, Grid, Alert } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import styles from './SubmitMoveForm.module.scss';

import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { formatSwaggerDate } from 'shared/formatters';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import CertificationText from 'scenes/Legalese/CertificationText';

const SubmitMoveForm = (props) => {
  const { onPrint, onSubmit, certificationText, error } = props;

  const validationSchema = Yup.object().shape({
    signature: Yup.string().required('Required'),
    date: Yup.date().required(),
  });

  const initialValues = {
    signature: '',
    date: formatSwaggerDate(new Date()),
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} validateOnBlur onSubmit={onSubmit}>
      {({ isValid, errors, touched, handleSubmit, isSubmitting }) => {
        const showSignatureError = !!(errors.signature && touched.signature);

        return (
          <Form className={`${formStyles.form} ${styles.SubmitMoveForm}`}>
            <h1>Now for the official part&hellip;</h1>
            <p>
              Please read this agreement, type your name in the <strong>Signature</strong> field to sign it, then tap
              the <strong>Complete</strong> button.
            </p>
            <p>This agreement covers the shipment of your personal property.</p>

            <SectionWrapper>
              <Button type="button" unstyled onClick={onPrint} className={styles.hideForPrint}>
                Print
              </Button>

              <CertificationText certificationText={certificationText} />

              <div className={styles.signatureBox}>
                <h3>Signature</h3>
                <p>
                  In consideration of said household goods or mobile homes being shipped at Government expense, I hereby
                  agree to the certifications stated above.
                </p>
                <Fieldset>
                  <Grid row gap>
                    <Grid tablet={{ col: 'fill' }}>
                      <FormGroup error={showSignatureError}>
                        <Label>Signature</Label>
                        {showSignatureError && (
                          <ErrorMessage id="signature-error-message">{errors.signature}</ErrorMessage>
                        )}
                        <Field
                          as={TextInput}
                          name="signature"
                          aria-describedby={showSignatureError && 'signature-error-message'}
                        />
                      </FormGroup>
                    </Grid>
                    <Grid tablet={{ col: 'auto' }}>
                      <FormGroup>
                        <Label>Date</Label>
                        <Field as={TextInput} name="date" disabled />
                      </FormGroup>
                    </Grid>
                  </Grid>
                </Fieldset>
              </div>

              {error && (
                <Alert type="error" heading="Server Error">
                  There was a problem saving your signature.
                </Alert>
              )}
            </SectionWrapper>
            <div className={formStyles.formActions}>
              <WizardNavigation isLastPage disableNext={!isValid || isSubmitting} onNextClick={handleSubmit} />
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
  onPrint: PropTypes.func,
  error: PropTypes.oneOf([PropTypes.bool, PropTypes.object]),
};

SubmitMoveForm.defaultProps = {
  certificationText: null,
  onPrint: () => window.print(),
  error: false,
};

export default SubmitMoveForm;
