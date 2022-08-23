import React from 'react';
import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { func, number, shape, string } from 'prop-types';
import { Button, Form, Radio } from '@trussworks/react-uswds';
import { FormGroup } from '@material-ui/core';

import styles from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm.module.scss';
import formStyles from 'styles/form.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';

const validationSchema = Yup.object().shape({
  selfProGear: Yup.string().required('Required'),
});

const ProGearForm = ({ initialValues, setNumber, onSubmit, onBack }) => {
  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, values }) => {
        return (
          <div className={classnames(ppmStyles.formContainer, styles.ProGearForm)}>
            <Form className={ppmStyles.form}>
              <SectionWrapper className={classnames(formStyles.formSection, styles.weightTicketSectionWrapper)}>
                <h2>Set {setNumber}</h2>
                <FormGroup>
                  <Fieldset className={styles.ownershipFieldset}>
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
  initialValues: shape({
    selfProGear: string,
  }),
};

ProGearForm.defaultProps = {
  setNumber: 1,
  initialValues: {
    selfProGear: 'true',
  },
};

export default ProGearForm;
