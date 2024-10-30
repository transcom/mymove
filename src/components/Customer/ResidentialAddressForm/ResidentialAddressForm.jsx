import { React, useState } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import classnames from 'classnames';
import * as Yup from 'yup';

import styles from './ResidentialAddressForm.module.scss';

import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';
import { ResidentialAddressShape } from 'types/address';

const ResidentialAddressForm = ({ formFieldsName, initialValues, onSubmit, onBack, validators }) => {
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredAddressSchema.required(),
  });
  const [isLookupErrorVisible, setIsLookupErrorVisible] = useState(false);

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validateOnChange={false}
      validateOnMount
      validationSchema={validationSchema}
    >
      {({ isValid, isSubmitting, handleSubmit, values, setValues }) => {
        const handleZipCityChange = (value) => {
          setValues(
            {
              ...values,
              current_residence: {
                ...values.current_residence,
                city: value.city ? value.city : '',
                state: value.state ? value.state : '',
                county: value.county ? value.county : '',
                postalCode: value.postalCode ? value.postalCode : '',
                usprcId: value.usPostRegionCitiesId ? value.usPostRegionCitiesId : '',
              },
            },
            { shouldValidate: true },
          );

          if (!value.city || !value.state || !value.county || !value.postalCode) {
            setIsLookupErrorVisible(true);
          } else {
            setIsLookupErrorVisible(false);
          }
        };

        return (
          <Form className={formStyles.form}>
            <h1>Current address</h1>
            <p className={styles.noBottomMargin}>Must be a physical address.</p>
            <SectionWrapper className={classnames(styles.noTopMargin, formStyles.formSection)}>
              <AddressFields
                labelHint="Required"
                name={formFieldsName}
                validators={validators}
                zipCityError={isLookupErrorVisible}
                handleZipCityChange={handleZipCityChange}
              />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
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

ResidentialAddressForm.propTypes = {
  formFieldsName: PropTypes.string.isRequired,
  initialValues: ResidentialAddressShape.isRequired,
  onBack: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
    county: PropTypes.func,
    usprcId: PropTypes.func,
  }),
};

ResidentialAddressForm.defaultProps = {
  validators: {},
};

export default ResidentialAddressForm;
