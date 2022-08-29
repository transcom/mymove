import React from 'react';
import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import { func, number } from 'prop-types';
import { Button, Form, Radio } from '@trussworks/react-uswds';
import { FormGroup } from '@material-ui/core';

import formStyles from 'styles/form.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';
import { ProGearTicketShape } from 'types/shipment';

const validationSchema = Yup.object().shape({
  selfProGear: Yup.string().required('Required'),
});

const ProGearForm = ({ proGear, setNumber, onSubmit, onBack }) => {
  const { selfProGear } = proGear || {};
  const initialValues = {
    selfProGear: selfProGear ? 'true' : 'false',
  };
  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, values }) => {
        return (
          <div className={ppmStyles.formContainer}>
            <Form className={ppmStyles.form}>
              <SectionWrapper className={formStyles.formSection}>
                <h2>Set {setNumber}</h2>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label margin-bottom-0">Pro-gear belongs to</legend>
                    <Field
                      as={Radio}
                      id="ownerOfProGearSelf"
                      label="Me"
                      name="selfProGear"
                      value="true"
                      checked={values.selfProGear === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="ownerOfProGearSpouse"
                      label="My spouse"
                      name="selfProGear"
                      value="false"
                      checked={values.selfProGear === 'false'}
                    />
                  </Fieldset>
                </FormGroup>
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Finish Later
                </Button>
                <Button
                  className={ppmStyles.saveButton}
                  type="button"
                  onClick={onSubmit}
                  disabled={!isValid || isSubmitting}
                >
                  Save &amp; Continue
                </Button>
              </div>
            </Form>
          </div>
        );
      }}
    </Formik>
  );
};

ProGearForm.propTypes = {
  setNumber: number,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  proGear: ProGearTicketShape,
};

ProGearForm.defaultProps = {
  setNumber: 1,
  proGear: {
    selfProGear: true,
  },
};

export default ProGearForm;
