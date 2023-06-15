import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Radio, FormGroup, Label, Tag, Textarea, Fieldset } from '@trussworks/react-uswds';
import moment from 'moment';

import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ReviewSITExtensionModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownInput, DatePickerInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { datePickerFormat, formatDateForDatePicker, utcDateFormat } from 'shared/dates';
import { SitStatusShape, LOCATION_TYPES } from 'types/sitStatusShape';
import { ShipmentShape } from 'types';

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

const SITHistoryItemHeader = ({ title, value }) => {
  let action = '';

  if (title.includes('approved')) {
    action = 'Approved';
  }

  if (title.includes('authorized')) {
    action = 'Authorized';
  }

  return (
    <div className={styles.sitHistoryItemHeader}>
      {title}
      <span className={styles.hintText}>
        {action} + Requested = {value}
      </span>
    </div>
  );
};

const SitStatusTables = ({ sitStatus, sitExtension, shipment }) => {
  const { totalSITDaysUsed, daysInSIT, location } = sitStatus;
  const sitEntryDate = moment(sitStatus.sitEntryDate, utcDateFormat);
  const daysInPreviousSIT = totalSITDaysUsed - daysInSIT;

  /**
   * @function
   * @description This function calculates the total remaining days. If the SIT entry date has started, a day is subtracted from the days remaining to account for the current day.
   * @returns {number|string} - The number of days remaming or 'expired' if there are no more days left.
   */
  const totalDaysRemaining = () => {
    const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
    // SIT not started
    if (!sitStatus && daysRemaining > 0) {
      return daysRemaining;
    }
    // SIT in has started
    if (sitStatus && daysRemaining > 0) {
      return daysRemaining - 1;
    }
    return 'Expired';
  };

  const approvedAndRequestedDaysCombined = shipment.sitDaysAllowance + sitExtension.requestedDays;
  const approvedAndRequestedDatesCombined = formatDateForDatePicker(
    moment()
      .add(sitExtension.requestedDays, 'days')
      .add(Number.isInteger(totalDaysRemaining()) ? totalDaysRemaining() - 1 : 0, 'days'),
  );

  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];
  // Currently active SIT
  const currentLocation = location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

  const currentDaysInSit = <p>{daysInSIT}</p>;
  const currentDateEnteredSit = <p>{formatDateForDatePicker(sitEntryDate)}</p>;

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
      moment.duration(moment(endDate, datePickerFormat).diff(sitEntryDate)).asDays() + daysInPreviousSIT,
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
    const calculatedSitEndDate = formatDateForDatePicker(sitEntryDate.add(daysApproved - daysInPreviousSIT, 'days'));
    endDateHelper.setTouched(true);
    endDateHelper.setValue(calculatedSitEndDate);
  };

  return (
    <>
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT)</p>
        <Tag>Additional Days Requested</Tag>
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        {/* Sit Total days table */}
        <DataTable
          custClass={styles.totalDaysTable}
          columnHeaders={[
            <SITHistoryItemHeader title="Total days of SIT approved" value={approvedAndRequestedDaysCombined} />,
            'Total days used',
            'Total days remaining',
          ]}
          dataRow={[
            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
            sitStatus.totalSITDaysUsed,
            totalDaysRemaining(),
          ]}
        />
      </div>
      <div className={styles.tableContainer}>
        {/* Sit Start and End table */}
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[
            `SIT start date`,
            <SITHistoryItemHeader title="SIT authorized end date" value={approvedAndRequestedDatesCombined} />,
          ]}
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

/**
 * @description This component contains a form that can be viewed from the SIT
 * Display on the MTO page when the Prime submits a SIT Extension for review of
 * the TOO.
 */
const ReviewSITExtensionsModal = ({ onClose, onSubmit, sitExtension, shipment, sitStatus }) => {
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
    acceptExtension: '',
    daysApproved: String(shipment.sitDaysAllowance),
    requestReason: sitExtension.requestReason,
    officeRemarks: '',
    sitEndDate: formatDateForDatePicker(moment(sitStatus.sitAllowanceEndDate, utcDateFormat)),
  };
  const minimumDaysAllowed = shipment.sitDaysAllowance + 1;
  const sitEntryDate = moment(sitStatus.sitEntryDate, utcDateFormat);
  const reviewSITExtensionSchema = Yup.object().shape({
    acceptExtension: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    requestReason: Yup.string().required('Required'),
    officeRemarks: Yup.string().nullable(),
    daysApproved: Yup.number().when('acceptExtension', {
      is: 'yes',
      then: Yup.number()
        .min(minimumDaysAllowed, `Total days of SIT approved must be ${minimumDaysAllowed} or more.`)
        .required('Required'),
    }),
    sitEndDate: Yup.date().min(
      formatDateForDatePicker(sitEntryDate.add(1, 'days')),
      'The end date must occur after the start date. Please select a new date.',
    ),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ReviewSITExtensionModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Review additional days requested</h2>
          </ModalTitle>
          <div>
            <Formik
              validationSchema={reviewSITExtensionSchema}
              onSubmit={(e) => onSubmit(sitExtension.id, e)}
              initialValues={initialValues}
            >
              {({ isValid, values, setValues }) => {
                const handleNoSelection = (e) => {
                  if (e.target.value === 'no') {
                    setValues({
                      ...values,
                      acceptExtension: 'no',
                    });
                  }
                };
                return (
                  <Form>
                    <DataTableWrapper
                      className={classnames('maxw-tablet', styles.sitDisplayForm)}
                      testID="sitExtensions"
                    >
                      <SitStatusTables sitStatus={sitStatus} sitExtension={sitExtension} shipment={shipment} />
                    </DataTableWrapper>
                    <div className={styles.ModalPanel}>
                      <dl className={styles.SITSummary}>
                        <div>
                          <dt>Additional days requested:</dt>
                          <dd>{sitExtension.requestedDays}</dd>
                        </div>
                        <div>
                          <dt>Reason:</dt>
                          <dd>{sitExtensionReasons[sitExtension.requestReason]}</dd>
                        </div>
                        <div>
                          <dt>Contractor remarks:</dt>
                          <dd>{sitExtension.contractorRemarks}</dd>
                        </div>
                      </dl>
                      <FormGroup>
                        <Fieldset legend="Accept request for extension?">
                          <Field
                            as={Radio}
                            label="Yes"
                            id="acceptExtension"
                            name="acceptExtension"
                            value="yes"
                            title="Yes, accept extension"
                            type="radio"
                          />
                          <Field
                            as={Radio}
                            label="No"
                            id="denyExtension"
                            name="acceptExtension"
                            value="no"
                            title="No, deny extension"
                            type="radio"
                            onChange={handleNoSelection}
                          />
                        </Fieldset>
                      </FormGroup>
                      {values.acceptExtension === 'yes' && (
                        <div className={styles.reasonDropdown}>
                          <DropdownInput
                            label="Reason for edit"
                            name="requestReason"
                            options={dropdownInputOptions(sitExtensionReasons)}
                          />
                        </div>
                      )}
                      <Label htmlFor="officeRemarks">Office remarks</Label>
                      <Field
                        as={Textarea}
                        data-testid="officeRemarks"
                        label="No"
                        name="officeRemarks"
                        id="officeRemarks"
                      />
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
                    </div>
                  </Form>
                );
              }}
            </Formik>
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ReviewSITExtensionsModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  sitExtension: SITExtensionShape.isRequired,
  sitStatus: SitStatusShape.isRequired,
  shipment: ShipmentShape.isRequired,
};
export default ReviewSITExtensionsModal;
