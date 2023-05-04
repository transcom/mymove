import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Radio, FormGroup, Label, Textarea, Fieldset } from '@trussworks/react-uswds';
import moment from 'moment';

import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ReviewSITExtensionModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownInput, DatePickerInput } from 'components/form/fields';
import { dropdownInputOptions, formatDate } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { utcDateFormat } from 'shared/dates';
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
  const { sitEntryDate, totalSITDaysUsed, daysInSIT } = sitStatus;
  const daysInPreviousSIT = totalSITDaysUsed - daysInSIT;

  const approvedAndRequestedDaysCombined = shipment.sitDaysAllowance + sitExtension.requestedDays;
  const approvedAndRequestedDatesCombined = moment(
    moment()
      .add(sitStatus.totalDaysRemaining - 1, 'days')
      .format('DD MMM YYYY'),
  )
    .add(sitExtension.requestedDays, 'days')
    .format('DD MMM YYYY');

  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];
  // Currently active SIT
  const currentLocation = sitStatus.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

  const currentDaysInSit = <p>{sitStatus.totalSITDaysUsed}</p>;
  const currentDateEnteredSit = <p>{formatDate(sitStatus.sitEntryDate, utcDateFormat, 'DD MMM YYYY')}</p>;

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
          columnHeaders={[
            <SITHistoryItemHeader title="Total days of SIT approved" value={approvedAndRequestedDaysCombined} />,
            'Total days used',
            'Total days remaining',
          ]}
          dataRow={[
            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
            sitStatus.totalSITDaysUsed,
            sitStatus.totalDaysRemaining,
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

const ReviewSITExtensionsModal = ({ onClose, onSubmit, sitExtension, shipment, sitStatus }) => {
  const initialValues = {
    acceptExtension: 'yes',
    daysApproved: String(shipment.sitDaysAllowance),
    requestReason: '',
    officeRemarks: '',
    sitEndDate: moment()
      .add(sitStatus.totalDaysRemaining - 1, 'days')
      .format('DD MMM YYYY'),
  };
  const minimumDaysAllowed = sitStatus.totalSITDaysUsed - sitStatus.daysInSIT + 1;
  const reviewSITExtensionSchema = Yup.object().shape({
    acceptExtension: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
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
