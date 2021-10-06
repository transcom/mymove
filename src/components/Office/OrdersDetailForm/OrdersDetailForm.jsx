import React from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { CheckboxField, DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import TextField from 'components/form/fields/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField';
import { DropdownArrayOf } from 'types/form';

const OrdersDetailForm = ({
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  tacWarning,
  validateTac,
  showDepartmentIndicator,
  showOrdersNumber,
  showOrdersTypeDetail,
  showTac,
  showSac,
  showOrdersAcknowledgement,
}) => {
  return (
    <div className={styles.OrdersDetailForm}>
      <DutyStationInput name="originDutyStation" label="Current duty station" displayAddress={false} />
      <DutyStationInput name="newDutyStation" label="New duty station" displayAddress={false} />
      <DatePickerInput name="issueDate" label="Date issued" />
      <DatePickerInput name="reportByDate" label="Report by date" />
      {showDepartmentIndicator && (
        <DropdownInput name="departmentIndicator" label="Department indicator" options={deptIndicatorOptions} />
      )}
      {showOrdersNumber && <TextField name="ordersNumber" label="Orders number" id="ordersNumberInput" />}
      <DropdownInput name="ordersType" label="Orders type" options={ordersTypeOptions} />
      {showOrdersTypeDetail && (
        <DropdownInput name="ordersTypeDetail" label="Orders type detail" options={ordersTypeDetailOptions} />
      )}
      {showTac && (
        <MaskedTextField name="tac" label="TAC" id="tacInput" mask="****" warning={tacWarning} validate={validateTac} />
      )}
      {showSac && <TextField name="sac" label="SAC" id="sacInput" />}
      {showOrdersAcknowledgement && (
        <div className={styles.wrappedCheckbox}>
          <CheckboxField
            id="ordersAcknowledgementInput"
            name="ordersAcknowledgement"
            label="I have read the new orders"
          />
        </div>
      )}
    </div>
  );
};

OrdersDetailForm.propTypes = {
  ordersTypeOptions: DropdownArrayOf.isRequired,
  deptIndicatorOptions: DropdownArrayOf,
  ordersTypeDetailOptions: DropdownArrayOf,
  tacWarning: string,
  validateTac: func,
  showDepartmentIndicator: bool,
  showOrdersNumber: bool,
  showOrdersTypeDetail: bool,
  showTac: bool,
  showSac: bool,
  showOrdersAcknowledgement: bool,
};

OrdersDetailForm.defaultProps = {
  tacWarning: '',
  deptIndicatorOptions: null,
  ordersTypeDetailOptions: null,
  validateTac: null,
  showDepartmentIndicator: true,
  showOrdersNumber: true,
  showOrdersTypeDetail: true,
  showTac: true,
  showSac: true,
  showOrdersAcknowledgement: false,
};

export default OrdersDetailForm;
