import React from 'react';
import { object, text } from '@storybook/addon-knobs';

import OrdersList from './OrdersList';

export default {
  title: 'Office Components/OrdersList',
  component: OrdersList,
};

export const Basic = () => (
  <div className="officeApp">
    <OrdersList
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: text('ordersInfo.departmentIndicator', 'NAVY_AND_MARINES'),
        ordersNumber: text('ordersInfo.ordersNumber', '999999999'),
        ordersType: text('ordersInfo.ordersType', 'PERMANENT_CHANGE_OF_STATION'),
        ordersTypeDetail: text('ordersInfo.ordersTypeDetail', 'HHG_PERMITTED'),
        tacMDC: text('ordersInfo.tacMDC', '9999'),
        sacSDN: text('ordersInfo.sacSDN', '999 999999 999'),
        NTSsac: text('ordersInfo.NTSsac', '999 999999 999'),
        NTStac: text('ordersInfo.NTStac', '9999'),
      }}
    />
  </div>
);

export const AsServiceCounselor = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings={false}
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: '',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);

export const AsServiceCounselorProcessingRetirement = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings={false}
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'RETIREMENT',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);

export const AsServiceCounselorProcessingSeparation = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings={false}
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'SEPARATION',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);

export const AsTOO = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: '',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);

export const AsTOOProcessingRetirement = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'RETIREMENT',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);

export const AsTOOProcessingSeparation = () => (
  <div className="officeApp">
    <OrdersList
      showMissingWarnings
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'SEPARATION',
        ordersTypeDetail: '',
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
      }}
    />
  </div>
);
