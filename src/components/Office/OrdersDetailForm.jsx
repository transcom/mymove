import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import classNames from 'classnames/bind';
import * as Yup from 'yup';
import styles from 'components/Office/OrdersDetailForm.module.scss';
import { Form } from '../form/Form';
import { TextInput, DropdownInput, DatePickerInput, DutyStationInput } from '../form/fields';

const cx = classNames.bind(styles);

export const OrdersDetailForm = ({
  initialValues,
  onSubmit,
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
}) => {
  return (
    <Formik
      onSubmit={(values) => {
        onSubmit(values);
      }}
      initialValues={initialValues}
      validationSchema={Yup.object({
        currentDutyStation: Yup.object().required('Required'),
        newDutyStation: Yup.object().required('Required'),
        dateIssued: Yup.date().required('Required'),
        reportByDate: Yup.date().required('Required'),
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
        <TextInput name="ordersNumber" label="Orders number" />
        <DropdownInput name="ordersType" label="Orders type" options={ordersTypeOptions} />
        <DropdownInput name="ordersTypeDetail" label="Orders type detail" options={ordersTypeDetailOptions} />
        <TextInput name="tac" label="TAC" />
        <TextInput name="sac" label="SAC" />
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
  onSubmit: PropTypes.func,
  deptIndicatorOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)),
  ordersTypeOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)),
  ordersTypeDetailOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)),
  initialValues: PropTypes.shape({
    currentDutyStation: PropTypes.shape({
      address: PropTypes.shape({
        city: PropTypes.string,
        id: PropTypes.string,
        postal_code: PropTypes.string,
        state: PropTypes.string,
        street_address_1: PropTypes.string,
      }),
      address_id: PropTypes.string,
      affiliation: PropTypes.string,
      created_at: PropTypes.string,
      id: PropTypes.string,
      name: PropTypes.string,
      updated_at: PropTypes.string,
    }),
    newDutyStation: PropTypes.shape({
      address: PropTypes.shape({
        city: PropTypes.string,
        id: PropTypes.string,
        postal_code: PropTypes.string,
        state: PropTypes.string,
        street_address_1: PropTypes.string,
      }),
      address_id: PropTypes.string,
      affiliation: PropTypes.string,
      created_at: PropTypes.string,
      id: PropTypes.string,
      name: PropTypes.string,
      updated_at: PropTypes.string,
    }),
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

export default OrdersDetailForm;
