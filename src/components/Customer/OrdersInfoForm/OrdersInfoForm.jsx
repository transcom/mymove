import React from 'react';
import { Field } from 'formik';
import { Radio, FormGroup, Label } from '@trussworks/react-uswds';

import { DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { Form } from 'components/form/Form';

const OrdersInfoForm = () => {
  const ordersTypeOptions = [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ];

  return (
    <Form>
      <Field as={DropdownInput} label="Orders type" name="orders_type" options={ordersTypeOptions} />
      <Field
        as={DatePickerInput}
        name="issue_date"
        label="Orders date"
        renderInput={(input) => (
          <>
            {input}
            <Hint>
              <p>Date your orders were issued.</p>
            </Hint>
          </>
        )}
      />
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
      <DutyStationInput name="new_duty_station" label="New duty station" displayAddress={false} />
    </Form>
  );
};

export default OrdersInfoForm;
