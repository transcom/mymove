import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';

import styles from './SubmitSITExtensionModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { DropdownInput, DatePickerInput } from 'components/form/fields';
import { dropdownInputOptions, formatDate } from 'utils/formatters';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

const SitDaysAllowanceForm = ({ onChange }) => (
  <MaskedTextField
    data-testid="daysApproved"
    defaultValue="1"
    id="daysApproved"
    name="daysApproved"
    mask={Number}
    lazy={false}
    scale={0}
    signed={false} // no negative numbers
    inputClassName={styles.approvedDaysInput}
    errorClassName={styles.errors}
    onChange={onChange}
  />
);

const SitEndDateForm = ({ onChange }) => (
  <DatePickerInput name="sitEndDate" label="" id="sitEndDate" onChange={onChange} />
);

const SitStatusTables = ({ sitStatus, shipment }) => {
  const { totalSITDaysUsed, calculatedTotalDaysInSIT } = sitStatus;
  const { daysInSIT } = sitStatus.currentSIT;
  const sitDepartureDate =
    formatDate(sitStatus.currentSIT?.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE;
  const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
  const daysInPreviousSIT = totalSITDaysUsed - daysInSIT;

  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];
  // Currently active SIT
  const currentLocation = sitStatus.currentSIT.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

  const currentDaysInSit = <p>{daysInSIT}</p>;
  const currentDateEnteredSit = <p>{formatDateForDatePicker(sitEntryDate)}</p>;
  const totalDaysRemaining = () => {
    const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
    if (daysRemaining > 0) {
      return daysRemaining;
    }
    return 'Expired';
  };

  /**
   * @function
   * @description This function is used to change the values of the Total Days
   * of SIT approved input when the End Date datepicker is modified. This is
   * being triggered on the `onChange` event for the SitEndDateForm component.
   * @param {Date} endDate A Moment.input representing the last day approved in the form.
   * @see handleDaysAllowanceChange
   * @see SitEndDateForm component
   */
  const handleSitEndDateChange = (endDate) => {
    // Calculate total allowance
    // Set dates to same time zone and strip of time information to calculate integer
    // days between them
    const endDay = moment(endDate).utcOffset(sitEntryDate.utcOffset(), true).startOf('day');
    const startDay = sitEntryDate.startOf('day');
    const sitDurationDays = moment.duration(endDay.diff(startDay)).asDays();
    const calculatedSitDaysAllowance = sitDurationDays + daysInPreviousSIT;

    // Update form values
    endDateHelper.setValue(endDate);
    sitAllowanceHelper.setValue(String(calculatedSitDaysAllowance));
  };

  /**
   * @function
   * @description This function is used to change the values of the End Date
   * datepicker when the Days Approved text input is modified. This is being
   * triggered on the `onChange` event for the SitDaysAllowanceForm component.
   * @param {number} daysApproved A number representing the number of days
   * approved in the form.
   * @see handleSitEndDateChange
   * @see SitDaysAllowanceForm component
   */
  const handleDaysAllowanceChange = (daysApproved) => {
    // Sit days allowance
    sitAllowanceHelper.setValue(daysApproved);
    // // // Sit End date
    const calculatedSitEndDate = formatDateForDatePicker(
      sitEntryDate.add(daysApproved - (calculatedTotalDaysInSIT - daysInSIT), 'days'),
    );
    endDateHelper.setTouched(true);
    endDateHelper.setValue(calculatedSitEndDate);
  };

  return (
    <>
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT)</p>
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        {/* Sit Total days table */}
        <DataTable
          custClass={styles.totalDaysTable}
          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
          dataRow={[
            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
            sitStatus.totalSITDaysUsed,
            totalDaysRemaining(),
          ]}
        />
      </div>
      <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
        {/* Sit Start and End table */}
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[`SIT start date`, 'SIT authorized end date', 'Calculated total SIT days']}
          dataRow={[
            currentDateEnteredSit,
            <SitEndDateForm
              onChange={(value) => {
                handleSitEndDateChange(value);
              }}
            />,
            sitStatus.calculatedTotalDaysInSIT,
          ]}
          custClass={styles.currentLocation}
        />
      </div>
      <div className={styles.tableContainer}>
        {/* Total days at current location */}
        <DataTable
          testID="currentSITDateData"
          columnHeaders={[`Total days in ${currentLocation}`, `SIT departure date`]}
          dataRow={[currentDaysInSit, sitDepartureDate]}
        />
      </div>
    </>
  );
};

const SubmitSITExtensionModal = ({ shipment, sitStatus, onClose, onSubmit }) => {
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
    sitEndDate: formatDateForDatePicker(moment(sitStatus.currentSIT.sitAuthorizedEndDate, swaggerDateFormat)),
  };
  const minimumDaysAllowed = sitStatus.calculatedTotalDaysInSIT - sitStatus.currentSIT.daysInSIT + 1;
  const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
  const reviewSITExtensionSchema = Yup.object().shape({
    requestReason: Yup.string().required('Required'),
    officeRemarks: Yup.string().nullable(),
    daysApproved: Yup.number()
      .min(minimumDaysAllowed, `Total days of SIT approved must be ${minimumDaysAllowed} or more.`)
      .required('Required'),
    sitEndDate: Yup.date().min(
      formatDateForDatePicker(sitEntryDate.add(1, 'days')),
      'The end date must occur after the start date. Please select a new date.',
    ),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.SubmitSITExtensionModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit SIT authorization</h2>
          </ModalTitle>
          <Formik
            validationSchema={reviewSITExtensionSchema}
            onSubmit={(e) => onSubmit(e)}
            initialValues={initialValues}
          >
            {({ isValid }) => {
              return (
                <Form>
                  <DataTableWrapper className={classnames('maxw-tablet', styles.sitDisplayForm)} testID="sitExtensions">
                    <SitStatusTables sitStatus={sitStatus} shipment={shipment} />
                  </DataTableWrapper>
                  <div className={styles.reasonDropdown}>
                    <DropdownInput
                      label="Reason for edit"
                      name="requestReason"
                      data-testid="reasonDropdown"
                      options={dropdownInputOptions(sitExtensionReasons)}
                    />
                  </div>
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

SubmitSITExtensionModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default SubmitSITExtensionModal;
