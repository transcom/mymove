import React from 'react';
import { string } from 'prop-types';

import { TextInput, DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import styles from 'components/Office/OrdersDetailForm/OrdersDetailForm.module.scss';

const OrdersDetailForm = ({ deptIndicatorOptions, ordersTypeOptions, ordersTypeDetailOptions, tacWarning }) => {
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
      <TextInput name="tac" label="TAC" id="tacInput" warning={tacWarning} />
      <TextInput name="sac" label="SAC" id="sacInput" />
    </div>
  );
};

OrdersDetailForm.propTypes = {
  deptIndicatorOptions: DropdownArrayOf.isRequired,
  ordersTypeOptions: DropdownArrayOf.isRequired,
  ordersTypeDetailOptions: DropdownArrayOf.isRequired,
  tacWarning: string,
};

OrdersDetailForm.defaultProps = {
  tacWarning: '',
};

export default OrdersDetailForm;
