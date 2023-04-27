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
import { utcDateFormat } from 'shared/dates';
import { LOCATION_TYPES } from 'types/sitStatusShape';

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

const SitStatusTables = ({ sitStatus }) => {
  const { sitEntryDate, totalSITDaysUsed, daysInSIT } = sitStatus;
  const daysInPreviousSIT = totalSITDaysUsed - daysInSIT;

  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];
  // Currently active SIT
  const currentLocation = sitStatus.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

  const currentDaysInSit = <p>{sitStatus.totalSITDaysUsed}</p>;
  const currentDateEnteredSit = <p>{formatDate(sitStatus.sitEntryDate, utcDateFormat, 'DD MMM YYYY')}</p>;
  const totalDaysRemaining = Number(sitStatus.totalDaysRemaining) < 0 ? 'Expired' : sitStatus.totalDaysRemaining;

  /**
   * @function
   * @description This function is used to change the values of the Total Days
   * of SIT approved input when the End Date datepicker is modified. This is
   * being triggered on the `onChange` event for the SitEndDateForm component.
   * @param {moment.input} endDate A Moment.input representing the last day approved in the form.
   * @see handleDaysAllowanceChange
   * @see SitEndDateForm component
   */
  const handleSitEndDateChange = (endDate) => {
    endDateHelper.setValue(endDate);
    // Total days of SIT
    const calculatedSitDaysAllowance = Math.ceil(
      moment.duration(moment(endDate).diff(moment(sitEntryDate))).asDays() + daysInPreviousSIT,
    );
    // Update form value
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
    const calculatedSitEndDate = moment(sitEntryDate)
      .add(daysApproved - daysInPreviousSIT, 'days')
      .format('DD MMM YYYY');
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
            totalDaysRemaining,
          ]}
        />
      </div>
      <div className={styles.tableContainer}>
        {/* Sit Start and End table */}
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[`SIT start date`, 'SIT authorized end date']}
          dataRow={[
            currentDateEnteredSit,
            <SitEndDateForm
              onChange={(value) => {
                handleSitEndDateChange(value);
              }}
            />,
          ]}
          custClass={styles.currentLocation}
        />
      </div>
      <div className={styles.tableContainer}>
        {/* Total days at current location */}
        <DataTable columnHeaders={[`Total days in ${currentLocation}`]} dataRow={[currentDaysInSit]} />
      </div>
    </>
  );
};

const SubmitSITExtensionModal = ({ shipment, sitStatus, onClose, onSubmit }) => {
  const initialValues = {
    requestReason: '',
    officeRemarks: '',
    daysApproved: String(shipment.sitDaysAllowance),
    sitEndDate: moment().add(sitStatus.totalDaysRemaining, 'days').format('DD MMM YYYY'),
  };
  const minimumDaysAllowed = sitStatus.totalSITDaysUsed - sitStatus.daysInSIT + 1;
  const reviewSITExtensionSchema = Yup.object().shape({
    requestReason: Yup.string().required('Required'),
    officeRemarks: Yup.string().nullable(),
    daysApproved: Yup.number()
      .min(minimumDaysAllowed, `Total days of SIT approved must be ${minimumDaysAllowed} or more.`)
      .required('Required'),
    sitEndDate: Yup.date().min(
      moment(sitStatus.sitEntryDate).add(1, 'days').format('DD MMM YYYY'),
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
