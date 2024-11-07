import React from 'react';
import { array, object, text } from '@storybook/addon-knobs';

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
        ordersDocuments: array('ordersInfo.ordersDocuments', [
          {
            'c0a22a98-a806-47a2-ab54-2dac938667b3': {
              bytes: 2202009,
              contentType: 'application/pdf',
              createdAt: '2024-10-23T16:31:21.085Z',
              filename: 'testFile.pdf',
              id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
              status: 'PROCESSING',
              updatedAt: '2024-10-23T16:31:21.085Z',
              uploadType: 'USER',
              url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
            },
          },
        ]),
        tacMDC: text('ordersInfo.tacMDC', '9999'),
        sacSDN: text('ordersInfo.sacSDN', '999 999999 999'),
        NTSsac: text('ordersInfo.NTSsac', '999 999999 999'),
        NTStac: text('ordersInfo.NTStac', '9999'),
        payGrade: text('ordersInfo.payGrade', 'E_5'),
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
        ordersDocuments: array('ordersInfo.ordersDocuments', [
          {
            'c0a22a98-a806-47a2-ab54-2dac938667b3': {
              bytes: 2202009,
              contentType: 'application/pdf',
              createdAt: '2024-10-23T16:31:21.085Z',
              filename: 'testFile.pdf',
              id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
              status: 'PROCESSING',
              updatedAt: '2024-10-23T16:31:21.085Z',
              uploadType: 'USER',
              url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
            },
          },
        ]),
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
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
        ordersDocuments: null,
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
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
        ordersDocuments: null,
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
      }}
    />
  </div>
);

export const AsTOO = () => (
  <div className="officeApp">
    <OrdersList
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: '',
        ordersTypeDetail: '',
        ordersDocuments: array('ordersInfo.ordersDocuments', [
          {
            'c0a22a98-a806-47a2-ab54-2dac938667b3': {
              bytes: 2202009,
              contentType: 'application/pdf',
              createdAt: '2024-10-23T16:31:21.085Z',
              filename: 'testFile.pdf',
              id: 'c0a22a98-a806-47a2-ab54-2dac938667b3',
              status: 'PROCESSING',
              updatedAt: '2024-10-23T16:31:21.085Z',
              uploadType: 'USER',
              url: '/storage/USER/uploads/c0a22a98-a806-47a2-ab54-2dac938667b3?contentType=application%2Fpdf',
            },
          },
        ]),
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
      }}
    />
  </div>
);

export const AsTOOProcessingRetirement = () => (
  <div className="officeApp">
    <OrdersList
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'RETIREMENT',
        ordersTypeDetail: '',
        ordersDocuments: null,
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
      }}
    />
  </div>
);

export const AsTOOProcessingSeparation = () => (
  <div className="officeApp">
    <OrdersList
      ordersInfo={{
        currentDutyLocation: object('ordersInfo.currentDutyLocation', { name: 'JBSA Lackland' }),
        newDutyLocation: object('ordersInfo.newDutyLocation', { name: 'JB Lewis-McChord' }),
        issuedDate: text('ordersInfo.issuedDate', '2020-03-08'),
        reportByDate: text('ordersInfo.reportByDate', '2020-04-01'),
        departmentIndicator: '',
        ordersNumber: '',
        ordersType: 'SEPARATION',
        ordersTypeDetail: '',
        ordersDocuments: null,
        tacMDC: '',
        sacSDN: '',
        NTSsac: '',
        NTStac: '',
        payGrade: text('ordersInfo.payGrade', 'E_5'),
      }}
    />
  </div>
);
