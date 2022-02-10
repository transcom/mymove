import React from 'react';

import AccountingCodesModal from './AccountingCodesModal';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/AccountingCodesModal',
  component: AccountingCodesModal,
};

export const minimal = () => (
  <div className="officeApp">
    <AccountingCodesModal isOpen shipmentType={SHIPMENT_OPTIONS.NTS} />
  </div>
);

export const withCodes = () => (
  <div className="officeApp">
    <AccountingCodesModal
      isOpen
      shipmentType={SHIPMENT_OPTIONS.NTSR}
      TACs={{ HHG: '1234', NTS: '2345' }}
      SACs={{ HHG: 'CF13', NTS: 'E8A1' }}
    />
  </div>
);

export const withCodesAndValues = () => (
  <div className="officeApp">
    <AccountingCodesModal
      isOpen
      shipmentType={SHIPMENT_OPTIONS.HHG}
      TACs={{ HHG: '1234', NTS: '2345' }}
      SACs={{ HHG: 'CF13', NTS: 'E8A1' }}
      tacType="NTS"
      sacType="HHG"
    />
  </div>
);
