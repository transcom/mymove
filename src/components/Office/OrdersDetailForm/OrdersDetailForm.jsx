import React, { useState } from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { formatLabelReportByDate } from 'utils/formatters';
import { CheckboxField, DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
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
  ordersType,
  setFieldValue,
}) => {
  const [formOrdersType, setFormOrdersType] = useState(ordersType);
  const reportDateRowLabel = formatLabelReportByDate(formOrdersType);
  // The text is different if the customer is retiring or separating.
  const newDutyLocationLabel = ['RETIREMENT', 'SEPARATION'].includes(formOrdersType)
    ? 'HOR, HOS or PLEAD'
    : 'New duty location';

  return (
    <div className={styles.OrdersDetailForm}>
      <DutyLocationInput name="originDutyLocation" label="Current duty location" displayAddress={false} />
      <DutyLocationInput name="newDutyLocation" label={newDutyLocationLabel} displayAddress={false} />
      <DatePickerInput name="issueDate" label="Date issued" />
      <DatePickerInput name="reportByDate" label={reportDateRowLabel} />
      {showDepartmentIndicator && (
        <DropdownInput name="departmentIndicator" label="Department indicator" options={deptIndicatorOptions} />
      )}
      {showOrdersNumber && <TextField name="ordersNumber" label="Orders number" id="ordersNumberInput" />}
      <DropdownInput
        name="ordersType"
        label="Orders type"
        options={ordersTypeOptions}
        onChange={(e) => {
          setFormOrdersType(e.target.value);
          setFieldValue('ordersType', e.target.value);
        }}
      />
      {showOrdersTypeDetail && (
        <DropdownInput name="ordersTypeDetail" label="Orders type detail" options={ordersTypeDetailOptions} />
      )}

      {showHHGTac && showHHGSac && <h3>HHG accounting codes</h3>}
      {showHHGTac && (
        <MaskedTextField
          name="tac"
          label="TAC"
          id="hhgTacInput"
          mask="****"
          inputTestId="hhgTacInput"
          warning={hhgTacWarning}
          validate={validateHHGTac}
        />
      )}
      {showHHGSac && <TextField name="sac" label="SAC" id="hhgSacInput" data-testid="hhgSacInput" optional />}

      {showNTSTac && showNTSSac && <h3>NTS accounting codes</h3>}
      {showNTSTac && (
        <MaskedTextField
          name="ntsTac"
          label="TAC"
          id="ntsTacInput"
          mask="****"
          inputTestId="ntsTacInput"
          warning={ntsTacWarning}
          validate={validateNTSTac}
          optional
        />
      )}
      {showNTSSac && <TextField name="ntsSac" label="SAC" id="ntsSacInput" data-testid="ntsSacInput" optional />}

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
  ordersType: string.isRequired,
  setFieldValue: func.isRequired,
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
