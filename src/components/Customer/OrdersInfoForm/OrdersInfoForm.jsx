import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink, ErrorMessage } from '@trussworks/react-uswds';

import styles from './OrdersInfoForm.module.scss';

import { ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';
import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { Form } from 'components/form/Form';
import { DropdownArrayOf } from 'types';
import { DutyLocationShape } from 'types/dutyLocation';
import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate, dropdownInputOptions } from 'utils/formatters';

let originMeta;
let newDutyMeta = '';
const OrdersInfoForm = ({ ordersTypeOptions, initialValues, onSubmit, onBack }) => {
  const payGradeOptions = dropdownInputOptions(ORDERS_PAY_GRADE_OPTIONS);

  const validationSchema = Yup.object().shape({
    orders_type: Yup.mixed()
      .oneOf(ordersTypeOptions.map((i) => i.key))
      .required('Required'),
    issue_date: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    report_by_date: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    has_dependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    new_duty_location: Yup.object().nullable().required('Required'),
    grade: Yup.mixed().oneOf(Object.keys(ORDERS_PAY_GRADE_OPTIONS)).required('Required'),
    origin_duty_location: Yup.object().nullable().required('Required'),
  });

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validationSchema={validationSchema}
      onSubmit={onSubmit}
      initialTouched={{ orders_type: true, issue_date: true, report_by_date: true, has_dependents: true, grade: true }}
    >
      {({ isValid, isSubmitting, handleSubmit, values, errors }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.orders_type);

        if (!values.origin_duty_location) originMeta = 'Required';
        else originMeta = null;

        if (!values.new_duty_location) newDutyMeta = 'Required';
        else newDutyMeta = null;

        return (
          <Form className={`${formStyles.form} ${styles.OrdersInfoForm}`}>
            <h1>Tell us about your move orders</h1>

            <SectionWrapper className={formStyles.formSection}>
              <DropdownInput label="Orders type" name="orders_type" options={ordersTypeOptions} required />
              <DatePickerInput
                name="issue_date"
                label="Orders date"
                required
                renderInput={(input) => (
                  <>
                    {input}
                    <Hint>
                      <p>Date your orders were issued.</p>
                    </Hint>
                  </>
                )}
              />
              <DatePickerInput name="report_by_date" label={formatLabelReportByDate(values.orders_type)} required />
              <FormGroup>
                <Label>Are dependents included in your orders?</Label>
                {errors.has_dependents ? <ErrorMessage>{errors.has_dependents}</ErrorMessage> : null}
                <div>
                  <Field
                    as={Radio}
                    label="Yes"
                    id="hasDependentsYes"
                    name="has_dependents"
                    value="yes"
                    title="Yes, dependents are included in my orders"
                    type="radio"
                  />
                  <Field
                    as={Radio}
                    label="No"
                    id="hasDependentsNo"
                    name="has_dependents"
                    value="no"
                    title="No, dependents are not included in my orders"
                    type="radio"
                  />
                </div>
              </FormGroup>

              <DutyLocationInput
                label="Current duty location"
                name="origin_duty_location"
                id="origin_duty_location"
                required
                metaOverride={originMeta}
              />

              {isRetirementOrSeparation ? (
                <>
                  <h3 className={styles.calloutLabel}>Where are you entitled to move?</h3>
                  <Callout>
                    <span>The government will pay for your move to:</span>
                    <ul>
                      <li>Home of record (HOR)</li>
                      <li>Place entered active duty (PLEAD)</li>
                    </ul>
                    <p>
                      It might pay for a move to your Home of selection (HOS), anywhere in CONUS. Check your orders.
                    </p>
                    <p>
                      Read more about where you are entitled to move when leaving the military on{' '}
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
                    name="new_duty_location"
                    label="HOR, PLEAD or HOS"
                    displayAddress={false}
                    hint="Enter the option closest to your destination. Your move counselor will identify if there might be a cost to you."
                    metaOverride={newDutyMeta}
                    placeholder="Enter a city or ZIP"
                  />
                </>
              ) : (
                <DutyLocationInput
                  name="new_duty_location"
                  label="New duty location"
                  displayAddress={false}
                  metaOverride={newDutyMeta}
                />
              )}

              <DropdownInput label="Pay grade" name="grade" id="grade" required options={payGradeOptions} />
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

OrdersInfoForm.propTypes = {
  ordersTypeOptions: DropdownArrayOf.isRequired,
  initialValues: PropTypes.shape({
    orders_type: PropTypes.string,
    issue_date: PropTypes.string,
    report_by_date: PropTypes.string,
    has_dependents: PropTypes.string,
    new_duty_location: PropTypes.shape({}),
    grade: PropTypes.string,
    origin_duty_location: DutyLocationShape,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
};

export default OrdersInfoForm;
