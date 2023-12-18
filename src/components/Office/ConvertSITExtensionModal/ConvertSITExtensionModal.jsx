import React from 'react';
// jeh import classnames from 'classnames';
// jeh import { Formik, Field, useField } from 'formik';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';

import styles from './ConvertSITExtensionModal.module.scss';

// jeh import DataTableWrapper from 'components/DataTableWrapper/index';
// jeh import DataTable from 'components/DataTable/index';
// jeh import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
// jeh import { DropdownInput, DatePickerInput, CheckboxField } from 'components/form/fields';
// jeh import { DatePickerInput, CheckboxField } from 'components/form/fields';
import { CheckboxField } from 'components/form/fields';
// jeh import { dropdownInputOptions } from 'utils/formatters';
// jeh import { sitExtensionReasons } from 'constants/sitExtensions';
// jeh import { LOCATION_TYPES } from 'types/sitStatusShape';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';

// jeh const SitDaysAllowanceForm = ({ onChange }) => (
// jeh   <MaskedTextField
// jeh     data-testid="daysApproved"
// jeh     defaultValue="1"
// jeh     id="daysApproved"
// jeh     name="daysApproved"
// jeh     mask={Number}
// jeh     lazy={false}
// jeh     scale={0}
// jeh     signed={false} // no negative numbers
// jeh     inputClassName={styles.approvedDaysInput}
// jeh     errorClassName={styles.errors}
// jeh     onChange={onChange}
// jeh   />
// jeh );

// jeh const SitEndDateForm = ({ onChange }) => (
// jeh   <DatePickerInput name="sitEndDate" label="" id="sitEndDate" onChange={onChange} />
// jeh );

// jeh const SitStatusTables = ({ sitStatus, shipment }) => {
// jeh  const { totalSITDaysUsed } = sitStatus;
// jeh  const { daysInSIT } = sitStatus.currentSIT;
// jeh  const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
// jeh  const daysInPreviousSIT = totalSITDaysUsed - daysInSIT;
// jeh
// jeh  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
// jeh  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];
// jeh  // Currently active SIT
// jeh  const currentLocation = sitStatus.currentSIT.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';
// jeh
// jeh  const currentDaysInSit = <p>{daysInSIT}</p>;
// jeh  const currentDateEnteredSit = <p>{formatDateForDatePicker(sitEntryDate)}</p>;
// jeh  const totalDaysRemaining = () => {
// jeh    const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
// jeh    if (daysRemaining > 0) {
// jeh      return daysRemaining;
// jeh    }
// jeh    return 'Expired';
// jeh  };
// jeh
// jeh  /**
// jeh   * @function
// jeh   * @description This function is used to change the values of the Total Days
// jeh   * of SIT approved input when the End Date datepicker is modified. This is
// jeh   * being triggered on the `onChange` event for the SitEndDateForm component.
// jeh   * @param {Date} endDate A Moment.input representing the last day approved in the form.
// jeh   * @see handleDaysAllowanceChange
// jeh   * @see SitEndDateForm component
// jeh   */
// jeh  const handleSitEndDateChange = (endDate) => {
// jeh    // Calculate total allowance
// jeh    // Set dates to same time zone and strip of time information to calculate integer
// jeh    // days between them
// jeh    const endDay = moment(endDate).utcOffset(sitEntryDate.utcOffset(), true).startOf('day');
// jeh    const startDay = sitEntryDate.startOf('day');
// jeh    const sitDurationDays = moment.duration(endDay.diff(startDay)).asDays();
// jeh    const calculatedSitDaysAllowance = sitDurationDays + daysInPreviousSIT;
// jeh
// jeh    // Update form values
// jeh    endDateHelper.setValue(endDate);
// jeh    sitAllowanceHelper.setValue(String(calculatedSitDaysAllowance));
// jeh  };
// jeh
// jeh  /**
// jeh   * @function
// jeh   * @description This function is used to change the values of the End Date
// jeh   * datepicker when the Days Approved text input is modified. This is being
// jeh   * triggered on the `onChange` event for the SitDaysAllowanceForm component.
// jeh   * @param {number} daysApproved A number representing the number of days
// jeh   * approved in the form.
// jeh   * @see handleSitEndDateChange
// jeh   * @see SitDaysAllowanceForm component
// jeh   */
// jeh  const handleDaysAllowanceChange = (daysApproved) => {
// jeh    // Sit days allowance
// jeh    sitAllowanceHelper.setValue(daysApproved);
// jeh    // // // Sit End date
// jeh    const calculatedSitEndDate = formatDateForDatePicker(sitEntryDate.add(daysApproved - daysInPreviousSIT, 'days'));
// jeh    endDateHelper.setTouched(true);
// jeh    endDateHelper.setValue(calculatedSitEndDate);
// jeh  };
// jeh
// jeh  return (
// jeh    <>
// jeh      <div className={styles.title}>
// jeh        <p>SIT (STORAGE IN TRANSIT)</p>
// jeh      </div>
// jeh      <div className={styles.tableContainer} data-testid="sitStatusTable">
// jeh        {/* Sit Total days table */}
// jeh        <DataTable
// jeh          custClass={styles.totalDaysTable}
// jeh          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
// jeh          dataRow={[
// jeh            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
// jeh            sitStatus.totalSITDaysUsed,
// jeh            totalDaysRemaining(),
// jeh          ]}
// jeh        />
// jeh      </div>
// jeh      <div className={styles.tableContainer}>
// jeh        {/* Sit Start and End table */}
// jeh        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
// jeh        <DataTable
// jeh          columnHeaders={[`SIT start date`, 'SIT authorized end date']}
// jeh          dataRow={[
// jeh            currentDateEnteredSit,
// jeh            <SitEndDateForm
// jeh              onChange={(value) => {
// jeh                handleSitEndDateChange(value);
// jeh              }}
// jeh            />,
// jeh          ]}
// jeh          custClass={styles.currentLocation}
// jeh        />
// jeh      </div>
// jeh      <div className={styles.tableContainer}>
// jeh        {/* Total days at current location */}
// jeh        <DataTable columnHeaders={[`Total days in ${currentLocation}`]} dataRow={[currentDaysInSit]} />
// jeh      </div>
// jeh    </>
// jeh  );
// jeh };

