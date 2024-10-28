import React from 'react';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import { FormGroup, Label, Radio, Link as USWDSLink } from '@trussworks/react-uswds';

import { DatePickerInput, DropdownInput, DutyLocationInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';
import { dropdownInputOptions } from 'utils/formatters';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';

const AddOrdersForm = ({ onSubmit, ordersTypeOptions, initialValues, onBack, isSafetyMoveSelected }) => {
  const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

  const validationSchema = Yup.object().shape({
    ordersType: Yup.mixed()
      .oneOf(ordersTypeOptions.map((i) => i.key))
      .required('Required'),
    issueDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    reportByDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    hasDependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    originDutyLocation: Yup.object().nullable().required('Required'),
    newDutyLocation: Yup.object().nullable().required('Required'),
    grade: Yup.mixed().oneOf(Object.keys(ORDERS_PAY_GRADE_OPTIONS)).required('Required'),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ values, isValid, isSubmitting, handleSubmit }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.ordersType);
        return (
          <Form className={`${formStyles.form}`}>
            <h1>Tell us about the orders</h1>

            <SectionWrapper className={formStyles.formSection}>
              <DropdownInput
                label="Orders type"
                name="ordersType"
                options={ordersTypeOptions}
                required
                isDisabled={isSafetyMoveSelected}
              />
              <DatePickerInput name="issueDate" label="Orders date" required />
              <DatePickerInput name="reportByDate" label="Report by date" required />
              <FormGroup>
                <Label>Are dependents included in the orders?</Label>
                <div>
                  <Field
                    as={Radio}
                    label="Yes"
                    id="hasDependentsYes"
                    data-testid="hasDependentsYes"
                    name="hasDependents"
                    value="yes"
                    title="Yes, dependents are included in my orders"
                    type="radio"
                  />
                  <Field
                    as={Radio}
                    label="No"
                    id="hasDependentsNo"
                    data-testid="hasDependentsNo"
                    name="hasDependents"
                    value="no"
                    title="No, dependents are not included in my orders"
                    type="radio"
                  />
                </div>
              </FormGroup>

              <DutyLocationInput
                label="Current duty location"
                name="originDutyLocation"
                id="originDutyLocation"
                required
              />

              {isRetirementOrSeparation ? (
                <>
                  <h3>Where are they entitled to move?</h3>
                  <Callout>
                    <span>The government will pay for their move to:</span>
                    <ul>
                      <li>Home of record (HOR)</li>
                      <li>Place entered active duty (PLEAD)</li>
                    </ul>
                    <p>
                      It might pay for a move to their Home of selection (HOS), anywhere in CONUS. Check their orders.
                    </p>
                    <p>
                      Read more about where they are entitled to move when leaving the military on{' '}
                      <USWDSLink
                        target="_blank"
                        rel="noopener noreferrer"
                        href="https://www.militaryonesource.mil/military-life-cycle/separation-transition/military-separation-retirement/deciding-where-to-live-when-you-leave-the-military/"
                      >
                        Military OneSource.
                      </USWDSLink>
                    </p>
                  </Callout>
                  <DutyLocationInput
                    name="newDutyLocation"
                    label="HOR, PLEAD or HOS"
                    displayAddress={false}
                    placeholder="Enter a city or ZIP"
                  />
                </>
              ) : (
                <DutyLocationInput name="newDutyLocation" label="New duty location" required />
              )}
              <DropdownInput label="Pay grade" name="grade" id="grade" required options={payGradeOptions} />
            </SectionWrapper>

            <div className={formStyles.formActions}>
              <WizardNavigation
                disableNext={!isValid || isSubmitting}
                onNextClick={handleSubmit}
                onBackClick={onBack}
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

export default AddOrdersForm;
