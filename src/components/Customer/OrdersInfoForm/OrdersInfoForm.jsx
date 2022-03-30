import React from 'react';
import PropTypes from 'prop-types';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import { Radio, FormGroup, Label, Link as USWDSLink } from '@trussworks/react-uswds';

import styles from './OrdersInfoForm.module.scss';

import { DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { Form } from 'components/form/Form';
import { DropdownArrayOf } from 'types';
import formStyles from 'styles/form.module.scss';
import { DutyLocationShape } from 'types/dutyLocation';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import Callout from 'components/Callout';
import { formatLabelReportByDate } from 'utils/formatters';

const OrdersInfoForm = ({ currentDutyLocation, ordersTypeOptions, initialValues, onSubmit, onBack }) => {
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
    new_duty_location: Yup.object()
      .shape({
        name: Yup.string().notOneOf(
          [currentDutyLocation?.name],
          'You entered the same duty location for your origin and destination. Please change one of them.',
        ),
      })
      .nullable()
      .required('Required'),
  });

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(values.orders_type);

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
                    placeholder="Enter a city or ZIP"
                  />
                </>
              ) : (
                <DutyLocationInput
                  name="new_duty_location"
                  label="New duty location"
                  displayAddress={false}
                  placeholder="Start typing a duty location..."
                />
              )}
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
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onBack: PropTypes.func.isRequired,
  currentDutyLocation: DutyLocationShape.isRequired,
};

export default OrdersInfoForm;