const ConvertSITExtensionModal = ({ shipment, sitStatus, onClose, onSubmit }) => {
  let sitStartDate = sitStatus?.sitEntryDate;
  if (!sitStartDate) {
    sitStartDate = shipment.mtoServiceItems?.reduce((item, acc) => {
      if (item.sitEntryDate < acc.sitEntryDate) {
        return item;
      }
      return acc;
    }).sitEntryDate;
  }

  const initialValues = {
    requestReason: '',
    officeRemarks: '',
    daysApproved: String(shipment.sitDaysAllowance),
    sitEndDate: formatDateForDatePicker(moment(sitStatus.currentSIT.sitAllowanceEndDate, swaggerDateFormat)),
  };
  // jeh const minimumDaysAllowed = sitStatus.totalSITDaysUsed - sitStatus.currentSIT.daysInSIT + 1;
  // jeh const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
  const convertSITExtensionSchema = Yup.object().shape({
    convertToCustomerExpense: Yup.bool().required(),
    officeRemarks: Yup.string().required(),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ConvertSITExtensionModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Convert SIT To Customer Expense</h2>
          </ModalTitle>
          <Formik
            validationSchema={convertSITExtensionSchema}
            onSubmit={(e) => onSubmit(e)}
            initialValues={initialValues}
          >
            {({ isValid }) => {
              return (
                <Form>
                  <CheckboxField
                    id="convertToCustomerExpense"
                    label="Convert to customer expense"
                    name="convert_to_customer_expense"
                  />
                  <Label htmlFor="officeRemarks">Office remarks</Label>
                  <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
                  <ModalActions>
                    <Button type="submit" disabled={!isValid}>
                      Save
                    </Button>
                    <Button
                      type="button"
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                      outline
                      className={styles.CancelButton}
                    >
                      Cancel
                    </Button>
                  </ModalActions>
                </Form>
              );
            }}
          </Formik>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ConvertSITExtensionModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default ConvertSITExtensionModal;
