import React from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';

import OrdersDetailForm from 'components/Office/OrdersDetailForm/OrdersDetailForm';
import { dropdownInputOptions } from 'utils/formatters';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS } from 'constants/orders';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';

const originDutyLocation = {
  address: {
    city: 'Dover AFB',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postalCode: '19902',
    state: 'DE',
    streetAddress1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Dover AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};
const newDutyLocation = {
  address: {
    city: 'Scott Air Force Base',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postalCode: '62225',
    state: 'IL',
    streetAddress1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Scott AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};

const deptIndicatorOptions = dropdownInputOptions(DEPARTMENT_INDICATOR_OPTIONS);
const ordersTypeOptions = dropdownInputOptions(ORDERS_TYPE_OPTIONS);
const ordersTypeDetailOptions = dropdownInputOptions(ORDERS_TYPE_DETAILS_OPTIONS);

export default {
  title: 'Office Components/OrdersDetailForm',
  component: OrdersDetailForm,
  decorators: [
    (Story) => (
      <div className="officeApp" style={{ padding: `20px`, background: `#f0f0f0` }}>
        <Story />
      </div>
    ),
  ],
  argTypes: {
    showDepartmentIndicator: { defaultValue: false },
    showOrdersNumber: { defaultValue: false },
    showOrdersTypeDetail: { defaultValue: false },
    showHHGTac: { defaultValue: false },
    showHHGSac: { defaultValue: false },
    showNTSTac: { defaultValue: false },
    showNTSSac: { defaultValue: false },

    showOrdersAcknowledgement: { defaultValue: false },
  },
};

export const EmptyValues = () => (
  <div style={{ width: '400px' }}>
    <Formik>
      <form>
        <OrdersDetailForm
          deptIndicatorOptions={deptIndicatorOptions}
          ordersTypeOptions={ordersTypeOptions}
          ordersTypeDetailOptions={ordersTypeDetailOptions}
          ordersType={ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION}
          setFieldValue={Formik.setFieldValue}
        />
      </form>
    </Formik>
  </div>
);

export const InitialValues = () => {
  return (
    <div style={{ width: '400px' }}>
      <Formik
        initialValues={{
          originDutyLocation,
          newDutyLocation,
          issueDate: '2020-03-08',
          reportByDate: '2020-04-01',
          departmentIndicator: 'NAVY_AND_MARINES',
          ordersNumber: '999999999',
          ordersType: 'PERMANENT_CHANGE_OF_STATION',
          ordersTypeDetail: 'HHG_PERMITTED',
          tac: 'Tac',
          sac: 'Sac',
          ntsTac: 'Tac',
          ntsSac: 'Sac',
          ordersAcknowledgement: true,
        }}
        validationSchema={Yup.object({
          originDutyLocation: Yup.object().defined('Required'),
          newDutyLocation: Yup.object().required('Required'),
          issueDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
          reportByDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
          departmentIndicator: Yup.string().required('Required'),
          ordersNumber: Yup.string().required('Required'),
          ordersType: Yup.string().required('Required'),
          ordersTypeDetail: Yup.string().required('Required'),
          tac: Yup.string().required('Required'),
          sac: Yup.string().required('Required'),
          ntsTac: Yup.string().required('Required'),
          ntsSac: Yup.string().required('Required'),
        })}
      >
        {(formik) => {
          return (
            <form>
              <OrdersDetailForm
                deptIndicatorOptions={deptIndicatorOptions}
                ordersTypeOptions={ordersTypeOptions}
                ordersTypeDetailOptions={ordersTypeDetailOptions}
                showOrdersAcknowledgement
                setFieldValue={formik.setFieldValue}
              />
            </form>
          );
        }}
      </Formik>
    </div>
  );
};

export const FieldsHidden = (args) => {
  return (
    <div style={{ width: '400px' }}>
      <Formik
        initialValues={{
          originDutyLocation,
          newDutyLocation,
          issueDate: '2020-03-08',
          reportByDate: '2020-04-01',
          departmentIndicator: 'NAVY_AND_MARINES',
          ordersNumber: '999999999',
          ordersType: 'PERMANENT_CHANGE_OF_STATION',
          ordersTypeDetail: 'HHG_PERMITTED',
          tac: 'Tac',
          sac: 'Sac',
          ntsTac: 'Tac',
          ntsSac: 'Sac',
        }}
        validationSchema={Yup.object({
          originDutyLocation: Yup.object().defined('Required'),
          newDutyLocation: Yup.object().required('Required'),
          issueDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
          reportByDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY').required('Required'),
          departmentIndicator: Yup.string().required('Required'),
          ordersNumber: Yup.string().required('Required'),
          ordersType: Yup.string().required('Required'),
          ordersTypeDetail: Yup.string().required('Required'),
          tac: Yup.string().required('Required'),
          sac: Yup.string().required('Required'),
          ntsTac: Yup.string().required('Required'),
          ntsSac: Yup.string().required('Required'),
        })}
      >
        {(formik) => {
          return (
            <form>
              <OrdersDetailForm
                deptIndicatorOptions={deptIndicatorOptions}
                ordersTypeOptions={ordersTypeOptions}
                ordersTypeDetailOptions={ordersTypeDetailOptions}
                setFieldValue={formik.setFieldValue}
                {...args}
              />
            </form>
          );
        }}
      </Formik>
    </div>
  );
};

export const Retiree = () => (
  <div style={{ width: '400px' }}>
    <Formik
      initialValues={{
        originDutyLocation,
        newDutyLocation,
        issueDate: '2020-03-08',
        reportByDate: '2020-04-01',
        departmentIndicator: 'NAVY_AND_MARINES',
        ordersNumber: '999999999',
        ordersType: 'RETIREMENT',
        ordersTypeDetail: 'HHG_PERMITTED',
        tac: 'Tac',
        sac: 'Sac',
        ntsTac: 'Tac',
        ntsSac: 'Sac',
        ordersAcknowledgement: true,
      }}
    >
      {(formik) => {
        return (
          <form>
            <OrdersDetailForm
              deptIndicatorOptions={deptIndicatorOptions}
              ordersTypeOptions={ordersTypeOptions}
              ordersTypeDetailOptions={ordersTypeDetailOptions}
              ordersType="RETIREMENT"
              setFieldValue={formik.setFieldValue}
            />
          </form>
        );
      }}
    </Formik>
  </div>
);

export const Separatee = () => (
  <div style={{ width: '400px' }}>
    <Formik
      initialValues={{
        originDutyLocation,
        newDutyLocation,
        issueDate: '2020-03-08',
        reportByDate: '2020-04-01',
        departmentIndicator: 'NAVY_AND_MARINES',
        ordersNumber: '999999999',
        ordersType: 'SEPARATION',
        ordersTypeDetail: 'HHG_PERMITTED',
        tac: 'Tac',
        sac: 'Sac',
        ntsTac: 'Tac',
        ntsSac: 'Sac',
        ordersAcknowledgement: true,
      }}
    >
      {(formik) => {
        return (
          <form>
            <OrdersDetailForm
              deptIndicatorOptions={deptIndicatorOptions}
              ordersTypeOptions={ordersTypeOptions}
              ordersTypeDetailOptions={ordersTypeDetailOptions}
              ordersType="SEPARATION"
              setFieldValue={formik.setFieldValue}
            />
          </form>
        );
      }}
    </Formik>
  </div>
);
