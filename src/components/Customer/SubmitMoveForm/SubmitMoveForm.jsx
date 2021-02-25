import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import { Button, TextInput, Label, FormGroup, Fieldset } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import { Form } from 'components/form/Form';
import SectionWrapper from 'components/Customer/SectionWrapper';
import styles from 'styles/form.module.scss';

const SubmitMoveForm = (props) => {
  const { handlePrint } = props;

  const validationSchema = Yup.object().shape({
    signature: Yup.string().required(),
    date: Yup.date().required(),
  });

  const initialValues = {
    signature: '',
    date: '',
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema}>
      {() => (
        <Form className={styles.form}>
          <h1>Now for the official part&hellip;</h1>
          <div className="instructions">
            <p>
              Please read this agreement, type your name in the <strong>Signature</strong> field to sign it, then tap
              the <strong>Complete</strong> button.
            </p>
            <p>This agreement covers the shipment of your personal property.</p>
          </div>

          <SectionWrapper>
            <Button type="button" unstyled onClick={handlePrint}>
              Print
            </Button>

            <div>
              <h3>SIGNATURE</h3>
              <p>
                In consideration of said household goods or mobile homes being shipped at Government expense, I hereby
                agree to the certifications stated above.
              </p>
              <Fieldset>
                <FormGroup>
                  <Label>Signature</Label>
                  <Field as={TextInput} name="signature" />
                </FormGroup>

                <FormGroup>
                  <Label>Date</Label>
                  <Field as={TextInput} name="date" disabled />
                </FormGroup>
              </Fieldset>
            </div>
          </SectionWrapper>
        </Form>
      )}
    </Formik>
  );
};

SubmitMoveForm.propTypes = {
  handlePrint: PropTypes.func,
};

SubmitMoveForm.defaultProps = {
  handlePrint: () => window.print(),
};

export default SubmitMoveForm;
