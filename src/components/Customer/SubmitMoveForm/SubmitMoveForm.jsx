import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import { Button, TextInput, Label, FormGroup, Fieldset, ErrorMessage } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Customer/SectionWrapper';
import styles from 'styles/form.module.scss';
import { formatSwaggerDate } from 'shared/formatters';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import CertificationText from 'scenes/Legalese/CertificationText';

const SubmitMoveForm = (props) => {
  const { onPrint, onSubmit } = props;

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
      {({ isValid, errors, touched, handleSubmit }) => {
        const showSignatureError = !!(errors.signature && touched.signature);

        return (
          <Form className={styles.form}>
            <h1>Now for the official part&hellip;</h1>
            <p>
              Please read this agreement, type your name in the <strong>Signature</strong> field to sign it, then tap
              the <strong>Complete</strong> button.
            </p>
            <p>This agreement covers the shipment of your personal property.</p>

            <SectionWrapper>
              <Button type="button" unstyled onClick={onPrint}>
                Print
              </Button>

              <CertificationText certificationText={completeCertificationText} />

              <div>
                <h3>SIGNATURE</h3>
                <p>
                  In consideration of said household goods or mobile homes being shipped at Government expense, I hereby
                  agree to the certifications stated above.
                </p>
                <Fieldset>
                  <FormGroup error={showSignatureError}>
                    <Label>Signature</Label>
                    {showSignatureError && <ErrorMessage id="signature-error-message">{errors.signature}</ErrorMessage>}
                    <Field
                      as={TextInput}
                      name="signature"
                      aria-describedby={showSignatureError && 'signature-error-message'}
                    />
                  </FormGroup>

                  <FormGroup>
                    <Label>Date</Label>
                    <Field as={TextInput} name="date" disabled />
                  </FormGroup>
                </Fieldset>
              </div>
            </SectionWrapper>
            <WizardNavigation isLastPage disableNext={!isValid} onNextClick={handleSubmit} />
          </Form>
        );
      }}
    </Formik>
  );
};

SubmitMoveForm.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  onPrint: PropTypes.func,
};

SubmitMoveForm.defaultProps = {
  onPrint: () => window.print(),
};

export default SubmitMoveForm;
