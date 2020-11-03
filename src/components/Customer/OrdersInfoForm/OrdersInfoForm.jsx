import React from 'react';
import { Formik, Field } from 'formik';
import { Radio, FormGroup, Label } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import { DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { Form } from 'components/form/Form';

/**
 * TODO:
 * - validation:
 *  - duty station != current duty station
 * - fetch latest orders on mount
 * - handle submit (create or update)
 * - initialize values based on current orders
 * - orders types feature flag
 * - display server error
 * -
 */

const ordersInfoSchema = Yup.object().shape({
  orders_type: Yup.mixed().oneOf(['PERMANENT_CHANGE_OF_STATION', 'RETIREMENT', 'SEPARATION']).required('Required'),
  issue_date: Yup.date().required('Required'),
  report_by_date: Yup.date().required('Required'),
  has_dependents: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
});

const OrdersInfoForm = () => {
  const initialValues = {
    orders_type: '', // required
    issue_date: '', // required
    report_by_date: '', // required
    has_dependents: '', // required
    new_duty_station: {},
  };

  const ordersTypeOptions = [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ];

  return (
    <Formik initialValues={initialValues} validateOnMount validationSchema={ordersInfoSchema}>
      {({ values, errors, isValid }) => {
        return (
          <>
            <h1>Tell us about your move orders</h1>
            <p>{JSON.stringify(values)}</p>
            <p>is valid? {JSON.stringify(isValid)}</p>
            <p>{JSON.stringify(errors)}</p>
            <Form>
              <Field as={DropdownInput} label="Orders type" name="orders_type" options={ordersTypeOptions} />
              <FormGroup>
                <Field as={DatePickerInput} name="issue_date" label="Orders date" />
                <Hint>
                  <p>Date your orders were issued.</p>
                </Hint>
              </FormGroup>
              <Field as={DatePickerInput} name="report_by_date" label="Report-by date" />
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
                  />
                  <Field
                    as={Radio}
                    label="No"
                    id="hasDependentsNo"
                    name="has_dependents"
                    value="no"
                    title="No, dependents are not included in my orders"
                  />
                </div>
              </FormGroup>
              <Field as={DutyStationInput} name="new_duty_station" label="New duty station" displayAddress={false} />
            </Form>
          </>
        );
      }}
    </Formik>
  );
};

export default OrdersInfoForm;
