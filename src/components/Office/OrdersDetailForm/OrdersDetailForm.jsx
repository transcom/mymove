import React from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { CheckboxField, DropdownInput, DatePickerInput, DutyStationInput } from 'components/form/fields';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownArrayOf } from 'types/form';

const OrdersDetailForm = ({
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  hhgTacWarning,
  ntsTacWarning,
  validateHHGTac,
  validateNTSTac,
  showDepartmentIndicator,
  showOrdersNumber,
  showOrdersTypeDetail,
  showHHGTac,
  showHHGSac,
  showNTSTac,
  showNTSSac,
  showOrdersAcknowledgement,
}) => {
  return (
    <div className={styles.OrdersDetailForm}>
      <DutyStationInput name="originDutyStation" label="Current duty location" displayAddress={false} />
      <DutyStationInput name="newDutyStation" label="New duty location" displayAddress={false} />
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

      {showHHGTac && showHHGSac && <h3>HHG accounting codes</h3>}
      {showHHGTac && (
        <MaskedTextField
          name="hhg_tac"
          label="TAC"
          id="hhgTacInput"
          mask="****"
          warning={hhgTacWarning}
          validate={validateHHGTac}
        />
      )}
      {showHHGSac && <TextField name="hhg_sac" label="SAC" id="hhgSacInput" optional />}

      {showNTSTac && showNTSSac && <h3>NTS accounting codes</h3>}
      {showNTSTac && (
        <MaskedTextField
          name="nts_tac"
          label="TAC"
          id="ntsTacInput"
          mask="****"
          warning={ntsTacWarning}
          validate={validateNTSTac}
          optional
        />
      )}
      {showNTSSac && <TextField name="nts_sac" label="SAC" id="ntsSacInput" optional />}

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
  hhgTacWarning: string,
  ntsTacWarning: string,
  validateHHGTac: func,
  validateNTSTac: func,
  showDepartmentIndicator: bool,
  showOrdersNumber: bool,
  showOrdersTypeDetail: bool,
  showHHGTac: bool,
  showHHGSac: bool,
  showNTSTac: bool,
  showNTSSac: bool,
  showOrdersAcknowledgement: bool,
};

OrdersDetailForm.defaultProps = {
  hhgTacWarning: '',
  ntsTacWarning: '',
  deptIndicatorOptions: null,
  ordersTypeDetailOptions: null,
  validateHHGTac: null,
  validateNTSTac: null,
  showDepartmentIndicator: true,
  showOrdersNumber: true,
  showOrdersTypeDetail: true,
  showHHGTac: true,
  showHHGSac: true,
  showNTSTac: true,
  showNTSSac: true,
  showOrdersAcknowledgement: false,
};

export default OrdersDetailForm;
