import React from 'react';
import { func, string } from 'prop-types';

import { TextInput, DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import { TextMaskedInput } from 'components/form/fields/TextInput';
import { DropdownArrayOf } from 'types/form';
import styles from 'components/Office/OrdersDetailForm/OrdersDetailForm.module.scss';

const OrdersDetailForm = ({
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  tacWarning,
  validateTac,
}) => {
  return (
    <div className={styles.OrdersDetailForm}>
      <DutyStationInput name="originDutyStation" label="Current duty station" displayAddress={false} />
      <DutyStationInput name="newDutyStation" label="New duty station" displayAddress={false} />
      <DatePickerInput name="issueDate" label="Date issued" />
      <DatePickerInput name="reportByDate" label="Report by date" />
      <DropdownInput name="departmentIndicator" label="Department indicator" options={deptIndicatorOptions} />
      <TextInput name="ordersNumber" label="Orders number" id="ordersNumberInput" />
      <DropdownInput name="ordersType" label="Orders type" options={ordersTypeOptions} />
      <DropdownInput name="ordersTypeDetail" label="Orders type detail" options={ordersTypeDetailOptions} />
      <TextMaskedInput name="tac" label="TAC" id="tacInput" mask="****" warning={tacWarning} validate={validateTac} />
      <TextInput name="sac" label="SAC" id="sacInput" />
    </div>
  );
};

OrdersDetailForm.propTypes = {
  deptIndicatorOptions: DropdownArrayOf.isRequired,
  ordersTypeOptions: DropdownArrayOf.isRequired,
  ordersTypeDetailOptions: DropdownArrayOf.isRequired,
  tacWarning: string,
  validateTac: func.isRequired,
};

OrdersDetailForm.defaultProps = {
  tacWarning: '',
};

export default OrdersDetailForm;
