import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import classNames from 'classnames/bind';
import * as Yup from 'yup';

import { Form } from '../form/Form';
import { TextInput, DropdownInput, DatePickerInput, DutyStationInput } from '../form/fields';

import { DutyStationShape } from 'types/dutyStation';
import { DropdownArrayOf } from 'types/form';
import styles from 'components/Office/OrdersDetailForm.module.scss';

const cx = classNames.bind(styles);

export const OrdersDetailForm = ({
  initialValues,
  onSubmit,
  onReset,
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
}) => {
  return (
    <Formik
      onSubmit={onSubmit}
      onReset={onReset}
      initialValues={initialValues}
      validationSchema={Yup.object({
        currentDutyStation: Yup.object().required('Required'),
        newDutyStation: Yup.object().required('Required'),
        dateIssued: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
        reportByDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
        departmentIndicator: Yup.string().required('Required'),
        ordersNumber: Yup.string().required('Required'),
        ordersType: Yup.string().required('Required'),
        ordersTypeDetail: Yup.string().required('Required'),
        tac: Yup.string().required('Required'),
        sac: Yup.string().required('Required'),
      })}
    >
      <Form>
        <DutyStationInput name="currentDutyStation" label="Current duty station" />
        <DutyStationInput name="newDutyStation" label="New duty station" />
        <DatePickerInput name="dateIssued" label="Date issued" />
        <DatePickerInput name="reportByDate" label="Report by date" />
        <DropdownInput name="departmentIndicator" label="Department indicator" options={deptIndicatorOptions} />
        <TextInput name="ordersNumber" label="Orders number" id="ordersNumberInput" />
        <DropdownInput name="ordersType" label="Orders type" options={ordersTypeOptions} />
        <DropdownInput name="ordersTypeDetail" label="Orders type detail" options={ordersTypeDetailOptions} />
        <TextInput name="tac" label="TAC" id="tacInput" />
        <TextInput name="sac" label="SAC" id="sacInput" />
        <div className={cx('form-buttons')}>
          <Button type="submit">Submit</Button>
          <Button type="reset" secondary>
            Cancel
          </Button>
        </div>
      </Form>
    </Formik>
  );
};

OrdersDetailForm.propTypes = {
  onReset: PropTypes.func,
  onSubmit: PropTypes.func.isRequired,
  deptIndicatorOptions: DropdownArrayOf.isRequired,
  ordersTypeOptions: DropdownArrayOf.isRequired,
  ordersTypeDetailOptions: DropdownArrayOf.isRequired,
  initialValues: PropTypes.shape({
    currentDutyStation: DutyStationShape,
    newDutyStation: DutyStationShape,
    dateIssued: PropTypes.string,
    reportByDate: PropTypes.string,
    departmentIndicator: PropTypes.string,
    ordersNumber: PropTypes.string,
    ordersType: PropTypes.string,
    ordersTypeDetail: PropTypes.string,
    tac: PropTypes.string,
    sac: PropTypes.string,
  }),
};

OrdersDetailForm.defaultProps = {
  onReset: null,
  initialValues: {},
};

export default OrdersDetailForm;
